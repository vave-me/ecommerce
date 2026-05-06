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

type messengerHandlers[T ddd.Event] struct {
	messenger domain.MessengerRepository
}

var _ ddd.EventHandler[ddd.Event] = (*messengerHandlers[ddd.Event])(nil)

func NewMessengerHandlers(messenger domain.MessengerRepository) ddd.EventHandler[ddd.Event] {
	return messengerHandlers[ddd.Event]{
		messenger: messenger,
	}
}

func RegisterMessengerHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.MessageSentEvent,
		domain.MessageDeletedEvent,
	)

}
func RegisterMessengerHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		messengerHandlers := di.Get(ctx, constants.MessengerHandlersKey).(ddd.EventHandler[ddd.Event])

		return messengerHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.Event])

	RegisterMessengerHandlers(subscriber, handlers)
}
func (h messengerHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {

	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling catalog event",
				trace.WithAttributes(errorsotel.ErrAttrs(err)...),
			)
		}
		span.AddEvent("Handled catalog event", trace.WithAttributes(
			attribute.Int64("TookMS", time.Since(started).Milliseconds()),
		))
	}(time.Now())

	span.AddEvent("Handling catalog event", trace.WithAttributes(
		attribute.String("Event", event.EventName()),
	))

	switch event.EventName() {
	case domain.MessageSentEvent:
		return h.onMessageSent(ctx, event)
	case domain.MessageDeletedEvent:
		return h.onMessageDeleted(ctx, event)

	}
	return nil
}
func (h messengerHandlers[T]) onMessageSent(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Message)
	return h.messenger.SendMessage(ctx,
		payload.ID(),
		payload.ConversationID,
		payload.SenderID,
		payload.RecipientID,
		payload.ItemID,
		payload.Body,
		payload.IsRead)
}
func (h messengerHandlers[T]) onMessageDeleted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Message)
	return h.messenger.Delete(ctx, payload.ID())
}
