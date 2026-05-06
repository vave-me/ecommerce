package queries

import (
	"context"
	"middleman/managers/internal/domain"
)

type GetConversationMessages struct {
	ConversationID string
	UserID         string
	Limit          int
	Offset         int
}

type GetConversationMessagesHandler struct {
	readRepo domain.ReadMessagesRepository
}

func NewGetConversationMessagesHandler(readMessages domain.ReadMessagesRepository) GetConversationMessagesHandler {
	return GetConversationMessagesHandler{
		readRepo: readMessages,
	}
}

func (h GetConversationMessagesHandler) GetConversationMessages(ctx context.Context, query GetConversationMessages) ([]*domain.ReadMessage, error) {
	// Set default pagination if not provided
	limit := query.Limit
	if limit <= 0 {
		limit = 50 // Default page size for messages
	}
	if limit > 200 {
		limit = 200 // Max page size for messages
	}

	offset := query.Offset
	if offset < 0 {
		offset = 0
	}

	// Execute query
	messages, err := h.readRepo.GetConversationMessages(ctx, query.ConversationID, query.UserID, limit, offset)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
