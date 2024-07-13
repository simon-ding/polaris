package server

import (
	"polaris/ent"
	"polaris/log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Activity struct {
	*ent.History
	InBackgroud bool `json:"in_backgroud"`
}

func (s *Server) GetAllActivities(c *gin.Context) (interface{}, error) {
	his := s.db.GetHistories()
	var activities = make([]Activity, 0, len(his))
	for _, h := range his {
		a := Activity{
			History: h,
		}
		for id, task := range s.tasks {
			if h.ID == id && task.Processing {
				a.InBackgroud = true
			}	
		}
		activities = append(activities, a)
	}

	return activities, nil
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
