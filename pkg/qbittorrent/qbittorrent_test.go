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
	log.Infof("new client success: %v", c)
	port, err := c.GetListenPort()
	if err != nil {
		log.Errorf("get listen port error: %v", err)
		t.Fail()
	} else {
		log.Infof("listen port: %d", port)
		err := c.SetListenPort(port + 1)
		if err!= nil {
			log.Errorf("set listen port error: %v", err)
			t.Fail()	
		}
	}
}
