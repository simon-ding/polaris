package storage

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"polaris/log"
	"polaris/pkg/gowebdav"
	"polaris/pkg/utils"

	"github.com/gabriel-vasile/mimetype"
	"github.com/pkg/errors"
)

type WebdavStorage struct {
	fs *gowebdav.Client
	dir string
	changeMediaHash bool
}

func NewWebdavStorage(url, user, password, path string, changeMediaHash bool) (*WebdavStorage, error) {
	c := gowebdav.NewClient(url, user, password)
	if err := c.Connect(); err != nil {
		return nil, errors.Wrap(err, "connect webdav")
	}
	return &WebdavStorage{
		fs: c,
		dir: path,
	}, nil
}

func (w *WebdavStorage) Move(local, remoteDir string) error {

	remoteBase := filepath.Join(w.dir,remoteDir, filepath.Base(local))
	info, err := os.Stat(local)
	if err != nil {
		return errors.Wrap(err, "read source dir")
	}
	if info.IsDir() { //如果是路径，则只移动路径里面的文件，不管当前路径, 行为类似 move dirname/* target_dir/
		remoteBase = filepath.Join(w.dir, remoteDir)
	}

	//log.Infof("remove all content in %s", remoteBase)
	//w.fs.RemoveAll(remoteBase)
	err = filepath.Walk(local, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "read file %v", path)
		}

		rel, err := filepath.Rel(local, path)
		if err != nil {
			return errors.Wrap(err, "path relation")
		}
		remoteName := filepath.Join(remoteBase, rel)

		if info.IsDir() {
			log.Infof("skip dir %v, webdav will mkdir automatically", info.Name())

			// if err := w.fs.Mkdir(remoteName, 0666); err != nil {
			// 	return errors.Wrapf(err, "mkdir %v", remoteName)
			// }

		} else { //is file
			if w.changeMediaHash {
				if err := utils.ChangeFileHash(path); err != nil {
					log.Errorf("change file %v hash error: %v", path, err)
				}
			}
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


func (w *WebdavStorage) ReadFile(name string) ([]byte, error) {
	return w.fs.Read(filepath.Join(w.dir, name))
}


func (w *WebdavStorage) WriteFile(name string, data []byte) error  {
	return w.fs.Write(filepath.Join(w.dir, name), data, os.ModePerm)
}