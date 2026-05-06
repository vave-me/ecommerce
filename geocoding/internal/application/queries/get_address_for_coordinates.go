package queries

import (
	"context"
	"middleman/geocoding/internal/domain"
)

type GetAddressForCoordinates struct {
	AddressID string
	Lat       float64
	Lng       float64
}

type GetAddressForCoordinatesHandler struct {
	catalog domain.CatalogRepository
}

func NewGetAddressForCoordinatesHandler(catalog domain.CatalogRepository) GetAddressForCoordinatesHandler {
	return GetAddressForCoordinatesHandler{
		catalog: catalog,
	}
}

func (h GetAddressForCoordinatesHandler) GetAddressForCoordinates(ctx context.Context, query GetAddressForCoordinates) (*domain.CatalogAddress, error) {

	return h.catalog.GetAddressForCoordinates(ctx, query.Lat, query.Lng)
}
