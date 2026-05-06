package constants

// ServiceName The name of this module/service
const ServiceName = "following"

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

	FollowRepliesHandlersKey = "followRepliesHandlers"
	MiddlemanHandlersKey     = "middlemanHandlers"
	FollowingRepoKey         = "followingRepo"
	FollowingCacheRepoKey    = "followingRepo"
	RepliesRepoKey           = "repliesRepo"
	FollowRepliesRepoKey     = "followRepliesRepo"
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

	FollowRepliesTableName  = ServiceName + ".replies"
	MiddlemanTableName      = ServiceName + ".following"
	MiddlemanCacheTableName = ServiceName + ".following_cache"
)
