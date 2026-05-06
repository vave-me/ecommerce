package constants

// ServiceName The name of this module/service
const ServiceName = "baskets"

// GRPC Service Names
const (
	UsersServiceName = "USERS"
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
	ApplicationKey              = "app"
	DomainEventHandlersKey      = "domainEventHandlers"
	IntegrationEventHandlersKey = "integrationEventHandlers"
	CommandHandlersKey          = "commandHandlers"
	ReplyHandlersKey            = "replyHandlers"
	CatalogRepoKey              = "catalogRepo"
	CatalogHandlersKey          = "catalogHandlers"
	BasketsRepoKey              = "basketsRepo"
	UsersRepoKey                = "usersRepo"
	ProductsRepoKey             = "productsRepo"
)

// Repository Table Names
const (
	OutboxTableName    = ServiceName + ".outbox"
	InboxTableName     = ServiceName + ".inbox"
	EventsTableName    = ServiceName + ".events"
	SnapshotsTableName = ServiceName + ".snapshots"
	SagasTableName     = ServiceName + ".sagas"

	CatalogTableName       = ServiceName + ".baskets"
	UsersCacheTableName    = ServiceName + ".users_cache"
	ProductsCacheTableName = ServiceName + ".products_cache"
)

// Metric Names
const (
	BasketsStartedCount    = "baskets_started_count"
	BasketsCheckedOutCount = "baskets_checked_out_count"
	BaksetsCanceledCount   = "baskets_canceled_count"
)
