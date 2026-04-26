package backend_test

import (
	"testing"

	"github.com/vaultrot/internal/backend"
	"github.com/vaultrot/internal/config"
)

func secretCfg(b string) config.SecretConfig {
	return config.SecretConfig{
		Name:    "test-secret",
		Backend: b,
		Key:     "/test/key",
	}
}

func TestNew_Vault(t *testing.T) {
	b, err := backend.New(secretCfg("vault"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if b.Name() != "vault" {
		t.Errorf("expected name %q, got %q", "vault", b.Name())
	}
}

func TestNew_AWSSSM(t *testing.T) {
	b, err := backend.New(secretCfg("aws_ssm"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if b.Name() != "aws_ssm" {
		t.Errorf("expected name %q, got %q", "aws_ssm", b.Name())
	}
}

func TestNew_Doppler(t *testing.T) {
	b, err := backend.New(secretCfg("doppler"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if b.Name() != "doppler" {
		t.Errorf("expected name %q, got %q", "doppler", b.Name())
	}
}

func TestNew_UnsupportedBackend(t *testing.T) {
	_, err := backend.New(secretCfg("unknown"))
	if err == nil {
		t.Fatal("expected error for unsupported backend, got nil")
	}
}

func TestVaultBackend_StubErrors(t *testing.T) {
	b, _ := backend.New(secretCfg("vault"))
	ctx := t.Context()

	if _, err := b.GetSecret(ctx, "key"); err == nil {
		t.Error("expected GetSecret to return stub error")
	}
	if err := b.SetSecret(ctx, "key", "value"); err == nil {
		t.Error("expected SetSecret to return stub error")
	}
}

func TestSSMBackend_StubErrors(t *testing.T) {
	b, _ := backend.New(secretCfg("aws_ssm"))
	ctx := t.Context()

	if _, err := b.GetSecret(ctx, "key"); err == nil {
		t.Error("expected GetSecret to return stub error")
	}
	if err := b.SetSecret(ctx, "key", "value"); err == nil {
		t.Error("expected SetSecret to return stub error")
	}
}

func TestDopplerBackend_StubErrors(t *testing.T) {
	b, _ := backend.New(secretCfg("doppler"))
	ctx := t.Context()

	if _, err := b.GetSecret(ctx, "key"); err == nil {
		t.Error("expected GetSecret to return stub error")
	}
	if err := b.SetSecret(ctx, "key", "value"); err == nil {
		t.Error("expected SetSecret to return stub error")
	}
}
