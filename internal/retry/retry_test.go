package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/yourusername/vaultrot/internal/retry"
)

func fastConfig(attempts int) retry.Config {
	return retry.Config{
		MaxAttempts:     attempts,
		InitialInterval: time.Millisecond,
		MaxInterval:     10 * time.Millisecond,
		Multiplier:      2.0,
	}
}

func TestNew_ValidConfig(t *testing.T) {
	_, err := retry.New(retry.DefaultConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_InvalidMaxAttempts(t *testing.T) {
	cfg := retry.DefaultConfig()
	cfg.MaxAttempts = 0
	_, err := retry.New(cfg)
	if err == nil {
		t.Fatal("expected error for MaxAttempts=0")
	}
}

func TestNew_InvalidMultiplier(t *testing.T) {
	cfg := retry.DefaultConfig()
	cfg.Multiplier = 0.5
	_, err := retry.New(cfg)
	if err == nil {
		t.Fatal("expected error for Multiplier < 1")
	}
}

func TestDo_SucceedsOnFirstAttempt(t *testing.T) {
	r, _ := retry.New(fastConfig(3))
	calls := 0
	err := r.Do(context.Background(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesOnFailure(t *testing.T) {
	r, _ := retry.New(fastConfig(3))
	sentinel := errors.New("transient")
	calls := 0
	err := r.Do(context.Background(), func() error {
		calls++
		if calls < 3 {
			return sentinel
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ExhaustsAttempts(t *testing.T) {
	r, _ := retry.New(fastConfig(3))
	sentinel := errors.New("permanent")
	calls := 0
	err := r.Do(context.Background(), func() error {
		calls++
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_RespectsContextCancellation(t *testing.T) {
	r, _ := retry.New(fastConfig(5))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := r.Do(ctx, func() error { return nil })
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}
