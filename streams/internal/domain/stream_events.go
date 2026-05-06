package domain

import (
	"time"
)

// Event names
const (
	StreamCreatedEvent         = "streams.StreamCreated"
	StreamUpdatedEvent         = "streams.StreamUpdated"
	StreamQualitiesSetEvent    = "streams.StreamQualitiesSet"
	SubtitleAddedEvent         = "streams.SubtitleAdded"
	AudioTrackAddedEvent       = "streams.AudioTrackAdded"
	StreamPricingSetEvent      = "streams.StreamPricingSet"
	StreamPublishedEvent       = "streams.StreamPublished"
	UserAccessGrantedEvent     = "streams.UserAccessGranted"
	UserAccessRevokedEvent     = "streams.UserAccessRevoked"
	WatchProgressUpdatedEvent  = "streams.WatchProgressUpdated"
	StreamViewedEvent          = "streams.StreamViewed"
	StreamRatedEvent           = "streams.StreamRated"
	StreamArchivedEvent        = "streams.StreamArchived"
	SeriesCreatedEvent         = "streams.SeriesCreated"
	SeasonAddedEvent           = "streams.SeasonAdded"
	EpisodeAddedToSeasonEvent  = "streams.EpisodeAddedToSeason"
)

// StreamCreated event payload
type StreamCreated struct {
	Title         string
	Description   string
	Synopsis      string
	StreamType    StreamType
	StreamURL     string
	ThumbnailURL  string
	TrailerURL    string
	Duration      int64
	ContentRating ContentRating
	AccessType    AccessType
	Genre         []string
	Language      string
	Country       string
	Studio        string
	Status        StreamStatus
	CreatedAt     time.Time
}

// Key implements registry.Registerable
func (StreamCreated) Key() string { return StreamCreatedEvent }

// StreamUpdated event payload
type StreamUpdated struct {
	StreamID      string
	Title         string
	Description   string
	Synopsis      string
	ThumbnailURL  string
	TrailerURL    string
	ContentRating ContentRating
	Genre         []string
	Tags          []string
	UpdatedAt     time.Time
}

// Key implements registry.Registerable
func (StreamUpdated) Key() string { return StreamUpdatedEvent }

// StreamQualitiesSet event payload
type StreamQualitiesSet struct {
	StreamID       string
	Qualities      []StreamQuality
	DefaultQuality StreamQuality
}

// Key implements registry.Registerable
func (StreamQualitiesSet) Key() string { return StreamQualitiesSetEvent }

// SubtitleAdded event payload
type SubtitleAdded struct {
	StreamID string
	Language string
	URL      string
	Default  bool
}

// Key implements registry.Registerable
func (SubtitleAdded) Key() string { return SubtitleAddedEvent }

// AudioTrackAdded event payload
type AudioTrackAdded struct {
	StreamID string
	Language string
	Type     string
	Default  bool
}

// Key implements registry.Registerable
func (AudioTrackAdded) Key() string { return AudioTrackAddedEvent }

// StreamPricingSet event payload
type StreamPricingSet struct {
	StreamID       string
	RentalPrice    int64
	RentalDuration int64
	PurchasePrice  int64
	PPVPrice       int64
}

// Key implements registry.Registerable
func (StreamPricingSet) Key() string { return StreamPricingSetEvent }

// StreamPublished event payload
type StreamPublished struct {
	StreamID    string
	PublishedAt time.Time
}

// Key implements registry.Registerable
func (StreamPublished) Key() string { return StreamPublishedEvent }

// UserAccessGranted event payload
type UserAccessGranted struct {
	StreamID   string
	UserID     string
	AccessType AccessType
	GrantedAt  time.Time
	ExpiresAt  time.Time
}

// Key implements registry.Registerable
func (UserAccessGranted) Key() string { return UserAccessGrantedEvent }

// UserAccessRevoked event payload
type UserAccessRevoked struct {
	StreamID  string
	UserID    string
	RevokedAt time.Time
}

// Key implements registry.Registerable
func (UserAccessRevoked) Key() string { return UserAccessRevokedEvent }

// WatchProgressUpdated event payload
type WatchProgressUpdated struct {
	StreamID      string
	UserID        string
	Progress      int64
	Completed     bool
	LastWatchedAt time.Time
}

// Key implements registry.Registerable
func (WatchProgressUpdated) Key() string { return WatchProgressUpdatedEvent }

// StreamViewed event payload
type StreamViewed struct {
	StreamID  string
	ViewedAt  time.Time
	ViewCount int64
}

// Key implements registry.Registerable
func (StreamViewed) Key() string { return StreamViewedEvent }

// StreamRated event payload
type StreamRated struct {
	StreamID string
	UserID   string
	Rating   int
	IsLike   bool
	RatedAt  time.Time
}

// Key implements registry.Registerable
func (StreamRated) Key() string { return StreamRatedEvent }

// StreamArchived event payload
type StreamArchived struct {
	StreamID   string
	ArchivedAt time.Time
}

// Key implements registry.Registerable
func (StreamArchived) Key() string { return StreamArchivedEvent }

// Series-related events

// SeriesCreated event payload
type SeriesCreated struct {
	SeriesID     string
	Title        string
	Description  string
	ThumbnailURL string
	Genre        []string
	Studio       string
	CreatedAt    time.Time
}

// Key implements registry.Registerable
func (SeriesCreated) Key() string { return SeriesCreatedEvent }

// SeasonAdded event payload
type SeasonAdded struct {
	SeriesID     string
	SeasonID     string
	SeasonNumber int
	Title        string
	Description  string
	ThumbnailURL string
	CreatedAt    time.Time
}

// Key implements registry.Registerable
func (SeasonAdded) Key() string { return SeasonAddedEvent }

// EpisodeAddedToSeason event payload
type EpisodeAddedToSeason struct {
	SeriesID      string
	SeasonID      string
	EpisodeID     string
	EpisodeNumber int
	StreamID      string // Reference to the actual stream
	AddedAt       time.Time
}

// Key implements registry.Registerable
func (EpisodeAddedToSeason) Key() string { return EpisodeAddedToSeasonEvent }