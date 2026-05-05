// Package lease provides lease tracking for rotated secrets.
//
// A Tracker records the issuance and expiry time of each managed secret,
// enabling the rotation engine to identify secrets that have expired or
// are approaching their renewal window before the next scheduled rotation.
//
// Usage:
//
//	tr := lease.New()
//	_ = tr.Record(lease.Entry{
//		SecretName: "db/password",
//		Backend:    "vault",
//		IssuedAt:   time.Now(),
//		ExpiresAt:  time.Now().Add(24 * time.Hour),
//		Renewable:  true,
//	})
//
//	due := tr.DueForRenewal(30 * time.Minute)
package lease
