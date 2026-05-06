package domain

import "context"

type MiddlemanMessage struct {
	ID             string
	ConversationID string
	SenderID       string
	RecipientID    string
	ItemID         string
	Body           string
	IsRead         bool
}

type MessengerRepository interface {
	SendMessage(ctx context.Context, messageID, conversationID, senderID, receiverID, itemID, body string, isRead bool) error
	Find(ctx context.Context, messageID string) (*MiddlemanMessage, error)
	All(ctx context.Context, conversationID string) ([]*MiddlemanMessage, error)
	Delete(ctx context.Context, messageID string) error
}
