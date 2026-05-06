package constants

// ServiceName The name of this module/service
const ServiceName = "payments"

// GRPC Service Names
const (
// (removed unused UsersServiceName constant)
)

// Dependency Injection Keys
const (
	RegistryKey                 = "registry"
	DomainDispatcherKey         = "domainDispatcher"
	Stripe                      = "stripe"
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

	InvoicesRepoKey  = "invoicesRepo"
	RecurringRepoKey = "recurringRepo"
	PaymentsRepoKey  = "paymentsRepo"
)

// Repository Table Names
const (
	OutboxTableName    = ServiceName + ".outbox"
	InboxTableName     = ServiceName + ".inbox"
	EventsTableName    = ServiceName + ".events"
	SnapshotsTableName = ServiceName + ".snapshots"
	SagasTableName     = ServiceName + ".sagas"

	InvoicesTableName  = ServiceName + ".invoices"
	PaymentsTableName  = ServiceName + ".payments"
	RecurringTableName = ServiceName + ".recurrings"
)
