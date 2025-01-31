package storage

import (
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"polaris/pkg/gowebdav"

	"github.com/gabriel-vasile/mimetype"
	"github.com/pkg/errors"
)

type WebdavStorage struct {
	fs              *gowebdav.Client
	dir             string
	changeMediaHash bool
	progresser      func() float64
	videoFormats    []string
	subtitleFormats []string
}

func NewWebdavStorage(url, user, password, path string, changeMediaHash bool, videoFormats []string, subtitleFormats []string) (*WebdavStorage, error) {
	c := gowebdav.NewClient(url, user, password)
	if err := c.Connect(); err != nil {
		return nil, errors.Wrap(err, "connect webdav")
	}
	return &WebdavStorage{
		fs:  c,
		dir: path,
		videoFormats: videoFormats,
		subtitleFormats: subtitleFormats,
	}, nil
}

func (w *WebdavStorage) Copy(local, remoteDir string, walkFn WalkFn) error {
	b, err := NewBase(local, w.videoFormats, w.subtitleFormats)
	if err != nil {
		return err
	}

	w.progresser = b.Progress

	uploadFunc := func(destPath string, destInfo fs.FileInfo, srcReader io.Reader, mtype *mimetype.MIME) error {
		callback := func(r *http.Request) {
			r.Header.Set("Content-Type", mtype.String())
			r.ContentLength = destInfo.Size()
		}

		if err := w.fs.WriteStream(destPath, srcReader, 0666, callback); err != nil {
			return errors.Wrap(err, "transmitting data error")
		}
		return nil

	}

	return b.Upload(filepath.Join(w.dir, remoteDir), false, true, w.changeMediaHash, uploadFunc, func(s string) error {
		return nil
	}, walkFn)
}

func (w *WebdavStorage) ReadDir(dir string) ([]fs.FileInfo, error) {
	return w.fs.ReadDir(filepath.Join(w.dir, dir))
}

func (w *WebdavStorage) ReadFile(name string) ([]byte, error) {
	return w.fs.Read(filepath.Join(w.dir, name))
}

func (w *WebdavStorage) WriteFile(name string, data []byte) error {
	return w.fs.Write(filepath.Join(w.dir, name), data, os.ModePerm)
}

func (w *WebdavStorage) UploadProgress() float64 {
	if w.progresser == nil {
		return 0
	}
	return w.progresser()
}

func (w *WebdavStorage) RemoveAll(path string) error {
	return w.fs.RemoveAll(filepath.Join(w.dir, path))
}
