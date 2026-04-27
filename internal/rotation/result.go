package rotation

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// Status represents the outcome of a single secret rotation.
type Status string

const (
	StatusSuccess Status = "success"
	StatusSkipped Status = "skipped"
	StatusFailed  Status = "failed"
)

// Result holds the outcome of rotating one secret.
type Result struct {
	SecretName string
	Backend    string
	Status     Status
	DryRun     bool
	Error      error
	RotatedAt  time.Time
}

// Summary aggregates all rotation results from a single run.
type Summary struct {
	Results  []Result
	DryRun   bool
	Started  time.Time
	Finished time.Time
}

// NewSummary creates an empty Summary with the current start time.
func NewSummary(dryRun bool) *Summary {
	return &Summary{
		DryRun:  dryRun,
		Started: time.Now(),
	}
}

// Add appends a Result to the Summary.
func (s *Summary) Add(r Result) {
	s.Results = append(s.Results, r)
}

// Finish marks the summary as complete.
func (s *Summary) Finish() {
	s.Finished = time.Now()
}

// Counts returns the number of successes, skips, and failures.
func (s *Summary) Counts() (success, skipped, failed int) {
	for _, r := range s.Results {
		switch r.Status {
		case StatusSuccess:
			success++
		case StatusSkipped:
			skipped++
		case StatusFailed:
			failed++
		}
	}
	return
}

// Print writes a human-readable table of results to w.
func (s *Summary) Print(w io.Writer) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	defer tw.Flush()

	dryTag := ""
	if s.DryRun {
		dryTag = " [DRY-RUN]"
	}
	fmt.Fprintf(tw, "Rotation Summary%s\n", dryTag)
	fmt.Fprintln(tw, "SECRET\tBACKEND\tSTATUS\tERROR")
	for _, r := range s.Results {
		errMsg := ""
		if r.Error != nil {
			errMsg = r.Error.Error()
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", r.SecretName, r.Backend, r.Status, errMsg)
	}
	success, skipped, failed := s.Counts()
	fmt.Fprintf(tw, "\nTotal: %d success, %d skipped, %d failed (elapsed: %s)\n",
		success, skipped, failed, s.Finished.Sub(s.Started).Round(time.Millisecond))
}
