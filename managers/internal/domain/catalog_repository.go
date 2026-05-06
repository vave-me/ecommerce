package domain

import (
	"context"
	"time"
)

type CatalogManager struct {
	ID           string
	Name         string
	Description  string
	UserID       string
	Type         ManagerType
	Capabilities []ManagerCapability
	Active       bool
	Temperature  float64
	MaxTokens    int
	SystemPrompt string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// CatalogRepository defines the interface for manager queries (Read Side)
type CatalogRepository interface {
	Add(ctx context.Context, id, name, description, userID string, managerType ManagerType, capabilities []ManagerCapability, active bool, temperature float64, maxTokens int, systemPrompt string) error
	Remove(ctx context.Context, id string) error
	Update(ctx context.Context, manager *CatalogManager) error
	FindAll(ctx context.Context, userID string) ([]*CatalogManager, error)
	Find(ctx context.Context, id string) (*CatalogManager, error)
	FindActiveByUser(ctx context.Context, userID string) ([]*CatalogManager, error)
	UpdateActiveStatus(ctx context.Context, id string, active bool, updatedAt time.Time) error
	UpdateConfiguration(ctx context.Context, id string, temperature float64, maxTokens int, systemPrompt string, updatedAt time.Time) error
	UpdateConfigurationWithCapabilities(ctx context.Context, id string, temperature float64, maxTokens int, systemPrompt string, capabilities []ManagerCapability, updatedAt time.Time) error
}
