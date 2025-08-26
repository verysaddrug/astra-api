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
	c.store.Store(key, CacheItem{Value: value, Expiration: time.Now().Add(c.ttl).UnixNano()})
}

func (c *Cache) Get(key string) (interface{}, bool) {
	item, ok := c.store.Load(key)
	if !ok {
		return nil, false
	}
	ci := item.(CacheItem)
	if time.Now().UnixNano() > ci.Expiration {
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
