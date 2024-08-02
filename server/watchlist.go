package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"polaris/db"
	"polaris/ent"
	"polaris/ent/episode"
	"polaris/ent/media"
	"polaris/log"
	"strconv"
	"time"

	tmdb "github.com/cyruzin/golang-tmdb"
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
}

func (s *Server) AddTv2Watchlist(c *gin.Context) (interface{}, error) {
	var in addWatchlistIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind query")
	}
	log.Debugf("add tv watchlist input %+v", in)
	if in.Folder == "" {
		return nil, errors.New("folder should be provided")
	}
	detailCn, err := s.MustTMDB().GetTvDetails(in.TmdbID, db.LanguageCN)
	if err != nil {
		return nil, errors.Wrap(err, "get tv detail")
	}
	var nameCn = detailCn.Name

	detailEn, _ := s.MustTMDB().GetTvDetails(in.TmdbID, db.LanguageEN)
	var nameEn = detailEn.Name
	var detail *tmdb.TVDetails
	if s.language == "" || s.language == db.LanguageCN {
		detail = detailCn
	} else {
		detail = detailEn
	}
	log.Infof("find detail for tv id %d: %v", in.TmdbID, detail)

	var epIds []int
	for _, season := range detail.Seasons {
		seasonId := season.SeasonNumber
		se, err := s.MustTMDB().GetSeasonDetails(int(detail.ID), seasonId, s.language)
		if err != nil {
			log.Errorf("get season detail (%s) error: %v", detail.Name, err)
			continue
		}
		for _, ep := range se.Episodes {
			epid, err := s.db.SaveEposideDetail(&ent.Episode{
				SeasonNumber:  seasonId,
				EpisodeNumber: ep.EpisodeNumber,
				Title:         ep.Name,
				Overview:      ep.Overview,
				AirDate:       ep.AirDate,
			})
			if err != nil {
				log.Errorf("save episode info error: %v", err)
				continue
			}
			epIds = append(epIds, epid)
		}
	}
	r, err := s.db.AddMediaWatchlist(&ent.Media{
		TmdbID:                  int(detail.ID),
		MediaType:               media.MediaTypeTv,
		NameCn:                  nameCn,
		NameEn:                  nameEn,
		OriginalName:            detail.OriginalName,
		Overview:                detail.Overview,
		AirDate:                 detail.FirstAirDate,
		Resolution:              media.Resolution(in.Resolution),
		StorageID:               in.StorageID,
		TargetDir:               in.Folder,
		DownloadHistoryEpisodes: in.DownloadHistoryEpisodes,
	}, epIds)
	if err != nil {
		return nil, errors.Wrap(err, "add to list")
	}
	go func() {
		if err := s.downloadPoster(detail.PosterPath, r.ID); err != nil {
			log.Errorf("download poster error: %v", err)
		}
		if err := s.downloadBackdrop(detail.BackdropPath, r.ID); err != nil {
			log.Errorf("download poster error: %v", err)
		}
		if err := s.core.CheckDownloadedSeriesFiles(r); err != nil {
			log.Errorf("check downloaded files error: %v", err)
		}

	}()

	log.Infof("add tv %s to watchlist success", detail.Name)
	return nil, nil
}

