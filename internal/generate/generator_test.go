package generate_test

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/yourusername/vaultrot/internal/generate"
)

func TestSecret_DefaultLength(t *testing.T) {
	g := generate.New(generate.Options{})
	s, err := g.Secret()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s) != generate.DefaultLength {
		t.Errorf("expected length %d, got %d", generate.DefaultLength, len(s))
	}
}

func TestSecret_CustomLength(t *testing.T) {
	g := generate.New(generate.Options{Length: 64})
	s, err := g.Secret()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s) != 64 {
		t.Errorf("expected length 64, got %d", len(s))
	}
}

func TestSecret_AlphaOnly(t *testing.T) {
	g := generate.New(generate.Options{Length: 128, AlphaOnly: true})
	s, err := g.Secret()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, c := range s {
		if strings.ContainsRune("!@#$%^&*()-_=+[]{}", c) {
			t.Errorf("AlphaOnly secret contains special char: %q", c)
		}
	}
}

func TestSecret_Base64Encoded(t *testing.T) {
	g := generate.New(generate.Options{Length: 32, Base64Encode: true})
	s, err := g.Secret()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := base64.StdEncoding.DecodeString(s); err != nil {
		t.Errorf("expected valid base64, got error: %v", err)
	}
}

func TestSecret_Uniqueness(t *testing.T) {
	g := generate.New(generate.Options{Length: 32})
	a, _ := g.Secret()
	b, _ := g.Secret()
	if a == b {
		t.Error("two generated secrets should not be equal (collision extremely unlikely)")
	}
}

func TestBytes_ReturnsBase64(t *testing.T) {
	s, err := generate.Bytes(32)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		t.Fatalf("expected valid base64: %v", err)
	}
	if len(decoded) != 32 {
		t.Errorf("expected 32 raw bytes, got %d", len(decoded))
	}
}

func TestBytes_DefaultsWhenZero(t *testing.T) {
	s, err := generate.Bytes(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == "" {
		t.Error("expected non-empty result for zero length input")
	}
}
