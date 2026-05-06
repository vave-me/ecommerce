package constants

// ServiceName The name of this module/service
const ServiceName = "search"

// GRPC Service Names
const (
	UsersServiceName      = "USERS"
	ProductsServiceName   = "PRODUCTS"
	PostsServiceName      = "POSTS"
	VehiclesServiceName   = "VEHICLES"
	PropertiesServiceName = "PROPERTIES"
	ServicesServiceName   = "SERVICES"
	DealsServiceName      = "DEALS"
	JobsServiceName       = "JOBS"
	MetricsServiceName    = "METRICS"
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
	ReplyPublisherKey           = "replyPublisher"
	SagaStoreKey                = "sagaStore"
	InboxStoreKey               = "inboxStore"
	ApplicationKey              = "app"
	DomainEventHandlersKey      = "domainEventHandlers"
	IntegrationEventHandlersKey = "integrationEventHandlers"
	CommandHandlersKey          = "commandHandlers"
	ReplyHandlersKey            = "replyHandlers"
	// New constants for Redis and Redisearch
	RedisClientKey      = "redisClient"
	RedisearchClientKey = "redisearchClient"
	RedisearchIndex     = "productsIndex"
	OrdersRepoKey       = "ordersRepo"
	VehiclesRepoKey     = "vehiclesRepo"
	PropertiesRepoKey   = "propertiesRepo"
	UsersRepoKey        = "usersRepo"
	ProductsRepoKey     = "productsRepo"
	ServicesRepoKey     = "servicesRepo"
	PostsRepoKey        = "postsRepo"
	DealsRepoKey        = "dealsRepo"
	JobsRepoKey         = "jobsRepo"
	VariantsRepoKey     = "variantsRepo"
	MetricsRepoKey      = "metricsRepo"
)

// Repository Table Names
const (
	OutboxTableName    = ServiceName + ".outbox"
	InboxTableName     = ServiceName + ".inbox"
	EventsTableName    = ServiceName + ".events"
	SnapshotsTableName = ServiceName + ".snapshots"
	SagasTableName     = ServiceName + ".sagas"

	OrdersTableName = ServiceName + ".orders"

	UsersCacheTableName      = ServiceName + ".users_cache"
	ProductsCacheTableName   = ServiceName + ".products_cache"
	PostsCacheTableName      = ServiceName + ".posts_cache"
	DealsCacheTableName      = ServiceName + ".deals_cache"
	JobsCacheTableName       = ServiceName + ".jobs_cache"
	VehiclesCacheTableName   = ServiceName + ".vehicles_cache"
	PropertiesCacheTableName = ServiceName + ".properties_cache"
	ServicesCacheTableName   = ServiceName + ".services_cache"
	VariantsCacheTableName   = ServiceName + ".variants_cache"
)

// Metric-based Sort Types for Advanced Workflow
const (
	SortByLikesHigh      = "likes_high"
	SortByLikesLow       = "likes_low"
	SortByCommentsHigh   = "comments_high"
	SortByCommentsLow    = "comments_low"
	SortByVisitedHigh    = "visited_high"
	SortByVisitedLow     = "visited_low"
	SortByRatingHigh     = "rating_high"
	SortByRatingLow      = "rating_low"
	SortByPopularityHigh = "popularity_high" // Combined metric
	SortByTrendingHigh   = "trending_high"   // Time-weighted engagement
)

// Metric Type Mappings for API
const (
	MetricTypeLikes    = "likesCount"
	MetricTypeComments = "commentsCount"
	MetricTypeVisited  = "visitedCount"
	MetricTypeRating   = "rating"
	MetricTypeShares   = "sharedCount"
	MetricTypeWishlist = "addedToWishlistCount"
)
