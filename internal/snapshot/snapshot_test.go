package snapshot_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultrot/internal/snapshot"
)

var fixedTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func newStore(t *testing.T) *snapshot.Store {
	t.Helper()
	return snapshot.New()
}

func TestCapture_And_Get_ReturnsEntry(t *testing.T) {
	s := newStore(t)
	s.Capture("vault", "db/password", "secret123")

	e, ok := s.Get("vault", "db/password")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Value != "secret123" {
		t.Errorf("got value %q, want %q", e.Value, "secret123")
	}
	if e.Backend != "vault" {
		t.Errorf("got backend %q, want %q", e.Backend, "vault")
	}
	if e.Key != "db/password" {
		t.Errorf("got key %q, want %q", e.Key, "db/password")
	}
}

func TestGet_MissingKey_ReturnsFalse(t *testing.T) {
	s := newStore(t)
	_, ok := s.Get("vault", "nonexistent")
	if ok {
		t.Error("expected ok=false for missing key")
	}
}

func TestCapture_CapturedAtIsUTC(t *testing.T) {
	s := newStore(t)
	s.Capture("ssm", "app/token", "tok")
	e, _ := s.Get("ssm", "app/token")
	if e.CapturedAt.Location() != time.UTC {
		t.Errorf("expected UTC, got %v", e.CapturedAt.Location())
	}
}

func TestCapture_OverwritesPreviousValue(t *testing.T) {
	s := newStore(t)
	s.Capture("vault", "key", "old")
	s.Capture("vault", "key", "new")
	e, _ := s.Get("vault", "key")
	if e.Value != "new" {
		t.Errorf("expected overwritten value %q, got %q", "new", e.Value)
	}
}

func TestLen_ReturnsCorrectCount(t *testing.T) {
	s := newStore(t)
	s.Capture("vault", "a", "1")
	s.Capture("vault", "b", "2")
	s.Capture("ssm", "a", "3")
	if s.Len() != 3 {
		t.Errorf("expected Len=3, got %d", s.Len())
	}
}

func TestAll_ReturnsAllEntries(t *testing.T) {
	s := newStore(t)
	s.Capture("vault", "x", "v1")
	s.Capture("doppler", "y", "v2")
	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}

func TestDifferentBackends_SameKey_StoredSeparately(t *testing.T) {
	s := newStore(t)
	s.Capture("vault", "shared/key", "from-vault")
	s.Capture("ssm", "shared/key", "from-ssm")

	v, _ := s.Get("vault", "shared/key")
	ss, _ := s.Get("ssm", "shared/key")

	if v.Value != "from-vault" {
		t.Errorf("vault value mismatch: %q", v.Value)
	}
	if ss.Value != "from-ssm" {
		t.Errorf("ssm value mismatch: %q", ss.Value)
	}
}
