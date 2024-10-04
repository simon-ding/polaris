package qbittorrent

import (
	"polaris/log"
	"testing"
)

func Test1(t *testing.T) {
	c, err := NewClient("http://10.0.0.8:8081/", "", "")
	if err != nil {
		log.Errorf("new client error: %v", err)
		t.Fail()
	}
	all, err := c.GetAll()
	for _, t := range all {
		name, _ := t.Name()
		log.Infof("torrent: %+v", name)
	}
	
}
