package commands

import (
	"context"

	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/scheduler/internal/domain"
)

type CompleteTask struct {
	TaskID string
	Result string
}

type CompleteTaskHandler struct {
	taskRepo  domain.TaskRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewCompleteTaskHandler(taskRepo domain.TaskRepository, publisher ddd.EventPublisher[ddd.Event]) CompleteTaskHandler {
	return CompleteTaskHandler{
		taskRepo:  taskRepo,
		publisher: publisher,
	}
}

func (h CompleteTaskHandler) CompleteTask(ctx context.Context, cmd CompleteTask) error {
	// Load the task
	task, err := h.taskRepo.Load(ctx, cmd.TaskID)
	if err != nil {
		return errors.Wrap(err, "loading task")
	}

	// Complete the task
	event, err := task.Complete(cmd.Result)
	if err != nil {
		return errors.Wrap(err, "completing task")
	}

	// Save the updated task
	if err = h.taskRepo.Save(ctx, task); err != nil {
		return errors.Wrap(err, "saving task")
	}

	// Publish domain event
	if err = h.publisher.Publish(ctx, event); err != nil {
		return errors.Wrap(err, "publishing task completed event")
	}

	return nil
}