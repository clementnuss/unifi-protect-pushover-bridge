package types

import (
	"encoding/json"
	"testing"
)

func TestUnifiWebhookParsing(t *testing.T) {
	payload := `{
		"alarm": {
			"name": "Front Door Person Detected",
			"sources": [{
				"device": "F4E2C60E6104",
				"type": "include"
			}],
			"conditions": [
				{
					"condition": {
						"type": "is",
						"source": "person"
					}
				}
			],
			"triggers": [
				{
					"key": "person",
					"device": "F4E2C60E6104"
				}
			]
		},
		"timestamp": 1725883107267
	}`

	var webhook UnifiWebhook
	err := json.Unmarshal([]byte(payload), &webhook)
	if err != nil {
		t.Fatalf("failed to unmarshal webhook: %v", err)
	}

	if webhook.Alarm.Name != "Front Door Person Detected" {
		t.Errorf("expected alarm name %q, got %q", "Front Door Person Detected", webhook.Alarm.Name)
	}

	if len(webhook.Alarm.Triggers) != 1 {
		t.Fatalf("expected 1 trigger, got %d", len(webhook.Alarm.Triggers))
	}

	if webhook.Alarm.Triggers[0].Key != "person" {
		t.Errorf("expected trigger key %q, got %q", "person", webhook.Alarm.Triggers[0].Key)
	}

	if webhook.Timestamp != 1725883107267 {
		t.Errorf("expected timestamp %d, got %d", 1725883107267, webhook.Timestamp)
	}

	if len(webhook.Alarm.Sources) != 1 {
		t.Fatalf("expected 1 source, got %d", len(webhook.Alarm.Sources))
	}

	if webhook.Alarm.Sources[0].Device != "F4E2C60E6104" {
		t.Errorf("expected device %q, got %q", "F4E2C60E6104", webhook.Alarm.Sources[0].Device)
	}
}

func TestMinimalWebhookParsing(t *testing.T) {
	payload := `{
		"alarm": {
			"name": "Front Door Person",
			"triggers": [{"key": "person"}]
		},
		"timestamp": 1725883107267
	}`

	var webhook UnifiWebhook
	err := json.Unmarshal([]byte(payload), &webhook)
	if err != nil {
		t.Fatalf("failed to unmarshal webhook: %v", err)
	}

	if webhook.Alarm.Name != "Front Door Person" {
		t.Errorf("expected alarm name %q, got %q", "Front Door Person", webhook.Alarm.Name)
	}

	if len(webhook.Alarm.Triggers) != 1 {
		t.Fatalf("expected 1 trigger, got %d", len(webhook.Alarm.Triggers))
	}

	if webhook.Alarm.Triggers[0].Key != "person" {
		t.Errorf("expected trigger key %q, got %q", "person", webhook.Alarm.Triggers[0].Key)
	}
}
