package domain

import (
	"time"

	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

// ReservationAggregate is the name of this aggregator.
const ReservationAggregate = "offers.Reservation"

// Potential domain errors.
var (
	ErrMissingOfferID          = errors.Wrap(errors.ErrBadRequest, "missing offerID")
	ErrLockedPriceInvalid      = errors.Wrap(errors.ErrBadRequest, "locked price must be > 0")
	ErrReservationNotActive    = errors.Wrap(errors.ErrConflict, "reservation not active, cannot proceed")
	ErrReservationAlreadyUsed  = errors.Wrap(errors.ErrConflict, "free reservation already used for this offer")
	ErrUnknownStatusTransition = errors.Wrap(errors.ErrConflict, "invalid reservation status transition")
)

type Reservation struct {
	es.Aggregate

	// Main fields
	OfferID             string
	LockBuyerID         string
	LockStartAt         time.Time
	LockExpiresAt       time.Time
	ReservationStatus   ReservationStatus
	LockedPrice         int64
	ReservationFee      int64
	IsFree              bool
	FreeReservationUsed bool
	NegotiationComments string
	NegotiatedPrice     int64
}

func (Reservation) Key() string {
	return ReservationAggregate
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Reservation)(nil)

func NewReservation(id string) *Reservation {
	return &Reservation{
		Aggregate: es.NewAggregate(id, ReservationAggregate),
		// Initialize default status if needed
	}
}

func (r *Reservation) CreateReservation(
	offerID string,
	lockedPrice int64,
	reservationFee int64,
	lockBuyerID string,
	lockDurationDays int,
) (ddd.Event, error) {

	if offerID == "" {
		return nil, ErrMissingOfferID
	}
	if lockedPrice <= 0 {
		return nil, ErrLockedPriceInvalid
	}

	isFree := (lockDurationDays == 1)
	if isFree && r.FreeReservationUsed {
		// The business rule: "1-day free reservation can happen only once per offer"
		return nil, ErrReservationAlreadyUsed
	}

	// Set fee = 0 if free, otherwise the "reservationFee" passed in
	finalReservationFee := int64(0)
	if !isFree {
		finalReservationFee = reservationFee
	}

	lockStartAt := time.Now()
	lockExpiresAt := lockStartAt.AddDate(0, 0, lockDurationDays)

	r.AddEvent(ReservationCreatedEvent, &ReservationCreated{
		OfferID:           offerID,
		LockedPrice:       lockedPrice,
		ReservationFee:    finalReservationFee,
		LockBuyerID:       lockBuyerID,
		LockStartAt:       lockStartAt,
		LockExpiresAt:     lockExpiresAt,
		ReservationStatus: ReservationStatusActive,
		IsFree:            isFree,
	})

	return ddd.NewEvent(ReservationCreatedEvent, r), nil
}

// RedeemReservation => user chooses to finalize/buy the item?
func (r *Reservation) RedeemReservation() (ddd.Event, error) {
	if r.ReservationStatus != ReservationStatusActive {
		return nil, ErrReservationNotActive
	}
	r.AddEvent(ReservationRedeemedEvent, &ReservationRedeemed{
		ReservationID: r.ID(),
	})
	return ddd.NewEvent(ReservationRedeemedEvent, r), nil
}

// ExpireReservation => the free or paid reservation expires after lockExpiresAt
func (r *Reservation) ExpireReservation() (ddd.Event, error) {
	if r.ReservationStatus != ReservationStatusActive {
		return nil, ErrReservationNotActive
	}
	r.AddEvent(ReservationExpiredEvent, &ReservationExpired{
		ReservationID: r.ID(),
	})
	return ddd.NewEvent(ReservationExpiredEvent, r), nil
}

// CancelReservation => user or system cancels reservation before it is redeemed/expired
func (r *Reservation) CancelReservation() (ddd.Event, error) {
	if r.ReservationStatus != ReservationStatusActive {
		return nil, errors.Wrap(errors.ErrConflict, "reservation not active or cannot cancel")
	}
	r.AddEvent(ReservationCanceledEvent, &ReservationCanceled{
		ReservationID: r.ID(),
	})
	return ddd.NewEvent(ReservationCanceledEvent, r), nil
}

// RequestNegotiation => buyer requests a new negotiated price
func (r *Reservation) RequestNegotiation(comments string) (ddd.Event, error) {
	// For example, only allow if status is Active
	if r.ReservationStatus != ReservationStatusActive {
		return nil, errors.Wrap(errors.ErrConflict, "cannot negotiate in current status")
	}
	r.AddEvent(ReservationNegotiationRequestedEvent, &ReservationNegotiationRequested{
		ReservationID: r.ID(),
		Comments:      comments,
	})
	return ddd.NewEvent(ReservationNegotiationRequestedEvent, r), nil
}

// AcceptNegotiation => the seller or system accepts the negotiation
func (r *Reservation) AcceptNegotiation(newPrice int64, userCustomerID string) (ddd.Event, error) {
	if r.ReservationStatus != ReservationStatusNegotiationPending {
		return nil, errors.Wrap(errors.ErrConflict, "no negotiation pending to accept")
	}
	if newPrice <= 0 {
		return nil, errors.Wrap(errors.ErrBadRequest, "invalid negotiated price")
	}
	r.AddEvent(ReservationNegotiationAcceptedEvent, &ReservationNegotiationAccepted{
		ReservationID:   r.ID(),
		NegotiatedPrice: newPrice,
		UserCustomerID:  userCustomerID,
	})
	return ddd.NewEvent(ReservationNegotiationAcceptedEvent, r), nil
}

// DeclineNegotiation => negotiation is refused
func (r *Reservation) DeclineNegotiation(reason string) (ddd.Event, error) {
	if r.ReservationStatus != ReservationStatusNegotiationPending {
		return nil, errors.Wrap(errors.ErrConflict, "no negotiation pending to decline")
	}
	r.AddEvent(ReservationNegotiationDeclinedEvent, &ReservationNegotiationDeclined{
		ReservationID: r.ID(),
		Reason:        reason,
	})
	return ddd.NewEvent(ReservationNegotiationDeclinedEvent, r), nil
}

func (r *Reservation) ApplyEvent(event ddd.Event) error {
	switch e := event.Payload().(type) {
	case *ReservationCreated:
		r.OfferID = e.OfferID
		r.LockedPrice = e.LockedPrice
		r.ReservationFee = e.ReservationFee
		r.LockBuyerID = e.LockBuyerID
		r.LockStartAt = e.LockStartAt
		r.LockExpiresAt = e.LockExpiresAt
		r.ReservationStatus = e.ReservationStatus
		r.IsFree = e.IsFree
		// If it's free => we mark that the free reservation was used (for that aggregator).
		if e.IsFree {
			r.FreeReservationUsed = true
		}

	case *ReservationRedeemed:
		r.ReservationStatus = ReservationStatusRedeemed

	case *ReservationExpired:
		r.ReservationStatus = ReservationStatusExpired

	case *ReservationCanceled:
		r.ReservationStatus = ReservationStatusCanceled

	case *ReservationNegotiationRequested:
		r.ReservationStatus = ReservationStatusNegotiationPending
		r.NegotiationComments = e.Comments

	case *ReservationNegotiationAccepted:
		r.ReservationStatus = ReservationStatusNegotiationAgreed
		r.NegotiatedPrice = e.NegotiatedPrice

	case *ReservationNegotiationDeclined:
		r.ReservationStatus = ReservationStatusNegotiationDeclined
		// Possibly revert to Active or something else?

	default:
		return errors.ErrInternal.Msgf("%T got unexpected event payload %T", r, e)
	}
	return nil
}

func (r *Reservation) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *ReservationV1:
		r.OfferID = ss.OfferID
		r.LockedPrice = ss.LockedPrice
		r.ReservationFee = ss.ReservationFee
		r.LockBuyerID = ss.LockBuyerID
		r.LockStartAt = ss.LockStartAt
		r.LockExpiresAt = ss.LockExpiresAt
		r.ReservationStatus = ss.ReservationStatus
		r.IsFree = ss.IsFree
		r.FreeReservationUsed = ss.FreeReservationUsed
		r.NegotiationComments = ss.NegotiationComments
		r.NegotiatedPrice = ss.NegotiatedPrice

	default:
		return errors.ErrInternal.
			Msgf("%T got unexpected snapshot %T", r, snapshot)
	}
	return nil
}

func (r *Reservation) ToSnapshot() es.Snapshot {
	return &ReservationV1{
		OfferID:             r.OfferID,
		LockedPrice:         r.LockedPrice,
		ReservationFee:      r.ReservationFee,
		LockBuyerID:         r.LockBuyerID,
		LockStartAt:         r.LockStartAt,
		LockExpiresAt:       r.LockExpiresAt,
		ReservationStatus:   r.ReservationStatus,
		IsFree:              r.IsFree,
		FreeReservationUsed: r.FreeReservationUsed,
		NegotiationComments: r.NegotiationComments,
		NegotiatedPrice:     r.NegotiatedPrice,
	}
}
