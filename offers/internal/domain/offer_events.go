package domain

// ---------------------
// EVENTS
// ---------------------

const (
	OfferCreatedEvent   = "offers.OfferCreated"
	OfferActivatedEvent = "offers.OfferActivated"
	OfferAcceptedEvent  = "offers.OfferAccepted"
	OfferClosedEvent    = "offers.OfferClosed"
)

// OfferCreated payload
type OfferCreated struct {
	UserSellerID string
	ProductID    string
	Price        int64
	Status       OfferStatus
}

func (OfferCreated) Key() string { return OfferCreatedEvent }

// OfferActivated payload
type OfferActivated struct {
	UserSellerID string
	ProductID    string
}

func (OfferActivated) Key() string { return OfferActivatedEvent }

// OfferAccepted payload
type OfferAccepted struct {
	OfferID        string
	UserCustomerID string
}

func (OfferAccepted) Key() string { return OfferAcceptedEvent }

// OfferClosed payload
type OfferClosed struct {
	OfferID string
	Reason  string
}

func (OfferClosed) Key() string { return OfferClosedEvent }
