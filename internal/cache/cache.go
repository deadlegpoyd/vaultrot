// Package cache provides a lightweight in-memory TTL cache for secret values,
// reducing redundant reads from remote backends during a rotation run.
package cache

import (
	"sync"
	"time"
)

// entry holds a cached value alongside its expiry time.
type entry struct {
	value     string
	expiresAt time.Time
}

// Cache is a thread-safe, TTL-based in-memory store.
type Cache struct {
	mu      sync.RWMutex
	items   map[string]entry
	ttl     time.Duration
	nowFunc func() time.Time // injectable for testing
}

// New creates a Cache whose entries expire after ttl.
// A ttl of zero means entries never expire.
func New(ttl time.Duration) *Cache {
	return &Cache{
		items:   make(map[string]entry),
		ttl:     ttl,
		nowFunc: time.Now,
	}
}

// Set stores value under key, overwriting any existing entry.
func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var exp time.Time
	if c.ttl > 0 {
		exp = c.nowFunc().Add(c.ttl)
	}
	c.items[key] = entry{value: value, expiresAt: exp}
}

// Get returns the cached value for key and whether it was found and still valid.
func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	e, ok := c.items[key]
	if !ok {
		return "", false
	}
	if !e.expiresAt.IsZero() && c.nowFunc().After(e.expiresAt) {
		return "", false
	}
	return e.value, true
}

// Delete removes the entry for key, if present.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Flush removes all entries from the cache.
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]entry)
}

// Len returns the number of entries currently held (including expired ones
// that have not yet been evicted).
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}
