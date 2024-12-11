package storage

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"polaris/pkg/alist"

	"github.com/gabriel-vasile/mimetype"
)

func NewAlist(cfg *alist.Config, dir string, videoFormats []string, subtitleFormats []string) (*Alist, error) {
	cl := alist.New(cfg)
	_, err := cl.Login()
	if err != nil {
		return nil, err
	}
	return &Alist{baseDir: dir, cfg: cfg, client: cl, videoFormats: videoFormats, subtitleFormats: subtitleFormats}, nil
}

type Alist struct {
	baseDir         string
	cfg             *alist.Config
	client          *alist.Client
	progresser      func() float64
	videoFormats    []string
	subtitleFormats []string
}

func (a *Alist) Move(src, dest string) error {
	if err := a.Copy(src, dest); err != nil {
		return err
	}
	return os.RemoveAll(src)
}

func (a *Alist) Copy(src, dest string) error {
	b, err := NewBase(src, a.videoFormats, a.subtitleFormats)
	if err != nil {
		return err
	}
	a.progresser = b.Progress

	uploadFunc := func(destPath string, destInfo fs.FileInfo, srcReader io.Reader, mimeType *mimetype.MIME) error {
		_, err := a.client.UploadStream(srcReader, destInfo.Size(), destPath)
		return err
	}
	mkdirFunc := func(dir string) error {
		return a.client.Mkdir(dir)
	}

	baseDest := filepath.Join(a.baseDir, dest)
	return b.Upload(baseDest, false, false, false, uploadFunc, mkdirFunc)
}

func (a *Alist) ReadDir(dir string) ([]fs.FileInfo, error) {
	return nil, nil
}

func (a *Alist) ReadFile(s string) ([]byte, error) {
	return nil, nil
}

func (a *Alist) WriteFile(s string, bytes []byte) error {
	return nil
}

func (a *Alist) UploadProgress() float64 {
	if a.progresser == nil {
		return 0
	}
	return a.progresser()
}

func (a *Alist) RemoveAll(path string) error {
	return nil
}