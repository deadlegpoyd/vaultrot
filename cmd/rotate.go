package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vaultrot/vaultrot/internal/audit"
	"github.com/vaultrot/vaultrot/internal/backend"
	"github.com/vaultrot/vaultrot/internal/config"
	"github.com/vaultrot/vaultrot/internal/filter"
	"github.com/vaultrot/vaultrot/internal/notify"
	"github.com/vaultrot/vaultrot/internal/rotation"
)

var (
	cfgFile    string
	dryRun     bool
	tagFilter  string
	nameFilter string
	auditFile  string
)

// rotateCmd is the primary subcommand that triggers secret rotation.
var rotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "Rotate secrets across configured backends",
	Long: `Rotate secrets defined in the vaultrot configuration file.

Supports dry-run mode to preview changes without writing to any backend.
Use --tag and --name flags to selectively rotate a subset of secrets.`,
	RunE: runRotate,
}

func init() {
	rotateCmd.Flags().StringVarP(&cfgFile, "config", "c", "configs/vaultrot.yaml", "Path to the vaultrot config file")
	rotateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview rotations without writing to backends")
	rotateCmd.Flags().StringVar(&tagFilter, "tag", "", "Only rotate secrets matching this tag")
	rotateCmd.Flags().StringVar(&nameFilter, "name", "", "Only rotate secrets matching this name pattern")
	rotateCmd.Flags().StringVar(&auditFile, "audit-log", "", "Path to write audit log (JSON lines)")

	rootCmd.AddCommand(rotateCmd)
}

func runRotate(cmd *cobra.Command, args []string) error {
	// Load and validate configuration.
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Build filter from CLI flags.
	filterOpts := filter.Options{}
	if nameFilter != "" {
		filterOpts.Patterns = []string{nameFilter}
	}
	if tagFilter != "" {
		filterOpts.Tags = []string{tagFilter}
	}
	f := filter.New(filterOpts)

	// Set up audit logger (optional).
	var auditor *audit.Logger
	if auditFile != "" {
		f2, err := os.OpenFile(auditFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
		if err != nil {
			return fmt.Errorf("opening audit log %q: %w", auditFile, err)
		}
		defer f2.Close()
		auditor = audit.New(f2)
	}

	// Build notifier from config webhooks.
	notifier := notify.New(cfg.Notify)

	// Construct the rotator and run.
	rot := rotation.New(rotation.Config{
		Secrets:  cfg.Secrets,
		DryRun:   dryRun,
		Filter:   f,
		Auditor:  auditor,
		Notifier: notifier,
		BackendFactory: func(kind string, opts map[string]string) (backend.Backend, error) {
			return backend.New(kind, opts)
		},
	})

	summary, err := rot.Run(cmd.Context())
	if err != nil {
		return fmt.Errorf("rotation failed: %w", err)
	}

	summary.Print(os.Stdout)

	if summary.HasErrors() {
		// Exit with a non-zero code so CI pipelines detect partial failures.
		os.Exit(1)
	}

	return nil
}
