package constants

// ServiceName The name of this module/service
const ServiceName = "messages"

// GRPC Service Names
const (
	MessagesServiceName = "MESSAGES"
)

// Dependency Injection Keys
const (
	RegistryKey                 = "registry"
	DomainDispatcherKey         = "domainDispatcher"
	DatabaseTransactionKey      = "tx"
	MessagePublisherKey         = "messagePublisher"
	MessageSubscriberKey        = "messageSubscriber"
	WebSocketPublisherKey       = "webSocketPublisher"
	EventPublisherKey           = "eventPublisher"
	CommandPublisherKey         = "commandPublisher"
	ReplyPublisherKey           = "replyPublisher"
	AggregateStoreKey           = "aggregateStore"
	Logger                      = "logger"
	SagaStoreKey                = "sagaStore"
	InboxStoreKey               = "inboxStore"
	ApplicationKey              = "app"
	DomainEventHandlersKey      = "domainEventHandlers"
	IntegrationEventHandlersKey = "integrationEventHandlers"
	WebsocketsEventHandlersKey  = "websocketCommandHandlers"
	WebSocketSubscriberKey      = "websocketSubscribers"
	CommandHandlersKey          = "commandHandlers"
	ReplyHandlersKey            = "replyHandlers"

	MessengerHandlersKey = "messengerHandlers"
	WebsocketHandlersKey = "websocketHandlers"
	MiddlemanHandlersKey = "middlemanHandlers"

	ConversationsRepoKey = "conversationsRepo"
	MessagesRepoKey      = "messagesRepo"
	MessengerRepoKey     = "messengerRepo"
	MiddlemanRepoKey     = "middlemanRepo"
)

// Repository Table Names
const (
	OutboxTableName    = ServiceName + ".outbox"
	InboxTableName     = ServiceName + ".inbox"
	EventsTableName    = ServiceName + ".events"
	SnapshotsTableName = ServiceName + ".snapshots"
	SagasTableName     = ServiceName + ".sagas"

	MessengerTableName = ServiceName + ".messages"
	MiddlemanTableName = ServiceName + ".conversations"
)
