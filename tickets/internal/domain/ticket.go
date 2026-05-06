package domain

import (
	"middleman/internal/ddd"
	"middleman/internal/es"
	"time"

	"github.com/stackus/errors"
)

const TicketAggregate = "tickets.Ticket"

var (
	ErrInvalidTicketData      = errors.Wrap(errors.ErrBadRequest, "invalid ticket data")
	ErrTicketAlreadyUsed      = errors.Wrap(errors.ErrBadRequest, "ticket has already been used")
	ErrTicketNotValid         = errors.Wrap(errors.ErrBadRequest, "ticket is not valid")
	ErrTicketExpired          = errors.Wrap(errors.ErrBadRequest, "ticket has expired")
	ErrTicketTransferBlocked  = errors.Wrap(errors.ErrBadRequest, "ticket transfer is blocked")
	ErrInvalidQRCode          = errors.Wrap(errors.ErrBadRequest, "invalid QR code")
	ErrTicketAlreadyCancelled = errors.Wrap(errors.ErrBadRequest, "ticket is already cancelled")
)

// TicketStatus represents the current status of a ticket
type TicketStatus string

const (
	TicketStatusActive      TicketStatus = "active"
	TicketStatusUsed        TicketStatus = "used"
	TicketStatusCancelled   TicketStatus = "cancelled"
	TicketStatusExpired     TicketStatus = "expired"
	TicketStatusTransferred TicketStatus = "transferred"
	TicketStatusRefunded    TicketStatus = "refunded"
)

// TicketType represents the type of ticket
type TicketType string

const (
	TicketTypeRegular    TicketType = "regular"
	TicketTypeSeason     TicketType = "season"
	TicketTypeVIP        TicketType = "vip"
	TicketTypePress      TicketType = "press"
	TicketTypeComplimentary TicketType = "complimentary"
	TicketTypeYouth      TicketType = "youth"
	TicketTypeSenior     TicketType = "senior"
	TicketTypeFamily     TicketType = "family"
)

// Ticket represents a ticket aggregate
type Ticket struct {
	es.Aggregate
	
	// Match Information
	MatchID      string
	MatchDate    time.Time
	HomeTeam     string
	AwayTeam     string
	Competition  string
	
	// Seat Information
	StadiumName  string
	SectorID     string
	SectorName   string
	RowID        string
	RowNumber    string
	SeatNumber   int
	EntranceGate string
	
	// Ticket Details
	Type         TicketType
	Status       TicketStatus
	Category     string
	Price        int64
	
	// Owner Information
	OwnerID      string
	OwnerName    string
	OwnerEmail   string
	OwnerPhone   string
	
	// Purchase Information
	PurchaserID  string
	PurchaseDate time.Time
	PaymentID    string
	OrderID      string
	
	// Security
	QRCode       string
	Barcode      string
	SecurityCode string
	
	// Transfer Information
	Transferable bool
	TransferHistory []TransferRecord
	
	// Usage Information
	UsedAt       time.Time
	UsedGate     string
	ValidationHistory []ValidationRecord
	
	// Timestamps
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ExpiresAt    time.Time
}

// TransferRecord represents a ticket transfer
type TransferRecord struct {
	FromUserID   string
	ToUserID     string
	TransferDate time.Time
	Reason       string
}

// ValidationRecord represents a ticket validation attempt
type ValidationRecord struct {
	ValidatedAt  time.Time
	Gate         string
	Result       string
	Reason       string
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Ticket)(nil)

func NewTicket(id string) *Ticket {
	return &Ticket{
		Aggregate:       es.NewAggregate(id, TicketAggregate),
		TransferHistory: []TransferRecord{},
		ValidationHistory: []ValidationRecord{},
	}
}

// Key implements registry.Registerable
func (Ticket) Key() string { return TicketAggregate }

