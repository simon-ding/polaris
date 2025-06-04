package server

import (
	"fmt"
	"polaris/internal/db"
	"polaris/internal/biz/engine"
	"polaris/ent/media"
	"polaris/log"
	"polaris/pkg/torznab"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (s *Server) searchAndDownloadSeasonPackage(seriesId, seasonNum int) (*string, error) {

	res, err := engine.SearchTvSeries(s.db, &engine.SearchParam{
		MediaId:         seriesId,
		SeasonNum:       seasonNum,
		Episodes:        nil,
		CheckResolution: true,
		CheckFileSize:   true,
	})
	if err != nil {
		return nil, err
	}

	r1 := res[0]
	log.Infof("found resource to download: %+v", r1)
	return s.core.DownloadEpisodeTorrent(r1, engine.DownloadOptions{
		SeasonNum: seasonNum,
		MediaId:   seriesId,
	})

}

type searchAndDownloadIn struct {
	ID      int `json:"id" binding:"required"`
	Season  int `json:"season"`
	Episode int `json:"episode"`
}

func (s *Server) SearchAvailableTorrents(c *gin.Context) (interface{}, error) {
	var in searchAndDownloadIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind json")
	}
	m, err := s.db.GetMedia(in.ID)
	if err != nil {
		return nil, errors.Wrap(err, "get media")
	}
	log.Infof("search torrents resources link: %+v", in)

	var res []torznab.Result
	if m.MediaType == media.MediaTypeTv {
		if in.Episode == 0 {
			//search season package
			log.Infof("search series season package S%02d", in.Season)
			res, err = engine.SearchTvSeries(s.db, &engine.SearchParam{
				MediaId:   in.ID,
				SeasonNum: in.Season,
				Episodes:  nil,
			})
			if err != nil {
				return nil, errors.Wrap(err, "search season package")
			}
		} else {
			log.Infof("search series episode S%02dE%02d", in.Season, in.Episode)
			res, err = engine.SearchTvSeries(s.db, &engine.SearchParam{
				MediaId:   in.ID,
				SeasonNum: in.Season,
				Episodes:  []int{in.Episode},
			})
			if err != nil {
				if err.Error() == "no resource found" {
					return []string{}, nil
				}
				return nil, errors.Wrap(err, "search episode")
			}

		}
	} else {
		log.Info("search movie %d", in.ID)
		qiangban := s.db.GetSetting(db.SettingAllowQiangban)
		allowQiangban := false
		if qiangban == "true" {
			allowQiangban = true
		}

		res, err = engine.SearchMovie(s.db, &engine.SearchParam{
			MediaId:        in.ID,
			FilterQiangban: !allowQiangban,
		})
		if err != nil {
			if err.Error() == "no resource found" {
				return []string{}, nil
			}
			return nil, err
		}
	}
	return res, nil
}

func (s *Server) SearchTvAndDownload(c *gin.Context) (interface{}, error) {
	var in searchAndDownloadIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind json")
	}
	log.Infof("search episode resources link: %v", in)
	var name string
	if in.Episode == 0 {
		log.Infof("season package search")
		//search season package
		name1, err := s.searchAndDownloadSeasonPackage(in.ID, in.Season)
		if err != nil {
			return nil, errors.Wrap(err, "download")
		}
		name = *name1
	} else {
		log.Infof("season episode search")
		name1, err := s.core.SearchAndDownload(in.ID, in.Season, in.Episode)
		if err != nil {
			return nil, errors.Wrap(err, "download")
		}
		if len(name1) == 0 {
			return nil, fmt.Errorf("no torrent found")
		}
		name = name1[0]
	}

	return gin.H{
		"name": name,
	}, nil
}

type downloadTorrentIn struct {
	MediaID int `json:"id" binding:"required"`
	Season  int `json:"season"`
	Episode int `json:"episode"`
	torznab.Result
}

func (s *Server) DownloadTorrent(c *gin.Context) (interface{}, error) {
	var in downloadTorrentIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind json")
	}
	log.Infof("download torrent input: %+v", in)

	m, err := s.db.GetMedia(in.MediaID)
	if err != nil {
		return nil, fmt.Errorf("no tv series of id %v", in.MediaID)
	}
	if m.MediaType == media.MediaTypeTv {
		if in.Episode == 0 {
			//download season package
			name := in.Name
			if name == "" {
				name = fmt.Sprintf("%v S%02d", m.OriginalName, in.Season)
			}
			res := torznab.Result{Name: name, Link: in.Link, Size: in.Size}
			return s.core.DownloadEpisodeTorrent(res, engine.DownloadOptions{
				SeasonNum: in.Season,
				MediaId:   in.MediaID,
			})
		}
		name := in.Name
		if name == "" {
			name = fmt.Sprintf("%v S%02dE%02d", m.OriginalName, in.Season, in.Episode)
		}
		res := torznab.Result{Name: name, Link: in.Link, Size: in.Size, IndexerId: in.IndexerId}
		return s.core.DownloadEpisodeTorrent(res, engine.DownloadOptions{
			SeasonNum: in.Season,
			MediaId:   in.MediaID,
			EpisodeNums: []int{in.Episode},
		})
	} else {
		//movie
		name := in.Name
		if name == "" {
			name = m.OriginalName
		}

		res := torznab.Result{Name: name, Link: in.Link, Size: in.Size, IndexerId: in.IndexerId}
		return s.core.DownloadMovie(m, res)
	}

}

func (s *Server) DownloadAll(c *gin.Context) (interface{}, error) {
	ids := c.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		return nil, errors.Wrap(err, "convert")
	}
	return s.downloadAllEpisodes(id)
}

func (s *Server) downloadAllEpisodes(id int) (interface{}, error) {
	m, err := s.db.GetMedia(id)
	if err != nil {
		return nil, errors.Wrap(err, "get media")
	}
	if m.MediaType == media.MediaTypeTv {
		return s.core.DownloadSeriesAllEpisodes(m.ID), nil
	}
	name, err := s.core.DownloadMovieByID(m.ID)

	return []string{name}, err
}

func (s *Server) DownloadAllTv(c *gin.Context) (interface{}, error) {
	tvs := s.db.GetMediaWatchlist(media.MediaTypeTv)
	var allNames []string
	for _, tv := range tvs {
		names, err := s.downloadAllEpisodes(tv.ID)
		if err == nil {
			allNames = append(allNames, names.([]string)...)
		}
	}
	return allNames, nil
}

func (s *Server) DownloadAllMovies(c *gin.Context) (interface{}, error) {
	movies := s.db.GetMediaWatchlist(media.MediaTypeMovie)
	var allNames []string
	for _, mv := range movies {
		names, err := s.downloadAllEpisodes(mv.ID)
		if err == nil {
			allNames = append(allNames, names.([]string)...)
		}
	}
	return allNames, nil
}
