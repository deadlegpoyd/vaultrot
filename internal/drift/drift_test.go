package drift_test

import (
	"testing"
	"time"

	"github.com/rodfernandez/vaultrot/internal/drift"
)

var (
	fixedNow = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	maxAge   = 7 * 24 * time.Hour
)

func newDetector() *drift.Detector {
	return drift.New(func() time.Time { return fixedNow })
}

func TestCheck_MatchStatus_WhenFresh(t *testing.T) {
	d := newDetector()
	lastRot := fixedNow.Add(-48 * time.Hour) // 2 days old, within 7-day max
	e := d.Check("api/key", "vault", lastRot, maxAge)
	if e.Status != drift.StatusMatch {
		t.Fatalf("expected StatusMatch, got %s", e.Status)
	}
}

func TestCheck_StaleStatus_WhenTooOld(t *testing.T) {
	d := newDetector()
	lastRot := fixedNow.Add(-10 * 24 * time.Hour) // 10 days old
	e := d.Check("db/pass", "aws-ssm", lastRot, maxAge)
	if e.Status != drift.StatusStale {
		t.Fatalf("expected StatusStale, got %s", e.Status)
	}
	if e.Note == "" {
		t.Error("expected non-empty note for stale entry")
	}
}

func TestCheck_MissingStatus_WhenZeroTime(t *testing.T) {
	d := newDetector()
	e := d.Check("svc/token", "doppler", time.Time{}, maxAge)
	if e.Status != drift.StatusMissing {
		t.Fatalf("expected StatusMissing, got %s", e.Status)
	}
}

func TestCheck_EntryFieldsPopulated(t *testing.T) {
	d := newDetector()
	lastRot := fixedNow.Add(-1 * time.Hour)
	e := d.Check("x/y", "vault", lastRot, maxAge)
	if e.Name != "x/y" || e.Backend != "vault" {
		t.Errorf("unexpected entry fields: %+v", e)
	}
}

func TestBuild_SetsGeneratedAt(t *testing.T) {
	d := newDetector()
	r := d.Build(nil)
	if !r.GeneratedAt.Equal(fixedNow) {
		t.Errorf("expected GeneratedAt=%v, got %v", fixedNow, r.GeneratedAt)
	}
}

func TestHasDrift_FalseWhenAllMatch(t *testing.T) {
	d := newDetector()
	entries := []drift.Entry{
		d.Check("a", "vault", fixedNow.Add(-1*time.Hour), maxAge),
		d.Check("b", "vault", fixedNow.Add(-2*time.Hour), maxAge),
	}
	r := d.Build(entries)
	if r.HasDrift() {
		t.Error("expected no drift")
	}
}

func TestHasDrift_TrueWhenStalePresent(t *testing.T) {
	d := newDetector()
	entries := []drift.Entry{
		d.Check("a", "vault", fixedNow.Add(-1*time.Hour), maxAge),
		d.Check("b", "vault", fixedNow.Add(-30*24*time.Hour), maxAge),
	}
	r := d.Build(entries)
	if !r.HasDrift() {
		t.Error("expected drift to be detected")
	}
}

func TestStatus_String(t *testing.T) {
	cases := map[drift.Status]string{
		drift.StatusMatch:   "OK",
		drift.StatusMissing: "MISSING",
		drift.StatusStale:   "STALE",
		drift.StatusUnknown: "UNKNOWN",
	}
	for s, want := range cases {
		if got := s.String(); got != want {
			t.Errorf("Status(%d).String() = %q, want %q", s, got, want)
		}
	}
}
