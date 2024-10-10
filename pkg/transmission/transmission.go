package transmission

import (
	"context"
	"fmt"
	"net/url"
	"polaris/log"
	"polaris/pkg"
	"polaris/pkg/utils"

	"github.com/hekmon/transmissionrpc/v3"
	"github.com/pkg/errors"
)

func NewClient(c Config) (*Client, error) {
	u, err := url.Parse(c.URL)
	if err != nil {
		return nil, errors.Wrap(err, "parse url")
	}
	if c.User != "" {
		log.Info("transmission login with user: ", c.User)
		u.User = url.UserPassword(c.User, c.Password)
	}
	u.Path = "/transmission/rpc"

	tbt, err := transmissionrpc.New(u, nil)
	if err != nil {
		return nil, errors.Wrap(err, "connect transmission")
	}
	_, err = tbt.TorrentGetAll(context.TODO())
	if err != nil {
		return nil, errors.Wrap(err, "transmission cannot connect")
	}
	return &Client{c: tbt, cfg: c}, nil
}

type Config struct {
	URL      string `json:"url"`
	User     string `json:"user"`
	Password string `json:"password"`
}
type Client struct {
	c   *transmissionrpc.Client
	cfg Config
}

func (c *Client) GetAll() ([]pkg.Torrent, error) {
	all, err := c.c.TorrentGetAll(context.TODO())
	if err != nil {
		return nil, errors.Wrap(err, "get all")
	}
	var torrents []pkg.Torrent
	for _, t := range all {
		torrents = append(torrents, &Torrent{
			hash: *t.HashString,
			c:    c.c,
			//cfg: c.cfg,
		})
	}
	return torrents, nil
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

	t, err := c.c.TorrentAdd(context.TODO(), transmissionrpc.TorrentAddPayload{
		Filename:    &magnet,
		DownloadDir: &dir,
	})
	log.Debugf("get torrent info: %+v", t)

	return &Torrent{
		hash: hash,
		c:    c.c,
		//cfg: c.cfg,
	}, err
}

func NewTorrent(cfg Config, link string) (*Torrent, error) {
	c, err := NewClient(cfg)
	if err != nil {
		return nil, err
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
		c:    c.c,
		hash: hash,
		//cfg: cfg,
	}
	if !t.Exists() {
		return nil, errors.Errorf("torrent not exist: %v", magnet)
	}
	return t, nil
}

type Torrent struct {
	//t *transmissionrpc.Torrent
	c    *transmissionrpc.Client
	hash string
	//cfg Config
}

func (t *Torrent) getTorrent() (transmissionrpc.Torrent, error) {
	r, err := t.c.TorrentGetAllForHashes(context.TODO(), []string{t.hash})
	if err != nil {
		log.Errorf("get torrent info for error: %v", err)
	}
	if len(r) == 0 {
		return transmissionrpc.Torrent{}, fmt.Errorf("no torrent")
	}
	return r[0], nil
}

func (t *Torrent) Exists() bool {
	r, err := t.c.TorrentGetAllForHashes(context.TODO(), []string{t.hash})
	if err != nil {
		log.Errorf("get torrent info for error: %v", err)
	}
	return len(r) > 0
}

func (t *Torrent) Name() (string, error) {
	tt, err := t.getTorrent()
	if err != nil {
		return "", err
	}
	return *tt.Name, nil
}

func (t *Torrent) Progress() (int, error) {
	tt, err := t.getTorrent()
	if err != nil {
		return 0, err
	}
	if tt.IsFinished != nil && *tt.IsFinished {
		return 100, nil
	}
	if tt.PercentComplete != nil && *tt.PercentComplete >= 1 {
		return 100, nil
	}

	if tt.PercentComplete != nil {
		p := int(*tt.PercentComplete * 100)
		if p == 100 {
			p = 99
		}
		return p, nil
	}
	return 0, nil
}

func (t *Torrent) Stop() error {
	return t.c.TorrentStopHashes(context.TODO(), []string{t.hash})
}

func (t *Torrent) SeedRatio() (float64, error) {
	tt, err := t.getTorrent()
	if err != nil {
		return 0, err
	}
	if tt.UploadRatio == nil {
		return 0, nil
	}
	return *tt.UploadRatio, nil
}

func (t *Torrent) Start() error {
	return t.c.TorrentStartHashes(context.TODO(), []string{t.hash})
}

func (t *Torrent) Remove() error {
	tt, err := t.getTorrent()
	if err != nil {
		return errors.Wrap(err, "get torrent")
	}
	return t.c.TorrentRemove(context.TODO(), transmissionrpc.TorrentRemovePayload{
		IDs:             []int64{*tt.ID},
		DeleteLocalData: true,
	})
}

func (t *Torrent) Size() (int, error) {
	tt, err := t.getTorrent()
	if err != nil {
		return 0, errors.Wrap(err, "get torrent")
	}
	return int(tt.TotalSize.Byte()), nil
}

func (t *Torrent) GetHash() string {
	return t.hash
}