func (s *Server) AddMovie2Watchlist(c *gin.Context) (interface{}, error) {
	var in addWatchlistIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind query")
	}
	detailCn, err := s.MustTMDB().GetMovieDetails(in.TmdbID, db.LanguageCN)
	if err != nil {
		return nil, errors.Wrap(err, "get movie detail")
	}
	var nameCn = detailCn.Title

	detailEn, _ := s.MustTMDB().GetMovieDetails(in.TmdbID, db.LanguageEN)
	var nameEn = detailEn.Title
	var detail *tmdb.MovieDetails
	if s.language == "" || s.language == db.LanguageCN {
		detail = detailCn
	} else {
		detail = detailEn
	}
	log.Infof("find detail for movie id %d: %v", in.TmdbID, detail)

	epid, err := s.db.SaveEposideDetail(&ent.Episode{
		SeasonNumber:  1,
		EpisodeNumber: 1,
		Title:         "dummy episode for movies",
		Overview:      "dummy episode for movies",
		AirDate:       detail.ReleaseDate,
	})
	if err != nil {
		return nil, errors.Wrap(err, "add dummy episode")
	}
	log.Infof("added dummy episode for movie: %v", nameEn)

	r, err := s.db.AddMediaWatchlist(&ent.Media{
		TmdbID:       int(detail.ID),
		MediaType:    media.MediaTypeMovie,
		NameCn:       nameCn,
		NameEn:       nameEn,
		OriginalName: detail.OriginalTitle,
		Overview:     detail.Overview,
		AirDate:      detail.ReleaseDate,
		Resolution:   media.Resolution(in.Resolution),
		StorageID:    in.StorageID,
		TargetDir:    in.Folder,
	}, []int{epid})
	if err != nil {
		return nil, errors.Wrap(err, "add to list")
	}
	go func() {
		if err := s.downloadPoster(detail.PosterPath, r.ID); err != nil {
			log.Errorf("download poster error: %v", err)
		}
		if err := s.downloadBackdrop(detail.BackdropPath, r.ID); err != nil {
			log.Errorf("download backdrop error: %v", err)
		}
	}()

	log.Infof("add movie %s to watchlist success", detail.Title)
	return nil, nil

}

func (s *Server) downloadBackdrop(path string, mediaID int) error {
	url := "https://image.tmdb.org/t/p/original" + path
	return s.downloadImage(url, mediaID, "backdrop.jpg")
}

func (s *Server) downloadPoster(path string, mediaID int) error {
	var url = "https://image.tmdb.org/t/p/original" + path

	return s.downloadImage(url, mediaID, "poster.jpg")
}

func (s *Server) downloadImage(url string, mediaID int, name string) error {

	log.Infof("try to download image: %v", url)
	var resp, err = http.Get(url)
	if err != nil {
		return errors.Wrap(err, "http get")
	}
	targetDir := fmt.Sprintf("%v/%d", db.ImgPath, mediaID)
	os.MkdirAll(targetDir, 0777)
	//ext := filepath.Ext(path)
	targetFile := filepath.Join(targetDir, name)
	f, err := os.Create(targetFile)
	if err != nil {
		return errors.Wrap(err, "new file")
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return errors.Wrap(err, "copy http response")
	}
	log.Infof("image successfully downlaoded: %v", targetFile)
	return nil

}

type MediaWithStatus struct {
	*ent.Media
	MonitoredNum  int    `json:"monitored_num"`
	DownloadedNum int    `json:"downloaded_num"`
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
			Media:      item,
			MonitoredNum:  0,
			DownloadedNum: 0,
		}

		details := s.db.GetMediaDetails(item.ID)

		for _, ep := range details.Episodes {
			monitored := false
			if ep.SeasonNumber == 0 {
				continue
			}
			if item.DownloadHistoryEpisodes {
				monitored = true
			} else {
				t, err := time.Parse("2006-01-02", ep.AirDate)
				if err != nil { //airdate not exist, maybe airdate not set yet
					monitored = true
				} else {
					if item.CreatedAt.Sub(t) > 24*time.Hour { //剧集在加入watchlist之前，不去下载
						continue
					}
					monitored = true
				}
			}
			if monitored {
				ms.MonitoredNum++
				if ep.Status == episode.StatusDownloaded {
					ms.DownloadedNum++
				}
			}

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
			Media:  item,
			MonitoredNum: 1,
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
	detail := s.db.GetMediaDetails(id)
	st := s.db.GetStorage(detail.StorageID)
	return MediaDetails{MediaDetails: detail, Storage: &st.Storage}, nil
}

func (s *Server) DeleteFromWatchlist(c *gin.Context) (interface{}, error) {
	ids := c.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		return nil, errors.Wrap(err, "convert")
	}
	if err := s.db.DeleteMedia(id); err != nil {
		return nil, errors.Wrap(err, "delete db")
	}
	os.RemoveAll(filepath.Join(db.ImgPath, ids)) //delete image related
	return "success", nil
}