// InitTicket initializes a new ticket
func (t *Ticket) InitTicket(
	matchID string,
	matchDate time.Time,
	homeTeam, awayTeam, competition string,
	stadiumName, sectorID, sectorName, rowID, rowNumber string,
	seatNumber int,
	entranceGate string,
	ticketType TicketType,
	category string,
	price int64,
	ownerID, ownerName, ownerEmail, ownerPhone string,
	purchaserID, paymentID, orderID string,
	transferable bool,
) (ddd.Event, error) {
	if matchID == "" || ownerID == "" {
		return nil, ErrInvalidTicketData
	}
	
	if price < 0 {
		return nil, errors.Wrap(errors.ErrBadRequest, "ticket price cannot be negative")
	}

	// Generate security codes
	qrCode := generateQRCode(t.ID(), matchID, sectorID, rowID, seatNumber)
	barcode := generateBarcode(t.ID())
	securityCode := generateSecurityCode()

	t.AddEvent(TicketCreatedEvent, &TicketCreated{
		MatchID:      matchID,
		MatchDate:    matchDate,
		HomeTeam:     homeTeam,
		AwayTeam:     awayTeam,
		Competition:  competition,
		StadiumName:  stadiumName,
		SectorID:     sectorID,
		SectorName:   sectorName,
		RowID:        rowID,
		RowNumber:    rowNumber,
		SeatNumber:   seatNumber,
		EntranceGate: entranceGate,
		Type:         ticketType,
		Category:     category,
		Price:        price,
		OwnerID:      ownerID,
		OwnerName:    ownerName,
		OwnerEmail:   ownerEmail,
		OwnerPhone:   ownerPhone,
		PurchaserID:  purchaserID,
		PurchaseDate: time.Now(),
		PaymentID:    paymentID,
		OrderID:      orderID,
		QRCode:       qrCode,
		Barcode:      barcode,
		SecurityCode: securityCode,
		Transferable: transferable,
		Status:       TicketStatusActive,
		CreatedAt:    time.Now(),
		ExpiresAt:    matchDate.Add(24 * time.Hour), // Expires 24 hours after match
	})
	return ddd.NewEvent(TicketCreatedEvent, t), nil
}

// TransferTicket transfers the ticket to another user
func (t *Ticket) TransferTicket(toUserID, toUserName, toUserEmail, toUserPhone, reason string) (ddd.Event, error) {
	if !t.Transferable {
		return nil, ErrTicketTransferBlocked
	}
	
	if t.Status != TicketStatusActive {
		return nil, ErrTicketNotValid
	}
	
	if time.Now().After(t.ExpiresAt) {
		return nil, ErrTicketExpired
	}

	t.AddEvent(TicketTransferredEvent, &TicketTransferred{
		TicketID:    t.ID(),
		FromUserID:  t.OwnerID,
		ToUserID:    toUserID,
		ToUserName:  toUserName,
		ToUserEmail: toUserEmail,
		ToUserPhone: toUserPhone,
		Reason:      reason,
		TransferredAt: time.Now(),
	})
	return ddd.NewEvent(TicketTransferredEvent, t), nil
}

// ValidateTicket validates the ticket at entry
func (t *Ticket) ValidateTicket(gate, qrCode string) (ddd.Event, error) {
	if t.Status == TicketStatusUsed {
		return nil, ErrTicketAlreadyUsed
	}
	
	if t.Status != TicketStatusActive {
		return nil, ErrTicketNotValid
	}
	
	if time.Now().After(t.ExpiresAt) {
		return nil, ErrTicketExpired
	}
	
	if qrCode != t.QRCode {
		return nil, ErrInvalidQRCode
	}

	t.AddEvent(TicketValidatedEvent, &TicketValidated{
		TicketID:    t.ID(),
		Gate:        gate,
		ValidatedAt: time.Now(),
		Success:     true,
	})
	return ddd.NewEvent(TicketValidatedEvent, t), nil
}

