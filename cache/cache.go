package cache

import (
	"../byteView"
	"../lru"
	"sync"
)

type Cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func NewMainCache(cacheBytes int64) *Cache {
	return  &Cache{
		lru:        lru.NewCache(cacheBytes, nil),
		cacheBytes: cacheBytes,
	}
}

func (c *Cache) Save(key string, value byteView.ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		// 这里回调函数暂时没用
		c.lru = lru.NewCache(c.cacheBytes, nil)
	}
	c.lru.Save(key, value)
}

func (c *Cache) Get(key string) (value byteView.ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock() //精简代码，使其在函数执行完之间执行此行代码
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(byteView.ByteView), ok
	}
	return
}

