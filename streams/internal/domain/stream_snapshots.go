package domain

import (
	"time"
)

// StreamV1 snapshot
type StreamV1 struct {
	// Basic Information
	Title           string
	Description     string
	Synopsis        string
	StreamType      StreamType
	Status          StreamStatus
	
	// Content Details
	StreamURL       string
	ThumbnailURL    string
	TrailerURL      string
	Duration        int64
	ReleaseDate     time.Time
	ContentRating   ContentRating
	
	// Technical Details
	AvailableQualities []StreamQuality
	DefaultQuality     StreamQuality
	Subtitles          []Subtitle
	AudioTracks        []AudioTrack
	
	// Access Control
	AccessType        AccessType
	SubscriptionTiers []string
	RentalPrice       int64
	RentalDuration    int64
	PurchasePrice     int64
	PPVPrice          int64
	
	// Metadata
	Genre         []string
	Tags          []string
	Cast          []CastMember
	Directors     []string
	Producers     []string
	Studio        string
	Language      string
	Country       string
	
	// Analytics
	ViewCount     int64
	LikeCount     int64
	DislikeCount  int64
	AverageRating float64
	TotalRevenue  int64
	
	// User Access Management
	UserAccess map[string]UserAccessInfo
	
	// Series Information
	SeriesID      string
	SeasonNumber  int
	EpisodeNumber int
	
	// Timestamps
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PublishedAt time.Time
}

// SnapshotName implements es.Snapshot
func (StreamV1) SnapshotName() string { return "streams.StreamV1" }

// SeriesV1 snapshot for TV series
type SeriesV1 struct {
	SeriesID     string
	Title        string
	Description  string
	ThumbnailURL string
	Genre        []string
	Studio       string
	Seasons      []SeasonInfo
	TotalSeasons int
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// SeasonInfo contains information about a season
type SeasonInfo struct {
	SeasonID      string
	SeasonNumber  int
	Title         string
	Description   string
	ThumbnailURL  string
	Episodes      []EpisodeInfo
	TotalEpisodes int
	CreatedAt     time.Time
}

// EpisodeInfo contains basic episode information
type EpisodeInfo struct {
	EpisodeID     string
	EpisodeNumber int
	StreamID      string // Reference to the actual stream
	Title         string
	Duration      int64
	AirDate       time.Time
}

// SnapshotName implements es.Snapshot
func (SeriesV1) SnapshotName() string { return "streams.SeriesV1" }