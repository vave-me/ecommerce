package domain

import (
	"time"

	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

const LeaseAggregate = "offers.Lease"

type Lease struct {
	es.Aggregate
	OfferID         string
	MonthlyPrice    int64
	LeaseTermMonths int64
	HasBuyout       bool
	BuyoutPrice     int64
	LeaseStartDate  time.Time
	LeaseEndDate    time.Time
	LeaseStatus     LeaseStatus
	SchufaRequired  bool
}

// Compile-time checks
var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Lease)(nil)

func NewLease(id string) *Lease {
	return &Lease{
		Aggregate: es.NewAggregate(id, LeaseAggregate),
	}
}

func (Lease) Key() string { return LeaseAggregate }

func (l *Lease) CreateLease(
	offerID string,
	monthlyPrice int64,
	leaseTermMonths int64,
	hasBuyout bool,
	buyoutPrice int64,
	startDate time.Time,
	endDate time.Time,
	status LeaseStatus,
) (ddd.Event, error) {

	if offerID == "" {
		return nil, errors.Wrap(errors.ErrBadRequest, "missing offerID")
	}
	if monthlyPrice <= 0 {
		return nil, errors.Wrap(errors.ErrBadRequest, "monthly price must be > 0")
	}
	if status == "" {
		status = LeaseStatusPending
	}

	l.AddEvent(LeaseCreatedEvent, &LeaseCreated{
		OfferID:         offerID,
		MonthlyPrice:    monthlyPrice,
		LeaseTermMonths: leaseTermMonths,
		HasBuyout:       hasBuyout,
		BuyoutPrice:     buyoutPrice,
		LeaseStatus:     status,
	})
	return ddd.NewEvent(LeaseCreatedEvent, l), nil
}

func (l *Lease) StartLease(leaseID string, startDate time.Time) (ddd.Event, error) {
	if l.LeaseStatus != LeaseStatusPending && l.LeaseStatus != LeaseStatusNegotiating {
		return nil, errors.Wrap(errors.ErrConflict, "lease must be pending or negotiating to start")
	}
	endDate := startDate.AddDate(0, int(l.LeaseTermMonths), 0)

	l.AddEvent(LeaseStartedEvent, &LeaseStarted{
		LeaseID:        leaseID,
		LeaseStartDate: startDate,
		LeaseEndDate:   endDate,
	})
	return ddd.NewEvent(LeaseStartedEvent, l), nil
}

func (l *Lease) MakeLeasePayment(leaseID string, amount int64, paymentDate time.Time) (ddd.Event, error) {
	if l.LeaseStatus != LeaseStatusActive {
		return nil, errors.Wrap(errors.ErrConflict, "lease not active")
	}
	if amount <= 0 {
		return nil, errors.Wrap(errors.ErrBadRequest, "payment must be > 0")
	}
	l.AddEvent(LeasePaymentMadeEvent, &LeasePaymentMade{
		LeaseID:     leaseID,
		Amount:      amount,
		PaymentDate: paymentDate,
	})
	return ddd.NewEvent(LeasePaymentMadeEvent, l), nil
}

func (l *Lease) ExecuteLeaseBuyout(leaseID string, buyoutAmount int64) (ddd.Event, error) {
	if !l.HasBuyout {
		return nil, errors.Wrap(errors.ErrConflict, "lease does not allow buyout")
	}
	if l.LeaseStatus != LeaseStatusActive {
		return nil, errors.Wrap(errors.ErrConflict, "lease not active")
	}
	l.AddEvent(LeaseBuyoutExecutedEvent, &LeaseBuyoutExecuted{
		LeaseID:      leaseID,
		BuyoutAmount: buyoutAmount,
	})
	return ddd.NewEvent(LeaseBuyoutExecutedEvent, l), nil
}

