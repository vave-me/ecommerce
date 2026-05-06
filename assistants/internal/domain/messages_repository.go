package domain

import (
	"context"
	"middleman/assistants/internal/models"
)

type MessagesRepository interface {
	// Conversation Management
	CreateNewConversation(ctx context.Context, senderID, recipientID, itemID string) (*models.ConversationResponse, error)
	RestoreArchivedConversation(ctx context.Context, conversationID string) (*models.ConversationStatusResponse, error)
	ArchiveExistingConversation(ctx context.Context, conversationID string) (*models.ConversationStatusResponse, error)
	GetConversationByID(ctx context.Context, conversationID string) (*models.Conversation, error)
	FindConversationByRecipientAndItem(ctx context.Context, recipientID, itemID string) (*models.Conversation, error)
	GetUserConversationList(ctx context.Context, userID string, page, limit int64) (*models.ConversationsResponse, error)
	GetUserActiveConversations(ctx context.Context, userID string, page, limit int64) (*models.ConversationsResponse, error)

	// Message Management
	SendNewMessage(ctx context.Context, messageID, conversationID, senderID, recipientID, itemID, body string, isRead bool) (*models.MessageSentResponse, error)
	DeleteMessageByID(ctx context.Context, messageID string) error
	GetMessageByID(ctx context.Context, messageID string) (*models.Message, error)
	GetConversationMessages(ctx context.Context, conversationID string, page, limit int64) (*models.MessagesResponse, error)

	// Additional methods needed by tool service
	MarkMessageAsReadByUser(ctx context.Context, messageID, userID string) error
	GetUserUnreadMessageCount(ctx context.Context, userID string) (int64, error)
}
