package domain

import (
	"time"
)

// Event names
const (
	MatchCreatedEvent        = "tickets.MatchCreated"
	MatchUpdatedEvent        = "tickets.MatchUpdated"
	SectorAddedEvent         = "tickets.SectorAdded"
	RowAddedToSectorEvent    = "tickets.RowAddedToSector"
	SeatsInitializedEvent    = "tickets.SeatsInitialized"
	DynamicPricingSetEvent   = "tickets.DynamicPricingSet"
	MatchPublishedEvent      = "tickets.MatchPublished"
	TicketSalesStartedEvent  = "tickets.TicketSalesStarted"
	MatchStatusUpdatedEvent  = "tickets.MatchStatusUpdated"
	SeatReservedEvent        = "tickets.SeatReserved"
	SeatReleasedEvent        = "tickets.SeatReleased"
	TicketPurchasedEvent     = "tickets.TicketPurchased"
	TicketCancelledEvent     = "tickets.TicketCancelled"
	MatchCancelledEvent      = "tickets.MatchCancelled"
	MatchPostponedEvent      = "tickets.MatchPostponed"
	AttendanceUpdatedEvent   = "tickets.AttendanceUpdated"
)

// MatchCreated event payload
type MatchCreated struct {
	HomeTeam       Team
	AwayTeam       Team
	Competition    Competition
	MatchDate      time.Time
	Stadium        Stadium
	TotalCapacity  int64
	SalesStartDate time.Time
	SalesEndDate   time.Time
	Status         MatchStatus
	CreatedAt      time.Time
}

// Key implements registry.Registerable
func (MatchCreated) Key() string { return MatchCreatedEvent }

// MatchUpdated event payload
type MatchUpdated struct {
	MatchID       string
	MatchDate     time.Time
	Referee       string
	Linesmen      []string
	VAR           string
	ThumbnailURL  string
	BannerURL     string
	UpdatedAt     time.Time
}

// Key implements registry.Registerable
func (MatchUpdated) Key() string { return MatchUpdatedEvent }

// SectorAdded event payload
type SectorAdded struct {
	MatchID       string
	SectorID      string
	Name          string
	Level         int
	Category      string
	TotalSeats    int64
	BasePrice     int64
	Amenities     []string
	EntranceGates []string
}

// Key implements registry.Registerable
func (SectorAdded) Key() string { return SectorAddedEvent }

// RowAddedToSector event payload
type RowAddedToSector struct {
	MatchID    string
	SectorID   string
	RowID      string
	RowNumber  string
	TotalSeats int
}

// Key implements registry.Registerable
func (RowAddedToSector) Key() string { return RowAddedToSectorEvent }

// SeatsInitialized event payload
type SeatsInitialized struct {
	MatchID  string
	SectorID string
	RowID    string
	Seats    []int // Seat numbers
	Price    int64
}

// Key implements registry.Registerable
func (SeatsInitialized) Key() string { return SeatsInitializedEvent }

// DynamicPricingSet event payload
type DynamicPricingSet struct {
	MatchID    string
	Enabled    bool
	BasePrices map[string]int64 // Category -> Base price
}

// Key implements registry.Registerable
func (DynamicPricingSet) Key() string { return DynamicPricingSetEvent }

// MatchPublished event payload
type MatchPublished struct {
	MatchID     string
	PublishedAt time.Time
}

// Key implements registry.Registerable
func (MatchPublished) Key() string { return MatchPublishedEvent }

// TicketSalesStarted event payload
type TicketSalesStarted struct {
	MatchID   string
	StartedAt time.Time
}

// Key implements registry.Registerable
func (TicketSalesStarted) Key() string { return TicketSalesStartedEvent }

// MatchStatusUpdated event payload
type MatchStatusUpdated struct {
	MatchID   string
	OldStatus MatchStatus
	NewStatus MatchStatus
	UpdatedAt time.Time
}

// Key implements registry.Registerable
func (MatchStatusUpdated) Key() string { return MatchStatusUpdatedEvent }

// SeatReserved event payload
type SeatReserved struct {
	MatchID       string
	SectorID      string
	RowID         string
	SeatNumber    int
	ReservedUntil time.Time
}

// Key implements registry.Registerable
func (SeatReserved) Key() string { return SeatReservedEvent }

// SeatReleased event payload
type SeatReleased struct {
	MatchID    string
	SectorID   string
	RowID      string
	SeatNumber int
	ReleasedAt time.Time
}

// Key implements registry.Registerable
func (SeatReleased) Key() string { return SeatReleasedEvent }

// TicketPurchased event payload
type TicketPurchased struct {
	MatchID     string
	TicketID    string
	SectorID    string
	RowID       string
	SeatNumber  int
	UserID      string
	Price       int64
	Category    string
	PurchasedAt time.Time
}

// Key implements registry.Registerable
func (TicketPurchased) Key() string { return TicketPurchasedEvent }

// TicketCancelled event payload
type TicketCancelled struct {
	MatchID     string
	TicketID    string
	SectorID    string
	RowID       string
	SeatNumber  int
	Reason      string
	RefundAmount int64
	CancelledAt time.Time
}

// Key implements registry.Registerable
func (TicketCancelled) Key() string { return TicketCancelledEvent }

// MatchCancelled event payload
type MatchCancelled struct {
	MatchID     string
	Reason      string
	CancelledAt time.Time
}

// Key implements registry.Registerable
func (MatchCancelled) Key() string { return MatchCancelledEvent }

// MatchPostponed event payload
type MatchPostponed struct {
	MatchID      string
	OldDate      time.Time
	NewDate      time.Time
	Reason       string
	PostponedAt  time.Time
}

// Key implements registry.Registerable
func (MatchPostponed) Key() string { return MatchPostponedEvent }

// AttendanceUpdated event payload
type AttendanceUpdated struct {
	MatchID    string
	Attendance int64
	UpdatedAt  time.Time
}

// Key implements registry.Registerable
func (AttendanceUpdated) Key() string { return AttendanceUpdatedEvent }