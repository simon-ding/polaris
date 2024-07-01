package transmission

import (
	"net/url"
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

func (c *Client) Download(magnet string) (*transmissionrpc.Torrent, error) {
	t, err := c.c.TorrentAdd(&transmissionrpc.TorrentAddPayload{Filename: &magnet})
	return t, err
}
