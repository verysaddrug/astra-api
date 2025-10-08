package cache

import (
	"sync"
	"time"
)

type CacheItem struct {
	Value      interface{}
	Expiration int64
}

type Cache struct {
	store sync.Map
	ttl   time.Duration
}

func NewCache(ttl time.Duration) *Cache {
	return &Cache{ttl: ttl}
}

func (c *Cache) Set(key string, value interface{}) {
	var expiration int64
	if c.ttl > 0 {
		expiration = time.Now().Add(c.ttl).UnixNano()
	} else {
		expiration = 0 // No expiration
	}
	c.store.Store(key, CacheItem{Value: value, Expiration: expiration})
}

func (c *Cache) Get(key string) (interface{}, bool) {
	item, ok := c.store.Load(key)
	if !ok {
		return nil, false
	}
	ci := item.(CacheItem)
	if ci.Expiration > 0 && time.Now().UnixNano() > ci.Expiration {
		c.store.Delete(key)
		return nil, false
	}
	return ci.Value, true
}

func (c *Cache) Invalidate(key string) {
	c.store.Delete(key)
}

func (c *Cache) InvalidateAll() {
	c.store = sync.Map{}
}
