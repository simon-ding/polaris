package main

import (
	"os"
	"polaris/internal/db"
	"polaris/log"
	"polaris/internal/biz/server"
)

func main() {
	if os.Getenv("GIN_MODE") == "release" {
		log.InitLogger(true)
	}

	log.Infof("------------------- Starting Polaris ---------------------")
	dbClient, err := db.Open()
	if err != nil {
		log.Panicf("init db error: %v", err)
	}

	s := server.NewServer(dbClient)
	if _, err := s.Start(":8080"); err != nil {
		log.Errorf("server start error: %v", err)
	}
	select {} //wait indefinitely
}
