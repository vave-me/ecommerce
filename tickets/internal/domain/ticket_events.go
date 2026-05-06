package domain

import (
	"time"
)

// Event names for ticket aggregate
const (
	TicketCreatedEvent          = "tickets.TicketCreated"
	TicketTransferredEvent      = "tickets.TicketTransferred"
	TicketValidatedEvent        = "tickets.TicketValidated"
	TicketUsedEvent             = "tickets.TicketUsed"
	TicketCancelledDomainEvent  = "tickets.TicketCancelled"
	TicketRefundedEvent         = "tickets.TicketRefunded"
	TicketExpiredEvent          = "tickets.TicketExpired"
)

// TicketCreated event payload
type TicketCreated struct {
	MatchID      string
	MatchDate    time.Time
	HomeTeam     string
	AwayTeam     string
	Competition  string
	StadiumName  string
	SectorID     string
	SectorName   string
	RowID        string
	RowNumber    string
	SeatNumber   int
	EntranceGate string
	Type         TicketType
	Category     string
	Price        int64
	OwnerID      string
	OwnerName    string
	OwnerEmail   string
	OwnerPhone   string
	PurchaserID  string
	PurchaseDate time.Time
	PaymentID    string
	OrderID      string
	QRCode       string
	Barcode      string
	SecurityCode string
	Transferable bool
	Status       TicketStatus
	CreatedAt    time.Time
	ExpiresAt    time.Time
}

// Key implements registry.Registerable
func (TicketCreated) Key() string { return TicketCreatedEvent }

// TicketTransferred event payload
type TicketTransferred struct {
	TicketID      string
	FromUserID    string
	ToUserID      string
	ToUserName    string
	ToUserEmail   string
	ToUserPhone   string
	Reason        string
	TransferredAt time.Time
}

// Key implements registry.Registerable
func (TicketTransferred) Key() string { return TicketTransferredEvent }

// TicketValidated event payload
type TicketValidated struct {
	TicketID    string
	Gate        string
	ValidatedAt time.Time
	Success     bool
	FailReason  string
}

// Key implements registry.Registerable
func (TicketValidated) Key() string { return TicketValidatedEvent }

// TicketUsed event payload
type TicketUsed struct {
	TicketID string
	UsedAt   time.Time
	Gate     string
}

// Key implements registry.Registerable
func (TicketUsed) Key() string { return TicketUsedEvent }

// TicketCancelledDomain event payload
type TicketCancelledDomain struct {
	TicketID     string
	Reason       string
	RefundAmount int64
	CancelledAt  time.Time
}

// Key implements registry.Registerable
func (TicketCancelledDomain) Key() string { return TicketCancelledDomainEvent }

// TicketRefunded event payload
type TicketRefunded struct {
	TicketID     string
	RefundAmount int64
	RefundReason string
	RefundedAt   time.Time
}

// Key implements registry.Registerable
func (TicketRefunded) Key() string { return TicketRefundedEvent }

// TicketExpired event payload
type TicketExpired struct {
	TicketID  string
	ExpiredAt time.Time
}

// Key implements registry.Registerable
func (TicketExpired) Key() string { return TicketExpiredEvent }