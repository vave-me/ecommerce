package domain

import (
	"middleman/internal/ddd"
	"middleman/internal/es"
	"time"

	"github.com/stackus/errors"
)

const StreamAggregate = "streams.Stream"

var (
	ErrStreamTitleIsBlank         = errors.Wrap(errors.ErrBadRequest, "the stream title cannot be blank")
	ErrInvalidStreamURL           = errors.Wrap(errors.ErrBadRequest, "the stream URL is invalid")
	ErrInvalidDuration            = errors.Wrap(errors.ErrBadRequest, "the stream duration is invalid")
	ErrInvalidQuality             = errors.Wrap(errors.ErrBadRequest, "the stream quality is invalid")
	ErrStreamAlreadyPublished     = errors.Wrap(errors.ErrBadRequest, "the stream is already published")
	ErrStreamNotPublished         = errors.Wrap(errors.ErrBadRequest, "the stream is not published")
	ErrInvalidSubscriptionPrice   = errors.Wrap(errors.ErrBadRequest, "the subscription price is invalid")
	ErrUserAlreadyHasAccess       = errors.Wrap(errors.ErrBadRequest, "user already has access to this stream")
	ErrUserDoesNotHaveAccess      = errors.Wrap(errors.ErrBadRequest, "user does not have access to this stream")
	ErrInvalidRentalPeriod        = errors.Wrap(errors.ErrBadRequest, "the rental period is invalid")
	ErrStreamNotAvailableForRental = errors.Wrap(errors.ErrBadRequest, "the stream is not available for rental")
)

// StreamType represents the type of content
type StreamType string

const (
	StreamTypeMovie      StreamType = "movie"
	StreamTypeSeries     StreamType = "series"
	StreamTypeDocumentary StreamType = "documentary"
	StreamTypeLive       StreamType = "live"
	StreamTypeEducational StreamType = "educational"
	StreamTypeMusic      StreamType = "music"
	StreamTypeSports     StreamType = "sports"
)

// StreamStatus represents the current status of the stream
type StreamStatus string

const (
	StreamStatusDraft      StreamStatus = "draft"
	StreamStatusProcessing StreamStatus = "processing"
	StreamStatusPublished  StreamStatus = "published"
	StreamStatusArchived   StreamStatus = "archived"
	StreamStatusDeleted    StreamStatus = "deleted"
)

// StreamQuality represents the available quality options
type StreamQuality string

const (
	StreamQualitySD    StreamQuality = "480p"
	StreamQualityHD    StreamQuality = "720p"
	StreamQualityFullHD StreamQuality = "1080p"
	StreamQuality4K    StreamQuality = "2160p"
)

// AccessType represents how users can access the stream
type AccessType string

const (
	AccessTypeFree         AccessType = "free"
	AccessTypeSubscription AccessType = "subscription"
	AccessTypeRental       AccessType = "rental"
	AccessTypePurchase     AccessType = "purchase"
	AccessTypePPV          AccessType = "pay_per_view"
)

// ContentRating represents the age rating for content
type ContentRating string

const (
	ContentRatingG     ContentRating = "G"     // General audiences
	ContentRatingPG    ContentRating = "PG"    // Parental guidance
	ContentRatingPG13  ContentRating = "PG-13" // Parents strongly cautioned
	ContentRatingR     ContentRating = "R"     // Restricted
	ContentRatingNC17  ContentRating = "NC-17" // Adults only
	ContentRatingUnrated ContentRating = "NR"   // Not rated
)

