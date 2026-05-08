// Package changelog provides an ordered, thread-safe log of secret rotation
// events for vaultrot.
//
// Each rotation event is captured as an [Entry] containing the backend,
// secret key, acting principal, action taken, timestamp, dry-run flag, and
// an optional human-readable note.
//
// Usage:
//
//	cl := changelog.New(nil) // nil uses time.Now
//	cl.Record("vault", "db/password", "ci-bot", "rotated", "scheduled", false)
//
//	// Pretty-print a table to stdout
//	changelog.Report(cl, os.Stdout)
//
// The changelog is safe for concurrent use.
package changelog
