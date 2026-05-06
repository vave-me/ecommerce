package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/shipping/internal/domain"
)

type (
	MarkShipmentAsDelivered struct {
		ID                 string
		SignedBy           string
		ProofOfDeliveryURL string
	}

	MarkShipmentAsDeliveredHandler struct {
		shipments domain.ShippingRepository
		publisher ddd.EventPublisher[ddd.Event]
	}
)

func NewMarkShipmentAsDeliveredHandler(shipments domain.ShippingRepository, publisher ddd.EventPublisher[ddd.Event]) MarkShipmentAsDeliveredHandler {
	return MarkShipmentAsDeliveredHandler{
		shipments: shipments,
		publisher: publisher,
	}
}

func (h MarkShipmentAsDeliveredHandler) MarkShipmentAsDelivered(ctx context.Context, cmd MarkShipmentAsDelivered) error {
	shipment, err := h.shipments.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := shipment.MarkAsDelivered(cmd.SignedBy, cmd.ProofOfDeliveryURL)
	if err != nil {
		return err
	}

	err = h.shipments.Save(ctx, shipment)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
