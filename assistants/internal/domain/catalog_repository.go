package domain

import (
	"context"
	"time"
)

type CatalogAssistant struct {
	ID           string
	Name         string
	Description  string
	UserID       string
	Type         AssistantType
	Capabilities []AssistantCapability
	Active       bool
	Temperature  float64
	MaxTokens    int
	SystemPrompt string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// CatalogRepository defines the interface for assistant queries (Read Side)
type CatalogRepository interface {
	Add(ctx context.Context, id, name, description, userID string, assistantType AssistantType, capabilities []AssistantCapability, active bool, temperature float64, maxTokens int, systemPrompt string) error
	Remove(ctx context.Context, id string) error
	Update(ctx context.Context, assistant *CatalogAssistant) error
	FindAll(ctx context.Context, userID string) ([]*CatalogAssistant, error)
	Find(ctx context.Context, id string) (*CatalogAssistant, error)
	FindActiveByUser(ctx context.Context, userID string) ([]*CatalogAssistant, error)
	UpdateActiveStatus(ctx context.Context, id string, active bool, updatedAt time.Time) error
	UpdateConfiguration(ctx context.Context, id string, temperature float64, maxTokens int, systemPrompt string, updatedAt time.Time) error
	UpdateConfigurationWithCapabilities(ctx context.Context, id string, temperature float64, maxTokens int, systemPrompt string, capabilities []AssistantCapability, updatedAt time.Time) error
}
