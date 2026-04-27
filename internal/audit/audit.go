package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Entry represents a single audit log record for a secret rotation event.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Backend   string    `json:"backend"`
	SecretKey string    `json:"secret_key"`
	Status    string    `json:"status"` // "rotated", "skipped", "failed"
	DryRun    bool      `json:"dry_run"`
	Message   string    `json:"message,omitempty"`
}

// Logger writes structured audit entries to a destination.
type Logger struct {
	writer  io.Writer
	dryRun  bool
}

// New creates a new audit Logger. If path is empty, stdout is used.
func New(path string, dryRun bool) (*Logger, error) {
	var w io.Writer = os.Stdout
	if path != "" {
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return nil, fmt.Errorf("audit: open log file: %w", err)
		}
		w = f
	}
	return &Logger{writer: w, dryRun: dryRun}, nil
}

// Record writes an audit entry for the given backend and secret key.
func (l *Logger) Record(backend, secretKey, status, message string) error {
	e := Entry{
		Timestamp: time.Now().UTC(),
		Backend:   backend,
		SecretKey: secretKey,
		Status:    status,
		DryRun:    l.dryRun,
		Message:   message,
	}
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	_, err = fmt.Fprintln(l.writer, string(data))
	return err
}
