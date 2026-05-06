package queries

import (
	"context"
	"middleman/assistants/internal/domain"
)

type GetAssistant struct {
	ID string
}

type GetAssistantHandler struct {
	assistants domain.CatalogRepository
}

func NewGetAssistantHandler(assistants domain.CatalogRepository) GetAssistantHandler {
	return GetAssistantHandler{
		assistants: assistants,
	}
}

func (h GetAssistantHandler) GetAssistant(ctx context.Context, query GetAssistant) (*domain.CatalogAssistant, error) {
	return h.assistants.Find(ctx, query.ID)
}
