package constants

// ServiceName The name of this module/service
const ServiceName = "notifications"

// GRPC Service Names
const (
	UsersServiceName    = "USERS"
	ProductsServiceName = "PRODUCTS"
)

// Dependency Injection Keys
const (
	RegistryKey                 = "registry"
	DomainDispatcherKey         = "domainDispatcher"
	DatabaseTransactionKey      = "tx"
	MessagePublisherKey         = "messagePublisher"
	MessageSubscriberKey        = "messageSubscriber"
	EventPublisherKey           = "eventPublisher"
	CommandPublisherKey         = "commandPublisher"
	ReplyPublisherKey           = "replyPublisher"
	SagaStoreKey                = "sagaStore"
	InboxStoreKey               = "inboxStore"
	AggregateStoreKey           = "aggregateStore"
	ApplicationKey              = "app"
	DomainEventHandlersKey      = "domainEventHandlers"
	IntegrationEventHandlersKey = "integrationEventHandlers"
	CommandHandlersKey          = "commandHandlers"
	ReplyHandlersKey            = "replyHandlers"

	MiddlemanHandlersKey       = "middlemanHandlers"
	MiddlemanAlertsHandlersKey = "middlemanAlertsHandlers"
	NotificationsRepoKey       = "notificationsRepo"
	AlertsRepoKey              = "alertsRepo"
	CatalogRepoKey             = "catalogRepo"
	UsersRepoKey               = "usersRepo"
	ProductsRepoKey            = "productsRepo"
	PreferencesRepoKey         = "preferencesRepo"
)

// Repository Table Names
const (
	OutboxTableName        = ServiceName + ".outbox"
	InboxTableName         = ServiceName + ".inbox"
	EventsTableName        = ServiceName + ".events"
	SnapshotsTableName     = ServiceName + ".snapshots"
	SagasTableName         = ServiceName + ".sagas"
	UsersCacheTableName    = ServiceName + ".users_cache"
	CatalogTableName       = ServiceName + ".alerts"
	ProductsCacheTableName = ServiceName + ".products_cache"
)

// Metric Names
const (
	NotificationsSentCount       = "notifications_sent_count"
	InteractionNotificationsSent = "interaction_notifications_sent"
)
