package metrics_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultrot/internal/metrics"
)

func TestNew_InitialisesEmpty(t *testing.T) {
	c := metrics.New()
	if got := c.Get("anything"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestInc_IncrementsCounter(t *testing.T) {
	c := metrics.New()
	c.Inc("rotated")
	c.Inc("rotated")
	if got := c.Get("rotated"); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestAdd_NegativeDelta(t *testing.T) {
	c := metrics.New()
	c.Add("errors", 5)
	c.Add("errors", -2)
	if got := c.Get("errors"); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

func TestRecordDuration_Overwrite(t *testing.T) {
	c := metrics.New()
	c.RecordDuration("vault_write", 100*time.Millisecond)
	c.RecordDuration("vault_write", 200*time.Millisecond)
	if got := c.GetDuration("vault_write"); got != 200*time.Millisecond {
		t.Fatalf("expected 200ms, got %s", got)
	}
}

func TestGetDuration_Missing(t *testing.T) {
	c := metrics.New()
	if got := c.GetDuration("nonexistent"); got != 0 {
		t.Fatalf("expected zero duration, got %s", got)
	}
}

func TestElapsed_IsPositive(t *testing.T) {
	c := metrics.New()
	time.Sleep(2 * time.Millisecond)
	if c.Elapsed() <= 0 {
		t.Fatal("expected positive elapsed duration")
	}
}

func TestPrint_ContainsCounterAndDuration(t *testing.T) {
	c := metrics.New()
	c.Inc("skipped")
	c.RecordDuration("ssm_write", 50*time.Millisecond)

	var buf bytes.Buffer
	c.Print(&buf)
	out := buf.String()

	if !strings.Contains(out, "skipped") {
		t.Error("expected 'skipped' counter in output")
	}
	if !strings.Contains(out, "ssm_write") {
		t.Error("expected 'ssm_write' duration in output")
	}
	if !strings.Contains(out, "elapsed") {
		t.Error("expected 'elapsed' in output")
	}
}

func TestConcurrentInc_NoRace(t *testing.T) {
	c := metrics.New()
	done := make(chan struct{})
	for i := 0; i < 50; i++ {
		go func() {
			c.Inc("concurrent")
			done <- struct{}{}
		}()
	}
	for i := 0; i < 50; i++ {
		<-done
	}
	if got := c.Get("concurrent"); got != 50 {
		t.Fatalf("expected 50, got %d", got)
	}
}
