package constants

// ServiceName The name of this module/service
const ServiceName = "support"

// GRPC Service Names
const (
	SupportServiceName = "SUPPORT"
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

	// Repository keys
	SupportChannelRepoKey          = "supportChannelRepo"
	TicketRepoKey                  = "ticketRepo"
	SupportChannelCatalogRepoKey   = "supportChannelCatalogRepo"
	TicketCatalogRepoKey           = "ticketCatalogRepo"
	KnowledgeArticleCatalogRepoKey = "knowledgeArticleCatalogRepo"
	CommunicationRepoKey           = "communicationRepo"
	AIConfigurationRepoKey         = "aiConfigurationRepo"
)

// Repository Table Names
const (
	OutboxTableName    = ServiceName + ".outbox"
	InboxTableName     = ServiceName + ".inbox"
	EventsTableName    = ServiceName + ".events"
	SnapshotsTableName = ServiceName + ".snapshots"
	SagasTableName     = ServiceName + ".sagas"

	// Domain tables
	SupportChannelsTableName          = ServiceName + ".support_channels"
	TicketsTableName                  = ServiceName + ".tickets"
	CommunicationsTableName           = ServiceName + ".communications"
	AttachmentsTableName              = ServiceName + ".attachments"
	KnowledgeArticlesTableName        = ServiceName + ".knowledge_articles"
	ArticleRatingsTableName           = ServiceName + ".article_ratings"
	AIConfigurationsTableName         = ServiceName + ".ai_configurations"
	
	// Catalog tables
	SupportChannelsCatalogTableName   = ServiceName + ".support_channels_catalog"
	TicketsCatalogTableName           = ServiceName + ".tickets_catalog"
	KnowledgeArticlesCatalogTableName = ServiceName + ".knowledge_articles_catalog"
)