package snapshot

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// Report writes a formatted table of all captured snapshots to w.
func Report(w io.Writer, s *Store) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "BACKEND\tKEY\tCAPTURED AT")
	fmt.Fprintln(tw, "-------\t---\t------------")

	entries := s.All()
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Backend != entries[j].Backend {
			return entries[i].Backend < entries[j].Backend
		}
		return entries[i].Key < entries[j].Key
	})

	for _, e := range entries {
		fmt.Fprintf(tw, "%s\t%s\t%s\n",
			e.Backend,
			e.Key,
			e.CapturedAt.Format("2006-01-02 15:04:05 UTC"),
		)
	}
	tw.Flush()
}
