package storage

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"polaris/log"
	"polaris/pkg/gowebdav"

	"github.com/gabriel-vasile/mimetype"
	"github.com/pkg/errors"
)

type WebdavStorage struct {
	fs *gowebdav.Client
	dir string
}

func NewWebdavStorage(url, user, password, path string) (*WebdavStorage, error) {
	c := gowebdav.NewClient(url, user, password)
	if err := c.Connect(); err != nil {
		return nil, errors.Wrap(err, "connect webdav")
	}
	return &WebdavStorage{
		fs: c,
		dir: path,
	}, nil
}

func (w *WebdavStorage) Move(local, remote string) error {
	baseLocal := filepath.Base(local)
	remoteBase := filepath.Join(w.dir,remote, baseLocal)

	log.Infof("remove all content in %s", remoteBase)
	w.fs.RemoveAll(remoteBase)
	err := filepath.Walk(local, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "read file %v", path)
		}

		rel, err := filepath.Rel(local, path)
		if err != nil {
			return errors.Wrap(err, "path relation")
		}
		remoteName := filepath.Join(remoteBase, rel)

		if info.IsDir() {

			if err := w.fs.Mkdir(remoteName, 0666); err != nil {
				return errors.Wrapf(err, "mkdir %v", remoteName)
			}

		} else { //is file
			if f, err := os.OpenFile(path, os.O_RDONLY, 0666); err != nil {
				return errors.Wrapf(err, "read file %v", path)
			} else { //open success
				defer f.Close()
				mtype, err := mimetype.DetectFile(path)
				if err != nil {
					return errors.Wrap(err, "mime type error")
				}

				callback := func(r *http.Request) {
					r.Header.Set("Content-Type", mtype.String())
					r.ContentLength = info.Size()
				}
			
				if err := w.fs.WriteStream(remoteName, f, 0666, callback); err != nil {
					return errors.Wrap(err, "transmitting data error")
				}
			}
		}
		log.Infof("file copy complete: %v", remoteName)
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "move file error")
	}
	return os.RemoveAll(local)
}

func (w *WebdavStorage) ReadDir(dir string) ([]fs.FileInfo, error) {
	return w.fs.ReadDir(filepath.Join(w.dir, dir))
}
