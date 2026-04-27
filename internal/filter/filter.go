// Package filter provides secret filtering based on name patterns and tags.
package filter

import (
	"regexp"
	"strings"
)

// Options holds the filtering criteria for secrets.
type Options struct {
	// Patterns is a list of glob-style or regex patterns to match secret names.
	Patterns []string
	// Tags filters secrets that contain all specified key=value tags.
	Tags map[string]string
	// ExcludePatterns is a list of patterns for secrets to skip.
	ExcludePatterns []string
}

// Filter evaluates secrets against the configured options.
type Filter struct {
	opts Options
}

// New creates a new Filter with the given options.
func New(opts Options) *Filter {
	return &Filter{opts: opts}
}

// Match returns true if the secret name and tags satisfy the filter criteria.
func (f *Filter) Match(name string, tags map[string]string) bool {
	if len(f.opts.ExcludePatterns) > 0 {
		for _, p := range f.opts.ExcludePatterns {
			if matchPattern(p, name) {
				return false
			}
		}
	}

	if len(f.opts.Patterns) > 0 {
		matched := false
		for _, p := range f.opts.Patterns {
			if matchPattern(p, name) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	for k, v := range f.opts.Tags {
		if tags[k] != v {
			return false
		}
	}

	return true
}

// matchPattern supports simple wildcard (*) patterns and falls back to regex.
func matchPattern(pattern, name string) bool {
	if strings.ContainsAny(pattern, "*?") {
		regexStr := "^" + regexp.QuoteMeta(pattern) + "$"
		regexStr = strings.ReplaceAll(regexStr, `\*`, ".*")
		regexStr = strings.ReplaceAll(regexStr, `\?`, ".")
		re, err := regexp.Compile(regexStr)
		if err != nil {
			return false
		}
		return re.MatchString(name)
	}
	return strings.EqualFold(pattern, name)
}
