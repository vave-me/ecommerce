package domain

import (
	"context"
	"time"
)

// ConnectorEntity represents a connector configuration in the database
type ConnectorEntity struct {
	ID                         string
	Name                       string
	Type                       string
	Status                     string
	AuthConfigEncrypted        []byte
	AuthConfigSalt             string
	BaseURL                    string
	Environment                string
	WebhookEnabled             bool
	WebhookSecretEncrypted     []byte
	WebhookURL                 string
	WebhookEvents              []string
	SyncEnabled                bool
	SyncIntervalSeconds        int
	BatchSize                  int
	RateLimitRequestsPerSecond int
	RateLimitBurst             int
	RetryMaxAttempts           int
	RetryInitialDelayMs        int
	RetryMaxDelayMs            int
	RetryMultiplier            float64
	CustomHeaders              map[string]string
	TimeoutSeconds             int
	LastHealthCheckAt          *time.Time
	LastHealthCheckStatus      string
	LastHealthCheckError       string
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
	CreatedBy                  string
	UpdatedBy                  string
	Version                    int
}

// ConnectorSyncEntity represents sync configuration for a specific entity type
type ConnectorSyncEntity struct {
	ID            string
	ConnectorID   string
	EntityType    string
	Enabled       bool
	SyncDirection string
	LastSyncAt    *time.Time
	Filters       map[string]interface{}
	FieldMapping  map[string]string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// ConnectorAPIKey represents an API key for a connector
type ConnectorAPIKey struct {
	ID                string
	ConnectorID       string
	KeyName           string
	KeyValueEncrypted []byte
	KeySalt           string
	KeyType           string
	ExpiresAt         *time.Time
	LastUsedAt        *time.Time
	CreatedAt         time.Time
}

// ConnectorAuditLog represents an audit log entry for connector changes
type ConnectorAuditLog struct {
	ID          string
	ConnectorID string
	Action      string
	ChangedBy   string
	ChangedAt   time.Time
	OldValues   map[string]interface{}
	NewValues   map[string]interface{}
	IPAddress   string
	UserAgent   string
}

// ConnectorRepository defines the interface for connector persistence
type ConnectorRepository interface {
	// Connector CRUD operations
	Create(ctx context.Context, connector *ConnectorEntity) error
	GetByID(ctx context.Context, id string) (*ConnectorEntity, error)
	GetByName(ctx context.Context, name string) (*ConnectorEntity, error)
	GetByType(ctx context.Context, connectorType string) ([]*ConnectorEntity, error)
	GetAll(ctx context.Context) ([]*ConnectorEntity, error)
	GetActive(ctx context.Context) ([]*ConnectorEntity, error)
	Update(ctx context.Context, connector *ConnectorEntity) error
	Delete(ctx context.Context, id string) error
	
	// Status operations
	UpdateStatus(ctx context.Context, id string, status string) error
	UpdateHealthCheck(ctx context.Context, id string, status string, errorMsg string) error
	
	// Sync entity operations
	CreateSyncEntity(ctx context.Context, entity *ConnectorSyncEntity) error
	GetSyncEntities(ctx context.Context, connectorID string) ([]*ConnectorSyncEntity, error)
	UpdateSyncEntity(ctx context.Context, entity *ConnectorSyncEntity) error
	DeleteSyncEntity(ctx context.Context, id string) error
	UpdateLastSyncTime(ctx context.Context, connectorID string, entityType string, syncTime time.Time) error
	
	// API key operations
	CreateAPIKey(ctx context.Context, key *ConnectorAPIKey) error
	GetAPIKeys(ctx context.Context, connectorID string) ([]*ConnectorAPIKey, error)
	GetAPIKeyByName(ctx context.Context, connectorID string, keyName string) (*ConnectorAPIKey, error)
	UpdateAPIKeyLastUsed(ctx context.Context, id string) error
	DeleteAPIKey(ctx context.Context, id string) error
	DeleteExpiredAPIKeys(ctx context.Context) error
	
	// Audit operations
	CreateAuditLog(ctx context.Context, log *ConnectorAuditLog) error
	GetAuditLogs(ctx context.Context, connectorID string, limit int) ([]*ConnectorAuditLog, error)
	GetAuditLogsByAction(ctx context.Context, action string, since time.Time) ([]*ConnectorAuditLog, error)
}