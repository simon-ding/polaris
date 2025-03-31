package engine

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/fs"
	"path/filepath"
	"polaris/db"
	"polaris/ent/media"
	storage1 "polaris/ent/storage"
	"polaris/log"
	"polaris/pkg/alist"
	"polaris/pkg/metadata"
	"polaris/pkg/notifier"
	"polaris/pkg/storage"
	"slices"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func (c *Engine) writeNfoFile(historyId int) error {
	if !c.nfoSupportEnabled() {
		return nil
	}

	his := c.db.GetHistory(historyId)

	md, err := c.db.GetMedia(his.MediaID)
	if err != nil {
		return err
	}

	if md.MediaType == media.MediaTypeTv { //tvshow.nfo
		st, err := c.GetStorage(md.StorageID, media.MediaTypeTv)
		if err != nil {
			return errors.Wrap(err, "get storage")
		}

		nfoPath := filepath.Join(md.TargetDir, "tvshow.nfo")
		_, err = st.ReadFile(nfoPath)
		if err != nil {
			log.Infof("tvshow.nfo file missing, create new one, tv series name: %s", md.NameEn)
			show := Tvshow{
				Title:         md.NameCn,
				Originaltitle: md.OriginalName,
				Showtitle:     md.NameCn,
				Plot:          md.Overview,
				ID:            strconv.Itoa(md.TmdbID),
				Uniqueid: []UniqueId{
					{
						Text:    strconv.Itoa(md.TmdbID),
						Type:    "tmdb",
						Default: "true",
					},
					{
						Text: md.ImdbID,
						Type: "imdb",
					},
				},
			}
			data, err := xml.MarshalIndent(&show, " ", "  ")
			if err != nil {
				return errors.Wrap(err, "xml marshal")
			}
			return st.WriteFile(nfoPath, []byte(xml.Header+string(data)))
		}

	} else if md.MediaType == media.MediaTypeMovie { //movie.nfo
		st, err := c.GetStorage(md.StorageID, media.MediaTypeMovie)
		if err != nil {
			return errors.Wrap(err, "get storage")
		}

		nfoPath := filepath.Join(md.TargetDir, "movie.nfo")
		_, err = st.ReadFile(nfoPath)
		if err != nil {
			log.Infof("movie.nfo file missing, create new one, tv series name: %s", md.NameEn)
			nfoData := Movie{
				Title:         md.NameCn,
				Originaltitle: md.OriginalName,
				Sorttitle:     md.NameCn,
				Plot:          md.Overview,
				ID:            strconv.Itoa(md.TmdbID),
				Uniqueid: []UniqueId{
					{
						Text:    strconv.Itoa(md.TmdbID),
						Type:    "tmdb",
						Default: "true",
					},
					{
						Text: md.ImdbID,
						Type: "imdb",
					},
				},
			}
			data, err := xml.MarshalIndent(&nfoData, " ", "  ")
			if err != nil {
				return errors.Wrap(err, "xml marshal")
			}
			return st.WriteFile(nfoPath, []byte(xml.Header+string(data)))
		}
	}
	return nil
}

func (c *Engine) writePlexmatch(historyId int) error {

	if !c.plexmatchEnabled() {
		return nil
	}

	his := c.db.GetHistory(historyId)

	series, err := c.db.GetMedia(his.MediaID)
	if err != nil {
		return err
	}
	if series.MediaType != media.MediaTypeTv { //.plexmatch only support tv series
		return nil
	}
	st, err := c.GetStorage(series.StorageID, media.MediaTypeTv)
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
	episodesIds := c.GetEpisodeIds(his)

	for _, id := range episodesIds {
		ep, err := c.db.GetEpisodeByID(id)
		if err != nil {
			log.Warnf("query episode: %v", err)
			continue
		}
		if strings.Contains(buff.String(), ep.TargetFile) {
			log.Debugf("already write plex episode line: %v", ep.TargetFile)
			return nil
		}
		buff.WriteString(fmt.Sprintf("\nep: %d: %s\n", ep.EpisodeNumber, ep.TargetFile))

	}

	log.Infof("write season plexmatch file content: %s", buff.String())
	return st.WriteFile(seasonPlex, buff.Bytes())
}

func (c *Engine) plexmatchEnabled() bool {
	return c.db.GetSetting(db.SettingPlexMatchEnabled) == "true"
}

func (c *Engine) nfoSupportEnabled() bool {
	return c.db.GetSetting(db.SettingNfoSupportEnabled) == "true"
}

func (c *Engine) GetStorage(storageId int, mediaType media.MediaType) (storage.Storage, error) {
	st := c.db.GetStorage(storageId)
	targetPath := st.TvPath
	if mediaType == media.MediaTypeMovie {
		targetPath = st.MoviePath
	}
	videoFormats, err := c.db.GetAcceptedVideoFormats()
	if err != nil {
		log.Warnf("get accepted video format error: %v", err)
	}
	subtitleFormats, err := c.db.GetAcceptedSubtitleFormats()
	if err != nil {
		log.Warnf("get accepted subtitle format error: %v", err)
	}

	switch st.Implementation {
	case storage1.ImplementationLocal:

		storageImpl1, err := storage.NewLocalStorage(targetPath, videoFormats, subtitleFormats)
		if err != nil {
			return nil, errors.Wrap(err, "new local")
		}
		return storageImpl1, nil

	case storage1.ImplementationWebdav:
		ws := st.ToWebDavSetting()
		storageImpl1, err := storage.NewWebdavStorage(ws.URL, ws.User, ws.Password, targetPath, ws.ChangeFileHash == "true", videoFormats, subtitleFormats)
		if err != nil {
			return nil, errors.Wrap(err, "new webdav")
		}
		return storageImpl1, nil
	case storage1.ImplementationAlist:
		cfg := st.ToWebDavSetting()
		storageImpl1, err := storage.NewAlist(&alist.Config{URL: cfg.URL, Username: cfg.User, Password: cfg.Password}, targetPath, videoFormats, subtitleFormats)
		if err != nil {
			return nil, errors.Wrap(err, "alist")
		}
		return storageImpl1, nil
	}
	return nil, errors.New("no storage found")
}

func (c *Engine) sendMsg(msg string) {
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

func (c *Engine) findEpisodeFilesPreMoving(historyId int) error {
	his := c.db.GetHistory(historyId)

	episodeIds := c.GetEpisodeIds(his)

	task := c.tasks[historyId]

	ff, err := c.db.GetAcceptedVideoFormats()
	if err != nil {
		return err
	}
	for _, id := range episodeIds {
		ep, _ := c.db.GetEpisode(his.MediaID, his.SeasonNum, id)
		task.WalkFunc()(func(path string, info fs.FileInfo) error {
			if info.IsDir() {
				return nil
			}
			ext := filepath.Ext(info.Name())
			if slices.Contains(ff, ext) {
				return nil
			}
			meta := metadata.ParseTv(info.Name())
			if meta.StartEpisode == meta.EndEpisode && meta.StartEpisode == ep.EpisodeNumber {
				if err := c.db.UpdateEpisodeTargetFile(id, info.Name()); err != nil {
					log.Errorf("writing downloaded file name to db error: %v", err)
				}

			}
			return nil
		})
	}
	return nil
}
