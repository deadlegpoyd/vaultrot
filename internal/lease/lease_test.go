package lease

import (
	"testing"
	"time"
)

var (
	fixedNow  = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	baseEntry = Entry{
		SecretName: "db/password",
		Backend:    "vault",
		IssuedAt:   fixedNow.Add(-1 * time.Hour),
		ExpiresAt:  fixedNow.Add(1 * time.Hour),
		Renewable:  true,
	}
)

func newTracker(now time.Time) *Tracker {
	t := New()
	t.now = func() time.Time { return now }
	return t
}

func TestRecord_And_Get(t *testing.T) {
	tr := newTracker(fixedNow)
	if err := tr.Record(baseEntry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := tr.Get("db/password")
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if e.Backend != "vault" {
		t.Errorf("got backend %q, want vault", e.Backend)
	}
}

func TestRecord_EmptyName_ReturnsError(t *testing.T) {
	tr := newTracker(fixedNow)
	err := tr.Record(Entry{ExpiresAt: fixedNow.Add(time.Hour)})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestRecord_ZeroExpiry_ReturnsError(t *testing.T) {
	tr := newTracker(fixedNow)
	err := tr.Record(Entry{SecretName: "x"})
	if err == nil {
		t.Fatal("expected error for zero expiry")
	}
}

func TestExpired_ReturnsExpiredEntries(t *testing.T) {
	tr := newTracker(fixedNow)
	_ = tr.Record(baseEntry) // expires in +1h, not expired
	expired := Entry{
		SecretName: "old/key",
		Backend:    "aws-ssm",
		IssuedAt:   fixedNow.Add(-3 * time.Hour),
		ExpiresAt:  fixedNow.Add(-1 * time.Hour),
	}
	_ = tr.Record(expired)
	list := tr.Expired()
	if len(list) != 1 || list[0].SecretName != "old/key" {
		t.Errorf("expected one expired entry, got %v", list)
	}
}

func TestDueForRenewal_WithinThreshold(t *testing.T) {
	tr := newTracker(fixedNow)
	// expires in 10 minutes — within a 30-minute threshold
	e := baseEntry
	e.ExpiresAt = fixedNow.Add(10 * time.Minute)
	_ = tr.Record(e)
	due := tr.DueForRenewal(30 * time.Minute)
	if len(due) != 1 {
		t.Errorf("expected 1 renewal candidate, got %d", len(due))
	}
}

func TestDueForRenewal_NotRenewable_Excluded(t *testing.T) {
	tr := newTracker(fixedNow)
	e := baseEntry
	e.ExpiresAt = fixedNow.Add(5 * time.Minute)
	e.Renewable = false
	_ = tr.Record(e)
	due := tr.DueForRenewal(30 * time.Minute)
	if len(due) != 0 {
		t.Errorf("expected 0 renewal candidates for non-renewable lease")
	}
}

func TestRemove_DeletesEntry(t *testing.T) {
	tr := newTracker(fixedNow)
	_ = tr.Record(baseEntry)
	tr.Remove("db/password")
	if _, ok := tr.Get("db/password"); ok {
		t.Error("expected entry to be removed")
	}
}

func TestLen_ReturnsCount(t *testing.T) {
	tr := newTracker(fixedNow)
	_ = tr.Record(baseEntry)
	if tr.Len() != 1 {
		t.Errorf("expected len 1, got %d", tr.Len())
	}
}
