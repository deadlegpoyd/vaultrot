package masking_test

import (
	"testing"

	"github.com/yourusername/vaultrot/internal/masking"
)

func rules() []masking.Rule {
	return []masking.Rule{
		{Field: "*password*"},
		{Field: "*token*", Mask: "[redacted]"},
		{Field: "api_key"},
	}
}

func TestApply_NoMatch_ReturnsOriginal(t *testing.T) {
	m, err := masking.New(rules(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := m.Apply("username", "alice")
	if got != "alice" {
		t.Errorf("expected 'alice', got %q", got)
	}
}

func TestApply_WildcardMatch_UsesDefaultMask(t *testing.T) {
	m, _ := masking.New(rules(), "***")
	got := m.Apply("db_password", "s3cr3t")
	if got != "***" {
		t.Errorf("expected '***', got %q", got)
	}
}

func TestApply_WildcardMatch_UsesCustomMask(t *testing.T) {
	m, _ := masking.New(rules(), "***")
	got := m.Apply("access_token", "tok_abc123")
	if got != "[redacted]" {
		t.Errorf("expected '[redacted]', got %q", got)
	}
}

func TestApply_ExactMatch_CaseInsensitive(t *testing.T) {
	m, _ := masking.New(rules(), "***")
	got := m.Apply("API_KEY", "key-value")
	if got != "***" {
		t.Errorf("expected '***', got %q", got)
	}
}

func TestApply_EmptyValue_MaskedWhenMatched(t *testing.T) {
	m, _ := masking.New(rules(), "***")
	got := m.Apply("password", "")
	if got != "***" {
		t.Errorf("expected '***', got %q", got)
	}
}

func TestApplyMap_MasksMatchingKeys(t *testing.T) {
	m, _ := masking.New(rules(), "***")
	input := map[string]string{
		"username":     "alice",
		"password":     "s3cr3t",
		"access_token": "tok_abc",
	}
	out := m.ApplyMap(input)
	if out["username"] != "alice" {
		t.Errorf("username should be unchanged")
	}
	if out["password"] != "***" {
		t.Errorf("password should be masked")
	}
	if out["access_token"] != "[redacted]" {
		t.Errorf("access_token should use custom mask")
	}
}

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := masking.New([]masking.Rule{{Field: "["}}, "")
	if err == nil {
		t.Error("expected error for invalid pattern, got nil")
	}
}

func TestNew_NoRules_IsNoop(t *testing.T) {
	m, err := masking.New(nil, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := m.Apply("password", "secret")
	if got != "secret" {
		t.Errorf("expected no masking with empty rules, got %q", got)
	}
}
