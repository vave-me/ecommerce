package domain

import (
	"context"
	"time"
)

// SchedulerRepository defines the interface for scheduler service operations
type SchedulerRepository interface {
	// ScheduleAction schedules a new action to be executed at a specific time
	ScheduleAction(ctx context.Context, action *ScheduledAction) error

	// CancelScheduledAction cancels a previously scheduled action
	CancelScheduledAction(ctx context.Context, actionID string) error

	// GetScheduledActions retrieves scheduled actions for a specific entity
	GetScheduledActions(ctx context.Context, entityID string, entityType string) ([]*ScheduledAction, error)

	// UpdateScheduledAction updates an existing scheduled action
	UpdateScheduledAction(ctx context.Context, actionID string, updates *ScheduledActionUpdate) error

	// Health checks if scheduler service is available
	Health(ctx context.Context) bool
}

// ScheduledAction represents an action to be executed at a scheduled time
type ScheduledAction struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	EntityID    string                 `json:"entity_id"`
	EntityType  string                 `json:"entity_type"`
	Action      string                 `json:"action"`
	Parameters  map[string]interface{} `json:"parameters"`
	ScheduledAt time.Time              `json:"scheduled_at"`
	Status      string                 `json:"status"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// ScheduledActionUpdate contains fields that can be updated
type ScheduledActionUpdate struct {
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	Action      *string                `json:"action,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	ScheduledAt *time.Time             `json:"scheduled_at,omitempty"`
	Status      *string                `json:"status,omitempty"`
}

// SchedulerActionStatus constants
const (
	SchedulerActionStatusPending   = "pending"
	SchedulerActionStatusExecuting = "executing"
	SchedulerActionStatusCompleted = "completed"
	SchedulerActionStatusFailed    = "failed"
	SchedulerActionStatusCancelled = "cancelled"
)