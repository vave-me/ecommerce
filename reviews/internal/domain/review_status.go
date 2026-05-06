package domain

// Enumeration of possible review statuses.
type ReviewStatus string

const (
	Unknown  ReviewStatus = "unknown"
	Pending  ReviewStatus = "pending"
	Approved ReviewStatus = "approved"
	Rejected ReviewStatus = "rejected"
)

func (s ReviewStatus) String() string {
	switch s {
	case Pending, Approved, Rejected:
		return string(s)
	default:
		return ""
	}
}
func ToReviewStatus(status string) ReviewStatus {
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
