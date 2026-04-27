package rotation_test

import (
	"testing"

	"github.com/yourusername/vaultrot/internal/config"
	"github.com/yourusername/vaultrot/internal/rotation"
)

func baseConfig(secrets []config.Secret) *config.Config {
	return &config.Config{Secrets: secrets}
}

func TestRun_DryRun_NoErrors(t *testing.T) {
	cfg := baseConfig([]config.Secret{
		{Name: "db-password", Backend: "vault", Path: "secret/db"},
		{Name: "api-key", Backend: "aws-ssm", Path: "/prod/api-key"},
	})

	rot := rotation.New(cfg, true)
	results := rot.Run()

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	for _, r := range results {
		if !r.DryRun {
			t.Errorf("expected DryRun=true for secret %s", r.SecretName)
		}
		if !r.Success {
			t.Errorf("expected Success=true for secret %s in dry-run", r.SecretName)
		}
		if r.Err != nil {
			t.Errorf("unexpected error for secret %s: %v", r.SecretName, r.Err)
		}
	}
}

func TestRun_DryRun_EmptySecrets(t *testing.T) {
	cfg := baseConfig([]config.Secret{})
	rot := rotation.New(cfg, true)
	results := rot.Run()

	if len(results) != 0 {
		t.Fatalf("expected 0 results for empty config, got %d", len(results))
	}
}

func TestRun_DryRun_ResultFields(t *testing.T) {
	secret := config.Secret{Name: "my-secret", Backend: "doppler", Path: "prod/MY_SECRET"}
	cfg := baseConfig([]config.Secret{secret})

	rot := rotation.New(cfg, true)
	results := rot.Run()

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	r := results[0]
	if r.SecretName != secret.Name {
		t.Errorf("expected SecretName=%s, got %s", secret.Name, r.SecretName)
	}
	if r.Backend != secret.Backend {
		t.Errorf("expected Backend=%s, got %s", secret.Backend, r.Backend)
	}
}
