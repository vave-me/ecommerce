package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/scheduler/internal/domain"
)

type (
	UpdateActionStatus struct {
		ID           string
		Status       string
		Result       string
		ErrorMessage string
	}

	UpdateActionStatusHandler struct {
		actions   domain.ActionRepository
		publisher ddd.EventPublisher[ddd.Event]
	}
)

func NewUpdateActionStatusHandler(actions domain.ActionRepository, publisher ddd.EventPublisher[ddd.Event]) UpdateActionStatusHandler {
	return UpdateActionStatusHandler{
		actions:   actions,
		publisher: publisher,
	}
}

func (h UpdateActionStatusHandler) UpdateActionStatus(ctx context.Context, cmd UpdateActionStatus) error {
	action, err := h.actions.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := action.UpdateStatus(cmd.Status, cmd.Result, cmd.ErrorMessage)
	if err != nil {
		return err
	}

	err = h.actions.Save(ctx, action)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}