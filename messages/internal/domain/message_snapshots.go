package domain

type MessageV1 struct {
	ID             string
	ConversationID string
	ItemID         string
	SenderID       string
	RecipientID    string
	Body           string
	IsRead         bool
}

func (MessageV1) SnapshotName() string { return "messages.MessageV1" }
