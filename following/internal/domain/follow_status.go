package domain

// Enumeration of possible follow statuses.
type FollowStatus string

const (
	Unknown  FollowStatus = "unknown"
	Pending  FollowStatus = "pending"
	Approved FollowStatus = "approved"
	Rejected FollowStatus = "rejected"
)

func (s FollowStatus) String() string {
	switch s {
	case Pending, Approved, Rejected:
		return string(s)
	default:
		return ""
	}
}
func ToFollowStatus(status string) FollowStatus {
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
