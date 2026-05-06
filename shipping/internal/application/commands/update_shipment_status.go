package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/shipping/internal/domain"
)

type (
	UpdateShipmentStatus struct {
		ID       string
		Status   string
		Location string
		Notes    string
	}

	UpdateShipmentStatusHandler struct {
		shipments domain.ShippingRepository
		publisher ddd.EventPublisher[ddd.Event]
	}
)

func NewUpdateShipmentStatusHandler(shipments domain.ShippingRepository, publisher ddd.EventPublisher[ddd.Event]) UpdateShipmentStatusHandler {
	return UpdateShipmentStatusHandler{
		shipments: shipments,
		publisher: publisher,
	}
}

func (h UpdateShipmentStatusHandler) UpdateShipmentStatus(ctx context.Context, cmd UpdateShipmentStatus) error {
	shipment, err := h.shipments.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := shipment.UpdateStatus(cmd.Status, cmd.Location, cmd.Notes)
	if err != nil {
		return err
	}

	err = h.shipments.Save(ctx, shipment)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}