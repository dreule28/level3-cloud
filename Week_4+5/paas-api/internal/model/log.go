package model

import "time"

const (
	LogTypeAudit   = "audit"
	LogTypeService = "service"
)

// LogEntry represents an audit or service log record.
type LogEntry struct {
	ID         string    `json:"id"`
	InstanceID string    `json:"instanceId"`
	Type       string    `json:"type"`   // "audit" | "service"
	Action     string    `json:"action"` // e.g. "instance.created", "status.changed"
	Message    string    `json:"message"`
	User       string    `json:"user,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}

type LogQuery struct {
	InstanceID string
	Type       string
	Limit      int
}
