package domain

import (
	"encoding/json"
	"time"
)

// SyncStatus represents the synchronization status for SAP entities
type SyncStatus struct {
	ID           string    `json:"id"`
	EntityType   string    `json:"entityType"`   // product, stock, price, order
	EntityID     string    `json:"entityId"`     // SAP entity ID
	LastSyncedAt time.Time `json:"lastSyncedAt"`
	Status       string    `json:"status"`       // success, failed, pending
	ErrorMessage string    `json:"errorMessage"`
	RetryCount   int       `json:"retryCount"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// SyncLog represents a log entry for synchronization activities
type SyncLog struct {
	ID            string          `json:"id"`
	EventID       string          `json:"eventId"`
	EventType     string          `json:"eventType"`
	Source        string          `json:"source"`
	Destination   string          `json:"destination"`
	Status        string          `json:"status"`
	Payload       json.RawMessage `json:"payload"`
	ResponseData  json.RawMessage `json:"responseData,omitempty"`
	ErrorMessage  string          `json:"errorMessage,omitempty"`
	ProcessedAt   time.Time       `json:"processedAt"`
	Duration      time.Duration   `json:"duration"`
}

// SyncConfiguration represents configuration for SAP synchronization
type SyncConfiguration struct {
	ID                string                 `json:"id"`
	EntityType        string                 `json:"entityType"`
	Enabled           bool                   `json:"enabled"`
	SyncInterval      time.Duration          `json:"syncInterval"`
	RetryPolicy       RetryPolicy            `json:"retryPolicy"`
	MappingRules      map[string]interface{} `json:"mappingRules"`
	LastExecutedAt    *time.Time             `json:"lastExecutedAt"`
	NextExecutionAt   *time.Time             `json:"nextExecutionAt"`
	CreatedAt         time.Time              `json:"createdAt"`
	UpdatedAt         time.Time              `json:"updatedAt"`
}

// RetryPolicy defines retry behavior for failed syncs
type RetryPolicy struct {
	MaxAttempts     int           `json:"maxAttempts"`
	InitialDelay    time.Duration `json:"initialDelay"`
	MaxDelay        time.Duration `json:"maxDelay"`
	BackoffMultiplier float64     `json:"backoffMultiplier"`
}

// WebhookEvent represents an incoming webhook event from SAP
type WebhookEvent struct {
	ID            string          `json:"id"`
	EventID       string          `json:"eventId"`
	EventType     string          `json:"eventType"`
	Source        string          `json:"source"`
	Signature     string          `json:"signature"`
	Payload       json.RawMessage `json:"payload"`
	ReceivedAt    time.Time       `json:"receivedAt"`
	ProcessedAt   *time.Time      `json:"processedAt,omitempty"`
	Status        string          `json:"status"` // received, processing, processed, failed
	ErrorMessage  string          `json:"errorMessage,omitempty"`
	RetryCount    int             `json:"retryCount"`
}