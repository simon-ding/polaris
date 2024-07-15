package storage

import (
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type Storage interface {
	Move(src, dest string) error
	ReadDir(dir string) ([]fs.FileInfo, error)
}

func NewLocalStorage(dir string) (*LocalStorage, error) {
	if _, err := os.Stat(dir); err != nil {
		return nil, errors.Wrap(err, "stat")
	}
	
	return &LocalStorage{dir: dir}, nil
}

type LocalStorage struct {
	dir string
}

func (l *LocalStorage) Move(src, dest string) error {
	targetDir := filepath.Join(l.dir, dest)
	os.MkdirAll(targetDir, 0655)
	err := filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err  := filepath.Rel(src, path)
		if err != nil {
			return errors.Wrapf(err, "relation between %s and %s", src, path)
		}
		destName := filepath.Join(targetDir, rel, info.Name())

		if info.IsDir() {

			if err := os.Mkdir(destName, 0666); err != nil {
				return errors.Wrapf(err, "mkdir %v", destName)
			}

		} else { //is file
			if writer, err := os.Create(destName); err != nil {
				return errors.Wrapf(err, "create file %s", destName)
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
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "move file error")
	}
	return os.RemoveAll(src)

}


func (l *LocalStorage) ReadDir(dir string) ([]fs.FileInfo, error) {
	 return ioutil.ReadDir(filepath.Join(l.dir, dir))
}