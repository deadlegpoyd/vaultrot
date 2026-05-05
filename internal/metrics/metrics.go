// Package metrics provides lightweight counters and gauges for tracking
// rotation run statistics across backends.
package metrics

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// Collector accumulates rotation metrics for a single run.
type Collector struct {
	mu        sync.Mutex
	counters  map[string]int64
	durations map[string]time.Duration
	start     time.Time
}

// New returns an initialised Collector with the run start time set to now.
func New() *Collector {
	return &Collector{
		counters:  make(map[string]int64),
		durations: make(map[string]time.Duration),
		start:     time.Now().UTC(),
	}
}

// Inc increments the named counter by 1.
func (c *Collector) Inc(name string) {
	c.Add(name, 1)
}

// Add adds delta to the named counter.
func (c *Collector) Add(name string, delta int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counters[name] += delta
}

// RecordDuration stores a duration sample under name, overwriting any
// previous value.
func (c *Collector) RecordDuration(name string, d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.durations[name] = d
}

// Get returns the current value of a counter (0 if not set).
func (c *Collector) Get(name string) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.counters[name]
}

// GetDuration returns the stored duration for name (0 if not set).
func (c *Collector) GetDuration(name string) time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.durations[name]
}

// Elapsed returns the time elapsed since the Collector was created.
func (c *Collector) Elapsed() time.Duration {
	return time.Since(c.start)
}

// Print writes a human-readable summary of all metrics to w.
func (c *Collector) Print(w io.Writer) {
	c.mu.Lock()
	defer c.mu.Unlock()

	fmt.Fprintf(w, "=== metrics ===\n")
	for k, v := range c.counters {
		fmt.Fprintf(w, "  %-30s %d\n", k, v)
	}
	for k, v := range c.durations {
		fmt.Fprintf(w, "  %-30s %s\n", k, v.Round(time.Millisecond))
	}
	fmt.Fprintf(w, "  %-30s %s\n", "elapsed", time.Since(c.start).Round(time.Millisecond))
}
