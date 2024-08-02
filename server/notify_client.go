package server

import (
	"polaris/ent"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (s *Server) GetAllNotificationClients(c *gin.Context) (interface{}, error) {
	return s.db.GetAllNotificationClients()
}

func (s *Server) GetNotificationClient(c *gin.Context) (interface{}, error) {
	ids := c.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		return nil, errors.Wrap(err, "convert")
	}
	return s.db.GetNotificationClient(id)
}

func (s *Server) DeleteNotificationClient(c *gin.Context) (interface{}, error) {
	ids := c.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		return nil, errors.Wrap(err, "convert")
	}
	return nil, s.db.DeleteNotificationClient(id)
}

func (s *Server) AddNotificationClient(c *gin.Context) (interface{}, error) {
	var in ent.NotificationClient
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "json")
	}

	err := s.db.AddNotificationClient(in.Name, in.Service, in.Settings, in.Enabled)
	if err != nil {
		return nil, errors.Wrap(err, "save db")
	}
	return nil, nil
}
