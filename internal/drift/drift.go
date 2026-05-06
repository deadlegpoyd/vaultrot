// Package drift detects configuration drift between the expected secret
// state defined in vaultrot config and the live values held by a backend.
package drift

import (
	"fmt"
	"time"
)

// Status represents the drift state of a single secret.
type Status int

const (
	StatusMatch   Status = iota // live value matches expected
	StatusMissing               // secret does not exist in the backend
	StatusStale                 // secret exists but is older than the max age
	StatusUnknown               // value could not be retrieved
)

func (s Status) String() string {
	switch s {
	case StatusMatch:
		return "OK"
	case StatusMissing:
		return "MISSING"
	case StatusStale:
		return "STALE"
	default:
		return "UNKNOWN"
	}
}

// Entry records the drift result for one secret.
type Entry struct {
	Name      string
	Backend   string
	Status    Status
	LastRotAt time.Time
	MaxAge    time.Duration
	Note      string
}

// Report holds all drift entries produced during a single detection run.
type Report struct {
	GeneratedAt time.Time
	Entries     []Entry
}

// Detector accumulates drift entries and produces a Report.
type Detector struct {
	now func() time.Time
}

// New returns a Detector. Passing a nil clock falls back to time.Now.
func New(clock func() time.Time) *Detector {
	if clock == nil {
		clock = time.Now
	}
	return &Detector{now: clock}
}

// Check evaluates whether a secret is drifted given its last-rotation
// timestamp and the allowed maximum age. An empty lastRotAt is treated as
// a missing secret.
func (d *Detector) Check(name, backend string, lastRotAt time.Time, maxAge time.Duration) Entry {
	e := Entry{
		Name:      name,
		Backend:   backend,
		LastRotAt: lastRotAt,
		MaxAge:    maxAge,
	}

	if lastRotAt.IsZero() {
		e.Status = StatusMissing
		e.Note = "no rotation record found"
		return e
	}

	age := d.now().Sub(lastRotAt)
	if age > maxAge {
		e.Status = StatusStale
		e.Note = fmt.Sprintf("age %s exceeds max %s", age.Round(time.Second), maxAge)
		return e
	}

	e.Status = StatusMatch
	return e
}

// Build finalises and returns a Report from the supplied entries.
func (d *Detector) Build(entries []Entry) Report {
	return Report{
		GeneratedAt: d.now(),
		Entries:     entries,
	}
}

// HasDrift returns true when at least one entry is not StatusMatch.
func (r *Report) HasDrift() bool {
	for _, e := range r.Entries {
		if e.Status != StatusMatch {
			return true
		}
	}
	return false
}

// Summary returns a concise breakdown of entry statuses in the report,
// mapping each Status to the number of entries carrying that status.
func (r *Report) Summary() map[Status]int {
	counts := make(map[Status]int)
	for _, e := range r.Entries {
		counts[e.Status]++
	}
	return counts
}
