package torznab

import (
	"polaris/log"
	"polaris/pkg/utils"

	"github.com/robfig/cron"

	"time"
)

var cache utils.Map[string, TimedResponse] = utils.Map[string, TimedResponse]{}

type TimedResponse struct {
	Response
	T time.Time
}

func init() {
	cr := cron.New()
	cr.AddFunc("@ervery 1m", func() {
		cache.Range(func(key string, value TimedResponse) bool {
			if time.Since(value.T) > 30*time.Minute {
				log.Debugf("delete old cache: %v", key)
				cache.Delete(key)
			}
			return true
		})
	})
	cr.Start()
}
