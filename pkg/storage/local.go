package storage

import (
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"polaris/log"

	"github.com/pkg/errors"
)

type Storage interface {
	Move(src, dest string) error
	Copy(src, dest string) error
	ReadDir(dir string) ([]fs.FileInfo, error)
	ReadFile(string) ([]byte, error)
	WriteFile(string, []byte) error
}

func NewLocalStorage(dir string) (*LocalStorage, error) {
	os.MkdirAll(dir, 0655)

	return &LocalStorage{dir: dir}, nil
}

type LocalStorage struct {
	dir string
}

func (l *LocalStorage) Copy(src, destDir string) error {
	os.MkdirAll(filepath.Join(l.dir, destDir), os.ModePerm)

	targetBase := filepath.Join(l.dir, destDir, filepath.Base(src)) //文件的场景，要加上文件名, move filename ./dir/
	info, err := os.Stat(src)
	if err != nil {
		return errors.Wrap(err, "read source dir")
	}
	if info.IsDir() { //如果是路径，则只移动路径里面的文件，不管当前路径, 行为类似 move dirname/* target_dir/
		targetBase = filepath.Join(l.dir, destDir)
	}
	log.Debugf("local storage target base dir is: %v", targetBase)

	err = filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return errors.Wrapf(err, "relation between %s and %s", src, path)
		}
		destName := filepath.Join(targetBase, rel)

		if info.IsDir() {
			os.Mkdir(destName, os.ModePerm)
		} else { //is file
			if err := os.Link(path, destName); err != nil {
				log.Warnf("hard link file error: %v, will try copy file, source: %s, dest: %s", err, path, destName)
				if writer, err := os.OpenFile(destName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm); err != nil {
					return errors.Wrapf(err, "create file %s", destName)
				} else {
					defer writer.Close()
					if f, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm); err != nil {
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
		}
		log.Infof("file copy complete: %v", destName)
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "move file error")
	}
	return nil
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
