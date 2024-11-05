package uploader

import (
	"fmt"
	"io"
	"os"
	"polaris/pkg/utils"
	"sync/atomic"
	"time"
)

type StreamWriter interface {
	WriteStream(path string, stream io.Reader, _ os.FileMode) error
}

type Uploader struct {
	sw       StreamWriter
	progress atomic.Int64
	dir      string
	size     int64
}

func NewUploader(dir string, sw StreamWriter) (*Uploader, error) {
	size, err := utils.DirSize(dir)
	if err != nil {
		return nil, err
	}
	return &Uploader{sw: sw, dir: dir, size: size, progress: atomic.Int64{}}, nil
}

func (u *Uploader) Upload() error {

	return nil
}

type ProgressReader struct {
	Reader   io.Reader
	Progress atomic.Int64
	Size     int64
	Name     string
	Once     bool
	Done     atomic.Bool
}

func (progressReader *ProgressReader) NewLoop() {
	ticker := time.NewTicker(time.Second)
	var op int64
	for range ticker.C {
		p := progressReader.Progress.Load()
		KB := (p - op) / 1024
		var percent int64
		if progressReader.Size != 0 {
			percent = p * 100 / progressReader.Size
		} else {
			percent = 100
		}
		if KB < 1024 {
			fmt.Printf("%s: %dKB/s %d%%\n", progressReader.Name, KB, percent)
		} else {
			fmt.Printf("%s: %.2fMB/s %d%%\n", progressReader.Name, float64(KB)/1024, percent)
		}

		if progressReader.Done.Load() {
			ticker.Stop()
			return
		}
	}
}

func (progressReader *ProgressReader) Read(p []byte) (int, error) {
	n, err := progressReader.Reader.Read(p)
	progressReader.Progress.Add(int64(n))
	if !progressReader.Once {
		progressReader.Once = true
		go progressReader.NewLoop()
	}
	if err != nil {
		progressReader.Done.Store(true)
	}
	return n, err
}

func (progressReader *ProgressReader) Close() error {
	progressReader.Done.Store(true)
	return nil
}
