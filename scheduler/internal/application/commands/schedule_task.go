package commands

import (
	"context"
	"time"

	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/scheduler/internal/domain"
)

type ScheduleTask struct {
	ID          string
	ManagerID   string
	TaskType    string
	ScheduledAt time.Time
	Payload     map[string]string
}

type ScheduleTaskHandler struct {
	taskRepo  domain.TaskRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewScheduleTaskHandler(taskRepo domain.TaskRepository, publisher ddd.EventPublisher[ddd.Event]) ScheduleTaskHandler {
	return ScheduleTaskHandler{
		taskRepo:  taskRepo,
		publisher: publisher,
	}
}

func (h ScheduleTaskHandler) ScheduleTask(ctx context.Context, cmd ScheduleTask) error {
	// Create new task aggregate
	task, err := domain.CreateTask(cmd.ID, cmd.ManagerID, cmd.TaskType, cmd.ScheduledAt, cmd.Payload)
	if err != nil {
		return errors.Wrap(err, "creating task")
	}

	// Save the task
	if err = h.taskRepo.Save(ctx, task); err != nil {
		return errors.Wrap(err, "saving task")
	}

	// Get the first event (the one returned from Schedule method)
	if len(task.Events()) > 0 {
		if err = h.publisher.Publish(ctx, task.Events()[0]); err != nil {
			return errors.Wrap(err, "publishing task scheduled event")
		}
	}

	return nil
}