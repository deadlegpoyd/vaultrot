package lease

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// Report writes a human-readable summary of all tracked leases to w.
func Report(w io.Writer, t *Tracker, now time.Time, renewalThreshold time.Duration) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(tw, "SECRET\tBACKEND\tEXPIRES\tSTATUS")

	t.mu.RLock()
	entries := make([]Entry, 0, len(t.entries))
	for _, e := range t.entries {
		entries = append(entries, e)
	}
	t.mu.RUnlock()

	for _, e := range entries {
		status := leaseStatus(e, now, renewalThreshold)
		_, _ = fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
			e.SecretName,
			e.Backend,
			e.ExpiresAt.UTC().Format(time.RFC3339),
			status,
		)
	}
	return tw.Flush()
}

func leaseStatus(e Entry, now time.Time, threshold time.Duration) string {
	switch {
	case e.IsExpired(now):
		return "EXPIRED"
	case e.DueForRenewal(now, threshold):
		return "RENEW_SOON"
	default:
		return "OK"
	}
}
