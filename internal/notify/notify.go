package notify

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Level represents the severity of a notification.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

// Event holds the data for a single rotation notification.
type Event struct {
	Secret  string
	Backend string
	Level   Level
	Message string
	DryRun  bool
	Time    time.Time
}

// Notifier writes rotation events to an output.
type Notifier struct {
	out    io.Writer
	prefix string
}

// New returns a Notifier writing to out. If out is nil, os.Stdout is used.
func New(out io.Writer, prefix string) *Notifier {
	if out == nil {
		out = os.Stdout
	}
	return &Notifier{out: out, prefix: prefix}
}

// Notify formats and writes an Event to the configured writer.
func (n *Notifier) Notify(e Event) error {
	if e.Time.IsZero() {
		e.Time = time.Now().UTC()
	}

	dryTag := ""
	if e.DryRun {
		dryTag = " [dry-run]"
	}

	prefix := ""
	if n.prefix != "" {
		prefix = fmt.Sprintf("[%s] ", strings.TrimSpace(n.prefix))
	}

	line := fmt.Sprintf("%s%s %-5s %s/%s: %s%s\n",
		prefix,
		e.Time.Format(time.RFC3339),
		string(e.Level),
		e.Backend,
		e.Secret,
		e.Message,
		dryTag,
	)

	_, err := fmt.Fprint(n.out, line)
	return err
}

// NotifyAll sends multiple events, returning the first error encountered.
func (n *Notifier) NotifyAll(events []Event) error {
	for _, e := range events {
		if err := n.Notify(e); err != nil {
			return err
		}
	}
	return nil
}
