package domain

import (
	"context"
	"time"
)

// MatchRepository interface for match persistence
type MatchRepository interface {
	// Find retrieves a match by ID
	Find(ctx context.Context, matchID string) (*Match, error)
	
	// FindByTeams retrieves matches by teams
	FindByTeams(ctx context.Context, teamID string, isHome bool) ([]*Match, error)
	
	// FindByDate retrieves matches on a specific date
	FindByDate(ctx context.Context, date time.Time) ([]*Match, error)
	
	// FindByDateRange retrieves matches within a date range
	FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*Match, error)
	
	// FindByCompetition retrieves matches for a competition
	FindByCompetition(ctx context.Context, competitionID string) ([]*Match, error)
	
	// FindByStadium retrieves matches at a specific stadium
	FindByStadium(ctx context.Context, stadiumID string) ([]*Match, error)
	
	// FindByStatus retrieves matches by status
	FindByStatus(ctx context.Context, status MatchStatus) ([]*Match, error)
	
	// Save persists a match
	Save(ctx context.Context, match *Match) error
	
	// Update updates an existing match
	Update(ctx context.Context, match *Match) error
	
	// Delete removes a match
	Delete(ctx context.Context, matchID string) error
	
	// GetUpcoming retrieves upcoming matches
	GetUpcoming(ctx context.Context, limit int) ([]*Match, error)
	
	// GetOnSale retrieves matches currently on sale
	GetOnSale(ctx context.Context) ([]*Match, error)
	
	// SearchMatches searches matches based on criteria
	SearchMatches(ctx context.Context, criteria MatchSearchCriteria) ([]*Match, error)
}

// MatchSearchCriteria for searching matches
type MatchSearchCriteria struct {
	TeamIDs       []string
	CompetitionID string
	StadiumID     string
	DateFrom      time.Time
	DateTo        time.Time
	Status        []MatchStatus
	MinPrice      int64
	MaxPrice      int64
	HasAvailability bool
	SortBy        string // "date", "price", "availability"
	SortOrder     string // "asc", "desc"
	Page          int
	PageSize      int
}

// MatchCatalogRepository interface for read-optimized match queries
type MatchCatalogRepository interface {
	// GetMatchCatalog retrieves matches for display
	GetMatchCatalog(ctx context.Context, filters MatchCatalogFilters) (*MatchCatalog, error)
	
	// GetMatchDetails retrieves detailed match information
	GetMatchDetails(ctx context.Context, matchID string) (*MatchDetails, error)
	
	// GetSectorAvailability retrieves availability for all sectors
	GetSectorAvailability(ctx context.Context, matchID string) ([]SectorAvailability, error)
	
	// GetSeatMap retrieves the complete seat map for a sector
	GetSeatMap(ctx context.Context, matchID, sectorID string) (*SectorSeatMap, error)
}

// MatchCatalogFilters for catalog queries
type MatchCatalogFilters struct {
	TeamID        string
	CompetitionID string
	Month         int
	Year          int
	ShowSoldOut   bool
	Page          int
	PageSize      int
}

// MatchCatalog represents a catalog response
type MatchCatalog struct {
	Matches    []*MatchSummary
	TotalCount int
	Page       int
	PageSize   int
	HasMore    bool
}

// MatchSummary represents a summary view of a match
type MatchSummary struct {
	MatchID         string
	HomeTeam        TeamInfo
	AwayTeam        TeamInfo
	Competition     CompetitionInfo
	MatchDate       time.Time
	Stadium         StadiumInfo
	Status          MatchStatus
	MinPrice        int64
	MaxPrice        int64
	AvailableTickets int64
	TotalCapacity   int64
	ThumbnailURL    string
}

// MatchDetails represents detailed match information
type MatchDetails struct {
	*Match
	Sectors         []SectorInfo
	PriceRange      PriceRange
	WeatherForecast WeatherInfo
	TransportInfo   []TransportOption
	NearbyEvents    []*MatchSummary
}

// TeamInfo represents team information for display
type TeamInfo struct {
	ID        string
	Name      string
	ShortName string
	LogoURL   string
}

// CompetitionInfo represents competition information for display
type CompetitionInfo struct {
	ID      string
	Name    string
	Type    CompetitionType
	LogoURL string
}

// StadiumInfo represents stadium information for display
type StadiumInfo struct {
	ID       string
	Name     string
	City     string
	Country  string
	Capacity int64
}

// SectorInfo represents sector information for display
type SectorInfo struct {
	ID             string
	Name           string
	Category       string
	Level          int
	TotalSeats     int64
	AvailableSeats int64
	MinPrice       int64
	MaxPrice       int64
	Amenities      []string
}

// SectorAvailability represents seat availability in a sector
type SectorAvailability struct {
	SectorID       string
	SectorName     string
	Category       string
	Available      int64
	Total          int64
	PercentageFull float64
	PriceRange     PriceRange
}

// SectorSeatMap represents the complete seat map for a sector
type SectorSeatMap struct {
	SectorID   string
	SectorName string
	Rows       []RowSeatMap
}

// RowSeatMap represents seats in a row
type RowSeatMap struct {
	RowID     string
	RowNumber string
	Seats     []SeatInfo
}

// SeatInfo represents individual seat information
type SeatInfo struct {
	Number    int
	Status    SeatStatus
	Category  string
	Price     int64
	Features  []string
}

// PriceRange represents a price range
type PriceRange struct {
	Min int64
	Max int64
}

// WeatherInfo represents weather information
type WeatherInfo struct {
	Temperature int
	Condition   string
	WindSpeed   int
	Humidity    int
}

// TransportOption represents transportation options
type TransportOption struct {
	Type        string // "metro", "bus", "parking"
	Name        string
	Distance    int // meters
	WalkingTime int // minutes
}