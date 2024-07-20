package server

import (
	"polaris/ent"
	"polaris/ent/episode"
	"polaris/log"
	"polaris/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Activity struct {
	*ent.History
	Progress int `json:"progress"`
}

func (s *Server) GetAllActivities(c *gin.Context) (interface{}, error) {
	his := s.db.GetHistories()
	var activities = make([]Activity, 0, len(his))
	for _, h := range his {
		a := Activity{
			History: h,
		}
		for id, task := range s.tasks {
			if h.ID == id && task.Exists() {
				a.Progress = task.Progress()
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
	if his.EpisodeID != 0 {
		s.db.SetEpisodeStatus(his.EpisodeID, episode.StatusMissing)

	} else {
		seasonNum, err := utils.SeasonId(his.TargetDir)
		if err != nil {
			log.Errorf("no season id: %v", his.TargetDir)
			seasonNum = -1
		}
		s.db.SetSeasonAllEpisodeStatus(his.MediaID, seasonNum, episode.StatusMissing)

	}

	err = s.db.DeleteHistory(id)
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}
	log.Infof("history record successful deleted: %v", his.SourceTitle)
	return nil, nil
}
