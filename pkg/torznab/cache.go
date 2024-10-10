package torznab

import (
	"polaris/log"
	"polaris/pkg/cache"
	"time"
)

var cc = cache.NewCache[string, []Result](time.Minute * 30)

func CleanCache() {
	log.Debugf("clean all torznab caches")
	cc = cache.NewCache[string, []Result](time.Minute * 30)
}