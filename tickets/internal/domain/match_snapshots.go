package domain

import (
	"time"
)

// MatchV1 snapshot
type MatchV1 struct {
	// Basic Information
	HomeTeam         Team
	AwayTeam         Team
	Competition      Competition
	MatchDate        time.Time
	Status           MatchStatus
	
	// Stadium Information
	Stadium          Stadium
	TotalCapacity    int64
	AvailableTickets int64
	SoldTickets      int64
	
	// Sectors and Seating
	Sectors          map[string]*Sector
	
	// Pricing
	BasePriceByCategory map[string]int64
	DynamicPricing      bool
	
	// Sales Information
	SalesStartDate   time.Time
	SalesEndDate     time.Time
	EarlySalesDate   time.Time
	
	// Match Details
	Referee          string
	Linesmen         []string
	VAR              string
	Weather          string
	Temperature      int
	Attendance       int64
	
	// Media
	ThumbnailURL     string
	BannerURL        string
	
	// Timestamps
	CreatedAt        time.Time
	UpdatedAt        time.Time
	PublishedAt      time.Time
}

// SnapshotName implements es.Snapshot
func (MatchV1) SnapshotName() string { return "tickets.MatchV1" }