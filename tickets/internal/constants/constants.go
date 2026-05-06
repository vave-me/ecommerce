package constants

// ServiceName The name of this module/service
const ServiceName = "tickets"

// GRPC Service Names
const (
	TicketsServiceName = "TICKETS"
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
	MatchHandlersKey            = "matchHandlers"
	TicketHandlersKey           = "ticketHandlers"

	MatchesRepoKey      = "matchesRepo"
	TicketsRepoKey      = "ticketsRepo"
	MatchCatalogRepoKey = "matchCatalogRepo"
)

// Repository Table Names
const (
	OutboxTableName    = ServiceName + ".outbox"
	InboxTableName     = ServiceName + ".inbox"
	EventsTableName    = ServiceName + ".events"
	SnapshotsTableName = ServiceName + ".snapshots"
	SagasTableName     = ServiceName + ".sagas"

	MatchesTableName = ServiceName + ".matches"
	TicketsTableName = ServiceName + ".tickets"
	SectorsTableName = ServiceName + ".sectors"
	SeatsTableName   = ServiceName + ".seats"
)

// Event store streams
const (
	MatchEventStream  = "match-events"
	TicketEventStream = "ticket-events"
)

// Default values
const (
	DefaultPageSize         = 20
	MaxPageSize             = 100
	DefaultReservationTime  = 15 * 60 // 15 minutes in seconds
	MaxTicketsPerOrder      = 10
	MaxTransfersPerTicket   = 5
	TicketExpirationHours   = 24 // Hours after match
)

// Cache TTL values
const (
	MatchCacheTTL    = 300  // 5 minutes
	SectorCacheTTL   = 600  // 10 minutes
	SeatMapCacheTTL  = 60   // 1 minute
	TicketCacheTTL   = 300  // 5 minutes
)

// Ticket validation
const (
	ValidationSuccess = "success"
	ValidationFailed  = "failed"
	ValidationExpired = "expired"
	ValidationUsed    = "used"
)

// Error messages
const (
	ErrMatchNotFound   = "match not found"
	ErrTicketNotFound  = "ticket not found"
	ErrSeatNotAvailable = "seat not available"
	ErrInvalidQRCode   = "invalid QR code"
	ErrAccessDenied    = "access denied"
)