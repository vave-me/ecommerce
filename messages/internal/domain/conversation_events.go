package domain

type ConversationStarted struct {
	ConversationID string
	SenderID       string
	RecipientID    string
	ItemID         string
}

type ConversationArchived struct {
	ConversationID string
	ArchivedBy     string
}

type ConversationDeleted struct {
	ConversationID string
	ParticipantID  string
	RemovedBy      string
}

// Key implements registry.Registerable

type ConversationMutted struct {
	ConversationID string
	MutedBy        string
}

type ConversationUnmutted struct {
	ConversationID string
	UnmutedBy      string
}

type ConversationClosed struct {
}

type ConversationReopened struct {
}
