package server

import (
	"polaris/db"
	"polaris/ent"
	"polaris/log"
	"strconv"

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
	TmdbID    int `json:"tmdb_id" binding:"required"`
	StorageID int `json:"storage_id" binding:"required"`
	Resolution db.ResolutionType `json:"resolution" binding:"required"`
}

func (s *Server) AddWatchlist(c *gin.Context) (interface{}, error) {
	var in addWatchlistIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind query")
	}
	detail, err := s.MustTMDB().GetTvDetails(in.TmdbID, s.language)
	if err != nil {
		return nil, errors.Wrap(err, "get tv detail")
	}
	log.Infof("find detail for tv id %d: %v", in.TmdbID, detail)
	detailEn, _ := s.MustTMDB().GetTvDetails(in.TmdbID, "en-US")
	var nameEn = detailEn.Name

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
	_, err = s.db.AddWatchlist(in.StorageID, nameEn, detail, epIds, db.R1080p)
	if err != nil {
		return nil, errors.Wrap(err, "add to list")
	}

	log.Infof("add tv %s to watchlist success", detail.Name)
	return nil, nil
}

func (s *Server) GetWatchlist(c *gin.Context) (interface{}, error) {
	list := s.db.GetWatchlist()
	return list, nil
}

func (s *Server) GetTvDetails(c *gin.Context) (interface{}, error) {
	ids := c.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		return nil, errors.Wrap(err, "convert")
	}
	detail := s.db.GetSeriesDetails(id)
	return detail, nil
}

func (s *Server) GetAvailableResolutions(c *gin.Context) (interface{}, error) {
	return []db.ResolutionType{
		db.R720p,
		db.R1080p,
		db.R4k,
	}, nil
}