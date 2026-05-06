package domain

import (
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

const OfferAggregate = "offers.Offer"

type Offer struct {
	es.Aggregate
	UserSellerID   string
	UserCustomerID string // optional if you track who eventually picks the offer
	ProductID      string
	Price          int64
	Status         OfferStatus
}

// Compile-time check
var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Offer)(nil)

func NewOffer(id string) *Offer {
	return &Offer{
		Aggregate: es.NewAggregate(id, OfferAggregate),
	}
}

func (Offer) Key() string { return OfferAggregate }

// Command: CreateOffer
func (o *Offer) CreateOffer(userSellerID, productID string, price int64) (ddd.Event, error) {
	// Basic validation
	if userSellerID == "" || productID == "" {
		return nil, errors.Wrap(errors.ErrBadRequest, "missing fields for creating offer")
	}
	if price <= 0 {
		return nil, errors.Wrap(errors.ErrBadRequest, "price must be greater than 0")
	}

	// Raise event
	o.AddEvent(OfferCreatedEvent, &OfferCreated{
		UserSellerID: userSellerID,
		ProductID:    productID,
		Price:        price,
		Status:       OfferStatusDraft,
	})
	return ddd.NewEvent(OfferCreatedEvent, o), nil
}

// Command: ActivateOffer
func (o *Offer) ActivateOffer(userID, productID string) (ddd.Event, error) {
	if o.Status != OfferStatusDraft {
		return nil, errors.Wrap(errors.ErrBadRequest, "only draft offers can be activated")
	}
	o.AddEvent(OfferActivatedEvent, &OfferActivated{
		UserSellerID: userID,
		ProductID:    o.ProductID,
	})
	return ddd.NewEvent(OfferActivatedEvent, o), nil
}

// Command: AcceptOffer
func (o *Offer) AcceptOffer(offerID string, userCustomerID string) (ddd.Event, error) {
	if o.Status != OfferStatusActive {
		return nil, errors.Wrap(errors.ErrConflict, "offer not in active state")
	}
	if userCustomerID == "" {
		return nil, errors.Wrap(errors.ErrBadRequest, "need a buyer ID to accept offer")
	}
	o.AddEvent(OfferAcceptedEvent, &OfferAccepted{
		OfferID:        offerID,
		UserCustomerID: userCustomerID,
	})
	return ddd.NewEvent(OfferAcceptedEvent, o), nil
}

// Command: CloseOffer
func (o *Offer) CloseOffer(offerID string, reason string) (ddd.Event, error) {
	if o.Status == OfferStatusClosed {
		return nil, errors.Wrap(errors.ErrBadRequest, "offer already closed")
	}
	o.AddEvent(OfferClosedEvent, &OfferClosed{
		OfferID: offerID,
		Reason:  reason,
	})
	return ddd.NewEvent(OfferClosedEvent, o), nil
}

func (o *Offer) ApplyEvent(event ddd.Event) error {
	switch e := event.Payload().(type) {

	case *OfferCreated:
		o.UserSellerID = e.UserSellerID
		o.ProductID = e.ProductID
		o.Price = e.Price
		o.Status = e.Status

	case *OfferActivated:
		o.Status = OfferStatusActive

	case *OfferAccepted:
		o.UserCustomerID = e.UserCustomerID
		o.Status = OfferStatusAccepted

	case *OfferClosed:
		o.Status = OfferStatusClosed

	default:
		return errors.ErrInternal.Msgf("%T got unexpected payload %T", o, e)
	}
	return nil
}

func (o *Offer) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *OfferV1:
		o.UserSellerID = ss.UserSellerID
		o.UserCustomerID = ss.UserCustomerID
		o.ProductID = ss.ProductID
		o.Price = ss.Price
		o.Status = ss.Status
	default:
		return errors.ErrInternal.Msgf("%T unexpected snapshot %T", o, snapshot)
	}
	return nil
}

func (o Offer) ToSnapshot() es.Snapshot {
	return &OfferV1{
		UserSellerID:   o.UserSellerID,
		UserCustomerID: o.UserCustomerID,
		ProductID:      o.ProductID,
		Price:          o.Price,
		Status:         o.Status,
	}
}
