package observe_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/yourusername/vaultrot/internal/observe"
)

func TestEmit_CallsAllHandlers(t *testing.T) {
	obs := observe.New()
	var got []observe.EventKind
	var mu sync.Mutex

	for i := 0; i < 3; i++ {
		obs.Register(func(e observe.Event) {
			mu.Lock()
			got = append(got, e.Kind)
			mu.Unlock()
		})
	}

	obs.Emit(observe.Event{Kind: observe.EventSecretRotated, Secret: "x", Backend: "vault"})

	if len(got) != 3 {
		t.Fatalf("expected 3 handler calls, got %d", len(got))
	}
	for _, k := range got {
		if k != observe.EventSecretRotated {
			t.Errorf("unexpected kind %q", k)
		}
	}
}

func TestEmit_SetsOccurredAt_WhenZero(t *testing.T) {
	obs := observe.New()
	var received observe.Event
	obs.Register(func(e observe.Event) { received = e })

	obs.Emit(observe.Event{Kind: observe.EventRotationStarted})

	if received.OccurredAt.IsZero() {
		t.Error("expected OccurredAt to be set automatically")
	}
}

func TestEmit_PreservesOccurredAt_WhenProvided(t *testing.T) {
	obs := observe.New()
	fixed := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	var received observe.Event
	obs.Register(func(e observe.Event) { received = e })

	obs.Emit(observe.Event{Kind: observe.EventSecretSkipped, OccurredAt: fixed})

	if !received.OccurredAt.Equal(fixed) {
		t.Errorf("expected %v, got %v", fixed, received.OccurredAt)
	}
}

func TestLen_ReturnsHandlerCount(t *testing.T) {
	obs := observe.New()
	if obs.Len() != 0 {
		t.Fatalf("expected 0, got %d", obs.Len())
	}
	obs.Register(func(observe.Event) {})
	obs.Register(func(observe.Event) {})
	if obs.Len() != 2 {
		t.Fatalf("expected 2, got %d", obs.Len())
	}
}

func TestEvent_String_WithError(t *testing.T) {
	e := observe.Event{
		Kind:    observe.EventSecretFailed,
		Secret:  "api/key",
		Backend: "doppler",
		Err:     errors.New("timeout"),
	}
	s := e.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
	for _, want := range []string{"secret.failed", "doppler", "api/key", "timeout"} {
		if !contains(s, want) {
			t.Errorf("expected %q in %q", want, s)
		}
	}
}

func TestEvent_String_NoError(t *testing.T) {
	e := observe.Event{Kind: observe.EventSecretRotated, Secret: "db/pass", Backend: "vault"}
	s := e.String()
	if contains(s, "error") {
		t.Errorf("unexpected 'error' in %q", s)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
