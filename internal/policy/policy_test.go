package policy_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultrot/internal/policy"
)

var fixedNow = time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)

func clock() func() time.Time { return func() time.Time { return fixedNow } }

func TestNew_ValidConfig(t *testing.T) {
	_, err := policy.New(policy.Config{MinAgeDays: 1, MaxAgeDays: 90}, clock())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_InvalidMinAge(t *testing.T) {
	_, err := policy.New(policy.Config{MinAgeDays: -1}, clock())
	if err == nil {
		t.Fatal("expected error for negative min_age_days")
	}
}

func TestNew_MinExceedsMax(t *testing.T) {
	_, err := policy.New(policy.Config{MinAgeDays: 30, MaxAgeDays: 10}, clock())
	if err == nil {
		t.Fatal("expected error when min_age_days > max_age_days")
	}
}

func TestCheck_TooYoung(t *testing.T) {
	p, _ := policy.New(policy.Config{MinAgeDays: 7}, clock())
	last := fixedNow.Add(-3 * 24 * time.Hour) // 3 days ago
	if err := p.Check("my-secret", last); err == nil {
		t.Fatal("expected violation for secret younger than min age")
	}
}

func TestCheck_OldEnough(t *testing.T) {
	p, _ := policy.New(policy.Config{MinAgeDays: 7}, clock())
	last := fixedNow.Add(-10 * 24 * time.Hour) // 10 days ago
	if err := p.Check("my-secret", last); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheck_ZeroLastRotated_SkipsMinAge(t *testing.T) {
	p, _ := policy.New(policy.Config{MinAgeDays: 30}, clock())
	if err := p.Check("new-secret", time.Time{}); err != nil {
		t.Fatalf("zero last-rotated should skip min-age check, got: %v", err)
	}
}

func TestCheck_ExceedsMaxAge(t *testing.T) {
	p, _ := policy.New(policy.Config{MaxAgeDays: 30}, clock())
	last := fixedNow.Add(-45 * 24 * time.Hour)
	if err := p.Check("old-secret", last); err == nil {
		t.Fatal("expected violation for secret exceeding max age")
	}
}

func TestCheck_WithinMaxAge(t *testing.T) {
	p, _ := policy.New(policy.Config{MaxAgeDays: 30}, clock())
	last := fixedNow.Add(-10 * 24 * time.Hour)
	if err := p.Check("ok-secret", last); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRequiresRotation_True(t *testing.T) {
	p, _ := policy.New(policy.Config{MaxAgeDays: 14}, clock())
	last := fixedNow.Add(-20 * 24 * time.Hour)
	if !p.RequiresRotation(last) {
		t.Fatal("expected RequiresRotation to return true")
	}
}

func TestRequiresRotation_False_NoMaxAge(t *testing.T) {
	p, _ := policy.New(policy.Config{}, clock())
	if p.RequiresRotation(time.Time{}) {
		t.Fatal("expected RequiresRotation false when no max_age_days set")
	}
}

func TestRequiresRotation_False_WithinMaxAge(t *testing.T) {
	p, _ := policy.New(policy.Config{MaxAgeDays: 30}, clock())
	last := fixedNow.Add(-15 * 24 * time.Hour) // 15 days ago, within 30-day max
	if p.RequiresRotation(last) {
		t.Fatal("expected RequiresRotation false when secret is within max age")
	}
}
