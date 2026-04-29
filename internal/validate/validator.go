// Package validate provides secret value validation before and after rotation.
package validate

import (
	"fmt"
	"regexp"
)

// Rule defines a single validation rule applied to a secret value.
type Rule struct {
	Name    string
	Pattern string
	MinLen  int
	MaxLen  int
}

// Result holds the outcome of validating a single secret.
type Result struct {
	SecretName string
	Passed     bool
	Errors     []string
}

// Validator applies a set of rules to secret values.
type Validator struct {
	rules []Rule
}

// New returns a Validator configured with the given rules.
func New(rules []Rule) *Validator {
	return &Validator{rules: rules}
}

// Validate checks value against all configured rules and returns a Result.
func (v *Validator) Validate(secretName, value string) Result {
	res := Result{SecretName: secretName, Passed: true}

	for _, r := range v.rules {
		if r.MinLen > 0 && len(value) < r.MinLen {
			res.Passed = false
			res.Errors = append(res.Errors, fmt.Sprintf("%s: value length %d below minimum %d", r.Name, len(value), r.MinLen))
		}
		if r.MaxLen > 0 && len(value) > r.MaxLen {
			res.Passed = false
			res.Errors = append(res.Errors, fmt.Sprintf("%s: value length %d exceeds maximum %d", r.Name, len(value), r.MaxLen))
		}
		if r.Pattern != "" {
			matched, err := regexp.MatchString(r.Pattern, value)
			if err != nil {
				res.Passed = false
				res.Errors = append(res.Errors, fmt.Sprintf("%s: invalid pattern %q: %v", r.Name, r.Pattern, err))
				continue
			}
			if !matched {
				res.Passed = false
				res.Errors = append(res.Errors, fmt.Sprintf("%s: value does not match pattern %q", r.Name, r.Pattern))
			}
		}
	}

	return res
}

// ValidateAll validates multiple secrets and returns all results.
func (v *Validator) ValidateAll(secrets map[string]string) []Result {
	results := make([]Result, 0, len(secrets))
	for name, val := range secrets {
		results = append(results, v.Validate(name, val))
	}
	return results
}