// Stream represents a video stream aggregate
type Stream struct {
	es.Aggregate
	
	// Basic Information
	Title           string
	Description     string
	Synopsis        string
	StreamType      StreamType
	Status          StreamStatus
	
	// Content Details
	StreamURL       string         // URL to the actual video stream
	ThumbnailURL    string         // URL to thumbnail image
	TrailerURL      string         // URL to trailer video
	Duration        int64          // Duration in seconds
	ReleaseDate     time.Time
	ContentRating   ContentRating
	
	// Technical Details
	AvailableQualities []StreamQuality
	DefaultQuality     StreamQuality
	Subtitles          []Subtitle
	AudioTracks        []AudioTrack
	
	// Access Control
	AccessType      AccessType
	SubscriptionTiers []string      // Which subscription tiers can access
	RentalPrice     int64          // Price in cents for rental
	RentalDuration  int64          // Rental duration in hours
	PurchasePrice   int64          // Price in cents for purchase
	PPVPrice        int64          // Price for pay-per-view events
	
	// Metadata
	Genre           []string
	Tags            []string
	Cast            []CastMember
	Directors       []string
	Producers       []string
	Studio          string
	Language        string        // Primary language
	Country         string        // Country of origin
	
	// Analytics
	ViewCount       int64
	LikeCount       int64
	DislikeCount    int64
	AverageRating   float64
	TotalRevenue    int64
	
	// User Access Management
	UserAccess      map[string]UserAccessInfo  // UserID -> Access info
	
	// Series Information (for TV series)
	SeriesID        string
	SeasonNumber    int
	EpisodeNumber   int
	
	// Timestamps
	CreatedAt       time.Time
	UpdatedAt       time.Time
	PublishedAt     time.Time
}

// Subtitle represents available subtitle tracks
type Subtitle struct {
	Language string
	URL      string
	Default  bool
}

// AudioTrack represents available audio tracks
type AudioTrack struct {
	Language string
	Type     string // e.g., "original", "dubbed"
	Default  bool
}

// CastMember represents an actor in the content
type CastMember struct {
	Name      string
	Role      string
	Character string
	ImageURL  string
}

// UserAccessInfo tracks individual user access to streams
type UserAccessInfo struct {
	UserID       string
	AccessType   AccessType
	GrantedAt    time.Time
	ExpiresAt    time.Time
	LastWatchedAt time.Time
	WatchProgress int64 // Progress in seconds
	Completed    bool
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Stream)(nil)

func NewStream(id string) *Stream {
	return &Stream{
		Aggregate:   es.NewAggregate(id, StreamAggregate),
		UserAccess:  make(map[string]UserAccessInfo),
		Subtitles:   []Subtitle{},
		AudioTracks: []AudioTrack{},
		Cast:        []CastMember{},
		AvailableQualities: []StreamQuality{},
	}
}

// Key implements registry.Registerable
func (Stream) Key() string { return StreamAggregate }

// InitStream initializes a new stream
func (s *Stream) InitStream(
	title, description, synopsis string,
	streamType StreamType,
	streamURL, thumbnailURL, trailerURL string,
	duration int64,
	contentRating ContentRating,
	accessType AccessType,
	genre []string,
	language, country, studio string,
) (ddd.Event, error) {
	if title == "" {
		return nil, ErrStreamTitleIsBlank
	}
	if streamURL == "" {
		return nil, ErrInvalidStreamURL
	}
	if duration <= 0 {
		return nil, ErrInvalidDuration
	}

	s.AddEvent(StreamCreatedEvent, &StreamCreated{
		Title:         title,
		Description:   description,
		Synopsis:      synopsis,
		StreamType:    streamType,
		StreamURL:     streamURL,
		ThumbnailURL:  thumbnailURL,
		TrailerURL:    trailerURL,
		Duration:      duration,
		ContentRating: contentRating,
		AccessType:    accessType,
		Genre:         genre,
		Language:      language,
		Country:       country,
		Studio:        studio,
		Status:        StreamStatusDraft,
		CreatedAt:     time.Now(),
	})
	return ddd.NewEvent(StreamCreatedEvent, s), nil
}

// SetStreamQualities sets the available quality options
func (s *Stream) SetStreamQualities(qualities []StreamQuality, defaultQuality StreamQuality) (ddd.Event, error) {
	if len(qualities) == 0 {
		return nil, ErrInvalidQuality
	}
	
	// Verify default quality is in the list
	found := false
	for _, q := range qualities {
		if q == defaultQuality {
			found = true
			break
		}
	}
	if !found {
		return nil, ErrInvalidQuality
	}

	s.AddEvent(StreamQualitiesSetEvent, &StreamQualitiesSet{
		StreamID:       s.ID(),
		Qualities:      qualities,
		DefaultQuality: defaultQuality,
	})
	return ddd.NewEvent(StreamQualitiesSetEvent, s), nil
}

