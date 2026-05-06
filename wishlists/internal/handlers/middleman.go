package handlers

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"middleman/internal/ddd"
	"middleman/internal/di"
	"middleman/internal/errorsotel"
	"middleman/wishlists/internal/constants"
	"middleman/wishlists/internal/domain"
	"time"
)

type middlemanHandlers[T ddd.Event] struct {
	middleman domain.MiddlemanRepository
}

var _ ddd.EventHandler[ddd.Event] = (*middlemanHandlers[ddd.Event])(nil)

func NewMiddlemanHandlers(middleman domain.MiddlemanRepository) ddd.EventHandler[ddd.Event] {
	return middlemanHandlers[ddd.Event]{
		middleman: middleman,
	}
}

func RegisterMiddlemanHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.WishlistCreatedEvent,
		domain.WishlistRemovedEvent,
	)
}

func RegisterMiddlemanHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		middlemanHandlers := di.Get(ctx, constants.MiddlemanHandlersKey).(ddd.EventHandler[ddd.Event])

		return middlemanHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.Event])

	RegisterMiddlemanHandlers(subscriber, handlers)
}

func (h middlemanHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling mall event",
				trace.WithAttributes(errorsotel.ErrAttrs(err)...),
			)
		}
		span.AddEvent("Handled mall event", trace.WithAttributes(
			attribute.Int64("TookMS", time.Since(started).Milliseconds()),
		))
	}(time.Now())

	span.AddEvent("Handling mall event", trace.WithAttributes(
		attribute.String("Event", event.EventName()),
	))

	switch event.EventName() {
	case domain.WishlistCreatedEvent:
		return h.onWishlistCreated(ctx, event)
	case domain.WishlistRemovedEvent:
		return h.onWishlistRemoved(ctx, event)
	}
	return nil
}

func (h middlemanHandlers[T]) onWishlistCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Wishlist)
	return h.middleman.AddWishlist(ctx, payload.ID(), payload.UserID, payload.Name)
}

func (h middlemanHandlers[T]) onWishlistRemoved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Wishlist)
	return h.middleman.Remove(ctx, payload.ID())
}
