package ratelimit_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultrot/internal/ratelimit"
)

func TestNew_ValidConfig(t *testing.T) {
	l, err := ratelimit.New(ratelimit.Config{MaxTokens: 5, RatePerSecond: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil Limiter")
	}
}

func TestNew_InvalidMaxTokens(t *testing.T) {
	_, err := ratelimit.New(ratelimit.Config{MaxTokens: 0, RatePerSecond: 1})
	if err == nil {
		t.Fatal("expected error for zero MaxTokens")
	}
}

func TestNew_InvalidRate(t *testing.T) {
	_, err := ratelimit.New(ratelimit.Config{MaxTokens: 5, RatePerSecond: -1})
	if err == nil {
		t.Fatal("expected error for negative RatePerSecond")
	}
}

func TestAllow_ConsumesTokens(t *testing.T) {
	l, _ := ratelimit.New(ratelimit.Config{MaxTokens: 3, RatePerSecond: 0.01})

	for i := 0; i < 3; i++ {
		if !l.Allow() {
			t.Fatalf("expected Allow()=true on call %d", i+1)
		}
	}
	// Bucket should now be empty
	if l.Allow() {
		t.Fatal("expected Allow()=false when bucket empty")
	}
}

func TestAllow_RefillsOverTime(t *testing.T) {
	// High rate so tokens replenish quickly in the test
	l, _ := ratelimit.New(ratelimit.Config{MaxTokens: 2, RatePerSecond: 1000})

	// Drain the bucket
	l.Allow()
	l.Allow()

	// Wait briefly for tokens to replenish
	time.Sleep(5 * time.Millisecond)

	if !l.Allow() {
		t.Fatal("expected token to be replenished after sleep")
	}
}

func TestAvailable_StartsAtMax(t *testing.T) {
	l, _ := ratelimit.New(ratelimit.Config{MaxTokens: 10, RatePerSecond: 1})
	if got := l.Available(); got != 10 {
		t.Fatalf("expected 10 available tokens, got %v", got)
	}
}

func TestAvailable_DecreasesAfterAllow(t *testing.T) {
	l, _ := ratelimit.New(ratelimit.Config{MaxTokens: 5, RatePerSecond: 0.01})
	l.Allow()
	l.Allow()
	if got := l.Available(); got > 3 {
		t.Fatalf("expected <=3 tokens after 2 Allow() calls, got %v", got)
	}
}

func TestWait_ReturnsWhenTokenAvailable(t *testing.T) {
	l, _ := ratelimit.New(ratelimit.Config{MaxTokens: 1, RatePerSecond: 200})
	l.Allow() // drain

	done := make(chan struct{})
	go func() {
		l.Wait()
		close(done)
	}()

	select {
	case <-done:
		// success
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Wait() did not return within 500ms")
	}
}
