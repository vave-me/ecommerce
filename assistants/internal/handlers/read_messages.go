package handlers

import (
	"context"
	"middleman/assistants/internal/constants"
	"middleman/assistants/internal/domain"
	"middleman/internal/ddd"
	"middleman/internal/di"
	"middleman/internal/errorsotel"

	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type readMessagesHandlers[T ddd.Event] struct {
	readMessagesRepo domain.ReadMessagesRepository
}

var _ ddd.EventHandler[ddd.Event] = (*readMessagesHandlers[ddd.Event])(nil)

func NewReadMessagesHandlers(readMessagesRepo domain.ReadMessagesRepository) ddd.EventHandler[ddd.Event] {
	return readMessagesHandlers[ddd.Event]{
		readMessagesRepo: readMessagesRepo,
	}
}

func RegisterReadMessagesHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,

		domain.MessageAddedEvent,
	)
}

func RegisterReadMessagesHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		readMessagesHandlers := di.Get(ctx, constants.ReadMessagesHandlersKey).(ddd.EventHandler[ddd.Event])
		return readMessagesHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get("domainDispatcher").(*ddd.EventDispatcher[ddd.Event])
	RegisterReadMessagesHandlers(subscriber, handlers)
}

func (h readMessagesHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
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
	case domain.MessageAddedEvent:
		return h.onMessageAdded(ctx, event)

	}
	return nil
}

func (h readMessagesHandlers[T]) onMessageAdded(ctx context.Context, event ddd.Event) error {
	// Event contains the MessageAdded payload
	messageAdded := event.Payload().(*domain.MessageAdded)

	// Add the message to the read model
	return h.readMessagesRepo.AddMessage(
		ctx,
		messageAdded.ConversationID,
		messageAdded.ID,
		messageAdded.AssistantID,
		messageAdded.Role,
		messageAdded.Content,
		messageAdded.Timestamp,
		messageAdded.Metadata,
		messageAdded.ActionsTaken,
	)
}
