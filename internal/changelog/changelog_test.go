package changelog_test

import {
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/nicholasgasior/vaultrot/internal/changelog"
)

var fixedTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func fixedClock() func() time.Time {
	return func() time.Time { return fixedTime }
}

func TestRecord_AppendsEntry(t *testing.T) {
	cl := changelog.New(fixedClock())
	cl.Record("vault", "db/password", "ci-bot", "rotated", "", false)
	if cl.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", cl.Len())
	}
}

func TestRecord_EntryFields(t *testing.T) {
	cl := changelog.New(fixedClock())
	cl.Record("aws-ssm", "/prod/token", "admin", "rotated", "scheduled", false)
	entries := cl.Entries()
	e := entries[0]
	if e.Backend != "aws-ssm" {
		t.Errorf("backend: got %q", e.Backend)
	}
	if e.Key != "/prod/token" {
		t.Errorf("key: got %q", e.Key)
	}
	if e.Actor != "admin" {
		t.Errorf("actor: got %q", e.Actor)
	}
	if e.Note != "scheduled" {
		t.Errorf("note: got %q", e.Note)
	}
	if !e.OccurredAt.Equal(fixedTime) {
		t.Errorf("time: got %v", e.OccurredAt)
	}
}

func TestRecord_DryRunFlag(t *testing.T) {
	cl := changelog.New(fixedClock())
	cl.Record("doppler", "API_KEY", "bot", "rotated", "", true)
	if !cl.Entries()[0].DryRun {
		t.Error("expected DryRun to be true")
	}
}

func TestLen_MultipleEntries(t *testing.T) {
	cl := changelog.New(fixedClock())
	for i := 0; i < 5; i++ {
		cl.Record("vault", "secret", "bot", "rotated", "", false)
	}
	if cl.Len() != 5 {
		t.Errorf("expected 5, got %d", cl.Len())
	}
}

func TestWrite_ContainsKeyFields(t *testing.T) {
	cl := changelog.New(fixedClock())
	cl.Record("vault", "db/pass", "ci", "rotated", "routine", false)
	var buf bytes.Buffer
	if err := cl.Write(&buf); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"vault/db/pass", "ci", "rotated", "routine", "2024-06-01"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\ngot: %s", want, out)
		}
	}
}

func TestWrite_DryRunTag(t *testing.T) {
	cl := changelog.New(fixedClock())
	cl.Record("vault", "key", "bot", "rotated", "", true)
	var buf bytes.Buffer
	_ = cl.Write(&buf)
	if !strings.Contains(buf.String(), "[dry-run]") {
		t.Errorf("expected [dry-run] tag in output")
	}
}

func TestWrite_EmptyChangelog(t *testing.T) {
	cl := changelog.New(fixedClock())
	var buf bytes.Buffer
	if err := cl.Write(&buf); err != nil {
		t.Fatalf("Write on empty changelog returned error: %v", err)
	}
	// An empty changelog should produce no output.
	if buf.Len() != 0 {
		t.Errorf("expected empty output for empty changelog, got: %q", buf.String())
	}
}

func TestEntries_ReturnsCopy(t *testing.T) {
	cl := changelog.New(fixedClock())
	cl.Record("vault", "k", "u", "rotated", "", false)
	a := cl.Entries()
	a[0].Actor = "tampered"
	if cl.Entries()[0].Actor == "tampered" {
		t.Error("Entries should return a copy, not a reference")
	}
}
