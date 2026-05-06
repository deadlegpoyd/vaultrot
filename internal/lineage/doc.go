// Package lineage provides rotation history tracking for secrets managed by
// vaultrot. Each time a secret is rotated, a lineage.Entry is appended to an
// in-memory chain for that backend/key pair. Entries are linked via PreviousID
// so the full rotation ancestry can be reconstructed at any time.
//
// Basic usage:
//
//	tr := lineage.New()
//	entry, err := tr.Record("vault", "db/password", "auto")
//	history := tr.History("vault", "db/password")
package lineage
