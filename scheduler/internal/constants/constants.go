package constants

// ServiceName The name of this module/service
const ServiceName = "scheduler"

// Default Assistant ID for scheduler tasks
const DefaultAssistantID = "scheduler-assistant"

// GRPC Service Names
const (
	AssistantsServiceName = "ASSISTANTS"
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
	CommandHandlersKey          = "commandHandlers"
	ReplyHandlersKey            = "replyHandlers"
	RedisPoolKey                = "redisPool"

	MiddlemanActionHandlersKey  = "middlemanActionHandlers"
	MiddlemanHandlersKey        = "middlemanHandlers"
	SchedulerRepoKey            = "schedulerRepo"
	ActionsRepoKey              = "interactionsRepo"
	MiddlemanRepoKey            = "middlemanRepo"
	MiddlemanCacheRepoKey       = "middlemanCacheRepo"
	MiddlemanActionRepoKey      = "middlemanActionRepo"
	MiddlemanCacheActionRepoKey = "middlemanCacheActionRepo"
	AssistantRepoKey            = "assistantRepo"
	SchedulerWorkerKey          = "schedulerWorker"
	TaskRepoKey                 = "taskRepo"
	CatalogTaskRepoKey          = "catalogTaskRepo"
	TaskHandlersKey             = "taskHandlers"
)

// Repository Table Names
const (
	OutboxTableName                = ServiceName + ".outbox"
	InboxTableName                 = ServiceName + ".inbox"
	EventsTableName                = ServiceName + ".events"
	SnapshotsTableName             = ServiceName + ".snapshots"
	SagasTableName                 = ServiceName + ".sagas"
	MiddlemanTableName             = ServiceName + ".scheduler"
	MiddlemanActionTableName       = ServiceName + ".actions"
	MiddlemanCacheTableName        = ServiceName + ".scheduler_cache"
	MiddlemanCacheActionTableName  = ServiceName + ".actions_cache"
	CatalogTaskTableName           = ServiceName + ".catalog_tasks"
)

// Metric Names
const (
	SchedulerSentCount = "user_scheduler_count"
	ActionSentCount    = "interactions_sent"
)
