package transmission

import (
	"context"
	"encoding/json"
	"net/url"
	"polaris/log"

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
	return &Client{c: tbt, cfg: c}, nil
}

type Config struct {
	URL      string `json:"url"`
	User     string `json:"user"`
	Password string `json:"password"`
}
type Client struct {
	c *transmissionrpc.Client
	cfg Config
}

func (c *Client) Download(magnet, dir string) (*Torrent, error) {
	t, err := c.c.TorrentAdd(context.TODO(), transmissionrpc.TorrentAddPayload{
		Filename:    &magnet,
		DownloadDir: &dir,
	})
	log.Infof("get torrent info: %+v", t)

	return &Torrent{
		ID: *t.ID,
		c:  c.c,
		Config: c.cfg,
	}, err
}

type Torrent struct {
	//t *transmissionrpc.Torrent
	c  *transmissionrpc.Client
	ID int64 `json: "id"`
	Config
}

func (t *Torrent) reloadClient() error {
	c, err := NewClient(t.Config)
	if err != nil {
		return err
	}
	t.c = c.c
	return nil
}

func (t *Torrent) getTorrent() transmissionrpc.Torrent {
	r, err := t.c.TorrentGetAllFor(context.TODO(), []int64{t.ID})
	if err != nil {
		log.Errorf("get torrent info for error: %v", err)
	}
	return r[0]
}

func (t *Torrent) Exists() bool {
	r, err := t.c.TorrentGetAllFor(context.TODO(), []int64{t.ID})
	if err != nil {
		log.Errorf("get torrent info for error: %v", err)
	}
	return len(r) > 0
}

func (t *Torrent) Name() string {
	return *t.getTorrent().Name
}

func (t *Torrent) Progress() int {
	if t.getTorrent().IsFinished != nil && *t.getTorrent().IsFinished {
		return 100
	}
	if t.getTorrent().PercentDone != nil {
		return int(*t.getTorrent().PercentDone * 100)
	}
	return 0
}

func (t *Torrent) Stop() error {
	return t.c.TorrentStopIDs(context.TODO(), []int64{t.ID})
}

func (t *Torrent) Start() error {
	return t.c.TorrentStartIDs(context.TODO(), []int64{t.ID})
}

func (t *Torrent) Remove() error {
	return t.c.TorrentRemove(context.TODO(), transmissionrpc.TorrentRemovePayload{
		IDs:             []int64{t.ID},
		DeleteLocalData: true,
	})
}

func (t *Torrent) Save() string {

	d, _ := json.Marshal(*t)
	return string(d)
}

func ReloadTorrent(s string) (*Torrent, error) {
	var torrent = Torrent{}
	err := json.Unmarshal([]byte(s), &torrent)
	if err != nil {
		return nil, err
	}

	err = torrent.reloadClient()
	if err != nil {
		return nil, errors.Wrap(err, "reload client")
	}
	return &torrent, nil
}