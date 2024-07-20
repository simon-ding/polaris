package server

import (
	"fmt"
	"polaris/db"
	"polaris/ent"
	"polaris/ent/episode"
	"polaris/ent/history"
	"polaris/log"
	"polaris/pkg/transmission"
	"polaris/pkg/utils"
	"polaris/server/core"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type addTorznabIn struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	ApiKey string `json:"api_key"`
}

func (s *Server) AddTorznabInfo(c *gin.Context) (interface{}, error) {
	var in addTorznabIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind json")
	}
	err := s.db.SaveTorznabInfo(in.Name, db.TorznabSetting{
		URL:    in.URL,
		ApiKey: in.ApiKey,
	})
	if err != nil {
		return nil, errors.Wrap(err, "add ")
	}
	return nil, nil
}

func (s *Server) DeleteTorznabInfo(c *gin.Context) (interface{}, error) {
	var ids = c.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		return nil, fmt.Errorf("id is not correct: %v", ids)
	}
	s.db.DeleteTorznab(id)
	return "success", nil
}

func (s *Server) GetAllIndexers(c *gin.Context) (interface{}, error) {
	indexers := s.db.GetAllTorznabInfo()
	if len(indexers) == 0 {
		return nil, nil
	}
	return indexers, nil
}

func (s *Server) getDownloadClient() (*transmission.Client, error) {
	tr := s.db.GetTransmission()
	trc, err := transmission.NewClient(transmission.Config{
		URL:      tr.URL,
		User:     tr.User,
		Password: tr.Password,
	})
	if err != nil {
		return nil, errors.Wrap(err, "connect transmission")
	}
	return trc, nil
}

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
	log.Infof("found resource to download: %v", r1)

	downloadDir := s.db.GetDownloadDir()
	size := utils.AvailableSpace(downloadDir)
	if size < uint64(r1.Size) {
		log.Errorf("space available %v, space needed %v", size, r1.Size)
		return nil, errors.New("no enough space")
	}

	torrent, err := trc.Download(r1.Magnet, s.db.GetDownloadDir())
	if err != nil {
		return nil, errors.Wrap(err, "downloading")
	}
	torrent.Start()

	series := s.db.GetMediaDetails(seriesId)
	if series == nil {
		return nil, fmt.Errorf("no tv series of id %v", seriesId)
	}
	dir := fmt.Sprintf("%s/Season %02d", series.TargetDir, seasonNum)

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
	log.Infof("found resource to download: %v", r1)
	torrent, err := trc.Download(r1.Magnet, s.db.GetDownloadDir())
	if err != nil {
		return nil, errors.Wrap(err, "downloading")
	}
	torrent.Start()

	dir := fmt.Sprintf("%s/Season %02d", series.TargetDir, seasonNum)

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
			Link:    r.Magnet,
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
}

func (s *Server) SearchAvailableMovies(c *gin.Context) (interface{}, error) {
	ids := c.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		return nil, errors.Wrap(err, "convert")
	}

	movieDetail := s.db.GetMediaDetails(id)
	if movieDetail == nil {
		return nil, errors.New("no media found of id " + ids)
	}

	res, err := core.SearchMovie(s.db, id, false)
	if err != nil {
		return nil, err
	}

	var searchResults []TorznabSearchResult
	for _, r := range res {
		searchResults = append(searchResults, TorznabSearchResult{
			Name:    r.Name,
			Size:    r.Size,
			Seeders: r.Seeders,
			Peers:   r.Peers,
			Link:    r.Magnet,
		})
	}
	if len(searchResults) == 0 {
		return nil, errors.New("no resource found")
	}
	return searchResults, nil
}

type downloadTorrentIn struct {
	MediaID int    `json:"media_id" binding:"required"`
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

	go func() {
		ep := media.Episodes[0]
		history, err := s.db.SaveHistoryRecord(ent.History{
			MediaID:     media.ID,
			EpisodeID:   ep.ID,
			SourceTitle: media.NameCn,
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

type downloadClientIn struct {
	Name           string `json:"name"`
	URL            string `json:"url"`
	User           string `json:"user"`
	Password       string `json:"password"`
	Implementation string `json:"implementation"`
}

func (s *Server) AddDownloadClient(c *gin.Context) (interface{}, error) {
	var in downloadClientIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind json")
	}
	if err := s.db.SaveTransmission(in.Name, in.URL, in.User, in.Password); err != nil {
		return nil, errors.Wrap(err, "save transmission")
	}
	return nil, nil
}

func (s *Server) GetAllDonloadClients(c *gin.Context) (interface{}, error) {
	res := s.db.GetAllDonloadClients()
	if len(res) == 0 {
		return nil, nil
	}
	return res, nil
}

func (s *Server) DeleteDownloadCLient(c *gin.Context) (interface{}, error) {
	var ids = c.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		return nil, fmt.Errorf("id is not correct: %v", ids)
	}
	s.db.DeleteDownloadCLient(id)
	return "success", nil
}
