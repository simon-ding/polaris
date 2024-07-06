package server

import (
	"polaris/ent"
	"polaris/log"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type searchTvParam struct {
	Query string `form:"query"`
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

type addWatchlistIn struct {
	ID         int    `json:"id" binding:"required"`
	RootFolder string `json:"folder" binding:"required"`
}

func (s *Server) AddWatchlist(c *gin.Context) (interface{}, error) {
	var in addWatchlistIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind query")
	}
	detail, err := s.MustTMDB().GetTvDetails(in.ID, s.language)
	if err != nil {
		return nil, errors.Wrap(err, "get tv detail")
	}
	log.Infof("find detail for tv id %d: %v", in.ID, detail)

	wl, err := s.db.AddWatchlist(in.RootFolder, detail)
	if err != nil {
		return nil, errors.Wrap(err, "add to list")
	}
	log.Infof("save watchlist success: %s", detail.Name)

	for _, season := range detail.Seasons {
		seasonId := season.SeasonNumber
		se, err := s.MustTMDB().GetSeasonDetails(int(detail.ID), seasonId, s.language)
		if err != nil {
			log.Errorf("get season detail (%s) error: %v", detail.Name, err)
			continue
		}
		for _, ep := range se.Episodes {
			err = s.db.SaveEposideDetail(&ent.Epidodes{
				SeriesID:      wl.ID,
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
		}
	}
	log.Infof("add tv %s to watchlist success", detail.Name)
	return nil, nil
}

func (s *Server) GetWatchlist(c *gin.Context) (interface{}, error) {
	list := s.db.GetWatchlist()
	return list, nil
}
