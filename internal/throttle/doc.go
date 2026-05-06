// Package throttle provides per-backend concurrency limiting for the vaultrot
// rotation pipeline.
//
// A Throttle is created with a maximum concurrency limit and controls how many
// rotation operations may run simultaneously against a single named backend
// (e.g. "vault", "aws-ssm", "doppler"). Each backend maintains its own
// semaphore, so operations against different backends are fully independent.
//
// Typical usage:
//
//	th, err := throttle.New(throttle.Config{MaxConcurrent: 5})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if err := th.Acquire(ctx, backendName); err != nil {
//		return err
//	}
//	defer th.Release(backendName)
//	// ... perform rotation ...
package throttle
