// Package ownership tracks which team or service owns a given secret,
// enabling targeted notifications and access-control checks during rotation.
package ownership

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

// Owner describes the team or service responsible for a secret.
type Owner struct {
	// Name is a human-readable label, e.g. "platform-team".
	Name string `json:"name"`
	// Contact is an email address or Slack channel for escalations.
	Contact string `json:"contact"`
	// Tags are arbitrary key/value metadata attached to the owner.
	Tags map[string]string `json:"tags,omitempty"`
}

// Registry maps secret keys to their owners.
type Registry struct {
	mu      sync.RWMutex
	entries map[string]Owner
}

// New returns an empty Registry.
func New() *Registry {
	return &Registry{entries: make(map[string]Owner)}
}

// Register associates secretKey with owner.
// Returns an error when secretKey or owner.Name is empty.
func (r *Registry) Register(secretKey string, owner Owner) error {
	if strings.TrimSpace(secretKey) == "" {
		return errors.New("ownership: secret key must not be empty")
	}
	if strings.TrimSpace(owner.Name) == "" {
		return errors.New("ownership: owner name must not be empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[secretKey] = owner
	return nil
}

// Lookup returns the Owner for secretKey.
// The boolean is false when no registration exists.
func (r *Registry) Lookup(secretKey string) (Owner, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	owner, ok := r.entries[secretKey]
	return owner, ok
}

// MustLookup returns the Owner for secretKey or panics.
func (r *Registry) MustLookup(secretKey string) Owner {
	owner, ok := r.Lookup(secretKey)
	if !ok {
		panic(fmt.Sprintf("ownership: no owner registered for %q", secretKey))
	}
	return owner
}

// Len returns the number of registered entries.
func (r *Registry) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.entries)
}

// All returns a shallow copy of all registered entries.
func (r *Registry) All() map[string]Owner {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]Owner, len(r.entries))
	for k, v := range r.entries {
		out[k] = v
	}
	return out
}
