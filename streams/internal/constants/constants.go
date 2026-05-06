package constants

// ServiceName The name of this module/service
const ServiceName = "streams"

// GRPC Service Names
const (
	StreamsServiceName = "STREAMS"
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
	SeriesHandlersKey           = "seriesHandlers"

	StreamsRepoKey = "streamsRepo"
	SeriesRepoKey  = "seriesRepo"
	CatalogRepoKey = "catalogRepo"
	
	// Live streaming keys
	LiveStreamsRepoKey = "liveStreamsRepo"
	WebhookSubscriptionRepoKey = "webhookSubscriptionRepo"
	WebhookDeliveryRepoKey = "webhookDeliveryRepo"
	WebhookClientKey = "webhookClient"
	WebhookDispatcherKey = "webhookDispatcher"
	StreamingServerKey = "streamingServer"
	CDNManagerKey = "cdnManager"
	DRMManagerKey = "drmManager"
	WebRTCServerKey = "webrtcServer"
)

// Repository Table Names
const (
	OutboxTableName    = ServiceName + ".outbox"
	InboxTableName     = ServiceName + ".inbox"
	EventsTableName    = ServiceName + ".events"
	SnapshotsTableName = ServiceName + ".snapshots"
	SagasTableName     = ServiceName + ".sagas"

	StreamsTableName = ServiceName + ".streams"
	SeriesTableName  = ServiceName + ".series"
	CatalogTableName = ServiceName + ".catalog"
)

// Event store streams
const (
	StreamEventStream = "stream-events"
	SeriesEventStream = "series-events"
)

// Default values
const (
	DefaultPageSize        = 20
	MaxPageSize            = 100
	DefaultRentalDuration  = 48  // hours
	DefaultPPVDuration     = 48  // hours
	MaxSearchResults       = 1000
	DefaultStreamQuality   = "1080p"
)

// Cache TTL values
const (
	CatalogCacheTTL    = 300  // 5 minutes
	StreamCacheTTL     = 600  // 10 minutes
	UserAccessCacheTTL = 60   // 1 minute
)
