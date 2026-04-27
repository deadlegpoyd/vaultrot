package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookConfig holds configuration for an HTTP webhook notifier.
type WebhookConfig struct {
	URL     string
	Headers map[string]string
	Timeout time.Duration
}

// WebhookPayload is the JSON body sent to the webhook endpoint.
type WebhookPayload struct {
	Secret  string `json:"secret"`
	Backend string `json:"backend"`
	Level   string `json:"level"`
	Message string `json:"message"`
	DryRun  bool   `json:"dry_run"`
	Time    string `json:"time"`
}

// WebhookNotifier sends rotation events to an HTTP endpoint.
type WebhookNotifier struct {
	cfg    WebhookConfig
	client *http.Client
}

// NewWebhook creates a WebhookNotifier. A zero Timeout defaults to 5s.
func NewWebhook(cfg WebhookConfig) *WebhookNotifier {
	if cfg.Timeout == 0 {
		cfg.Timeout = 5 * time.Second
	}
	return &WebhookNotifier{
		cfg:    cfg,
		client: &http.Client{Timeout: cfg.Timeout},
	}
}

// Send marshals the Event and POSTs it to the configured URL.
func (w *WebhookNotifier) Send(e Event) error {
	if e.Time.IsZero() {
		e.Time = time.Now().UTC()
	}

	payload := WebhookPayload{
		Secret:  e.Secret,
		Backend: e.Backend,
		Level:   string(e.Level),
		Message: e.Message,
		DryRun:  e.DryRun,
		Time:    e.Time.Format(time.RFC3339),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("notify/webhook: marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, w.cfg.URL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("notify/webhook: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range w.cfg.Headers {
		req.Header.Set(k, v)
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("notify/webhook: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("notify/webhook: unexpected status %d", resp.StatusCode)
	}
	return nil
}
