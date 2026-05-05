package snapshot_test

import (
	"strings"
	"testing"

	"github.com/yourusername/vaultrot/internal/snapshot"
)

func TestReport_ContainsHeaders(t *testing.T) {
	s := snapshot.New()
	var buf strings.Builder
	snapshot.Report(&buf, s)
	out := buf.String()
	for _, h := range []string{"BACKEND", "KEY", "CAPTURED AT"} {
		if !strings.Contains(out, h) {
			t.Errorf("expected header %q in output", h)
		}
	}
}

func TestReport_EmptyStore_OnlyHeaders(t *testing.T) {
	s := snapshot.New()
	var buf strings.Builder
	snapshot.Report(&buf, s)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// header + separator = 2 lines
	if len(lines) != 2 {
		t.Errorf("expected 2 lines for empty store, got %d", len(lines))
	}
}

func TestReport_ShowsBackendAndKey(t *testing.T) {
	s := snapshot.New()
	s.Capture("vault", "db/password", "secret")
	var buf strings.Builder
	snapshot.Report(&buf, s)
	out := buf.String()
	if !strings.Contains(out, "vault") {
		t.Error("expected backend 'vault' in report")
	}
	if !strings.Contains(out, "db/password") {
		t.Error("expected key 'db/password' in report")
	}
}

func TestReport_MultipleEntries_Sorted(t *testing.T) {
	s := snapshot.New()
	s.Capture("vault", "z/key", "v1")
	s.Capture("vault", "a/key", "v2")
	s.Capture("ssm", "m/key", "v3")
	var buf strings.Builder
	snapshot.Report(&buf, s)
	out := buf.String()
	ssmIdx := strings.Index(out, "ssm")
	vaultIdx := strings.Index(out, "vault")
	if ssmIdx == -1 || vaultIdx == -1 {
		t.Fatal("expected both backends in output")
	}
	if ssmIdx > vaultIdx {
		t.Error("expected ssm to appear before vault (alphabetical sort)")
	}
}
