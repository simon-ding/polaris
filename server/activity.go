package server

import (
	"fmt"
	"polaris/ent"
	"polaris/ent/blacklist"
	"polaris/ent/episode"
	"polaris/ent/history"
	"polaris/ent/schema"
	"polaris/log"
	"polaris/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Activity struct {
	*ent.History
	Progress  int     `json:"progress"`
	SeedRatio float64 `json:"seed_ratio"`
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
					p, err := task.Progress()
					if err != nil {
						log.Warnf("get task progress error: %v", err)
					} else {
						a.Progress = p
					}
					r, err := task.SeedRatio()
					if err != nil {
						log.Warnf("get task seed ratio error: %v", err)
					} else {
						a.SeedRatio = r
					}
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

type removeActivityIn struct {
	ID            int  `json:"id"`
	Add2Blacklist bool `json:"add_2_blacklist"`
}

func (s *Server) RemoveActivity(c *gin.Context) (interface{}, error) {
	var in removeActivityIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind json")
	}

	his := s.db.GetHistory(in.ID)
	if his == nil {
		log.Errorf("no record of id: %d", in.ID)
		return nil, nil
	}
	if in.Add2Blacklist && his.Link != "" {
		//should add to blacklist
		if err := s.addTorrent2Blacklist(his.Link); err != nil {
			return nil, errors.Errorf("add to blacklist: %v", err)
		} else {
			log.Infof("success add magnet link to blacklist: %v", his.Link)
		}
	}

	if err := s.core.RemoveTaskAndTorrent(his.ID); err != nil {
		return nil, errors.Wrap(err, "remove torrent")
	}
	err := s.db.DeleteHistory(in.ID)
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}

	if his.EpisodeID != 0 {
		if !s.db.IsEpisodeDownloadingOrDownloaded(his.EpisodeID) {
			s.db.SetEpisodeStatus(his.EpisodeID, episode.StatusMissing)
		}

	} else {
		seasonNum, err := utils.SeasonId(his.TargetDir)
		if err != nil {
			log.Errorf("no season id: %v", his.TargetDir)
			seasonNum = -1
		}
		if his.Status == history.StatusRunning || his.Status == history.StatusUploading {
			s.db.SetSeasonAllEpisodeStatus(his.MediaID, seasonNum, episode.StatusMissing)
		}
	}

	log.Infof("history record successful deleted: %v", his.SourceTitle)
	return nil, nil
}

func (s *Server) addTorrent2Blacklist(link string) error {
	if link == "" {
		return nil
	}
	if hash, err := utils.MagnetHash(link); err != nil {
		return err
	} else {
		item := ent.Blacklist{
			Type: blacklist.TypeTorrent,
			Value: schema.BlacklistValue{
				TorrentHash: hash,
			},
		}
		err := s.db.AddBlacklistItem(&item)
		if err != nil {
			return errors.Wrap(err, "add to db")
		}
	}
	return nil
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
	ID        string  `json:"id"`
	SeedRatio float32 `json:"seed_ratio"`
	Progress  int     `json:"progress"`
}

func (s *Server) GetAllTorrents(c *gin.Context) (interface{}, error) {
	trc, _, err := s.core.GetDownloadClient()
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
		name, _ := t.Name()
		p, _ := t.Progress()
		infos = append(infos, TorrentInfo{
			Name:     name,
			ID:       t.GetHash(),
			Progress: p,
		})
	}
	return infos, nil
}
