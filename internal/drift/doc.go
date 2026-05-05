// Package drift provides secret-drift detection for vaultrot.
//
// A Detector compares the last-rotation timestamp of each managed secret
// against a configurable maximum age. Secrets that have never been rotated
// are reported as MISSING; secrets whose age exceeds the threshold are
// reported as STALE; all others are OK.
//
// Usage:
//
//	detector := drift.New(nil) // nil clock → time.Now
//
//	entry := detector.Check("db/password", "vault", lastRotated, 7*24*time.Hour)
//	report := detector.Build([]drift.Entry{entry})
//	if report.HasDrift() {
//	    // surface findings via notify or audit
//	}
package drift
