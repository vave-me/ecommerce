package domain

import (
	"middleman/internal/ddd"
	"middleman/internal/es"
	"time"

	"github.com/stackus/errors"
)

const MatchAggregate = "tickets.Match"

var (
	ErrMatchNameIsBlank        = errors.Wrap(errors.ErrBadRequest, "the match name cannot be blank")
	ErrInvalidMatchDate        = errors.Wrap(errors.ErrBadRequest, "the match date is invalid")
	ErrInvalidStadium          = errors.Wrap(errors.ErrBadRequest, "the stadium information is invalid")
	ErrInvalidCapacity         = errors.Wrap(errors.ErrBadRequest, "the stadium capacity is invalid")
	ErrMatchAlreadyPublished   = errors.Wrap(errors.ErrBadRequest, "the match is already published")
	ErrMatchNotPublished       = errors.Wrap(errors.ErrBadRequest, "the match is not published")
	ErrMatchAlreadyStarted     = errors.Wrap(errors.ErrBadRequest, "the match has already started")
	ErrMatchCancelled          = errors.Wrap(errors.ErrBadRequest, "the match has been cancelled")
	ErrInvalidSectorCapacity   = errors.Wrap(errors.ErrBadRequest, "the sector capacity is invalid")
	ErrSectorNotFound          = errors.Wrap(errors.ErrBadRequest, "the sector not found")
	ErrSectorAlreadyExists     = errors.Wrap(errors.ErrBadRequest, "the sector already exists")
)

// MatchStatus represents the current status of a match
type MatchStatus string

const (
	MatchStatusDraft      MatchStatus = "draft"
	MatchStatusPublished  MatchStatus = "published"
	MatchStatusOnSale     MatchStatus = "on_sale"
	MatchStatusSoldOut    MatchStatus = "sold_out"
	MatchStatusInProgress MatchStatus = "in_progress"
	MatchStatusCompleted  MatchStatus = "completed"
	MatchStatusCancelled  MatchStatus = "cancelled"
	MatchStatusPostponed  MatchStatus = "postponed"
)

// CompetitionType represents the type of competition
type CompetitionType string

const (
	CompetitionTypeLeague        CompetitionType = "league"
	CompetitionTypeCup           CompetitionType = "cup"
	CompetitionTypeChampionsLeague CompetitionType = "champions_league"
	CompetitionTypeEuropaLeague  CompetitionType = "europa_league"
	CompetitionTypeFriendly      CompetitionType = "friendly"
	CompetitionTypeInternational CompetitionType = "international"
)

// Match represents a football match aggregate
type Match struct {
	es.Aggregate
	
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
	Sectors          map[string]*Sector // SectorID -> Sector
	
	// Pricing
	BasePriceByCategory map[string]int64 // Category -> Base price in cents
	DynamicPricing      bool
	
	// Sales Information
	SalesStartDate   time.Time
	SalesEndDate     time.Time
	EarlySalesDate   time.Time // For season ticket holders, members, etc.
	
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

// Team represents a football team
type Team struct {
	ID           string
	Name         string
	ShortName    string
	LogoURL      string
	HomeStadium  string
	City         string
	Country      string
	Founded      int
}

// Competition represents a football competition
type Competition struct {
	ID           string
	Name         string
	Type         CompetitionType
	Season       string
	Round        int
	Stage        string // e.g., "Group Stage", "Round of 16", "Final"
	LogoURL      string
}

// Stadium represents a football stadium
type Stadium struct {
	ID           string
	Name         string
	City         string
	Country      string
	Address      string
	Latitude     float64
	Longitude    float64
	Capacity     int64
	YearBuilt    int
	ImageURLs    []string
	Facilities   []string
}

// Sector represents a section of the stadium
type Sector struct {
	ID           string
	Name         string
	Level        int // 0 = ground, 1 = first tier, 2 = second tier, etc.
	Category     string // e.g., "VIP", "Premium", "Standard", "Away"
	TotalSeats   int64
	AvailableSeats int64
	SoldSeats    int64
	BasePrice    int64
	Rows         map[string]*Row // RowID -> Row
	Amenities    []string
	EntranceGates []string
}

// Row represents a row within a sector
type Row struct {
	ID           string
	Number       string
	TotalSeats   int
	Seats        map[int]*Seat // Seat number -> Seat
}

// Seat represents an individual seat
type Seat struct {
	Number       int
	Status       SeatStatus
	Category     string // Can override sector category (e.g., "Wheelchair", "Restricted View")
	TicketID     string // Reference to ticket if sold
	ReservedUntil time.Time // For temporary reservations
	Price        int64 // Can override sector price
	Features     []string // e.g., "Cup Holder", "Extra Legroom"
}

// SeatStatus represents the status of a seat
type SeatStatus string

const (
	SeatStatusAvailable SeatStatus = "available"
	SeatStatusReserved  SeatStatus = "reserved"
	SeatStatusSold      SeatStatus = "sold"
	SeatStatusBlocked   SeatStatus = "blocked"
	SeatStatusMaintenance SeatStatus = "maintenance"
)

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Match)(nil)

