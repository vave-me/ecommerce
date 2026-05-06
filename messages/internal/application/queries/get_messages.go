package queries

import (
	"context"
	"middleman/messages/internal/domain"
)

type GetMessages struct {
	ConversationID string
}
type GetMessagesHandler struct {
	messages domain.MessengerRepository
}

func NewGetMessagesHandler(messages domain.MessengerRepository) GetMessagesHandler {
	return GetMessagesHandler{messages: messages}
}

func (h GetMessagesHandler) GetMessages(ctx context.Context, query GetMessages) ([]*domain.MiddlemanMessage, error) {
	return h.messages.All(ctx, query.ConversationID)
}
