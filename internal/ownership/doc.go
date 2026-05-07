// Package ownership provides a thread-safe registry that maps secret keys to
// the teams or services responsible for them.
//
// During a rotation run, vaultrot consults the registry to:
//   - Route post-rotation notifications to the correct contact.
//   - Attach owner metadata to audit log entries.
//   - Enforce access-control policies that restrict which backends a given
//     owner is permitted to rotate.
//
// Usage:
//
//	reg := ownership.New()
//	_ = reg.Register("db/password", ownership.Owner{
//	    Name:    "platform-team",
//	    Contact: "platform@example.com",
//	})
//	owner, ok := reg.Lookup("db/password")
package ownership
