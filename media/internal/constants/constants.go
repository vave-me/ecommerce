package constants

// ServiceName The name of this module/service
const ServiceName = "media"

// GRPC Service Names
const (
	UsersServiceName     = "USERS"
	ProductsServiceName  = "PRODUCTS"
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
	MiddlemanMediaHandlersKey   = "middlemanMediaHandlers"
	MiddlemanImageHandlersKey   = "middlemanImageHandlers"
	MiddlemanVideoHandlersKey   = "middlemanVideoHandlers"
	IntegrationEventHandlersKey = "integrationEventHandlers"
	CommandHandlersKey          = "commandHandlers"
	ReplyHandlersKey            = "replyHandlers"
	AwsClient                   = "awsClient"
	MinioClient                 = "minioClient"
	MediaRepoKey                = "mediaRepo"
	ImagesRepoKey               = "imagesRepo"
	VideosRepoKey               = "videosRepo"
	MiddlemanRepoKey            = "middlemanRepo"
	MiddlemanMediaRepoKey       = "middlemanMediaRepo"
	MiddlemanImageRepoKey       = "middlemanImageRepo"
	MiddlemanVideoRepoKey       = "middlemanVideoRepo"
	ImportSessionRepoKey        = "importSessionRepo"
	ImportItemRepoKey           = "importItemRepo"
	ImporterRepoKey             = "importerRepo"
	ProductRepoKey              = "productRepo"
)

// Repository Table Names
const (
	OutboxTableName    = ServiceName + ".outbox"
	InboxTableName     = ServiceName + ".inbox"
	EventsTableName    = ServiceName + ".events"
	SnapshotsTableName = ServiceName + ".snapshots"
	SagasTableName     = ServiceName + ".sagas"

	MiddlemanMediaTableName = ServiceName + ".media"
	MiddlemanImageTableName = ServiceName + ".images"
	MiddlemanVideoTableName = ServiceName + ".videos"
	ImportSessionTableName  = ServiceName + ".import_sessions"
	ImportItemTableName     = ServiceName + ".import_items"
)
