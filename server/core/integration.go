package core

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"polaris/db"
	"polaris/ent/media"
	storage1 "polaris/ent/storage"
	"polaris/log"
	"polaris/pkg/metadata"
	"polaris/pkg/notifier"
	"polaris/pkg/storage"
	"polaris/pkg/utils"
	"slices"
	"strings"
)

func (c *Client) writePlexmatch(historyId int) error {

	if !c.plexmatchEnabled() {
		return nil
	}

	his := c.db.GetHistory(historyId)

	series, err := c.db.GetMedia(his.MediaID)
	if err != nil {
		return err
	}
	if series.MediaType != media.MediaTypeTv {
		return nil
	}
	st, err := c.getStorage(series.StorageID, media.MediaTypeTv)
	if err != nil {
		return errors.Wrap(err, "get storage")
	}

	//series plexmatch file
	_, err = st.ReadFile(filepath.Join(series.TargetDir, ".plexmatch"))
	if err != nil {
		//create new
		buff := bytes.Buffer{}
		if series.ImdbID != "" {
			buff.WriteString(fmt.Sprintf("imdbid: %s\n", series.ImdbID))
		}
		buff.WriteString(fmt.Sprintf("tmdbid: %d\n", series.TmdbID))
		log.Warnf(".plexmatch file not found, create new one: %s", series.NameEn)
		if err := st.WriteFile(filepath.Join(series.TargetDir, ".plexmatch"), buff.Bytes()); err != nil {
			return errors.Wrap(err, "series plexmatch")
		}
	}

	buff := bytes.Buffer{}
	seasonPlex := filepath.Join(his.TargetDir, ".plexmatch")
	data, err := st.ReadFile(seasonPlex)
	if err != nil {
		log.Infof("read season plexmatch: %v", err)
	} else {
		buff.Write(data)
	}

	if his.EpisodeID > 0 {
		//single episode download
		ep, err := c.db.GetEpisodeByID(his.EpisodeID)
		if err != nil {
			return errors.Wrap(err, "query episode")
		}
		if strings.Contains(buff.String(), ep.TargetFile) {
			log.Debugf("already write plex episode line: %v", ep.TargetFile)
			return nil
		}
		buff.WriteString(fmt.Sprintf("\nep: %d: %s\n", ep.EpisodeNumber, ep.TargetFile))
	} else {
		seasonNum, err := utils.SeasonId(his.TargetDir)
		if err != nil {
			return errors.Wrap(err, "no season id")
		}
		allEpisodes, err := c.db.GetSeasonEpisodes(his.MediaID, seasonNum)
		if err != nil {
			return errors.Wrap(err, "query season episode")
		}
		for _, ep := range allEpisodes {
			if strings.Contains(buff.String(), ep.TargetFile) {
				log.Debugf("already write plex episode line: %v", ep.TargetFile)
				continue
			}
			buff.WriteString(fmt.Sprintf("\nep: %d: %s\n", ep.EpisodeNumber, ep.TargetFile))
		}

	}
	log.Infof("write season plexmatch file content: %s", buff.String())
	return st.WriteFile(seasonPlex, buff.Bytes())
}

func (c *Client) plexmatchEnabled() bool {
	return c.db.GetSetting(db.SettingPlexMatchEnabled) == "true"
}

func (c *Client) getStorage(storageId int, mediaType media.MediaType) (storage.Storage, error) {
	st := c.db.GetStorage(storageId)
	targetPath := st.TvPath
	if mediaType == media.MediaTypeMovie {
		targetPath = st.MoviePath
	}

	switch st.Implementation {
	case storage1.ImplementationLocal:

		storageImpl1, err := storage.NewLocalStorage(targetPath)
		if err != nil {
			return nil, errors.Wrap(err, "new local")
		}
		return storageImpl1, nil

	case storage1.ImplementationWebdav:
		ws := st.ToWebDavSetting()
		storageImpl1, err := storage.NewWebdavStorage(ws.URL, ws.User, ws.Password, targetPath, ws.ChangeFileHash == "true")
		if err != nil {
			return nil, errors.Wrap(err, "new webdav")
		}
		return storageImpl1, nil
	}
	return nil, errors.New("no storage found")
}

func (c *Client) sendMsg(msg string) {
	clients, err := c.db.GetAllNotificationClients2()
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

func (c *Client) findEpisodeFilesPreMoving(historyId int) error {
	his := c.db.GetHistory(historyId)

	isSingleEpisode := his.EpisodeID > 0
	downloadDir := c.db.GetDownloadDir()
	task := c.tasks[historyId]
	target := filepath.Join(downloadDir, task.Name())
	fi, err := os.Stat(target)
	if err != nil {
		return errors.Wrapf(err, "read dir %v", target)
	}
	if isSingleEpisode {
		if fi.IsDir() {
			//download single episode in dir
			//TODO
		} else {
			//is file
			if err := c.db.UpdateEpisodeTargetFile(his.EpisodeID, fi.Name()); err != nil {
				log.Errorf("writing downloaded file name to db error: %v", err)
			}
		}
	} else {
		if !fi.IsDir() {
			return fmt.Errorf("not season pack downloaded")
		}
		seasonNum, err := utils.SeasonId(his.TargetDir)
		if err != nil {
			return errors.Wrap(err, "no season id")
		}

		files, err := os.ReadDir(target)
		if err != nil {
			return err
		}
		for _, f := range files {
			if f.IsDir() { //want media file
				continue
			}
			excludedExt := []string{".txt", ".srt", ".ass", ".sub"}
			ext := filepath.Ext(f.Name())
			if slices.Contains(excludedExt, strings.ToLower(ext)) {
				continue
			}

			meta := metadata.ParseTv(f.Name())
			if meta.Episode > 0 {
				//episode exists
				ep, err := c.db.GetEpisode(his.MediaID, seasonNum, meta.Episode)
				if err != nil {
					return err
				}
				if err := c.db.UpdateEpisodeTargetFile(ep.ID, f.Name()); err != nil {
					return errors.Wrap(err, "update episode file")
				}
			}
		}
	}
	return nil
}
