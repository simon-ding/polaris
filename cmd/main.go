package main

import (
	"polaris/db"
	"polaris/log"
	"polaris/server"
	"syscall"
)

func main() {
	log.Infof("------------------- Starting Polaris ---------------------")

	syscall.Umask(0) //max permission 0777

	dbClient, err := db.Open()
	if err != nil {
		log.Panicf("init db error: %v", err)
	}

	s := server.NewServer(dbClient)
	if err := s.Serve(); err != nil {
		log.Errorf("server start error: %v", err)
	}
}
