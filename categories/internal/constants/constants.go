package constants

// ServiceName The name of this module/service
const ServiceName = "categories"

// GRPC Service Names
const (
	CategoriesServiceName = "USERS"
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
	AggregateStoreKey           = "aggregateStore"
	SagaStoreKey                = "sagaStore"
	InboxStoreKey               = "inboxStore"
	ApplicationKey              = "app"
	DomainEventHandlersKey      = "domainEventHandlers"
	IntegrationEventHandlersKey = "integrationEventHandlers"
	CommandHandlersKey          = "commandHandlers"
	ReplyHandlersKey            = "replyHandlers"
	RedisPoolKey                = "redisPool"
	CatalogHandlersKey          = "catalogHandlers"
	CatalogFilterHandlersKey    = "catalogFilterHandlers"
	MiddlemanHandlersKey        = "middlemanHandlers"

	CategoriesRepoKey = "categoriesRepo"
	CatalogRepoKey    = "catalogRepo"

	FiltersRepoKey       = "filtersRepo"
	CatalogFilterRepoKey = "catalogFilterRepo"
)

// Repository Table Names
const (
	OutboxTableName    = ServiceName + ".outbox"
	InboxTableName     = ServiceName + ".inbox"
	EventsTableName    = ServiceName + ".events"
	SnapshotsTableName = ServiceName + ".snapshots"
	SagasTableName     = ServiceName + ".sagas"

	CatalogTableName       = ServiceName + ".categories"
	CatalogFilterTableName = ServiceName + ".filters"
)
