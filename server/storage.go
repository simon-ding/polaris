package server

import (
	"fmt"
	"polaris/db"
	"polaris/log"
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
