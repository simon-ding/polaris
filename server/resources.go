package server

import (
	"fmt"
	"polaris/db"
	"polaris/ent"
	"polaris/ent/episode"
	"polaris/ent/history"
	"polaris/log"
	"polaris/pkg/torznab"
	"polaris/pkg/transmission"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (s *Server) searchWithTorznab(q string) []torznab.Result {

	var res []torznab.Result
	allTorznab := s.db.GetAllTorznabInfo()
	for _, tor := range allTorznab {
		resp, err := torznab.Search(tor.URL, tor.ApiKey, q)
		if err != nil {
			log.Errorf("search %s error: %v", tor.Name, err)
			continue
		}
		res = append(res, resp...)
	}
	return res
}

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

	q := fmt.Sprintf("%s S%02dE%02d", series.OriginalName, seasonNum, episodeNum)

	res := s.searchWithTorznab(q)
	if len(res) == 0 {
		return nil, fmt.Errorf("no resource found")
	}
	r1 := s.findBestMatchTv(res, seasonNum, episodeNum, series)
	log.Infof("found resource to download: %v", r1)
	torrent, err := trc.Download(r1.Magnet, s.db.GetDownloadDir())
	if err != nil {
		return nil, errors.Wrap(err, "downloading")
	}
	torrent.Start()

	dir := fmt.Sprintf("%s/Season %02d", series.TargetDir, ep.SeasonNumber)

	history, err := s.db.SaveHistoryRecord(ent.History{
		MediaID:     ep.MediaID,
		EpisodeID:   ep.ID,
		SourceTitle: r1.Name,
		TargetDir:   dir,
		Status:      history.StatusRunning,
		Size:        r1.Size,
		Saved:       torrent.Save(),
	})
	s.db.SetEpisodeStatus(ep.ID, episode.StatusDownloading)
	if err != nil {
		return nil, errors.Wrap(err, "save record")
	}
	s.tasks[history.ID] = &Task{Torrent: torrent}

	log.Infof("success add %s to download task", r1.Name)
	return &r1.Name, nil
}

func (s *Server) findBestMatchTv(resources []torznab.Result, season, episode int, series *db.MediaDetails) torznab.Result {
	var filtered []torznab.Result
	for _, r := range resources {
		if !(series.NameEn != "" && strings.Contains(r.Name, series.NameEn)) && !strings.Contains(r.Name, series.OriginalName) {
			//name not match
			continue
		}

		se := fmt.Sprintf("S%02dE%02d", season, episode)
		if !strings.Contains(r.Name, se) {
			//season or episode not match
			continue
		}
		if !strings.Contains(strings.ToLower(r.Name), series.Resolution) {
			//resolution not match
			continue
		}
		filtered = append(filtered, r)
	}

	sort.Slice(filtered, func(i, j int) bool {
		var s1 = filtered[i]
		var s2 = filtered[2]
		return s1.Seeders > s2.Seeders
	})

	return filtered[0]
}

type searchAndDownloadIn struct {
	ID      int `json:"id" binding:"required"`
	Season  int `json:"season"`
	Episode int `json:"episode"`
}

func (s *Server) SearchTvAndDownload(c *gin.Context) (interface{}, error) {
	var in searchAndDownloadIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind json")
	}
	log.Infof("search episode resources link: %v", in)
	name, err := s.searchAndDownload(in.ID, in.Season, in.Episode)
	if err != nil {
		return nil, errors.Wrap(err, "download")
	}

	return gin.H{
		"name": *name,
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

	res := s.searchWithTorznab(movieDetail.NameEn)

	res1 := s.searchWithTorznab(movieDetail.NameCn)
	res = append(res, res1...)

	if len(res) == 0 {
		return nil, fmt.Errorf("no resource found")
	}
	ss := strings.Split(movieDetail.AirDate, "-")[0]
	year, _ := strconv.Atoi(ss)

	var searchResults []TorznabSearchResult
	for _, r := range res {
		if !strings.Contains(r.Name, strconv.Itoa(year)) && !strings.Contains(r.Name, strconv.Itoa(year+1)) && !strings.Contains(r.Name, strconv.Itoa(year-1)) {
			continue //not the same movie, if year is not correct
		}
		if !strings.Contains(r.Name, movieDetail.NameCn) && !strings.Contains(r.Name, movieDetail.NameEn) {
			continue //name not match
		}
		searchResults = append(searchResults, TorznabSearchResult{
			Name: r.Name,
			Size: r.Size,
			Seeders: r.Seeders,
			Peers: r.Peers,
			Link: r.Magnet,
		})
	}

	return searchResults, nil

}

type downloadTorrentIn struct {
	MediaID int `json:"media_id" binding:"required"`
	Link string `json:"link" binding:"required"`
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

	go func ()  {
		for {
			if !torrent.Exists() {
				continue
			}
			history, err := s.db.SaveHistoryRecord(ent.History{
				MediaID:     media.ID,
				SourceTitle: torrent.Name(),
				TargetDir:   "./",
				Status:      history.StatusRunning,
				Size:        torrent.Size(),
				Saved:       torrent.Save(),
			})
			if err != nil {
				log.Errorf("save history error: %v", err)
			}

			s.tasks[history.ID] = &Task{Torrent: torrent}
		
			break
		}	
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
