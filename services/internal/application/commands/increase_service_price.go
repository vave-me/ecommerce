package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/services/internal/domain"
)

type IncreaseServicePrice struct {
	ID    string
	Price int64
}

type IncreaseServicePriceHandler struct {
	services  domain.ServiceRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewIncreaseServicePriceHandler(services domain.ServiceRepository, publisher ddd.EventPublisher[ddd.Event]) IncreaseServicePriceHandler {
	return IncreaseServicePriceHandler{
		services:  services,
		publisher: publisher,
	}
}

func (h IncreaseServicePriceHandler) IncreaseServicePrice(ctx context.Context, cmd IncreaseServicePrice) error {
	service, err := h.services.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := service.IncreaseServicePrice(cmd.Price)
	if err != nil {
		return err
	}

	err = h.services.Save(ctx, service)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
