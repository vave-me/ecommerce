package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/shipping/internal/domain"
)

type (
	CancelShipment struct {
		ID     string
		Reason string
	}

	CancelShipmentHandler struct {
		shipments domain.ShippingRepository
		publisher ddd.EventPublisher[ddd.Event]
	}
)

func NewCancelShipmentHandler(shipments domain.ShippingRepository, publisher ddd.EventPublisher[ddd.Event]) CancelShipmentHandler {
	return CancelShipmentHandler{
		shipments: shipments,
		publisher: publisher,
	}
}

func (h CancelShipmentHandler) CancelShipment(ctx context.Context, cmd CancelShipment) error {
	shipment, err := h.shipments.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := shipment.CancelShipment(cmd.Reason)
	if err != nil {
		return err
	}

	err = h.shipments.Save(ctx, shipment)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