// AddSubtitle adds a subtitle track to the stream
func (s *Stream) AddSubtitle(language, url string, isDefault bool) (ddd.Event, error) {
	s.AddEvent(SubtitleAddedEvent, &SubtitleAdded{
		StreamID: s.ID(),
		Language: language,
		URL:      url,
		Default:  isDefault,
	})
	return ddd.NewEvent(SubtitleAddedEvent, s), nil
}

// AddAudioTrack adds an audio track to the stream
func (s *Stream) AddAudioTrack(language, trackType string, isDefault bool) (ddd.Event, error) {
	s.AddEvent(AudioTrackAddedEvent, &AudioTrackAdded{
		StreamID: s.ID(),
		Language: language,
		Type:     trackType,
		Default:  isDefault,
	})
	return ddd.NewEvent(AudioTrackAddedEvent, s), nil
}

// SetPricing sets the pricing for different access types
func (s *Stream) SetPricing(rentalPrice, rentalDuration, purchasePrice, ppvPrice int64) (ddd.Event, error) {
	if s.AccessType == AccessTypeFree {
		return nil, errors.Wrap(errors.ErrBadRequest, "cannot set pricing for free content")
	}
	
	if rentalPrice < 0 || purchasePrice < 0 || ppvPrice < 0 {
		return nil, ErrInvalidSubscriptionPrice
	}
	
	if rentalDuration <= 0 && rentalPrice > 0 {
		return nil, ErrInvalidRentalPeriod
	}

	s.AddEvent(StreamPricingSetEvent, &StreamPricingSet{
		StreamID:       s.ID(),
		RentalPrice:    rentalPrice,
		RentalDuration: rentalDuration,
		PurchasePrice:  purchasePrice,
		PPVPrice:       ppvPrice,
	})
	return ddd.NewEvent(StreamPricingSetEvent, s), nil
}

// PublishStream publishes the stream making it available to users
func (s *Stream) PublishStream() (ddd.Event, error) {
	if s.Status == StreamStatusPublished {
		return nil, ErrStreamAlreadyPublished
	}

	s.AddEvent(StreamPublishedEvent, &StreamPublished{
		StreamID:    s.ID(),
		PublishedAt: time.Now(),
	})
	return ddd.NewEvent(StreamPublishedEvent, s), nil
}

// GrantUserAccess grants a user access to the stream
func (s *Stream) GrantUserAccess(userID string, accessType AccessType, duration int64) (ddd.Event, error) {
	if _, exists := s.UserAccess[userID]; exists {
		if s.UserAccess[userID].ExpiresAt.After(time.Now()) {
			return nil, ErrUserAlreadyHasAccess
		}
	}

	var expiresAt time.Time
	switch accessType {
	case AccessTypeRental:
		if s.AccessType != AccessTypeRental && s.AccessType != AccessTypePurchase {
			return nil, ErrStreamNotAvailableForRental
		}
		expiresAt = time.Now().Add(time.Duration(duration) * time.Hour)
	case AccessTypePurchase:
		expiresAt = time.Now().AddDate(100, 0, 0) // Effectively permanent
	case AccessTypeSubscription:
		expiresAt = time.Now().Add(time.Duration(duration) * time.Hour)
	case AccessTypePPV:
		expiresAt = time.Now().Add(48 * time.Hour) // 48 hours for PPV
	}

	s.AddEvent(UserAccessGrantedEvent, &UserAccessGranted{
		StreamID:   s.ID(),
		UserID:     userID,
		AccessType: accessType,
		GrantedAt:  time.Now(),
		ExpiresAt:  expiresAt,
	})
	return ddd.NewEvent(UserAccessGrantedEvent, s), nil
}

// RevokeUserAccess revokes a user's access to the stream
func (s *Stream) RevokeUserAccess(userID string) (ddd.Event, error) {
	if _, exists := s.UserAccess[userID]; !exists {
		return nil, ErrUserDoesNotHaveAccess
	}

	s.AddEvent(UserAccessRevokedEvent, &UserAccessRevoked{
		StreamID:  s.ID(),
		UserID:    userID,
		RevokedAt: time.Now(),
	})
	return ddd.NewEvent(UserAccessRevokedEvent, s), nil
}

