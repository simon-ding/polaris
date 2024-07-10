package pkg

type Torrent interface {
	Name() string
	Progress() int
	Stop() error
	Start() error
	Remove() error
	Save() string
}


type Storage interface {

}