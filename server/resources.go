package server

import (
	"fmt"
	"polaris/db"
	"polaris/ent"
	"polaris/log"
	"polaris/pkg/torznab"
	"polaris/pkg/transmission"
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
		return nil, fmt.Errorf("no indexer found")
	}
	return indexers, nil
}

type searchAndDownloadIn struct {
	ID      int `json:"id"`
	Season  int `json:"season"`
	Episode int `json:"episode"`
}

func (s *Server) SearchAndDownload(c *gin.Context) (interface{}, error) {
	var in searchAndDownloadIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind json")
	}
	tr := s.db.GetTransmission()
	trc, err := transmission.NewClient(tr.URL, tr.User, tr.Password)
	if err != nil {
		return nil, errors.Wrap(err, "connect transmission")
	}
	log.Infof("search episode resources link: %v", in)
	series := s.db.GetSeriesDetails(in.ID)
	if series == nil {
		return nil, fmt.Errorf("no tv series of id %v", in.ID)
	}
	var ep *ent.Episode
	for _, e := range series.Episodes {
		if e.SeasonNumber == in.Season && e.EpisodeNumber == in.Episode {
			ep = e
		}
	}
	if ep == nil {
		return nil, errors.Errorf("no episode of season %d episode %d", in.Season, in.Episode)
	}

	res := s.searchTvWithTorznab(series.OriginalName, in.Season, in.Episode)
	if len(res) == 0 {
		return "", fmt.Errorf("no resource found")
	}
	r1 := res[0]
	log.Infof("found resource to download: %v", r1)
	torrent, err := trc.Download(r1.Magnet, s.db.GetDownloadDir())
	if err != nil {
		return nil, errors.Wrap(err, "downloading")
	}
	torrent.Start()

	var name = series.NameEn
	if name == "" {
		name = series.OriginalName
	}
	var year = strings.Split(series.AirDate, "-")[0]

	dir := fmt.Sprintf("%s (%s)/Season %02d", name, year, ep.SeasonNumber)

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
	s.tasks[history.ID] = torrent

	// t, err := downloader.DownloadByMagnet(r1.Magnet, "~")
	// if err != nil {
	// 	return nil, errors.Wrap(err, "download torrent")
	// }
	log.Errorf("success add %s to download task", r1.Name)
	return gin.H{
		"name": r1.Name,
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
		return nil, fmt.Errorf("no download client")
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