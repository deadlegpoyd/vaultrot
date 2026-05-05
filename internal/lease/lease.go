// Package lease tracks secret lease expiry and renewal windows.
package lease

import (
	"errors"
	"sync"
	"time"
)

// Entry holds metadata about a single secret lease.
type Entry struct {
	SecretName string
	Backend    string
	IssuedAt   time.Time
	ExpiresAt  time.Time
	Renewable  bool
}

// IsExpired reports whether the lease has passed its expiry time.
func (e Entry) IsExpired(now time.Time) bool {
	return now.After(e.ExpiresAt)
}

// DueForRenewal reports whether the lease is within the given threshold of expiry.
func (e Entry) DueForRenewal(now time.Time, threshold time.Duration) bool {
	return e.Renewable && now.After(e.ExpiresAt.Add(-threshold))
}

// Tracker stores and queries lease entries.
type Tracker struct {
	mu      sync.RWMutex
	entries map[string]Entry
	now     func() time.Time
}

// New returns a new Tracker.
func New() *Tracker {
	return &Tracker{
		entries: make(map[string]Entry),
		now:     time.Now,
	}
}

// Record stores or overwrites the lease entry for a secret.
func (t *Tracker) Record(e Entry) error {
	if e.SecretName == "" {
		return errors.New("lease: secret name must not be empty")
	}
	if e.ExpiresAt.IsZero() {
		return errors.New("lease: expiry time must be set")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries[e.SecretName] = e
	return nil
}

// Get returns the lease entry for the named secret.
func (t *Tracker) Get(name string) (Entry, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	e, ok := t.entries[name]
	return e, ok
}

// Expired returns all entries whose leases have expired.
func (t *Tracker) Expired() []Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	now := t.now()
	var out []Entry
	for _, e := range t.entries {
		if e.IsExpired(now) {
			out = append(out, e)
		}
	}
	return out
}

// DueForRenewal returns entries within threshold of expiry that are renewable.
func (t *Tracker) DueForRenewal(threshold time.Duration) []Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	now := t.now()
	var out []Entry
	for _, e := range t.entries {
		if e.DueForRenewal(now, threshold) {
			out = append(out, e)
		}
	}
	return out
}

// Remove deletes the lease entry for the named secret.
func (t *Tracker) Remove(name string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, name)
}

// Len returns the number of tracked leases.
func (t *Tracker) Len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.entries)
}
