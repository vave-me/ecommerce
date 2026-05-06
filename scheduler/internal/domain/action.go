package domain

import (
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
	"time"
)

const ActionAggregate = "scheduler.Action"

var (
	ErrActionAlreadyCreated = errors.Wrap(errors.ErrBadRequest, "the action cannot be recreated")
	ErrActionHasNoPayload   = errors.Wrap(errors.ErrBadRequest, "the action has no payload")
	ErrTaskCannotBeBlank    = errors.Wrap(errors.ErrBadRequest, "the task cannot be blank")
	ErrInvalidExecutionTime = errors.Wrap(errors.ErrBadRequest, "execution time must be in the future")
)

type Action struct {
	es.Aggregate
	SchedulerID          string
	NaturalLanguageTask  string    // Natural language description of the task
	ExecutionTime        time.Time // When to execute the task
	Status               string    // pending, executing, completed, failed
	CreatedAt            time.Time
	ExecutedAt           *time.Time
	Result               string // Result from LLM/Assistant service
	ErrorMessage         string // Error message if failed
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Action)(nil)

func NewAction(id string) *Action {
	return &Action{
		Aggregate: es.NewAggregate(id, ActionAggregate),
	}
}

func (i *Action) AddAction(schedulerID string, task string, executionTime time.Time) (ddd.Event, error) {
	if task == "" {
		return nil, ErrTaskCannotBeBlank
	}
	
	if executionTime.Before(time.Now()) {
		return nil, ErrInvalidExecutionTime
	}

	i.AddEvent(ActionAddedEvent, &ActionAdded{
		ID:                  i.ID(),
		SchedulerID:         schedulerID,
		NaturalLanguageTask: task,
		ExecutionTime:       executionTime,
		Status:              "pending",
		CreatedAt:           time.Now(),
	})

	return ddd.NewEvent(ActionAddedEvent, i), nil
}

func (i *Action) Remove(schedulerID string) (ddd.Event, error) {
	if i.Status == "completed" || i.Status == "executing" {
		return nil, errors.Wrap(errors.ErrBadRequest, "cannot remove completed or executing action")
	}

	i.AddEvent(ActionRemovedEvent, &ActionRemoved{
		SchedulerID: schedulerID,
		ActionID:    i.ID(),
	})

	return ddd.NewEvent(ActionRemovedEvent, i), nil
}

func (i *Action) UpdateStatus(status string, result string, errorMessage string) (ddd.Event, error) {
	validStatuses := map[string]bool{
		"pending":   true,
		"executing": true,
		"completed": true,
		"failed":    true,
	}
	
	if !validStatuses[status] {
		return nil, errors.Wrap(errors.ErrBadRequest, "invalid status")
	}

	// Validate state transitions
	switch i.Status {
	case "completed":
		return nil, errors.Wrap(errors.ErrBadRequest, "cannot update a completed action")
	case "failed":
		if status != "pending" {
			return nil, errors.Wrap(errors.ErrBadRequest, "failed actions can only be reset to pending")
		}
	case "executing":
		if status == "pending" {
			return nil, errors.Wrap(errors.ErrBadRequest, "executing actions cannot go back to pending")
		}
	case "pending":
		if status == "completed" || status == "failed" {
			return nil, errors.Wrap(errors.ErrBadRequest, "pending actions must go through executing state first")
		}
	}

	// Validate error message is only set for failed status
	if status != "failed" && errorMessage != "" {
		return nil, errors.Wrap(errors.ErrBadRequest, "error message can only be set for failed status")
	}
	
	// Validate result is only set for completed status
	if status != "completed" && result != "" {
		return nil, errors.Wrap(errors.ErrBadRequest, "result can only be set for completed status")
	}

	now := time.Now()
	i.AddEvent(ActionUpdatedEvent, &ActionUpdated{
		ActionID:     i.ID(),
		SchedulerID:  i.SchedulerID,
		Status:       status,
		ExecutedAt:   &now,
		Result:       result,
		ErrorMessage: errorMessage,
	})

	return ddd.NewEvent(ActionUpdatedEvent, i), nil
}
func (b *Action) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *ActionAdded:
		b.SchedulerID = payload.SchedulerID
		b.NaturalLanguageTask = payload.NaturalLanguageTask
		b.ExecutionTime = payload.ExecutionTime
		b.Status = payload.Status
		b.CreatedAt = payload.CreatedAt
	case *ActionUpdated:
		b.Status = payload.Status
		b.ExecutedAt = payload.ExecutedAt
		b.Result = payload.Result
		b.ErrorMessage = payload.ErrorMessage
	case *ActionRemoved:
		// Action is marked for removal
	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", b, event.EventName(), payload)
	}

	return nil
}

func (b *Action) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *ActionVi:
		b.SchedulerID = ss.SchedulerID
		b.NaturalLanguageTask = ss.NaturalLanguageTask
		b.ExecutionTime = ss.ExecutionTime
		b.Status = ss.Status
		b.CreatedAt = ss.CreatedAt
		b.ExecutedAt = ss.ExecutedAt
		b.Result = ss.Result
		b.ErrorMessage = ss.ErrorMessage
	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", b, snapshot)
	}

	return nil
}

func (b *Action) ToSnapshot() es.Snapshot {
	return &ActionVi{
		SchedulerID:         b.SchedulerID,
		NaturalLanguageTask: b.NaturalLanguageTask,
		ExecutionTime:       b.ExecutionTime,
		Status:              b.Status,
		CreatedAt:           b.CreatedAt,
		ExecutedAt:          b.ExecutedAt,
		Result:              b.Result,
		ErrorMessage:        b.ErrorMessage,
	}
}
