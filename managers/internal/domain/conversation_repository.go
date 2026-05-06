package domain

import (
	"context"
)

// ConversationRepository defines the interface for conversation persistence (Event Sourcing)
// This repository handles only write operations for the conversation aggregate
type ConversationRepository interface {
	Load(ctx context.Context, id string) (*Conversation, error)
	Save(ctx context.Context, conversation *Conversation) error
}
