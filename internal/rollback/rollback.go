// Package rollback provides snapshot and restore functionality for secrets
// rotated by vaultrot, enabling recovery when a rotation partially fails.
package rollback

import (
	"fmt"
	"sync"
	"time"
)

// Snapshot holds the previous value of a secret before rotation.
type Snapshot struct {
	Backend   string
	SecretKey string
	OldValue  string
	CapturedAt time.Time
}

// Store holds snapshots indexed by a composite key.
type Store struct {
	mu        sync.RWMutex
	snapshots map[string]Snapshot
}

// New returns an initialised rollback Store.
func New() *Store {
	return &Store{
		snapshots: make(map[string]Snapshot),
	}
}

// Record saves a snapshot of a secret's current value before it is rotated.
func (s *Store) Record(backend, key, oldValue string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	composite := compositeKey(backend, key)
	s.snapshots[composite] = Snapshot{
		Backend:    backend,
		SecretKey:  key,
		OldValue:   oldValue,
		CapturedAt: time.Now().UTC(),
	}
}

// Get retrieves a previously recorded snapshot. The second return value
// indicates whether the snapshot exists.
func (s *Store) Get(backend, key string) (Snapshot, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snap, ok := s.snapshots[compositeKey(backend, key)]
	return snap, ok
}

// Keys returns all composite keys currently held in the store.
func (s *Store) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]string, 0, len(s.snapshots))
	for k := range s.snapshots {
		keys = append(keys, k)
	}
	return keys
}

// Len returns the number of snapshots stored.
func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.snapshots)
}

func compositeKey(backend, key string) string {
	return fmt.Sprintf("%s::%s", backend, key)
}
