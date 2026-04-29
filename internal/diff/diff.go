// Package diff provides utilities for comparing secret values
// before and after rotation, supporting dry-run previews.
package diff

import (
	"fmt"
	"strings"
)

// Entry represents a single secret diff between old and new values.
type Entry struct {
	Key      string
	OldValue string
	NewValue string
	Masked   bool
}

// Result holds a collection of diff entries for a rotation run.
type Result struct {
	entries []Entry
}

// New creates a new empty diff Result.
func New() *Result {
	return &Result{}
}

// Add appends a diff entry. If masked is true, values are redacted in output.
func (r *Result) Add(key, oldVal, newVal string, masked bool) {
	r.entries = append(r.entries, Entry{
		Key:      key,
		OldValue: oldVal,
		NewValue: newVal,
		Masked:   masked,
	})
}

// Entries returns a copy of all recorded diff entries.
func (r *Result) Entries() []Entry {
	out := make([]Entry, len(r.entries))
	copy(out, r.entries)
	return out
}

// HasChanges returns true if any entry differs between old and new.
func (r *Result) HasChanges() bool {
	for _, e := range r.entries {
		if e.OldValue != e.NewValue {
			return true
		}
	}
	return false
}

// Format returns a human-readable diff string suitable for dry-run output.
func (r *Result) Format() string {
	if len(r.entries) == 0 {
		return "(no changes)"
	}
	var sb strings.Builder
	for _, e := range r.entries {
		old := mask(e.OldValue, e.Masked)
		new_ := mask(e.NewValue, e.Masked)
		if e.OldValue == e.NewValue {
			fmt.Fprintf(&sb, "  ~ %s: %s (unchanged)\n", e.Key, old)
		} else {
			fmt.Fprintf(&sb, "  - %s: %s\n", e.Key, old)
			fmt.Fprintf(&sb, "  + %s: %s\n", e.Key, new_)
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

func mask(val string, masked bool) string {
	if !masked {
		return val
	}
	if len(val) == 0 {
		return "(empty)"
	}
	return "***"
}
