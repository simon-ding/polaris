package utils

import (
	"polaris/log"
	"testing"
)

func TestLink2Magnet(t *testing.T) {
	s, err := Link2Magnet("https://api.m-team.cc/api/rss/dlv2?useHttps=true&type=ipv6&sign=2ecfdb9d1317fce1edc123d024be1d65&t=1738309528&tid=900434&uid=346577")
	log.Errorf("%v", err)
	log.Infof("%v", s)
}
