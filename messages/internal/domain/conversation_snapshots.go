package domain

type ConversationV1 struct {
	SenderID    string
	RecipientID string
	ItemID      string
}

func (ConversationV1) SnapshotName() string { return "messages.ConversationV1" }
