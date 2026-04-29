// Package schedule provides cron-based scheduling support for secret rotation.
package schedule

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

// Entry represents a scheduled rotation job.
type Entry struct {
	Name     string
	CronExpr string
	NextRun  time.Time
}

// Scheduler wraps a cron runner and tracks registered jobs.
type Scheduler struct {
	c       *cron.Cron
	entries []Entry
}

// New creates a new Scheduler using UTC time.
func New() *Scheduler {
	return &Scheduler{
		c: cron.New(cron.WithLocation(time.UTC)),
	}
}

// Add registers a named cron job. fn is called on each trigger.
// Returns an error if the cron expression is invalid.
func (s *Scheduler) Add(name, expr string, fn func()) error {
	if name == "" {
		return fmt.Errorf("schedule: job name must not be empty")
	}
	_, err := cron.ParseStandard(expr)
	if err != nil {
		return fmt.Errorf("schedule: invalid cron expression %q for job %q: %w", expr, name, err)
	}
	_, err = s.c.AddFunc(expr, fn)
	if err != nil {
		return fmt.Errorf("schedule: failed to add job %q: %w", name, err)
	}
	s.entries = append(s.entries, Entry{Name: name, CronExpr: expr})
	return nil
}

// Start begins the scheduler in the background.
func (s *Scheduler) Start() {
	s.c.Start()
}

// Stop halts the scheduler, waiting for running jobs to finish.
func (s *Scheduler) Stop() {
	s.c.Stop()
}

// Entries returns a snapshot of registered schedule entries with next-run times.
func (s *Scheduler) Entries() []Entry {
	cronEntries := s.c.Entries()
	result := make([]Entry, len(s.entries))
	copy(result, s.entries)
	for i := range result {
		if i < len(cronEntries) {
			result[i].NextRun = cronEntries[i].Next
		}
	}
	return result
}
