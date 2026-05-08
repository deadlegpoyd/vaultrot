// Package changelog records a human-readable history of secret rotation
// events, providing an ordered log of what changed, when, and by whom.
package changelog

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// Entry represents a single rotation event in the changelog.
type Entry struct {
	OccurredAt time.Time
	Backend    string
	Key        string
	Actor      string
	Action     string
	DryRun     bool
	Note       string
}

// Changelog holds an ordered list of rotation entries.
type Changelog struct {
	mu      sync.RWMutex
	entries []Entry
	clock   func() time.Time
}

// New returns a new Changelog. An optional clock function may be injected
// for testing; pass nil to use time.Now.
func New(clock func() time.Time) *Changelog {
	if clock == nil {
		clock = time.Now
	}
	return &Changelog{clock: clock}
}

// Record appends a new entry to the changelog.
func (c *Changelog) Record(backend, key, actor, action, note string, dryRun bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = append(c.entries, Entry{
		OccurredAt: c.clock().UTC(),
		Backend:    backend,
		Key:        key,
		Actor:      actor,
		Action:     action,
		DryRun:     dryRun,
		Note:       note,
	})
}

// Entries returns a copy of all recorded entries.
func (c *Changelog) Entries() []Entry {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]Entry, len(c.entries))
	copy(out, c.entries)
	return out
}

// Len returns the number of recorded entries.
func (c *Changelog) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}

// Write formats all entries as plain text and writes them to w.
func (c *Changelog) Write(w io.Writer) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, e := range c.entries {
		dryTag := ""
		if e.DryRun {
			dryTag = " [dry-run]"
		}
		line := fmt.Sprintf("%s%s  %-10s  %-8s  %s/%s",
			e.OccurredAt.Format(time.RFC3339), dryTag,
			e.Action, e.Actor, e.Backend, e.Key)
		if e.Note != "" {
			line += "  (" + e.Note + ")"
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}
