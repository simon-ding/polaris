package storage

import (
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
	"github.com/pkg/errors"
)


func NewLocalStorage(dir string) (*LocalStorage, error) {
	os.MkdirAll(dir, 0655)
	return &LocalStorage{dir: dir}, nil
}

type LocalStorage struct {
	dir string
}

func (l *LocalStorage) Copy(src, destDir string) error {
	b, err := NewBase(src)
	if err != nil {
		return err
	}

	baseDest := filepath.Join(l.dir, destDir)
	uploadFunc := func(destPath string, destInfo fs.FileInfo, srcReader io.Reader, mimeType *mimetype.MIME) error {
		if writer, err := os.OpenFile(destPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm); err != nil {
			return errors.Wrapf(err, "create file %s", destPath)
		} else {
			defer writer.Close()
			_, err := io.Copy(writer, srcReader)
			if err != nil {
				return errors.Wrap(err, "transmitting data error")
			}

		}
		return nil
	}
	return b.Upload(baseDest, true, false, false, uploadFunc, func(s string) error {
		return os.Mkdir(s, os.ModePerm)
	})
}

func (l *LocalStorage) Move(src, destDir string) error {
	if err := l.Copy(src, destDir); err != nil {
		return err
	}
	return os.RemoveAll(src)
}

func (l *LocalStorage) ReadDir(dir string) ([]fs.FileInfo, error) {
	return ioutil.ReadDir(filepath.Join(l.dir, dir))
}

func (l *LocalStorage) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(filepath.Join(l.dir, name))
}

func (l *LocalStorage) WriteFile(name string, data []byte) error {
	path := filepath.Join(l.dir, name)
	os.MkdirAll(filepath.Dir(path), os.ModePerm)
	return os.WriteFile(path, data, os.ModePerm)
}

func (l *LocalStorage) UploadProgress() float64 {
	return 0
}