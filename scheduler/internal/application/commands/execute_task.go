package commands

import (
	"context"

	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/scheduler/internal/domain"
)

type ExecuteTask struct {
	TaskID string
}

type ExecuteTaskHandler struct {
	taskRepo  domain.TaskRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewExecuteTaskHandler(taskRepo domain.TaskRepository, publisher ddd.EventPublisher[ddd.Event]) ExecuteTaskHandler {
	return ExecuteTaskHandler{
		taskRepo:  taskRepo,
		publisher: publisher,
	}
}

func (h ExecuteTaskHandler) ExecuteTask(ctx context.Context, cmd ExecuteTask) error {
	// Load the task
	task, err := h.taskRepo.Load(ctx, cmd.TaskID)
	if err != nil {
		return errors.Wrap(err, "loading task")
	}

	// Start execution
	event, err := task.StartExecution()
	if err != nil {
		return errors.Wrap(err, "starting task execution")
	}

	// Save the updated task
	if err = h.taskRepo.Save(ctx, task); err != nil {
		return errors.Wrap(err, "saving task")
	}

	// Publish domain event
	if err = h.publisher.Publish(ctx, event); err != nil {
		return errors.Wrap(err, "publishing task execution started event")
	}

	return nil
}