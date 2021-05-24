package WinCache

import (
	"./byteView"
	"./lru"
	"sync"
)

type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *cache) save(key string, value byteView.ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.NewCache(c.cacheBytes, nil)
	}
	c.lru.Save(key, value)
}

func (c *cache) get(key string) (value byteView.ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock() // 精简代码，使其在最后函数进行关闭
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(byteView.ByteView), ok
	}
	return
}

