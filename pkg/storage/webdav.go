package storage

import (
	"io/fs"
	"os"
	"path/filepath"
	"polaris/log"

	"github.com/pkg/errors"
	"github.com/studio-b12/gowebdav"
)

type WebdavStorage struct {
	fs *gowebdav.Client
}

func NewWebdavStorage(url, user, password string) (*WebdavStorage, error) {
	c := gowebdav.NewClient(url, user, password)
	if err := c.Connect(); err != nil {
		return nil, errors.Wrap(err, "connect webdav")
	}
	return &WebdavStorage{
		fs: c,
	}, nil
}

func (w *WebdavStorage) Move(local, remote string) error {
	baseLocal := filepath.Base(local)

	err := filepath.Walk(local, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "read file %v", path)
		}

		rel, err := filepath.Rel(local, path)
		if err != nil {
			return errors.Wrap(err, "path relation")
		}
		remoteName := filepath.Join(remote, baseLocal, rel)

		if info.IsDir() {

			if err := w.fs.Mkdir(remoteName, 0666); err != nil {
				return errors.Wrapf(err, "mkdir %v", remoteName)
			}

		} else { //is file
			if f, err := os.OpenFile(path, os.O_RDONLY, 0666); err != nil {
				return errors.Wrapf(err, "read file %v", path)
			} else { //open success
				defer f.Close()

				if err := w.fs.WriteStream(remoteName, f, 0666); err != nil {
					return errors.Wrap(err, "transmitting data error")
				}
			}
		}
		log.Infof("file copy complete: %d", remoteName)
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "move file error")
	}
	return os.RemoveAll(local)
}

func (w *WebdavStorage) ReadDir(dir string) ([]fs.FileInfo, error) {
	return w.fs.ReadDir(dir)
}
