package backend

import (
	"context"
	"fmt"

	"github.com/vaultrot/internal/config"
)

// Secret represents a secret value with metadata.
type Secret struct {
	Key   string
	Value string
}

// Backend defines the interface for all secret backends.
type Backend interface {
	// GetSecret retrieves a secret by key.
	GetSecret(ctx context.Context, key string) (*Secret, error)
	// SetSecret writes a secret value for the given key.
	SetSecret(ctx context.Context, key, value string) error
	// Name returns the backend identifier.
	Name() string
}

// New constructs the appropriate Backend implementation based on the
// secret config's backend field.
func New(cfg config.SecretConfig) (Backend, error) {
	switch cfg.Backend {
	case "vault":
		return NewVaultBackend(cfg)
	case "aws_ssm":
		return NewSSMBackend(cfg)
	case "doppler":
		return NewDopplerBackend(cfg)
	default:
		return nil, fmt.Errorf("unsupported backend: %q", cfg.Backend)
	}
}

// --- stub implementations kept in this file for now ---

// VaultBackend is a stub for HashiCorp Vault.
type VaultBackend struct{ cfg config.SecretConfig }

func NewVaultBackend(cfg config.SecretConfig) (*VaultBackend, error) {
	return &VaultBackend{cfg: cfg}, nil
}
func (b *VaultBackend) Name() string { return "vault" }
func (b *VaultBackend) GetSecret(ctx context.Context, key string) (*Secret, error) {
	return nil, fmt.Errorf("vault: GetSecret not yet implemented")
}
func (b *VaultBackend) SetSecret(ctx context.Context, key, value string) error {
	return fmt.Errorf("vault: SetSecret not yet implemented")
}

// SSMBackend is a stub for AWS SSM Parameter Store.
type SSMBackend struct{ cfg config.SecretConfig }

func NewSSMBackend(cfg config.SecretConfig) (*SSMBackend, error) {
	return &SSMBackend{cfg: cfg}, nil
}
func (b *SSMBackend) Name() string { return "aws_ssm" }
func (b *SSMBackend) GetSecret(ctx context.Context, key string) (*Secret, error) {
	return nil, fmt.Errorf("aws_ssm: GetSecret not yet implemented")
}
func (b *SSMBackend) SetSecret(ctx context.Context, key, value string) error {
	return fmt.Errorf("aws_ssm: SetSecret not yet implemented")
}

// DopplerBackend is a stub for Doppler.
type DopplerBackend struct{ cfg config.SecretConfig }

func NewDopplerBackend(cfg config.SecretConfig) (*DopplerBackend, error) {
	return &DopplerBackend{cfg: cfg}, nil
}
func (b *DopplerBackend) Name() string { return "doppler" }
func (b *DopplerBackend) GetSecret(ctx context.Context, key string) (*Secret, error) {
	return nil, fmt.Errorf("doppler: GetSecret not yet implemented")
}
func (b *DopplerBackend) SetSecret(ctx context.Context, key, value string) error {
	return fmt.Errorf("doppler: SetSecret not yet implemented")
}
