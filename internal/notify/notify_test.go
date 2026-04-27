package notify

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func fixedTime() time.Time {
	t, _ := time.Parse(time.RFC3339, "2024-06-01T12:00:00Z")
	return t
}

func TestNotify_InfoEvent(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf, "")

	err := n.Notify(Event{
		Secret:  "db/password",
		Backend: "vault",
		Level:   LevelInfo,
		Message: "rotated successfully",
		Time:    fixedTime(),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "vault/db/password") {
		t.Errorf("expected backend/secret in output, got: %s", out)
	}
	if !strings.Contains(out, "INFO") {
		t.Errorf("expected INFO level in output, got: %s", out)
	}
}

func TestNotify_DryRunTag(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf, "")

	_ = n.Notify(Event{
		Secret:  "api/key",
		Backend: "doppler",
		Level:   LevelInfo,
		Message: "would rotate",
		DryRun:  true,
		Time:    fixedTime(),
	})

	if !strings.Contains(buf.String(), "[dry-run]") {
		t.Errorf("expected [dry-run] tag in output, got: %s", buf.String())
	}
}

func TestNotify_WithPrefix(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf, "vaultrot")

	_ = n.Notify(Event{
		Secret:  "s",
		Backend: "aws-ssm",
		Level:   LevelWarn,
		Message: "retrying",
		Time:    fixedTime(),
	})

	if !strings.Contains(buf.String(), "[vaultrot]") {
		t.Errorf("expected prefix in output, got: %s", buf.String())
	}
}

func TestNotifyAll_MultipleEvents(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf, "")

	events := []Event{
		{Secret: "a", Backend: "vault", Level: LevelInfo, Message: "ok", Time: fixedTime()},
		{Secret: "b", Backend: "vault", Level: LevelError, Message: "failed", Time: fixedTime()},
	}

	if err := n.NotifyAll(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	n := New(nil, "")
	if n.out == nil {
		t.Error("expected non-nil writer")
	}
}
