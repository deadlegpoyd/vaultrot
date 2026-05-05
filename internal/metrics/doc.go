// Package metrics provides a lightweight, thread-safe Collector for tracking
// rotation run statistics such as the number of secrets rotated, skipped, or
// failed, as well as per-backend write durations.
//
// Usage:
//
//	c := metrics.New()
//	c.Inc("rotated")
//	c.Inc("skipped")
//	c.Add("errors", 1)
//	c.RecordDuration("vault_write", elapsed)
//	c.Print(os.Stdout)
//
// The Collector is safe for concurrent use. All counter and duration names
// are arbitrary strings chosen by the caller.
package metrics
