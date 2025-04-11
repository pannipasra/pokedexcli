package pokecache

import (
	"sync"
	"time"
)

// Cache represents an in-memory cache with expiration
type Cache struct {
	cache map[string]cacheEntry
	mutex sync.RWMutex
}

// cacheEntry represents a single entry in the cache
type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// =====================================================
// =====================================================

// creates a new cache with a configurable interval
func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		cache: make(map[string]cacheEntry),
	}

	// Start a background goroutine to clean up expired entries
	go c.reapLoop(interval)

	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.cache[key]
	if !exists {
		return nil, false
	}

	return entry.val, true
}

// reapLoop periodically removes expired entries from the cache
func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.reap(interval)
	}
}

// reap removes entries that have been in the cache longer than the specified duration
func (c *Cache) reap(interval time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()
	for key, entry := range c.cache {
		if now.Sub(entry.createdAt) > interval {
			delete(c.cache, key)
		}
	}
}
