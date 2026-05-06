package queries

import (
	"context"
	"middleman/geocoding/internal/domain"
)

type SuggestAddress struct {
	Address string
}

type SuggestAddressHandler struct {
	catalog domain.CatalogRepository
}

func NewSuggestAddressHandler(catalog domain.CatalogRepository) SuggestAddressHandler {
	return SuggestAddressHandler{
		catalog: catalog,
	}
}

func (h SuggestAddressHandler) SuggestAddress(ctx context.Context, query SuggestAddress) ([]*domain.CatalogAddress, error) {

	return h.catalog.SuggestAddress(ctx, query.Address)
}
