package douban

import (
	"polaris/log"
	"testing"
)

func TestParseDoulist(t *testing.T) {
	r, err := ParseDoulist("https://www.douban.com/doulist/81580/")
	log.Info(r, err)
}


func Test111(t *testing.T) {
	d := NewDoubanWishlist("69894889")
	_, err := d.GetWishlist(1)
	log.Infof("err: %v", err)
}