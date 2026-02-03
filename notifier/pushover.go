package notifier

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/clementnuss/unifi-protect-pushover-bridge/config"
	"github.com/clementnuss/unifi-protect-pushover-bridge/types"
	"github.com/gregdel/pushover"
)

// Notifier handles sending notifications to Pushover
type Notifier struct {
	app       *pushover.Pushover
	recipient *pushover.Recipient
	config    *config.Config
	logger    *slog.Logger
}

// New creates a new Pushover notifier
func New(cfg *config.Config, logger *slog.Logger) *Notifier {
	app := pushover.New(cfg.PushoverToken)
	recipient := pushover.NewRecipient(cfg.PushoverUserKey)

	return &Notifier{
		app:       app,
		recipient: recipient,
		config:    cfg,
		logger:    logger,
	}
}

// SendAlert sends an emergency notification for a UniFi webhook event
func (n *Notifier) SendAlert(webhook *types.UnifiWebhook) error {
	triggerType := "unknown"
	if len(webhook.Alarm.Triggers) > 0 {
		triggerType = webhook.Alarm.Triggers[0].Key
	}

	// Convert timestamp from milliseconds to time
	timestamp := time.UnixMilli(webhook.Timestamp)
	timeStr := timestamp.Format("03:04 PM")

	// Format the message
	messageText := fmt.Sprintf("🚨 %s - %s detected at %s",
		webhook.Alarm.Name,
		triggerType,
		timeStr,
	)

	message := pushover.NewMessageWithTitle(messageText, "UniFi Protect Alert")
	message.Priority = n.config.PushoverPriority
	message.Timestamp = timestamp.Unix()

	// Emergency priority requires retry and expire
	if n.config.PushoverPriority == pushover.PriorityEmergency {
		message.Retry = n.config.PushoverRetry
		message.Expire = n.config.PushoverExpire
	}

	n.logger.Info("sending pushover notification",
		"alarm", webhook.Alarm.Name,
		"trigger", triggerType,
		"priority", n.config.PushoverPriority,
	)

	response, err := n.app.SendMessage(message, n.recipient)
	if err != nil {
		n.logger.Error("failed to send pushover notification",
			"error", err,
			"alarm", webhook.Alarm.Name,
		)
		return fmt.Errorf("failed to send notification: %w", err)
	}

	n.logger.Info("pushover notification sent successfully",
		"alarm", webhook.Alarm.Name,
		"request_id", response.ID,
	)

	return nil
}
