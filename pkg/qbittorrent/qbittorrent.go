package qbittorrent

import (
	"encoding/json"
	"fmt"
	"polaris/pkg"
	"polaris/pkg/go-qbittorrent/qbt"
	"time"

	"github.com/pkg/errors"
)

type Info struct {
	URL      string
	User     string
	Password string
}

type Client struct {
	c *qbt.Client
	Info
}

func NewClient(url, user, pass string) (*Client, error) {
	// connect to qbittorrent client
	qb := qbt.NewClient(url)

	// login to the client
	loginOpts := qbt.LoginOptions{
		Username: user,
		Password: pass,
	}
	err := qb.Login(loginOpts)
	if err != nil {
		return nil, err
	}

	return &Client{c: qb, Info: Info{URL: url, User: user, Password: pass}}, nil
}

func (c *Client) GetAll() ([]pkg.Torrent, error) {
	tt, err :=c.c.Torrents(qbt.TorrentsOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "get torrents")
	}
	var res []pkg.Torrent
	for _, t := range tt {
		t1 := &Torrent{
			c: c.c,
			Hash: t.Hash,
			Info: c.Info,
		}
		res = append(res, t1)
	}
	return res, nil
}

func (c *Client) Download(link, dir string) (pkg.Torrent, error) {
	all, err := c.c.Torrents(qbt.TorrentsOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "get old torrents")
	}
	allHash := make(map[string]bool, len(all))
	for _, t := range all {
		allHash[t.Hash] = true
	}
	err = c.c.DownloadLinks([]string{link}, qbt.DownloadOptions{Savepath: &dir})
	if err != nil {
		return nil, errors.Wrap(err, "qbt download")
	}
	var newHash string

loop:
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		all, err = c.c.Torrents(qbt.TorrentsOptions{})
		if err != nil {
			return nil, errors.Wrap(err, "get new torrents")
		}

		for _, t := range all {
			if !allHash[t.Hash] {
				newHash = t.Hash
				break loop
			}
		}
	}

	if newHash == "" {
		return nil, fmt.Errorf("download torrent fail: timeout")
	}
	return &Torrent{Hash: newHash, c: c.c, Info: c.Info}, nil

}

type Torrent struct {
	c    *qbt.Client
	Hash string
	Info
}

func (t *Torrent) GetHash() string {
	return t.Hash
}

func (t *Torrent) getTorrent() (*qbt.TorrentInfo, error) {
	all, err := t.c.Torrents(qbt.TorrentsOptions{Hashes: []string{t.Hash}})
	if err != nil {
		return nil, err
	}
	if len(all) == 0 {
		return nil, fmt.Errorf("no such torrent: %v", t.Hash)
	}
	return &all[0], nil
}

func (t *Torrent) Name() (string, error) {
	qb, err := t.getTorrent()
	if err != nil {
		return "", err
	}

	return qb.Name, nil
}

func (t *Torrent) Progress() (int, error) {
	qb, err := t.getTorrent()
	if err != nil {
		return 0, err
	}
	return int(qb.Progress), nil
}

func (t *Torrent) Stop() error {
	return t.c.Pause([]string{t.Hash})
}

func (t *Torrent) Start() error {
	ok, err := t.c.Resume([]string{t.Hash})
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("status not 200")
	}
	return nil
}

func (t *Torrent) Remove() error {
	ok, err := t.c.Delete([]string{t.Hash}, true)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("status not 200")
	}
	return nil
}

func (t *Torrent) Save() string {
	data, _ := json.Marshal(t)
	return string(data)
}

func (t *Torrent) Exists() bool {
	_, err := t.getTorrent()
	return err == nil
}

func (t *Torrent) SeedRatio() (float64, error) {
	qb, err := t.getTorrent()
	if err != nil {
		return 0, err
	}
	return qb.Ratio, nil
}
