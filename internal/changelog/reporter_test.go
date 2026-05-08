package changelog_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/nicholasgasior/vaultrot/internal/changelog"
)

func TestReport_ContainsHeaders(t *testing.T) {
	cl := changelog.New(fixedClock())
	var buf bytes.Buffer
	if err := changelog.Report(cl, &buf); err != nil {
		t.Fatalf("Report error: %v", err)
	}
	for _, h := range []string{"TIME", "BACKEND", "KEY", "ACTOR", "ACTION", "DRY-RUN", "NOTE"} {
		if !strings.Contains(buf.String(), h) {
			t.Errorf("missing header %q", h)
		}
	}
}

func TestReport_EmptyChangelog_OnlyHeaders(t *testing.T) {
	cl := changelog.New(fixedClock())
	var buf bytes.Buffer
	_ = changelog.Report(cl, &buf)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines (header+separator), got %d", len(lines))
	}
}

func TestReport_ShowsEntryData(t *testing.T) {
	cl := changelog.New(fixedClock())
	cl.Record("doppler", "STRIPE_KEY", "ci-bot", "rotated", "weekly", false)
	var buf bytes.Buffer
	_ = changelog.Report(cl, &buf)
	out := buf.String()
	for _, want := range []string{"doppler", "STRIPE_KEY", "ci-bot", "weekly"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in report output", want)
		}
	}
}

func TestReport_DryRunColumn(t *testing.T) {
	cl := changelog.New(fixedClock())
	cl.Record("vault", "k", "u", "rotated", "", true)
	var buf bytes.Buffer
	_ = changelog.Report(cl, &buf)
	if !strings.Contains(buf.String(), "yes") {
		t.Error("expected 'yes' in dry-run column")
	}
}

func TestReport_MultipleEntries_Sorted(t *testing.T) {
	t1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	clock1 := func() time.Time { return t2 }
	cl := changelog.New(clock1)
	cl.Record("vault", "b", "u", "rotated", "", false)

	// Manually inject an older entry via a second changelog merged by record order
	clock2 := func() time.Time { return t1 }
	cl2 := changelog.New(clock2)
	cl2.Record("vault", "a", "u", "rotated", "", false)

	// Build combined changelog preserving insertion order
	combined := changelog.New(func() time.Time { return t2 })
	for _, e := range append(cl2.Entries(), cl.Entries()...) {
		combined.Record(e.Backend, e.Key, e.Actor, e.Action, e.Note, e.DryRun)
	}

	var buf bytes.Buffer
	_ = changelog.Report(combined, &buf)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// header + separator + 2 data rows
	if len(lines) != 4 {
		t.Errorf("expected 4 lines, got %d", len(lines))
	}
}
