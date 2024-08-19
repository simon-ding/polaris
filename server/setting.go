package server

import (
	"encoding/json"
	"fmt"
	"polaris/db"
	"polaris/ent"
	"polaris/log"
	"polaris/pkg/transmission"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type GeneralSettings struct {
	TmdbApiKey         string `json:"tmdb_api_key"`
	DownloadDir        string `json:"download_dir"`
	LogLevel           string `json:"log_level"`
	Proxy              string `json:"proxy"`
	EnablePlexmatch    bool   `json:"enable_plexmatch"`
	EnableNfo          bool   `json:"enable_nfo"`
	AllowQiangban      bool   `json:"allow_qiangban"`
	EnableAdultContent bool   `json:"enable_adult_content"`
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
		log.Info("set download dir to %s", in.DownloadDir)
		if err := s.db.SetSetting(db.SettingDownloadDir, in.DownloadDir); err != nil {
			return nil, errors.Wrap(err, "save download dir")
		}
	}
	if in.LogLevel != "" {
		log.SetLogLevel(in.LogLevel)
		if err := s.db.SetSetting(db.SettingLogLevel, in.LogLevel); err != nil {
			return nil, errors.Wrap(err, "save log level")
		}

	}

	plexmatchEnabled := s.db.GetSetting(db.SettingPlexMatchEnabled)
	if in.EnablePlexmatch && plexmatchEnabled != "true" {
		s.db.SetSetting(db.SettingPlexMatchEnabled, "true")
	} else if !in.EnablePlexmatch && plexmatchEnabled != "false" {
		s.db.SetSetting(db.SettingPlexMatchEnabled, "false")
	}

	s.db.SetSetting(db.SettingProxy, in.Proxy)

	if in.AllowQiangban {
		s.db.SetSetting(db.SettingAllowQiangban, "true")
	} else {
		s.db.SetSetting(db.SettingAllowQiangban, "false")
	}

	if in.EnableNfo {
		s.db.SetSetting(db.SettingNfoSupportEnabled, "true")
	} else {
		s.db.SetSetting(db.SettingNfoSupportEnabled, "false")
	}
	if in.EnableAdultContent {
		s.db.SetSetting(db.SettingEnableTmdbAdultContent, "true")
	} else {
		s.db.SetSetting(db.SettingEnableTmdbAdultContent, "false")
	}

	return nil, nil
}

func (s *Server) GetSetting(c *gin.Context) (interface{}, error) {
	tmdb := s.db.GetSetting(db.SettingTmdbApiKey)
	downloadDir := s.db.GetSetting(db.SettingDownloadDir)
	logLevel := s.db.GetSetting(db.SettingLogLevel)
	plexmatchEnabled := s.db.GetSetting(db.SettingPlexMatchEnabled)
	allowQiangban := s.db.GetSetting(db.SettingAllowQiangban)
	enableNfo := s.db.GetSetting(db.SettingNfoSupportEnabled)
	enableAdult := s.db.GetSetting(db.SettingEnableTmdbAdultContent)
	return &GeneralSettings{
		TmdbApiKey:      tmdb,
		DownloadDir:     downloadDir,
		LogLevel:        logLevel,
		Proxy:           s.db.GetSetting(db.SettingProxy),
		EnablePlexmatch: plexmatchEnabled == "true",
		AllowQiangban:   allowQiangban == "true",
		EnableNfo:       enableNfo == "true",
		EnableAdultContent: enableAdult == "true",
	}, nil
}

type addTorznabIn struct {
	ID        int     `json:"id"`
	Name      string  `json:"name" binding:"required"`
	URL       string  `json:"url" binding:"required"`
	ApiKey    string  `json:"api_key" binding:"required"`
	Disabled  bool    `json:"disabled"`
	Priority  int     `json:"priority"`
	SeedRatio float32 `json:"seed_ratio"`
}

func (s *Server) AddTorznabInfo(c *gin.Context) (interface{}, error) {
	var in addTorznabIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind json")
	}

	log.Infof("add indexer settings: %+v", in)
	setting := db.TorznabSetting{
		URL:    in.URL,
		ApiKey: in.ApiKey,
	}

	data, err := json.Marshal(setting)
	if err != nil {
		return nil, errors.Wrap(err, "marshal json")
	}
	if in.Priority > 128 {
		in.Priority = 128
	}

	indexer := ent.Indexers{
		ID:             in.ID,
		Name:           in.Name,
		Implementation: "torznab",
		Settings:       string(data),
		Priority:       in.Priority,
		Disabled:       in.Disabled,
		SeedRatio:      in.SeedRatio,
	}
	err = s.db.SaveIndexer(&indexer)
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

func (s *Server) getDownloadClient() (*transmission.Client, *ent.DownloadClients, error) {
	tr := s.db.GetTransmission()
	trc, err := transmission.NewClient(transmission.Config{
		URL:      tr.URL,
		User:     tr.User,
		Password: tr.Password,
	})
	if err != nil {
		return nil, nil, errors.Wrap(err, "connect transmission")
	}
	return trc, tr, nil
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
		URL:      in.URL,
		User:     in.User,
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

type episodeMonitoringIn struct {
	EpisodeID int  `json:"episode_id"`
	Monitor   bool `json:"monitor"`
}

func (s *Server) ChangeEpisodeMonitoring(c *gin.Context) (interface{}, error) {
	var in episodeMonitoringIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind")
	}
	s.db.SetEpisodeMonitoring(in.EpisodeID, in.Monitor)
	return "success", nil
}

func (s *Server) EditMediaMetadata(c *gin.Context) (interface{}, error) {
	var in db.EditMediaData
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind")
	}
	err := s.db.EditMediaMetadata(in)
	if err != nil {
		return nil, errors.Wrap(err, "save db")
	}
	return "success", nil
}
