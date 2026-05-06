package domain

import (
	"time"
)

const (
	TaskScheduledEvent          = "scheduler.TaskScheduled"
	TaskCancelledEvent          = "scheduler.TaskCancelled"
	TaskUpdatedEvent            = "scheduler.TaskUpdated"
	TaskExecutionStartedEvent   = "scheduler.TaskExecutionStarted"
	TaskExecutionCompletedEvent = "scheduler.TaskExecutionCompleted"
	TaskExecutionFailedEvent    = "scheduler.TaskExecutionFailed"
)

// TaskScheduled is raised when a new task is scheduled
type TaskScheduled struct {
	TaskID      string
	ManagerID   string
	TaskType    string
	ScheduledAt time.Time
	Payload     map[string]string
	CreatedAt   time.Time
}

// TaskCancelled is raised when a task is cancelled
type TaskCancelled struct {
	TaskID      string
	CancelledAt time.Time
}

// TaskUpdated is raised when a task is updated
type TaskUpdated struct {
	TaskID      string
	ScheduledAt *time.Time
	Payload     map[string]string
	UpdatedAt   time.Time
}

// TaskExecutionStarted is raised when task execution begins
type TaskExecutionStarted struct {
	TaskID    string
	StartedAt time.Time
}

// TaskExecutionCompleted is raised when task execution completes successfully
type TaskExecutionCompleted struct {
	TaskID      string
	CompletedAt time.Time
	Result      string
}

// TaskExecutionFailed is raised when task execution fails
type TaskExecutionFailed struct {
	TaskID       string
	FailedAt     time.Time
	ErrorMessage string
}

// Key implementations for event store
func (TaskScheduled) Key() string          { return TaskScheduledEvent }
func (TaskCancelled) Key() string          { return TaskCancelledEvent }
func (TaskUpdated) Key() string            { return TaskUpdatedEvent }
func (TaskExecutionStarted) Key() string   { return TaskExecutionStartedEvent }
func (TaskExecutionCompleted) Key() string { return TaskExecutionCompletedEvent }
func (TaskExecutionFailed) Key() string    { return TaskExecutionFailedEvent }
