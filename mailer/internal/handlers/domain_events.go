package handlers

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/mailer/internal/domain"
	"middleman/mailer/mailerpb"
	"time"
)

type domainHandlers[T ddd.Event] struct {
	publisher am.EventPublisher
}

var _ ddd.EventHandler[ddd.Event] = (*domainHandlers[ddd.Event])(nil)

func NewDomainEventHandlers(publisher am.EventPublisher) ddd.EventHandler[ddd.Event] {
	return &domainHandlers[ddd.Event]{
		publisher: publisher,
	}
}

func RegisterDomainEventHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.EmailCreatedEvent,
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
	case domain.EmailCreatedEvent:
		return h.onMailerCreated(ctx, event)

	}
	return nil
}

func (h domainHandlers[T]) onMailerCreated(ctx context.Context, event ddd.Event) error {
	mailer := event.Payload().(*domain.Email)
	return h.publisher.Publish(ctx, mailerpb.EmailAggregateChannel,
		ddd.NewEvent(mailerpb.EmailCreatedEvent, &mailerpb.EmailCreated{
			Id:        mailer.ID(),
			SenderId:  mailer.SenderID,
			Recipient: mailer.Recipient,
			Subject:   mailer.Subject,
			Body:      mailer.Body,
			Status:    string(mailer.Status),
		}),
	)
}
