package domain

import (
	"time"
)

// Event names for live streaming
const (
	LiveStreamCreatedEvent        = "streams.LiveStreamCreated"
	StreamingConfiguredEvent      = "streams.StreamingConfigured"
	CDNConfiguredEvent            = "streams.CDNConfigured"
	DRMConfiguredEvent            = "streams.DRMConfigured"
	IngestionConfiguredEvent      = "streams.IngestionConfigured"
	StreamStartedEvent            = "streams.StreamStarted"
	StreamEndedEvent              = "streams.StreamEnded"
	ViewerJoinedEvent             = "streams.ViewerJoined"
	ViewerLeftEvent               = "streams.ViewerLeft"
	StatisticsUpdatedEvent        = "streams.StatisticsUpdated"
	QualityMetricRecordedEvent    = "streams.QualityMetricRecorded"
	ManifestURLSetEvent           = "streams.ManifestURLSet"
	RecordingEnabledEvent         = "streams.RecordingEnabled"
	StreamErrorEvent              = "streams.StreamError"
	CDNFailoverEvent              = "streams.CDNFailover"
	BitrateAdaptationEvent        = "streams.BitrateAdaptation"
	AdBreakStartedEvent           = "streams.AdBreakStarted"
	AdBreakEndedEvent             = "streams.AdBreakEnded"
)

// LiveStreamCreated is raised when a new live stream is created
type LiveStreamCreated struct {
	StreamID           string
	Title              string
	Description        string
	EventType          string
	HomeTeam           string
	AwayTeam           string
	Competition        string
	Season             string
	MatchDay           int
	Stadium            string
	ScheduledStartTime time.Time
	ScheduledEndTime   time.Time
	Status             LiveStreamStatus
	CreatedAt          time.Time
}

// StreamingConfigured is raised when streaming configuration is set
type StreamingConfigured struct {
	StreamID           string
	Protocols          []StreamingProtocol
	QualityProfiles    []StreamingQualityProfile
	AdaptiveBitrate    bool
	LowLatencyMode     bool
	DVREnabled         bool
	DVRWindowMinutes   int
}

// CDNConfigured is raised when CDN configuration is set
type CDNConfigured struct {
	StreamID        string
	CDNEndpoints    []CDNEndpoint
	PrimaryCDN      string
	FailoverEnabled bool
}

// DRMConfigured is raised when DRM configuration is set
type DRMConfigured struct {
	StreamID    string
	DRMConfigs  map[string]DRMConfig
	RequiresDRM bool
}

// IngestionConfigured is raised when ingestion configuration is set
type IngestionConfigured struct {
	StreamID        string
	Protocol        StreamingProtocol
	IngestURL       string
	BackupIngestURL string
	StreamKey       string
}

// StreamStarted is raised when a stream goes live
type StreamStarted struct {
	StreamID        string
	ActualStartTime time.Time
}

// StreamEnded is raised when a stream ends
type StreamEnded struct {
	StreamID      string
	ActualEndTime time.Time
	FinalStats    LiveStreamStatistics
}

// ViewerJoined is raised when a viewer joins the stream
type ViewerJoined struct {
	StreamID   string
	SessionID  string
	UserID     string
	JoinTime   time.Time
	ClientInfo ClientInfo
	CDNNode    string
}

// ViewerLeft is raised when a viewer leaves the stream
type ViewerLeft struct {
	StreamID        string
	SessionID       string
	LeaveTime       time.Time
	Duration        int64 // seconds
	QualityChanges  int
	BufferingEvents int
}

// StatisticsUpdated is raised when stream statistics are updated
type StatisticsUpdated struct {
	StreamID   string
	Statistics LiveStreamStatistics
	UpdatedAt  time.Time
}

// QualityMetricRecorded is raised when a quality metric is recorded
type QualityMetricRecorded struct {
	StreamID  string
	Metric    string
	Value     float64
	Timestamp time.Time
}

// ManifestURLSet is raised when a manifest URL is set for a protocol
type ManifestURLSet struct {
	StreamID string
	Protocol StreamingProtocol
	URL      string
}

// RecordingEnabled is raised when recording is enabled for a stream
type RecordingEnabled struct {
	StreamID        string
	RecordingURL    string
	RecordingFormat string
}

// StreamError is raised when a stream encounters an error
type StreamError struct {
	StreamID    string
	ErrorType   string
	ErrorMessage string
	Severity    string // "warning", "error", "critical"
	Timestamp   time.Time
}

// CDNFailover is raised when CDN failover occurs
type CDNFailover struct {
	StreamID    string
	FromCDN     string
	ToCDN       string
	Reason      string
	Timestamp   time.Time
}

// BitrateAdaptation is raised when a viewer's bitrate is adapted
type BitrateAdaptation struct {
	StreamID      string
	SessionID     string
	FromBitrate   int
	ToBitrate     int
	FromQuality   string
	ToQuality     string
	Reason        string // "network_congestion", "bandwidth_available", etc.
	Timestamp     time.Time
}

// AdBreakStarted is raised when an ad break starts
type AdBreakStarted struct {
	StreamID  string
	AdBreakID string
	Duration  int // seconds
	Timestamp time.Time
}

// AdBreakEnded is raised when an ad break ends
type AdBreakEnded struct {
	StreamID  string
	AdBreakID string
	Timestamp time.Time
}