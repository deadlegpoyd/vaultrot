// Package policy provides rotation policy enforcement for vaultrot.
//
// A Policy is constructed from a Config that specifies:
//
//   - MinAgeDays: the minimum number of days that must have elapsed since the
//     last rotation before a new rotation is permitted. Useful for preventing
//     accidental rapid re-rotation.
//
//   - MaxAgeDays: the maximum number of days allowed between rotations. When
//     exceeded, RequiresRotation returns true and Check returns an error.
//
// Example:
//
//	p, err := policy.New(policy.Config{MinAgeDays: 1, MaxAgeDays: 90}, nil)
//	if err != nil { ... }
//	if err := p.Check(secretName, lastRotatedAt); err != nil { ... }
package policy
