package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestWebhook_Send_Success(t *testing.T) {
	var received WebhookPayload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	wh := NewWebhook(WebhookConfig{URL: ts.URL})
	err := wh.Send(Event{
		Secret:  "db/pass",
		Backend: "vault",
		Level:   LevelInfo,
		Message: "rotated",
		Time:    fixedTime(),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Secret != "db/pass" {
		t.Errorf("expected secret db/pass, got %s", received.Secret)
	}
	if received.Backend != "vault" {
		t.Errorf("expected backend vault, got %s", received.Backend)
	}
}

func TestWebhook_Send_CustomHeader(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Token") != "secret" {
			t.Errorf("expected X-Token header")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	wh := NewWebhook(WebhookConfig{
		URL:     ts.URL,
		Headers: map[string]string{"X-Token": "secret"},
	})
	if err := wh.Send(Event{Secret: "s", Backend: "b", Level: LevelInfo, Message: "m", Time: fixedTime()}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWebhook_Send_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	wh := NewWebhook(WebhookConfig{URL: ts.URL})
	err := wh.Send(Event{Secret: "s", Backend: "b", Level: LevelError, Message: "fail", Time: fixedTime()})
	if err == nil {
		t.Error("expected error for 5xx response")
	}
}

func TestWebhook_DefaultTimeout(t *testing.T) {
	wh := NewWebhook(WebhookConfig{URL: "http://localhost"})
	if wh.client.Timeout != 5*time.Second {
		t.Errorf("expected 5s default timeout, got %v", wh.client.Timeout)
	}
}
