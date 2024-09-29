//go:build !linux
// +build !linux

package utils

import (
	"math"
)

func AvailableSpace(dir string) uint64 {
	return math.MaxUint64
}

func MaxPermission() {
	return
}
