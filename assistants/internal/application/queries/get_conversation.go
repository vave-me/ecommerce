package queries

import (
	"context"
	"middleman/assistants/internal/domain"
)

type GetConversation struct {
	ConversationID string
	UserID         string
}

type GetConversationHandler struct {
	readRepo domain.ReadConversationRepository
}

func NewGetConversationHandler(readRepo domain.ReadConversationRepository) GetConversationHandler {
	return GetConversationHandler{
		readRepo: readRepo,
	}
}

func (h GetConversationHandler) GetConversation(ctx context.Context, query GetConversation) (*domain.ReadConversation, error) {
	// Execute query
	conversation, err := h.readRepo.GetConversation(ctx, query.ConversationID, query.UserID)
	if err != nil {
		return nil, err
	}

	return conversation, nil
}
