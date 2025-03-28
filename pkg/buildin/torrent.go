package buildin

import (
	"fmt"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/pkg/errors"
	"io/fs"
	"net/http"
	"os"
	"polaris/pkg"
	"strings"
)

type Downloader struct {
	cl *torrent.Client
}

func (d *Downloader) GetAll() ([]pkg.Torrent, error) {
	ts := d.cl.Torrents()
	var res []pkg.Torrent
	for _, t := range ts {
		res = append(res, &Torrent{
			t:  t,
			cl: d.cl,
		})
	}
	return res, nil
}

func (d *Downloader) Download(link, hash, dir string) (pkg.Torrent, error) {

	if strings.HasPrefix(strings.ToLower(link), "magnet:") {
		t, err := d.cl.AddMagnet(link)
		if err != nil {
			return nil, fmt.Errorf("failed to add magnet: %v", err)
		}
		return &Torrent{
			t:    t,
			cl:   d.cl,
			hash: hash,
		}, nil
	}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse //do not follow redirects
		},
	}

	resp, err := client.Get(link)
	if err != nil {
		return nil, errors.Wrap(err, "get link")
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		//redirects
		tourl := resp.Header.Get("Location")
		return d.Download(tourl, hash, dir)
	}
	info, err := metainfo.Load(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to load metainfo: %v", err)
	}
	t, err := d.cl.AddTorrent(info)
	if err != nil {
		return nil, fmt.Errorf("failed to add torrent: %v", err)
	}
	return &Torrent{
		t:    t,
		cl:   d.cl,
		hash: hash,
	}, nil
}

type Torrent struct {
	t    *torrent.Torrent
	cl   *torrent.Client
	hash string
}

func (t *Torrent) Name() (string, error) {
	return t.t.Name(), nil
}

func (t *Torrent) TotalSize() int64 {
	var c int64
	for _, f := range t.t.Files() {
		c += f.FileInfo().Length
	}
	return c
}

func (t *Torrent) Progress() (int, error) {
	if t.t.Complete().Bool() {
		return 100, nil
	}
	return int(t.t.BytesCompleted() * 100 / t.TotalSize()), nil
}

func (t *Torrent) Stop() error {
	return nil
}

func (t *Torrent) Start() error {
	<-t.t.GotInfo()
	t.t.DownloadAll()
	return nil
}

func (t *Torrent) Remove() error {
	t.t.Drop()
	return nil
}

func (t *Torrent) Exists() bool {
	return true
}

func (t *Torrent) SeedRatio() (float64, error) {
	return 0, nil
}

func (t *Torrent) GetHash() string {
	return t.hash
}

func (t *Torrent) WalkFunc() func(fn func(path string, info fs.FileInfo) error) error {
	files := t.t.Files()

	return func(fn func(path string, info fs.FileInfo) error) error {
		for _, file := range files {
			name := file.Path()
			info, err := os.Stat(name)
			if err != nil {
				return err
			}

			if err := fn(name, info); err != nil {
				return errors.Errorf("proccess file (%s) error: %v", file.Path(), err)
			}
		}
		return nil

	}
}
