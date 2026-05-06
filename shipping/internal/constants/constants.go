package constants

// ServiceName The name of this module/service
const ServiceName = "shipping"

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
	AggregateStoreKey           = "aggregateStore"
	SagaStoreKey                = "sagaStore"
	InboxStoreKey               = "inboxStore"
	ApplicationKey              = "app"
	DomainEventHandlersKey      = "domainEventHandlers"
	IntegrationEventHandlersKey = "integrationEventHandlers"
	CommandHandlersKey          = "commandHandlers"
	ReplyHandlersKey            = "replyHandlers"
	CatalogHandlersKey          = "catalogHandlers"
	DHLClient                   = "dhlClient"
	ShippingRepoKey             = "shippingRepo"
	ShippingCatalogRepoKey      = "shippingCatalogRepo"
	ImagesRepoKey               = "imagesRepo"
	VideosRepoKey               = "videosRepo"
	MiddlemanRepoKey            = "middlemanRepo"
)

// Repository Table Names
const (
	OutboxTableName    = ServiceName + ".outbox"
	InboxTableName     = ServiceName + ".inbox"
	EventsTableName    = ServiceName + ".events"
	SnapshotsTableName = ServiceName + ".snapshots"
	SagasTableName     = ServiceName + ".sagas"
	
	CatalogTableName   = ServiceName + ".shipping"
)
