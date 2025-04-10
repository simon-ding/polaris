package main

import "C"
import (
	"os"
	"polaris/cmd"
	"polaris/log"
)

func main() {}

//export Start
func Start() {
	cmd.Start(true)
}

//export Stop
func Stop() {
	log.Infof("stop polaris")
	os.Exit(0)
}
