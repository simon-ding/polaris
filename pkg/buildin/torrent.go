package buildin

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"polaris/log"
	"polaris/pkg"
	"strings"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/pkg/errors"
)

func NewDownloader(downloadDir string) (*Downloader, error) {
	cfg := torrent.NewDefaultClientConfig()
	cfg.DataDir = downloadDir
	cfg.ListenPort = 51243
	t, err := torrent.NewClient(cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "create torrent client")
	}
	return &Downloader{cl: t, dir: downloadDir}, nil
}

type Downloader struct {
	cl  *torrent.Client
	dir string
}

func (d *Downloader) GetAll() ([]pkg.Torrent, error) {
	ts := d.cl.Torrents()
	var res []pkg.Torrent
	for _, t := range ts {
		res = append(res, &Torrent{
			t:    t,
			cl:   d.cl,
			hash: t.InfoHash().HexString(),
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
		<-t.GotInfo()
		return &Torrent{
			t:    t,
			cl:   d.cl,
			hash: hash,
			dir:  d.dir,
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
	<-t.GotInfo()
	return &Torrent{
		t:    t,
		cl:   d.cl,
		hash: hash,
		dir:  d.dir,
	}, nil
}

func NewTorrentFromHash(hash string, downloadDir string) (*Torrent, error) {
	cl, err := NewDownloader(downloadDir)
	if err != nil {
		return nil, errors.Wrap(err, "create downloader")
	}
	ttt := cl.cl.Torrents()
	log.Infof("all torrents: %+v", ttt)
	t, _ := cl.cl.AddTorrentInfoHash(metainfo.NewHashFromHex(hash))
	// if new {
	// 	return nil, fmt.Errorf("torrent not found")
	// }
	<-t.GotInfo()
	return &Torrent{
		t:    t,
		cl:   cl.cl,
		hash: hash,
		dir:  downloadDir,
	}, nil
}

type Torrent struct {
	t    *torrent.Torrent
	cl   *torrent.Client
	hash string
	dir  string
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
	files := t.t.Files()
	for _, file := range files {
		name := file.Path()
		if err := os.RemoveAll(filepath.Join(t.dir, name)); err != nil {
			return errors.Errorf("remove file (%s) error: %v", file.Path(), err)
		}
	}
	t.t.Drop()
	return nil
}

func (t *Torrent) Exists() bool {
	tors := t.cl.Torrents()
	for _, to := range tors {
		if to.InfoHash().HexString() == t.hash {
			return true
		}
	}
	return false
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
			name := filepath.Join(t.dir, file.Path())
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
