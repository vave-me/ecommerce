package domain

import (
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

// The aggregate name
const BuyNowAggregate = "offers.BuyNow"

type BuyNow struct {
	es.Aggregate

	OfferID    string
	FinalPrice int64
	Status     BuyNowStatus

	// Extra fields to track negotiation
	NegotiatedPrice     int64
	NegotiationComments string
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*BuyNow)(nil)

func NewBuyNow(id string) *BuyNow {
	return &BuyNow{
		Aggregate: es.NewAggregate(id, BuyNowAggregate),
	}
}

func (BuyNow) Key() string { return BuyNowAggregate }

//--------------------
// Commands on BuyNow
//--------------------

// CreateBuyNow
func (b *BuyNow) CreateBuyNow(offerID string, finalPrice int64) (ddd.Event, error) {
	if offerID == "" {
		return nil, errors.Wrap(errors.ErrBadRequest, "missing offerID")
	}
	if finalPrice <= 0 {
		return nil, errors.Wrap(errors.ErrBadRequest, "invalid price")
	}

	b.AddEvent(BuyNowCreatedEvent, &BuyNowCreated{
		OfferID:    offerID,
		FinalPrice: finalPrice,
		Status:     BuyNowStatusPending,
	})
	return ddd.NewEvent(BuyNowCreatedEvent, b), nil
}

// ConfirmBuyNow
func (b *BuyNow) ConfirmBuyNow() (ddd.Event, error) {
	if b.Status != BuyNowStatusPending {
		return nil, errors.Wrap(errors.ErrConflict, "BuyNow not pending")
	}
	b.AddEvent(BuyNowConfirmedEvent, &BuyNowConfirmed{
		BuyNowID: b.ID(),
	})
	return ddd.NewEvent(BuyNowConfirmedEvent, b), nil
}

// CancelBuyNow
func (b *BuyNow) CancelBuyNow() (ddd.Event, error) {
	if b.Status != BuyNowStatusPending {
		return nil, errors.Wrap(errors.ErrConflict, "BuyNow not pending or already done")
	}
	b.AddEvent(BuyNowCanceledEvent, &BuyNowCanceled{
		BuyNowID: b.ID(),
	})
	return ddd.NewEvent(BuyNowCanceledEvent, b), nil
}

// RequestNegotiation: The buyer wants a price negotiation.
func (b *BuyNow) RequestNegotiation(comments string) (ddd.Event, error) {
	// Possibly only allowed if current status is "pending" or "canceled"?
	if b.Status != BuyNowStatusPending {
		return nil, errors.Wrap(errors.ErrConflict, "BuyNow cannot be negotiated in current status")
	}
	b.AddEvent(BuyNowNegotiationRequestedEvent, &BuyNowNegotiationRequested{
		BuyNowID: b.ID(),
		Comments: comments,
	})
	return ddd.NewEvent(BuyNowNegotiationRequestedEvent, b), nil
}

// AcceptNegotiation: The seller or system accepts the new negotiated price
func (b *BuyNow) AcceptNegotiation(newPrice int64, userCustomerID string) (ddd.Event, error) {
	if b.Status != BuyNowStatusNegotiationPending {
		return nil, errors.Wrap(errors.ErrConflict, "no negotiation pending or mismatch status")
	}
	if newPrice <= 0 {
		return nil, errors.Wrap(errors.ErrBadRequest, "invalid negotiated price")
	}

	b.AddEvent(BuyNowNegotiationAcceptedEvent, &BuyNowNegotiationAccepted{
		BuyNowID:        b.ID(),
		NegotiatedPrice: newPrice,
		UserCustomerID:  userCustomerID,
	})
	return ddd.NewEvent(BuyNowNegotiationAcceptedEvent, b), nil
}

// DeclineNegotiation: The seller or system declines the negotiation
func (b *BuyNow) DeclineNegotiation(reason string) (ddd.Event, error) {
	if b.Status != BuyNowStatusNegotiationPending {
		return nil, errors.Wrap(errors.ErrConflict, "no negotiation pending or mismatch status")
	}
	b.AddEvent(BuyNowNegotiationDeclinedEvent, &BuyNowNegotiationDeclined{
		BuyNowID: b.ID(),
		Reason:   reason,
	})
	return ddd.NewEvent(BuyNowNegotiationDeclinedEvent, b), nil
}

func (b *BuyNow) ApplyEvent(event ddd.Event) error {
	switch e := event.Payload().(type) {

	case *BuyNowCreated:
		b.OfferID = e.OfferID
		b.FinalPrice = e.FinalPrice
		b.Status = e.Status

	case *BuyNowConfirmed:
		b.Status = BuyNowStatusConfirmed

	case *BuyNowCanceled:
		b.Status = BuyNowStatusCanceled

	case *BuyNowNegotiationRequested:
		b.Status = BuyNowStatusNegotiationPending
		b.NegotiationComments = e.Comments

	case *BuyNowNegotiationAccepted:
		b.Status = BuyNowStatusNegotiationAccepted
		b.NegotiatedPrice = e.NegotiatedPrice
		// Possibly update b.FinalPrice = e.NegotiatedPrice if that is your logic.

	case *BuyNowNegotiationDeclined:
		b.Status = BuyNowStatusNegotiationDeclined
		// Possibly revert to pending or fully close? up to you.

	default:
		return errors.ErrInternal.Msgf("%T got unexpected payload %T", b, e)
	}
	return nil
}

func (b *BuyNow) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *BuyNowV1:
		b.OfferID = ss.OfferID
		b.FinalPrice = ss.FinalPrice
		b.Status = ss.Status
		b.NegotiatedPrice = ss.NegotiatedPrice
		b.NegotiationComments = ss.NegotiationComments

	default:
		return errors.ErrInternal.Msgf("%T got unexpected snapshot %T", b, snapshot)
	}
	return nil
}

func (b BuyNow) ToSnapshot() es.Snapshot {
	return &BuyNowV1{
		OfferID:             b.OfferID,
		FinalPrice:          b.FinalPrice,
		Status:              b.Status,
		NegotiatedPrice:     b.NegotiatedPrice,
		NegotiationComments: b.NegotiationComments,
	}
}
