package handlers

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"middleman/activity/activitypb"
	"middleman/activity/internal/domain"
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
		domain.ActivityCreatedEvent,
		domain.InteractionAddedEvent,
		domain.InteractionRemovedEvent,
		domain.InteractionUpdatedEvent,
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
	case domain.ActivityCreatedEvent:
		return h.onActivityCreated(ctx, event)
	case domain.InteractionAddedEvent:
		return h.onInteractionAdded(ctx, event)
	case domain.InteractionRemovedEvent:
		return h.onActivityCreated(ctx, event)
	case domain.InteractionUpdatedEvent:
		return h.onActivityCreated(ctx, event)
	}
	return nil
}

func (h domainHandlers[T]) onInteractionAdded(ctx context.Context, event ddd.Event) error {
	interaciton := event.Payload().(*domain.Interaction)
	return h.publisher.Publish(ctx, activitypb.InteractionAggregateChannel,
		ddd.NewEvent(activitypb.InteractionAddedEvent, &activitypb.InteractionAdded{
			ItemId:     interaciton.ItemID,
			ActivityId: interaciton.ActivityID,
			ItemType:   interaciton.ItemType,
			ActionType: interaciton.ActionType,
		}),
	)
}
func (h domainHandlers[T]) onActivityCreated(ctx context.Context, event ddd.Event) error {
	activity := event.Payload().(*domain.Activity)
	return h.publisher.Publish(ctx, activitypb.ActivityAggregateChannel,
		ddd.NewEvent(activitypb.ActivityCreatedEvent, &activitypb.ActivityCreated{
			Id:     activity.ID(),
			UserId: activity.UserID,
		}),
	)
}
