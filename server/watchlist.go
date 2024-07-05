package server

import (
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
	ID         int `json:"id" binding:"required"`
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
	 
	if err := s.db.AddWatchlist(in.RootFolder, detail); err != nil {
		return nil, errors.Wrap(err, "add to list")
	}
	log.Infof("save watchlist success: %s", detail.Name)

	for _, season := range detail.Seasons {
		seasonId := season.SeasonNumber
		for  i := 1; i <= season.EpisodeCount; i++ {
			ep, err := s.MustTMDB().GetEposideDetail(int(detail.ID), seasonId, i, s.language)
			if err != nil {
				log.Errorf("get eposide detail: %v", err)
				return nil, errors.Wrap(err, "get eposide detail")
			}
			err = s.db.SaveEposideDetail(int(detail.ID), ep)
			if err != nil {
				return nil, errors.Wrap(err, "save episode")
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