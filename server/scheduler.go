package server

import (
	"path/filepath"
	"polaris/ent"
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
	s.mustAddCron("@every 10m", s.checkAllFiles)
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
		if t.Processing {
			continue
		}

		log.Infof("task (%s) percentage done: %d%%", t.Name(), t.Progress())
		if t.Progress() == 100 {
			log.Infof("task is done: %v", t.Name())
			t.Processing = true
			go func() {
				if err := s.moveCompletedTask(id); err != nil {
					log.Infof("post tasks for id %v fail: %v", id, err)
				}
			}()
		}
	}
}

func (s *Server) moveCompletedTask(id int) error {
	torrent := s.tasks[id]
	r := s.db.GetHistory(id)

	series := s.db.GetSeriesDetails(r.SeriesID)
	if series == nil {
		return nil
	}
	st := s.db.GetStorage(series.StorageID)
	log.Infof("move task files to target dir: %v", r.TargetDir)
	if st.Implementation == storage1.ImplementationWebdav {
		ws := st.ToWebDavSetting()
		storageImpl, err := storage.NewWebdavStorage(ws.Path, ws.User, ws.Password)
		if err != nil {
			return errors.Wrap(err, "new webdav")
		}
		if err := storageImpl.Move(filepath.Join(s.db.GetDownloadDir(), torrent.Name()), r.TargetDir); err != nil {
			return errors.Wrap(err, "move webdav")
		}
	} else if st.Implementation == storage1.ImplementationLocal {
		ls := st.ToLocalSetting()
		storageImpl, err := storage.NewLocalStorage(ls.Path)
		if err != nil {
			return errors.Wrap(err, "new storage")
		}

		if err := storageImpl.Move(filepath.Join(s.db.GetDownloadDir(), torrent.Name()), r.TargetDir); err != nil {
			return errors.Wrap(err, "move webdav")
		}

	}
	log.Infof("move downloaded files to target dir success, file: %v, target dir: %v", torrent.Name(), r.TargetDir)
	torrent.Remove()
	delete(s.tasks, r.ID)
	s.db.SetHistoryComplete(r.ID)
	return nil
}

func (s *Server) updateSeriesEpisodes(seriesId int) {

}

func (s *Server) checkAllFiles() {
	var tvs = s.db.GetWatchlist()
	for _, se := range tvs {
		if err := s.checkFileExists(se); err != nil {
			log.Errorf("check files for %s error: %v", se.NameCn, err)
		}
	}
}

func (s *Server) checkFileExists(series *ent.Series) error{
	log.Infof("check files in directory: %s", series.TargetDir)
	st := s.db.GetStorage(series.StorageID)
	var storageImpl storage.Storage

	switch st.Implementation {
	case storage1.ImplementationLocal:
		ls := st.ToLocalSetting()
		storageImpl1, err := storage.NewLocalStorage(ls.Path)
		if err != nil {
			return errors.Wrap(err, "new local")
		}
		storageImpl = storageImpl1

	case storage1.ImplementationWebdav:
		ws := st.ToWebDavSetting()
		storageImpl1, err := storage.NewWebdavStorage(ws.Path, ws.User, ws.Password)
		if err != nil {
			return errors.Wrap(err, "new webdav")
		}
		storageImpl = storageImpl1
	} 
	files, err := storageImpl.ReadDir(series.TargetDir)
	if err != nil {
		return errors.Wrapf(err, "read dir %s", series.TargetDir)
	}
	numRe := regexp.MustCompile("[0-9]+")
	epRe := regexp.MustCompile("E[0-9]+")
	for _, in := range files {
		if !in.IsDir() {//season dir, ignore file
			continue
		}
		nums := numRe.FindAllString(in.Name(), -1)
		if len(nums) == 0 {
			continue
		}
		seasonNum := nums[0]
		seasonNum1, _ := strconv.Atoi(seasonNum)
		dir := filepath.Join(series.TargetDir, in.Name())
		epFiles, err := storageImpl.ReadDir(dir)
		if err != nil {
			log.Errorf("read dir %s error: %v", dir, err)
			continue
		}
		for _, ep := range epFiles {
			match := epRe.FindAllString(ep.Name(), -1)
			if len(match) == 0 {
				continue
			}
			epNum := strings.TrimPrefix(match[0], "E")
			epNum1, _ := strconv.Atoi(epNum)
			err := s.db.UpdateEpisodeFile(series.ID, seasonNum1, epNum1, filepath.Join(in.Name(), ep.Name()))
			if err != nil {
				log.Error("update episode: %v", err)
			}
		}
	}
	return nil

}

type Task struct {
	Processing bool
	pkg.Torrent
}