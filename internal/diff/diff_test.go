package diff

import (
	"strings"
	"testing"
)

func TestNew_EmptyResult(t *testing.T) {
	r := New()
	if len(r.Entries()) != 0 {
		t.Errorf("expected 0 entries, got %d", len(r.Entries()))
	}
}

func TestAdd_StoresEntry(t *testing.T) {
	r := New()
	r.Add("db/password", "old123", "new456", false)
	entries := r.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Key != "db/password" || e.OldValue != "old123" || e.NewValue != "new456" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestHasChanges_ReturnsTrueOnDiff(t *testing.T) {
	r := New()
	r.Add("key", "old", "new", false)
	if !r.HasChanges() {
		t.Error("expected HasChanges to return true")
	}
}

func TestHasChanges_ReturnsFalseWhenSame(t *testing.T) {
	r := New()
	r.Add("key", "same", "same", false)
	if r.HasChanges() {
		t.Error("expected HasChanges to return false")
	}
}

func TestFormat_NoEntries(t *testing.T) {
	r := New()
	out := r.Format()
	if out != "(no changes)" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormat_MaskedValues(t *testing.T) {
	r := New()
	r.Add("api/key", "secret-old", "secret-new", true)
	out := r.Format()
	if strings.Contains(out, "secret-old") || strings.Contains(out, "secret-new") {
		t.Error("masked values should not appear in formatted output")
	}
	if !strings.Contains(out, "***") {
		t.Error("expected masked placeholder '***' in output")
	}
}

func TestFormat_ShowsUnchangedLabel(t *testing.T) {
	r := New()
	r.Add("token", "abc", "abc", false)
	out := r.Format()
	if !strings.Contains(out, "unchanged") {
		t.Errorf("expected 'unchanged' label, got: %s", out)
	}
}

func TestFormat_ShowsPlusMinus(t *testing.T) {
	r := New()
	r.Add("pass", "old", "new", false)
	out := r.Format()
	if !strings.Contains(out, "- pass") || !strings.Contains(out, "+ pass") {
		t.Errorf("expected diff markers in output: %s", out)
	}
}

func TestEntries_ReturnsCopy(t *testing.T) {
	r := New()
	r.Add("k", "a", "b", false)
	e1 := r.Entries()
	e1[0].Key = "mutated"
	e2 := r.Entries()
	if e2[0].Key == "mutated" {
		t.Error("Entries should return a copy, not a reference")
	}
}
