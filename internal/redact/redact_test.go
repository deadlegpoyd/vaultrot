package redact_test

import (
	"strings"
	"testing"

	"github.com/yourusername/vaultrot/internal/redact"
)

func TestMask_DefaultFullyRedacts(t *testing.T) {
	r := redact.New()
	got := r.Mask("supersecret")
	if got != "***REDACTED***" {
		t.Fatalf("expected default mask, got %q", got)
	}
}

func TestMask_EmptyValuePassesThrough(t *testing.T) {
	r := redact.New()
	if got := r.Mask(""); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestMask_CustomMask(t *testing.T) {
	r := redact.New(redact.WithMask("[hidden]"))
	if got := r.Mask("abc"); got != "[hidden]" {
		t.Fatalf("unexpected mask %q", got)
	}
}

func TestMask_PeekPrefixAndSuffix(t *testing.T) {
	r := redact.New(redact.WithMask("***"), redact.WithPeek(2, 2))
	got := r.Mask("supersecret")
	if !strings.HasPrefix(got, "su") {
		t.Fatalf("expected prefix 'su', got %q", got)
	}
	if !strings.HasSuffix(got, "et") {
		t.Fatalf("expected suffix 'et', got %q", got)
	}
	if !strings.Contains(got, "***") {
		t.Fatalf("expected mask in middle, got %q", got)
	}
}

func TestMask_PeekExceedsLength_FullyRedacts(t *testing.T) {
	r := redact.New(redact.WithPeek(5, 5))
	got := r.Mask("short")
	if got != "***REDACTED***" {
		t.Fatalf("expected full redaction for short value, got %q", got)
	}
}

func TestMaskMap_RedactsAllValues(t *testing.T) {
	r := redact.New()
	input := map[string]string{
		"DB_PASS": "hunter2",
		"API_KEY": "abc123xyz",
	}
	out := r.MaskMap(input)
	for k, v := range out {
		if v != "***REDACTED***" {
			t.Errorf("key %q: expected redacted value, got %q", k, v)
		}
	}
	// original map must not be mutated
	if input["DB_PASS"] != "hunter2" {
		t.Fatal("original map was mutated")
	}
}

func TestMaskMap_EmptyMap(t *testing.T) {
	r := redact.New()
	out := r.MaskMap(map[string]string{})
	if len(out) != 0 {
		t.Fatalf("expected empty map, got %v", out)
	}
}
