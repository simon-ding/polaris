package server

import (
	"fmt"
	"polaris/engine"
	"polaris/ent"
	"polaris/ent/episode"
	"polaris/ent/history"
	"polaris/log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Activity struct {
	*ent.History
	Progress       int     `json:"progress"`
	SeedRatio      float64 `json:"seed_ratio"`
	UploadProgress float64 `json:"upload_progress"`
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
			tasks := s.core.GetTasks()
			tasks.Range(func(id int, task *engine.Task) bool {
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
					if task.UploadProgresser != nil {
						a.UploadProgress = task.UploadProgresser()
					}
				}
				return true
			})

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
		if err := s.addTorrent2Blacklist(his); err != nil {
			return nil, errors.Errorf("add to blacklist: %v", err)
		} else {
			log.Infof("success add magnet link to blacklist: %v", his.Link)
		}
	}

	if err := s.core.RemoveTaskAndTorrent(his.ID); err != nil {
		return nil, errors.Wrap(err, "remove torrent")
	}
	if his.Status == history.StatusSeeding {
		//seeding, will mark as complete
		log.Infof("history is now seeding, will only mark history as success: (%d) %s", his.ID, his.SourceTitle)
		if err := s.db.SetHistoryStatus(his.ID, history.StatusSuccess); err!= nil {
			return nil, errors.Wrap(err, "set status")
		}
		return nil, nil
	}

	err := s.db.DeleteHistory(in.ID)
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}

	episodeIds := s.core.GetEpisodeIds(his)

	for _, id := range episodeIds {
		ep, err := s.db.GetEpisodeByID(id)
		if err != nil {
			log.Warnf("get episode (%d) error: %v", id, err)
			continue
		}
		if !s.db.IsEpisodeDownloadingOrDownloaded(id) && ep.Status != episode.StatusDownloaded {
			//没有正在下载中或者下载完成的任务，并且episode状态不是已经下载完成
			log.Debugf("set episode (%d) status to missing", id)
			s.db.SetEpisodeStatus(id, episode.StatusMissing)
		}
	}

	log.Infof("history record successful deleted: %v", his.SourceTitle)
	return nil, nil
}

func (s *Server) addTorrent2Blacklist(h *ent.History) error {
	return s.db.AddTorrent2Blacklist(h.Hash, h.SourceTitle, h.MediaID)
}

func (s *Server) GetAllBlacklistItems(c *gin.Context) (interface{}, error) {
	list, err := s.db.GetTorrentBlacklist()
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}
	return list, nil
}
func (s *Server) RemoveBlacklistItem(c *gin.Context) (interface{}, error) {
	id := c.Param("id")
	if id == "" {
		return nil, fmt.Errorf("id is empty")
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("id is not int: %v", id)
	}
	if err := s.db.DeleteTorrentBlacklist(idInt); err != nil {
		return nil, errors.Wrap(err, "db")
	}
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
