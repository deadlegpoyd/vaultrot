// Package validate provides pre- and post-rotation secret validation for vaultrot.
//
// A Validator is constructed with a slice of Rule values, each specifying
// optional minimum/maximum length constraints and an optional regex pattern.
// Rules are evaluated independently; all failures are collected into the
// returned Result so callers receive a complete picture of why a value failed.
//
// Example usage:
//
//	rules := []validate.Rule{
//		{Name: "length", MinLen: 16, MaxLen: 128},
//		{Name: "complexity", Pattern: `[A-Z].*[0-9]`},
//	}
//	v := validate.New(rules)
//	result := v.Validate("db-password", newValue)
//	if !result.Passed {
//		log.Printf("validation failed: %v", result.Errors)
//	}
package validate
