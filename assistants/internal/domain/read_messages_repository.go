package domain

import (
	"context"
	"time"
)

type ReadMessage struct {
	ID             string
	ConversationID string
	AssistantID    string
	Role           MessageRole
	Content        string
	Timestamp      time.Time
	Metadata       map[string]interface{}
	ActionsTaken   []AssistantAction
}
type ReadMessagesRepository interface {
	GetConversationMessages(ctx context.Context, conversationID string, userID string, limit, offset int) ([]*ReadMessage, error)
	GetLatestMessage(ctx context.Context, conversationID string, userID string) (*ReadMessage, error)
	GetUserMessageCount(ctx context.Context, userID string, dateRange string) (int64, error)
	AddMessage(ctx context.Context, conversationID, assistantID, id string, role MessageRole, content string, timestamp time.Time, metadata map[string]interface{}, ActionsTaken []AssistantAction) error
}
