package constants

// ServiceName The name of this module/service
const ServiceName = "wishlists"

// GRPC Service Names
const (
	UsersServiceName    = "USERS"
	ProductsServiceName = "PRODUCTS"
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

	CatalogHandlersKey   = "catalogHandlers"
	MiddlemanHandlersKey = "middlemanHandlers"

	WishlistsRepoKey     = "wishlistsRepo"
	WishlistItemsRepoKey = "wishlistItemsRepo"
	CatalogRepoKey       = "catalogRepo"
	MiddlemanRepoKey     = "middlemanRepo"
	ProductsRepoKey      = "productsRepo"
)

// Repository Table Names
const (
	OutboxTableName        = ServiceName + ".outbox"
	InboxTableName         = ServiceName + ".inbox"
	EventsTableName        = ServiceName + ".events"
	SnapshotsTableName     = ServiceName + ".snapshots"
	SagasTableName         = ServiceName + ".sagas"
	ProductsCacheTableName = ServiceName + ".products_cache"

	CatalogTableName   = ServiceName + ".wishlist_items"
	MiddlemanTableName = ServiceName + ".wishlists"
)
