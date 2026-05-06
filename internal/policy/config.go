package policy

import (
	"fmt"

	"github.com/yourusername/vaultrot/internal/config"
)

// FromSecretConfig builds a Policy from the rotation policy block of a secret
// config entry. If no policy block is present, a zero-constraint Policy is
// returned that allows rotation at any time.
func FromSecretConfig(s config.Secret) (*Policy, error) {
	cfg := Config{}
	if s.Policy != nil {
		cfg.MinAgeDays = s.Policy.MinAgeDays
		cfg.MaxAgeDays = s.Policy.MaxAgeDays
	}
	p, err := New(cfg, nil)
	if err != nil {
		return nil, fmt.Errorf("secret %q: %w", s.Name, err)
	}
	return p, nil
}

// MustFromSecretConfig is like FromSecretConfig but panics on error.
// Intended for use in tests or init paths where misconfiguration is fatal.
func MustFromSecretConfig(s config.Secret) *Policy {
	p, err := FromSecretConfig(s)
	if err != nil {
		panic(err)
	}
	return p
}
