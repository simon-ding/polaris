package alist

import (
	"os"
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

	f, err := os.Open("/Users/simonding/Downloads/Steam Link_1.3.9_APKPure.apk")
	if err != nil {
		log.Errorf("openfile: %v", err)
		t.Fail()
	} else {
		defer f.Close()
		ss, _ := f.Stat()
		log.Infof("upload file size %d", ss.Size())
		info, err := c.UploadStream(f, ss.Size(), "/aliyun/Steam Link_1.3.9_APKPure.apk")
		if err != nil {
			log.Errorf("upload error: %v", err)
			t.Fail()
		} else {
			log.Infof("upload success: %+v", info)
		}
	}
}
