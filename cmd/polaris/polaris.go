package main

import (
	"flag"
	"fmt"
	"os"
	"polaris/db"
	"polaris/log"
	"polaris/server"
)

func main() {
	port := flag.Int("port", 8080, "port to listen on")
	flag.Parse()

	if os.Getenv("GIN_MODE") == "release" {
		log.InitLogger(true)
	}

	log.Infof("------------------- Starting Polaris ---------------------")
	dbClient, err := db.Open()
	if err != nil {
		log.Panicf("init db error: %v", err)
	}

	s := server.NewServer(dbClient)
	if _, err := s.Start(fmt.Sprintf(":%d", *port)); err != nil {
		log.Errorf("server start error: %v", err)
	}
	select {} //wait indefinitely
}