func NewMatch(id string) *Match {
	return &Match{
		Aggregate:           es.NewAggregate(id, MatchAggregate),
		Sectors:             make(map[string]*Sector),
		BasePriceByCategory: make(map[string]int64),
	}
}

// Key implements registry.Registerable
func (Match) Key() string { return MatchAggregate }

// InitMatch initializes a new match
func (m *Match) InitMatch(
	homeTeam, awayTeam Team,
	competition Competition,
	matchDate time.Time,
	stadium Stadium,
	salesStartDate, salesEndDate time.Time,
) (ddd.Event, error) {
	if homeTeam.Name == "" || awayTeam.Name == "" {
		return nil, ErrMatchNameIsBlank
	}
	if matchDate.Before(time.Now()) {
		return nil, ErrInvalidMatchDate
	}
	if stadium.Capacity <= 0 {
		return nil, ErrInvalidCapacity
	}

	m.AddEvent(MatchCreatedEvent, &MatchCreated{
		HomeTeam:       homeTeam,
		AwayTeam:       awayTeam,
		Competition:    competition,
		MatchDate:      matchDate,
		Stadium:        stadium,
		TotalCapacity:  stadium.Capacity,
		SalesStartDate: salesStartDate,
		SalesEndDate:   salesEndDate,
		Status:         MatchStatusDraft,
		CreatedAt:      time.Now(),
	})
	return ddd.NewEvent(MatchCreatedEvent, m), nil
}

// AddSector adds a sector to the stadium
func (m *Match) AddSector(
	sectorID, name string,
	level int,
	category string,
	totalSeats int64,
	basePrice int64,
	amenities []string,
	entranceGates []string,
) (ddd.Event, error) {
	if _, exists := m.Sectors[sectorID]; exists {
		return nil, ErrSectorAlreadyExists
	}
	
	if totalSeats <= 0 {
		return nil, ErrInvalidSectorCapacity
	}
	
	if basePrice < 0 {
		return nil, errors.Wrap(errors.ErrBadRequest, "sector price cannot be negative")
	}

	m.AddEvent(SectorAddedEvent, &SectorAdded{
		MatchID:       m.ID(),
		SectorID:      sectorID,
		Name:          name,
		Level:         level,
		Category:      category,
		TotalSeats:    totalSeats,
		BasePrice:     basePrice,
		Amenities:     amenities,
		EntranceGates: entranceGates,
	})
	return ddd.NewEvent(SectorAddedEvent, m), nil
}

// AddRowToSector adds a row to a specific sector
func (m *Match) AddRowToSector(
	sectorID, rowID, rowNumber string,
	totalSeats int,
) (ddd.Event, error) {
	sector, exists := m.Sectors[sectorID]
	if !exists {
		return nil, ErrSectorNotFound
	}
	
	if _, exists := sector.Rows[rowID]; exists {
		return nil, errors.Wrap(errors.ErrBadRequest, "row already exists")
	}

	m.AddEvent(RowAddedToSectorEvent, &RowAddedToSector{
		MatchID:    m.ID(),
		SectorID:   sectorID,
		RowID:      rowID,
		RowNumber:  rowNumber,
		TotalSeats: totalSeats,
	})
	return ddd.NewEvent(RowAddedToSectorEvent, m), nil
}

// InitializeSeatsInRow initializes all seats in a row
func (m *Match) InitializeSeatsInRow(
	sectorID, rowID string,
	startNumber, endNumber int,
) (ddd.Event, error) {
	sector, exists := m.Sectors[sectorID]
	if !exists {
		return nil, ErrSectorNotFound
	}
	
	row, exists := sector.Rows[rowID]
	if !exists {
		return nil, errors.Wrap(errors.ErrBadRequest, "row not found")
	}

	seats := make([]int, 0, endNumber-startNumber+1)
	for i := startNumber; i <= endNumber; i++ {
		seats = append(seats, i)
	}

	m.AddEvent(SeatsInitializedEvent, &SeatsInitialized{
		MatchID:  m.ID(),
		SectorID: sectorID,
		RowID:    rowID,
		Seats:    seats,
		Price:    sector.BasePrice,
	})
	return ddd.NewEvent(SeatsInitializedEvent, m), nil
}

