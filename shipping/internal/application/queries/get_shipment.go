package queries

import (
	"context"
	"middleman/shipping/internal/domain"
)

type (
	GetShipment struct {
		ID string
	}

	GetShipmentHandler struct {
		catalog domain.ShippingCatalogRepository
	}
)

func NewGetShipmentHandler(catalog domain.ShippingCatalogRepository) GetShipmentHandler {
	return GetShipmentHandler{
		catalog: catalog,
	}
}

func (h GetShipmentHandler) GetShipment(ctx context.Context, query GetShipment) (*domain.CatalogShipment, error) {
	return h.catalog.Find(ctx, query.ID)
}
