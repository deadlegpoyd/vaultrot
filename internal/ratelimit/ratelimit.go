// Package ratelimit provides a simple token-bucket rate limiter for
// controlling how frequently secret rotations are performed against backends.
package ratelimit

import (
	"fmt"
	"sync"
	"time"
)

// Limiter controls the rate of operations using a token-bucket strategy.
type Limiter struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	rate     float64 // tokens per second
	lastTick time.Time
	clock    func() time.Time
}

// Config holds configuration for constructing a Limiter.
type Config struct {
	// MaxTokens is the bucket capacity (burst size).
	MaxTokens float64
	// RatePerSecond is how many tokens are replenished per second.
	RatePerSecond float64
}

// New creates a Limiter from the given Config.
// Returns an error if MaxTokens or RatePerSecond are not positive.
func New(cfg Config) (*Limiter, error) {
	if cfg.MaxTokens <= 0 {
		return nil, fmt.Errorf("ratelimit: MaxTokens must be positive, got %v", cfg.MaxTokens)
	}
	if cfg.RatePerSecond <= 0 {
		return nil, fmt.Errorf("ratelimit: RatePerSecond must be positive, got %v", cfg.RatePerSecond)
	}
	now := time.Now()
	return &Limiter{
		tokens:   cfg.MaxTokens,
		max:      cfg.MaxTokens,
		rate:     cfg.RatePerSecond,
		lastTick: now,
		clock:    time.Now,
	}, nil
}

// Allow returns true and consumes one token if the operation is permitted.
// Returns false without consuming a token when the bucket is empty.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.clock()
	elapsed := now.Sub(l.lastTick).Seconds()
	l.lastTick = now

	l.tokens += elapsed * l.rate
	if l.tokens > l.max {
		l.tokens = l.max
	}

	if l.tokens < 1 {
		return false
	}
	l.tokens--
	return true
}

// Wait blocks until a token is available, then consumes it.
func (l *Limiter) Wait() {
	for {
		if l.Allow() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// Available returns the current number of available tokens (approximate).
func (l *Limiter) Available() float64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.tokens
}
