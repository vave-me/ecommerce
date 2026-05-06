package constants

// ServiceName The name of this module/service
const ServiceName = "products"

// GRPC Service Names
const (
	ProductsServiceName = "USERS"
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

	ProductsRepoKey = "productsRepo"
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

	CatalogTableName        = ServiceName + ".products"
	CatalogVariantTableName = ServiceName + ".variants"
)
