package domain

// ---------------
// EVENTS
// ---------------

const (
	BuyNowCreatedEvent              = "offers.BuyNowCreated"
	BuyNowConfirmedEvent            = "offers.BuyNowConfirmed"
	BuyNowCanceledEvent             = "offers.BuyNowCanceled"
	BuyNowNegotiationRequestedEvent = "offers.BuyNowNegotiationRequested"
	BuyNowNegotiationAcceptedEvent  = "offers.BuyNowNegotiationAccepted"
	BuyNowNegotiationDeclinedEvent  = "offers.BuyNowNegotiationDeclined"
)

type BuyNowCreated struct {
	OfferID    string
	FinalPrice int64
	Status     BuyNowStatus
}

func (BuyNowCreated) Key() string { return BuyNowCreatedEvent }

type BuyNowConfirmed struct {
	BuyNowID string
}

func (BuyNowConfirmed) Key() string { return BuyNowConfirmedEvent }

type BuyNowCanceled struct {
	BuyNowID string
}

func (BuyNowCanceled) Key() string { return BuyNowCanceledEvent }

// Negotiation events
type BuyNowNegotiationRequested struct {
	BuyNowID string
	Comments string
}

func (BuyNowNegotiationRequested) Key() string { return BuyNowNegotiationRequestedEvent }

type BuyNowNegotiationAccepted struct {
	BuyNowID        string
	NegotiatedPrice int64
	UserCustomerID  string
}

func (BuyNowNegotiationAccepted) Key() string { return BuyNowNegotiationAcceptedEvent }

type BuyNowNegotiationDeclined struct {
	BuyNowID string
	Reason   string
}

func (BuyNowNegotiationDeclined) Key() string { return BuyNowNegotiationDeclinedEvent }
