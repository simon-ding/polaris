package server

import (
	"polaris/log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (s *Server) GetAllActivities(c *gin.Context) (interface{}, error) {
	his := s.db.GetHistories()

	return his, nil
}

func (s *Server) RemoveActivity(c *gin.Context) (interface{}, error) {
	ids := c.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		return nil, errors.Wrap(err, "convert")
	}
	his := s.db.GetHistory(id)
	if his == nil {
		log.Errorf("no record of id: %d", id)
		return nil, nil
	}
	torrent := s.tasks[his.ID]
	if torrent != nil {
		if err := torrent.Remove(); err != nil {
			return nil, errors.Wrap(err, "remove torrent")
		}
		delete(s.tasks, his.ID)
	}
	
	err = s.db.DeleteHistory(id)
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}
	log.Infof("history record successful deleted: %v", his.SourceTitle)
	return nil, nil
}
