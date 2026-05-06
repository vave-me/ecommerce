package domain

import (
	"time"
)

// Manager Event Names
const (
	ManagerCreatedEvent              = "managers.ManagerCreated"
	ManagerActivatedEvent            = "managers.ManagerActivated"
	ManagerDeactivatedEvent          = "managers.ManagerDeactivated"
	ManagerConfigurationUpdatedEvent = "managers.ManagerConfigurationUpdated"
	ManagerRequestProcessedEvent     = "managers.ManagerRequestProcessed"
)

// ManagerCreated event
type ManagerCreated struct {
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	Type         ManagerType         `json:"type"`
	Model        string              `json:"model"`
	UserID       string              `json:"user_id"`
	Capabilities []ManagerCapability `json:"capabilities"`
	Temperature  float64             `json:"temperature"`
	MaxTokens    int                 `json:"max_tokens"`
	SystemPrompt string              `json:"system_prompt"`
	Active       bool                `json:"active"`
	CreatedAt    time.Time           `json:"created_at"`
}

// Key implements registry.Registerable
func (ManagerCreated) Key() string { return ManagerCreatedEvent }

// ManagerActivated event
type ManagerActivated struct {
	Active    bool      `json:"active"`
	Timestamp time.Time `json:"timestamp"`
}

// Key implements registry.Registerable
func (ManagerActivated) Key() string { return ManagerActivatedEvent }

// ManagerDeactivated event
type ManagerDeactivated struct {
	Active    bool      `json:"active"`
	Timestamp time.Time `json:"timestamp"`
}

// Key implements registry.Registerable
func (ManagerDeactivated) Key() string { return ManagerDeactivatedEvent }

// ManagerConfigurationUpdated event
type ManagerConfigurationUpdated struct {
	Temperature  float64             `json:"temperature"`
	MaxTokens    int                 `json:"max_tokens"`
	SystemPrompt string              `json:"system_prompt"`
	Capabilities []ManagerCapability `json:"capabilities,omitempty"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

// Key implements registry.Registerable
func (ManagerConfigurationUpdated) Key() string { return ManagerConfigurationUpdatedEvent }

// ManagerRequestProcessed event
type ManagerRequestProcessed struct {
	RequestID          string                 `json:"request_id"`
	UserID             string                 `json:"user_id"`
	Message            string                 `json:"message"`
	Context            map[string]interface{} `json:"context,omitempty"`
	RequestType        string                 `json:"request_type"`
	RequestTimestamp   time.Time              `json:"request_timestamp"`
	ResponseID         string                 `json:"response_id"`
	ResponseMessage    string                 `json:"response_message"`
	ResponseData       map[string]interface{} `json:"response_data,omitempty"`
	ResponseActions    []ManagerAction        `json:"response_actions,omitempty"`
	ResponseTimestamp  time.Time              `json:"response_timestamp"`
	ResponseStatus     string                 `json:"response_status"`
	ResponseConfidence float64                `json:"response_confidence"`
}

// Key implements registry.Registerable
func (ManagerRequestProcessed) Key() string { return ManagerRequestProcessedEvent }
