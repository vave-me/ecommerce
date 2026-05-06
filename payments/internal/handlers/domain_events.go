package handlers

import (
	"context"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/payments/internal/models"
	"middleman/payments/paymentspb"
	"time"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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
		models.InvoicePaidEvent,
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

	log.Debug().Str("event", event.EventName()).Msg("Domain event received for handling")

	switch event.EventName() {
	case models.InvoicePaidEvent:
		return h.onInvoicePaid(ctx, event)

	}
	return nil
}

func (h domainHandlers[T]) onInvoicePaid(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*models.InvoicePaid)
	log.Debug().Str("event", paymentspb.InvoicePaidEvent).Msg("Publishing integration event")
	return h.publisher.Publish(ctx, paymentspb.InvoiceAggregateChannel,
		ddd.NewEvent(paymentspb.InvoicePaidEvent, &paymentspb.InvoicePaid{
			Id:      payload.InvoiceID,
			OrderId: payload.OrderID,
		}),
	)
}
