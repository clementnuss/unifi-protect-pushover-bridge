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

// triggerInfo holds French translation and gender
type triggerInfo struct {
	name     string
	feminine bool
}

// triggerTypeFR translates trigger types to French with gender
var triggerTypeFR = map[string]triggerInfo{
	"person":  {"personne", true},
	"vehicle": {"véhicule", false},
	"animal":  {"animal", false},
	"motion":  {"mouvement", false},
	"package": {"colis", false},
	"ring":    {"sonnette", true},
}

// SendAlert sends an emergency notification for a UniFi webhook event
func (n *Notifier) SendAlert(webhook *types.UnifiWebhook) error {
	triggerType := "inconnu"
	detected := "détecté"

	if len(webhook.Alarm.Triggers) > 0 {
		key := webhook.Alarm.Triggers[0].Key
		if info, ok := triggerTypeFR[key]; ok {
			triggerType = info.name
			if info.feminine {
				detected = "détectée"
			}
		} else {
			triggerType = key
		}
	}

	// Convert timestamp from milliseconds to time in Swiss timezone
	loc, _ := time.LoadLocation("Europe/Zurich")
	timestamp := time.UnixMilli(webhook.Timestamp).In(loc)
	timeStr := timestamp.Format("15:04")

	// Format the message in French
	messageText := fmt.Sprintf("🚨 %s - %s %s à %s",
		webhook.Alarm.Name,
		triggerType,
		detected,
		timeStr,
	)

	message := pushover.NewMessageWithTitle(messageText, "Alerte UniFi Protect")
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
