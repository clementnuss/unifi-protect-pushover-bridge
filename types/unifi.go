package types

// UnifiWebhook represents the webhook payload from UniFi Protect Alarm Manager
type UnifiWebhook struct {
	Alarm     Alarm `json:"alarm"`
	Timestamp int64 `json:"timestamp"`
}

// Alarm represents the alarm configuration
type Alarm struct {
	Name       string      `json:"name"`
	Sources    []Source    `json:"sources,omitempty"`
	Conditions []Condition `json:"conditions,omitempty"`
	Triggers   []Trigger   `json:"triggers"`
}

// Source represents a device source
type Source struct {
	Device string `json:"device"`
	Type   string `json:"type"`
}

// Condition represents an alarm condition
type Condition struct {
	Condition ConditionDetail `json:"condition"`
}

// ConditionDetail represents the condition details
type ConditionDetail struct {
	Type   string `json:"type"`
	Source string `json:"source"`
}

// Trigger represents what triggered the alarm
type Trigger struct {
	Key    string `json:"key"`
	Device string `json:"device,omitempty"`
}
