package constants

// ServiceName The name of this module/service
const ServiceName = "comments"

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

	CommentRepliesHandlersKey = "commentRepliesHandlers"
	MiddlemanHandlersKey      = "middlemanHandlers"
	CommentsRepoKey           = "commentsRepo"
	CommentsCacheRepoKey      = "commentsRepo"
	RepliesRepoKey            = "repliesRepo"
	CommentRepliesRepoKey     = "commentRepliesRepo"
	MiddlemanRepoKey          = "middlemanRepo"
	MiddlemanCacheRepoKey     = "middlemanCacheRepo"
)

// Repository Table Names
const (
	OutboxTableName    = ServiceName + ".outbox"
	InboxTableName     = ServiceName + ".inbox"
	EventsTableName    = ServiceName + ".events"
	SnapshotsTableName = ServiceName + ".snapshots"
	SagasTableName     = ServiceName + ".sagas"

	CommentRepliesTableName = ServiceName + ".replies"
	MiddlemanTableName      = ServiceName + ".comments"
	MiddlemanCacheTableName = ServiceName + ".comments_cache"
)
