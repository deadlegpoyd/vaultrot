package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultrot/internal/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print build and runtime version information",
	Long: `Display the current vaultrot version alongside commit hash,
build timestamp, Go runtime version, and target platform.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		version.Print(os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
