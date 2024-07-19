package server

import (
	"path/filepath"
	"polaris/ent"
	"polaris/ent/episode"
	"polaris/ent/history"
	"polaris/ent/media"
	storage1 "polaris/ent/storage"
	"polaris/log"
	"polaris/pkg"
	"polaris/pkg/storage"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func (s *Server) scheduler() {
	s.mustAddCron("@every 1m", s.checkTasks)
	//s.mustAddCron("@every 1h", s.checkAllFiles)
	s.mustAddCron("@every 1h", s.downloadTvSeries)
	s.cron.Start()
}

func (s *Server) mustAddCron(spec string, cmd func()) {
	if err := s.cron.AddFunc(spec, cmd); err != nil {
		log.Errorf("add func error: %v", err)
		panic(err)
	}
}

func (s *Server) checkTasks() {
	log.Infof("begin check tasks...")
	for id, t := range s.tasks {
		if !t.Exists() {
			log.Infof("task no longer exists: %v", id)
			continue
		}
		log.Infof("task (%s) percentage done: %d%%", t.Name(), t.Progress())
		if t.Progress() == 100 {
			log.Infof("task is done: %v", t.Name())
			go func() {
				if err := s.moveCompletedTask(id); err != nil {
					log.Infof("post tasks for id %v fail: %v", id, err)
				}
			}()
		}
	}
}

func (s *Server) moveCompletedTask(id int) (err error) {
	torrent := s.tasks[id]
	r := s.db.GetHistory(id)
	if r.Status == history.StatusUploading {
		log.Infof("task %d is already uploading, skip", id)
		return nil
	}
	s.db.SetHistoryStatus(r.ID, history.StatusUploading)

	defer func() {
		if err != nil {
			s.db.SetHistoryStatus(r.ID, history.StatusFail)
			if r.EpisodeID != 0 {
				s.db.SetEpisodeStatus(r.EpisodeID, episode.StatusMissing)
			}

		} else {
			delete(s.tasks, r.ID)
			s.db.SetHistoryStatus(r.ID, history.StatusSuccess)
			if r.EpisodeID != 0 {
				s.db.SetEpisodeStatus(r.EpisodeID, episode.StatusDownloaded)
			}

			torrent.Remove()
		}
	}()

	series := s.db.GetMediaDetails(r.MediaID)
	if series == nil {
		return nil
	}
	st := s.db.GetStorage(series.StorageID)
	log.Infof("move task files to target dir: %v", r.TargetDir)
	var stImpl storage.Storage
	if st.Implementation == storage1.ImplementationWebdav {
		ws := st.ToWebDavSetting()
		targetPath := ws.TvPath
		if series.MediaType == media.MediaTypeMovie {
			targetPath = ws.MoviePath
		}
		storageImpl, err := storage.NewWebdavStorage(ws.URL, ws.User, ws.Password, targetPath)
		if err != nil {
			return errors.Wrap(err, "new webdav")
		}
		stImpl = storageImpl

	} else if st.Implementation == storage1.ImplementationLocal {
		ls := st.ToLocalSetting()
		targetPath := ls.TvPath
		if series.MediaType == media.MediaTypeMovie {
			targetPath = ls.MoviePath
		}

		storageImpl, err := storage.NewLocalStorage(targetPath)
		if err != nil {
			return errors.Wrap(err, "new storage")
		}
		stImpl = storageImpl

	}
	if err := stImpl.Move(filepath.Join(s.db.GetDownloadDir(), torrent.Name()), r.TargetDir); err != nil {
		return errors.Wrap(err, "move file")
	}

	log.Infof("move downloaded files to target dir success, file: %v, target dir: %v", torrent.Name(), r.TargetDir)
	return nil
}

func (s *Server) checkDownloadedSeriesFiles(m *ent.Media) error {
	if m.MediaType != media.MediaTypeTv {
		return nil
	}
	log.Infof("check files in directory: %s", m.TargetDir)
	st := s.db.GetStorage(m.StorageID)

	var storageImpl storage.Storage

	switch st.Implementation {
	case storage1.ImplementationLocal:
		ls := st.ToLocalSetting()
		targetPath := ls.TvPath
		storageImpl1, err := storage.NewLocalStorage(targetPath)
		if err != nil {
			return errors.Wrap(err, "new local")
		}
		storageImpl = storageImpl1

	case storage1.ImplementationWebdav:
		ws := st.ToWebDavSetting()
		targetPath := ws.TvPath
		storageImpl1, err := storage.NewWebdavStorage(ws.URL, ws.User, ws.Password, targetPath)
		if err != nil {
			return errors.Wrap(err, "new webdav")
		}
		storageImpl = storageImpl1
	}
	files, err := storageImpl.ReadDir(m.TargetDir)
	if err != nil {
		return errors.Wrapf(err, "read dir %s", m.TargetDir)
	}
	seRe := regexp.MustCompile(`S\d+`)
	epRe := regexp.MustCompile(`E\d+`)
	for _, in := range files {
		if !in.IsDir() { //season dir, ignore file
			continue
		}
		dir := filepath.Join(m.TargetDir, in.Name())
		epFiles, err := storageImpl.ReadDir(dir)
		if err != nil {
			log.Errorf("read dir %s error: %v", dir, err)
			continue
		}
		for _, ep := range epFiles {
			log.Infof("found file: %v", ep.Name())
			matchEp := epRe.FindAllString(ep.Name(), -1)
			if len(matchEp) == 0 {
				continue
			}
			matchSe := seRe.FindAllString(ep.Name(), -1)
			if len(matchSe) == 0 {
				continue
			}

			epNum := strings.TrimPrefix(matchEp[0], "E")
			epNum1, _ := strconv.Atoi(epNum)
			seNum := strings.TrimPrefix(matchSe[0], "S")
			seNum1, _ := strconv.Atoi(seNum)
			var dirname = filepath.Join(in.Name(), ep.Name())
			log.Infof("found match, season num %d, episode num %d", seNum1, epNum1)
			err := s.db.UpdateEpisodeFile(m.ID, seNum1, epNum1, dirname)
			if err != nil {
				log.Error("update episode: %v", err)
			}
		}
	}
	return nil

}

type Task struct {
	//Processing bool
	pkg.Torrent
}

func (s *Server) downloadTvSeries() {
	log.Infof("begin check all tv series resources")
	allSeries := s.db.GetMediaWatchlist(media.MediaTypeTv)
	for _, series := range allSeries {
		detail, err := s.MustTMDB().GetTvDetails(series.TmdbID, s.language)
		if err != nil {
			log.Errorf("get tv details error: %v", err)
			continue
		}

		lastEpisode, err := s.db.GetEpisode(series.ID, detail.LastEpisodeToAir.SeasonNumber, detail.LastEpisodeToAir.EpisodeNumber)
		if err != nil {
			log.Errorf("get last episode error: %v", err)
			continue
		}
		if lastEpisode.Title != detail.LastEpisodeToAir.Name {
			s.db.UpdateEpiode(lastEpisode.ID, detail.LastEpisodeToAir.Name, detail.LastEpisodeToAir.Overview)
		}
		if lastEpisode.Status  == episode.StatusMissing {
			name, err := s.searchAndDownload(series.ID, lastEpisode.SeasonNumber, lastEpisode.EpisodeNumber)
			if err != nil {
				log.Infof("cannot find resource to download for %s: %v", lastEpisode.Title, err)
			} else {
				log.Infof("begin download torrent resource: %v",name)
			}
		}

		nextEpisode, err := s.db.GetEpisode(series.ID, detail.NextEpisodeToAir.SeasonNumber, detail.NextEpisodeToAir.EpisodeNumber)
		if err == nil {
			if nextEpisode.Title != detail.NextEpisodeToAir.Name {
				s.db.UpdateEpiode(nextEpisode.ID, detail.NextEpisodeToAir.Name, detail.NextEpisodeToAir.Overview)
				log.Errorf("updated next episode name to %v", detail.NextEpisodeToAir.Name)
			}
		}

	}
}
