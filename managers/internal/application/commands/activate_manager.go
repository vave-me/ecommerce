package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/managers/internal/domain"
)

type ActivateManager struct {
	ID string `json:"id"`
}

type ActivateManagerHandler struct {
	managers  domain.ManagerRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewActivateManagerHandler(managers domain.ManagerRepository, publisher ddd.EventPublisher[ddd.Event]) ActivateManagerHandler {
	return ActivateManagerHandler{
		managers:  managers,
		publisher: publisher,
	}
}

func (h ActivateManagerHandler) ActivateManager(ctx context.Context, cmd ActivateManager) error {
	// Step 1: Load manager aggregate
	manager, err := h.managers.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	// Step 2: Activate manager in aggregate
	event, err := manager.Activate()
	if err != nil {
		return err
	}

	// Step 3: Save manager aggregate
	if err = h.managers.Save(ctx, manager); err != nil {
		return err
	}

	// Step 4: Publish domain event
	err = h.publisher.Publish(ctx, event)
	if err != nil {
		return err
	}

	return nil
}
