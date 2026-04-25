package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultrot/internal/config"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "vaultrot.yaml")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}
	return p
}

func TestLoad_ValidConfig(t *testing.T) {
	content := `
dry_run: true
secrets:
  - name: db-password
    path: /prod/db/password
    backend: aws_ssm
vault:
  address: http://127.0.0.1:8200
  token: root
`
	p := writeTempConfig(t, content)
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.DryRun {
		t.Error("expected dry_run to be true")
	}
	if len(cfg.Secrets) != 1 {
		t.Fatalf("expected 1 secret, got %d", len(cfg.Secrets))
	}
	if cfg.Secrets[0].Backend != config.BackendAWSSSM {
		t.Errorf("expected backend aws_ssm, got %q", cfg.Secrets[0].Backend)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path/vaultrot.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestValidate_MissingName(t *testing.T) {
	cfg := &config.Config{
		Secrets: []config.SecretEntry{
			{Path: "/some/path", Backend: config.BackendVault},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for missing name")
	}
}

func TestValidate_UnsupportedBackend(t *testing.T) {
	cfg := &config.Config{
		Secrets: []config.SecretEntry{
			{Name: "test", Path: "/test", Backend: "unknown"},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for unsupported backend")
	}
}

func TestValidate_Empty(t *testing.T) {
	cfg := &config.Config{}
	if err := cfg.Validate(); err != nil {
		t.Errorf("empty config should be valid, got: %v", err)
	}
}
