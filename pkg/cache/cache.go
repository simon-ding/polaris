package cache

import (
	"polaris/log"
	"polaris/pkg/utils"
	"time"

	"github.com/robfig/cron"
)

func NewCache[T comparable, S any](timeout time.Duration) *Cache[T, S] {
	c := &Cache[T, S]{
		m:       utils.Map[T, inner[S]]{},
		timeout: timeout,
		cr:      cron.New(),
	}

	c.cr.AddFunc("@ervery 1m", func() {
		c.m.Range(func(key T, value inner[S]) bool {
			if time.Since(value.t) > c.timeout {
				log.Debugf("delete old cache: %v", key)
				c.m.Delete(key)

			}
			return true
		})
	})
	c.cr.Start()
	return c
}

type Cache[T comparable, S any] struct {
	m       utils.Map[T, inner[S]]
	timeout time.Duration
	cr      *cron.Cron
}

type inner[S any] struct {
	t time.Time
	s S
}

func (c *Cache[T, S]) Set(key T, value S) {
	c.m.Store(key, inner[S]{t: time.Now(), s: value})
}

func (c *Cache[T, S]) Get(key T) (S, bool) {
	v, ok := c.m.Load(key)
	if !ok {
		return getZero[S](), ok
	}
	return v.s, ok
}

func getZero[T any]() T {
	var result T
	return result
}
