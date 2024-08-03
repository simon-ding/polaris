package server

import (
	"fmt"
	"polaris/ent/media"
	"polaris/log"
	"polaris/pkg/torznab"
	"polaris/server/core"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (s *Server) searchAndDownloadSeasonPackage(seriesId, seasonNum int) (*string, error) {

	res, err := core.SearchTvSeries(s.db, seriesId, seasonNum, nil, true, true)
	if err != nil {
		return nil, err
	}

	r1 := res[0]
	log.Infof("found resource to download: %+v", r1)
	return s.core.DownloadSeasonPackage(r1, seriesId, seasonNum)

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
			res, err = core.SearchTvSeries(s.db, in.ID, in.Season, nil, false, false)
			if err != nil {
				return nil, errors.Wrap(err, "search season package")
			}
		} else {
			log.Infof("search series episode S%02dE%02d", in.Season, in.Episode)
			res, err = core.SearchTvSeries(s.db, in.ID, in.Season, []int{in.Episode}, false, false)
			if err != nil {
				if err.Error() == "no resource found" {
					return []string{}, nil
				}
				return nil, errors.Wrap(err, "search episode")
			}

		}
	} else {
		log.Info("search movie %d", in.ID)
		res, err = core.SearchMovie(s.db, in.ID, false, false)
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
		name = *name1
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
			return s.core.DownloadSeasonPackage(res, in.MediaID, in.Season)
		}
		name := in.Name
		if name == "" {
			name = fmt.Sprintf("%v S%02dE%02d", m.OriginalName, in.Season, in.Episode)
		}
		res := torznab.Result{Name: name, Link: in.Link, Size: in.Size, IndexerId: in.IndexerId}
		return s.core.DownloadEpisodeTorrent(res, in.MediaID, in.Season, in.Episode)
	} else {
		//movie
		return s.core.DownloadMovie(m, in.Link, in.Name, in.Size, in.IndexerId)
	}

}
