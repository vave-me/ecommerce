package constants

// ServiceName The name of this module/service
const ServiceName = "newsletters"

// GRPC Service Names
const (
	StoresServiceName = "STORES"
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

	// Repository keys
	NewslettersRepoKey         = "newslettersRepo"
	SubscriptionsRepoKey       = "subscriptionsRepo"
	EditionsRepoKey            = "editionsRepo"
	TemplatesRepoKey           = "templatesRepo"
	NewsletterCatalogRepoKey   = "newsletterCatalogRepo"
	SubscriptionCatalogRepoKey = "subscriptionCatalogRepo"
	EditionCatalogRepoKey      = "editionCatalogRepo"
	TemplateCatalogRepoKey     = "templateCatalogRepo"
	MiddlemanRepoKey           = "middlemanRepo"
)

// Repository Table Names
const (
	OutboxTableName    = ServiceName + ".outbox"
	InboxTableName     = ServiceName + ".inbox"
	EventsTableName    = ServiceName + ".events"
	SnapshotsTableName = ServiceName + ".snapshots"
	SagasTableName     = ServiceName + ".sagas"

	// Catalog table names
	NewslettersTableName   = ServiceName + ".newsletters"
	SubscriptionsTableName = ServiceName + ".newsletter_subscriptions"
	EditionsTableName      = ServiceName + ".newsletter_editions"
	TemplatesTableName     = ServiceName + ".newsletter_templates"
	MiddlemanTableName     = ServiceName + ".newsletters"
)