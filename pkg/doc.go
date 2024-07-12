package pkg

type Torrent interface {
	Name() string
	Progress() int
	Stop() error
	Start() error
	Remove(deleteData bool) error
	Save() string
	Exists() bool
}


type Storage interface {

}