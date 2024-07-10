package transmission

import (
	"net/url"
	"strconv"

	"github.com/hekmon/transmissionrpc"
	"github.com/pkg/errors"
)

func NewClient(url1, user, password string) (*Client, error) {
	u, err := url.Parse(url1)
	if err != nil {
		return nil, errors.Wrap(err, "parse url")
	}
	tbt, err := transmissionrpc.New(u.Hostname(), user, password, nil)
	if err != nil {
		return nil, errors.Wrap(err, "connect transmission")
	}
	return &Client{c: tbt}, nil
}

type Client struct {
	c *transmissionrpc.Client
}

func (c *Client) Download(magnet, dir string) (*Torrent, error) {
	t, err := c.c.TorrentAdd(&transmissionrpc.TorrentAddPayload{
		Filename: &magnet,
		DownloadDir: &dir,
	})
	return &Torrent{
		t: t,
		c: c.c,
	}, err
}

type Torrent struct {
	t *transmissionrpc.Torrent
	c *transmissionrpc.Client
}

func (t *Torrent) Name() string {
	return *t.t.Name
}

func (t *Torrent) Progress() int {
	if *t.t.IsFinished {
		return 100
	}
	return int(*t.t.PercentDone*100)
}

func (t *Torrent) Stop() error {
	return t.c.TorrentStopIDs([]int64{*t.t.ID})
}

func (t *Torrent) Start() error {
	return t.c.TorrentStartIDs([]int64{*t.t.ID})
}

func (t *Torrent) Remove() error {
	return t.c.TorrentRemove(&transmissionrpc.TorrentRemovePayload{
		IDs: []int64{*t.t.ID},
		DeleteLocalData: true,
	})
}

func (t *Torrent) Save() string {
	return strconv.Itoa(int(*t.t.ID))
}