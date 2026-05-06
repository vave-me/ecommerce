package domain

import (
	"context"
	"middleman/managers/internal/models"
)

type MessagesRepository interface {
	// Conversation Management
	StartConversation(ctx context.Context, senderID, recipientID, itemID string) (*models.ConversationResponse, error)
	RestoreConversation(ctx context.Context, conversationID string) (*models.ConversationStatusResponse, error)
	ArchiveConversation(ctx context.Context, conversationID string) (*models.ConversationStatusResponse, error)
	GetConversation(ctx context.Context, conversationID string) (*models.Conversation, error)
	GetConversationByRecipientAndItem(ctx context.Context, recipientID, itemID string) (*models.Conversation, error)
	GetConversations(ctx context.Context, userID string, page, limit int64) (*models.ConversationsResponse, error)
	GetActiveConversations(ctx context.Context, userID string, page, limit int64) (*models.ConversationsResponse, error)

	// Message Management
	SendMessage(ctx context.Context, messageID, conversationID, senderID, recipientID, itemID, body string, isRead bool) (*models.MessageSentResponse, error)
	DeleteMessage(ctx context.Context, messageID string) error
	GetMessage(ctx context.Context, messageID string) (*models.Message, error)
	GetMessages(ctx context.Context, conversationID string, page, limit int64) (*models.MessagesResponse, error)

	// Additional methods needed by tool service
	MarkMessageAsRead(ctx context.Context, messageID, userID string) error
	GetUnreadMessagesCount(ctx context.Context, userID string) (int64, error)
}
