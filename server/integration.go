package server

import (
	"fmt"
	"path/filepath"
	"polaris/db"
	"polaris/ent/media"
	"polaris/log"

	"github.com/pkg/errors"
)

func (s *Server) createPlexmatchIfNotExists(seriesId int) error {

	if !s.plexmatchEnabled() {
		return nil
	}
	series, err := s.db.GetMedia(seriesId)
	if err != nil {
		return err
	}
	if series.MediaType != media.MediaTypeTv {
		return nil
	}
	st, err := s.getStorage(series.StorageID, media.MediaTypeTv)
	if err != nil {
		return errors.Wrap(err, "get storage")
	}

	_, err = st.ReadFile(filepath.Join(series.TargetDir, ".plexmatch"))
	if err != nil { 
		//create new
		log.Warnf(".plexmatch file not found, create new one: %s", series.NameEn)
		return st.WriteFile(filepath.Join(series.TargetDir, ".plexmatch"), []byte(fmt.Sprintf("tmdbid: %d\n",series.TmdbID)))
	} 
	return nil
}

func (s *Server) plexmatchEnabled() bool {
	return s.db.GetSetting(db.SettingPlexMatchEnabled) == "true"
}
