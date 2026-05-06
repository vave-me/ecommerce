package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/services/internal/domain"
)

type DecreaseServicePrice struct {
	ID    string // service_id
	Price int64
}

type DecreaseServicePriceHandler struct {
	services  domain.ServiceRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewDecreaseServicePriceHandler(services domain.ServiceRepository, publisher ddd.EventPublisher[ddd.Event]) DecreaseServicePriceHandler {
	return DecreaseServicePriceHandler{
		services:  services,
		publisher: publisher,
	}
}

func (h DecreaseServicePriceHandler) DecreaseServicePrice(ctx context.Context, cmd DecreaseServicePrice) error {
	service, err := h.services.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := service.DecreaseServicePrice(cmd.Price)
	if err != nil {
		return err
	}

	err = h.services.Save(ctx, service)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
