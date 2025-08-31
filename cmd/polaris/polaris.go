package main

import (
	"flag"
	"fmt"
	"os"
	"polaris/db"
	"polaris/log"
	"polaris/pkg/utils"
	"polaris/server"
)

func main() {
	port := flag.Int("port", 3322, "port to listen on")
	flag.Parse()

	if os.Getenv("GIN_MODE") == "release" {
		log.InitLogger(true)
	}

	log.Infof("------------------- Starting Polaris ---------------------")
	dbClient, err := db.Open()
	if err != nil {
		log.Panicf("init db error: %v", err)
	}
	if !utils.IsRunningInDocker() {
		go utils.OpenURL(fmt.Sprintf("http://127.0.0.1:%d", *port))
	}

	s := server.NewServer(dbClient)
	if _, err := s.Start(fmt.Sprintf(":%d", *port)); err != nil {
		log.Errorf("server start error: %v", err)
	}
	select {} //wait indefinitely
}
