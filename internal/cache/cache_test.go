package cache

import (
	"testing"
	"time"
)

// fixedNow returns a function that always returns t, used to freeze time.
func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestSet_And_Get_ReturnsValue(t *testing.T) {
	c := New(5 * time.Minute)
	c.Set("mykey", "myvalue")

	got, ok := c.Get("mykey")
	if !ok {
		t.Fatal("expected key to be present")
	}
	if got != "myvalue" {
		t.Fatalf("expected 'myvalue', got %q", got)
	}
}

func TestGet_MissingKey_ReturnsFalse(t *testing.T) {
	c := New(time.Minute)
	_, ok := c.Get("nonexistent")
	if ok {
		t.Fatal("expected missing key to return false")
	}
}

func TestGet_ExpiredEntry_ReturnsFalse(t *testing.T) {
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	c := New(10 * time.Second)
	c.nowFunc = fixedNow(base)
	c.Set("secret", "value")

	// Advance time beyond TTL
	c.nowFunc = fixedNow(base.Add(11 * time.Second))
	_, ok := c.Get("secret")
	if ok {
		t.Fatal("expected expired entry to return false")
	}
}

func TestGet_ZeroTTL_NeverExpires(t *testing.T) {
	c := New(0)
	c.Set("k", "v")

	// Simulate a very large time advance
	c.nowFunc = fixedNow(time.Now().Add(100 * 24 * time.Hour))
	_, ok := c.Get("k")
	if !ok {
		t.Fatal("zero-TTL entry should never expire")
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	c := New(time.Minute)
	c.Set("del", "gone")
	c.Delete("del")

	_, ok := c.Get("del")
	if ok {
		t.Fatal("expected deleted key to be absent")
	}
}

func TestFlush_ClearsAll(t *testing.T) {
	c := New(time.Minute)
	c.Set("a", "1")
	c.Set("b", "2")
	c.Flush()

	if c.Len() != 0 {
		t.Fatalf("expected empty cache after flush, got %d entries", c.Len())
	}
}

func TestLen_ReturnsCorrectCount(t *testing.T) {
	c := New(time.Minute)
	if c.Len() != 0 {
		t.Fatal("new cache should be empty")
	}
	c.Set("x", "1")
	c.Set("y", "2")
	if c.Len() != 2 {
		t.Fatalf("expected 2, got %d", c.Len())
	}
}

func TestSet_Overwrite_UpdatesValue(t *testing.T) {
	c := New(time.Minute)
	c.Set("key", "old")
	c.Set("key", "new")

	got, ok := c.Get("key")
	if !ok || got != "new" {
		t.Fatalf("expected 'new', got %q (ok=%v)", got, ok)
	}
}
