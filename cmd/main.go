package main

import (
	"polaris/db"
	"polaris/log"
	"polaris/pkg/utils"
	"polaris/server"
	"time"
)

func main() {
	log.Infof("------------------- Starting Polaris ---------------------")

	//utils.MaxPermission()

	dbClient, err := db.Open()
	if err != nil {
		log.Panicf("init db error: %v", err)
	}

	go func() {
		time.Sleep(2 * time.Second)
		if err := utils.OpenURL("http://127.0.0.1:8080"); err != nil {
			log.Errorf("open url error: %v", err)
		}

	}()
	s := server.NewServer(dbClient)
	if err := s.Serve(); err != nil {
		log.Errorf("server start error: %v", err)
	}
}