// UseTicket marks the ticket as used
func (t *Ticket) UseTicket(gate string) (ddd.Event, error) {
	if t.Status == TicketStatusUsed {
		return nil, ErrTicketAlreadyUsed
	}
	
	if t.Status != TicketStatusActive {
		return nil, ErrTicketNotValid
	}

	t.AddEvent(TicketUsedEvent, &TicketUsed{
		TicketID: t.ID(),
		UsedAt:   time.Now(),
		Gate:     gate,
	})
	return ddd.NewEvent(TicketUsedEvent, t), nil
}

// CancelTicket cancels the ticket
func (t *Ticket) CancelTicket(reason string, refundAmount int64) (ddd.Event, error) {
	if t.Status == TicketStatusCancelled || t.Status == TicketStatusRefunded {
		return nil, ErrTicketAlreadyCancelled
	}
	
	if t.Status == TicketStatusUsed {
		return nil, errors.Wrap(errors.ErrBadRequest, "cannot cancel used ticket")
	}

	t.AddEvent(TicketCancelledDomainEvent, &TicketCancelledDomain{
		TicketID:     t.ID(),
		Reason:       reason,
		RefundAmount: refundAmount,
		CancelledAt:  time.Now(),
	})
	return ddd.NewEvent(TicketCancelledDomainEvent, t), nil
}

// ApplyEvent implements es.EventApplier
func (t *Ticket) ApplyEvent(event ddd.Event) error {
	switch e := event.Payload().(type) {
	case *TicketCreated:
		t.MatchID = e.MatchID
		t.MatchDate = e.MatchDate
		t.HomeTeam = e.HomeTeam
		t.AwayTeam = e.AwayTeam
		t.Competition = e.Competition
		t.StadiumName = e.StadiumName
		t.SectorID = e.SectorID
		t.SectorName = e.SectorName
		t.RowID = e.RowID
		t.RowNumber = e.RowNumber
		t.SeatNumber = e.SeatNumber
		t.EntranceGate = e.EntranceGate
		t.Type = e.Type
		t.Category = e.Category
		t.Price = e.Price
		t.OwnerID = e.OwnerID
		t.OwnerName = e.OwnerName
		t.OwnerEmail = e.OwnerEmail
		t.OwnerPhone = e.OwnerPhone
		t.PurchaserID = e.PurchaserID
		t.PurchaseDate = e.PurchaseDate
		t.PaymentID = e.PaymentID
		t.OrderID = e.OrderID
		t.QRCode = e.QRCode
		t.Barcode = e.Barcode
		t.SecurityCode = e.SecurityCode
		t.Transferable = e.Transferable
		t.Status = e.Status
		t.CreatedAt = e.CreatedAt
		t.UpdatedAt = e.CreatedAt
		t.ExpiresAt = e.ExpiresAt

	case *TicketTransferred:
		t.TransferHistory = append(t.TransferHistory, TransferRecord{
			FromUserID:   e.FromUserID,
			ToUserID:     e.ToUserID,
			TransferDate: e.TransferredAt,
			Reason:       e.Reason,
		})
		t.OwnerID = e.ToUserID
		t.OwnerName = e.ToUserName
		t.OwnerEmail = e.ToUserEmail
		t.OwnerPhone = e.ToUserPhone
		t.UpdatedAt = e.TransferredAt

	case *TicketValidated:
		t.ValidationHistory = append(t.ValidationHistory, ValidationRecord{
			ValidatedAt: e.ValidatedAt,
			Gate:        e.Gate,
			Result:      "success",
			Reason:      "",
		})

	case *TicketUsed:
		t.Status = TicketStatusUsed
		t.UsedAt = e.UsedAt
		t.UsedGate = e.Gate
		t.UpdatedAt = e.UsedAt

	case *TicketCancelledDomain:
		if e.RefundAmount > 0 {
			t.Status = TicketStatusRefunded
		} else {
			t.Status = TicketStatusCancelled
		}
		t.UpdatedAt = e.CancelledAt

	default:
		return errors.ErrInternal.Msgf(
			"%T received the event %s with unexpected payload %T",
			t, event.EventName(), e)
	}
	return nil
}

