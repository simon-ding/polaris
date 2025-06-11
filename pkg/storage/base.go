package storage

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"polaris/log"
	"polaris/pkg/utils"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/pkg/errors"
)

type WalkFn func(fn func(path string, info fs.FileInfo) error) error
type Storage interface {
	//Move(src, dest string) error
	Copy(src, dest string, walkFn WalkFn) error
	ReadDir(dir string) ([]fs.FileInfo, error)
	ReadFile(string) ([]byte, error)
	WriteFile(string, []byte) error
	UploadProgress() float64
	RemoveAll(path string) error
}

type uploadFunc func(destPath string, destInfo fs.FileInfo, srcReader io.Reader, mimeType *mimetype.MIME) error

type Base struct {
	src             string
	videoFormats    []string
	subtitleFormats []string
	totalSize       int64
	uploadedSize    int64
}

func NewBase(src string, videoFormats []string, subtitleFormats []string) (*Base, error) {
	b := &Base{src: src, videoFormats: videoFormats, subtitleFormats: subtitleFormats}
	err := b.calculateSize()
	return b, err
}

func (b *Base) checkVideoFilesExist() bool {
	if len(b.videoFormats) == 0 { // do not check
		return true
	}
	hasVideo := false
	filepath.Walk(b.src, func(path string, info fs.FileInfo, err error) error {
		ext := filepath.Ext(strings.ToLower(info.Name()))

		for _, f := range b.videoFormats {
			if f == ext {
				hasVideo = true
			}
		}
		return nil
	})
	return hasVideo
}

func (b *Base) isFileNeeded(name string) bool {
	ext := filepath.Ext(strings.ToLower(name))
	if len(b.videoFormats) == 0 {
		return true
	} else {
		for _, f := range b.videoFormats {
			if f == ext {
				return true
			}
		}
	}
	if len(b.subtitleFormats) > 0 {
		for _, f := range b.subtitleFormats {
			if f == ext {
				return true
			}
		}
	}
	return false

}

func (b *Base) Upload(destDir string, tryLink, detectMime, changeMediaHash bool, upload uploadFunc, mkdir func(string) error, walkFn WalkFn) error {
	if !b.checkVideoFilesExist() {
		return errors.Errorf("torrent has no video file(s)")
	}
	os.MkdirAll(destDir, os.ModePerm)

	targetBase := filepath.Join(destDir, filepath.Base(b.src)) //文件的场景，要加上文件名, move filename ./dir/
	info, err := os.Stat(b.src)
	if err != nil {
		return errors.Wrap(err, "read source dir")
	}
	if info.IsDir() { //如果是路径，则只移动路径里面的文件，不管当前路径, 行为类似 move dirname/* target_dir/
		targetBase = destDir
	}
	log.Debugf("local storage target base dir is: %v", targetBase)

	err = walkFn(func(path string, info fs.FileInfo) (err error) {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(b.src, path)
		if err != nil {
			return errors.Wrapf(err, "relation between %s and %s", b.src, path)
		}
		destName := filepath.Clean(filepath.Join(targetBase, rel))
		if !strings.HasPrefix(destName, targetBase) {
			//如果目标路径不是在目标目录下，则报错
			return errors.Errorf("destination: %s is not in target dir: %s", destName, targetBase)
		}

		if info.IsDir() {
			mkdir(destName)
		} else { //is file
			if !b.isFileNeeded(info.Name()) {
				log.Debugf("file is not needed, skip: %s", info.Name())
				return nil
			}
			
			defer func ()  {
				if err == nil {
					log.Infof("copy file success, filename %s, destination %s", rel, destName)
				}
			}()

			if tryLink {
				if err := os.Link(path, destName); err == nil {
					return nil //link success
				}
				log.Warnf("hard link file error: %v, will try copy file, source: %s, dest: %s", err, path, destName)
			}
			if changeMediaHash {
				if err := utils.ChangeFileHash(path); err != nil {
					log.Errorf("change file %v hash error: %v", path, err)
				}
			}

			if f, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm); err != nil {
				return errors.Wrapf(err, "read file %v", path)
			} else { //open success
				defer f.Close()
				var mtype *mimetype.MIME
				if detectMime {
					mtype, err = mimetype.DetectFile(path)
					if err != nil {
						return errors.Wrap(err, "mime type error")
					}
				}
				return upload(destName, info, &progressReader{R: f, Add: func(i int) {
					b.uploadedSize += int64(i)
				}}, mtype)
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

func (b *Base) calculateSize() error {
	var size int64
	err := filepath.Walk(b.src, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	b.totalSize = size
	return err
}

func (b *Base) Progress() float64 {
	return float64(b.uploadedSize) / float64(b.totalSize)
}

type progressReader struct {
	R   io.Reader
	Add func(int)
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.R.Read(p)
	pr.Add(n)
	return n, err
}
