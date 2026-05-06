package queries

import (
	"context"
	"middleman/posts/internal/domain"
)

type GetPost struct {
	ID string
}

type GetPostHandler struct {
	catalog domain.CatalogRepository
}

func NewGetPostHandler(catalog domain.CatalogRepository) GetPostHandler {
	return GetPostHandler{catalog: catalog}
}

func (h GetPostHandler) GetPost(ctx context.Context, query GetPost) (*domain.CatalogPost, error) {
	return h.catalog.Find(ctx, query.ID)
}
