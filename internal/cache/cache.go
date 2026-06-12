package cache

import (
	"context"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
)

type cacheEntry struct {
	Data any
	Ttl  time.Duration
	At   time.Time
}

type Cache struct {
	mu sync.Mutex

	CachesEntries map[string]cacheEntry
	maxEntries    int
}

func NewCache(ctx context.Context, maxEntries int) *Cache {
	cache := &Cache{
		CachesEntries: make(map[string]cacheEntry),
		maxEntries:    maxEntries,
	}

	go cache.EntriesCleaner(ctx)

	return cache
}

func (c *Cache) Get(key string) any {
	defer c.mu.Unlock()
	c.mu.Lock()

	entry, exists := c.CachesEntries[key]
	if !exists {
		return nil
	} else if c.isEntryExpired(entry) {
		delete(c.CachesEntries, key)
		return nil
	}

	return entry.Data
}

func (c *Cache) Set(key string, data any, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.maxEntries > 0 && len(c.CachesEntries) >= c.maxEntries {
		var oldestKey string
		var oldestTime time.Time

		for k, entry := range c.CachesEntries {
			if oldestKey == "" || entry.At.Before(oldestTime) {
				oldestKey, oldestTime = k, entry.At
			}
		}
		if oldestKey != "" {
			delete(c.CachesEntries, oldestKey)
		}
	}

	c.CachesEntries[key] = cacheEntry{
		Data: data,
		Ttl:  ttl,
		At:   time.Now().UTC(),
	}
}

func (c *Cache) SetHeader(ctx fiber.Ctx, key string) {
	defer c.mu.Unlock()
	c.mu.Lock()

	entry, exists := c.CachesEntries[key]
	if !exists {
		return
	}

	if !c.isEntryExpired(entry) {
		ctx.Set("X-Cache", entry.At.Format(time.RFC3339))
	}

}

func (c *Cache) EntriesCleaner(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			for key, entry := range c.CachesEntries {
				if c.isEntryExpired(entry) {
					delete(c.CachesEntries, key)
				}
			}
			c.mu.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

func (c *Cache) isEntryExpired(entry cacheEntry) bool {
	return time.Since(entry.At) > entry.Ttl
}
