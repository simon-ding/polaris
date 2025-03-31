package server

import (
	"os"
	"path/filepath"
	"polaris/db"
	"polaris/engine"
	"polaris/ent"
	"polaris/ent/episode"
	"polaris/ent/media"
	"polaris/log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type searchTvParam struct {
	Query string `form:"query"`
	Page  int    `form:"page"`
}

func (s *Server) SearchTvSeries(c *gin.Context) (interface{}, error) {
	var q searchTvParam
	if err := c.ShouldBindQuery(&q); err != nil {
		return nil, errors.Wrap(err, "bind query")
	}
	log.Infof("search tv series with keyword: %v", q.Query)
	r, err := s.MustTMDB().SearchTvShow(q.Query, "")
	if err != nil {
		return nil, errors.Wrap(err, "search tv")
	}
	return r, nil
}

func (s *Server) SearchMedia(c *gin.Context) (interface{}, error) {
	var q searchTvParam
	if err := c.ShouldBindQuery(&q); err != nil {
		return nil, errors.Wrap(err, "bind query")
	}
	log.Infof("search media with keyword: %v", q.Query)
	tmdb, err := s.TMDB()
	if err != nil {
		return nil, err
	}
	r, err := tmdb.SearchMedia(q.Query, s.language, q.Page)
	if err != nil {
		return nil, errors.Wrap(err, "search tv")
	}
	for i, res := range r.Results {
		if s.db.TmdbIdInWatchlist(int(res.ID)) {
			r.Results[i].InWatchlist = true
		}
	}
	return r, nil

}

type addWatchlistIn struct {
	TmdbID                  int    `json:"tmdb_id" binding:"required"`
	StorageID               int    `json:"storage_id" `
	Resolution              string `json:"resolution" binding:"required"`
	Folder                  string `json:"folder" binding:"required"`
	DownloadHistoryEpisodes bool   `json:"download_history_episodes"` //for tv
	SizeMin                 int    `json:"size_min"`
	SizeMax                 int    `json:"size_max"`
}

func (s *Server) AddTv2Watchlist(c *gin.Context) (interface{}, error) {
	var in engine.AddWatchlistIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind query")
	}
	return s.core.AddTv2Watchlist(in)
}

func (s *Server) AddMovie2Watchlist(c *gin.Context) (interface{}, error) {
	var in engine.AddWatchlistIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind query")
	}
	return s.core.AddMovie2Watchlist(in)
}

type MediaWithStatus struct {
	*ent.Media
	MonitoredNum  int `json:"monitored_num"`
	DownloadedNum int `json:"downloaded_num"`
}

//missing: episode aired missing
//downloaded: all monitored episode downloaded
//monitoring: episode aired downloaded, but still has not aired episode
//for movie, only monitoring/downloaded

func (s *Server) GetTvWatchlist(c *gin.Context) (interface{}, error) {
	list := s.db.GetMediaWatchlist(media.MediaTypeTv)
	res := make([]MediaWithStatus, len(list))
	for i, item := range list {
		var ms = MediaWithStatus{
			Media:         item,
			MonitoredNum:  0,
			DownloadedNum: 0,
		}
		mon, ok1 := s.monitorNumCache.Get(item.ID)
		dow, ok2 := s.downloadNumCache.Get(item.ID)
		if ok1 && ok2 {
			ms.MonitoredNum = mon
			ms.DownloadedNum = dow
		} else {
			details, err := s.db.GetMediaDetails(item.ID)
			if err != nil {
				return nil, errors.Wrap(err, "get details")
			}
			for _, ep := range details.Episodes {
				if ep.Monitored {
					ms.MonitoredNum++
					if ep.Status == episode.StatusDownloaded {
						ms.DownloadedNum++
					}
				}
			}
			s.monitorNumCache.Set(item.ID, ms.MonitoredNum)
			s.downloadNumCache.Set(item.ID, ms.DownloadedNum)
		}

		res[i] = ms
	}
	return res, nil
}

func (s *Server) GetMovieWatchlist(c *gin.Context) (interface{}, error) {
	list := s.db.GetMediaWatchlist(media.MediaTypeMovie)
	res := make([]MediaWithStatus, len(list))
	for i, item := range list {
		var ms = MediaWithStatus{
			Media:         item,
			MonitoredNum:  1,
			DownloadedNum: 0,
		}
		dummyEp, err := s.db.GetMovieDummyEpisode(item.ID)
		if err != nil {
			log.Errorf("get dummy episode: %v", err)
		} else {
			if dummyEp.Status == episode.StatusDownloaded {
				ms.DownloadedNum++
			}
		}
		res[i] = ms
	}
	return res, nil
}

type MediaDetails struct {
	*db.MediaDetails
	Storage *ent.Storage `json:"storage"`
}

func (s *Server) GetMediaDetails(c *gin.Context) (interface{}, error) {
	ids := c.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		return nil, errors.Wrap(err, "convert")
	}
	detail, err := s.db.GetMediaDetails(id)
	if err != nil {
		return nil, errors.Wrap(err, "get details")
	}
	st := s.db.GetStorage(detail.StorageID)
	return MediaDetails{MediaDetails: detail, Storage: &st.Storage}, nil
}

func (s *Server) DeleteFromWatchlist(c *gin.Context) (interface{}, error) {
	ids := c.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		return nil, errors.Wrap(err, "convert")
	}

	deleteFiles := c.Query("delete_files")
	if strings.ToLower(deleteFiles) == "true" {
		//will delete local media file
		log.Infof("will delete local media files for %d", id)
		m, err := s.db.GetMedia(id)
		if err != nil {
			log.Warnf("get media: %v", err)
		} else {
			st, err := s.core.GetStorage(m.StorageID, m.MediaType)
			if err != nil {
				log.Warnf("get storage error: %v", err)
			} else {
				if err := st.RemoveAll(m.TargetDir); err != nil {
					log.Warnf("remove all : %v", err)
				} else {
					log.Infof("delete media files success: %v", m.TargetDir)
				}
			}
		}
	}

	if err := s.db.DeleteMedia(id); err != nil {
		return nil, errors.Wrap(err, "delete db")
	}
	os.RemoveAll(filepath.Join(db.ImgPath, ids)) //delete image related

	return "success", nil
}
