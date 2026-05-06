package commands

import (
	"context"
	"time"
	"middleman/internal/ddd"
	"middleman/scheduler/internal/domain"
)

type (
	AddAction struct {
		ID                  string
		SchedulerID         string
		NaturalLanguageTask string
		ExecutionTime       time.Time
	}

	AddActionHandler struct {
		actions   domain.ActionRepository
		publisher ddd.EventPublisher[ddd.Event]
	}
)

func NewAddActionHandler(actions domain.ActionRepository, publisher ddd.EventPublisher[ddd.Event]) AddActionHandler {
	return AddActionHandler{
		actions:   actions,
		publisher: publisher,
	}
}

func (h AddActionHandler) AddAction(ctx context.Context, cmd AddAction) error {
	action, err := h.actions.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := action.AddAction(cmd.SchedulerID, cmd.NaturalLanguageTask, cmd.ExecutionTime)
	if err != nil {
		return err
	}

	err = h.actions.Save(ctx, action)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
