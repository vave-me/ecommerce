package queries

import (
	"context"
	"middleman/geocoding/internal/domain"
)

type GetGeocodingCache struct {
	Address string
}

type GetGeocodingCacheHandler struct {
	catalog domain.CatalogRepository
}

func NewGetGeocodingCacheHandler(catalog domain.CatalogRepository) GetGeocodingCacheHandler {
	return GetGeocodingCacheHandler{
		catalog: catalog,
	}
}

func (h GetGeocodingCacheHandler) GetGeocodingCache(ctx context.Context, query GetGeocodingCache) (*domain.CatalogAddress, error) {

	return h.catalog.GetGeocodingCache(ctx, query.Address)
}
