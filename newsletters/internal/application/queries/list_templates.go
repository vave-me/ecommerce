package queries

import (
	"context"
	"middleman/newsletters/internal/domain"
)

type ListTemplates struct {
	UserID     string
	PublicOnly bool
	Page       int
	Limit      int
}

type ListTemplatesHandler struct {
	catalog domain.TemplateCatalogRepository
}

func NewListTemplatesHandler(catalog domain.TemplateCatalogRepository) ListTemplatesHandler {
	return ListTemplatesHandler{catalog: catalog}
}

func (h ListTemplatesHandler) ListTemplates(ctx context.Context, query ListTemplates) ([]*domain.CatalogTemplate, int, error) {
	if query.Limit == 0 {
		query.Limit = 20
	}
	
	offset := query.Page * query.Limit
	
	if query.PublicOnly {
		return h.catalog.FindPublic(ctx, query.Limit, offset)
	}
	
	if query.UserID != "" {
		return h.catalog.FindByUser(ctx, query.UserID, query.Limit, offset)
	}
	
	// Default to public templates
	return h.catalog.FindPublic(ctx, query.Limit, offset)
}