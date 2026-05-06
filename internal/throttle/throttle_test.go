package throttle_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/yourusername/vaultrot/internal/throttle"
)

func TestNew_ValidConfig(t *testing.T) {
	th, err := throttle.New(throttle.Config{MaxConcurrent: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th == nil {
		t.Fatal("expected non-nil throttle")
	}
}

func TestNew_InvalidMaxConcurrent(t *testing.T) {
	_, err := throttle.New(throttle.Config{MaxConcurrent: 0})
	if err == nil {
		t.Fatal("expected error for MaxConcurrent=0")
	}
}

func TestAcquire_And_Release(t *testing.T) {
	th, _ := throttle.New(throttle.Config{MaxConcurrent: 2})
	ctx := context.Background()

	if err := th.Acquire(ctx, "vault"); err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}
	if th.Active("vault") != 1 {
		t.Fatalf("expected active=1, got %d", th.Active("vault"))
	}

	th.Release("vault")
	if th.Active("vault") != 0 {
		t.Fatalf("expected active=0 after release, got %d", th.Active("vault"))
	}
}

func TestAcquire_BlocksAtLimit(t *testing.T) {
	th, _ := throttle.New(throttle.Config{MaxConcurrent: 1})
	ctx := context.Background()

	if err := th.Acquire(ctx, "ssm"); err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := th.Acquire(ctxTimeout, "ssm")
	if err == nil {
		t.Fatal("expected error when limit reached and context expires")
	}
}

func TestAcquire_IsolatedPerBackend(t *testing.T) {
	th, _ := throttle.New(throttle.Config{MaxConcurrent: 1})
	ctx := context.Background()

	if err := th.Acquire(ctx, "vault"); err != nil {
		t.Fatalf("vault acquire failed: %v", err)
	}
	if err := th.Acquire(ctx, "doppler"); err != nil {
		t.Fatalf("doppler acquire should succeed independently: %v", err)
	}
	if th.Active("vault") != 1 {
		t.Errorf("expected vault active=1")
	}
	if th.Active("doppler") != 1 {
		t.Errorf("expected doppler active=1")
	}
}

func TestActive_UnknownBackend_ReturnsZero(t *testing.T) {
	th, _ := throttle.New(throttle.Config{MaxConcurrent: 2})
	if got := th.Active("unknown"); got != 0 {
		t.Errorf("expected 0 for unknown backend, got %d", got)
	}
}

func TestAcquire_ConcurrentGoroutines(t *testing.T) {
	const limit = 3
	th, _ := throttle.New(throttle.Config{MaxConcurrent: limit})
	ctx := context.Background()

	var wg sync.WaitGroup
	for i := 0; i < limit; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := th.Acquire(ctx, "vault"); err != nil {
				t.Errorf("concurrent acquire failed: %v", err)
			}
		}()
	}
	wg.Wait()
	if th.Active("vault") != limit {
		t.Errorf("expected active=%d, got %d", limit, th.Active("vault"))
	}
}
