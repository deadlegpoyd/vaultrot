// Package masking provides field-level secret masking rules that can be
// applied before secrets are logged, printed, or sent to notification hooks.
package masking

import (
	"regexp"
	"strings"
)

// Rule describes how a single field should be masked.
type Rule struct {
	// Field is a case-insensitive field name or glob pattern (e.g. "*password*").
	Field string
	// Mask is the replacement string. Defaults to "***".
	Mask string
}

// Masker applies a set of Rules to key/value pairs.
type Masker struct {
	rules   []compiled
	default_ string
}

type compiled struct {
	re   *regexp.Regexp
	mask string
}

// New creates a Masker from the provided Rules.
// If no rules are given the Masker is a no-op.
func New(rules []Rule, defaultMask string) (*Masker, error) {
	if defaultMask == "" {
		defaultMask = "***"
	}
	m := &Masker{default_: defaultMask}
	for _, r := range rules {
		pattern := "(?i)^" + strings.ReplaceAll(regexp.QuoteMeta(r.Field), `\*`, `.*`) + "$"
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		mask := r.Mask
		if mask == "" {
			mask = defaultMask
		}
		m.rules = append(m.rules, compiled{re: re, mask: mask})
	}
	return m, nil
}

// Apply returns the masked value for the given field name.
// If no rule matches, the original value is returned unchanged.
func (m *Masker) Apply(field, value string) string {
	for _, r := range m.rules {
		if r.re.MatchString(field) {
			return r.mask
		}
	}
	return value
}

// ApplyMap masks all matching keys in the provided map, returning a new map.
func (m *Masker) ApplyMap(fields map[string]string) map[string]string {
	out := make(map[string]string, len(fields))
	for k, v := range fields {
		out[k] = m.Apply(k, v)
	}
	return out
}
