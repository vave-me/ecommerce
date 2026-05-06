package handlers

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/shipping/internal/domain"
	"middleman/shipping/shippingpb"
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
		domain.ShipmentCreatedEvent,
		domain.CarrierAssignedEvent,
		domain.ShipmentStartedEvent,
		domain.ShipmentStatusUpdatedEvent,
		domain.ShipmentCancelledEvent,
		domain.ShipmentDeliveredEvent,
		domain.PickupScheduledEvent,
		domain.ShipmentReturnedEvent,
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
	case domain.ShipmentCreatedEvent:
		return h.onShippingCreated(ctx, event)
	case domain.CarrierAssignedEvent:
		return h.onCarrierAssigned(ctx, event)
	case domain.ShipmentStartedEvent:
		return h.onShipmentStarted(ctx, event)
	case domain.ShipmentStatusUpdatedEvent:
		return h.onShipmentStatusUpdated(ctx, event)
	case domain.ShipmentCancelledEvent:
		return h.onShipmentCancelled(ctx, event)
	case domain.ShipmentDeliveredEvent:
		return h.onShipmentDelivered(ctx, event)
	case domain.PickupScheduledEvent:
		return h.onPickupScheduled(ctx, event)
	case domain.ShipmentReturnedEvent:
		return h.onShipmentReturned(ctx, event)
	}
	return nil
}

func (h domainHandlers[T]) onShippingCreated(ctx context.Context, event ddd.Event) error {
	shipment := event.Payload().(*domain.Shipment)
	return h.publisher.Publish(ctx, shippingpb.ShipmentAggregateChannel,
		ddd.NewEvent(shippingpb.ShipmentCreatedEvent, &shippingpb.ShipmentCreated{
			Id:             shipment.ID(),
			ProductId:      shipment.ProductID,
			TrackingNumber: shipment.TrackingNumber,
		}),
	)
}

func (h domainHandlers[T]) onCarrierAssigned(ctx context.Context, event ddd.Event) error {
	shipment := event.Payload().(*domain.Shipment)
	return h.publisher.Publish(ctx, shippingpb.ShipmentAggregateChannel,
		ddd.NewEvent(shippingpb.CarrierAssignedEvent, &shippingpb.CarrierAssigned{
			ShipmentId:  shipment.ID(),
			CarrierId:   shipment.CarrierID,
			CarrierName: shipment.CarrierName,
		}),
	)
}

func (h domainHandlers[T]) onShipmentStarted(ctx context.Context, event ddd.Event) error {
	shipment := event.Payload().(*domain.Shipment)
	return h.publisher.Publish(ctx, shippingpb.ShipmentAggregateChannel,
		ddd.NewEvent(shippingpb.ShipmentStartedEvent, &shippingpb.ShipmentStarted{
			ShipmentId: shipment.ID(),
		}),
	)
}

func (h domainHandlers[T]) onShipmentStatusUpdated(ctx context.Context, event ddd.Event) error {
	shipment := event.Payload().(*domain.Shipment)
	return h.publisher.Publish(ctx, shippingpb.ShipmentAggregateChannel,
		ddd.NewEvent(shippingpb.ShipmentStatusUpdatedEvent, &shippingpb.ShipmentStatusUpdated{
			ShipmentId: shipment.ID(),
			Status:     shipment.Status,
		}),
	)
}

func (h domainHandlers[T]) onShipmentCancelled(ctx context.Context, event ddd.Event) error {
	shipment := event.Payload().(*domain.Shipment)
	return h.publisher.Publish(ctx, shippingpb.ShipmentAggregateChannel,
		ddd.NewEvent(shippingpb.ShipmentCancelledEvent, &shippingpb.ShipmentCancelled{
			ShipmentId: shipment.ID(),
		}),
	)
}

func (h domainHandlers[T]) onShipmentDelivered(ctx context.Context, event ddd.Event) error {
	shipment := event.Payload().(*domain.Shipment)
	return h.publisher.Publish(ctx, shippingpb.ShipmentAggregateChannel,
		ddd.NewEvent(shippingpb.ShipmentDeliveredEvent, &shippingpb.ShipmentDelivered{
			ShipmentId: shipment.ID(),
		}),
	)
}

func (h domainHandlers[T]) onPickupScheduled(ctx context.Context, event ddd.Event) error {
	shipment := event.Payload().(*domain.Shipment)
	return h.publisher.Publish(ctx, shippingpb.ShipmentAggregateChannel,
		ddd.NewEvent(shippingpb.PickupScheduledEvent, &shippingpb.PickupScheduled{
			ShipmentId: shipment.ID(),
		}),
	)
}

func (h domainHandlers[T]) onShipmentReturned(ctx context.Context, event ddd.Event) error {
	shipment := event.Payload().(*domain.Shipment)
	return h.publisher.Publish(ctx, shippingpb.ShipmentAggregateChannel,
		ddd.NewEvent(shippingpb.ShipmentReturnedEvent, &shippingpb.ShipmentReturned{
			ShipmentId: shipment.ID(),
		}),
	)
}
