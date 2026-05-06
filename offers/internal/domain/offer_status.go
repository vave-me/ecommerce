package domain

// Example OfferStatus
type OfferStatus string

const (
	OfferStatusDraft    OfferStatus = "draft"
	OfferStatusActive   OfferStatus = "active"
	OfferStatusAccepted OfferStatus = "accepted"
	OfferStatusBuyBack  OfferStatus = "buy-back"
	OfferStatusLeased   OfferStatus = "leased"
	OfferStatusSold     OfferStatus = "sold"
	OfferStatusClosed   OfferStatus = "closed"
	OfferStatusError    OfferStatus = "error"
)

func (s OfferStatus) String() string {
	switch s {
	case OfferStatusAccepted, OfferStatusDraft, OfferStatusActive, OfferStatusClosed:
		return string(s)
	default:
		return ""
	}
}
func ToOfferStatus(status string) OfferStatus {
	switch status {
	case OfferStatusActive.String():
		return OfferStatusActive
	case OfferStatusDraft.String():
		return OfferStatusDraft
	case OfferStatusAccepted.String():
		return OfferStatusAccepted
	case OfferStatusClosed.String():
		return OfferStatusClosed
	case OfferStatusBuyBack.String():
		return OfferStatusBuyBack
	case OfferStatusLeased.String():
		return OfferStatusLeased
	case OfferStatusSold.String():
		return OfferStatusSold

	default:
		return OfferStatusError
	}
}
