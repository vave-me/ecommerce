package queries

import (
	"context"
	"middleman/newsletters/internal/domain"
)

type GetTemplate struct {
	ID string
}

type GetTemplateHandler struct {
	catalog domain.TemplateCatalogRepository
}

func NewGetTemplateHandler(catalog domain.TemplateCatalogRepository) GetTemplateHandler {
	return GetTemplateHandler{catalog: catalog}
}

func (h GetTemplateHandler) GetTemplate(ctx context.Context, query GetTemplate) (*domain.CatalogTemplate, error) {
	return h.catalog.Find(ctx, query.ID)
}