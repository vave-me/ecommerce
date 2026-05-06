package handlers

import (
	"context"
	"middleman/internal/ddd"
	"middleman/internal/di"
	"middleman/internal/errorsotel"
	"middleman/wishlists/internal/constants"
	"middleman/wishlists/internal/domain"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type catalogHandlers[T ddd.Event] struct {
	catalog domain.CatalogRepository
}

var _ ddd.EventHandler[ddd.Event] = (*catalogHandlers[ddd.Event])(nil)

func NewCatalogHandlers(catalog domain.CatalogRepository) ddd.EventHandler[ddd.Event] {
	return catalogHandlers[ddd.Event]{
		catalog: catalog,
	}
}

func RegisterCatalogHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.WishlistItemAddedEvent,
		domain.WishlistItemRemovedEvent,
	)
}

func RegisterCatalogHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		catalogHandlers := di.Get(ctx, constants.CatalogHandlersKey).(ddd.EventHandler[ddd.Event])

		return catalogHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.Event])

	RegisterCatalogHandlers(subscriber, handlers)
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
	case domain.WishlistItemAddedEvent:
		return h.onWishlistItemAdded(ctx, event)
	case domain.WishlistItemRemovedEvent:
		return h.onWishlistItemRemoved(ctx, event)

	}
	return nil
}

func (h catalogHandlers[T]) onWishlistItemAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.WishlistItem)
	return h.catalog.AddWishlistItem(ctx, payload.ID(), payload.WishlistID, payload.ItemID, payload.EntityType)
}

func (h catalogHandlers[T]) onWishlistItemRemoved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.WishlistItem)
	return h.catalog.RemoveWishlistItem(ctx, payload.ID())
}
