package domain

type SupportStatus string

const (
	SupportUnknown    SupportStatus = ""
	SupportIsOpen     SupportStatus = "open"
	SupportIsResolved SupportStatus = "resolved"
	SupportIsClosed   SupportStatus = "closed"
)

func (s SupportStatus) String() string {
	switch s {
	case SupportIsOpen, SupportIsResolved, SupportIsClosed:
		return string(s)
	default:
		return ""
	}
}

func ToSupportStatus(status string) SupportStatus {
	switch status {
	case SupportIsOpen.String():
		return SupportIsOpen
	case SupportIsResolved.String():
		return SupportIsResolved
	case SupportIsClosed.String():
		return SupportIsClosed
	default:
		return SupportUnknown
	}
}
