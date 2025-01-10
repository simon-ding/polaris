package prowlarr

import (
	"polaris/log"
	"testing"
)

func Test111(t *testing.T) {
	c := New("", "http://10.0.0.8:9696/")
	apis , err := c.GetIndexers()
	log.Infof("errors: %v", err)
	log.Infof("indexers: %+v", apis[0])
}
