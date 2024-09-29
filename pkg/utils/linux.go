//go:build linux
// +build linux

package utils

import (
	"golang.org/x/sys/unix"
	"math"
	"runtime"
	"syscall"
)

func AvailableSpace(dir string) uint64 {
	if runtime.GOOS != "linux" {
		return math.MaxUint64
	}
	var stat unix.Statfs_t

	unix.Statfs(dir, &stat)
	return stat.Bavail * uint64(stat.Bsize)
}

func MaxPermission() {
	syscall.Umask(0) //max permission 0777
}
