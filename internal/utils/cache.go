package utils

import (
	"sync"
	"time"
)

// CacheEntry represents a single cached item with an expiration time.
type CacheEntry struct {
	Value     interface{}
	ExpiresAt time.Time
}

// Cache is a simple in-memory cache with TTL support.
type Cache struct {
	mu      sync.Mutex
	entries map[string]CacheEntry
	ttl     time.Duration
}

// NewCache creates a new Cache with the given TTL.
func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		entries: make(map[string]CacheEntry),
		ttl:     ttl,
	}
}

// Set adds a value to the cache with the configured TTL.
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = CacheEntry{
		Value:     value,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

// Get retrieves a value from the cache. Returns the value and true if found
// and not expired, otherwise returns nil and false.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	if time.Now().After(entry.ExpiresAt) {
		// BUG: We detect the entry is expired but don't delete it from the map,
		// causing a memory leak over time as expired entries accumulate.
		return nil, false
	}
	return entry.Value, true
}

// Delete removes a value from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Size returns the number of entries in the cache (including expired ones).
func (c *Cache) Size() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.entries)
}
