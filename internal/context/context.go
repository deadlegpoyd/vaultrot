// Package context provides a rotation context that carries per-run metadata
// such as the run ID, dry-run flag, and start time through the call chain.
package context

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type key int

const rotationKey key = iota

// RotationContext holds metadata associated with a single rotation run.
type RotationContext struct {
	RunID     string
	DryRun    bool
	StartedAt time.Time
	Operator  string
}

// New creates a RotationContext with a generated run ID and the current UTC time.
func New(dryRun bool, operator string) *RotationContext {
	return &RotationContext{
		RunID:     newRunID(),
		DryRun:    dryRun,
		StartedAt: time.Now().UTC(),
		Operator:  operator,
	}
}

// Inject stores the RotationContext inside a standard context.Context.
func Inject(ctx context.Context, rc *RotationContext) context.Context {
	return context.WithValue(ctx, rotationKey, rc)
}

// Extract retrieves the RotationContext from a standard context.Context.
// Returns nil, false if no RotationContext is present.
func Extract(ctx context.Context) (*RotationContext, bool) {
	rc, ok := ctx.Value(rotationKey).(*RotationContext)
	return rc, ok
}

// MustExtract retrieves the RotationContext or panics if it is absent.
func MustExtract(ctx context.Context) *RotationContext {
	rc, ok := Extract(ctx)
	if !ok {
		panic("vaultrot: RotationContext not found in context")
	}
	return rc
}

// String returns a human-readable summary of the rotation context.
func (rc *RotationContext) String() string {
	mode := "live"
	if rc.DryRun {
		mode = "dry-run"
	}
	return fmt.Sprintf("run=%s mode=%s operator=%s started=%s",
		rc.RunID, mode, rc.Operator, rc.StartedAt.Format(time.RFC3339))
}

func newRunID() string {
	return uuid.NewString()
}
