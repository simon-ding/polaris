package pkg

import "io/fs"

type Torrent interface {
	Name() (string, error)
	Progress() (int, error)
	Stop() error
	Start() error
	Remove() error
	//Save() string
	Exists() bool
	SeedRatio() (float64, error)
	GetHash() string
	//Reload() error
	WalkFunc() func(fn func(path string, info fs.FileInfo) error) error
}

type Downloader interface {
	GetAll() ([]Torrent, error)
	Download(link, hash, dir string) (Torrent, error)
}
