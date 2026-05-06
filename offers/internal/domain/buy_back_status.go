package domain // BuyBackStatus
type BuyBackStatus string

const (
	BuyBackStatusActive   BuyBackStatus = "active"
	BuyBackStatusRedeemed BuyBackStatus = "redeemed"
	BuyBackStatusExpired  BuyBackStatus = "expired"
	BuyBackStatusCanceled BuyBackStatus = "canceled"

	// Potential statuses for negotiation
	BuyBackStatusNegotiationPending  BuyBackStatus = "negotiation_pending"
	BuyBackStatusNegotiationAgreed   BuyBackStatus = "negotiation_agreed"
	BuyBackStatusNegotiationDeclined BuyBackStatus = "negotiation_declined"
)

func (s BuyBackStatus) String() string {
	switch s {
	case BuyBackStatusActive, BuyBackStatusRedeemed, BuyBackStatusExpired, BuyBackStatusCanceled, BuyBackStatusNegotiationPending, BuyBackStatusNegotiationAgreed, BuyBackStatusNegotiationDeclined:
		return string(s)
	default:
		return ""
	}
}