func (l *Lease) EndLease(leaseID string) (ddd.Event, error) {
	if l.LeaseStatus != LeaseStatusActive {
		return nil, errors.Wrap(errors.ErrConflict, "lease not active")
	}
	l.AddEvent(LeaseEndedEvent, &LeaseEnded{
		LeaseID: leaseID,
	})
	return ddd.NewEvent(LeaseEndedEvent, l), nil
}

func (l *Lease) DefaultLease(leaseID string, reason string) (ddd.Event, error) {
	if l.LeaseStatus != LeaseStatusActive {
		return nil, errors.Wrap(errors.ErrConflict, "lease not active")
	}
	l.AddEvent(LeaseDefaultedEvent, &LeaseDefaulted{
		LeaseID: leaseID,
		Reason:  reason,
	})
	return ddd.NewEvent(LeaseDefaultedEvent, l), nil
}

func (l *Lease) RequestNegotiation(leaseID string, userCustomerID string, proposedMonthly int64, proposedTerm int64) (ddd.Event, error) {
	if l.LeaseStatus != LeaseStatusPending {
		return nil, errors.Wrap(errors.ErrConflict, "lease must be pending to request negotiation")
	}
	if proposedMonthly <= 0 {
		return nil, errors.Wrap(errors.ErrBadRequest, "proposed monthly must be > 0")
	}

	l.AddEvent(LeaseNegotiationRequestedEvent, &LeaseNegotiationRequested{
		LeaseID:            leaseID,
		UserCustomerID:     userCustomerID,
		ProposedMonthly:    proposedMonthly,
		ProposedTermMonths: proposedTerm,
	})
	return ddd.NewEvent(LeaseNegotiationRequestedEvent, l), nil
}

func (l *Lease) UpdateNegotiation(leaseID string, userCustomerID string, newMonthlyPrice int64, newTermMonths int64) (ddd.Event, error) {
	if l.LeaseStatus != LeaseStatusNegotiating && l.LeaseStatus != LeaseStatusPending {
		return nil, errors.Wrap(errors.ErrConflict, "lease not in a negotiable state")
	}
	if newMonthlyPrice <= 0 {
		return nil, errors.Wrap(errors.ErrBadRequest, "new monthly price must be > 0")
	}

	l.AddEvent(LeaseNegotiationUpdatedEvent, &LeaseNegotiationUpdated{
		LeaseID:         leaseID,
		UserCustomerID:  userCustomerID,
		NewMonthlyPrice: newMonthlyPrice,
		NewTermMonths:   newTermMonths,
	})
	return ddd.NewEvent(LeaseNegotiationUpdatedEvent, l), nil
}

func (l *Lease) DeclineLease(leaseID, reason string) (ddd.Event, error) {
	if l.LeaseStatus == LeaseStatusActive || l.LeaseStatus == LeaseStatusCompleted {
		return nil, errors.Wrap(errors.ErrConflict, "cannot decline an active or completed lease")
	}
	l.AddEvent(LeaseDeclinedEvent, &LeaseDeclined{
		LeaseID:       leaseID,
		DeclineReason: reason,
	})
	return ddd.NewEvent(LeaseDeclinedEvent, l), nil
}

func (l *Lease) RejectLease(leaseID, reason string) (ddd.Event, error) {
	if l.LeaseStatus == LeaseStatusActive || l.LeaseStatus == LeaseStatusCompleted {
		return nil, errors.Wrap(errors.ErrConflict, "cannot reject an active or completed lease")
	}
	l.AddEvent(LeaseRejectedEvent, &LeaseRejected{
		LeaseID:      leaseID,
		RejectReason: reason,
	})
	return ddd.NewEvent(LeaseRejectedEvent, l), nil
}

// CancelBuyBack
func (l *Lease) CancelLease() (ddd.Event, error) {
	if l.LeaseStatus != LeaseStatusActive {
		return nil, errors.Wrap(errors.ErrConflict, "buyBack not active or cannot cancel")
	}
	l.AddEvent(LeaseCanceledEvent, &LeaseCanceled{})
	return ddd.NewEvent(LeaseCanceledEvent, l), nil
}

