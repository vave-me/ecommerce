package queries

import (
	"context"
	"middleman/geocoding/internal/domain"
)

type SuggestCity struct {
	Address string
}

type SuggestCityHandler struct {
	catalog domain.CatalogRepository
}

func NewSuggestCityHandler(catalog domain.CatalogRepository) SuggestCityHandler {
	return SuggestCityHandler{
		catalog: catalog,
	}
}

func (h SuggestCityHandler) SuggestCity(ctx context.Context, query SuggestCity) ([]*domain.CatalogAddress, error) {

	return h.catalog.SuggestCity(ctx, query.Address)
}
