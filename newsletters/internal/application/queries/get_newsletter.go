package queries

import (
	"context"
	"middleman/newsletters/internal/domain"
)

type GetNewsletter struct {
	ID string
}

type GetNewsletterHandler struct {
	catalog domain.NewsletterCatalogRepository
}

func NewGetNewsletterHandler(catalog domain.NewsletterCatalogRepository) GetNewsletterHandler {
	return GetNewsletterHandler{catalog: catalog}
}

func (h GetNewsletterHandler) GetNewsletter(ctx context.Context, query GetNewsletter) (*domain.CatalogNewsletter, error) {
	return h.catalog.Find(ctx, query.ID)
}