// SetDynamicPricing enables or disables dynamic pricing
func (m *Match) SetDynamicPricing(enabled bool, basePrices map[string]int64) (ddd.Event, error) {
	m.AddEvent(DynamicPricingSetEvent, &DynamicPricingSet{
		MatchID:        m.ID(),
		Enabled:        enabled,
		BasePrices:     basePrices,
	})
	return ddd.NewEvent(DynamicPricingSetEvent, m), nil
}

// PublishMatch publishes the match making it visible to fans
func (m *Match) PublishMatch() (ddd.Event, error) {
	if m.Status == MatchStatusPublished || m.Status == MatchStatusOnSale {
		return nil, ErrMatchAlreadyPublished
	}
	
	if len(m.Sectors) == 0 {
		return nil, errors.Wrap(errors.ErrBadRequest, "cannot publish match without sectors")
	}

	m.AddEvent(MatchPublishedEvent, &MatchPublished{
		MatchID:     m.ID(),
		PublishedAt: time.Now(),
	})
	return ddd.NewEvent(MatchPublishedEvent, m), nil
}

// StartTicketSales starts the ticket sales
func (m *Match) StartTicketSales() (ddd.Event, error) {
	if m.Status != MatchStatusPublished {
		return nil, ErrMatchNotPublished
	}
	
	if time.Now().After(m.MatchDate) {
		return nil, ErrMatchAlreadyStarted
	}

	m.AddEvent(TicketSalesStartedEvent, &TicketSalesStarted{
		MatchID:   m.ID(),
		StartedAt: time.Now(),
	})
	return ddd.NewEvent(TicketSalesStartedEvent, m), nil
}

// UpdateMatchStatus updates the match status
func (m *Match) UpdateMatchStatus(status MatchStatus) (ddd.Event, error) {
	if m.Status == MatchStatusCancelled {
		return nil, ErrMatchCancelled
	}

	m.AddEvent(MatchStatusUpdatedEvent, &MatchStatusUpdated{
		MatchID:   m.ID(),
		OldStatus: m.Status,
		NewStatus: status,
		UpdatedAt: time.Now(),
	})
	return ddd.NewEvent(MatchStatusUpdatedEvent, m), nil
}

// ReserveSeat temporarily reserves a seat
func (m *Match) ReserveSeat(sectorID, rowID string, seatNumber int, duration time.Duration) (ddd.Event, error) {
	sector, exists := m.Sectors[sectorID]
	if !exists {
		return nil, ErrSectorNotFound
	}
	
	row, exists := sector.Rows[rowID]
	if !exists {
		return nil, errors.Wrap(errors.ErrBadRequest, "row not found")
	}
	
	seat, exists := row.Seats[seatNumber]
	if !exists {
		return nil, errors.Wrap(errors.ErrBadRequest, "seat not found")
	}
	
	if seat.Status != SeatStatusAvailable {
		return nil, errors.Wrap(errors.ErrBadRequest, "seat is not available")
	}

	m.AddEvent(SeatReservedEvent, &SeatReserved{
		MatchID:       m.ID(),
		SectorID:      sectorID,
		RowID:         rowID,
		SeatNumber:    seatNumber,
		ReservedUntil: time.Now().Add(duration),
	})
	return ddd.NewEvent(SeatReservedEvent, m), nil
}

