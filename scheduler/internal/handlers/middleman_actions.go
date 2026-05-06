package handlers

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"middleman/internal/ddd"
	"middleman/internal/di"
	"middleman/internal/errorsotel"
	"middleman/scheduler/internal/constants"
	"middleman/scheduler/internal/domain"
	"time"
)

type catalogHandlers[T ddd.Event] struct {
	catalog domain.MiddlemanCacheActionRepository
}

var _ ddd.EventHandler[ddd.Event] = (*catalogHandlers[ddd.Event])(nil)

func NewMiddlemanActionHandlers(catalog domain.MiddlemanCacheActionRepository) ddd.EventHandler[ddd.Event] {
	return catalogHandlers[ddd.Event]{
		catalog: catalog,
	}
}

func RegisterMiddlemanActionHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.ActionAddedEvent,
		domain.ActionUpdatedEvent,
		domain.ActionRemovedEvent,
	)
}

func RegisterMiddlemanActionHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		catalogHandlers := di.Get(ctx, constants.MiddlemanActionHandlersKey).(ddd.EventHandler[ddd.Event])

		return catalogHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.Event])

	RegisterMiddlemanActionHandlers(subscriber, handlers)
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
	case domain.ActionAddedEvent:
		return h.onActionAdded(ctx, event)
	case domain.ActionUpdatedEvent:
		return h.onActionUpdated(ctx, event)
	case domain.ActionRemovedEvent:
		return h.onActionRemoved(ctx, event)
	}
	return nil
}

func (h catalogHandlers[T]) onActionAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Action)
	return h.catalog.Add(ctx, payload.ID(), payload.SchedulerID, payload.NaturalLanguageTask, payload.ExecutionTime)
}

func (h catalogHandlers[T]) onActionUpdated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Action)
	return h.catalog.UpdateStatus(ctx, payload.ID(), payload.Status, payload.Result, payload.ErrorMessage)
}

func (h catalogHandlers[T]) onActionRemoved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Action)
	return h.catalog.Remove(ctx, payload.ID())
}
