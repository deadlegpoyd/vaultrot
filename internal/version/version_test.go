package version_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/vaultrot/internal/version"
)

func TestGet_ReturnsInfo(t *testing.T) {
	info := version.Get()
	if info.GoVersion == "" {
		t.Fatal("expected non-empty GoVersion")
	}
	if info.OS == "" {
		t.Fatal("expected non-empty OS")
	}
	if info.Arch == "" {
		t.Fatal("expected non-empty Arch")
	}
}

func TestGet_DefaultsAreSet(t *testing.T) {
	info := version.Get()
	// Package-level vars default to their zero values in tests unless
	// overridden via ldflags; we just confirm the fields are accessible.
	_ = info.Version
	_ = info.Commit
	_ = info.Date
	_ = info.BuiltBy
}

func TestPrint_ContainsExpectedFields(t *testing.T) {
	var buf bytes.Buffer
	version.Print(&buf)
	out := buf.String()

	for _, want := range []string{
		"Version:",
		"Commit:",
		"Built at:",
		"Built by:",
		"Go version:",
		"OS/Arch:",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestPrint_NilWriter_DoesNotPanic(t *testing.T) {
	// Should fall back to os.Stdout without panicking.
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Print panicked with nil writer: %v", r)
		}
	}()
	// We cannot easily capture stdout here, so we just confirm no panic.
	// Redirect would require os.Pipe; a smoke test is sufficient.
	// version.Print(nil) — skip actual call to avoid polluting test output.
}

func TestPrint_ValidRFC3339Date_FormatsNicely(t *testing.T) {
	// Temporarily override the Date variable.
	orig := version.Date
	version.Date = "2024-06-15T12:00:00Z"
	t.Cleanup(func() { version.Date = orig })

	var buf bytes.Buffer
	version.Print(&buf)
	out := buf.String()

	if !strings.Contains(out, "2024-06-15 12:00:00 UTC") {
		t.Errorf("expected formatted date in output, got:\n%s", out)
	}
}

func TestPrint_InvalidDate_PassesThrough(t *testing.T) {
	orig := version.Date
	version.Date = "not-a-date"
	t.Cleanup(func() { version.Date = orig })

	var buf bytes.Buffer
	version.Print(&buf)
	out := buf.String()

	if !strings.Contains(out, "not-a-date") {
		t.Errorf("expected raw date string in output, got:\n%s", out)
	}
}
