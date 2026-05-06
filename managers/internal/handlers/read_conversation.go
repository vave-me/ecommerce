package handlers

import (
	"context"
	"middleman/internal/ddd"
	"middleman/internal/di"
	"middleman/internal/errorsotel"
	"middleman/managers/internal/constants"
	"middleman/managers/internal/domain"

	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type readConversationHandlers[T ddd.Event] struct {
	readConversationRepo domain.ReadConversationRepository
}

var _ ddd.EventHandler[ddd.Event] = (*readConversationHandlers[ddd.Event])(nil)

func NewReadConversationHandlers(readConversationRepo domain.ReadConversationRepository) ddd.EventHandler[ddd.Event] {
	return readConversationHandlers[ddd.Event]{
		readConversationRepo: readConversationRepo,
	}
}

func RegisterReadConversationHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.ConversationCreatedEvent,
		domain.ConversationContextUpdatedEvent,
		domain.ConversationArchivedEvent,
		domain.MessageAddedEvent,
	)
}

func RegisterReadConversationHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		readConversationHandlers := di.Get(ctx, constants.ReadConversationHandlersKey).(ddd.EventHandler[ddd.Event])
		return readConversationHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get("domainDispatcher").(*ddd.EventDispatcher[ddd.Event])
	RegisterReadConversationHandlers(subscriber, handlers)
}

func (h readConversationHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling conversation read model event",
				trace.WithAttributes(errorsotel.ErrAttrs(err)...),
			)
		}
		span.AddEvent("Handled conversation read model event", trace.WithAttributes(
			attribute.Int64("TookMS", time.Since(started).Milliseconds()),
		))
	}(time.Now())

	span.AddEvent("Handling conversation read model event", trace.WithAttributes(
		attribute.String("Event", event.EventName()),
	))

	switch event.EventName() {
	case domain.ConversationCreatedEvent:
		return h.onConversationCreated(ctx, event)
	case domain.ConversationContextUpdatedEvent:
		return h.onConversationContextUpdated(ctx, event)
	case domain.ConversationArchivedEvent:
		return h.onConversationArchived(ctx, event)
	case domain.MessageAddedEvent:
		return h.onMessageAdded(ctx, event)
	}
	return nil
}

// Helper method to extract conversation ID from event
func (h readConversationHandlers[T]) getConversationIDFromEvent(event ddd.Event) string {
	// Try type assertion to get AggregateID if the concrete type supports it
	if eventWithAggregateID, ok := event.(interface{ AggregateID() string }); ok {
		return eventWithAggregateID.AggregateID()
	}

	return event.ID()
}

func (h readConversationHandlers[T]) onConversationCreated(ctx context.Context, event ddd.Event) error {
	// Event contains the ConversationCreated payload
	conversationCreated := event.Payload().(*domain.ConversationCreated)

	// Add conversation to the read model with correct data
	return h.readConversationRepo.AddConversation(
		ctx,
		conversationCreated.ConversationID,
		conversationCreated.UserID,
		conversationCreated.ManagerID,
		conversationCreated.CreatedAt,
		conversationCreated.Context,
	)
}

func (h readConversationHandlers[T]) onConversationContextUpdated(ctx context.Context, event ddd.Event) error {
	// Event contains the ConversationContextUpdated payload
	contextUpdated := event.Payload().(*domain.ConversationContextUpdated)

	// Update conversation context in the read model
	return h.readConversationRepo.UpdateConversationContext(
		ctx,
		contextUpdated.ConversationID,
		contextUpdated.Context,
		contextUpdated.UpdatedAt,
	)
}

func (h readConversationHandlers[T]) onConversationArchived(ctx context.Context, event ddd.Event) error {
	// Event contains the ConversationArchived payload
	archived := event.Payload().(*domain.ConversationArchived)

	// Archive conversation in the read model
	return h.readConversationRepo.ArchiveConversation(
		ctx,
		archived.ConversationID,
		archived.ArchivedAt,
	)
}

func (h readConversationHandlers[T]) onMessageAdded(ctx context.Context, event ddd.Event) error {
	// Event contains the MessageAdded payload
	messageAdded := event.Payload().(*domain.MessageAdded)

	// Update only the timestamp of the conversation (updated_at)
	return h.readConversationRepo.UpdateConversationTimestamp(
		ctx,
		messageAdded.ConversationID,
		messageAdded.Timestamp,
	)
}
