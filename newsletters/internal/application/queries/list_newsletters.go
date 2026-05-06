package queries

import (
	"context"
	"middleman/newsletters/internal/domain"
)

type ListNewsletters struct {
	UserID     string
	Category   string
	ActiveOnly bool
	Page       int
	Limit      int
}

type ListNewslettersHandler struct {
	catalog domain.NewsletterCatalogRepository
}

func NewListNewslettersHandler(catalog domain.NewsletterCatalogRepository) ListNewslettersHandler {
	return ListNewslettersHandler{catalog: catalog}
}

func (h ListNewslettersHandler) ListNewsletters(ctx context.Context, query ListNewsletters) ([]*domain.CatalogNewsletter, int, error) {
	if query.Limit == 0 {
		query.Limit = 20
	}
	
	offset := query.Page * query.Limit
	
	if query.UserID != "" {
		return h.catalog.FindByUser(ctx, query.UserID, query.Limit, offset)
	}
	
	if query.Category != "" {
		return h.catalog.FindByCategory(ctx, query.Category, query.ActiveOnly, query.Limit, offset)
	}
	
	return h.catalog.FindAll(ctx, query.ActiveOnly, query.Limit, offset)
}