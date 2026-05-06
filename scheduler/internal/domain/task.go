package domain

import (
	"time"

	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

const TaskAggregate = "scheduler.Task"

var (
	ErrTaskAlreadyExists    = errors.Wrap(errors.ErrBadRequest, "task already exists")
	ErrTaskNotFound         = errors.Wrap(errors.ErrNotFound, "task not found")
	ErrInvalidTaskType      = errors.Wrap(errors.ErrBadRequest, "task type cannot be empty")
	ErrInvalidManagerID     = errors.Wrap(errors.ErrBadRequest, "manager ID cannot be empty")
	ErrInvalidScheduledTime = errors.Wrap(errors.ErrBadRequest, "scheduled time must be in the future")
	ErrTaskAlreadyExecuted  = errors.Wrap(errors.ErrBadRequest, "task has already been executed")
	ErrTaskCancelled        = errors.Wrap(errors.ErrBadRequest, "task has been cancelled")
)

// TaskStatus represents the status of a scheduled task
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
	TaskStatusCancelled TaskStatus = "cancelled"
)

// Task represents a scheduled task for managers
type Task struct {
	es.Aggregate
	ManagerID    string
	TaskType     string
	ScheduledAt  time.Time
	Payload      map[string]string
	Status       TaskStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ExecutedAt   *time.Time
	Result       string
	ErrorMessage string
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Task)(nil)

// NewTask creates a new task aggregate
func NewTask(id string) *Task {
	return &Task{
		Aggregate: es.NewAggregate(id, TaskAggregate),
		Status:    TaskStatusPending,
	}
}

// CreateTask creates a new task with initial schedule
func CreateTask(id, managerID, taskType string, scheduledAt time.Time, payload map[string]string) (*Task, error) {
	task := NewTask(id)
	
	// Use the Schedule method to create the task
	_, err := task.Schedule(managerID, taskType, scheduledAt, payload)
	if err != nil {
		return nil, err
	}
	
	return task, nil
}

// Schedule creates a new scheduled task
func (t *Task) Schedule(managerID, taskType string, scheduledAt time.Time, payload map[string]string) (ddd.Event, error) {
	if t.ManagerID != "" {
		return nil, ErrTaskAlreadyExists
	}

	if managerID == "" {
		return nil, ErrInvalidManagerID
	}

	if taskType == "" {
		return nil, ErrInvalidTaskType
	}

	if scheduledAt.Before(time.Now()) {
		return nil, ErrInvalidScheduledTime
	}

	t.AddEvent(TaskScheduledEvent, &TaskScheduled{
		TaskID:      t.ID(),
		ManagerID:   managerID,
		TaskType:    taskType,
		ScheduledAt: scheduledAt,
		Payload:     payload,
		CreatedAt:   time.Now(),
	})

	return ddd.NewEvent(TaskScheduledEvent, t), nil
}

// Cancel cancels a scheduled task
func (t *Task) Cancel() (ddd.Event, error) {
	if t.Status == TaskStatusCancelled {
		return nil, nil // Already cancelled
	}

	if t.Status == TaskStatusCompleted || t.Status == TaskStatusFailed {
		return nil, ErrTaskAlreadyExecuted
	}

	t.AddEvent(TaskCancelledEvent, &TaskCancelled{
		TaskID:      t.ID(),
		CancelledAt: time.Now(),
	})

	return ddd.NewEvent(TaskCancelledEvent, t), nil
}

// Update updates the scheduled time and/or payload
func (t *Task) Update(scheduledAt *time.Time, payload map[string]string) (ddd.Event, error) {
	if t.Status != TaskStatusPending {
		return nil, ErrTaskAlreadyExecuted
	}

	if scheduledAt != nil && scheduledAt.Before(time.Now()) {
		return nil, ErrInvalidScheduledTime
	}

	event := &TaskUpdated{
		TaskID:    t.ID(),
		UpdatedAt: time.Now(),
	}

	if scheduledAt != nil {
		event.ScheduledAt = scheduledAt
	}

	if payload != nil {
		event.Payload = payload
	}

	t.AddEvent(TaskUpdatedEvent, event)

	return ddd.NewEvent(TaskUpdatedEvent, t), nil
}

