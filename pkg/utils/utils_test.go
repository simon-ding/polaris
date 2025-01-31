package utils

import (
	"polaris/log"
	"testing"
)

func TestLink2Magnet(t *testing.T) {
	s, err := Link2Magnet("")
	log.Errorf("%v", err)
	log.Infof("%v", s)
}
