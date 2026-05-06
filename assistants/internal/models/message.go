package models

import "time"

// Message represents an individual message within a conversation
type Message struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	SenderID       string    `json:"sender_id"`
	RecipientID    string    `json:"recipient_id"`
	ItemID         string    `json:"item_id"`
	Body           string    `json:"body"`
	IsRead         bool      `json:"is_read"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
}

// Conversation represents a dialogue between users regarding a specific item
type Conversation struct {
	ID                 string    `json:"id"`
	SenderID           string    `json:"sender_id"`
	RecipientID        string    `json:"recipient_id"`
	ItemID             string    `json:"item_id"`
	ConversationStatus string    `json:"conversation_status"`
	Active             bool      `json:"active,omitempty"` // For backward compatibility
	CreatedAt          time.Time `json:"created_at,omitempty"`
	UpdatedAt          time.Time `json:"updated_at,omitempty"`
}

// Response models matching protobuf specifications

// ConversationResponse for StartConversation
type ConversationResponse struct {
	ID string `json:"id"`
}

// ConversationStatusResponse for conversation status operations
type ConversationStatusResponse struct {
	ID                 string `json:"id"`
	ConversationStatus string `json:"conversation_status"`
}

// ConversationsResponse for listing conversations
type ConversationsResponse struct {
	Conversations []*Conversation `json:"conversations"`
	Total         int64           `json:"total"`
	Page          int64           `json:"page"`
	Limit         int64           `json:"limit"`
}

// MessageSentResponse for SendMessage
type MessageSentResponse struct {
	ID     string    `json:"id"`
	SentAt time.Time `json:"sent_at"`
}

// MessagesResponse for listing messages
type MessagesResponse struct {
	Messages []*Message `json:"messages"`
	Total    int64      `json:"total"`
	Page     int64      `json:"page"`
	Limit    int64      `json:"limit"`
}

// ErrorResponse provides a standardized error structure
type ErrorResponse struct {
	Code    int64    `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details"`
}
