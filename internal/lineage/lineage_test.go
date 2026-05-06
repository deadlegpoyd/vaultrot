package lineage

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func newTracker(t *testing.T) *Tracker {
	t.Helper()
	tr := New()
	tr.nowFn = fixedClock(time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC))
	return tr
}

func TestRecord_StoresEntry(t *testing.T) {
	tr := newTracker(t)
	e, err := tr.Record("vault", "db/password", "auto")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Backend != "vault" || e.Key != "db/password" {
		t.Errorf("unexpected entry fields: %+v", e)
	}
	if e.GeneratedBy != "auto" {
		t.Errorf("expected GeneratedBy=auto, got %q", e.GeneratedBy)
	}
}

func TestRecord_FirstEntry_HasNoPreviousID(t *testing.T) {
	tr := newTracker(t)
	e, _ := tr.Record("vault", "api/key", "manual")
	if e.PreviousID != "" {
		t.Errorf("expected empty PreviousID for first rotation, got %q", e.PreviousID)
	}
}

func TestRecord_SecondEntry_LinksToPrevious(t *testing.T) {
	tr := newTracker(t)
	first, _ := tr.Record("vault", "api/key", "auto")
	second, _ := tr.Record("vault", "api/key", "auto")
	if second.PreviousID != first.ID {
		t.Errorf("expected PreviousID=%q, got %q", first.ID, second.PreviousID)
	}
}

func TestRecord_TimestampIsUTC(t *testing.T) {
	tr := newTracker(t)
	e, _ := tr.Record("ssm", "/prod/token", "scheduled")
	if e.RotatedAt.Location() != time.UTC {
		t.Errorf("expected UTC timestamp, got %v", e.RotatedAt.Location())
	}
}

func TestRecord_EmptyBackend_ReturnsError(t *testing.T) {
	tr := newTracker(t)
	_, err := tr.Record("", "key", "auto")
	if err == nil {
		t.Fatal("expected error for empty backend")
	}
}

func TestRecord_EmptyKey_ReturnsError(t *testing.T) {
	tr := newTracker(t)
	_, err := tr.Record("vault", "", "auto")
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestHistory_ReturnsChronologicalEntries(t *testing.T) {
	tr := newTracker(t)
	tr.Record("vault", "secret", "auto")
	tr.Record("vault", "secret", "manual")
	tr.Record("vault", "secret", "scheduled")

	h := tr.History("vault", "secret")
	if len(h) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(h))
	}
	if h[0].GeneratedBy != "auto" || h[2].GeneratedBy != "scheduled" {
		t.Errorf("unexpected order: %v", h)
	}
}

func TestHistory_MissingKey_ReturnsNil(t *testing.T) {
	tr := newTracker(t)
	if tr.History("vault", "nonexistent") != nil {
		t.Error("expected nil for unknown key")
	}
}

func TestLen_CountsAllEntries(t *testing.T) {
	tr := newTracker(t)
	tr.Record("vault", "a", "auto")
	tr.Record("vault", "b", "auto")
	tr.Record("vault", "a", "manual")
	if tr.Len() != 3 {
		t.Errorf("expected Len=3, got %d", tr.Len())
	}
}
