// Package observe provides a lightweight event observer that fans out
// rotation lifecycle events to registered listeners.
package observe

import (
	"fmt"
	"sync"
	"time"
)

// EventKind classifies a lifecycle event.
type EventKind string

const (
	EventRotationStarted  EventKind = "rotation.started"
	EventRotationFinished EventKind = "rotation.finished"
	EventSecretRotated    EventKind = "secret.rotated"
	EventSecretSkipped    EventKind = "secret.skipped"
	EventSecretFailed     EventKind = "secret.failed"
)

// Event carries metadata about a single lifecycle moment.
type Event struct {
	Kind      EventKind
	Secret    string
	Backend   string
	DryRun    bool
	Err       error
	OccurredAt time.Time
}

// String returns a human-readable representation of the event.
func (e Event) String() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s/%s error=%v", e.Kind, e.Backend, e.Secret, e.Err)
	}
	return fmt.Sprintf("[%s] %s/%s", e.Kind, e.Backend, e.Secret)
}

// Handler is a function that receives an Event.
type Handler func(Event)

// Observer fans out events to all registered handlers.
type Observer struct {
	mu       sync.RWMutex
	handlers []Handler
	clock    func() time.Time
}

// New returns an Observer ready for use.
func New() *Observer {
	return &Observer{clock: time.Now}
}

// Register adds a Handler to the observer.
// Handlers are called in registration order.
func (o *Observer) Register(h Handler) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.handlers = append(o.handlers, h)
}

// Emit dispatches e to every registered handler.
// OccurredAt is set automatically if it is zero.
func (o *Observer) Emit(e Event) {
	if e.OccurredAt.IsZero() {
		e.OccurredAt = o.clock()
	}
	o.mu.RLock()
	handlers := make([]Handler, len(o.handlers))
	copy(handlers, o.handlers)
	o.mu.RUnlock()

	for _, h := range handlers {
		h(e)
	}
}

// Len returns the number of registered handlers.
func (o *Observer) Len() int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return len(o.handlers)
}
