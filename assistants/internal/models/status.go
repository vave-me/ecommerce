package models

type Status string

const (
	StatusActive    Status = "active"
	StatusLocked    Status = "locked"
	StatusSold      Status = "sold"
	StatusExpired   Status = "expired"
	StatusPaused    Status = "paused"
	StatusRent      Status = "rent"
	StatusDraft     Status = "draft"
	StatusArchived  Status = "archived"
	StatusReference Status = "reference"
	StatusUnknown   Status = ""
)

func (s Status) String() string {
	switch s {
	case StatusActive, StatusLocked, StatusSold, StatusExpired, StatusRent, StatusDraft, StatusPaused, StatusArchived, StatusReference:
		return string(s)
	default:
		return ""
	}
}

func ToStatus(s string) Status {
	switch s {
	case StatusActive.String():
		return StatusActive
	case StatusLocked.String():
		return StatusLocked
	case StatusSold.String():
		return StatusSold
	case StatusExpired.String():
		return StatusExpired
	case StatusPaused.String():
		return StatusPaused
	case StatusRent.String():
		return StatusRent
	case StatusArchived.String():
		return StatusArchived
	case StatusDraft.String():
		return StatusDraft
	case StatusReference.String():
		return StatusReference
	default:
		return StatusUnknown
	}
}
