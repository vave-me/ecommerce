package queries

import (
	"context"
	"middleman/shipping/internal/domain"
)

type (
	GetShipmentHistory struct {
		ID string
	}

	GetShipmentHistoryHandler struct {
		shippingRepo domain.ShippingRepository
	}
)

func NewGetShipmentHistoryHandler(shippingRepo domain.ShippingRepository) GetShipmentHistoryHandler {
	return GetShipmentHistoryHandler{
		shippingRepo: shippingRepo,
	}
}

func (h GetShipmentHistoryHandler) GetShipmentHistory(ctx context.Context, query GetShipmentHistory) ([]*domain.ShipmentEvent, error) {
	// TODO: Implement event history retrieval from event store
	// For now, return empty array
	return []*domain.ShipmentEvent{}, nil
}
