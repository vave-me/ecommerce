package domain

import (
	"time"

	"middleman/internal/es"
)

// TaskSnapshot is a snapshot of the Task aggregate
type TaskSnapshot struct {
	es.Snapshot
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

// SnapshotName implements es.Snapshot
func (s TaskSnapshot) SnapshotName() string {
	return "scheduler.TaskSnapshot"
}