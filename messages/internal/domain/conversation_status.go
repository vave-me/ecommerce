package domain

// Enumeration of possible conversation statuses.
type ConversationStatus string

const (
	Unknown  ConversationStatus = "unknown"
	Pending  ConversationStatus = "pending"
	Active   ConversationStatus = "active"
	Aborted  ConversationStatus = "aborted"
	Archived ConversationStatus = "archived"
	Reopened ConversationStatus = "reopened"
)

func (s ConversationStatus) String() string {
	switch s {
	case Pending, Active, Aborted, Archived, Reopened:
		return string(s)
	default:
		return ""
	}
}

func ToConversationStatus(status string) ConversationStatus {
	switch status {
	case Pending.String():
		return Pending
	case Active.String():
		return Active
	case Aborted.String():
		return Aborted
	case Archived.String():
		return Archived
	case Reopened.String():
		return Reopened
	default:
		return Unknown
	}
}

// Enumeration of possible message statuses.
type MessageStatus string

const (
	Other    MessageStatus = "other"
	Sent     MessageStatus = "sent"
	Received MessageStatus = "received"
	Read     MessageStatus = "read"
	Answered MessageStatus = "answered"
	Deleted  MessageStatus = "deleted"
)

func (s MessageStatus) String() string {
	switch s {
	case Sent, Received, Read, Answered, Deleted, Other:
		return string(s)
	default:
		return ""
	}
}
func ToMessageStatus(status string) MessageStatus {
	switch status {
	case Sent.String():
		return Sent
	case Received.String():
		return Received
	case Read.String():
		return Read
	case Answered.String():
		return Answered
	case Deleted.String():
		return Deleted
	default:
		return Other
	}
}
