package domain

import (
	"context"
	"time"
)

// CatalogTask represents a task in the read model
type CatalogTask struct {
	ID           string
	ManagerID    string
	TaskType     string
	ScheduledAt  time.Time
	Payload      map[string]string
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ExecutedAt   *time.Time
	Result       string
	ErrorMessage string
}

// CatalogTaskRepository defines the interface for task read model operations
type CatalogTaskRepository interface {
	// Add adds a new task to the catalog
	Add(ctx context.Context, task *CatalogTask) error
	
	// Update updates an existing task in the catalog
	Update(ctx context.Context, taskID string, updates map[string]interface{}) error
	
	// Find retrieves a task by ID
	Find(ctx context.Context, taskID string) (*CatalogTask, error)
	
	// FindByManagerID retrieves tasks for a specific manager
	FindByManagerID(ctx context.Context, managerID string, filter TaskFilter) ([]*CatalogTask, error)
	
	// FindPendingTasks retrieves tasks that are due for execution
	FindPendingTasks(ctx context.Context, beforeTime time.Time, limit int) ([]*CatalogTask, error)
	
	// FindByStatus retrieves tasks by status
	FindByStatus(ctx context.Context, status string, limit int) ([]*CatalogTask, error)
	
	// CountByManagerID counts tasks for a specific manager
	CountByManagerID(ctx context.Context, managerID string, filter TaskFilter) (int, error)
	
	// Delete removes a task from the catalog
	Delete(ctx context.Context, taskID string) error
}

// TaskFilter provides filtering options for task queries
type TaskFilter struct {
	Status          *string
	ScheduledAfter  *time.Time
	ScheduledBefore *time.Time
	TaskTypes       []string
	Limit           int
	Offset          int
}