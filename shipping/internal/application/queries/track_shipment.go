package queries

import (
	"context"
	"middleman/shipping/internal/domain"
)

type (
	TrackShipment struct {
		TrackingNumber string
	}

	TrackShipmentHandler struct {
		catalog domain.ShippingCatalogRepository
	}
)

func NewTrackShipmentHandler(catalog domain.ShippingCatalogRepository) TrackShipmentHandler {
	return TrackShipmentHandler{
		catalog: catalog,
	}
}

func (h TrackShipmentHandler) TrackShipment(ctx context.Context, query TrackShipment) (*domain.CatalogShipment, error) {
	return h.catalog.GetByTrackingNumber(ctx, query.TrackingNumber)
}
