package queries

import (
	"context"
	"middleman/geocoding/internal/domain"
)

type GetCoordinatesForAddress struct {
	Address string
}

type GetCoordinatesForAddressHandler struct {
	catalog domain.CatalogRepository
}

func NewGetCoordinatesForAddressHandler(catalog domain.CatalogRepository) GetCoordinatesForAddressHandler {
	return GetCoordinatesForAddressHandler{
		catalog: catalog,
	}
}

func (h GetCoordinatesForAddressHandler) GetCoordinatesForAddress(ctx context.Context, query GetCoordinatesForAddress) (*domain.CatalogAddress, error) {

	return h.catalog.GetCoordinatesForAddress(ctx, query.Address)
}
