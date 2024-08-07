package cache

import (
	"polaris/log"
	"polaris/pkg/utils"
	"time"
)

func NewCache[T comparable, S any](timeout time.Duration) *Cache[T, S] {
	c := &Cache[T, S]{
		m:       utils.Map[T, inner[S]]{},
		timeout: timeout,
	}

	return c
}

type Cache[T comparable, S any] struct {
	m       utils.Map[T, inner[S]]
	timeout time.Duration
}

type inner[S any] struct {
	t time.Time
	s S
}

func (c *Cache[T, S]) Set(key T, value S) {
	c.m.Store(key, inner[S]{t: time.Now(), s: value})
}

func (c *Cache[T, S]) Get(key T) (S, bool) {
	c.m.Range(func(key T, value inner[S]) bool {
		if time.Since(value.t) > c.timeout {
			log.Debugf("delete old cache: %v", key)
			c.m.Delete(key)

		}
		return true
	})

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
