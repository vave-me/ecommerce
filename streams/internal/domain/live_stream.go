package domain

import (
	"middleman/internal/ddd"
	"middleman/internal/es"
	"time"

	"github.com/stackus/errors"
)

const LiveStreamAggregate = "streams.LiveStream"

var (
	ErrStreamNotLive            = errors.Wrap(errors.ErrBadRequest, "stream is not live")
	ErrStreamAlreadyLive        = errors.Wrap(errors.ErrBadRequest, "stream is already live")
	ErrInvalidStreamingProtocol = errors.Wrap(errors.ErrBadRequest, "invalid streaming protocol")
	ErrInvalidBitrate           = errors.Wrap(errors.ErrBadRequest, "invalid bitrate configuration")
	ErrMaxConcurrentViewers     = errors.Wrap(errors.ErrBadRequest, "maximum concurrent viewers reached")
	ErrStreamingServerError     = errors.Wrap(errors.ErrInternalServerError, "streaming server error")
	ErrInvalidDRMConfig         = errors.Wrap(errors.ErrBadRequest, "invalid DRM configuration")
	ErrGeoBlocked               = errors.Wrap(errors.ErrForbidden, "content is geo-blocked in your region")
)

// StreamingProtocol represents the streaming protocol used
type StreamingProtocol string

const (
	ProtocolHLS    StreamingProtocol = "HLS"     // HTTP Live Streaming (Apple)
	ProtocolDASH   StreamingProtocol = "DASH"    // Dynamic Adaptive Streaming over HTTP
	ProtocolWebRTC StreamingProtocol = "WebRTC"  // For ultra-low latency
	ProtocolRTMP   StreamingProtocol = "RTMP"    // Real-Time Messaging Protocol (input)
	ProtocolSRT    StreamingProtocol = "SRT"     // Secure Reliable Transport
)

// LiveStreamStatus represents the current status of a live stream
type LiveStreamStatus string

const (
	LiveStreamStatusScheduled   LiveStreamStatus = "scheduled"
	LiveStreamStatusStarting    LiveStreamStatus = "starting"
	LiveStreamStatusLive        LiveStreamStatus = "live"
	LiveStreamStatusEnding      LiveStreamStatus = "ending"
	LiveStreamStatusEnded       LiveStreamStatus = "ended"
	LiveStreamStatusError       LiveStreamStatus = "error"
	LiveStreamStatusCancelled   LiveStreamStatus = "cancelled"
)

// StreamingQualityProfile represents different quality profiles for adaptive streaming
type StreamingQualityProfile struct {
	Name       string
	Resolution string  // e.g., "1920x1080"
	Bitrate    int     // in kbps
	Framerate  int     // fps
	Codec      string  // e.g., "h264", "h265"
}

// CDNEndpoint represents a CDN endpoint for content delivery
type CDNEndpoint struct {
	Provider   string   // e.g., "Cloudflare", "Akamai", "Fastly"
	Region     string   // e.g., "eu-west-1", "us-east-1"
	URL        string   // CDN endpoint URL
	EdgeNodes  []string // List of edge node URLs
	Active     bool
}

// DRMConfig represents Digital Rights Management configuration
type DRMConfig struct {
	Provider      string            // e.g., "Widevine", "FairPlay", "PlayReady"
	LicenseURL    string
	CertificateURL string
	Headers       map[string]string
	Enabled       bool
}

// LiveStreamStatistics tracks real-time statistics
type LiveStreamStatistics struct {
	CurrentViewers     int64
	PeakViewers        int64
	TotalViews         int64
	AverageDuration    int64  // in seconds
	BufferingEvents    int64
	BitrateChanges     int64
	ErrorCount         int64
	BandwidthConsumed  int64  // in MB
	LastUpdated        time.Time
}

