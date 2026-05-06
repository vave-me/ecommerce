package domain

import (
	"time"
)

// TicketV1 snapshot
type TicketV1 struct {
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

// SnapshotName implements es.Snapshot
func (TicketV1) SnapshotName() string { return "tickets.TicketV1" }