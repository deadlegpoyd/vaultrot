package validate

import (
	"fmt"
	"strings"
)

// RuleConfig is the YAML-serialisable form of a Rule, matching the shape used
// in vaultrot.yaml under each secret's `validation` block.
type RuleConfig struct {
	Name    string `yaml:"name"`
	Pattern string `yaml:"pattern"`
	MinLen  int    `yaml:"min_length"`
	MaxLen  int    `yaml:"max_length"`
}

// FromConfig converts a slice of RuleConfig (from YAML) into Rule values
// suitable for constructing a Validator.
func FromConfig(cfgs []RuleConfig) ([]Rule, error) {
	rules := make([]Rule, 0, len(cfgs))
	for _, c := range cfgs {
		if strings.TrimSpace(c.Name) == "" {
			return nil, fmt.Errorf("validate: rule name must not be empty")
		}
		if c.MinLen < 0 {
			return nil, fmt.Errorf("validate: rule %q has negative min_length", c.Name)
		}
		if c.MaxLen > 0 && c.MinLen > c.MaxLen {
			return nil, fmt.Errorf("validate: rule %q min_length %d exceeds max_length %d", c.Name, c.MinLen, c.MaxLen)
		}
		rules = append(rules, Rule{
			Name:    c.Name,
			Pattern: c.Pattern,
			MinLen:  c.MinLen,
			MaxLen:  c.MaxLen,
		})
	}
	return rules, nil
}

// MustFromConfig is like FromConfig but panics on error; intended for tests.
func MustFromConfig(cfgs []RuleConfig) []Rule {
	rules, err := FromConfig(cfgs)
	if err != nil {
		panic(err)
	}
	return rules
}
