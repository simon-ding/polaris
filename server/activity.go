package server

import (
	"fmt"
	"polaris/ent"
	"polaris/ent/episode"
	"polaris/ent/history"
	"polaris/log"
	"polaris/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Activity struct {
	*ent.History
	Progress  int     `json:"progress"`
	SeedRatio float32 `json:"seed_ratio"`
}

func (s *Server) GetAllActivities(c *gin.Context) (interface{}, error) {
	q := c.Query("status")
	var activities = make([]Activity, 0)
	if q == "active" {
		his := s.db.GetRunningHistories()
		for _, h := range his {
			a := Activity{
				History: h,
			}
			for id, task := range s.core.GetTasks() {
				if h.ID == id && task.Exists() {
					a.Progress = task.Progress()
					a.SeedRatio = float32(*task.SeedRatio())
				}
			}
			activities = append(activities, a)
		}
	} else {
		his := s.db.GetHistories()
		for _, h := range his {
			if h.Status == history.StatusRunning || h.Status == history.StatusUploading || h.Status == history.StatusSeeding {
				continue //archived downloads
			}

			a := Activity{
				History: h,
			}
			activities = append(activities, a)
		}

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

	if err := s.core.RemoveTaskAndTorrent(his.ID); err != nil {
		return nil, errors.Wrap(err, "remove torrent")
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
func (s *Server) GetMediaDownloadHistory(c *gin.Context) (interface{}, error) {
	var ids = c.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		return nil, fmt.Errorf("id is not correct: %v", ids)
	}
	his, err := s.db.GetDownloadHistory(id)
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}
	return his, nil
}

type TorrentInfo struct {
	Name      string  `json:"name"`
	ID        int64   `json:"id"`
	SeedRatio float32 `json:"seed_ratio"`
	Progress  int     `json:"progress"`
}

func (s *Server) GetAllTorrents(c *gin.Context) (interface{}, error) {
	trc, _, err := s.getDownloadClient()
	if err != nil {
		return nil, errors.Wrap(err, "connect transmission")
	}
	all, err := trc.GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "get all")
	}
	var infos []TorrentInfo
	for _, t := range all {
		if !t.Exists() {
			continue
		}
		infos = append(infos, TorrentInfo{
			Name:     t.Name(),
			ID:       t.ID,
			Progress: t.Progress(),
		})
	}
	return infos, nil
}
