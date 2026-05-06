package application

import (
	"context"
	"time"
	
	"middleman/internal/ddd"
	"middleman/sap/internal/domain"
	"middleman/sap/internal/sap"
	"middleman/sap/internal/sap/transformer"
)

// SAPConnectorDomain defines the interface for SAP connector operations
type SAPConnectorDomain interface {
	// Webhook event storage
	StoreWebhookEvent(ctx context.Context, event *domain.WebhookEvent) error
	UpdateWebhookEventStatus(ctx context.Context, id, status, errorMessage string) error
	
	// IDoc processing
	ProcessMaterialMaster(ctx context.Context, cmd ProcessMaterialMasterCommand) error
	ProcessInventoryUpdate(ctx context.Context, cmd ProcessInventoryUpdateCommand) error
	ProcessPricingUpdate(ctx context.Context, cmd ProcessPricingUpdateCommand) error
	
	// JSON event processing
	ProcessProductEvent(ctx context.Context, cmd ProcessProductEventCommand) error
	ProcessStockEvent(ctx context.Context, cmd ProcessStockEventCommand) error
	ProcessPriceEvent(ctx context.Context, cmd ProcessPriceEventCommand) error
	ProcessProductDeletedEvent(ctx context.Context, cmd ProcessProductDeletedEventCommand) error
	
	// Sync operations
	SyncProductsFromSAP(ctx context.Context, cmd SyncProductsCommand) error
	SyncStockLevelsFromSAP(ctx context.Context, cmd SyncStockLevelsCommand) error
	SyncPricesFromSAP(ctx context.Context, cmd SyncPricesCommand) error
	
	// Query operations
	GetSyncStatus(ctx context.Context, entityType, entityID string) (*domain.SyncStatus, error)
	GetRecentSyncLogs(ctx context.Context, limit int) ([]*domain.SyncLog, error)
	GetSyncConfiguration(ctx context.Context, entityType string) (*domain.SyncConfiguration, error)
}

// Commands

type ProcessMaterialMasterCommand struct {
	IDocData       interface{}
	CorrelationID  string
	WebhookEventID string
}

type ProcessInventoryUpdateCommand struct {
	IDocData       interface{}
	CorrelationID  string
	WebhookEventID string
}

type ProcessPricingUpdateCommand struct {
	IDocData       interface{}
	CorrelationID  string
	WebhookEventID string
}

type ProcessProductEventCommand struct {
	Event          *sap.SAPEvent
	WebhookEventID string
}

type ProcessStockEventCommand struct {
	Event          *sap.SAPEvent
	WebhookEventID string
}

type ProcessPriceEventCommand struct {
	Event          *sap.SAPEvent
	WebhookEventID string
}

type ProcessProductDeletedEventCommand struct {
	Event          *sap.SAPEvent
	WebhookEventID string
}

type SyncProductsCommand struct {
	Since *time.Time
}

type SyncStockLevelsCommand struct {
	ProductIDs []string
}

type SyncPricesCommand struct {
	ProductIDs  []string
	PriceListID string
}

// Application implements the SAP connector domain
type Application struct {
	sapClient            *sap.SAPClient
	eventPublisher       ddd.EventPublisher[ddd.Event]
	syncStatusRepo       domain.SyncStatusRepository
	syncLogRepo          domain.SyncLogRepository
	syncConfigRepo       domain.SyncConfigurationRepository
	webhookEventRepo     domain.WebhookEventRepository
	transformer          *transformer.SAPToCanonicalTransformer
}

// NewApplication creates a new SAP connector application
func NewApplication(
	sapClient *sap.SAPClient,
	eventPublisher ddd.EventPublisher[ddd.Event],
	syncStatusRepo domain.SyncStatusRepository,
	syncLogRepo domain.SyncLogRepository,
	syncConfigRepo domain.SyncConfigurationRepository,
	webhookEventRepo domain.WebhookEventRepository,
) *Application {
	return &Application{
		sapClient:        sapClient,
		eventPublisher:   eventPublisher,
		syncStatusRepo:   syncStatusRepo,
		syncLogRepo:      syncLogRepo,
		syncConfigRepo:   syncConfigRepo,
		webhookEventRepo: webhookEventRepo,
		transformer:      transformer.NewSAPToCanonicalTransformer("SAP"),
	}
}

// Ensure Application implements SAPConnectorDomain
var _ SAPConnectorDomain = (*Application)(nil)