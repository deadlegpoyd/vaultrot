package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultrot/internal/config"
)

var (
	cfgFile string
	dryRun  bool
	loadedConfig *config.Config
)

// rootCmd is the base command for vaultrot.
var rootCmd = &cobra.Command{
	Use:   "vaultrot",
	Short: "Rotate secrets across Vault, AWS SSM, and Doppler",
	Long: `vaultrot is a CLI tool for rotating secrets across multiple backends.
Supported backends: vault, aws_ssm, doppler.

Use --dry-run to preview changes without applying them.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}
		// CLI flag overrides config file setting.
		if dryRun {
			cfg.DryRun = true
		}
		loadedConfig = cfg
		return nil
	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "vaultrot.yaml", "path to config file")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "preview rotations without applying changes")
}
