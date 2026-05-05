package lease

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestReport_ContainsHeaders(t *testing.T) {
	tr := newTracker(fixedNow)
	_ = tr.Record(baseEntry)
	var buf bytes.Buffer
	if err := Report(&buf, tr, fixedNow, 30*time.Minute); err != nil {
		t.Fatalf("Report returned error: %v", err)
	}
	out := buf.String()
	for _, col := range []string{"SECRET", "BACKEND", "EXPIRES", "STATUS"} {
		if !strings.Contains(out, col) {
			t.Errorf("expected column %q in output", col)
		}
	}
}

func TestReport_ShowsOKStatus(t *testing.T) {
	tr := newTracker(fixedNow)
	_ = tr.Record(baseEntry) // expires in +1h, threshold 30m
	var buf bytes.Buffer
	_ = Report(&buf, tr, fixedNow, 30*time.Minute)
	if !strings.Contains(buf.String(), "OK") {
		t.Errorf("expected OK status, got: %s", buf.String())
	}
}

func TestReport_ShowsExpiredStatus(t *testing.T) {
	tr := newTracker(fixedNow)
	e := baseEntry
	e.ExpiresAt = fixedNow.Add(-10 * time.Minute)
	_ = tr.Record(e)
	var buf bytes.Buffer
	_ = Report(&buf, tr, fixedNow, 30*time.Minute)
	if !strings.Contains(buf.String(), "EXPIRED") {
		t.Errorf("expected EXPIRED status, got: %s", buf.String())
	}
}

func TestReport_ShowsRenewSoonStatus(t *testing.T) {
	tr := newTracker(fixedNow)
	e := baseEntry
	e.ExpiresAt = fixedNow.Add(10 * time.Minute) // within 30m threshold
	_ = tr.Record(e)
	var buf bytes.Buffer
	_ = Report(&buf, tr, fixedNow, 30*time.Minute)
	if !strings.Contains(buf.String(), "RENEW_SOON") {
		t.Errorf("expected RENEW_SOON status, got: %s", buf.String())
	}
}

func TestReport_EmptyTracker_OnlyHeaders(t *testing.T) {
	tr := newTracker(fixedNow)
	var buf bytes.Buffer
	_ = Report(&buf, tr, fixedNow, 30*time.Minute)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Errorf("expected only header line for empty tracker, got %d lines", len(lines))
	}
}
