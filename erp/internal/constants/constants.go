package constants

// ServiceName The name of this module/service
const ServiceName = "erp"

// GRPC Service Names
const (
	ERPServiceName      = "ERP"
	ProductsServiceName = "products"
)

// Stream Names
const (
	CommandStream = "commands"
	EventStream   = "events"
)

// Message Prefixes
const (
	CommandPrefix      = "CMD"
	IntegrationPrefix  = "INT"
	DomainEventPrefix  = "DOM"
)

// Dependency Injection Keys
const (
	RegistryKey                 = "registry"
	DomainDispatcherKey         = "domainDispatcher"
	DatabaseTransactionKey      = "tx"
	MessagePublisherKey         = "messagePublisher"
	MessageSubscriberKey        = "messageSubscriber"
	EventPublisherKey           = "eventPublisher"
	DomainEventPublisherKey     = "domainEventPublisher"
	CommandPublisherKey         = "commandPublisher"
	ReplyPublisherKey           = "replyPublisher"
	AggregateStoreKey           = "aggregateStore"
	SagaStoreKey                = "sagaStore"
	InboxStoreKey               = "inboxStore"
	ApplicationKey              = "app"
	ERPApplicationKey           = "erpApplication"
	DomainEventHandlersKey      = "domainEventHandlers"
	IntegrationEventHandlersKey = "integrationEventHandlers"
	CommandHandlersKey          = "commandHandlers"
	ReplyHandlersKey            = "replyHandlers"
	
	// ERP specific keys
	ConnectorRegistryKey    = "connectorRegistry"
	WebhookHandlerKey       = "webhookHandler"
	SyncHandlerKey          = "syncHandler"
	
	// Connector keys
	SAPConnectorKey         = "sapConnector"
	OdooConnectorKey        = "odooConnector"
	Dynamics365ConnectorKey = "dynamics365Connector"
	NetSuiteConnectorKey    = "netsuiteConnector"
	
	// Repository keys
	WebhookEventRepositoryKey      = "webhookEventRepository"
	SyncStatusRepositoryKey        = "syncStatusRepository"
	SyncLogRepositoryKey           = "syncLogRepository"
	SyncConfigurationRepositoryKey = "syncConfigurationRepository"
	OrderSyncRepositoryKey         = "orderSyncRepository"
	InvoiceSyncRepositoryKey       = "invoiceSyncRepository"
	ConnectorRepositoryKey         = "connectorRepository"
	ProductRepositoryKey           = "productRepository"
	
	// Factory and registry keys
	ConnectorFactoryKey               = "connectorFactory"
	DomainEventPublisherAdapterKey    = "domainEventPublisherAdapter"
	InvoiceRepositoryKey              = "invoiceRepository"
	ReturnRepositoryKey               = "returnRepository"
	InventoryReservationRepositoryKey = "inventoryReservationRepository"
	EncryptorKey                      = "encryptor"
)

// Repository Table Names
const (
	OutboxTableName    = "outbox"
	InboxTableName     = "inbox"
	EventsTableName    = "events"
	SnapshotsTableName = "snapshots"
	SagasTableName     = "sagas"
	
	// ERP specific tables
	WebhookEventsTableName      = "webhook_events"
	SyncStatusTableName         = "sync_status"
	SyncLogsTableName           = "sync_logs"
	SyncConfigurationsTableName = "sync_configurations"
	OrderSyncTableName          = "order_sync"
	InvoiceSyncTableName        = "invoice_sync"
)
