package domain

import "time"

type ConversationMessage struct {
	ID             string
	ConversationID string
	ManagerID      string
	Role           MessageRole
	Content        string
	Timestamp      time.Time
	Metadata       map[string]interface{}
	ActionsTaken   []ManagerAction
}
