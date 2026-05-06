package handlers

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"middleman/categories/categoriespb"
	"middleman/categories/internal/domain"
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
	return &domainHandlers[ddd.Event]{publisher: publisher}
}

func RegisterDomainEventHandlers(
	subscriber ddd.EventSubscriber[ddd.Event],
	handlers ddd.EventHandler[ddd.Event],
) {
	subscriber.Subscribe(handlers,
		domain.CategoryAddedEvent,
		domain.CategoryUpdatedEvent,
		domain.CategoryRebrandedEvent,
		domain.CategoryArchivedEvent,
		domain.CategoryRemovedEvent,

		domain.FilterAddedEvent,
		domain.FilterRebrandedEvent,
		domain.FilterArchivedEvent,
		domain.FilterRemovedEvent,
	)
}

func (h domainHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent("Error handling domain event",
				trace.WithAttributes(errorsotel.ErrAttrs(err)...),
			)
		}
		span.AddEvent("Handled domain event", trace.WithAttributes(
			attribute.Int64("DurationMS", time.Since(started).Milliseconds()),
		))
	}(time.Now())

	switch event.EventName() {

	// Category events
	case domain.CategoryAddedEvent:
		err = h.onCategoryAdded(ctx, event)
	case domain.CategoryRebrandedEvent:
		err = h.onCategoryRebranded(ctx, event)
	case domain.CategoryArchivedEvent:
		err = h.onCategoryArchived(ctx, event)
	case domain.CategoryRemovedEvent:
		err = h.onCategoryRemoved(ctx, event)

	// Filter events
	case domain.FilterAddedEvent:
		err = h.onFilterAdded(ctx, event)
	case domain.FilterArchivedEvent:
		err = h.onFilterArchived(ctx, event)
	case domain.FilterRemovedEvent:
		err = h.onFilterRemoved(ctx, event)
	}
	return err
}

// Publish CategoryAdded on the "CategoryAggregateChannel"
func (h domainHandlers[T]) onCategoryAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.CategoryAdded)
	return h.publisher.Publish(
		ctx,
		categoriespb.CategoryAggregateChannel,
		ddd.NewEvent(categoriespb.CategoryAddedEvent, &categoriespb.CategoryAdded{
			Id:          payload.CategoryID,
			Description: payload.Description,
		}),
	)
}

func (h domainHandlers[T]) onCategoryRebranded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.CategoryRebranded)
	return h.publisher.Publish(
		ctx,
		categoriespb.CategoryAggregateChannel,
		ddd.NewEvent(categoriespb.CategoryRebrandedEvent, &categoriespb.CategoryRebranded{
			Id: payload.Description,
		}),
	)
}

func (h domainHandlers[T]) onCategoryArchived(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.CategoryArchived)
	return h.publisher.Publish(
		ctx,
		categoriespb.CategoryAggregateChannel,
		ddd.NewEvent(categoriespb.CategoryArchivedEvent, &categoriespb.CategoryArchived{
			Id:     payload.CategoryID,
			UserId: payload.CategoryID,
		}),
	)
}

func (h domainHandlers[T]) onCategoryRemoved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.CategoryRemoved)
	return h.publisher.Publish(
		ctx,
		categoriespb.CategoryAggregateChannel,
		ddd.NewEvent(categoriespb.CategoryRemovedEvent, &categoriespb.CategoryRemoved{
			Id: payload.CategoryID,
		}),
	)
}

// Filter events

func (h domainHandlers[T]) onFilterAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.FilterAdded)
	return h.publisher.Publish(
		ctx,
		categoriespb.FilterAggregateChannel,
		ddd.NewEvent(categoriespb.FilterAddedEvent, &categoriespb.FilterAdded{
			Id:         payload.CategoryID,
			CategoryId: payload.CategoryID,
			Name:       payload.Name,
			// etc
			IsActive: payload.IsActive,
		}),
	)
}

func (h domainHandlers[T]) onFilterArchived(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.FilterArchived)
	return h.publisher.Publish(
		ctx,
		categoriespb.FilterAggregateChannel,
		ddd.NewEvent(categoriespb.FilterArchivedEvent, &categoriespb.FilterArchived{
			Id: payload.FilterID,
		}),
	)
}

func (h domainHandlers[T]) onFilterRemoved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.FilterRemoved)
	return h.publisher.Publish(
		ctx,
		categoriespb.FilterAggregateChannel,
		ddd.NewEvent(categoriespb.FilterRemovedEvent, &categoriespb.FilterRemoved{
			Id: payload.FilterID,
		}),
	)
}
