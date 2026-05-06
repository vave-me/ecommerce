package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/managers/internal/domain"
)

type DeactivateManager struct {
	ID string `json:"id"`
}

type DeactivateManagerHandler struct {
	managers  domain.ManagerRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewDeactivateManagerHandler(managers domain.ManagerRepository, publisher ddd.EventPublisher[ddd.Event]) DeactivateManagerHandler {
	return DeactivateManagerHandler{
		managers:  managers,
		publisher: publisher,
	}
}

func (h DeactivateManagerHandler) DeactivateManager(ctx context.Context, cmd DeactivateManager) error {

	manager, err := h.managers.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := manager.Deactivate()
	if err != nil {
		return err
	}

	if err = h.managers.Save(ctx, manager); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
