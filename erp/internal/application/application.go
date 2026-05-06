package application

import (
	"context"
	"middleman/internal/ddd"
	"middleman/internal/erp"
	"middleman/internal/es"

	"middleman/erp/internal/application/commands"
	"middleman/erp/internal/application/queries"
	"middleman/erp/internal/crypto"
	"middleman/erp/internal/domain"
)

type (
	// App defines the application interface combining Commands and Queries
	App interface {
		Commands
		Queries
	}

	// Commands defines all available commands
	Commands interface {
		// Connector management
		RegisterConnector(ctx context.Context, cmd commands.RegisterConnector) error
		AddConnector(ctx context.Context, cmd commands.AddConnector) error
		UpdateConnector(ctx context.Context, cmd commands.UpdateConnector) error
		RemoveConnector(ctx context.Context, cmd commands.RemoveConnector) error
		ToggleConnector(ctx context.Context, cmd commands.ToggleConnector) error
		
		// Sync commands
		SyncProducts(ctx context.Context, cmd commands.SyncProducts) error
		SyncStock(ctx context.Context, cmd commands.SyncStock) error
		SyncPrices(ctx context.Context, cmd commands.SyncPrices) error
		SyncCustomers(ctx context.Context, cmd commands.SyncCustomers) error
		SendOrder(ctx context.Context, cmd commands.SendOrder) error
		UpdateInventoryReservation(ctx context.Context, cmd commands.UpdateInventoryReservation) error
		
		// Invoice commands
		CreateInvoice(ctx context.Context, cmd commands.CreateInvoice) error
		ApproveInvoice(ctx context.Context, cmd commands.ApproveInvoice) error
		SendInvoice(ctx context.Context, cmd commands.SendInvoice) error
		VoidInvoice(ctx context.Context, cmd commands.VoidInvoice) error
		RecordInvoicePayment(ctx context.Context, cmd commands.RecordInvoicePayment) error
		SyncInvoiceToERP(ctx context.Context, cmd commands.SyncInvoiceToERP) error
		
		// Return commands
		CreateReturn(ctx context.Context, cmd commands.CreateReturn) error
		ApproveReturn(ctx context.Context, cmd commands.ApproveReturn) error
		RejectReturn(ctx context.Context, cmd commands.RejectReturn) error
		ProcessReturnStart(ctx context.Context, cmd commands.ProcessReturnStart) error
		CompleteReturn(ctx context.Context, cmd commands.CompleteReturn) error
		RestockReturnItems(ctx context.Context, cmd commands.RestockReturnItems) error
		SyncReturnToERP(ctx context.Context, cmd commands.SyncReturnToERP) error
		
		// Inventory Reservation commands
		CreateInventoryReservation(ctx context.Context, cmd commands.CreateInventoryReservation) error
		ReleaseInventoryReservation(ctx context.Context, cmd commands.ReleaseInventoryReservation) error
		FulfillInventoryReservation(ctx context.Context, cmd commands.FulfillInventoryReservation) error
		TransferInventoryReservation(ctx context.Context, cmd commands.TransferInventoryReservation) error
		
		// Webhook commands
		ProcessWebhook(ctx context.Context, cmd commands.ProcessWebhook) error
	}

	// Queries defines all available queries
	Queries interface {
		GetConnectorStatus(ctx context.Context, query queries.GetConnectorStatus) (*queries.ConnectorStatus, error)
		ListConnectors(ctx context.Context, query queries.ListConnectors) ([]queries.ConnectorListItem, error)
		GetSyncHistory(ctx context.Context, query queries.GetSyncHistory) ([]queries.SyncHistoryItem, error)
	}

	// Application implements the ERP connector application using CQRS
	Application struct {
		appCommands
		appQueries
	}

	// appCommands embeds all command handlers
	appCommands struct {
		// Connector management
		commands.RegisterConnectorHandler
		commands.AddConnectorHandler
		commands.UpdateConnectorHandler
		commands.RemoveConnectorHandler
		commands.ToggleConnectorHandler

		// Sync handlers
		commands.SyncProductsHandler
		commands.SyncStockHandler
		commands.SyncPricesHandler
		commands.SyncCustomersHandler

		// Order handlers
		commands.SendOrderHandler
		commands.UpdateInventoryReservationHandler

		// Invoice handlers (granular)
		commands.CreateInvoiceHandler
		commands.ApproveInvoiceHandler
		commands.SendInvoiceHandler
		commands.VoidInvoiceHandler
		commands.RecordInvoicePaymentHandler
		commands.SyncInvoiceToERPHandler

		// Return handlers (granular)
		commands.CreateReturnHandler
		commands.ApproveReturnHandler
		commands.RejectReturnHandler
		commands.ProcessReturnStartHandler
		commands.CompleteReturnHandler
		commands.RestockReturnItemsHandler
		commands.SyncReturnToERPHandler

		// Inventory Reservation handlers
		commands.CreateInventoryReservationHandler
		commands.ReleaseInventoryReservationHandler
		commands.FulfillInventoryReservationHandler
		commands.TransferInventoryReservationHandler
		
		// Webhook handlers
		commands.ProcessWebhookHandler
	}

	// appQueries embeds all query handlers
	appQueries struct {
		queries.GetConnectorStatusHandler
		queries.ListConnectorsHandler
		queries.GetSyncHistoryHandler
	}
)

