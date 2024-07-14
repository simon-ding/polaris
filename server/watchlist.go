package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"polaris/db"
	"polaris/ent"
	"polaris/log"
	"strconv"

	tmdb "github.com/cyruzin/golang-tmdb"
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
	detailCn, err := s.MustTMDB().GetTvDetails(in.TmdbID, db.LanguageCN)
	if err != nil {
		return nil, errors.Wrap(err, "get tv detail")
	}
	var nameCn = detailCn.Name

	detailEn, _ := s.MustTMDB().GetTvDetails(in.TmdbID, db.LanguageEN)
	var nameEn = detailEn.Name
	var detail *tmdb.TVDetails
	if s.language == "" || s.language ==db.LanguageCN {
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
	r, err := s.db.AddWatchlist(in.StorageID, nameCn, nameEn, detail, epIds, db.R1080p)
	if err != nil {
		return nil, errors.Wrap(err, "add to list")
	}
	go func ()  {
		if err := s.downloadPoster(detail.PosterPath, r.ID); err != nil {
			log.Errorf("download poster error: %v", err)
		}	
	}()

	log.Infof("add tv %s to watchlist success", detail.Name)
	return nil, nil
}

func (s *Server) downloadPoster(path string, seriesId int)  error{
	var tmdbImgBaseUrl = "https://image.tmdb.org/t/p/w500/"
	url := tmdbImgBaseUrl + path
	log.Infof("try to download poster: %v", url)
	var resp, err = http.Get(url)
	if err != nil {
		return errors.Wrap(err, "http get")
	}
	targetDir := fmt.Sprintf("%v/%d", db.ImgPath, seriesId)
	os.MkdirAll(targetDir, 0777)
	ext := filepath.Ext(path)
	targetFile := filepath.Join(targetDir, "poster"+ ext)
	f, err := os.Create(targetFile)
	if err != nil {
		return errors.Wrap(err, "new file")
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return errors.Wrap(err, "copy http response")
	}
	log.Infof("poster successfully downlaoded: %v", targetFile)
	return nil
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

func (s *Server) DeleteFromWatchlist(c *gin.Context) (interface{}, error) {
	ids := c.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		return nil, errors.Wrap(err, "convert")
	}
	if err := s.db.DeleteSeries(id); err != nil {
		return nil, errors.Wrap(err, "delete db")
	}
	return "success", nil
}