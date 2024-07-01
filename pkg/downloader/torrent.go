package downloader

import (
	"github.com/anacrolix/torrent"
	"github.com/pkg/errors"
)



func DownloadByMagnet(magnet string, dir string) (*torrent.Torrent, error) {
	c, err := torrent.NewClient(nil)
	if err != nil {
		return nil, errors.Wrap(err, "new torrent")
	}
	defer c.Close()
	t, err := c.AddMagnet(magnet)
	if err != nil {
		return nil, errors.Wrap(err, "add torrent")
	}

	<-t.GotInfo()
	t.DownloadAll()
	c.WaitAll()
	return t, nil
}