// AcceptNegotiation: The seller or system accepts the new negotiated price
func (l *Lease) AcceptNegotiation(newPrice int64, userCustomerID string) (ddd.Event, error) {
	if l.LeaseStatus != LeaseStatusPending {
		return nil, errors.Wrap(errors.ErrConflict, "no negotiation pending or mismatch status")
	}
	if newPrice <= 0 {
		return nil, errors.Wrap(errors.ErrBadRequest, "invalid negotiated price")
	}

	l.AddEvent(LeaseNegotiationAcceptedEvent, &LeaseNegotiationAccepted{
		LeaseID:         l.ID(),
		NegotiatedPrice: newPrice,
		UserCustomerID:  userCustomerID,
	})
	return ddd.NewEvent(BuyNowNegotiationAcceptedEvent, l), nil
}
func (l *Lease) ApplyEvent(event ddd.Event) error {
	switch e := event.Payload().(type) {

	case *LeaseCreated:
		l.OfferID = e.OfferID
		l.MonthlyPrice = e.MonthlyPrice
		l.LeaseTermMonths = e.LeaseTermMonths
		l.HasBuyout = e.HasBuyout
		l.BuyoutPrice = e.BuyoutPrice
		l.LeaseStatus = e.LeaseStatus // typically "pending" or whatever was set

	case *LeaseStarted:
		l.LeaseStartDate = e.LeaseStartDate
		l.LeaseEndDate = e.LeaseEndDate
		l.LeaseStatus = LeaseStatusActive

	case *LeasePaymentMade:
		// Possibly track total paid somewhere

	case *LeaseBuyoutExecuted:
		l.LeaseStatus = LeaseStatusCompleted

	case *LeaseEnded:
		l.LeaseStatus = LeaseStatusCompleted

	case *LeaseDefaulted:
		l.LeaseStatus = LeaseStatusDefaulted

	case *LeaseNegotiationAccepted:
		l.LeaseStatus = LeaseStatusNegotiating
		// Optionally record proposed monthly price/term if aggregator wants to
	case *LeaseNegotiationRequested:
		l.LeaseStatus = LeaseStatusNegotiating
		// Optionally record proposed monthly price/term if aggregator wants to

	case *LeaseNegotiationUpdated:
		// If aggregator tracks the new terms:
		l.MonthlyPrice = e.NewMonthlyPrice
		l.LeaseTermMonths = e.NewTermMonths
		l.LeaseStatus = LeaseStatusNegotiating // remains negotiating

	case *LeaseDeclined:
		l.LeaseStatus = LeaseStatusDeclined

	case *LeaseRejected:
		l.LeaseStatus = LeaseStatusRejected

	default:
		return errors.ErrInternal.Msgf("%T got unexpected payload %T", l, e)
	}
	return nil
}

func (l *Lease) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *LeaseV1:
		l.OfferID = ss.OfferID
		l.MonthlyPrice = ss.MonthlyPrice
		l.LeaseTermMonths = ss.LeaseTermMonths
		l.HasBuyout = ss.HasBuyout
		l.BuyoutPrice = ss.BuyoutPrice
		l.LeaseStartDate = ss.LeaseStartDate
		l.LeaseEndDate = ss.LeaseEndDate
		l.LeaseStatus = ss.LeaseStatus
	default:
		return errors.ErrInternal.Msgf("%T got unexpected snapshot %T", l, snapshot)
	}
	return nil
}

func (l Lease) ToSnapshot() es.Snapshot {
	return &LeaseV1{
		OfferID:         l.OfferID,
		MonthlyPrice:    l.MonthlyPrice,
		LeaseTermMonths: l.LeaseTermMonths,
		HasBuyout:       l.HasBuyout,
		BuyoutPrice:     l.BuyoutPrice,
		LeaseStartDate:  l.LeaseStartDate,
		LeaseEndDate:    l.LeaseEndDate,
		LeaseStatus:     l.LeaseStatus,
	}
}
