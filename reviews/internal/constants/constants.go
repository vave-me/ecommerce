package constants

// ServiceName The name of this module/service
const ServiceName = "reviews"

// GRPC Service Names
const (
	UsersServiceName = "USERS"
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
	WebSocketPublisherKey       = "webSocketPublisher"
	Logger                      = "logger"
	ReplyPublisherKey           = "replyPublisher"
	AggregateStoreKey           = "aggregateStore"
	SagaStoreKey                = "sagaStore"
	InboxStoreKey               = "inboxStore"
	ApplicationKey              = "app"
	DomainEventHandlersKey      = "domainEventHandlers"
	WebsocketsEventHandlersKey  = "websocketCommandHandlers"
	WebSocketSubscriberKey      = "websocketSubscribers"
	IntegrationEventHandlersKey = "integrationEventHandlers"
	CommandHandlersKey          = "commandHandlers"
	ReplyHandlersKey            = "replyHandlers"
	RedisClientKey              = "redisClient"
	RedisPoolKey                = "redisPool"

	ReviewRepliesHandlersKey = "reviewRepliesHandlers"
	MiddlemanHandlersKey     = "middlemanHandlers"
	ReviewsRepoKey           = "reviewsRepo"
	ReviewsCacheRepoKey      = "reviewsRepo"
	RepliesRepoKey           = "repliesRepo"
	ReviewRepliesRepoKey     = "reviewRepliesRepo"
	MiddlemanRepoKey         = "middlemanRepo"
	MiddlemanCacheRepoKey    = "middlemanCacheRepo"
)

// Repository Table Names
const (
	OutboxTableName    = ServiceName + ".outbox"
	InboxTableName     = ServiceName + ".inbox"
	EventsTableName    = ServiceName + ".events"
	SnapshotsTableName = ServiceName + ".snapshots"
	SagasTableName     = ServiceName + ".sagas"

	ReviewRepliesTableName  = ServiceName + ".replies"
	MiddlemanTableName      = ServiceName + ".reviews"
	MiddlemanCacheTableName = ServiceName + ".reviews_cache"
)
