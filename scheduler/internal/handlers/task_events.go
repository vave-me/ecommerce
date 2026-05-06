package handlers

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/di"
	"middleman/scheduler/internal/constants"
	"middleman/scheduler/internal/domain"
)

type taskHandlers[T ddd.Event] struct {
	catalogTasks domain.CatalogTaskRepository
	logger       zerolog.Logger
}

var _ ddd.EventHandler[ddd.Event] = (*taskHandlers[ddd.Event])(nil)

func NewTaskHandlers(catalogTasks domain.CatalogTaskRepository, logger zerolog.Logger) ddd.EventHandler[ddd.Event] {
	return taskHandlers[ddd.Event]{
		catalogTasks: catalogTasks,
		logger:       logger.With().Str("handlers", "taskEvents").Logger(),
	}
}

func RegisterTaskHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		taskHandlers := di.Get(ctx, constants.TaskHandlersKey).(ddd.EventHandler[ddd.Event])
		return taskHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get(constants.DomainDispatcherKey).(ddd.EventSubscriber[ddd.Event])

	subscriber.Subscribe(handlers,
		domain.TaskScheduledEvent,
		domain.TaskCancelledEvent,
		domain.TaskUpdatedEvent,
		domain.TaskExecutionStartedEvent,
		domain.TaskExecutionCompletedEvent,
		domain.TaskExecutionFailedEvent,
	)
}

func (h taskHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.TaskScheduledEvent:
		return h.onTaskScheduled(ctx, event)
	case domain.TaskCancelledEvent:
		return h.onTaskCancelled(ctx, event)
	case domain.TaskUpdatedEvent:
		return h.onTaskUpdated(ctx, event)
	case domain.TaskExecutionStartedEvent:
		return h.onTaskExecutionStarted(ctx, event)
	case domain.TaskExecutionCompletedEvent:
		return h.onTaskExecutionCompleted(ctx, event)
	case domain.TaskExecutionFailedEvent:
		return h.onTaskExecutionFailed(ctx, event)
	}

	return nil
}

func (h taskHandlers[T]) onTaskScheduled(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TaskScheduled)

	task := &domain.CatalogTask{
		ID:          payload.TaskID,
		ManagerID:   payload.ManagerID,
		TaskType:    payload.TaskType,
		ScheduledAt: payload.ScheduledAt,
		Payload:     payload.Payload,
		Status:      string(domain.TaskStatusPending),
		CreatedAt:   payload.CreatedAt,
		UpdatedAt:   payload.CreatedAt,
	}

	err := h.catalogTasks.Add(ctx, task)
	if err != nil {
		h.logger.Error().Err(err).Str("taskID", task.ID).Msg("failed to add task to catalog")
		return errors.Wrap(err, "adding task to catalog")
	}

	h.logger.Info().Str("taskID", task.ID).Msg("task scheduled")
	return nil
}

func (h taskHandlers[T]) onTaskCancelled(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TaskCancelled)

	updates := map[string]interface{}{
		"status": string(domain.TaskStatusCancelled),
	}

	err := h.catalogTasks.Update(ctx, payload.TaskID, updates)
	if err != nil {
		h.logger.Error().Err(err).Str("taskID", payload.TaskID).Msg("failed to update task status to cancelled")
		return errors.Wrap(err, "updating task status")
	}

	h.logger.Info().Str("taskID", payload.TaskID).Msg("task cancelled")
	return nil
}

func (h taskHandlers[T]) onTaskUpdated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TaskUpdated)

	updates := map[string]interface{}{}

	if payload.ScheduledAt != nil {
		updates["scheduled_at"] = *payload.ScheduledAt
	}

	if len(payload.Payload) > 0 {
		updates["payload"] = payload.Payload
	}

	if len(updates) == 0 {
		return nil
	}

	err := h.catalogTasks.Update(ctx, payload.TaskID, updates)
	if err != nil {
		h.logger.Error().Err(err).Str("taskID", payload.TaskID).Msg("failed to update task")
		return errors.Wrap(err, "updating task")
	}

	h.logger.Info().Str("taskID", payload.TaskID).Msg("task updated")
	return nil
}

func (h taskHandlers[T]) onTaskExecutionStarted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TaskExecutionStarted)

	updates := map[string]interface{}{
		"status":      string(domain.TaskStatusRunning),
		"executed_at": payload.StartedAt,
	}

	err := h.catalogTasks.Update(ctx, payload.TaskID, updates)
	if err != nil {
		h.logger.Error().Err(err).Str("taskID", payload.TaskID).Msg("failed to update task execution start")
		return errors.Wrap(err, "updating task execution start")
	}

	h.logger.Info().Str("taskID", payload.TaskID).Msg("task execution started")
	return nil
}

func (h taskHandlers[T]) onTaskExecutionCompleted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TaskExecutionCompleted)

	updates := map[string]interface{}{
		"status": string(domain.TaskStatusCompleted),
		"result": payload.Result,
	}

	err := h.catalogTasks.Update(ctx, payload.TaskID, updates)
	if err != nil {
		h.logger.Error().Err(err).Str("taskID", payload.TaskID).Msg("failed to update task completion")
		return errors.Wrap(err, "updating task completion")
	}

	h.logger.Info().Str("taskID", payload.TaskID).Msg("task execution completed")
	return nil
}

func (h taskHandlers[T]) onTaskExecutionFailed(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TaskExecutionFailed)

	updates := map[string]interface{}{
		"status":        string(domain.TaskStatusFailed),
		"error_message": payload.ErrorMessage,
	}

	err := h.catalogTasks.Update(ctx, payload.TaskID, updates)
	if err != nil {
		h.logger.Error().Err(err).Str("taskID", payload.TaskID).Msg("failed to update task failure")
		return errors.Wrap(err, "updating task failure")
	}

	h.logger.Error().Str("taskID", payload.TaskID).Str("error", payload.ErrorMessage).Msg("task execution failed")
	return nil
}
