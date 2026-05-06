package domain

import (
	"context"
)

// TaskRepository defines the interface for task aggregate persistence
type TaskRepository interface {
	// Load retrieves a task aggregate by ID
	Load(ctx context.Context, taskID string) (*Task, error)
	
	// Save persists a task aggregate
	Save(ctx context.Context, task *Task) error
}