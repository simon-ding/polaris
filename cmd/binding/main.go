package main

import "C"
import (
	"polaris/internal/db"
	"polaris/log"
	"polaris/internal/biz/server"
)

func main() {}

var srv *server.Server
var port int

//export Start
func Start() (C.int, *C.char) {
	if srv != nil {
		return C.int(port), nil
	}
	log.InitLogger(true)

	log.Infof("------------------- Starting Polaris ---------------------")
	dbClient, err := db.Open()
	if err != nil {
		log.Panicf("init db error: %v", err)
		return C.int(0), C.CString(err.Error())
	}

	s := server.NewServer(dbClient)
	if p, err := s.Start(""); err != nil {
		return C.int(0), C.CString(err.Error())
	} else {
		port = p
		srv = s
		return C.int(p), C.CString("")
	}

}

//export Stop
func Stop() {
	if srv != nil {
		srv.Stop()
	}
	srv = nil
}
