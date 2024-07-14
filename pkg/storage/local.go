package storage

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type Storage interface {
	Move(src, dest string) error
	ReadDir(dir string) ([]fs.FileInfo, error)
}

func NewLocalStorage(dir string) *LocalStorage {
	return &LocalStorage{dir: dir}
}

type LocalStorage struct {
	dir string
}

func (l *LocalStorage) Move(src, dest string) error {
	targetDir := filepath.Join(l.dir, dest)
	os.MkdirAll(targetDir, 0655)
	err := filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		destName := filepath.Join(targetDir, info.Name())
		srcName := filepath.Join(src, info.Name())
		if info.IsDir() {

			if err := os.Mkdir(destName, 0666); err != nil {
				return errors.Wrapf(err, "mkdir %v", destName)
			}

		} else { //is file
			if writer, err := os.Create(destName); err != nil {
				return errors.Wrapf(err, "create file %s", destName)
			} else {
				defer writer.Close()
				if f, err := os.OpenFile(srcName, os.O_RDONLY, 0666); err != nil {
					return errors.Wrapf(err, "read file %v", srcName)
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
