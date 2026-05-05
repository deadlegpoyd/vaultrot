// Package retry provides a simple exponential-backoff retry mechanism
// used when communicating with secret backends.
package retry

import (
	"context"
	"errors"
	"math"
	"time"
)

// Config holds tuning parameters for the retry loop.
type Config struct {
	// MaxAttempts is the total number of tries (including the first).
	MaxAttempts int
	// InitialInterval is the wait time before the second attempt.
	InitialInterval time.Duration
	// MaxInterval caps the computed back-off duration.
	MaxInterval time.Duration
	// Multiplier is applied to the interval after each failure.
	Multiplier float64
}

// DefaultConfig returns sensible defaults suitable for backend API calls.
func DefaultConfig() Config {
	return Config{
		MaxAttempts:     4,
		InitialInterval: 250 * time.Millisecond,
		MaxInterval:     10 * time.Second,
		Multiplier:      2.0,
	}
}

// Retrier executes operations with exponential back-off.
type Retrier struct {
	cfg   Config
	sleep func(time.Duration) // injectable for tests
}

// New creates a Retrier from cfg. Pass retry.DefaultConfig() for production use.
func New(cfg Config) (*Retrier, error) {
	if cfg.MaxAttempts < 1 {
		return nil, errors.New("retry: MaxAttempts must be >= 1")
	}
	if cfg.Multiplier < 1.0 {
		return nil, errors.New("retry: Multiplier must be >= 1.0")
	}
	return &Retrier{
		cfg:   cfg,
		sleep: time.Sleep,
	}, nil
}

// Do calls fn up to MaxAttempts times. It stops early if ctx is cancelled or
// fn returns nil. The last non-nil error is returned.
func (r *Retrier) Do(ctx context.Context, fn func() error) error {
	var err error
	interval := r.cfg.InitialInterval

	for attempt := 1; attempt <= r.cfg.MaxAttempts; attempt++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		err = fn()
		if err == nil {
			return nil
		}
		if attempt == r.cfg.MaxAttempts {
			break
		}
		wait := time.Duration(math.Min(float64(interval), float64(r.cfg.MaxInterval)))
		r.sleep(wait)
		interval = time.Duration(float64(interval) * r.cfg.Multiplier)
	}
	return err
}
