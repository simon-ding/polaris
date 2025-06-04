package server

import (
	"os"
	"polaris/internal/db"
	"polaris/log"
	"polaris/pkg/metadata"
	"polaris/pkg/uptime"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type LogFile struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

func (s *Server) GetAllLogs(c *gin.Context) (interface{}, error) {
	fs, err := os.ReadDir(db.LogPath)
	if err != nil {
		return []LogFile{}, nil
	}
	var logs []LogFile
	for _, f := range fs {
		if f.IsDir() {
			continue
		}
		info, err := f.Info()
		if err != nil {
			log.Warnf("get log file error: %v", err)
			continue
		}
		l := LogFile{
			Name: f.Name(),
			Size: info.Size(),
		}
		logs = append(logs, l)
	}
	return logs, nil
}

func (s *Server) About(c *gin.Context) (interface{}, error) {

	return gin.H{
		"intro":      "Polaris © Simon Ding",
		"homepage":   "https://github.com/simon-ding/polaris",
		"uptime":     uptime.Uptime(),
		"chat_group": "https://t.me/+8R2nzrlSs2JhMDgx",
		"go_version": runtime.Version(),
		"version":    db.Version,
	}, nil
}

type parseIn struct {
	S         string `json:"s" binding:"required"`
}

func (s *Server) ParseTv(c *gin.Context) (interface{}, error) {
	var in parseIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind")
	}
	return metadata.ParseTv(in.S), nil
}

func (s *Server) ParseMovie(c *gin.Context) (interface{}, error) {
	var in parseIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind")
	}
	return metadata.ParseMovie(in.S), nil
}
