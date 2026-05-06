package domain

// Enumeration of possible comment statuses.
type CommentStatus string

const (
	Unknown  CommentStatus = "unknown"
	Pending  CommentStatus = "pending"
	Approved CommentStatus = "approved"
	Rejected CommentStatus = "rejected"
)

func (s CommentStatus) String() string {
	switch s {
	case Pending, Approved, Rejected:
		return string(s)
	default:
		return ""
	}
}
func ToCommentStatus(status string) CommentStatus {
	switch status {
	case Pending.String():
		return Pending
	case Approved.String():
		return Approved
	case Rejected.String():
		return Rejected
	default:
		return Unknown
	}
}
