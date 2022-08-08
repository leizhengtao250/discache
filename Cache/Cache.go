package Cache

import (
	"DisCache/ByteView"
	"DisCache/LRU"
	"sync"
)

//cache 是对LRU的封装
type cache struct {
	mu         sync.Mutex
	LRU        *LRU.Cache
	cacheBytes int64
}

func (c *cache) add(key string, value ByteView.ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.LRU == nil {
		c.LRU = LRU.New(c.cacheBytes, nil)
	}
	c.LRU.Add(key, value)
}

func (c *cache) GetC(key string) (value ByteView.ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.LRU == nil {
		return
	}
	if ele, ok := c.LRU.Get(key); ok {
		return ele.(ByteView.ByteView), ok
	}
	return
}
