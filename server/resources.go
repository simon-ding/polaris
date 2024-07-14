package server

import (
	"fmt"
	"polaris/db"
	"polaris/ent"
	"polaris/log"
	"polaris/pkg/torznab"
	"polaris/pkg/transmission"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (s *Server) searchTvWithTorznab(name string, season, episode int) []torznab.Result {
	q := fmt.Sprintf("%s S%02dE%02d", name, season, episode)

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

type searchAndDownloadIn struct {
	ID      int `json:"id"`
	Season  int `json:"season"`
	Episode int `json:"episode"`
}

func (s *Server) searchAndDownload(seriesId, seasonNum, episodeNum int) (*string, error) {
	tr := s.db.GetTransmission()
	trc, err := transmission.NewClient(transmission.Config{
		URL: tr.URL,
		User: tr.User,
		Password: tr.Password,
	})
	if err != nil {
		return nil, errors.Wrap(err, "connect transmission")
	}
	series := s.db.GetSeriesDetails(seriesId)
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

	res := s.searchTvWithTorznab(series.OriginalName, seasonNum, episodeNum)
	if len(res) == 0 {
		return nil, fmt.Errorf("no resource found")
	}
	r1 := s.findBestMatch(res, seasonNum, episodeNum, series)
	log.Infof("found resource to download: %v", r1)
	torrent, err := trc.Download(r1.Magnet, s.db.GetDownloadDir())
	if err != nil {
		return nil, errors.Wrap(err, "downloading")
	}
	torrent.Start()

	dir := fmt.Sprintf("%s/Season %02d", series.TargetDir, ep.SeasonNumber)

	history, err :=s.db.SaveHistoryRecord(ent.History{
		SeriesID: ep.SeriesID,
		EpisodeID: ep.ID,
		SourceTitle: r1.Name,
		TargetDir: dir,
		Completed: false,
		Saved: torrent.Save(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "save record")
	}
	s.tasks[history.ID] = &Task{Torrent: torrent}

	log.Infof("success add %s to download task", r1.Name)
	return &r1.Name, nil
}

func (s *Server) findBestMatch(resources []torznab.Result,season, episode int, series *db.SeriesDetails) torznab.Result {
	var filtered []torznab.Result
	for _, r := range resources {
		if !(series.NameEn != "" && strings.Contains(r.Name,series.NameEn)) && !strings.Contains(r.Name, series.OriginalName) {
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

func (s *Server) SearchAndDownload(c *gin.Context) (interface{}, error) {
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