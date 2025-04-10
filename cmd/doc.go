package cmd

import (
	"os"
	"polaris/db"
	"polaris/log"
	"polaris/server"
)

func Start(sharedLib bool) {
	if sharedLib || os.Getenv("GIN_MODE") == "release" {
		log.InitLogger(true)
	} else {
		log.InitLogger(false)
	}

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
