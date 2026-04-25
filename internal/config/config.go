package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Backend represents a supported secret backend type.
type Backend string

const (
	BackendVault  Backend = "vault"
	BackendAWSSSM Backend = "aws_ssm"
	BackendDoppler Backend = "doppler"
)

// SecretEntry defines a single secret to rotate.
type SecretEntry struct {
	Name    string  `yaml:"name"`
	Path    string  `yaml:"path"`
	Backend Backend `yaml:"backend"`
}

// Config is the top-level configuration for vaultrot.
type Config struct {
	DryRun  bool          `yaml:"dry_run"`
	Secrets []SecretEntry `yaml:"secrets"`
	Vault   VaultConfig   `yaml:"vault"`
	AWSSSM  AWSSMMConfig  `yaml:"aws_ssm"`
	Doppler DopplerConfig `yaml:"doppler"`
}

// VaultConfig holds Vault-specific connection settings.
type VaultConfig struct {
	Address string `yaml:"address"`
	Token   string `yaml:"token"`
}

// AWSSMMConfig holds AWS SSM-specific settings.
type AWSSMMConfig struct {
	Region string `yaml:"region"`
}

// DopplerConfig holds Doppler-specific settings.
type DopplerConfig struct {
	Token   string `yaml:"token"`
	Project string `yaml:"project"`
	Config  string `yaml:"config"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// Validate checks that required fields are present and values are valid.
func (c *Config) Validate() error {
	for i, s := range c.Secrets {
		if s.Name == "" {
			return fmt.Errorf("secret[%d]: name is required", i)
		}
		if s.Path == "" {
			return fmt.Errorf("secret[%d] %q: path is required", i, s.Name)
		}
		switch s.Backend {
		case BackendVault, BackendAWSSSM, BackendDoppler:
		default:
			return fmt.Errorf("secret[%d] %q: unsupported backend %q", i, s.Name, s.Backend)
		}
	}
	return nil
}
