package validate_test

import (
	"testing"

	"github.com/yourusername/vaultrot/internal/validate"
)

func baseRules() []validate.Rule {
	return []validate.Rule{
		{Name: "min-length", MinLen: 8},
		{Name: "max-length", MaxLen: 64},
		{Name: "alphanumeric", Pattern: `^[A-Za-z0-9]+$`},
	}
}

func TestValidate_PassesValidValue(t *testing.T) {
	v := validate.New(baseRules())
	res := v.Validate("mySecret", "SecurePass123")
	if !res.Passed {
		t.Fatalf("expected pass, got errors: %v", res.Errors)
	}
	if len(res.Errors) != 0 {
		t.Errorf("expected no errors, got %d", len(res.Errors))
	}
}

func TestValidate_FailsMinLength(t *testing.T) {
	v := validate.New(baseRules())
	res := v.Validate("short", "abc")
	if res.Passed {
		t.Fatal("expected failure for short value")
	}
	if len(res.Errors) == 0 {
		t.Error("expected at least one error")
	}
}

func TestValidate_FailsMaxLength(t *testing.T) {
	v := validate.New(baseRules())
	long := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz123"
	res := v.Validate("longSecret", long)
	if res.Passed {
		t.Fatal("expected failure for long value")
	}
}

func TestValidate_FailsPattern(t *testing.T) {
	v := validate.New(baseRules())
	res := v.Validate("withSymbols", "Pass!@#word1")
	if res.Passed {
		t.Fatal("expected failure for non-alphanumeric value")
	}
}

func TestValidate_NoRules_AlwaysPasses(t *testing.T) {
	v := validate.New(nil)
	res := v.Validate("anything", "")
	if !res.Passed {
		t.Errorf("expected pass with no rules, got errors: %v", res.Errors)
	}
}

func TestValidateAll_ReturnsAllResults(t *testing.T) {
	v := validate.New(baseRules())
	secrets := map[string]string{
		"good": "GoodPass99",
		"bad":  "x",
	}
	results := v.ValidateAll(secrets)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	passed := 0
	for _, r := range results {
		if r.Passed {
			passed++
		}
	}
	if passed != 1 {
		t.Errorf("expected 1 passing result, got %d", passed)
	}
}
