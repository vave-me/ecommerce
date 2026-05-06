package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/shipping/internal/domain"
)

type (
	AssignCarrier struct {
		ID          string
		CarrierID   string
		CarrierName string
	}

	AssignCarrierHandler struct {
		shipments domain.ShippingRepository
		publisher ddd.EventPublisher[ddd.Event]
	}
)

func NewAssignCarrierHandler(shipments domain.ShippingRepository, publisher ddd.EventPublisher[ddd.Event]) AssignCarrierHandler {
	return AssignCarrierHandler{
		shipments: shipments,
		publisher: publisher,
	}
}

func (h AssignCarrierHandler) AssignCarrier(ctx context.Context, cmd AssignCarrier) error {
	shipment, err := h.shipments.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := shipment.AssignCarrier(cmd.CarrierID, cmd.CarrierName)
	if err != nil {
		return err
	}

	err = h.shipments.Save(ctx, shipment)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