// UpdateWatchProgress updates a user's watch progress
func (s *Stream) UpdateWatchProgress(userID string, progress int64, completed bool) (ddd.Event, error) {
	if _, exists := s.UserAccess[userID]; !exists {
		return nil, ErrUserDoesNotHaveAccess
	}

	s.AddEvent(WatchProgressUpdatedEvent, &WatchProgressUpdated{
		StreamID:      s.ID(),
		UserID:        userID,
		Progress:      progress,
		Completed:     completed,
		LastWatchedAt: time.Now(),
	})
	return ddd.NewEvent(WatchProgressUpdatedEvent, s), nil
}

// IncrementViewCount increments the view count
func (s *Stream) IncrementViewCount() (ddd.Event, error) {
	s.AddEvent(StreamViewedEvent, &StreamViewed{
		StreamID:  s.ID(),
		ViewedAt:  time.Now(),
		ViewCount: s.ViewCount + 1,
	})
	return ddd.NewEvent(StreamViewedEvent, s), nil
}

// UpdateRating updates the stream's rating
func (s *Stream) UpdateRating(userID string, rating int, isLike bool) (ddd.Event, error) {
	s.AddEvent(StreamRatedEvent, &StreamRated{
		StreamID: s.ID(),
		UserID:   userID,
		Rating:   rating,
		IsLike:   isLike,
		RatedAt:  time.Now(),
	})
	return ddd.NewEvent(StreamRatedEvent, s), nil
}

// ArchiveStream archives the stream
func (s *Stream) ArchiveStream() (ddd.Event, error) {
	if s.Status == StreamStatusArchived {
		return nil, errors.Wrap(errors.ErrBadRequest, "stream is already archived")
	}

	s.AddEvent(StreamArchivedEvent, &StreamArchived{
		StreamID:   s.ID(),
		ArchivedAt: time.Now(),
	})
	return ddd.NewEvent(StreamArchivedEvent, s), nil
}

// ApplyEvent implements es.EventApplier
func (s *Stream) ApplyEvent(event ddd.Event) error {
	switch e := event.Payload().(type) {
	case *StreamCreated:
		s.Title = e.Title
		s.Description = e.Description
		s.Synopsis = e.Synopsis
		s.StreamType = e.StreamType
		s.StreamURL = e.StreamURL
		s.ThumbnailURL = e.ThumbnailURL
		s.TrailerURL = e.TrailerURL
		s.Duration = e.Duration
		s.ContentRating = e.ContentRating
		s.AccessType = e.AccessType
		s.Genre = e.Genre
		s.Language = e.Language
		s.Country = e.Country
		s.Studio = e.Studio
		s.Status = e.Status
		s.CreatedAt = e.CreatedAt
		s.UpdatedAt = e.CreatedAt

	case *StreamQualitiesSet:
		s.AvailableQualities = e.Qualities
		s.DefaultQuality = e.DefaultQuality

	case *SubtitleAdded:
		s.Subtitles = append(s.Subtitles, Subtitle{
			Language: e.Language,
			URL:      e.URL,
			Default:  e.Default,
		})

	case *AudioTrackAdded:
		s.AudioTracks = append(s.AudioTracks, AudioTrack{
			Language: e.Language,
			Type:     e.Type,
			Default:  e.Default,
		})

	case *StreamPricingSet:
		s.RentalPrice = e.RentalPrice
		s.RentalDuration = e.RentalDuration
		s.PurchasePrice = e.PurchasePrice
		s.PPVPrice = e.PPVPrice

	case *StreamPublished:
		s.Status = StreamStatusPublished
		s.PublishedAt = e.PublishedAt

	case *UserAccessGranted:
		s.UserAccess[e.UserID] = UserAccessInfo{
			UserID:     e.UserID,
			AccessType: e.AccessType,
			GrantedAt:  e.GrantedAt,
			ExpiresAt:  e.ExpiresAt,
		}

	case *UserAccessRevoked:
		delete(s.UserAccess, e.UserID)

	case *WatchProgressUpdated:
		if access, exists := s.UserAccess[e.UserID]; exists {
			access.WatchProgress = e.Progress
			access.Completed = e.Completed
			access.LastWatchedAt = e.LastWatchedAt
			s.UserAccess[e.UserID] = access
		}

	case *StreamViewed:
		s.ViewCount = e.ViewCount

	case *StreamRated:
		if e.IsLike {
			s.LikeCount++
		} else {
			s.DislikeCount++
		}
		// Note: Average rating would need to be calculated from all ratings

	case *StreamArchived:
		s.Status = StreamStatusArchived

	default:
		return errors.ErrInternal.Msgf(
			"%T received the event %s with unexpected payload %T",
			s, event.EventName(), e)
	}
	return nil
}

