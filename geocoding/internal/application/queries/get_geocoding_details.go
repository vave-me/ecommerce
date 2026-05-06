package queries

import (
	"context"
	"middleman/geocoding/internal/domain"
)

type GetGeocodingDetails struct {
	Address string
}

type GetGeocodingDetailsHandler struct {
	catalog domain.CatalogRepository
}

func NewGetGeocodingDetailsHandler(catalog domain.CatalogRepository) GetGeocodingDetailsHandler {
	return GetGeocodingDetailsHandler{
		catalog: catalog,
	}
}

func (h GetGeocodingDetailsHandler) GetGeocodingDetails(ctx context.Context, query GetGeocodingDetails) (*domain.CatalogAddress, error) {

	return h.catalog.GetGeocodingDetails(ctx, query.Address)
}