// ApplyEvent implements es.EventApplier
func (m *Match) ApplyEvent(event ddd.Event) error {
	switch e := event.Payload().(type) {
	case *MatchCreated:
		m.HomeTeam = e.HomeTeam
		m.AwayTeam = e.AwayTeam
		m.Competition = e.Competition
		m.MatchDate = e.MatchDate
		m.Stadium = e.Stadium
		m.TotalCapacity = e.TotalCapacity
		m.AvailableTickets = e.TotalCapacity
		m.SalesStartDate = e.SalesStartDate
		m.SalesEndDate = e.SalesEndDate
		m.Status = e.Status
		m.CreatedAt = e.CreatedAt
		m.UpdatedAt = e.CreatedAt

	case *SectorAdded:
		m.Sectors[e.SectorID] = &Sector{
			ID:             e.SectorID,
			Name:           e.Name,
			Level:          e.Level,
			Category:       e.Category,
			TotalSeats:     e.TotalSeats,
			AvailableSeats: e.TotalSeats,
			BasePrice:      e.BasePrice,
			Rows:           make(map[string]*Row),
			Amenities:      e.Amenities,
			EntranceGates:  e.EntranceGates,
		}

	case *RowAddedToSector:
		if sector, exists := m.Sectors[e.SectorID]; exists {
			sector.Rows[e.RowID] = &Row{
				ID:         e.RowID,
				Number:     e.RowNumber,
				TotalSeats: e.TotalSeats,
				Seats:      make(map[int]*Seat),
			}
		}

	case *SeatsInitialized:
		if sector, exists := m.Sectors[e.SectorID]; exists {
			if row, exists := sector.Rows[e.RowID]; exists {
				for _, seatNum := range e.Seats {
					row.Seats[seatNum] = &Seat{
						Number:   seatNum,
						Status:   SeatStatusAvailable,
						Category: sector.Category,
						Price:    e.Price,
					}
				}
			}
		}

	case *DynamicPricingSet:
		m.DynamicPricing = e.Enabled
		m.BasePriceByCategory = e.BasePrices

	case *MatchPublished:
		m.Status = MatchStatusPublished
		m.PublishedAt = e.PublishedAt

	case *TicketSalesStarted:
		m.Status = MatchStatusOnSale

	case *MatchStatusUpdated:
		m.Status = e.NewStatus
		m.UpdatedAt = e.UpdatedAt

	case *SeatReserved:
		if sector, exists := m.Sectors[e.SectorID]; exists {
			if row, exists := sector.Rows[e.RowID]; exists {
				if seat, exists := row.Seats[e.SeatNumber]; exists {
					seat.Status = SeatStatusReserved
					seat.ReservedUntil = e.ReservedUntil
				}
			}
		}

	default:
		return errors.ErrInternal.Msgf(
			"%T received the event %s with unexpected payload %T",
			m, event.EventName(), e)
	}
	return nil
}

// ToSnapshot implements es.Snapshotter
func (m Match) ToSnapshot() es.Snapshot {
	return MatchV1{
		HomeTeam:            m.HomeTeam,
		AwayTeam:            m.AwayTeam,
		Competition:         m.Competition,
		MatchDate:           m.MatchDate,
		Status:              m.Status,
		Stadium:             m.Stadium,
		TotalCapacity:       m.TotalCapacity,
		AvailableTickets:    m.AvailableTickets,
		SoldTickets:         m.SoldTickets,
		Sectors:             m.Sectors,
		BasePriceByCategory: m.BasePriceByCategory,
		DynamicPricing:      m.DynamicPricing,
		SalesStartDate:      m.SalesStartDate,
		SalesEndDate:        m.SalesEndDate,
		EarlySalesDate:      m.EarlySalesDate,
		Referee:             m.Referee,
		Linesmen:            m.Linesmen,
		VAR:                 m.VAR,
		Weather:             m.Weather,
		Temperature:         m.Temperature,
		Attendance:          m.Attendance,
		ThumbnailURL:        m.ThumbnailURL,
		BannerURL:           m.BannerURL,
		CreatedAt:           m.CreatedAt,
		UpdatedAt:           m.UpdatedAt,
		PublishedAt:         m.PublishedAt,
	}
}

// ApplySnapshot implements es.Snapshotter
func (m *Match) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *MatchV1:
		m.HomeTeam = ss.HomeTeam
		m.AwayTeam = ss.AwayTeam
		m.Competition = ss.Competition
		m.MatchDate = ss.MatchDate
		m.Status = ss.Status
		m.Stadium = ss.Stadium
		m.TotalCapacity = ss.TotalCapacity
		m.AvailableTickets = ss.AvailableTickets
		m.SoldTickets = ss.SoldTickets
		m.Sectors = ss.Sectors
		m.BasePriceByCategory = ss.BasePriceByCategory
		m.DynamicPricing = ss.DynamicPricing
		m.SalesStartDate = ss.SalesStartDate
		m.SalesEndDate = ss.SalesEndDate
		m.EarlySalesDate = ss.EarlySalesDate
		m.Referee = ss.Referee
		m.Linesmen = ss.Linesmen
		m.VAR = ss.VAR
		m.Weather = ss.Weather
		m.Temperature = ss.Temperature
		m.Attendance = ss.Attendance
		m.ThumbnailURL = ss.ThumbnailURL
		m.BannerURL = ss.BannerURL
		m.CreatedAt = ss.CreatedAt
		m.UpdatedAt = ss.UpdatedAt
		m.PublishedAt = ss.PublishedAt

	default:
		return errors.ErrInternal.Msgf(
			"%T received the unexpected snapshot %T", m, snapshot)
	}
	return nil
}