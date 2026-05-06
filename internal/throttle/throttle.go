// Package throttle provides per-backend concurrency limiting for secret rotation.
// It ensures that no more than a configured number of rotation operations run
// simultaneously against a single backend, preventing API quota exhaustion.
package throttle

import (
	"context"
	"fmt"
	"sync"
)

// Throttle controls concurrent access to a named backend.
type Throttle struct {
	mu      sync.Mutex
	sems    map[string]chan struct{}
	limit   int
}

// Config holds configuration for the throttle.
type Config struct {
	// MaxConcurrent is the maximum number of simultaneous operations per backend.
	MaxConcurrent int
}

// New creates a Throttle with the given configuration.
// Returns an error if MaxConcurrent is less than 1.
func New(cfg Config) (*Throttle, error) {
	if cfg.MaxConcurrent < 1 {
		return nil, fmt.Errorf("throttle: MaxConcurrent must be at least 1, got %d", cfg.MaxConcurrent)
	}
	return &Throttle{
		sems:  make(map[string]chan struct{}),
		limit: cfg.MaxConcurrent,
	}, nil
}

// Acquire blocks until a slot is available for the named backend or the context
// is cancelled. Returns an error if the context expires before a slot is acquired.
func (t *Throttle) Acquire(ctx context.Context, backend string) error {
	sem := t.semaphore(backend)
	select {
	case sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("throttle: acquire cancelled for backend %q: %w", backend, ctx.Err())
	}
}

// Release frees a slot for the named backend. It is a no-op if the backend has
// no active acquisitions.
func (t *Throttle) Release(backend string) {
	sem := t.semaphore(backend)
	select {
	case <-sem:
	default:
	}
}

// Active returns the number of currently held slots for the named backend.
func (t *Throttle) Active(backend string) int {
	t.mu.Lock()
	sem, ok := t.sems[backend]
	t.mu.Unlock()
	if !ok {
		return 0
	}
	return len(sem)
}

func (t *Throttle) semaphore(backend string) chan struct{} {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.sems[backend]; !ok {
		t.sems[backend] = make(chan struct{}, t.limit)
	}
	return t.sems[backend]
}
