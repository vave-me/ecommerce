package queries

import (
	"context"
	"middleman/messages/internal/domain"
)

type GetConversationByRecipientAndItem struct {
	SenderID    string
	RecipientID string
	ItemID      string
}
type GetConversationByRecipientAndItemHandler struct {
	conversations domain.MiddlemanRepository
}

func NewGetConversationByRecipientAndItemHandler(conversations domain.MiddlemanRepository) GetConversationByRecipientAndItemHandler {
	return GetConversationByRecipientAndItemHandler{conversations: conversations}
}

func (h GetConversationByRecipientAndItemHandler) GetConversationByRecipientAndItem(ctx context.Context, query GetConversationByRecipientAndItem) (*domain.MiddlemanConversation, error) {
	return h.conversations.FindByRecipientAndItem(ctx, query.SenderID, query.RecipientID, query.ItemID)
}
