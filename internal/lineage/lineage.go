// Package lineage tracks the rotation history of secrets, recording each
// rotation event with its predecessor so callers can reconstruct a full
// audit trail for any given secret key.
package lineage

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Entry describes a single rotation event for a secret.
type Entry struct {
	Backend     string
	Key         string
	RotatedAt   time.Time
	GeneratedBy string // e.g. "auto", "manual", "scheduled"
	PreviousID  string // opaque ID of the preceding entry, empty for first rotation
	ID          string // opaque identifier for this entry
}

// Tracker stores rotation lineage in memory.
type Tracker struct {
	mu      sync.RWMutex
	entries map[string][]Entry // keyed by "backend/key"
	nowFn   func() time.Time
}

// New returns a new Tracker.
func New() *Tracker {
	return &Tracker{
		entries: make(map[string][]Entry),
		nowFn:   time.Now,
	}
}

// Record appends a rotation entry for the given backend and key.
// generatedBy describes the rotation trigger (e.g. "auto", "manual").
func (t *Tracker) Record(backend, key, generatedBy string) (Entry, error) {
	if backend == "" {
		return Entry{}, errors.New("lineage: backend must not be empty")
	}
	if key == "" {
		return Entry{}, errors.New("lineage: key must not be empty")
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	mapKey := backend + "/" + key
	chain := t.entries[mapKey]

	var prevID string
	if len(chain) > 0 {
		prevID = chain[len(chain)-1].ID
	}

	e := Entry{
		Backend:     backend,
		Key:         key,
		RotatedAt:   t.nowFn().UTC(),
		GeneratedBy: generatedBy,
		PreviousID:  prevID,
		ID:          makeID(backend, key, len(chain)),
	}

	t.entries[mapKey] = append(chain, e)
	return e, nil
}

// History returns all recorded entries for the given backend and key,
// oldest first. Returns nil if no history exists.
func (t *Tracker) History(backend, key string) []Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()

	chain := t.entries[backend+"/"+key]
	if len(chain) == 0 {
		return nil
	}
	out := make([]Entry, len(chain))
	copy(out, chain)
	return out
}

// Len returns the total number of rotation events recorded across all keys.
func (t *Tracker) Len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()

	total := 0
	for _, chain := range t.entries {
		total += len(chain)
	}
	return total
}

func makeID(backend, key string, index int) string {
	return fmt.Sprintf("%s/%s#%d", backend, key, index)
}
