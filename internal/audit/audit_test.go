package audit

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func newBufferedLogger(dryRun bool) (*Logger, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	return &Logger{writer: buf, dryRun: dryRun}, buf
}

func TestRecord_WritesJSON(t *testing.T) {
	l, buf := newBufferedLogger(false)
	if err := l.Record("vault", "db/password", "rotated", ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var e Entry
	if err := json.Unmarshal(buf.Bytes(), &e); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if e.Backend != "vault" {
		t.Errorf("expected backend vault, got %s", e.Backend)
	}
	if e.SecretKey != "db/password" {
		t.Errorf("expected secret_key db/password, got %s", e.SecretKey)
	}
	if e.Status != "rotated" {
		t.Errorf("expected status rotated, got %s", e.Status)
	}
	if e.DryRun {
		t.Error("expected dry_run false")
	}
}

func TestRecord_DryRunFlag(t *testing.T) {
	l, buf := newBufferedLogger(true)
	_ = l.Record("aws-ssm", "/app/secret", "skipped", "dry-run mode")
	var e Entry
	_ = json.Unmarshal(buf.Bytes(), &e)
	if !e.DryRun {
		t.Error("expected dry_run true")
	}
	if e.Message != "dry-run mode" {
		t.Errorf("expected message 'dry-run mode', got %q", e.Message)
	}
}

func TestRecord_TimestampIsUTC(t *testing.T) {
	l, buf := newBufferedLogger(false)
	before := time.Now().UTC()
	_ = l.Record("doppler", "API_KEY", "rotated", "")
	after := time.Now().UTC()
	var e Entry
	_ = json.Unmarshal(buf.Bytes(), &e)
	if e.Timestamp.Before(before) || e.Timestamp.After(after) {
		t.Errorf("timestamp %v not in expected range [%v, %v]", e.Timestamp, before, after)
	}
}

func TestRecord_MultipleEntries(t *testing.T) {
	l, buf := newBufferedLogger(false)
	_ = l.Record("vault", "secret/one", "rotated", "")
	_ = l.Record("vault", "secret/two", "failed", "connection timeout")
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 log lines, got %d", len(lines))
	}
}

func TestNew_InvalidPath(t *testing.T) {
	_, err := New("/nonexistent/dir/audit.log", false)
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}
