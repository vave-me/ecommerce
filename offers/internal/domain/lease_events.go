package domain

import "time"

const (
	LeaseCreatedEvent              = "offers.LeaseCreated"
	LeaseStartedEvent              = "offers.LeaseStarted"
	LeasePaymentMadeEvent          = "offers.LeasePaymentMade"
	LeaseBuyoutExecutedEvent       = "offers.LeaseBuyoutExecuted"
	LeaseEndedEvent                = "offers.LeaseEnded"
	LeaseDefaultedEvent            = "offers.LeaseDefaulted"
	LeaseNegotiationRequestedEvent = "offers.LeaseNegotiationRequested"
	LeaseNegotiationAcceptedEvent  = "offers.LeaseNegotiationAccepted"
	LeaseNegotiationUpdatedEvent   = "offers.LeaseNegotiationUpdated"
	LeaseDeclinedEvent             = "offers.LeaseDeclined"
	LeaseRejectedEvent             = "offers.LeaseRejected"
	LeaseCanceledEvent             = "offers.LeaseCanceled"
)

// LeaseCreated is raised when a new lease is created (pending).
type LeaseCreated struct {
	OfferID         string
	MonthlyPrice    int64
	LeaseTermMonths int64
	HasBuyout       bool
	BuyoutPrice     int64
	LeaseStatus     LeaseStatus
}

func (LeaseCreated) Key() string { return LeaseCreatedEvent }

// LeaseStarted is raised when the lease actually begins (active).
type LeaseStarted struct {
	LeaseID        string
	LeaseStartDate time.Time
	LeaseEndDate   time.Time
}

func (LeaseStarted) Key() string { return LeaseStartedEvent }

type LeaseCanceled struct {
	LeaseID string
}

func (LeaseCanceled) Key() string { return LeaseCanceledEvent }

// LeasePaymentMade is raised whenever a payment is made for the lease.
type LeasePaymentMade struct {
	LeaseID     string
	Amount      int64
	PaymentDate time.Time
}

func (LeasePaymentMade) Key() string { return LeasePaymentMadeEvent }

// LeaseBuyoutExecuted is raised when the buyer exercises the buyout option.
type LeaseBuyoutExecuted struct {
	LeaseID      string
	BuyoutAmount int64
}

func (LeaseBuyoutExecuted) Key() string { return LeaseBuyoutExecutedEvent }

// LeaseEnded is raised if the lease ended normally.
type LeaseEnded struct {
	LeaseID string
}

func (LeaseEnded) Key() string { return LeaseEndedEvent }

// LeaseDefaulted is raised if the lease is defaulted.
type LeaseDefaulted struct {
	LeaseID string
	Reason  string
}

func (LeaseDefaulted) Key() string { return LeaseDefaultedEvent }
type LeaseNegotiationRequested struct {
	LeaseID            string
	UserCustomerID     string
	ProposedMonthly    int64
	ProposedTermMonths int64
}

func (LeaseNegotiationRequested) Key() string { return LeaseNegotiationRequestedEvent }

type LeaseNegotiationUpdated struct {
	LeaseID         string
	UserCustomerID  string
	NewMonthlyPrice int64
	NewTermMonths   int64
}

func (LeaseNegotiationUpdated) Key() string { return LeaseNegotiationUpdatedEvent }

type LeaseDeclined struct {
	LeaseID       string
	DeclineReason string
}

func (LeaseDeclined) Key() string { return LeaseDeclinedEvent }

type LeaseRejected struct {
	LeaseID      string
	RejectReason string
}

func (LeaseRejected) Key() string { return LeaseRejectedEvent }

type LeaseNegotiationAccepted struct {
	LeaseID         string
	NegotiatedPrice int64
	UserCustomerID  string
}

func (LeaseNegotiationAccepted) Key() string { return LeaseNegotiationAcceptedEvent }
