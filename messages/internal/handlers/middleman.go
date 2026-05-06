package handlers

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"middleman/internal/ddd"
	"middleman/internal/di"
	"middleman/internal/errorsotel"
	"middleman/messages/internal/constants"
	"middleman/messages/internal/domain"
	"time"
)

type middlemanHandlers[T ddd.Event] struct {
	middleman domain.MiddlemanRepository
}

var _ ddd.EventHandler[ddd.Event] = (*middlemanHandlers[ddd.Event])(nil)

func NewMiddlemanHandlers(middleman domain.MiddlemanRepository) ddd.EventHandler[ddd.Event] {
	return middlemanHandlers[ddd.Event]{
		middleman: middleman,
	}
}
func RegisterMiddlemanHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.ConversationStartedEvent,
		domain.ConversationDeletedEvent,
		domain.ConversationClosedEvent,
		domain.ConversationArchivedEvent,
	)
}
func RegisterMiddlemanHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		middlemanHandlers := di.Get(ctx, constants.MiddlemanHandlersKey).(ddd.EventHandler[ddd.Event])

		return middlemanHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.Event])

	RegisterMiddlemanHandlers(subscriber, handlers)
}

func (h middlemanHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling mall event",
				trace.WithAttributes(errorsotel.ErrAttrs(err)...),
			)
		}
		span.AddEvent("Handled mall event", trace.WithAttributes(
			attribute.Int64("TookMS", time.Since(started).Milliseconds()),
		))
	}(time.Now())

	span.AddEvent("Handling mall event", trace.WithAttributes(
		attribute.String("Event", event.EventName()),
	))

	switch event.EventName() {
	case domain.ConversationStartedEvent:
		return h.onConversationStarted(ctx, event)
	case domain.ConversationDeletedEvent:
		return h.onConversationDeleted(ctx, event)
	case domain.ConversationArchivedEvent:
		return h.ConversationArchived(ctx, event)
	case domain.ConversationClosedEvent:
		return h.onConversationClosed(ctx, event)

	}
	return nil
}
func (h middlemanHandlers[T]) onConversationStarted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Conversation)
	return h.middleman.Add(ctx, payload.ID(), payload.SenderID, payload.RecipientID, payload.ItemID)
}
func (h middlemanHandlers[T]) onConversationDeleted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Conversation)
	return h.middleman.Add(ctx, payload.ID(), payload.SenderID, payload.RecipientID, payload.ItemID)
}

func (h middlemanHandlers[T]) ConversationArchived(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Conversation)
	return h.middleman.Add(ctx, payload.ID(), payload.SenderID, payload.RecipientID, payload.ItemID)
}
func (h middlemanHandlers[T]) onConversationClosed(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Conversation)
	return h.middleman.Add(ctx, payload.ID(), payload.SenderID, payload.RecipientID, payload.ItemID)
}
