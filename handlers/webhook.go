package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/clementnuss/unifi-protect-pushover-bridge/notifier"
	"github.com/clementnuss/unifi-protect-pushover-bridge/types"
)

// WebhookHandler handles incoming UniFi Protect webhooks
type WebhookHandler struct {
	notifier *notifier.Notifier
	logger   *slog.Logger
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(n *notifier.Notifier, logger *slog.Logger) *WebhookHandler {
	return &WebhookHandler{
		notifier: n,
		logger:   logger,
	}
}

// ServeHTTP handles the webhook POST request
func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.Warn("invalid method", "method", r.Method)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var webhook types.UnifiWebhook
	if err := json.NewDecoder(r.Body).Decode(&webhook); err != nil {
		h.logger.Error("failed to decode webhook payload", "error", err)
		http.Error(w, "invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if webhook.Alarm.Name == "" {
		h.logger.Warn("missing alarm name in webhook")
		http.Error(w, "missing alarm.name", http.StatusBadRequest)
		return
	}

	if len(webhook.Alarm.Triggers) == 0 {
		h.logger.Warn("missing triggers in webhook")
		http.Error(w, "missing alarm.triggers", http.StatusBadRequest)
		return
	}

	h.logger.Info("received webhook",
		"alarm", webhook.Alarm.Name,
		"triggers", len(webhook.Alarm.Triggers),
		"timestamp", webhook.Timestamp,
	)

	if err := h.notifier.SendAlert(&webhook); err != nil {
		http.Error(w, "failed to send notification", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

// HealthHandler handles health check requests
type HealthHandler struct{}

// ServeHTTP handles the health check request
func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}
