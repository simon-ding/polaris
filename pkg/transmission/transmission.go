package transmission

import (
	"context"
	"net/url"
	"polaris/log"
	"strconv"

	"github.com/hekmon/transmissionrpc/v3"
	"github.com/pkg/errors"
)

func NewClient(url1, user, password string) (*Client, error) {
	u, err := url.Parse(url1)
	if err != nil {
		return nil, errors.Wrap(err, "parse url")
	}
	if user != "" {
		log.Info("transmission login with user: ", user)
		u.User = url.UserPassword(user, password)
	}
	u.Path = "/transmission/rpc"
	
	tbt, err := transmissionrpc.New(u, nil)
	if err != nil {
		return nil, errors.Wrap(err, "connect transmission")
	}
	return &Client{c: tbt}, nil
}

type Client struct {
	c *transmissionrpc.Client
}

func (c *Client) Download(magnet, dir string) (*Torrent, error) {
	t, err := c.c.TorrentAdd(context.TODO(), transmissionrpc.TorrentAddPayload{
		Filename:    &magnet,
		DownloadDir: &dir,
	})
	log.Infof("get torrent info: %+v", t)

	return &Torrent{
		id: *t.ID,
		c:  c.c,
	}, err
}

type Torrent struct {
	//t *transmissionrpc.Torrent
	c  *transmissionrpc.Client
	id int64
}

func (t *Torrent) getTorrent() transmissionrpc.Torrent {
	r, err := t.c.TorrentGetAllFor(context.TODO(), []int64{t.id})
	if err != nil {
		log.Errorf("get torrent info for error: %v", err)
	}
	return r[0]
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
	return t.c.TorrentStopIDs(context.TODO(), []int64{t.id})
}

func (t *Torrent) Start() error {
	return t.c.TorrentStartIDs(context.TODO(), []int64{t.id})
}

func (t *Torrent) Remove() error {
	return t.c.TorrentRemove(context.TODO(), transmissionrpc.TorrentRemovePayload{
		IDs:             []int64{t.id},
		DeleteLocalData: true,
	})
}

func (t *Torrent) Save() string {
	return strconv.Itoa(int(t.id))
}
