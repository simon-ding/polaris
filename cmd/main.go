package main

import (
	"polaris/db"
	"polaris/log"
	"polaris/server"
)

func main() {
	log.Infof("------------------- Starting Polaris ---------------------")
	dbClient, err := db.Open()
	if err != nil {
		log.Panicf("init db error: %v", err)
	}

	s := server.NewServer(dbClient)
	if err := s.Serve(); err != nil {
		log.Errorf("server start error: %v", err)
	}
}
