package alist

import (
	"polaris/log"
	"testing"
)

func TestLogin(t *testing.T) {
	c := New(&Config{
		URL: "http://10.0.0.8:5244/",
		Username: "",
		Password: "",
	})
	cre, err := c.Login()
	if err != nil {
		log.Errorf("login fail: %v", err)
		t.Fail()
	} else {
		log.Errorf("login success: %s", cre)
	}
	info, err := c.Ls("/aliyun")
	if err != nil {
		log.Errorf("ls fail: %v", err)
		t.Fail()
	} else {
		log.Infof("ls results: %+v", info)
	}
	err = c.Mkdir("/aliyun/test1")
	log.Errorf("mkdir: %v", err)
}
