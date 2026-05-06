package constants

// ServiceName The name of this module/service
const ServiceName = "geocoding"

// GRPC Service Names
const (
	UsersServiceName = "USERS"
)

// Dependency Injection Keys
const (
	RegistryKey                 = "registry"
	DomainDispatcherKey         = "domainDispatcher"
	DatabaseTransactionKey      = "tx"
	GoogleGeocode               = "googleGeocode"
	NominatimGeocode            = "nominatimGeocode"
	RedisTransactionKey         = "redis"
	MessagePublisherKey         = "messagePublisher"
	MessageSubscriberKey        = "messageSubscriber"
	EventPublisherKey           = "eventPublisher"
	AggregateStoreKey           = "aggregateStore"
	CommandPublisherKey         = "commandPublisher"
	ReplyPublisherKey           = "replyPublisher"
	SagaStoreKey                = "sagaStore"
	InboxStoreKey               = "inboxStore"
	ApplicationKey              = "app"
	DomainEventHandlersKey      = "domainEventHandlers"
	IntegrationEventHandlersKey = "integrationEventHandlers"
	CommandHandlersKey          = "commandHandlers"
	ReplyHandlersKey            = "replyHandlers"
	// New constants for Redis and Redigeocoding
	RedisClientKey             = "redisClient"
	RedigeocodingClientKey     = "redigeocodingClient"
	RedigeocodingIndex         = "productsIndex"
	AddressesRepoKey           = "addressesRepo"
	LocationsRepoKey           = "locationsRepo"
	ProductsRepoKey            = "productsRepo"
	VariantsRepoKey            = "variantsRepo"
	CatalogRepoKey             = "catalogRepo"
	CatalogLocationRepoKey     = "catalogLocationRepo"
	CatalogHandlersKey         = "catalogHandlers"
	CatalogLocationHandlersKey = "catalogLocationHandlers"
)

// Repository Table Names
const (
	OutboxTableName          = ServiceName + ".outbox"
	InboxTableName           = ServiceName + ".inbox"
	EventsTableName          = ServiceName + ".events"
	SnapshotsTableName       = ServiceName + ".snapshots"
	SagasTableName           = ServiceName + ".sagas"
	CatalogTableName         = ServiceName + ".addresses"
	CatalogLocationTableName = ServiceName + ".locations"
	LocationsTableName       = ServiceName + ".locations"
	AddressesTableName       = ServiceName + ".addresses"
)