var _ App = (*Application)(nil)

// eventPublisherAdapter adapts ddd.EventPublisher[ddd.Event] to domain.EventPublisher
type eventPublisherAdapter struct {
	publisher ddd.EventPublisher[ddd.Event]
}

func (a eventPublisherAdapter) Publish(ctx context.Context, event domain.Event) error {
	// Convert domain.Event to ddd.Event
	if dddEvent, ok := event.(ddd.Event); ok {
		return a.publisher.Publish(ctx, dddEvent)
	}
	// If it's not a ddd.Event, we can't publish it
	return nil
}

// New creates a new ERP application with event sourcing support
func New(
	invoices es.AggregateRepository[*domain.Invoice],
	returns es.AggregateRepository[*domain.Return],
	reservations es.AggregateRepository[*domain.InventoryReservation],
	publisher ddd.EventPublisher[ddd.Event],
	factory erp.ConnectorFactory,
	registry erp.ConnectorRegistry,
	connectorRepo domain.ConnectorRepository,
	syncLogRepo domain.SyncLogRepository,
	webhookEventRepo domain.WebhookEventRepository,
	orderSyncRepo domain.OrderSyncRepository,
	invoiceSyncRepo domain.InvoiceSyncRepository,
	productRepo domain.ProductRepository,
	encryptor *crypto.Encryptor,
) *Application {
	return &Application{
		appCommands: appCommands{
			// Connector management
			RegisterConnectorHandler:            commands.NewRegisterConnectorHandler(factory, registry),
			AddConnectorHandler:                 commands.NewAddConnectorHandler(connectorRepo, factory, registry, encryptor),
			UpdateConnectorHandler:              commands.NewUpdateConnectorHandler(connectorRepo, factory, registry, encryptor),
			RemoveConnectorHandler:              commands.NewRemoveConnectorHandler(connectorRepo, registry, syncLogRepo, orderSyncRepo, invoiceSyncRepo),
			ToggleConnectorHandler:              commands.NewToggleConnectorHandler(connectorRepo, registry, factory, encryptor),
			SyncProductsHandler:                 commands.NewSyncProductsHandler(registry, syncLogRepo, productRepo),
			SyncStockHandler:                    commands.NewSyncStockHandler(registry, syncLogRepo, productRepo),
			SyncPricesHandler:                   commands.NewSyncPricesHandler(registry, syncLogRepo, productRepo),
			SyncCustomersHandler:                commands.NewSyncCustomersHandler(registry, syncLogRepo),
			SendOrderHandler:                    commands.NewSendOrderHandler(registry, orderSyncRepo),
			UpdateInventoryReservationHandler:   commands.NewUpdateInventoryReservationHandler(registry, orderSyncRepo),
			CreateInvoiceHandler:                commands.NewCreateInvoiceHandler(invoices, publisher),
			ApproveInvoiceHandler:               commands.NewApproveInvoiceHandler(invoices, publisher),
			SendInvoiceHandler:                  commands.NewSendInvoiceHandler(invoices, publisher),
			VoidInvoiceHandler:                  commands.NewVoidInvoiceHandler(invoices, publisher),
			RecordInvoicePaymentHandler:         commands.NewRecordInvoicePaymentHandler(invoices, publisher),
			SyncInvoiceToERPHandler:             commands.NewSyncInvoiceToERPHandler(invoices, registry, invoiceSyncRepo, publisher),
			CreateReturnHandler:                 commands.NewCreateReturnHandler(returns, publisher),
			ApproveReturnHandler:                commands.NewApproveReturnHandler(returns, publisher),
			RejectReturnHandler:                 commands.NewRejectReturnHandler(returns, publisher),
			ProcessReturnStartHandler:           commands.NewProcessReturnStartHandler(returns, publisher),
			CompleteReturnHandler:               commands.NewCompleteReturnHandler(returns, publisher),
			RestockReturnItemsHandler:           commands.NewRestockReturnItemsHandler(returns, publisher),
			SyncReturnToERPHandler:              commands.NewSyncReturnToERPHandler(returns, registry, publisher),
			CreateInventoryReservationHandler:   commands.NewCreateInventoryReservationHandler(reservations, publisher),
			ReleaseInventoryReservationHandler:  commands.NewReleaseInventoryReservationHandler(reservations, publisher),
			FulfillInventoryReservationHandler:  commands.NewFulfillInventoryReservationHandler(reservations, publisher),
			TransferInventoryReservationHandler: commands.NewTransferInventoryReservationHandler(reservations, publisher),
			ProcessWebhookHandler:               commands.NewProcessWebhookHandler(registry, connectorRepo, webhookEventRepo, eventPublisherAdapter{publisher}),
		},
		appQueries: appQueries{
			GetConnectorStatusHandler: queries.NewGetConnectorStatusHandler(registry, syncLogRepo, webhookEventRepo),
			ListConnectorsHandler:     queries.NewListConnectorsHandler(registry),
			GetSyncHistoryHandler:     queries.NewGetSyncHistoryHandler(syncLogRepo),
		},
	}
}
