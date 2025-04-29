package server

import (
	"fmt"
	"os"
	"polaris/db"
	"strings"

	"polaris/log"
	"polaris/pkg/alist"
	"polaris/pkg/storage"
	"polaris/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (s *Server) GetAllStorage(c *gin.Context) (interface{}, error) {
	data := s.db.GetAllStorage()
	return data, nil
}

func (s *Server) AddStorage(c *gin.Context) (interface{}, error) {
	var in db.StorageInfo
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind json")
	}
	utils.TrimFields(&in)

	if in.Implementation == "webdav" {
		//test webdav
		wd := in.ToWebDavSetting()
		st, err := storage.NewWebdavStorage(wd.URL, wd.User, wd.Password, in.TvPath, false, nil, nil)
		if err != nil {
			return nil, errors.Wrap(err, "new webdav")
		}
		fs, err := st.ReadDir(".")
		if err != nil {
			return nil, errors.Wrap(err, "test read")
		}
		for _, f := range fs {
			log.Infof("file name: %v", f.Name())
		}
	} else if in.Implementation == "alist" {
		cfg := in.ToAlistSetting()
		_, err := storage.NewAlist(&alist.Config{URL: cfg.URL, Username: cfg.User, Password: cfg.Password}, in.TvPath, nil, nil)
		if err != nil {
			return nil, errors.Wrap(err, "alist")
		}
	} else if in.Implementation == "local" {
		_, err := os.Stat(in.TvPath)
		if err != nil {
			return nil, err
		}
		if !strings.HasSuffix(in.TvPath, string(os.PathSeparator)) {
			in.TvPath = in.TvPath + string(os.PathSeparator)
		}
		_, err = os.Stat(in.MoviePath)
		if err != nil {
			return nil, err
		}
		if !strings.HasSuffix(in.MoviePath, string(os.PathSeparator)) {
			in.MoviePath = in.MoviePath + string(os.PathSeparator)
		}

	}
	log.Infof("received add storage input: %v", in)
	err := s.db.AddStorage(&in)
	return nil, err
}

func (s *Server) DeleteStorage(c *gin.Context) (interface{}, error) {
	ids := c.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		return nil, fmt.Errorf("id is not int: %v", ids)
	}
	err = s.db.DeleteStorage(id)
	return nil, err
}

func (s *Server) SuggestedSeriesFolderName(c *gin.Context) (interface{}, error) {
	ids := c.Param("tmdb_id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		return nil, fmt.Errorf("id is not int: %v", ids)
	}
	name, err := s.core.SuggestedSeriesFolderName(id)
	if err != nil {
		return nil, err
	}
	return gin.H{"name": name}, nil
}

func (s *Server) SuggestedMovieFolderName(c *gin.Context) (interface{}, error) {
	ids := c.Param("tmdb_id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		return nil, fmt.Errorf("id is not int: %v", ids)
	}
	name, err := s.core.SuggestedMovieFolderName(id)
	if err != nil {
		return nil, err
	}
	return gin.H{"name": name}, nil
}
