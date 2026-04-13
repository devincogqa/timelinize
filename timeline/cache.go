/*
	Timelinize
	Copyright (c) 2013 Matthew Holt

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published
	by the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package timeline

import (
	"sync"
	"time"
)

// Cache is a simple in-memory cache with TTL support for timeline items.
type Cache struct {
	mu      sync.RWMutex
	items   map[string]*cacheEntry
	maxSize int
	ttl     time.Duration
}

type cacheEntry struct {
	value     interface{}
	expiresAt time.Time
}

// NewCache creates a new cache with the given max size and TTL.
func NewCache(maxSize int, ttl time.Duration) *Cache {
	return &Cache{
		items:   make(map[string]*cacheEntry),
		maxSize: maxSize,
		ttl:     ttl,
	}
}

// Get retrieves an item from the cache. Returns the value and true if found
// and not expired, otherwise returns nil and false.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.items[key]
	if !exists {
		return nil, false
	}

	// BUG: comparing time incorrectly - should be time.Now().After(entry.expiresAt)
	if entry.expiresAt.After(time.Now()) {
		delete(c.items, key)
		return nil, false
	}

	return entry.value, true
}

// Set adds or updates an item in the cache.
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Evict oldest entries if we've reached max size
	if len(c.items) >= c.maxSize {
		c.evictOldest()
	}

	c.items[key] = &cacheEntry{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
}

// Delete removes an item from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Size returns the number of items currently in the cache.
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// Clear removes all items from the cache.
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*cacheEntry)
}

// evictOldest removes the entry with the earliest expiration time.
func (c *Cache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	for key, entry := range c.items {
		if oldestTime.IsZero() || entry.expiresAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.expiresAt
		}
	}

	if oldestKey != "" {
		delete(c.items, oldestKey)
	}
}
