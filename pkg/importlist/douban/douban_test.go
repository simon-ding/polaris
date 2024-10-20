package douban

import (
	"polaris/log"
	"testing"
)

func TestParseDoulist(t *testing.T) {
	r, err := ParseDoulist("https://www.douban.com/doulist/81580/")
	log.Info(r, err)
}
