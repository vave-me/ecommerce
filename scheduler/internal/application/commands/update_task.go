package commands

import (
	"context"
	"time"

	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/scheduler/internal/domain"
)

type UpdateTask struct {
	TaskID      string
	ScheduledAt *time.Time
	Payload     map[string]string
}

type UpdateTaskHandler struct {
	taskRepo  domain.TaskRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewUpdateTaskHandler(taskRepo domain.TaskRepository, publisher ddd.EventPublisher[ddd.Event]) UpdateTaskHandler {
	return UpdateTaskHandler{
		taskRepo:  taskRepo,
		publisher: publisher,
	}
}

func (h UpdateTaskHandler) UpdateTask(ctx context.Context, cmd UpdateTask) error {
	// Load the task
	task, err := h.taskRepo.Load(ctx, cmd.TaskID)
	if err != nil {
		return errors.Wrap(err, "loading task")
	}

	// Update the task
	event, err := task.Update(cmd.ScheduledAt, cmd.Payload)
	if err != nil {
		return errors.Wrap(err, "updating task")
	}

	// Save the updated task
	if err = h.taskRepo.Save(ctx, task); err != nil {
		return errors.Wrap(err, "saving task")
	}

	// Publish domain event
	if err = h.publisher.Publish(ctx, event); err != nil {
		return errors.Wrap(err, "publishing task updated event")
	}

	return nil
}