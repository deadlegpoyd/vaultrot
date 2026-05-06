// Package policy enforces rotation policies such as minimum age and
// maximum age constraints before allowing a secret to be rotated.
package policy

import (
	"errors"
	"fmt"
	"time"
)

// ErrPolicyViolation is returned when a secret does not satisfy a policy.
var ErrPolicyViolation = errors.New("policy violation")

// Config holds rotation policy constraints for a single secret.
type Config struct {
	// MinAgeDays prevents rotation if the secret is younger than this many days.
	MinAgeDays int `yaml:"min_age_days"`
	// MaxAgeDays requires rotation if the secret is older than this many days.
	MaxAgeDays int `yaml:"max_age_days"`
}

// Policy evaluates whether a secret is eligible for rotation.
type Policy struct {
	cfg  Config
	now  func() time.Time
}

// New returns a Policy using the provided Config.
// An optional clock function can be supplied for testing; pass nil to use time.Now.
func New(cfg Config, now func() time.Time) (*Policy, error) {
	if cfg.MinAgeDays < 0 {
		return nil, fmt.Errorf("%w: min_age_days must be >= 0", ErrPolicyViolation)
	}
	if cfg.MaxAgeDays < 0 {
		return nil, fmt.Errorf("%w: max_age_days must be >= 0", ErrPolicyViolation)
	}
	if cfg.MaxAgeDays > 0 && cfg.MinAgeDays > cfg.MaxAgeDays {
		return nil, fmt.Errorf("%w: min_age_days (%d) cannot exceed max_age_days (%d)",
			ErrPolicyViolation, cfg.MinAgeDays, cfg.MaxAgeDays)
	}
	if now == nil {
		now = time.Now
	}
	return &Policy{cfg: cfg, now: now}, nil
}

// Check returns an error if lastRotated violates the policy.
// A zero lastRotated is treated as "never rotated" and satisfies any MinAge check.
func (p *Policy) Check(name string, lastRotated time.Time) error {
	age := p.age(lastRotated)

	if p.cfg.MinAgeDays > 0 && !lastRotated.IsZero() {
		min := time.Duration(p.cfg.MinAgeDays) * 24 * time.Hour
		if age < min {
			return fmt.Errorf("%w: secret %q is only %s old (min %dd)",
				ErrPolicyViolation, name, age.Round(time.Minute), p.cfg.MinAgeDays)
		}
	}

	if p.cfg.MaxAgeDays > 0 {
		max := time.Duration(p.cfg.MaxAgeDays) * 24 * time.Hour
		if age > max {
			return fmt.Errorf("%w: secret %q is %s old (max %dd)",
				ErrPolicyViolation, name, age.Round(time.Minute), p.cfg.MaxAgeDays)
		}
	}

	return nil
}

// RequiresRotation returns true when the secret has exceeded MaxAgeDays.
func (p *Policy) RequiresRotation(lastRotated time.Time) bool {
	if p.cfg.MaxAgeDays <= 0 {
		return false
	}
	max := time.Duration(p.cfg.MaxAgeDays) * 24 * time.Hour
	return p.age(lastRotated) > max
}

func (p *Policy) age(lastRotated time.Time) time.Duration {
	if lastRotated.IsZero() {
		return 0
	}
	return p.now().UTC().Sub(lastRotated.UTC())
}
