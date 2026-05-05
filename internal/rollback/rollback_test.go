package rollback_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultrot/internal/rollback"
)

func TestRecord_StoresSnapshot(t *testing.T) {
	s := rollback.New()
	s.Record("vault", "db/password", "old-secret")

	snap, ok := s.Get("vault", "db/password")
	if !ok {
		t.Fatal("expected snapshot to exist")
	}
	if snap.OldValue != "old-secret" {
		t.Errorf("got OldValue %q, want %q", snap.OldValue, "old-secret")
	}
	if snap.Backend != "vault" {
		t.Errorf("got Backend %q, want %q", snap.Backend, "vault")
	}
	if snap.SecretKey != "db/password" {
		t.Errorf("got SecretKey %q, want %q", snap.SecretKey, "db/password")
	}
}

func TestRecord_CapturedAtIsUTC(t *testing.T) {
	s := rollback.New()
	before := time.Now().UTC()
	s.Record("aws-ssm", "/prod/api-key", "v1")
	after := time.Now().UTC()

	snap, _ := s.Get("aws-ssm", "/prod/api-key")
	if snap.CapturedAt.Before(before) || snap.CapturedAt.After(after) {
		t.Errorf("CapturedAt %v not within expected range", snap.CapturedAt)
	}
	if snap.CapturedAt.Location() != time.UTC {
		t.Error("CapturedAt should be UTC")
	}
}

func TestGet_MissingKeyReturnsFalse(t *testing.T) {
	s := rollback.New()
	_, ok := s.Get("doppler", "nonexistent")
	if ok {
		t.Error("expected ok=false for missing key")
	}
}

func TestRecord_OverwritesPreviousSnapshot(t *testing.T) {
	s := rollback.New()
	s.Record("vault", "secret/token", "first")
	s.Record("vault", "secret/token", "second")

	snap, _ := s.Get("vault", "secret/token")
	if snap.OldValue != "second" {
		t.Errorf("expected overwritten value %q, got %q", "second", snap.OldValue)
	}
}

func TestLen_ReturnsCorrectCount(t *testing.T) {
	s := rollback.New()
	if s.Len() != 0 {
		t.Errorf("expected 0, got %d", s.Len())
	}
	s.Record("vault", "a", "1")
	s.Record("aws-ssm", "b", "2")
	if s.Len() != 2 {
		t.Errorf("expected 2, got %d", s.Len())
	}
}

func TestKeys_ReturnsAllCompositeKeys(t *testing.T) {
	s := rollback.New()
	s.Record("vault", "x", "val")
	s.Record("doppler", "y", "val")

	keys := s.Keys()
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
}
