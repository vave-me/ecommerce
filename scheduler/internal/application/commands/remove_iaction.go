package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/scheduler/internal/domain"
)

type (
	RemoveAction struct {
		ID          string
		SchedulerID string
	}

	RemoveActionHandler struct {
		actions   domain.ActionRepository
		publisher ddd.EventPublisher[ddd.Event]
	}
)

func NewRemoveActionHandler(actions domain.ActionRepository, publisher ddd.EventPublisher[ddd.Event]) RemoveActionHandler {
	return RemoveActionHandler{
		actions:   actions,
		publisher: publisher,
	}
}

func (h RemoveActionHandler) RemoveAction(ctx context.Context, cmd RemoveAction) error {
	action, err := h.actions.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := action.Remove(cmd.SchedulerID)
	if err != nil {
		return err
	}

	err = h.actions.Save(ctx, action)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
