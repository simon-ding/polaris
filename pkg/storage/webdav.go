package storage

import (
	"context"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"polaris/log"
	"time"

	"github.com/emersion/go-webdav"
	"github.com/pkg/errors"
)

type FileInfo struct {
	Path     string
	Size     int64
	ModTime  time.Time
	IsDir    bool
	MIMEType string
	ETag     string
}

type WebdavStorage struct {
	fs *webdav.Client
}

func NewWebdavStorage(url, user, password string) (*WebdavStorage, error) {
	c, err := webdav.NewClient(webdav.HTTPClientWithBasicAuth(http.DefaultClient, user, password), url)
	if err != nil {
		return nil, errors.Wrap(err, "new webdav")
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

			if err := w.fs.Mkdir(context.TODO(), remoteName); err != nil {
				return errors.Wrapf(err, "mkdir %v", remoteName)
			}

		} else { //is file
			if writer, err := w.fs.Create(context.TODO(), remoteName); err != nil {
				return errors.Wrapf(err, "create file %s", remoteName)
			} else {
				defer writer.Close()
				if f, err := os.OpenFile(path, os.O_RDONLY, 0666); err != nil {
					return errors.Wrapf(err, "read file %v", path)
				} else { //open success
					defer f.Close()
					_, err := io.Copy(writer, f)
					if err != nil {
						return errors.Wrap(err, "transmitting data error")
					}
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

func (w *WebdavStorage) ReadDir(dir string) ([]FileInfo, error) {
	fi, err := w.fs.ReadDir(context.TODO(), dir, false)
	if err != nil {
		return nil, err
	}
	var res []FileInfo = make([]FileInfo, 0, len(fi))
	for _, f := range fi {
		res = append(res, FileInfo{
			Path:     f.Path,
			Size:     f.Size,
			ModTime:  f.ModTime,
			IsDir:    f.IsDir,
			MIMEType: f.MIMEType,
			ETag:     f.ETag,
		})
	}
	return res, nil
}
