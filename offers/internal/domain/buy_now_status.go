package domain // BuyNowStatus indicates the stage of a lease
type BuyNowStatus string

const (
	BuyNowStatusPending             BuyNowStatus = "pending"
	BuyNowStatusConfirmed           BuyNowStatus = "confirmed"
	BuyNowStatusCanceled            BuyNowStatus = "canceled"
	BuyNowStatusNegotiationPending  BuyNowStatus = "negotiation_pending"
	BuyNowStatusNegotiationAccepted BuyNowStatus = "negotiation_accepted"
	BuyNowStatusNegotiationDeclined BuyNowStatus = "negotiation_declined"
)
