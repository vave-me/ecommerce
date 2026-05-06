package handlers

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"middleman/geocoding/geocodingpb"
	"middleman/geocoding/internal/domain"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
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
		domain.AddressCreatedEvent,
		domain.LocationAddedEvent,
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

	case domain.AddressCreatedEvent:
		return h.onAddressCreated(ctx, event)
	case domain.LocationAddedEvent:
		return h.onLocationAdded(ctx, event)
	}
	return nil

}

func (h domainHandlers[T]) onAddressCreated(ctx context.Context, event ddd.Event) error {
	//address := event.Payload().(*domain.Address)
	return h.publisher.Publish(ctx, geocodingpb.AddressAggregateChannel,
		ddd.NewEvent(geocodingpb.AddressCreatedEvent, &geocodingpb.AddressCreated{}),
	)
}

func (h domainHandlers[T]) onLocationAdded(ctx context.Context, event ddd.Event) error {
	//	location := event.Payload().(*domain.Location)
	return h.publisher.Publish(ctx, geocodingpb.AddressAggregateChannel,
		ddd.NewEvent(geocodingpb.LocationAddedEvent, &geocodingpb.LocationAdded{}),
	)
}
