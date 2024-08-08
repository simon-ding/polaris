package core

import (
	"bytes"
	"fmt"
	"path/filepath"
	"polaris/db"
	"polaris/ent/media"
	storage1 "polaris/ent/storage"
	"polaris/log"
	"polaris/pkg/notifier"
	"polaris/pkg/storage"
	"strings"

	"github.com/pkg/errors"
)

func (c *Client) writePlexmatch(seriesId int, episodeId int, targetDir, name string) error {

	if !c.plexmatchEnabled() {
		return nil
	}
	series, err := c.db.GetMedia(seriesId)
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

	//season plexmatch file
	ep, err := c.db.GetEpisodeByID(episodeId)
	if err != nil {
		return errors.Wrap(err, "query episode")
	}
	buff := bytes.Buffer{}
	seasonPlex := filepath.Join(targetDir, ".plexmatch")
	data, err := st.ReadFile(seasonPlex)
	if err != nil {
		log.Infof("read season plexmatch: %v", err)
	} else {
		buff.Write(data)
	}
	if strings.Contains(buff.String(), name) {
		log.Debugf("already write plex episode line: %v", name)
		return nil
	}
	buff.WriteString(fmt.Sprintf("\nep: %d: %s\n", ep.EpisodeNumber, name))
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
