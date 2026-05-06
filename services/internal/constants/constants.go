package constants

// ServiceName The name of this module/service
const ServiceName = "services"

// GRPC Service Names
const (
	ServicesServiceName = "USERS"
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
	CatalogVariantHandlersKey   = "catalogVariantHandlers"
	MiddlemanHandlersKey        = "middlemanHandlers"

	ServicesRepoKey = "servicesRepo"
	CatalogRepoKey  = "catalogRepo"

	VariantsRepoKey       = "variantsRepo"
	CatalogVariantRepoKey = "catalogVariantRepo"
)

// Repository Table Names
const (
	OutboxTableName    = ServiceName + ".outbox"
	InboxTableName     = ServiceName + ".inbox"
	EventsTableName    = ServiceName + ".events"
	SnapshotsTableName = ServiceName + ".snapshots"
	SagasTableName     = ServiceName + ".sagas"

	CatalogTableName        = ServiceName + ".services"
	CatalogVariantTableName = ServiceName + ".variants"
)
