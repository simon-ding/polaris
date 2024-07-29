package server

import (
	"polaris/ent"
	"polaris/log"
	"polaris/pkg/notifier"
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

func (s *Server) sendMsg(msg string)  {
	clients, err := s.db.GetAllNotificationClients2()
	if err != nil {
		log.Errorf("query notification clients: %v", err)
		return 
	}
	for _, cl := range clients {
		if !cl.Enabled {
			continue
		}
		handler, ok := notifier.Gethandler(cl.Service)
		if !ok {
			log.Errorf("no notification implementation of service %s", cl.Service)
			continue
		}
		noCl, err := handler(cl.Settings)
		if err != nil {
			log.Errorf("handle setting for name %s error: %v", cl.Name, err)
			continue
		}
		err = noCl.SendMsg(msg)
		if err != nil {
			log.Errorf("send message error: %v", err)
			continue
		}
		log.Debugf("send message to %s success, msg is %s", cl.Name, msg)
	}
}