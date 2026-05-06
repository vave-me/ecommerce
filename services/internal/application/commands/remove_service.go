package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/services/internal/domain"
)

type RemoveService struct {
	ID     string
	UserID string
}

type RemoveServiceHandler struct {
	services  domain.ServiceRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewRemoveServiceHandler(services domain.ServiceRepository, publisher ddd.EventPublisher[ddd.Event]) RemoveServiceHandler {
	return RemoveServiceHandler{
		services:  services,
		publisher: publisher,
	}
}

func (h RemoveServiceHandler) RemoveService(ctx context.Context, cmd RemoveService) error {
	service, err := h.services.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := service.RemoveService(cmd.ID, cmd.UserID)
	if err != nil {
		return err
	}

	err = h.services.Save(ctx, service)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
