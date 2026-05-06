package domain

import "context"

type MiddlemanConversation struct {
	ID          string
	RecipientID string
	SenderID    string
	ItemID      string
	Active      bool
}

type MiddlemanRepository interface {
	Add(ctx context.Context, conversationID, recipientID, senderID, itemID string) error
	Find(ctx context.Context, conversationID string) (*MiddlemanConversation, error)
	FindByRecipientAndItem(ctx context.Context, senderID, recipientID, itemID string) (*MiddlemanConversation, error)
	All(ctx context.Context, userID string) ([]*MiddlemanConversation, error)
}
