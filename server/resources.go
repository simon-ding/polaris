package server

import (
	"fmt"
	"polaris/ent"
	"polaris/ent/episode"
	"polaris/ent/history"
	"polaris/log"
	"polaris/pkg/utils"
	"polaris/server/core"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (s *Server) searchAndDownloadSeasonPackage(seriesId, seasonNum int) (*string, error) {
	trc, err := s.getDownloadClient()
	if err != nil {
		return nil, errors.Wrap(err, "connect transmission")
	}

	res, err := core.SearchSeasonPackage(s.db, seriesId, seasonNum, true)
	if err != nil {
		return nil, err
	}

	r1 := res[0]
	log.Infof("found resource to download: %+v", r1)

	downloadDir := s.db.GetDownloadDir()
	size := utils.AvailableSpace(downloadDir)
	if size < uint64(r1.Size) {
		log.Errorf("space available %v, space needed %v", size, r1.Size)
		return nil, errors.New("no enough space")
	}

	torrent, err := trc.Download(r1.Link, s.db.GetDownloadDir())
	if err != nil {
		return nil, errors.Wrap(err, "downloading")
	}
	torrent.Start()

	series := s.db.GetMediaDetails(seriesId)
	if series == nil {
		return nil, fmt.Errorf("no tv series of id %v", seriesId)
	}
	dir := fmt.Sprintf("%s/Season %02d/", series.TargetDir, seasonNum)

	history, err := s.db.SaveHistoryRecord(ent.History{
		MediaID:     seriesId,
		EpisodeID:   0,
		SourceTitle: r1.Name,
		TargetDir:   dir,
		Status:      history.StatusRunning,
		Size:        r1.Size,
		Saved:       torrent.Save(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "save record")
	}
	s.db.SetSeasonAllEpisodeStatus(seriesId, seasonNum, episode.StatusDownloading)

	s.tasks[history.ID] = &Task{Torrent: torrent}
	return &r1.Name, nil
}

func (s *Server) searchAndDownload(seriesId, seasonNum, episodeNum int) (*string, error) {
	trc, err := s.getDownloadClient()
	if err != nil {
		return nil, errors.Wrap(err, "connect transmission")
	}
	series := s.db.GetMediaDetails(seriesId)
	if series == nil {
		return nil, fmt.Errorf("no tv series of id %v", seriesId)
	}
	var ep *ent.Episode
	for _, e := range series.Episodes {
		if e.SeasonNumber == seasonNum && e.EpisodeNumber == episodeNum {
			ep = e
		}
	}
	if ep == nil {
		return nil, errors.Errorf("no episode of season %d episode %d", seasonNum, episodeNum)
	}

	res, err := core.SearchEpisode(s.db, seriesId, seasonNum, episodeNum, true)
	if err != nil {
		return nil, err
	}
	r1 := res[0]
	log.Infof("found resource to download: %+v", r1)
	torrent, err := trc.Download(r1.Link, s.db.GetDownloadDir())
	if err != nil {
		return nil, errors.Wrap(err, "downloading")
	}
	torrent.Start()

	dir := fmt.Sprintf("%s/Season %02d/", series.TargetDir, seasonNum)

	history, err := s.db.SaveHistoryRecord(ent.History{
		MediaID:     ep.MediaID,
		EpisodeID:   ep.ID,
		SourceTitle: r1.Name,
		TargetDir:   dir,
		Status:      history.StatusRunning,
		Size:        r1.Size,
		Saved:       torrent.Save(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "save record")
	}
	s.db.SetEpisodeStatus(ep.ID, episode.StatusDownloading)

	s.tasks[history.ID] = &Task{Torrent: torrent}

	log.Infof("success add %s to download task", r1.Name)
	return &r1.Name, nil
}

type searchAndDownloadIn struct {
	ID      int `json:"id" binding:"required"`
	Season  int `json:"season"`
	Episode int `json:"episode"`
}

func (s *Server) SearchAvailableEpisodeResource(c *gin.Context) (interface{}, error) {
	var in searchAndDownloadIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind json")
	}
	log.Infof("search episode resources link: %v", in)
	res, err := core.SearchEpisode(s.db, in.ID, in.Season, in.Episode, true)
	if err != nil {
		return nil, errors.Wrap(err, "search episode")
	}
	var searchResults []TorznabSearchResult
	for _, r := range res {
		searchResults = append(searchResults, TorznabSearchResult{
			Name:    r.Name,
			Size:    r.Size,
			Seeders: r.Seeders,
			Peers:   r.Peers,
			Link:    r.Link,
		})
	}
	if len(searchResults) == 0 {
		return nil, errors.New("no resource found")
	}
	return searchResults, nil
}

func (s *Server) SearchTvAndDownload(c *gin.Context) (interface{}, error) {
	var in searchAndDownloadIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind json")
	}
	log.Infof("search episode resources link: %v", in)
	var name string
	if in.Episode == 0 {
		log.Infof("season package search")
		//search season package
		name1, err := s.searchAndDownloadSeasonPackage(in.ID, in.Season)
		if err != nil {
			return nil, errors.Wrap(err, "download")
		}
		name = *name1
	} else {
		log.Infof("season episode search")
		name1, err := s.searchAndDownload(in.ID, in.Season, in.Episode)
		if err != nil {
			return nil, errors.Wrap(err, "download")
		}
		name = *name1
	}

	return gin.H{
		"name": name,
	}, nil
}

type TorznabSearchResult struct {
	Name    string `json:"name"`
	Size    int    `json:"size"`
	Link    string `json:"link"`
	Seeders int    `json:"seeders"`
	Peers   int    `json:"peers"`
	Source  string `json:"source"`
}

func (s *Server) SearchAvailableMovies(c *gin.Context) (interface{}, error) {
	ids := c.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		return nil, errors.Wrap(err, "convert")
	}

	res, err := core.SearchMovie(s.db, id, false)
	if err != nil {
		if err.Error() == "no resource found" {
			return []TorznabSearchResult{}, nil
		}
		return nil, err
	}

	var searchResults []TorznabSearchResult
	for _, r := range res {
		searchResults = append(searchResults, TorznabSearchResult{
			Name:    r.Name,
			Size:    r.Size,
			Seeders: r.Seeders,
			Peers:   r.Peers,
			Link:    r.Link,
			Source:  r.Source,
		})
	}
	if len(searchResults) == 0 {
		return []TorznabSearchResult{}, nil
	}
	return searchResults, nil
}

