package commands

import (
	"context"

	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/scheduler/internal/domain"
)

type CancelTask struct {
	TaskID string
}

type CancelTaskHandler struct {
	taskRepo  domain.TaskRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewCancelTaskHandler(taskRepo domain.TaskRepository, publisher ddd.EventPublisher[ddd.Event]) CancelTaskHandler {
	return CancelTaskHandler{
		taskRepo:  taskRepo,
		publisher: publisher,
	}
}

func (h CancelTaskHandler) CancelTask(ctx context.Context, cmd CancelTask) error {
	// Load the task
	task, err := h.taskRepo.Load(ctx, cmd.TaskID)
	if err != nil {
		return errors.Wrap(err, "loading task")
	}

	// Cancel the task
	event, err := task.Cancel()
	if err != nil {
		return errors.Wrap(err, "cancelling task")
	}

	// Save the updated task
	if err = h.taskRepo.Save(ctx, task); err != nil {
		return errors.Wrap(err, "saving task")
	}

	// Publish domain event if one was returned
	if event != nil {
		if err = h.publisher.Publish(ctx, event); err != nil {
			return errors.Wrap(err, "publishing task cancelled event")
		}
	}

	return nil
}