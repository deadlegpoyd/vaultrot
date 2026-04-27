package rotation

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestSummary_Counts(t *testing.T) {
	s := NewSummary(false)
	s.Add(Result{Status: StatusSuccess})
	s.Add(Result{Status: StatusSuccess})
	s.Add(Result{Status: StatusSkipped})
	s.Add(Result{Status: StatusFailed, Error: errors.New("boom")})

	success, skipped, failed := s.Counts()
	if success != 2 {
		t.Errorf("expected 2 successes, got %d", success)
	}
	if skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", skipped)
	}
	if failed != 1 {
		t.Errorf("expected 1 failed, got %d", failed)
	}
}

func TestSummary_Print_DryRun(t *testing.T) {
	s := NewSummary(true)
	s.Add(Result{
		SecretName: "my-secret",
		Backend:    "vault",
		Status:     StatusSkipped,
		DryRun:     true,
	})
	s.Finish()

	var buf bytes.Buffer
	s.Print(&buf)
	out := buf.String()

	if !strings.Contains(out, "DRY-RUN") {
		t.Error("expected DRY-RUN tag in output")
	}
	if !strings.Contains(out, "my-secret") {
		t.Error("expected secret name in output")
	}
	if !strings.Contains(out, string(StatusSkipped)) {
		t.Error("expected skipped status in output")
	}
}

func TestSummary_Print_WithError(t *testing.T) {
	s := NewSummary(false)
	s.Add(Result{
		SecretName: "broken-secret",
		Backend:    "aws-ssm",
		Status:     StatusFailed,
		Error:      errors.New("access denied"),
	})
	s.Finish()

	var buf bytes.Buffer
	s.Print(&buf)
	out := buf.String()

	if !strings.Contains(out, "access denied") {
		t.Error("expected error message in output")
	}
	if !strings.Contains(out, "failed") {
		t.Error("expected failed count in summary line")
	}
}

func TestSummary_Finish_SetsTime(t *testing.T) {
	s := NewSummary(false)
	before := time.Now()
	s.Finish()
	after := time.Now()

	if s.Finished.Before(before) || s.Finished.After(after) {
		t.Error("Finished timestamp is out of expected range")
	}
}