type downloadTorrentIn struct {
	MediaID int `json:"media_id" binding:"required"`
	TorznabSearchResult
}

func (s *Server) DownloadMovieTorrent(c *gin.Context) (interface{}, error) {
	var in downloadTorrentIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind json")
	}
	log.Infof("download torrent input: %+v", in)

	trc, err := s.getDownloadClient()
	if err != nil {
		return nil, errors.Wrap(err, "connect transmission")
	}
	media := s.db.GetMediaDetails(in.MediaID)
	if media == nil {
		return nil, fmt.Errorf("no tv series of id %v", in.MediaID)
	}

	torrent, err := trc.Download(in.Link, s.db.GetDownloadDir())
	if err != nil {
		return nil, errors.Wrap(err, "downloading")
	}
	torrent.Start()
	name := in.Name
	if name == "" {
		name = media.OriginalName
	}
	go func() {
		ep := media.Episodes[0]
		history, err := s.db.SaveHistoryRecord(ent.History{
			MediaID:     media.ID,
			EpisodeID:   ep.ID,
			SourceTitle: name,
			TargetDir:   "./",
			Status:      history.StatusRunning,
			Size:        in.Size,
			Saved:       torrent.Save(),
		})
		if err != nil {
			log.Errorf("save history error: %v", err)
		}

		s.tasks[history.ID] = &Task{Torrent: torrent}

		s.db.SetEpisodeStatus(ep.ID, episode.StatusDownloading)
	}()

	log.Infof("success add %s to download task", media.NameEn)
	return media.NameEn, nil

}
