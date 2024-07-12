package server

import (
	"path/filepath"
	"polaris/db"
	"polaris/log"
	"polaris/pkg/storage"

	"github.com/pkg/errors"
)

func (s *Server) scheduler() {
	s.mustAddCron("@every 1m", s.checkTasks)
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

func (s *Server) moveCompletedTask(id int) error {
	torrent := s.tasks[id]
	r := s.db.GetHistory(id)
	s.db.SetHistoryComplete(r.ID)

	delete(s.tasks, r.ID)

	series := s.db.GetSeriesDetails(r.SeriesID)
	st := s.db.GetStorage(series.StorageID)
	log.Infof("move task files to target dir: %v", r.TargetDir)
	if st.Implementation == db.ImplWebdav {
		storageImpl, err := storage.NewWebdavStorage(st.Path, st.User, st.Password)
		if err != nil {
			return errors.Wrap(err, "new webdav")
		}
		if err := storageImpl.Move(filepath.Join(s.db.GetDownloadDir(), torrent.Name()), r.TargetDir); err != nil {
			return errors.Wrap(err, "move webdav")
		}
	} else if st.Implementation == db.ImplLocal {
		storageImpl := storage.NewLocalStorage(st.Path)

		if err := storageImpl.Move(filepath.Join(s.db.GetDownloadDir(), torrent.Name()), r.TargetDir); err != nil {
			return errors.Wrap(err, "move webdav")
		}

	}
	log.Infof("move downloaded files to target dir success, file: %v, target dir: %v", torrent.Name(), r.TargetDir)
	torrent.Remove()
	return nil
}
