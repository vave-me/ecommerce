package domain

type MessageSent struct {
	ID             string
	ConversationID string
	SenderID       string
	RecipientID    string
	Body           string
	ItemID         string
	IsRead         bool
}

type MessageRead struct {
	MessageID string
	UserID    string
}

type MessageDeleted struct {
	MessageID string
	UserID    string
}

type MessengerSent struct {
	Text string
}

type MessageReceived struct {
	ID             string
	ConversationID string
	RecipientID    string
}

// Key implements registry.Registerable
