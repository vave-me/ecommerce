package commands

import (
	"context"

	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/scheduler/internal/domain"
)

type FailTask struct {
	TaskID       string
	ErrorMessage string
}

type FailTaskHandler struct {
	taskRepo  domain.TaskRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewFailTaskHandler(taskRepo domain.TaskRepository, publisher ddd.EventPublisher[ddd.Event]) FailTaskHandler {
	return FailTaskHandler{
		taskRepo:  taskRepo,
		publisher: publisher,
	}
}

func (h FailTaskHandler) FailTask(ctx context.Context, cmd FailTask) error {
	// Load the task
	task, err := h.taskRepo.Load(ctx, cmd.TaskID)
	if err != nil {
		return errors.Wrap(err, "loading task")
	}

	// Fail the task
	event, err := task.Fail(cmd.ErrorMessage)
	if err != nil {
		return errors.Wrap(err, "failing task")
	}

	// Save the updated task
	if err = h.taskRepo.Save(ctx, task); err != nil {
		return errors.Wrap(err, "saving task")
	}

	// Publish domain event
	if err = h.publisher.Publish(ctx, event); err != nil {
		return errors.Wrap(err, "publishing task failed event")
	}

	return nil
}