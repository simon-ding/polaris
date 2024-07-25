package server

import (
	"fmt"
	"polaris/db"
	"polaris/log"
	"polaris/pkg/transmission"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type GeneralSettings struct {
	TmdbApiKey  string `json:"tmdb_api_key"`
	DownloadDir string `json:"download_dir"`
}

func (s *Server) SetSetting(c *gin.Context) (interface{}, error) {
	var in GeneralSettings
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind json")
	}
	log.Infof("set setting input: %+v", in)
	if in.TmdbApiKey != "" {
		if err := s.db.SetSetting(db.SettingTmdbApiKey, in.TmdbApiKey); err != nil {
			return nil, errors.Wrap(err, "save tmdb api")
		}
	}
	if in.DownloadDir != "" {
		if err := s.db.SetSetting(db.SettingDownloadDir, in.DownloadDir); err != nil {
			return nil, errors.Wrap(err, "save download dir")
		}
	}
	return nil, nil
}

func (s *Server) GetSetting(c *gin.Context) (interface{}, error) {
	tmdb := s.db.GetSetting(db.SettingTmdbApiKey)
	downloadDir := s.db.GetSetting(db.SettingDownloadDir)

	return &GeneralSettings{
		TmdbApiKey:  tmdb,
		DownloadDir: downloadDir,
	}, nil
}

type addTorznabIn struct {
	Name   string `json:"name" binding:"required"`
	URL    string `json:"url" binding:"required"`
	ApiKey string `json:"api_key" binding:"required"`
}

func (s *Server) AddTorznabInfo(c *gin.Context) (interface{}, error) {
	var in addTorznabIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind json")
	}
	err := s.db.SaveTorznabInfo(in.Name, db.TorznabSetting{
		URL:    in.URL,
		ApiKey: in.ApiKey,
	})
	if err != nil {
		return nil, errors.Wrap(err, "add ")
	}
	return nil, nil
}

func (s *Server) DeleteTorznabInfo(c *gin.Context) (interface{}, error) {
	var ids = c.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		return nil, fmt.Errorf("id is not correct: %v", ids)
	}
	s.db.DeleteTorznab(id)
	return "success", nil
}

func (s *Server) GetAllIndexers(c *gin.Context) (interface{}, error) {
	indexers := s.db.GetAllTorznabInfo()
	if len(indexers) == 0 {
		return nil, nil
	}
	return indexers, nil
}

func (s *Server) getDownloadClient() (*transmission.Client, error) {
	tr := s.db.GetTransmission()
	trc, err := transmission.NewClient(transmission.Config{
		URL:      tr.URL,
		User:     tr.User,
		Password: tr.Password,
	})
	if err != nil {
		return nil, errors.Wrap(err, "connect transmission")
	}
	return trc, nil
}


type downloadClientIn struct {
	Name           string `json:"name" binding:"required"`
	URL            string `json:"url" binding:"required"`
	User           string `json:"user"`
	Password       string `json:"password"`
	Implementation string `json:"implementation" binding:"required"`
}

func (s *Server) AddDownloadClient(c *gin.Context) (interface{}, error) {
	var in downloadClientIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind json")
	}
	//test connection
	_, err := transmission.NewClient(transmission.Config{
		URL: in.URL,
		User: in.User,
		Password: in.Password,
	})
	if err != nil {
		return nil, errors.Wrap(err, "tranmission setting")
	}
	if err := s.db.SaveTransmission(in.Name, in.URL, in.User, in.Password); err != nil {
		return nil, errors.Wrap(err, "save transmission")
	}
	return nil, nil
}

func (s *Server) GetAllDonloadClients(c *gin.Context) (interface{}, error) {
	res := s.db.GetAllDonloadClients()
	if len(res) == 0 {
		return nil, nil
	}
	return res, nil
}

func (s *Server) DeleteDownloadCLient(c *gin.Context) (interface{}, error) {
	var ids = c.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		return nil, fmt.Errorf("id is not correct: %v", ids)
	}
	s.db.DeleteDownloadCLient(id)
	return "success", nil
}

type logLovelIn struct {
	Level string `json:"level"`
}

func (s *Server) SetLogLevel(c *gin.Context) (interface{}, error) {
	var in logLovelIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind json")
	}
	log.SetLogLevel(in.Level)
	return "success", nil
}