package queries

import (
	"context"
	"middleman/messages/internal/domain"
)

type GetConversation struct {
	ID string
}
type GetConversationHandler struct {
	conversations domain.MiddlemanRepository
}

func NewGetConversationHandler(conversations domain.MiddlemanRepository) GetConversationHandler {
	return GetConversationHandler{conversations: conversations}
}

func (h GetConversationHandler) GetConversation(ctx context.Context, query GetConversation) (*domain.MiddlemanConversation, error) {
	return h.conversations.Find(ctx, query.ID)
}
