package schedule_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/mikelorant/vaultrot/internal/schedule"
)

func TestNew_ReturnsScheduler(t *testing.T) {
	s := schedule.New()
	if s == nil {
		t.Fatal("expected non-nil scheduler")
	}
}

func TestAdd_ValidExpression(t *testing.T) {
	s := schedule.New()
	err := s.Add("test-job", "@every 1m", func() {})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAdd_InvalidExpression(t *testing.T) {
	s := schedule.New()
	err := s.Add("bad-job", "not-a-cron", func() {})
	if err == nil {
		t.Fatal("expected error for invalid cron expression")
	}
}

func TestAdd_EmptyName(t *testing.T) {
	s := schedule.New()
	err := s.Add("", "@every 1m", func() {})
	if err == nil {
		t.Fatal("expected error for empty job name")
	}
}

func TestEntries_ReturnsRegistered(t *testing.T) {
	s := schedule.New()
	_ = s.Add("job-a", "@every 5m", func() {})
	_ = s.Add("job-b", "@every 10m", func() {})

	entries := s.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Name != "job-a" {
		t.Errorf("expected job-a, got %s", entries[0].Name)
	}
	if entries[1].CronExpr != "@every 10m" {
		t.Errorf("unexpected cron expr: %s", entries[1].CronExpr)
	}
}

func TestStart_Stop_JobFires(t *testing.T) {
	s := schedule.New()
	var count atomic.Int32

	err := s.Add("frequent", "* * * * * *", func() {
		count.Add(1)
	})
	if err != nil {
		t.Fatalf("add error: %v", err)
	}

	s.Start()
	time.Sleep(1100 * time.Millisecond)
	s.Stop()

	if count.Load() < 1 {
		t.Errorf("expected job to fire at least once, fired %d times", count.Load())
	}
}
