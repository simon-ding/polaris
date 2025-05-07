package main

import (
	"net"
	"polaris/log"
	"polaris/pkg/nat"
)

func main() {
	// This is a placeholder for the main function.
	// The actual implementation will depend on the specific requirements of the application.
	src, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := src.Accept()
		if err != nil {
			panic(err)
		}
		log.Infof("new connection: %+v", conn)
		dest, err := net.Dial("tcp", "10.0.0.8:8080")
		if err != nil {
			panic(err)
		}
	
		go nat.ReverseProxy(conn, dest)
	}
	select {}
}