// ToSnapshot implements es.Snapshotter
func (t Ticket) ToSnapshot() es.Snapshot {
	return TicketV1{
		MatchID:           t.MatchID,
		MatchDate:         t.MatchDate,
		HomeTeam:          t.HomeTeam,
		AwayTeam:          t.AwayTeam,
		Competition:       t.Competition,
		StadiumName:       t.StadiumName,
		SectorID:          t.SectorID,
		SectorName:        t.SectorName,
		RowID:             t.RowID,
		RowNumber:         t.RowNumber,
		SeatNumber:        t.SeatNumber,
		EntranceGate:      t.EntranceGate,
		Type:              t.Type,
		Status:            t.Status,
		Category:          t.Category,
		Price:             t.Price,
		OwnerID:           t.OwnerID,
		OwnerName:         t.OwnerName,
		OwnerEmail:        t.OwnerEmail,
		OwnerPhone:        t.OwnerPhone,
		PurchaserID:       t.PurchaserID,
		PurchaseDate:      t.PurchaseDate,
		PaymentID:         t.PaymentID,
		OrderID:           t.OrderID,
		QRCode:            t.QRCode,
		Barcode:           t.Barcode,
		SecurityCode:      t.SecurityCode,
		Transferable:      t.Transferable,
		TransferHistory:   t.TransferHistory,
		UsedAt:            t.UsedAt,
		UsedGate:          t.UsedGate,
		ValidationHistory: t.ValidationHistory,
		CreatedAt:         t.CreatedAt,
		UpdatedAt:         t.UpdatedAt,
		ExpiresAt:         t.ExpiresAt,
	}
}

// ApplySnapshot implements es.Snapshotter
func (t *Ticket) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *TicketV1:
		t.MatchID = ss.MatchID
		t.MatchDate = ss.MatchDate
		t.HomeTeam = ss.HomeTeam
		t.AwayTeam = ss.AwayTeam
		t.Competition = ss.Competition
		t.StadiumName = ss.StadiumName
		t.SectorID = ss.SectorID
		t.SectorName = ss.SectorName
		t.RowID = ss.RowID
		t.RowNumber = ss.RowNumber
		t.SeatNumber = ss.SeatNumber
		t.EntranceGate = ss.EntranceGate
		t.Type = ss.Type
		t.Status = ss.Status
		t.Category = ss.Category
		t.Price = ss.Price
		t.OwnerID = ss.OwnerID
		t.OwnerName = ss.OwnerName
		t.OwnerEmail = ss.OwnerEmail
		t.OwnerPhone = ss.OwnerPhone
		t.PurchaserID = ss.PurchaserID
		t.PurchaseDate = ss.PurchaseDate
		t.PaymentID = ss.PaymentID
		t.OrderID = ss.OrderID
		t.QRCode = ss.QRCode
		t.Barcode = ss.Barcode
		t.SecurityCode = ss.SecurityCode
		t.Transferable = ss.Transferable
		t.TransferHistory = ss.TransferHistory
		t.UsedAt = ss.UsedAt
		t.UsedGate = ss.UsedGate
		t.ValidationHistory = ss.ValidationHistory
		t.CreatedAt = ss.CreatedAt
		t.UpdatedAt = ss.UpdatedAt
		t.ExpiresAt = ss.ExpiresAt

	default:
		return errors.ErrInternal.Msgf(
			"%T received the unexpected snapshot %T", t, snapshot)
	}
	return nil
}

// Helper functions for generating codes
func generateQRCode(ticketID, matchID, sectorID, rowID string, seatNumber int) string {
	// In production, use proper QR code generation
	return "QR-" + ticketID + "-" + matchID
}

func generateBarcode(ticketID string) string {
	// In production, use proper barcode generation
	return "BAR-" + ticketID
}

func generateSecurityCode() string {
	// In production, use cryptographically secure random generation
	return "SEC-" + time.Now().Format("20060102150405")
}