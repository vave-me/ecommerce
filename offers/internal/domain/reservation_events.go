package domain

import "time"

// Event constants
const (
	ReservationCreatedEvent              = "offers.ReservationCreated"
	ReservationRedeemedEvent             = "offers.ReservationRedeemed"
	ReservationExpiredEvent              = "offers.ReservationExpired"
	ReservationCanceledEvent             = "offers.ReservationCanceled"
	ReservationNegotiationRequestedEvent = "offers.ReservationNegotiationRequested"
	ReservationNegotiationAcceptedEvent  = "offers.ReservationNegotiationAccepted"
	ReservationNegotiationDeclinedEvent  = "offers.ReservationNegotiationDeclined"
)

// ReservationCreated is fired when a new reservation is created.
type ReservationCreated struct {
	OfferID           string
	LockedPrice       int64
	ReservationFee    int64
	LockBuyerID       string
	LockStartAt       time.Time
	LockExpiresAt     time.Time
	ReservationStatus ReservationStatus
	IsFree            bool
}

func (ReservationCreated) Key() string { return ReservationCreatedEvent }

// ReservationRedeemed
type ReservationRedeemed struct {
	ReservationID string
}

func (ReservationRedeemed) Key() string { return ReservationRedeemedEvent }

// ReservationExpired
type ReservationExpired struct {
	ReservationID string
}

func (ReservationExpired) Key() string { return ReservationExpiredEvent }

// ReservationCanceled
type ReservationCanceled struct {
	ReservationID string
}

func (ReservationCanceled) Key() string { return ReservationCanceledEvent }

// ReservationNegotiationRequested
type ReservationNegotiationRequested struct {
	ReservationID string
	Comments      string
}

func (ReservationNegotiationRequested) Key() string { return ReservationNegotiationRequestedEvent }

// ReservationNegotiationAccepted
type ReservationNegotiationAccepted struct {
	ReservationID   string
	NegotiatedPrice int64
	UserCustomerID  string
}

func (ReservationNegotiationAccepted) Key() string { return ReservationNegotiationAcceptedEvent }

// ReservationNegotiationDeclined
type ReservationNegotiationDeclined struct {
	ReservationID string
	Reason        string
}

func (ReservationNegotiationDeclined) Key() string { return ReservationNegotiationDeclinedEvent }
