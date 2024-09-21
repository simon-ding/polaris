package pkg

type Torrent interface {
	Name() (string, error)
	Progress() (int, error)
	Stop() error
	Start() error
	Remove() error
	Save() string
	Exists() bool
	SeedRatio() (float64, error)
}


type Storage interface {

}