// Package redact provides utilities for masking sensitive secret values
// before they are written to logs, audit trails, or notification payloads.
package redact

import "strings"

const defaultMask = "***REDACTED***"

// Option configures a Redactor.
type Option func(*Redactor)

// Redactor masks sensitive string values.
type Redactor struct {
	mask    string
	prefix  int // visible prefix characters
	suffix  int // visible suffix characters
}

// New returns a Redactor with optional configuration applied.
func New(opts ...Option) *Redactor {
	r := &Redactor{mask: defaultMask}
	for _, o := range opts {
		o(r)
	}
	return r
}

// WithMask overrides the replacement string.
func WithMask(m string) Option {
	return func(r *Redactor) { r.mask = m }
}

// WithPeek keeps n leading and n trailing characters visible.
// e.g. "supersecret" with peek 2 → "su***et"
func WithPeek(prefix, suffix int) Option {
	return func(r *Redactor) {
		r.prefix = prefix
		r.suffix = suffix
	}
}

// Mask replaces value with the configured mask, optionally preserving
// a short prefix/suffix for identification purposes.
func (r *Redactor) Mask(value string) string {
	if value == "" {
		return value
	}
	n := len(value)
	p := r.prefix
	s := r.suffix
	if p+s >= n {
		// Not enough characters — fully redact.
		return r.mask
	}
	if p == 0 && s == 0 {
		return r.mask
	}
	var b strings.Builder
	if p > 0 {
		b.WriteString(value[:p])
	}
	b.WriteString(r.mask)
	if s > 0 {
		b.WriteString(value[n-s:])
	}
	return b.String()
}

// MaskMap returns a copy of m with every value replaced by the mask.
func (r *Redactor) MaskMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = r.Mask(v)
	}
	return out
}
