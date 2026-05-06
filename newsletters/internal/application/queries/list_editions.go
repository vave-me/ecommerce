package queries

import (
	"context"
	"middleman/newsletters/internal/domain"
)

type ListEditions struct {
	NewsletterID string
	Status       string
	Page         int
	Limit        int
}

type ListEditionsHandler struct {
	catalog domain.EditionCatalogRepository
}

func NewListEditionsHandler(catalog domain.EditionCatalogRepository) ListEditionsHandler {
	return ListEditionsHandler{catalog: catalog}
}

func (h ListEditionsHandler) ListEditions(ctx context.Context, query ListEditions) ([]*domain.CatalogEdition, int, error) {
	if query.Limit == 0 {
		query.Limit = 20
	}
	
	offset := query.Page * query.Limit
	
	return h.catalog.FindByNewsletter(ctx, query.NewsletterID, query.Status, query.Limit, offset)
}