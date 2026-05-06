package domain

import (
	"time"
)

// WebhookSubscription represents a webhook subscription
type WebhookSubscription struct {
	ID          string
	URL         string
	Secret      string
	Events      []string
	Headers     map[string]string
	RetryPolicy RetryPolicy
	Active      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// RetryPolicy defines retry behavior for webhook deliveries
type RetryPolicy struct {
	MaxRetries    int           `json:"max_retries"`
	BackoffFactor float64       `json:"backoff_factor"`
	InitialDelay  time.Duration `json:"initial_delay"`
	MaxBackoff    time.Duration `json:"max_backoff"`
}

// WebhookDelivery represents a webhook delivery attempt
type WebhookDelivery struct {
	ID             string
	SubscriptionID string
	EventID        string
	EventType      string
	Payload        []byte
	Signature      string
	Status         DeliveryStatus
	Attempts       int
	NextRetryAt    *time.Time
	ResponseStatus int
	ResponseBody   string
	Error          string
	CreatedAt      time.Time
	LastAttemptAt  *time.Time
	CompletedAt    *time.Time
}

// DeliveryStatus represents the status of a webhook delivery
type DeliveryStatus string

const (
	DeliveryStatusPending   DeliveryStatus = "pending"
	DeliveryStatusDelivered DeliveryStatus = "delivered"
	DeliveryStatusFailed    DeliveryStatus = "failed"
	DeliveryStatusRetrying  DeliveryStatus = "retrying"
)

// WebhookEvent represents an event to be delivered via webhook
type WebhookEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	StreamID  string                 `json:"stream_id"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// WebhookEventType defines supported webhook event types
type WebhookEventType string

const (
	// Stream lifecycle events
	WebhookEventStreamCreated   WebhookEventType = "stream.created"
	WebhookEventStreamStarted   WebhookEventType = "stream.started"
	WebhookEventStreamEnded     WebhookEventType = "stream.ended"
	WebhookEventStreamError     WebhookEventType = "stream.error"
	WebhookEventStreamArchived  WebhookEventType = "stream.archived"
	
	// Viewer events
	WebhookEventViewerJoined    WebhookEventType = "viewer.joined"
	WebhookEventViewerLeft      WebhookEventType = "viewer.left"
	
	// Quality events
	WebhookEventQualityDegraded WebhookEventType = "quality.degraded"
	WebhookEventBufferingHigh   WebhookEventType = "buffering.high"
	
	// CDN events
	WebhookEventCDNFailover     WebhookEventType = "cdn.failover"
	
	// Recording events
	WebhookEventRecordingStarted   WebhookEventType = "recording.started"
	WebhookEventRecordingCompleted WebhookEventType = "recording.completed"
	
	// Access control events
	WebhookEventAccessGranted   WebhookEventType = "access.granted"
	WebhookEventAccessRevoked   WebhookEventType = "access.revoked"
)

// WebhookSubscriptionRepository defines the interface for webhook subscription persistence
type WebhookSubscriptionRepository interface {
	Create(subscription *WebhookSubscription) error
	Update(subscription *WebhookSubscription) error
	Delete(id string) error
	Find(id string) (*WebhookSubscription, error)
	FindAll() ([]*WebhookSubscription, error)
	FindActiveByEvent(eventType string) ([]*WebhookSubscription, error)
}

// WebhookDeliveryRepository defines the interface for webhook delivery persistence
type WebhookDeliveryRepository interface {
	Create(delivery *WebhookDelivery) error
	Update(delivery *WebhookDelivery) error
	Find(id string) (*WebhookDelivery, error)
	FindBySubscription(subscriptionID string) ([]*WebhookDelivery, error)
	FindPendingDeliveries(limit int) ([]*WebhookDelivery, error)
}

// MapDomainEventToWebhookEvent maps domain events to webhook event types
func MapDomainEventToWebhookEvent(domainEvent string) WebhookEventType {
	mappings := map[string]WebhookEventType{
		LiveStreamCreatedEvent:   WebhookEventStreamCreated,
		StreamStartedEvent:       WebhookEventStreamStarted,
		StreamEndedEvent:         WebhookEventStreamEnded,
		StreamErrorEvent:         WebhookEventStreamError,
		StreamArchivedEvent:      WebhookEventStreamArchived,
		ViewerJoinedEvent:        WebhookEventViewerJoined,
		ViewerLeftEvent:          WebhookEventViewerLeft,
		CDNFailoverEvent:         WebhookEventCDNFailover,
		RecordingEnabledEvent:    WebhookEventRecordingStarted,
		UserAccessGrantedEvent:   WebhookEventAccessGranted,
		UserAccessRevokedEvent:   WebhookEventAccessRevoked,
	}
	
	if eventType, exists := mappings[domainEvent]; exists {
		return eventType
	}
	return ""
}