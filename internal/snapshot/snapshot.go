// Package snapshot captures and compares secret values at a point in time,
// enabling before/after diffing during rotation runs.
package snapshot

import (
	"fmt"
	"sync"
	"time"
)

// Entry holds a captured secret value alongside metadata.
type Entry struct {
	Backend   string
	Key       string
	Value     string
	CapturedAt time.Time
}

// Store holds a collection of named snapshots.
type Store struct {
	mu      sync.RWMutex
	entries map[string]Entry
	now     func() time.Time
}

// New returns an initialised Store.
func New() *Store {
	return &Store{
		entries: make(map[string]Entry),
		now:     time.Now,
	}
}

// Capture records the current value of a secret identified by backend and key.
func (s *Store) Capture(backend, key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[entryKey(backend, key)] = Entry{
		Backend:    backend,
		Key:        key,
		Value:      value,
		CapturedAt: s.now().UTC(),
	}
}

// Get retrieves a previously captured entry. Returns false if not found.
func (s *Store) Get(backend, key string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[entryKey(backend, key)]
	return e, ok
}

// Len returns the number of captured entries.
func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.entries)
}

// All returns a copy of all captured entries.
func (s *Store) All() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, e)
	}
	return out
}

func entryKey(backend, key string) string {
	return fmt.Sprintf("%s::%s", backend, key)
}
