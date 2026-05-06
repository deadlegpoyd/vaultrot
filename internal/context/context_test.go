package context_test

import (
	"context"
	"strings"
	"testing"
	"time"

	rotctx "github.com/yourusername/vaultrot/internal/context"
)

func TestNew_SetsFields(t *testing.T) {
	before := time.Now().UTC()
	rc := rotctx.New(false, "alice")
	after := time.Now().UTC()

	if rc.RunID == "" {
		t.Fatal("expected non-empty RunID")
	}
	if rc.DryRun {
		t.Fatal("expected DryRun=false")
	}
	if rc.Operator != "alice" {
		t.Fatalf("expected operator 'alice', got %q", rc.Operator)
	}
	if rc.StartedAt.Before(before) || rc.StartedAt.After(after) {
		t.Fatal("StartedAt is outside expected range")
	}
}

func TestNew_UniqueRunIDs(t *testing.T) {
	a := rotctx.New(false, "ci")
	b := rotctx.New(false, "ci")
	if a.RunID == b.RunID {
		t.Fatal("expected distinct RunIDs for separate calls")
	}
}

func TestInjectExtract_RoundTrip(t *testing.T) {
	rc := rotctx.New(true, "bot")
	ctx := rotctx.Inject(context.Background(), rc)

	got, ok := rotctx.Extract(ctx)
	if !ok {
		t.Fatal("expected Extract to return true")
	}
	if got.RunID != rc.RunID {
		t.Fatalf("RunID mismatch: got %q, want %q", got.RunID, rc.RunID)
	}
}

func TestExtract_MissingContext_ReturnsFalse(t *testing.T) {
	_, ok := rotctx.Extract(context.Background())
	if ok {
		t.Fatal("expected Extract to return false for empty context")
	}
}

func TestMustExtract_Panics_WhenAbsent(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic but did not get one")
		}
	}()
	rotctx.MustExtract(context.Background())
}

func TestString_ContainsExpectedFields(t *testing.T) {
	rc := rotctx.New(true, "deployer")
	s := rc.String()

	for _, want := range []string{"run=", "dry-run", "operator=deployer", "started="} {
		if !strings.Contains(s, want) {
			t.Errorf("String() missing %q, got: %s", want, s)
		}
	}
}

func TestString_LiveMode(t *testing.T) {
	rc := rotctx.New(false, "human")
	if strings.Contains(rc.String(), "dry-run") {
		t.Errorf("expected 'live' mode, got: %s", rc.String())
	}
	if !strings.Contains(rc.String(), "live") {
		t.Errorf("expected 'live' in string, got: %s", rc.String())
	}
}
