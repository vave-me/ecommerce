package domain

import "context"

type ConversationRepository interface {
	Load(ctx context.Context, conversationID string) (*Conversation, error)
	Save(ctx context.Context, conversation *Conversation) error
}
