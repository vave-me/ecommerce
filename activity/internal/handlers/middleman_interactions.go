package handlers

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"middleman/activity/internal/constants"
	"middleman/activity/internal/domain"
	"middleman/internal/ddd"
	"middleman/internal/di"
	"middleman/internal/errorsotel"
	"time"
)

type catalogHandlers[T ddd.Event] struct {
	catalog domain.MiddlemanCacheInteractionRepository
}

var _ ddd.EventHandler[ddd.Event] = (*catalogHandlers[ddd.Event])(nil)

func NewMiddlemanInteractionHandlers(catalog domain.MiddlemanCacheInteractionRepository) ddd.EventHandler[ddd.Event] {
	return catalogHandlers[ddd.Event]{
		catalog: catalog,
	}
}

func RegisterMiddlemanInteractionHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.InteractionAddedEvent,
		domain.InteractionUpdatedEvent,
		domain.InteractionRemovedEvent,
	)
}

func RegisterMiddlemanInteractionHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		catalogHandlers := di.Get(ctx, constants.MiddlemanInteractionHandlersKey).(ddd.EventHandler[ddd.Event])

		return catalogHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.Event])

	RegisterMiddlemanInteractionHandlers(subscriber, handlers)
}

func (h catalogHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
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
	case domain.InteractionAddedEvent:
		return h.onInteractionAdded(ctx, event)
	case domain.InteractionUpdatedEvent:
		return h.onInteractionUpdated(ctx, event)
	case domain.InteractionRemovedEvent:
		return h.onInteractionRemoved(ctx, event)
	}
	return nil
}

func (h catalogHandlers[T]) onInteractionAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Interaction)
	return h.catalog.Add(ctx, payload.ID(), payload.ActivityID, payload.ItemID, payload.ItemType, payload.ActionType)
}

func (h catalogHandlers[T]) onInteractionUpdated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Interaction)
	return h.catalog.Update(ctx, payload.ID(), payload.ActionType)
}

func (h catalogHandlers[T]) onInteractionRemoved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Interaction)
	return h.catalog.Remove(ctx, payload.ID())
}
