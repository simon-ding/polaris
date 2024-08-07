package torznab

import (
	"polaris/pkg/cache"
	"time"
)

var cc = cache.NewCache[string, Response](time.Minute * 30)
