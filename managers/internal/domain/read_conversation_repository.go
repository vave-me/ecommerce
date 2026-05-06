package domain

import (
	"context"
	"time"
)

type ReadConversation struct {
	ID        string
	UserID    string
	ManagerID string
	Messages  []ConversationMessage
	CreatedAt time.Time
	UpdatedAt time.Time
	Active    bool
	Context   map[string]interface{}
}

type ReadConversationRepository interface {
	// Read operations
	GetConversation(ctx context.Context, id string, userID string) (*ReadConversation, error)
	GetUserConversations(ctx context.Context, userID string, activeOnly bool, limit, offset int) ([]*ReadConversation, int, error)
	GetManagerConversations(ctx context.Context, managerID string, limit, offset int) ([]*ReadConversation, int, error)
	GetRecentConversations(ctx context.Context, userID string, limit int) ([]*ReadConversation, error)
	GetConversationsByDateRange(ctx context.Context, userID string, startDate, endDate time.Time) ([]*ReadConversation, error)

	// Write operations for read model maintenance
	AddConversation(ctx context.Context, id, userID, managerID string, createdAt time.Time, context map[string]interface{}) error
	UpdateConversationContext(ctx context.Context, id string, context map[string]interface{}, updatedAt time.Time) error
	ArchiveConversation(ctx context.Context, id string, archivedAt time.Time) error
	UpdateConversationTimestamp(ctx context.Context, id string, updatedAt time.Time) error
}