// LiveStream represents a live streaming event (e.g., football match)
type LiveStream struct {
	es.Aggregate
	
	// Basic Information
	Title              string
	Description        string
	EventType          string // "football_match", "basketball_game", etc.
	Status             LiveStreamStatus
	
	// Match Information (for sports events)
	HomeTeam           string
	AwayTeam           string
	Competition        string // e.g., "Premier League", "Champions League"
	Season             string
	MatchDay           int
	Stadium            string
	
	// Scheduling
	ScheduledStartTime time.Time
	ActualStartTime    time.Time
	ScheduledEndTime   time.Time
	ActualEndTime      time.Time
	PreShowDuration    int // minutes before match
	PostShowDuration   int // minutes after match
	
	// Streaming Configuration
	StreamingProtocols []StreamingProtocol
	QualityProfiles    []StreamingQualityProfile
	AdaptiveBitrate    bool
	LowLatencyMode     bool // For near real-time streaming
	DVREnabled         bool // Allow rewinding during live stream
	DVRWindowMinutes   int  // How far back users can rewind
	
	// CDN Configuration
	CDNEndpoints       []CDNEndpoint
	PrimaryCDN         string
	FailoverEnabled    bool
	
	// DRM and Security
	DRMConfigs         map[string]DRMConfig // Provider -> Config
	RequiresDRM        bool
	TokenAuthentication bool
	GeoRestrictions    []string // List of allowed country codes
	DRMKeyID           string   // Content key ID for DRM
	
	// Ingestion Settings
	IngestProtocol     StreamingProtocol
	IngestURL          string
	BackupIngestURL    string
	StreamKey          string
	
	// Playback URLs
	ManifestURLs       map[StreamingProtocol]string // Protocol -> URL
	PreviewImageURL    string
	ThumbnailURL       string
	
	// Interactive Features
	ChatEnabled        bool
	PollsEnabled       bool
	BettingIntegration bool
	StatsOverlay       bool
	MultiAngleViews    bool
	
	// Monetization
	RequiresSubscription bool
	SubscriptionTiers    []string
	PPVPrice            int64 // Pay-per-view price in cents
	AdInsertionEnabled  bool
	MidrollAdBreaks     []int // Timestamps for ad breaks
	
	// Analytics and Monitoring
	Statistics         LiveStreamStatistics
	QualityMetrics     map[string]float64 // Metric name -> value
	ViewerSessions     map[string]ViewerSession // SessionID -> Session
	
	// Technical Metadata
	MaxConcurrentViewers int
	MaxBitrate          int    // in kbps
	BufferSize          int    // in seconds
	ChunkDuration       int    // in seconds (for HLS/DASH)
	KeyFrameInterval    int    // in seconds
	
	// Recording
	RecordingEnabled    bool
	RecordingURL        string
	RecordingFormat     string // "mp4", "ts"
	
	// Timestamps
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// ViewerSession tracks individual viewer sessions
type ViewerSession struct {
	SessionID         string
	UserID            string
	StartTime         time.Time
	EndTime           time.Time
	Duration          int64 // in seconds
	Quality           string // Current quality being watched
	BufferingEvents   int
	QualityChanges    int
	AverageBitrate    int
	ClientInfo        ClientInfo
	CDNNode           string
	LastHeartbeat     time.Time
}

// ClientInfo contains viewer client information
type ClientInfo struct {
	UserAgent      string
	IP             string
	Country        string
	City           string
	ISP            string
	DeviceType     string // "mobile", "desktop", "tv", "tablet"
	OS             string
	Browser        string
	PlayerVersion  string
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*LiveStream)(nil)

func NewLiveStream(id string) *LiveStream {
	return &LiveStream{
		Aggregate:          es.NewAggregate(id, LiveStreamAggregate),
		StreamingProtocols: []StreamingProtocol{},
		QualityProfiles:    []StreamingQualityProfile{},
		CDNEndpoints:       []CDNEndpoint{},
		DRMConfigs:         make(map[string]DRMConfig),
		ManifestURLs:       make(map[StreamingProtocol]string),
		ViewerSessions:     make(map[string]ViewerSession),
		QualityMetrics:     make(map[string]float64),
		Statistics:         LiveStreamStatistics{},
	}
}

// Key implements registry.Registerable
func (LiveStream) Key() string { return LiveStreamAggregate }

// InitLiveStream initializes a new live stream for a sports event
func (ls *LiveStream) InitLiveStream(
	title, description, eventType string,
	homeTeam, awayTeam, competition, season string,
	matchDay int, stadium string,
	scheduledStart, scheduledEnd time.Time,
) (ddd.Event, error) {
	if title == "" {
		return nil, ErrStreamTitleIsBlank
	}
	if scheduledStart.After(scheduledEnd) {
		return nil, errors.Wrap(errors.ErrBadRequest, "invalid schedule times")
	}

	ls.AddEvent(LiveStreamCreatedEvent, &LiveStreamCreated{
		StreamID:           ls.ID(),
		Title:              title,
		Description:        description,
		EventType:          eventType,
		HomeTeam:           homeTeam,
		AwayTeam:           awayTeam,
		Competition:        competition,
		Season:             season,
		MatchDay:           matchDay,
		Stadium:            stadium,
		ScheduledStartTime: scheduledStart,
		ScheduledEndTime:   scheduledEnd,
		Status:             LiveStreamStatusScheduled,
		CreatedAt:          time.Now(),
	})
	return ddd.NewEvent(LiveStreamCreatedEvent, ls), nil
}

// ConfigureStreaming sets up the streaming configuration
func (ls *LiveStream) ConfigureStreaming(
	protocols []StreamingProtocol,
	qualityProfiles []StreamingQualityProfile,
	adaptiveBitrate, lowLatency, dvrEnabled bool,
	dvrWindow int,
) (ddd.Event, error) {
	if len(protocols) == 0 {
		return nil, ErrInvalidStreamingProtocol
	}
	if len(qualityProfiles) == 0 {
		return nil, ErrInvalidBitrate
	}

	ls.AddEvent(StreamingConfiguredEvent, &StreamingConfigured{
		StreamID:           ls.ID(),
		Protocols:          protocols,
		QualityProfiles:    qualityProfiles,
		AdaptiveBitrate:    adaptiveBitrate,
		LowLatencyMode:     lowLatency,
		DVREnabled:         dvrEnabled,
		DVRWindowMinutes:   dvrWindow,
	})
	return ddd.NewEvent(StreamingConfiguredEvent, ls), nil
}

// ConfigureCDN sets up CDN endpoints
func (ls *LiveStream) ConfigureCDN(
	endpoints []CDNEndpoint,
	primaryCDN string,
	failoverEnabled bool,
) (ddd.Event, error) {
	if len(endpoints) == 0 {
		return nil, errors.Wrap(errors.ErrBadRequest, "at least one CDN endpoint required")
	}

	// Verify primary CDN exists in endpoints
	found := false
	for _, ep := range endpoints {
		if ep.Provider == primaryCDN {
			found = true
			break
		}
	}
	if !found {
		return nil, errors.Wrap(errors.ErrBadRequest, "primary CDN not found in endpoints")
	}

	ls.AddEvent(CDNConfiguredEvent, &CDNConfigured{
		StreamID:        ls.ID(),
		CDNEndpoints:    endpoints,
		PrimaryCDN:      primaryCDN,
		FailoverEnabled: failoverEnabled,
	})
	return ddd.NewEvent(CDNConfiguredEvent, ls), nil
}

// ConfigureDRM sets up Digital Rights Management
func (ls *LiveStream) ConfigureDRM(configs map[string]DRMConfig, required bool) (ddd.Event, error) {
	if required && len(configs) == 0 {
		return nil, ErrInvalidDRMConfig
	}

	ls.AddEvent(DRMConfiguredEvent, &DRMConfigured{
		StreamID:    ls.ID(),
		DRMConfigs:  configs,
		RequiresDRM: required,
	})
	return ddd.NewEvent(DRMConfiguredEvent, ls), nil
}

// SetIngestionConfig configures the stream ingestion
func (ls *LiveStream) SetIngestionConfig(
	protocol StreamingProtocol,
	ingestURL, backupURL, streamKey string,
) (ddd.Event, error) {
	if ingestURL == "" || streamKey == "" {
		return nil, errors.Wrap(errors.ErrBadRequest, "invalid ingestion configuration")
	}

	ls.AddEvent(IngestionConfiguredEvent, &IngestionConfigured{
		StreamID:        ls.ID(),
		Protocol:        protocol,
		IngestURL:       ingestURL,
		BackupIngestURL: backupURL,
		StreamKey:       streamKey,
	})
	return ddd.NewEvent(IngestionConfiguredEvent, ls), nil
}

// StartStream transitions the stream to live status
func (ls *LiveStream) StartStream() (ddd.Event, error) {
	if ls.Status == LiveStreamStatusLive {
		return nil, ErrStreamAlreadyLive
	}
	if ls.Status != LiveStreamStatusScheduled && ls.Status != LiveStreamStatusStarting {
		return nil, errors.Wrap(errors.ErrBadRequest, "stream cannot be started from current status")
	}

	ls.AddEvent(StreamStartedEvent, &StreamStarted{
		StreamID:        ls.ID(),
		ActualStartTime: time.Now(),
	})
	return ddd.NewEvent(StreamStartedEvent, ls), nil
}

// EndStream transitions the stream to ended status
func (ls *LiveStream) EndStream() (ddd.Event, error) {
	if ls.Status != LiveStreamStatusLive {
		return nil, ErrStreamNotLive
	}

	ls.AddEvent(StreamEndedEvent, &StreamEnded{
		StreamID:      ls.ID(),
		ActualEndTime: time.Now(),
		FinalStats:    ls.Statistics,
	})
	return ddd.NewEvent(StreamEndedEvent, ls), nil
}

// AddViewerSession records a new viewer joining the stream
func (ls *LiveStream) AddViewerSession(
	sessionID, userID string,
	clientInfo ClientInfo,
	cdnNode string,
) (ddd.Event, error) {
	if ls.Status != LiveStreamStatusLive {
		return nil, ErrStreamNotLive
	}
	if ls.MaxConcurrentViewers > 0 && int(ls.Statistics.CurrentViewers) >= ls.MaxConcurrentViewers {
		return nil, ErrMaxConcurrentViewers
	}

	ls.AddEvent(ViewerJoinedEvent, &ViewerJoined{
		StreamID:   ls.ID(),
		SessionID:  sessionID,
		UserID:     userID,
		JoinTime:   time.Now(),
		ClientInfo: clientInfo,
		CDNNode:    cdnNode,
	})
	return ddd.NewEvent(ViewerJoinedEvent, ls), nil
}

// RemoveViewerSession records a viewer leaving the stream
func (ls *LiveStream) RemoveViewerSession(sessionID string) (ddd.Event, error) {
	session, exists := ls.ViewerSessions[sessionID]
	if !exists {
		return nil, errors.Wrap(errors.ErrNotFound, "viewer session not found")
	}

	duration := int64(time.Since(session.StartTime).Seconds())
	
	ls.AddEvent(ViewerLeftEvent, &ViewerLeft{
		StreamID:        ls.ID(),
		SessionID:       sessionID,
		LeaveTime:       time.Now(),
		Duration:        duration,
		QualityChanges:  session.QualityChanges,
		BufferingEvents: session.BufferingEvents,
	})
	return ddd.NewEvent(ViewerLeftEvent, ls), nil
}

// UpdateStatistics updates the stream statistics
func (ls *LiveStream) UpdateStatistics(stats LiveStreamStatistics) (ddd.Event, error) {
	ls.AddEvent(StatisticsUpdatedEvent, &StatisticsUpdated{
		StreamID:   ls.ID(),
		Statistics: stats,
		UpdatedAt:  time.Now(),
	})
	return ddd.NewEvent(StatisticsUpdatedEvent, ls), nil
}

// RecordQualityMetric records a quality metric for monitoring
func (ls *LiveStream) RecordQualityMetric(metric string, value float64) (ddd.Event, error) {
	ls.AddEvent(QualityMetricRecordedEvent, &QualityMetricRecorded{
		StreamID:  ls.ID(),
		Metric:    metric,
		Value:     value,
		Timestamp: time.Now(),
	})
	return ddd.NewEvent(QualityMetricRecordedEvent, ls), nil
}

// SetManifestURL sets the manifest URL for a specific protocol
func (ls *LiveStream) SetManifestURL(protocol StreamingProtocol, url string) (ddd.Event, error) {
	if url == "" {
		return nil, errors.Wrap(errors.ErrBadRequest, "manifest URL cannot be empty")
	}

	ls.AddEvent(ManifestURLSetEvent, &ManifestURLSet{
		StreamID: ls.ID(),
		Protocol: protocol,
		URL:      url,
	})
	return ddd.NewEvent(ManifestURLSetEvent, ls), nil
}

// EnableRecording enables stream recording
func (ls *LiveStream) EnableRecording(url, format string) (ddd.Event, error) {
	if ls.Status != LiveStreamStatusScheduled && ls.Status != LiveStreamStatusStarting {
		return nil, errors.Wrap(errors.ErrBadRequest, "recording must be enabled before stream starts")
	}

	ls.AddEvent(RecordingEnabledEvent, &RecordingEnabled{
		StreamID:        ls.ID(),
		RecordingURL:    url,
		RecordingFormat: format,
	})
	return ddd.NewEvent(RecordingEnabledEvent, ls), nil
}

// ApplyEvent implements es.EventApplier
func (ls *LiveStream) ApplyEvent(event ddd.Event) error {
	switch e := event.Payload().(type) {
	case *LiveStreamCreated:
		ls.Title = e.Title
		ls.Description = e.Description
		ls.EventType = e.EventType
		ls.HomeTeam = e.HomeTeam
		ls.AwayTeam = e.AwayTeam
		ls.Competition = e.Competition
		ls.Season = e.Season
		ls.MatchDay = e.MatchDay
		ls.Stadium = e.Stadium
		ls.ScheduledStartTime = e.ScheduledStartTime
		ls.ScheduledEndTime = e.ScheduledEndTime
		ls.Status = e.Status
		ls.CreatedAt = e.CreatedAt
		ls.UpdatedAt = e.CreatedAt

	case *StreamingConfigured:
		ls.StreamingProtocols = e.Protocols
		ls.QualityProfiles = e.QualityProfiles
		ls.AdaptiveBitrate = e.AdaptiveBitrate
		ls.LowLatencyMode = e.LowLatencyMode
		ls.DVREnabled = e.DVREnabled
		ls.DVRWindowMinutes = e.DVRWindowMinutes

	case *CDNConfigured:
		ls.CDNEndpoints = e.CDNEndpoints
		ls.PrimaryCDN = e.PrimaryCDN
		ls.FailoverEnabled = e.FailoverEnabled

	case *DRMConfigured:
		ls.DRMConfigs = e.DRMConfigs
		ls.RequiresDRM = e.RequiresDRM

	case *IngestionConfigured:
		ls.IngestProtocol = e.Protocol
		ls.IngestURL = e.IngestURL
		ls.BackupIngestURL = e.BackupIngestURL
		ls.StreamKey = e.StreamKey

	case *StreamStarted:
		ls.Status = LiveStreamStatusLive
		ls.ActualStartTime = e.ActualStartTime

	case *StreamEnded:
		ls.Status = LiveStreamStatusEnded
		ls.ActualEndTime = e.ActualEndTime
		ls.Statistics = e.FinalStats

	case *ViewerJoined:
		ls.ViewerSessions[e.SessionID] = ViewerSession{
			SessionID:  e.SessionID,
			UserID:     e.UserID,
			StartTime:  e.JoinTime,
			ClientInfo: e.ClientInfo,
			CDNNode:    e.CDNNode,
		}
		ls.Statistics.CurrentViewers++
		ls.Statistics.TotalViews++
		if ls.Statistics.CurrentViewers > ls.Statistics.PeakViewers {
			ls.Statistics.PeakViewers = ls.Statistics.CurrentViewers
		}

	case *ViewerLeft:
		delete(ls.ViewerSessions, e.SessionID)
		ls.Statistics.CurrentViewers--

	case *StatisticsUpdated:
		ls.Statistics = e.Statistics

	case *QualityMetricRecorded:
		ls.QualityMetrics[e.Metric] = e.Value

	case *ManifestURLSet:
		ls.ManifestURLs[e.Protocol] = e.URL

	case *RecordingEnabled:
		ls.RecordingEnabled = true
		ls.RecordingURL = e.RecordingURL
		ls.RecordingFormat = e.RecordingFormat

	default:
		return errors.ErrInternal.Msgf(
			"%T received the event %s with unexpected payload %T",
			ls, event.EventName(), e)
	}
	return nil
}

// ToSnapshot implements es.Snapshotter
func (ls LiveStream) ToSnapshot() es.Snapshot {
	return LiveStreamV1{
		Title:                ls.Title,
		Description:          ls.Description,
		EventType:            ls.EventType,
		Status:               ls.Status,
		HomeTeam:             ls.HomeTeam,
		AwayTeam:             ls.AwayTeam,
		Competition:          ls.Competition,
		Season:               ls.Season,
		MatchDay:             ls.MatchDay,
		Stadium:              ls.Stadium,
		ScheduledStartTime:   ls.ScheduledStartTime,
		ActualStartTime:      ls.ActualStartTime,
		ScheduledEndTime:     ls.ScheduledEndTime,
		ActualEndTime:        ls.ActualEndTime,
		PreShowDuration:      ls.PreShowDuration,
		PostShowDuration:     ls.PostShowDuration,
		StreamingProtocols:   ls.StreamingProtocols,
		QualityProfiles:      ls.QualityProfiles,
		AdaptiveBitrate:      ls.AdaptiveBitrate,
		LowLatencyMode:       ls.LowLatencyMode,
		DVREnabled:           ls.DVREnabled,
		DVRWindowMinutes:     ls.DVRWindowMinutes,
		CDNEndpoints:         ls.CDNEndpoints,
		PrimaryCDN:           ls.PrimaryCDN,
		FailoverEnabled:      ls.FailoverEnabled,
		DRMConfigs:           ls.DRMConfigs,
		RequiresDRM:          ls.RequiresDRM,
		TokenAuthentication:  ls.TokenAuthentication,
		GeoRestrictions:      ls.GeoRestrictions,
		IngestProtocol:       ls.IngestProtocol,
		IngestURL:            ls.IngestURL,
		BackupIngestURL:      ls.BackupIngestURL,
		StreamKey:            ls.StreamKey,
		ManifestURLs:         ls.ManifestURLs,
		PreviewImageURL:      ls.PreviewImageURL,
		ThumbnailURL:         ls.ThumbnailURL,
		ChatEnabled:          ls.ChatEnabled,
		PollsEnabled:         ls.PollsEnabled,
		BettingIntegration:   ls.BettingIntegration,
		StatsOverlay:         ls.StatsOverlay,
		MultiAngleViews:      ls.MultiAngleViews,
		RequiresSubscription: ls.RequiresSubscription,
		SubscriptionTiers:    ls.SubscriptionTiers,
		PPVPrice:             ls.PPVPrice,
		AdInsertionEnabled:   ls.AdInsertionEnabled,
		MidrollAdBreaks:      ls.MidrollAdBreaks,
		Statistics:           ls.Statistics,
		QualityMetrics:       ls.QualityMetrics,
		ViewerSessions:       ls.ViewerSessions,
		MaxConcurrentViewers: ls.MaxConcurrentViewers,
		MaxBitrate:           ls.MaxBitrate,
		BufferSize:           ls.BufferSize,
		ChunkDuration:        ls.ChunkDuration,
		KeyFrameInterval:     ls.KeyFrameInterval,
		RecordingEnabled:     ls.RecordingEnabled,
		RecordingURL:         ls.RecordingURL,
		RecordingFormat:      ls.RecordingFormat,
		CreatedAt:            ls.CreatedAt,
		UpdatedAt:            ls.UpdatedAt,
	}
}

// ApplySnapshot implements es.Snapshotter
func (ls *LiveStream) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *LiveStreamV1:
		ls.Title = ss.Title
		ls.Description = ss.Description
		ls.EventType = ss.EventType
		ls.Status = ss.Status
		ls.HomeTeam = ss.HomeTeam
		ls.AwayTeam = ss.AwayTeam
		ls.Competition = ss.Competition
		ls.Season = ss.Season
		ls.MatchDay = ss.MatchDay
		ls.Stadium = ss.Stadium
		ls.ScheduledStartTime = ss.ScheduledStartTime
		ls.ActualStartTime = ss.ActualStartTime
		ls.ScheduledEndTime = ss.ScheduledEndTime
		ls.ActualEndTime = ss.ActualEndTime
		ls.PreShowDuration = ss.PreShowDuration
		ls.PostShowDuration = ss.PostShowDuration
		ls.StreamingProtocols = ss.StreamingProtocols
		ls.QualityProfiles = ss.QualityProfiles
		ls.AdaptiveBitrate = ss.AdaptiveBitrate
		ls.LowLatencyMode = ss.LowLatencyMode
		ls.DVREnabled = ss.DVREnabled
		ls.DVRWindowMinutes = ss.DVRWindowMinutes
		ls.CDNEndpoints = ss.CDNEndpoints
		ls.PrimaryCDN = ss.PrimaryCDN
		ls.FailoverEnabled = ss.FailoverEnabled
		ls.DRMConfigs = ss.DRMConfigs
		ls.RequiresDRM = ss.RequiresDRM
		ls.TokenAuthentication = ss.TokenAuthentication
		ls.GeoRestrictions = ss.GeoRestrictions
		ls.IngestProtocol = ss.IngestProtocol
		ls.IngestURL = ss.IngestURL
		ls.BackupIngestURL = ss.BackupIngestURL
		ls.StreamKey = ss.StreamKey
		ls.ManifestURLs = ss.ManifestURLs
		ls.PreviewImageURL = ss.PreviewImageURL
		ls.ThumbnailURL = ss.ThumbnailURL
		ls.ChatEnabled = ss.ChatEnabled
		ls.PollsEnabled = ss.PollsEnabled
		ls.BettingIntegration = ss.BettingIntegration
		ls.StatsOverlay = ss.StatsOverlay
		ls.MultiAngleViews = ss.MultiAngleViews
		ls.RequiresSubscription = ss.RequiresSubscription
		ls.SubscriptionTiers = ss.SubscriptionTiers
		ls.PPVPrice = ss.PPVPrice
		ls.AdInsertionEnabled = ss.AdInsertionEnabled
		ls.MidrollAdBreaks = ss.MidrollAdBreaks
		ls.Statistics = ss.Statistics
		ls.QualityMetrics = ss.QualityMetrics
		ls.ViewerSessions = ss.ViewerSessions
		ls.MaxConcurrentViewers = ss.MaxConcurrentViewers
		ls.MaxBitrate = ss.MaxBitrate
		ls.BufferSize = ss.BufferSize
		ls.ChunkDuration = ss.ChunkDuration
		ls.KeyFrameInterval = ss.KeyFrameInterval
		ls.RecordingEnabled = ss.RecordingEnabled
		ls.RecordingURL = ss.RecordingURL
		ls.RecordingFormat = ss.RecordingFormat
		ls.CreatedAt = ss.CreatedAt
		ls.UpdatedAt = ss.UpdatedAt

	default:
		return errors.ErrInternal.Msgf(
			"%T received the unexpected snapshot %T", ls, snapshot)
	}
	return nil
}

// HasAccess checks if a user has access to the stream
func (ls *LiveStream) HasAccess(userID string, userTiers []string) bool {
	// Free streams
	if !ls.RequiresSubscription && ls.PPVPrice == 0 {
		return true
	}

	// Check subscription tiers
	if ls.RequiresSubscription {
		for _, required := range ls.SubscriptionTiers {
			for _, userTier := range userTiers {
				if required == userTier {
					return true
				}
			}
		}
	}

	// PPV would need separate purchase check
	return false
}

// IsGeoBlocked checks if the stream is blocked in a country
func (ls *LiveStream) IsGeoBlocked(countryCode string) bool {
	if len(ls.GeoRestrictions) == 0 {
		return false // No restrictions
	}
	
	for _, allowed := range ls.GeoRestrictions {
		if allowed == countryCode {
			return false
		}
	}
	return true
}