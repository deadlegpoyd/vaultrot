// Package version provides build-time version information and a
// structured reporter for displaying it in the CLI.
package version

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"text/tabwriter"
	"time"
)

// These variables are populated at build time via -ldflags.
var (
	Version   = "dev"
	Commit    = "none"
	Date      = "unknown"
	BuiltBy   = "unknown"
)

// Info holds structured version metadata.
type Info struct {
	Version   string
	Commit    string
	Date      string
	BuiltBy   string
	GoVersion string
	OS        string
	Arch      string
}

// Get returns the current build information.
func Get() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		Date:      Date,
		BuiltBy:   BuiltBy,
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}

// String returns a compact single-line representation of the version info,
// suitable for use in log output or --version flags.
func (i Info) String() string {
	return fmt.Sprintf("%s (commit=%s, built=%s, go=%s, %s/%s)",
		i.Version, i.Commit, formatDate(i.Date), i.GoVersion, i.OS, i.Arch)
}

// Print writes a formatted version table to the given writer.
// If w is nil it falls back to os.Stdout.
func Print(w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	info := Get()
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Version:\t%s\n", info.Version)
	fmt.Fprintf(tw, "Commit:\t%s\n", info.Commit)
	fmt.Fprintf(tw, "Built at:\t%s\n", formatDate(info.Date))
	fmt.Fprintf(tw, "Built by:\t%s\n", info.BuiltBy)
	fmt.Fprintf(tw, "Go version:\t%s\n", info.GoVersion)
	fmt.Fprintf(tw, "OS/Arch:\t%s/%s\n", info.OS, info.Arch)
	_ = tw.Flush()
}

// formatDate attempts to parse an RFC3339 date string and re-format it
// as a human-readable UTC timestamp; returns the raw string on failure.
func formatDate(raw string) string {
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return raw
	}
	return t.UTC().Format("2006-01-02 15:04:05 UTC")
}
