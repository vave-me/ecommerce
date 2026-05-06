package application

import (
	"context"
	
	"middleman/internal/ddd"
	"middleman/sap/internal/application/commands"
	"middleman/sap/internal/domain"
	"middleman/sap/internal/sap"
	"middleman/sap/internal/sap/transformer"
)

// Implementation methods for Application

// StoreWebhookEvent stores a webhook event
func (a *Application) StoreWebhookEvent(ctx context.Context, event *domain.WebhookEvent) error {
	return a.webhookHandler.StoreWebhookEvent(ctx, event)
}

// UpdateWebhookEventStatus updates webhook event status
func (a *Application) UpdateWebhookEventStatus(ctx context.Context, id, status, errorMessage string) error {
	return a.webhookHandler.UpdateWebhookEventStatus(ctx, id, status, errorMessage)
}

// ProcessMaterialMaster processes MATMAS IDoc
func (a *Application) ProcessMaterialMaster(ctx context.Context, cmd ProcessMaterialMasterCommand) error {
	return a.webhookHandler.ProcessMaterialMaster(ctx, cmd.IDocData, cmd.CorrelationID, cmd.WebhookEventID)
}

// ProcessInventoryUpdate processes INVCON IDoc
func (a *Application) ProcessInventoryUpdate(ctx context.Context, cmd ProcessInventoryUpdateCommand) error {
	return a.webhookHandler.ProcessInventoryUpdate(ctx, cmd.IDocData, cmd.CorrelationID, cmd.WebhookEventID)
}

// ProcessPricingUpdate processes COND_A IDoc
func (a *Application) ProcessPricingUpdate(ctx context.Context, cmd ProcessPricingUpdateCommand) error {
	return a.webhookHandler.ProcessPricingUpdate(ctx, cmd.IDocData, cmd.CorrelationID, cmd.WebhookEventID)
}

// ProcessProductEvent processes JSON product events
func (a *Application) ProcessProductEvent(ctx context.Context, cmd ProcessProductEventCommand) error {
	return a.webhookHandler.ProcessProductEvent(ctx, cmd.Event, cmd.WebhookEventID)
}

// ProcessStockEvent processes JSON stock events
func (a *Application) ProcessStockEvent(ctx context.Context, cmd ProcessStockEventCommand) error {
	return a.webhookHandler.ProcessStockEvent(ctx, cmd.Event, cmd.WebhookEventID)
}

// ProcessPriceEvent processes JSON price events
func (a *Application) ProcessPriceEvent(ctx context.Context, cmd ProcessPriceEventCommand) error {
	return a.webhookHandler.ProcessPriceEvent(ctx, cmd.Event, cmd.WebhookEventID)
}

// ProcessProductDeletedEvent processes product deletion events
func (a *Application) ProcessProductDeletedEvent(ctx context.Context, cmd ProcessProductDeletedEventCommand) error {
	return a.webhookHandler.ProcessProductDeletedEvent(ctx, cmd.Event, cmd.WebhookEventID)
}

// SyncProductsFromSAP synchronizes products from SAP
func (a *Application) SyncProductsFromSAP(ctx context.Context, cmd SyncProductsCommand) error {
	return a.syncHandler.SyncProductsFromSAP(ctx, cmd.Since)
}

// SyncStockLevelsFromSAP synchronizes stock levels from SAP
func (a *Application) SyncStockLevelsFromSAP(ctx context.Context, cmd SyncStockLevelsCommand) error {
	return a.syncHandler.SyncStockLevelsFromSAP(ctx, cmd.ProductIDs)
}

// SyncPricesFromSAP synchronizes prices from SAP
func (a *Application) SyncPricesFromSAP(ctx context.Context, cmd SyncPricesCommand) error {
	return a.syncHandler.SyncPricesFromSAP(ctx, cmd.ProductIDs, cmd.PriceListID)
}

// GetSyncStatus gets sync status for an entity
func (a *Application) GetSyncStatus(ctx context.Context, entityType, entityID string) (*domain.SyncStatus, error) {
	return a.syncStatusRepo.GetByEntityID(ctx, entityType, entityID)
}

// GetRecentSyncLogs gets recent sync logs
func (a *Application) GetRecentSyncLogs(ctx context.Context, limit int) ([]*domain.SyncLog, error) {
	return a.syncLogRepo.GetRecentLogs(ctx, limit)
}

// GetSyncConfiguration gets sync configuration for an entity type
func (a *Application) GetSyncConfiguration(ctx context.Context, entityType string) (*domain.SyncConfiguration, error) {
	return a.syncConfigRepo.GetByEntityType(ctx, entityType)
}

// Update Application struct to include handlers
type Application struct {
	sapClient        *sap.EnhancedSAPClient
	eventPublisher   ddd.EventPublisher[ddd.Event]
	syncStatusRepo   domain.SyncStatusRepository
	syncLogRepo      domain.SyncLogRepository
	syncConfigRepo   domain.SyncConfigurationRepository
	webhookEventRepo domain.WebhookEventRepository
	transformer      *transformer.SAPToCanonicalTransformer
	
	// Command handlers
	webhookHandler *commands.ProcessWebhookHandler
	syncHandler    *commands.SyncFromSAPHandler
}

// NewApplication creates a new SAP connector application
func NewApplication(
	sapClient *sap.EnhancedSAPClient,
	eventPublisher ddd.EventPublisher[ddd.Event],
	syncStatusRepo domain.SyncStatusRepository,
	syncLogRepo domain.SyncLogRepository,
	syncConfigRepo domain.SyncConfigurationRepository,
	webhookEventRepo domain.WebhookEventRepository,
) *Application {
	// Create command handlers
	webhookHandler := commands.NewProcessWebhookHandler(
		webhookEventRepo,
		syncLogRepo,
		syncStatusRepo,
		eventPublisher,
	)
	
	syncHandler := commands.NewSyncFromSAPHandler(
		sapClient,
		syncStatusRepo,
		syncLogRepo,
		eventPublisher,
	)
	
	return &Application{
		sapClient:        sapClient,
		eventPublisher:   eventPublisher,
		syncStatusRepo:   syncStatusRepo,
		syncLogRepo:      syncLogRepo,
		syncConfigRepo:   syncConfigRepo,
		webhookEventRepo: webhookEventRepo,
		transformer:      transformer.NewSAPToCanonicalTransformer("SAP"),
		webhookHandler:   webhookHandler,
		syncHandler:      syncHandler,
	}
}