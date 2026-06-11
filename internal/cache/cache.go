package cache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	Data any
	Ttl  time.Time
}

type Cache struct {
	mu sync.Mutex

	CachesEntries map[string]CacheEntry
}

func NewCache() *Cache {
	return &Cache{
		CachesEntries: make(map[string]CacheEntry),
	}
}

func (c *Cache) Get(key string) any {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.CachesEntries[key]
	if !exists || time.Now().After(entry.Ttl) {
		return nil
	}

	return entry.Data
}

func (c *Cache) Set(key string, data any, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.CachesEntries[key] = CacheEntry{
		Data: data,
		Ttl:  time.Now().Add(ttl),
	}
}
