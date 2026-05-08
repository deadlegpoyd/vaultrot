package changelog

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
	"time"
)

// Report writes a formatted table of changelog entries to w.
// Entries are sorted by OccurredAt ascending.
func Report(c *Changelog, w io.Writer) error {
	entries := c.Entries()
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].OccurredAt.Before(entries[j].OccurredAt)
	})

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIME\tBACKEND\tKEY\tACTOR\tACTION\tDRY-RUN\tNOTE")
	fmt.Fprintln(tw, "----\t-------\t---\t-----\t------\t-------\t----")
	for _, e := range entries {
		dry := "no"
		if e.DryRun {
			dry = "yes"
		}
		note := e.Note
		if note == "" {
			note = "-"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			e.OccurredAt.Format(time.RFC3339),
			e.Backend, e.Key, e.Actor, e.Action, dry, note)
	}
	return tw.Flush()
}
