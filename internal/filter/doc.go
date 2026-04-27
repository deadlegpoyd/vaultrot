// Package filter provides name-pattern and tag-based filtering for secrets
// during rotation. It supports wildcard glob patterns (using * and ?),
// exact name matching, tag key=value filtering, and exclusion patterns.
//
// Example usage:
//
//	f := filter.New(filter.Options{
//		Patterns:        []string{"prod/*"},
//		ExcludePatterns: []string{"prod/temp-*"},
//		Tags:            map[string]string{"env": "prod"},
//	})
//
//	if f.Match(secretName, secretTags) {
//		// rotate this secret
//	}
package filter