// StartExecution marks the task as running
func (t *Task) StartExecution() (ddd.Event, error) {
	if t.Status != TaskStatusPending {
		return nil, ErrTaskAlreadyExecuted
	}

	t.AddEvent(TaskExecutionStartedEvent, &TaskExecutionStarted{
		TaskID:    t.ID(),
		StartedAt: time.Now(),
	})

	return ddd.NewEvent(TaskExecutionStartedEvent, t), nil
}

// Complete marks the task as completed with a result
func (t *Task) Complete(result string) (ddd.Event, error) {
	if t.Status != TaskStatusRunning {
		return nil, errors.Wrap(errors.ErrBadRequest, "task is not running")
	}

	t.AddEvent(TaskExecutionCompletedEvent, &TaskExecutionCompleted{
		TaskID:      t.ID(),
		CompletedAt: time.Now(),
		Result:      result,
	})

	return ddd.NewEvent(TaskExecutionCompletedEvent, t), nil
}

// Fail marks the task as failed with an error
func (t *Task) Fail(errorMessage string) (ddd.Event, error) {
	if t.Status != TaskStatusRunning {
		return nil, errors.Wrap(errors.ErrBadRequest, "task is not running")
	}

	t.AddEvent(TaskExecutionFailedEvent, &TaskExecutionFailed{
		TaskID:       t.ID(),
		FailedAt:     time.Now(),
		ErrorMessage: errorMessage,
	})

	return ddd.NewEvent(TaskExecutionFailedEvent, t), nil
}

// ApplyEvent implements es.EventApplier
func (t *Task) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *TaskScheduled:
		t.ManagerID = payload.ManagerID
		t.TaskType = payload.TaskType
		t.ScheduledAt = payload.ScheduledAt
		t.Payload = payload.Payload
		t.Status = TaskStatusPending
		t.CreatedAt = payload.CreatedAt
		t.UpdatedAt = payload.CreatedAt

	case *TaskCancelled:
		t.Status = TaskStatusCancelled
		t.UpdatedAt = payload.CancelledAt

	case *TaskUpdated:
		if payload.ScheduledAt != nil {
			t.ScheduledAt = *payload.ScheduledAt
		}
		if payload.Payload != nil {
			t.Payload = payload.Payload
		}
		t.UpdatedAt = payload.UpdatedAt

	case *TaskExecutionStarted:
		t.Status = TaskStatusRunning
		now := payload.StartedAt
		t.ExecutedAt = &now
		t.UpdatedAt = payload.StartedAt

	case *TaskExecutionCompleted:
		t.Status = TaskStatusCompleted
		t.Result = payload.Result
		t.UpdatedAt = payload.CompletedAt

	case *TaskExecutionFailed:
		t.Status = TaskStatusFailed
		t.ErrorMessage = payload.ErrorMessage
		t.UpdatedAt = payload.FailedAt

	default:
		return errors.ErrInternal.Msgf("unknown event type: %T", event)
	}

	return nil
}

// ApplySnapshot implements es.Snapshotter
func (t *Task) ApplySnapshot(snapshot es.Snapshot) error {
	switch s := snapshot.(type) {
	case *TaskSnapshot:
		t.ManagerID = s.ManagerID
		t.TaskType = s.TaskType
		t.ScheduledAt = s.ScheduledAt
		t.Payload = s.Payload
		t.Status = TaskStatus(s.Status)
		t.CreatedAt = s.CreatedAt
		t.UpdatedAt = s.UpdatedAt
		t.ExecutedAt = s.ExecutedAt
		t.Result = s.Result
		t.ErrorMessage = s.ErrorMessage
	default:
		return errors.ErrInternal.Msgf("unknown snapshot type: %T", snapshot)
	}

	return nil
}

// ToSnapshot implements es.Snapshotter
func (t *Task) ToSnapshot() es.Snapshot {
	return &TaskSnapshot{
		ID:           t.ID(),
		ManagerID:    t.ManagerID,
		TaskType:     t.TaskType,
		ScheduledAt:  t.ScheduledAt,
		Payload:      t.Payload,
		Status:       string(t.Status),
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
		ExecutedAt:   t.ExecutedAt,
		Result:       t.Result,
		ErrorMessage: t.ErrorMessage,
	}
}