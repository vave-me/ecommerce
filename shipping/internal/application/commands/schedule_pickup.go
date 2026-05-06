package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/shipping/internal/domain"
)

type (
	SchedulePickup struct {
		ID           string
		PickupTime   string
		Instructions string
	}

	SchedulePickupHandler struct {
		shipments domain.ShippingRepository
		publisher ddd.EventPublisher[ddd.Event]
	}
)

func NewSchedulePickupHandler(shipments domain.ShippingRepository, publisher ddd.EventPublisher[ddd.Event]) SchedulePickupHandler {
	return SchedulePickupHandler{
		shipments: shipments,
		publisher: publisher,
	}
}

func (h SchedulePickupHandler) SchedulePickup(ctx context.Context, cmd SchedulePickup) error {
	shipment, err := h.shipments.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := shipment.SchedulePickup(cmd.PickupTime, cmd.Instructions)
	if err != nil {
		return err
	}

	err = h.shipments.Save(ctx, shipment)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
