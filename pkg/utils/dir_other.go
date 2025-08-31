//go:build !lib

package utils

import (
	"os"
	"path/filepath"
)

func GetUserDataDir() string {
	if IsRunningInDocker() {
		return "./data"
	}
	homeDir, _ := os.UserHomeDir()
	dir := filepath.Join(homeDir, ".polaris")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	return dir
}
