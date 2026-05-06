package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/shipping/internal/domain"
)

type (
	ReturnShipment struct {
		ID                   string
		Reason               string
		ReturnTrackingNumber string
	}

	ReturnShipmentHandler struct {
		shipments domain.ShippingRepository
		publisher ddd.EventPublisher[ddd.Event]
	}
)

func NewReturnShipmentHandler(shipments domain.ShippingRepository, publisher ddd.EventPublisher[ddd.Event]) ReturnShipmentHandler {
	return ReturnShipmentHandler{
		shipments: shipments,
		publisher: publisher,
	}
}

func (h ReturnShipmentHandler) ReturnShipment(ctx context.Context, cmd ReturnShipment) error {
	shipment, err := h.shipments.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := shipment.InitiateReturn(cmd.Reason, cmd.ReturnTrackingNumber)
	if err != nil {
		return err
	}

	err = h.shipments.Save(ctx, shipment)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
