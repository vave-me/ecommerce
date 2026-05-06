package constants

// ServiceName The name of this module/service
const ServiceName = "activity"

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
	RedisPoolKey                = "redisPool"

	MiddlemanInteractionHandlersKey  = "middlemanInteractionHandlers"
	MiddlemanHandlersKey             = "middlemanHandlers"
	ActivityRepoKey                  = "activityRepo"
	InteractionsRepoKey              = "interactionsRepo"
	MiddlemanRepoKey                 = "middlemanRepo"
	MiddlemanCacheRepoKey            = "middlemanCacheRepo"
	MiddlemanInteractionRepoKey      = "middlemanInteractionRepo"
	MiddlemanCacheInteractionRepoKey = "middlemanCacheInteractionRepo"
	UsersRepoKey                     = "usersRepo"
	ProductsRepoKey                  = "productsRepo"
)

// Repository Table Names
const (
	OutboxTableName                     = ServiceName + ".outbox"
	InboxTableName                      = ServiceName + ".inbox"
	EventsTableName                     = ServiceName + ".events"
	SnapshotsTableName                  = ServiceName + ".snapshots"
	SagasTableName                      = ServiceName + ".sagas"
	UsersCacheTableName                 = ServiceName + ".users_cache"
	MiddlemanTableName                  = ServiceName + ".activity"
	MiddlemanInteractionTableName       = ServiceName + ".interactions"
	MiddlemanInteractionCountsTableName = ServiceName + ".item_interaction_counts"

	ProductsCacheTableName                   = ServiceName + ".products_cache"
	MiddlemanCacheTableName                  = ServiceName + ".activity_cache"
	MiddlemanCacheInteractionTableName       = ServiceName + ".interactions_cache"
	MiddlemanCacheInteractionCountsTableName = ServiceName + ".item_interaction_counts_cache"
)

// Metric Names
const (
	ActivitySentCount    = "user_activity_count"
	InteractionSentCount = "interactions_sent"
)
