package qbittorrent

import (
	"encoding/json"
	"fmt"
	"polaris/pkg"
	"polaris/pkg/go-qbittorrent/qbt"
	"polaris/pkg/utils"

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
	tt, err := c.c.Torrents(qbt.TorrentsOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "get torrents")
	}
	var res []pkg.Torrent
	for _, t := range tt {
		t1 := &Torrent{
			c:    c.c,
			hash: t.Hash,
			//Info: c.Info,
		}
		res = append(res, t1)
	}
	return res, nil
}

func (c *Client) Download(link, dir string) (pkg.Torrent, error) {
	magnet, err := utils.Link2Magnet(link)
	if err != nil {
		return nil, errors.Errorf("converting link to magnet error, link: %v, error: %v", link, err)
	}

	hash, err := utils.MagnetHash(magnet)
	if err != nil {
		return nil, errors.Wrap(err, "get hash")
	}
	err = c.c.DownloadLinks([]string{magnet}, qbt.DownloadOptions{Savepath: &dir})
	if err != nil {
		return nil, errors.Wrap(err, "qbt download")
	}
	return &Torrent{hash: hash, c: c.c, }, nil

}


func NewTorrent(info Info, link string) (*Torrent, error) {
	c, err := NewClient(info.URL, info.User, info.Password)
	if err != nil {
		return nil,  err
	}
	magnet, err := utils.Link2Magnet(link)
	if err != nil {
		return nil, errors.Errorf("converting link to magnet error, link: %v, error: %v", link, err)
	}

	hash, err := utils.MagnetHash(magnet)
	if err != nil {
		return nil, err
	}
	t := &Torrent{
		c: c.c,
		hash: hash,
	}
	if !t.Exists() {
		return nil, errors.Errorf("torrent not exist: %v", magnet)
	}
	return t, nil
}
type Torrent struct {
	c    *qbt.Client
	hash string
	//info Info
}

func (t *Torrent) GetHash() string {
	return t.hash
}

func (t *Torrent) getTorrent() (*qbt.TorrentInfo, error) {
	all, err := t.c.Torrents(qbt.TorrentsOptions{Hashes: []string{t.hash}})
	if err != nil {
		return nil, err
	}
	if len(all) == 0 {
		return nil, fmt.Errorf("no such torrent: %v", t.hash)
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
	p := qb.Progress * 100
	if p >= 100 {
		return 100, nil
	}
	if int(p) == 100 {
		return 99, nil
	}

	return int(p), nil
}

func (t *Torrent) Stop() error {
	return t.c.Pause([]string{t.hash})
}

func (t *Torrent) Start() error {
	ok, err := t.c.Resume([]string{t.hash})
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("status not 200")
	}
	return nil
}

func (t *Torrent) Remove() error {
	ok, err := t.c.Delete([]string{t.hash}, true)
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
