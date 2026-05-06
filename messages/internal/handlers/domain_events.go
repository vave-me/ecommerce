package handlers

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"log"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/messages/internal/domain"
	"middleman/messages/messagespb"
	"time"
)

type domainHandlers[T ddd.Event] struct {
	publisher am.EventPublisher
}

var _ ddd.EventHandler[ddd.Event] = (*domainHandlers[ddd.Event])(nil)

// NewDomainEventHandlers initializes the domain event handlers with both publishers.
func NewDomainEventHandlers(publisher am.EventPublisher) ddd.EventHandler[ddd.Event] {
	return &domainHandlers[ddd.Event]{
		publisher: publisher,
	}
}

// RegisterDomainEventHandlers subscribes the handlers to the domain events.
func RegisterDomainEventHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	log.Println("[DomainHandlers] Registering Domain Event Handlers")
	subscriber.Subscribe(handlers,
		domain.MessageSentEvent,
		domain.MessageDeletedEvent,
	)
}

func (h domainHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling domain event",
				trace.WithAttributes(errorsotel.ErrAttrs(err)...),
			)
		}
		span.AddEvent("Handled domain event", trace.WithAttributes(
			attribute.Int64("TookMS", time.Since(started).Milliseconds()),
		))
	}(time.Now())

	span.AddEvent("Handling domain event", trace.WithAttributes(
		attribute.String("Event", event.EventName()),
	))

	switch event.EventName() {
	case domain.MessageSentEvent:
		err = h.onMessageSent(ctx, event)
	case domain.MessageDeletedEvent:
		err = h.onMessageDeleted(ctx, event)
	default:
		log.Printf("[DomainHandlers] Unhandled event: %s", event.EventName())
	}

	return err
}

// onMessageSent handles the MessageSentEvent.
func (h domainHandlers[T]) onMessageSent(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Message)

	return h.publisher.Publish(ctx, messagespb.MessageAggregateChannel,
		ddd.NewEvent(messagespb.MessageSentEvent, &messagespb.MessageSent{
			Id:             payload.ID(),
			ConversationId: payload.ConversationID,
			SenderId:       payload.SenderID,
			RecipientId:    payload.RecipientID,
			ItemId:         payload.ItemID,
			Body:           payload.Body,
			IsRead:         payload.IsRead,
		}),
	)

}

// onMessageDeleted handles the MessageDeletedEvent.
func (h domainHandlers[T]) onMessageDeleted(ctx context.Context, event ddd.Event) error {

	payload := event.Payload().(*domain.Message)
	return h.publisher.Publish(ctx, messagespb.MessageAggregateChannel,
		ddd.NewEvent(messagespb.MessageDeletedEvent, &messagespb.MessageDeleted{
			Id: payload.ID(),
		}),
	)

}
