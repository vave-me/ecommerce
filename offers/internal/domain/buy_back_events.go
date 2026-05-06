package domain

import "time"

// ---------------------
// EVENTS
// ---------------------

const (
	BuyBackCreatedEvent              = "offers.BuyBackCreated"
	BuyBackRedeemedEvent             = "offers.BuyBackRedeemed"
	BuyBackExpiredEvent              = "offers.BuyBackExpired"
	BuyBackCanceledEvent             = "offers.BuyBackCanceled"
	BuyBackNegotiationRequestedEvent = "offers.BuyBackNegotiationRequested"
	BuyBackNegotiationAcceptedEvent  = "offers.BuyBackNegotiationAccepted"
	BuyBackNegotiationDeclinedEvent  = "offers.BuyBackNegotiationDeclined"
)

type BuyBackCreated struct {
	OfferID       string
	LockedPrice   int64
	RedemptionFee int64
	LockBuyerID   string
	LockStartAt   time.Time
	LockExpiresAt time.Time
	BuyBackStatus BuyBackStatus
}

func (BuyBackCreated) Key() string { return BuyBackCreatedEvent }

type BuyBackRedeemed struct {
	BuyBackID string
}

func (BuyBackRedeemed) Key() string { return BuyBackRedeemedEvent }

type BuyBackExpired struct {
	BuyBackID string
}

func (BuyBackExpired) Key() string { return BuyBackExpiredEvent }

type BuyBackCanceled struct {
	BuyBackID string
}

func (BuyBackCanceled) Key() string { return BuyBackCanceledEvent }

type BuyBackNegotiationRequested struct {
	BuyBackID string
	Comments  string
}

func (BuyBackNegotiationRequested) Key() string { return BuyBackNegotiationRequestedEvent }

type BuyBackNegotiationAccepted struct {
	BuyBackID       string
	NegotiatedPrice int64
	UserCustomerID  string
}

func (BuyBackNegotiationAccepted) Key() string { return BuyBackNegotiationAcceptedEvent }

type BuyBackNegotiationDeclined struct {
	BuyBackID string
	Reason    string
}

func (BuyBackNegotiationDeclined) Key() string { return BuyBackNegotiationDeclinedEvent }
