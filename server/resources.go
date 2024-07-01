package server

import (
	"fmt"
	"polaris/db"
	"polaris/log"
	"polaris/pkg/downloader"
	"polaris/pkg/torznab"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (s *Server) searchTvWithTorznab(name string, season, episode int) []torznab.Result {
	q := fmt.Sprintf("%s S%02dE%02d", name, season, episode)

	var res []torznab.Result
	allTorznab := s.db.GetAllTorznabInfo()
	for name, setting := range allTorznab {
		resp, err := torznab.Search(setting.URL, setting.ApiKey, q)
		if err != nil {
			log.Errorf("search %s error: %v", name, err)
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

type searchAndDownloadIn struct {
	Title string `json:"title"`
	Season int `json:"season"`
	Episode int `json:"episode"`
}

func (s *Server) SearchAndDownload(c *gin.Context) (interface{}, error) {
	var in searchAndDownloadIn
	if err := c.ShouldBindJSON(&in); err != nil {
		return nil, errors.Wrap(err, "bind json")
	}
	res := s.searchTvWithTorznab(in.Title, in.Season, in.Episode)
	r1 := res[0]
	t, err := downloader.DownloadByMagnet(r1.Magnet, "~")
	if err != nil {
		return nil, errors.Wrap(err, "download torrent")
	}
	s.tasks[r1.Name] = t
	log.Errorf("success add %s to download task", r1.Name)
	return nil, nil
}