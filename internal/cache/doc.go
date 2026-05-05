// Package cache implements a lightweight, thread-safe, TTL-based in-memory
// cache used by vaultrot to avoid redundant reads from secret backends during
// a single rotation run.
//
// Usage:
//
//	c := cache.New(30 * time.Second)
//	c.Set("prod/db/password", "s3cr3t")
//
//	if val, ok := c.Get("prod/db/password"); ok {
//		// use cached value
//	}
//
// A TTL of zero disables expiry so entries live for the lifetime of the
// cache object. Call Flush to reset the cache between rotation cycles when
// running in scheduled (daemon) mode.
package cache
