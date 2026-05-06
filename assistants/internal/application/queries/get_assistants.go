package queries

import (
	"context"
	"middleman/assistants/internal/domain"
)

type GetAssistants struct {
	UserID string
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

type GetAssistantsHandler struct {
	readModel domain.CatalogRepository
}

func NewGetAssistantsHandler(readModel domain.CatalogRepository) GetAssistantsHandler {
	return GetAssistantsHandler{
		readModel: readModel,
	}
}

func (h GetAssistantsHandler) GetAssistants(ctx context.Context, query GetAssistants) ([]*domain.CatalogAssistant, error) {
	// Call the read model to get all assistants
	assistants, err := h.readModel.FindAll(ctx, query.UserID)
	if err != nil {
		return nil, err
	}

	return assistants, nil
}
