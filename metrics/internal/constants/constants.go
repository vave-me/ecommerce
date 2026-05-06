package constants

// ServiceName The name of this module/service
const ServiceName = "metrics"

// GRPC Service Names
const (
	UsersServiceName      = "USERS"
	ProductsServiceName   = "PRODUCTS"
	PostsServiceName      = "POSTS"
	VehiclesServiceName   = "VEHICLES"
	PropertiesServiceName = "PROPERTIES"
	ServicesServiceName   = "SERVICES"
	DealsServiceName      = "DEALS"
	JobsServiceName       = "JOBS"
)

// Dependency Injection Keys
const (
	RegistryKey                 = "registry"
	DomainDispatcherKey         = "domainDispatcher"
	DatabaseTransactionKey      = "tx"
	RedisTransactionKey         = "redis"
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
	// New constants for Redis and Redisearch
	RedisClientKey           = "redisClient"
	RedisearchClientKey      = "redisearchClient"
	RedisearchIndex          = "metricsIndex"
	OrdersRepoKey            = "ordersRepo"
	VehiclesRepoKey          = "vehiclesRepo"
	PropertiesRepoKey        = "propertiesRepo"
	UsersMetricsCacheRepoKey = "userMetricsCacheRepo"
	ItemMetricsCacheRepoKey  = "itemMetricsCacheRepo"
	UserMetricsRepoKey       = "itemMetricsRepo"
	ItemMetricsRepoKey       = "usersMetricsRepo"
)

// Repository Table Names
const (
	OutboxTableName    = ServiceName + ".outbox"
	InboxTableName     = ServiceName + ".inbox"
	EventsTableName    = ServiceName + ".events"
	SnapshotsTableName = ServiceName + ".snapshots"
	SagasTableName     = ServiceName + ".sagas"

	OrdersTableName = ServiceName + ".orders"

	UsersMetricsCacheTableName = ServiceName + ".users_metrics_cache"
	ItemMetricsCacheTableName  = ServiceName + ".items_metrics_cache"
)
