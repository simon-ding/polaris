package server

import (
	"polaris/db"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)


type GeneralSettings struct {
	TmdbApiKey string `json:"tmdb_api_key"`
	DownloadDir string `json:"download_dir"`
}
func (s *Server) SetSetting(c *gin.Context) (interface{}, error) {
	var in GeneralSettings
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind json")
	}
	if in.TmdbApiKey != "" {
		if err := s.db.SetSetting(db.SettingTmdbApiKey, in.TmdbApiKey); err != nil {
			return nil, errors.Wrap(err, "save tmdb api")
		}
	}
	if in.DownloadDir == "" {
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
		TmdbApiKey: tmdb,
		DownloadDir: downloadDir,
	}, nil
}
