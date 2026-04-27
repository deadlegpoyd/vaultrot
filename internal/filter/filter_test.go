package filter_test

import (
	"testing"

	"github.com/yourusername/vaultrot/internal/filter"
)

func TestMatch_NoOptions_AcceptsAll(t *testing.T) {
	f := filter.New(filter.Options{})
	if !f.Match("any-secret", nil) {
		t.Error("expected match with empty options")
	}
}

func TestMatch_PatternWildcard(t *testing.T) {
	f := filter.New(filter.Options{
		Patterns: []string{"prod/*"},
	})
	if !f.Match("prod/db-password", nil) {
		t.Error("expected prod/db-password to match prod/*")
	}
	if f.Match("staging/db-password", nil) {
		t.Error("expected staging/db-password not to match prod/*")
	}
}

func TestMatch_ExactPattern(t *testing.T) {
	f := filter.New(filter.Options{
		Patterns: []string{"my-secret"},
	})
	if !f.Match("my-secret", nil) {
		t.Error("expected exact match")
	}
	if f.Match("other-secret", nil) {
		t.Error("expected no match for different name")
	}
}

func TestMatch_ExcludePattern(t *testing.T) {
	f := filter.New(filter.Options{
		ExcludePatterns: []string{"temp-*"},
	})
	if f.Match("temp-token", nil) {
		t.Error("expected temp-token to be excluded")
	}
	if !f.Match("prod-token", nil) {
		t.Error("expected prod-token to pass exclusion filter")
	}
}

func TestMatch_TagFilter(t *testing.T) {
	f := filter.New(filter.Options{
		Tags: map[string]string{"env": "prod"},
	})
	if !f.Match("secret", map[string]string{"env": "prod", "team": "ops"}) {
		t.Error("expected match when tag is present")
	}
	if f.Match("secret", map[string]string{"env": "staging"}) {
		t.Error("expected no match when tag value differs")
	}
}

func TestMatch_TagFilterMissingKey(t *testing.T) {
	f := filter.New(filter.Options{
		Tags: map[string]string{"env": "prod"},
	})
	if f.Match("secret", map[string]string{"team": "ops"}) {
		t.Error("expected no match when required tag is absent")
	}
}

func TestMatch_PatternAndExcludeCombined(t *testing.T) {
	f := filter.New(filter.Options{
		Patterns:        []string{"prod/*"},
		ExcludePatterns: []string{"prod/temp-*"},
	})
	if !f.Match("prod/db-pass", nil) {
		t.Error("expected prod/db-pass to match")
	}
	if f.Match("prod/temp-key", nil) {
		t.Error("expected prod/temp-key to be excluded")
	}
}