// ToSnapshot implements es.Snapshotter
func (s Stream) ToSnapshot() es.Snapshot {
	return StreamV1{
		Title:              s.Title,
		Description:        s.Description,
		Synopsis:           s.Synopsis,
		StreamType:         s.StreamType,
		Status:             s.Status,
		StreamURL:          s.StreamURL,
		ThumbnailURL:       s.ThumbnailURL,
		TrailerURL:         s.TrailerURL,
		Duration:           s.Duration,
		ReleaseDate:        s.ReleaseDate,
		ContentRating:      s.ContentRating,
		AvailableQualities: s.AvailableQualities,
		DefaultQuality:     s.DefaultQuality,
		Subtitles:          s.Subtitles,
		AudioTracks:        s.AudioTracks,
		AccessType:         s.AccessType,
		SubscriptionTiers:  s.SubscriptionTiers,
		RentalPrice:        s.RentalPrice,
		RentalDuration:     s.RentalDuration,
		PurchasePrice:      s.PurchasePrice,
		PPVPrice:           s.PPVPrice,
		Genre:              s.Genre,
		Tags:               s.Tags,
		Cast:               s.Cast,
		Directors:          s.Directors,
		Producers:          s.Producers,
		Studio:             s.Studio,
		Language:           s.Language,
		Country:            s.Country,
		ViewCount:          s.ViewCount,
		LikeCount:          s.LikeCount,
		DislikeCount:       s.DislikeCount,
		AverageRating:      s.AverageRating,
		TotalRevenue:       s.TotalRevenue,
		UserAccess:         s.UserAccess,
		SeriesID:           s.SeriesID,
		SeasonNumber:       s.SeasonNumber,
		EpisodeNumber:      s.EpisodeNumber,
		CreatedAt:          s.CreatedAt,
		UpdatedAt:          s.UpdatedAt,
		PublishedAt:        s.PublishedAt,
	}
}

// ApplySnapshot implements es.Snapshotter
func (s *Stream) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *StreamV1:
		s.Title = ss.Title
		s.Description = ss.Description
		s.Synopsis = ss.Synopsis
		s.StreamType = ss.StreamType
		s.Status = ss.Status
		s.StreamURL = ss.StreamURL
		s.ThumbnailURL = ss.ThumbnailURL
		s.TrailerURL = ss.TrailerURL
		s.Duration = ss.Duration
		s.ReleaseDate = ss.ReleaseDate
		s.ContentRating = ss.ContentRating
		s.AvailableQualities = ss.AvailableQualities
		s.DefaultQuality = ss.DefaultQuality
		s.Subtitles = ss.Subtitles
		s.AudioTracks = ss.AudioTracks
		s.AccessType = ss.AccessType
		s.SubscriptionTiers = ss.SubscriptionTiers
		s.RentalPrice = ss.RentalPrice
		s.RentalDuration = ss.RentalDuration
		s.PurchasePrice = ss.PurchasePrice
		s.PPVPrice = ss.PPVPrice
		s.Genre = ss.Genre
		s.Tags = ss.Tags
		s.Cast = ss.Cast
		s.Directors = ss.Directors
		s.Producers = ss.Producers
		s.Studio = ss.Studio
		s.Language = ss.Language
		s.Country = ss.Country
		s.ViewCount = ss.ViewCount
		s.LikeCount = ss.LikeCount
		s.DislikeCount = ss.DislikeCount
		s.AverageRating = ss.AverageRating
		s.TotalRevenue = ss.TotalRevenue
		s.UserAccess = ss.UserAccess
		s.SeriesID = ss.SeriesID
		s.SeasonNumber = ss.SeasonNumber
		s.EpisodeNumber = ss.EpisodeNumber
		s.CreatedAt = ss.CreatedAt
		s.UpdatedAt = ss.UpdatedAt
		s.PublishedAt = ss.PublishedAt

	default:
		return errors.ErrInternal.Msgf(
			"%T received the unexpected snapshot %T", s, snapshot)
	}
	return nil
}