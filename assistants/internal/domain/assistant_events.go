package domain

import (
	"time"
)

// Assistant Event Names
const (
	AssistantCreatedEvent              = "assistants.AssistantCreated"
	AssistantActivatedEvent            = "assistants.AssistantActivated"
	AssistantDeactivatedEvent          = "assistants.AssistantDeactivated"
	AssistantConfigurationUpdatedEvent = "assistants.AssistantConfigurationUpdated"
	AssistantRequestProcessedEvent     = "assistants.AssistantRequestProcessed"
)

// AssistantCreated event
type AssistantCreated struct {
	Name         string                `json:"name"`
	Description  string                `json:"description"`
	Type         AssistantType         `json:"type"`
	Model        string                `json:"model"`
	UserID       string                `json:"user_id"`
	Capabilities []AssistantCapability `json:"capabilities"`
	Temperature  float64               `json:"temperature"`
	MaxTokens    int                   `json:"max_tokens"`
	SystemPrompt string                `json:"system_prompt"`
	Active       bool                  `json:"active"`
	CreatedAt    time.Time             `json:"created_at"`
}

// Key implements registry.Registerable
func (AssistantCreated) Key() string { return AssistantCreatedEvent }

// AssistantActivated event
type AssistantActivated struct {
	Active    bool      `json:"active"`
	Timestamp time.Time `json:"timestamp"`
}

// Key implements registry.Registerable
func (AssistantActivated) Key() string { return AssistantActivatedEvent }

// AssistantDeactivated event
type AssistantDeactivated struct {
	Active    bool      `json:"active"`
	Timestamp time.Time `json:"timestamp"`
}

// Key implements registry.Registerable
func (AssistantDeactivated) Key() string { return AssistantDeactivatedEvent }

// AssistantConfigurationUpdated event
type AssistantConfigurationUpdated struct {
	Temperature  float64               `json:"temperature"`
	MaxTokens    int                   `json:"max_tokens"`
	SystemPrompt string                `json:"system_prompt"`
	Capabilities []AssistantCapability `json:"capabilities,omitempty"`
	UpdatedAt    time.Time             `json:"updated_at"`
}

// Key implements registry.Registerable
func (AssistantConfigurationUpdated) Key() string { return AssistantConfigurationUpdatedEvent }

// AssistantRequestProcessed event
type AssistantRequestProcessed struct {
	RequestID          string                 `json:"request_id"`
	UserID             string                 `json:"user_id"`
	Message            string                 `json:"message"`
	Context            map[string]interface{} `json:"context,omitempty"`
	RequestType        string                 `json:"request_type"`
	RequestTimestamp   time.Time              `json:"request_timestamp"`
	ResponseID         string                 `json:"response_id"`
	ResponseMessage    string                 `json:"response_message"`
	ResponseData       map[string]interface{} `json:"response_data,omitempty"`
	ResponseActions    []AssistantAction      `json:"response_actions,omitempty"`
	ResponseTimestamp  time.Time              `json:"response_timestamp"`
	ResponseStatus     string                 `json:"response_status"`
	ResponseConfidence float64                `json:"response_confidence"`
}

// Key implements registry.Registerable
func (AssistantRequestProcessed) Key() string { return AssistantRequestProcessedEvent }
