package domain

import (
	"context"
	"time"
)

// TicketRepository interface for ticket persistence
type TicketRepository interface {
	// Find retrieves a ticket by ID
	Find(ctx context.Context, ticketID string) (*Ticket, error)
	
	// FindByMatch retrieves all tickets for a match
	FindByMatch(ctx context.Context, matchID string) ([]*Ticket, error)
	
	// FindByOwner retrieves tickets owned by a user
	FindByOwner(ctx context.Context, ownerID string) ([]*Ticket, error)
	
	// FindByPurchaser retrieves tickets purchased by a user
	FindByPurchaser(ctx context.Context, purchaserID string) ([]*Ticket, error)
	
	// FindByOrder retrieves tickets in an order
	FindByOrder(ctx context.Context, orderID string) ([]*Ticket, error)
	
	// FindBySeat retrieves a ticket for a specific seat
	FindBySeat(ctx context.Context, matchID, sectorID, rowID string, seatNumber int) (*Ticket, error)
	
	// FindByQRCode retrieves a ticket by QR code
	FindByQRCode(ctx context.Context, qrCode string) (*Ticket, error)
	
	// FindByBarcode retrieves a ticket by barcode
	FindByBarcode(ctx context.Context, barcode string) (*Ticket, error)
	
	// Save persists a ticket
	Save(ctx context.Context, ticket *Ticket) error
	
	// Update updates an existing ticket
	Update(ctx context.Context, ticket *Ticket) error
	
	// Delete removes a ticket
	Delete(ctx context.Context, ticketID string) error
	
	// GetUserUpcomingTickets retrieves upcoming match tickets for a user
	GetUserUpcomingTickets(ctx context.Context, userID string) ([]*Ticket, error)
	
	// GetUserPastTickets retrieves past match tickets for a user
	GetUserPastTickets(ctx context.Context, userID string, limit int) ([]*Ticket, error)
	
	// GetTransferableTickets retrieves tickets that can be transferred
	GetTransferableTickets(ctx context.Context, userID string) ([]*Ticket, error)
	
	// CountTicketsByStatus counts tickets by status for a match
	CountTicketsByStatus(ctx context.Context, matchID string) (map[TicketStatus]int64, error)
	
	// GetExpiringTickets retrieves tickets expiring soon
	GetExpiringTickets(ctx context.Context, before time.Time) ([]*Ticket, error)
}

// TicketSearchRepository interface for ticket search operations
type TicketSearchRepository interface {
	// SearchTickets searches tickets based on criteria
	SearchTickets(ctx context.Context, criteria TicketSearchCriteria) ([]*Ticket, error)
	
	// GetTicketHistory retrieves ticket history for reporting
	GetTicketHistory(ctx context.Context, filters TicketHistoryFilters) (*TicketHistory, error)
	
	// GetSalesStatistics retrieves sales statistics
	GetSalesStatistics(ctx context.Context, matchID string) (*SalesStatistics, error)
	
	// GetRevenueReport retrieves revenue information
	GetRevenueReport(ctx context.Context, filters RevenueFilters) (*RevenueReport, error)
}

// TicketSearchCriteria for searching tickets
type TicketSearchCriteria struct {
	MatchID      string
	OwnerID      string
	PurchaserID  string
	Status       []TicketStatus
	Type         []TicketType
	DateFrom     time.Time
	DateTo       time.Time
	MinPrice     int64
	MaxPrice     int64
	SectorID     string
	Transferable *bool
	Page         int
	PageSize     int
}

// TicketHistoryFilters for ticket history queries
type TicketHistoryFilters struct {
	UserID   string
	DateFrom time.Time
	DateTo   time.Time
	Status   []TicketStatus
	SortBy   string
}

// TicketHistory represents ticket history
type TicketHistory struct {
	Tickets      []*TicketHistoryItem
	TotalCount   int
	TotalSpent   int64
	MatchesCount int
}

// TicketHistoryItem represents a ticket in history
type TicketHistoryItem struct {
	TicketID     string
	MatchInfo    MatchInfo
	SeatInfo     SeatDetails
	PurchaseDate time.Time
	Price        int64
	Status       TicketStatus
	Used         bool
}

// MatchInfo represents match information for tickets
type MatchInfo struct {
	MatchID     string
	HomeTeam    string
	AwayTeam    string
	MatchDate   time.Time
	Stadium     string
	Competition string
}

// SeatDetails represents seat details for tickets
type SeatDetails struct {
	SectorName string
	RowNumber  string
	SeatNumber int
	Category   string
	Gate       string
}

// SalesStatistics represents ticket sales statistics
type SalesStatistics struct {
	MatchID          string
	TotalTickets     int64
	SoldTickets      int64
	AvailableTickets int64
	Revenue          int64
	AveragePrice     int64
	SalesByCategory  map[string]CategorySales
	SalesByType      map[TicketType]int64
	SalesTimeline    []SalesTimepoint
}

// CategorySales represents sales by category
type CategorySales struct {
	Sold      int64
	Available int64
	Revenue   int64
}

// SalesTimepoint represents sales at a point in time
type SalesTimepoint struct {
	Time         time.Time
	TotalSold    int64
	Revenue      int64
	SalesRate    float64 // tickets per hour
}

// RevenueFilters for revenue queries
type RevenueFilters struct {
	DateFrom      time.Time
	DateTo        time.Time
	TeamID        string
	CompetitionID string
	GroupBy       string // "day", "week", "month", "match"
}

// RevenueReport represents revenue information
type RevenueReport struct {
	TotalRevenue    int64
	TicketsSold     int64
	AveragePrice    int64
	RefundAmount    int64
	NetRevenue      int64
	RevenueByPeriod []RevenuePeriod
	TopMatches      []MatchRevenue
}

// RevenuePeriod represents revenue for a period
type RevenuePeriod struct {
	Period       string
	Revenue      int64
	TicketsSold  int64
	MatchesCount int
}

// MatchRevenue represents revenue for a match
type MatchRevenue struct {
	MatchID      string
	MatchInfo    MatchInfo
	Revenue      int64
	TicketsSold  int64
	Capacity     int64
	Attendance   float64 // percentage
}