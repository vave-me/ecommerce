package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/shipping/internal/domain"
)

type (
	StartShipment struct {
		ID string
	}

	StartShipmentHandler struct {
		shipments domain.ShippingRepository
		publisher ddd.EventPublisher[ddd.Event]
	}
)

func NewStartShipmentHandler(shipments domain.ShippingRepository, publisher ddd.EventPublisher[ddd.Event]) StartShipmentHandler {
	return StartShipmentHandler{
		shipments: shipments,
		publisher: publisher,
	}
}

func (h StartShipmentHandler) StartShipment(ctx context.Context, cmd StartShipment) error {
	shipment, err := h.shipments.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := shipment.StartShipment()
	if err != nil {
		return err
	}

	err = h.shipments.Save(ctx, shipment)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
