package queries

import (
	"context"
	"middleman/messages/internal/domain"
)

type GetConversations struct {
	UserID string
}
type GetConversationsHandler struct {
	conversations domain.MiddlemanRepository
}

func NewGetConversationsHandler(conversations domain.MiddlemanRepository) GetConversationsHandler {
	return GetConversationsHandler{conversations: conversations}
}

func (h GetConversationsHandler) GetConversations(ctx context.Context, query GetConversations) ([]*domain.MiddlemanConversation, error) {
	return h.conversations.All(ctx, query.UserID)
}
