//go:build lib

package utils

import (
	"os"
	"path/filepath"
)


func GetUserDataDir() string  {
	d, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	d = filepath.Join(d, ".polaris")
	if _, err := os.Stat(d); os.IsNotExist(err) {
		os.MkdirAll(d, os.ModePerm)
	}
	return d
}