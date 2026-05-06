package domain

import (
	"time"

	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

const BuyBackAggregate = "offers.BuyBack"

type BuyBack struct {
	es.Aggregate
	OfferID             string
	LockedPrice         int64
	RedemptionFee       int64
	LockBuyerID         string
	LockStartAt         time.Time
	LockExpiresAt       time.Time
	BuyBackStatus       BuyBackStatus
	NegotiatedPrice     int64
	NegotiationComments string
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*BuyBack)(nil)

func NewBuyBack(id string) *BuyBack {
	return &BuyBack{
		Aggregate: es.NewAggregate(id, BuyBackAggregate),
	}
}

func (BuyBack) Key() string { return BuyBackAggregate }

// CreateBuyBack
func (p *BuyBack) CreateBuyBack(
	offerID string,
	lockedPrice int64,
	redemptionFee int64,
	lockBuyerID string,
	lockDurationDays int,
) (ddd.Event, error) {

	if offerID == "" {
		return nil, errors.Wrap(errors.ErrBadRequest, "missing offerID")
	}
	if lockedPrice <= 0 {
		return nil, errors.Wrap(errors.ErrBadRequest, "locked price must be > 0")
	}
	lockStartAt := time.Now()
	lockExpiresAt := lockStartAt.AddDate(0, 0, lockDurationDays)

	p.AddEvent(BuyBackCreatedEvent, &BuyBackCreated{
		OfferID:       offerID,
		LockedPrice:   lockedPrice,
		RedemptionFee: redemptionFee,
		LockBuyerID:   lockBuyerID,
		LockStartAt:   lockStartAt,
		LockExpiresAt: lockExpiresAt,
		BuyBackStatus: BuyBackStatusActive,
	})
	return ddd.NewEvent(BuyBackCreatedEvent, p), nil
}

// RedeemBuyBack
func (p *BuyBack) RedeemBuyBack() (ddd.Event, error) {
	if p.BuyBackStatus != BuyBackStatusActive {
		return nil, errors.Wrap(errors.ErrConflict, "buyBack not active, cannot redeem")
	}
	p.AddEvent(BuyBackRedeemedEvent, &BuyBackRedeemed{
		BuyBackID: p.ID(),
	})
	return ddd.NewEvent(BuyBackRedeemedEvent, p), nil
}

// ExpireBuyBack
func (p *BuyBack) ExpireBuyBack() (ddd.Event, error) {
	if p.BuyBackStatus != BuyBackStatusActive {
		return nil, errors.Wrap(errors.ErrConflict, "buyBack not active, cannot expire")
	}
	p.AddEvent(BuyBackExpiredEvent, &BuyBackExpired{
		BuyBackID: p.ID(),
	})
	return ddd.NewEvent(BuyBackExpiredEvent, p), nil
}

// CancelBuyBack
func (p *BuyBack) CancelBuyBack() (ddd.Event, error) {
	if p.BuyBackStatus != BuyBackStatusActive {
		return nil, errors.Wrap(errors.ErrConflict, "buyBack not active or cannot cancel")
	}
	p.AddEvent(BuyBackCanceledEvent, &BuyBackCanceled{
		BuyBackID: p.ID(),
	})
	return ddd.NewEvent(BuyBackCanceledEvent, p), nil
}

// RequestNegotiation
func (p *BuyBack) RequestNegotiation(comments string) (ddd.Event, error) {
	// Possibly only allowed if status is "active"? Or "pending"? Up to your logic
	if p.BuyBackStatus != BuyBackStatusActive {
		return nil, errors.Wrap(errors.ErrConflict, "cannot negotiate now")
	}

	p.AddEvent(BuyBackNegotiationRequestedEvent, &BuyBackNegotiationRequested{
		BuyBackID: p.ID(),
		Comments:  comments,
	})
	return ddd.NewEvent(BuyBackNegotiationRequestedEvent, p), nil
}

// AcceptNegotiation
func (p *BuyBack) AcceptNegotiation(newPrice int64, userCustomerID string) (ddd.Event, error) {
	if p.BuyBackStatus != BuyBackStatusNegotiationPending {
		return nil, errors.Wrap(errors.ErrConflict, "no negotiation pending")
	}
	if newPrice <= 0 {
		return nil, errors.Wrap(errors.ErrBadRequest, "invalid negotiated price")
	}

	p.AddEvent(BuyBackNegotiationAcceptedEvent, &BuyBackNegotiationAccepted{
		NegotiatedPrice: newPrice,
		UserCustomerID:  userCustomerID,
	})
	return ddd.NewEvent(BuyBackNegotiationAcceptedEvent, p), nil
}

// DeclineNegotiation
func (p *BuyBack) DeclineNegotiation(reason string) (ddd.Event, error) {
	if p.BuyBackStatus != BuyBackStatusNegotiationPending {
		return nil, errors.Wrap(errors.ErrConflict, "no negotiation pending")
	}
	p.AddEvent(BuyBackNegotiationDeclinedEvent, &BuyBackNegotiationDeclined{
		BuyBackID: p.ID(),
		Reason:    reason,
	})
	return ddd.NewEvent(BuyBackNegotiationDeclinedEvent, p), nil
}

// ----------------------
// Apply Events
// ----------------------
func (p *BuyBack) ApplyEvent(event ddd.Event) error {
	switch e := event.Payload().(type) {

	case *BuyBackCreated:
		p.OfferID = e.OfferID
		p.LockedPrice = e.LockedPrice
		p.RedemptionFee = e.RedemptionFee
		p.LockBuyerID = e.LockBuyerID
		p.LockStartAt = e.LockStartAt
		p.LockExpiresAt = e.LockExpiresAt
		p.BuyBackStatus = e.BuyBackStatus

	case *BuyBackRedeemed:
		p.BuyBackStatus = BuyBackStatusRedeemed

	case *BuyBackExpired:
		p.BuyBackStatus = BuyBackStatusExpired

	case *BuyBackCanceled:
		p.BuyBackStatus = BuyBackStatusCanceled

	case *BuyBackNegotiationRequested:
		// Switch to negotiation pending
		p.BuyBackStatus = BuyBackStatusNegotiationPending
		p.NegotiationComments = e.Comments

	case *BuyBackNegotiationAccepted:
		p.BuyBackStatus = BuyBackStatusNegotiationAgreed
		p.NegotiatedPrice = e.NegotiatedPrice

	case *BuyBackNegotiationDeclined:
		p.BuyBackStatus = BuyBackStatusNegotiationDeclined
		// Possibly reverts to active or fully declines the entire buyBack? Up to your logic

	default:
		return errors.ErrInternal.
			Msgf("%T got unexpected event payload %T", p, e)
	}
	return nil
}

func (p *BuyBack) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *BuyBackV1:
		p.OfferID = ss.OfferID
		p.LockedPrice = ss.LockedPrice
		p.RedemptionFee = ss.RedemptionFee
		p.LockBuyerID = ss.LockBuyerID
		p.LockStartAt = ss.LockStartAt
		p.LockExpiresAt = ss.LockExpiresAt
		p.BuyBackStatus = ss.BuyBackStatus
		p.NegotiatedPrice = ss.NegotiatedPrice
		p.NegotiationComments = ss.NegotiationComments

	default:
		return errors.ErrInternal.Msgf("%T got unexpected snapshot %T", p, snapshot)
	}
	return nil
}

func (p BuyBack) ToSnapshot() es.Snapshot {
	return &BuyBackV1{
		OfferID:             p.OfferID,
		LockedPrice:         p.LockedPrice,
		RedemptionFee:       p.RedemptionFee,
		LockBuyerID:         p.LockBuyerID,
		LockStartAt:         p.LockStartAt,
		LockExpiresAt:       p.LockExpiresAt,
		BuyBackStatus:       p.BuyBackStatus,
		NegotiatedPrice:     p.NegotiatedPrice,
		NegotiationComments: p.NegotiationComments,
	}
}
