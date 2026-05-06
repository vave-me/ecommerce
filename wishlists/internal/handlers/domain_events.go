package handlers

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/wishlists/internal/domain"
	"middleman/wishlists/wishlistspb"
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
		domain.WishlistCreatedEvent,
		domain.WishlistRemovedEvent,
		domain.WishlistItemAddedEvent,
		domain.WishlistItemRemovedEvent,
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
	case domain.WishlistCreatedEvent:
		return h.onWishlistCreated(ctx, event)
	case domain.WishlistRemovedEvent:
		return h.onWishlistRemoved(ctx, event)
	case domain.WishlistItemAddedEvent:
		return h.onWishlistItemAdded(ctx, event)
	case domain.WishlistItemRemovedEvent:
		return h.onWishlistItemRemoved(ctx, event)

	}
	return nil
}

func (h domainHandlers[T]) onWishlistCreated(ctx context.Context, event ddd.Event) error {
	wishlist := event.Payload().(*domain.Wishlist)
	return h.publisher.Publish(ctx, wishlistspb.WishlistAggregateChannel,
		ddd.NewEvent(wishlistspb.WishlistCreatedEvent, &wishlistspb.WishlistCreated{
			Id:     wishlist.ID(),
			UserId: wishlist.UserID,
		}),
	)
}

func (h domainHandlers[T]) onWishlistItemAdded(ctx context.Context, event ddd.Event) error {
	item := event.Payload().(*domain.WishlistItem)
	return h.publisher.Publish(ctx, wishlistspb.WishlistItemAggregateChannel,
		ddd.NewEvent(wishlistspb.WishlistItemAddedEvent, &wishlistspb.WishlistItemAdded{
			Id:         item.ID(),
			WishlistId: item.WishlistID,
			ItemId:     item.ItemID,
		}),
	)
}

func (h domainHandlers[T]) onWishlistItemRemoved(ctx context.Context, event ddd.Event) error {
	item := event.Payload().(*domain.WishlistItem)
	return h.publisher.Publish(ctx, wishlistspb.WishlistItemAggregateChannel,
		ddd.NewEvent(wishlistspb.WishlistItemRemovedEvent, &wishlistspb.WishlistItemRemoved{
			Id: item.ID(),
		}),
	)
}

func (h domainHandlers[T]) onWishlistRemoved(ctx context.Context, event ddd.Event) error {
	item := event.Payload().(*domain.Wishlist)
	return h.publisher.Publish(ctx, wishlistspb.WishlistAggregateChannel,
		ddd.NewEvent(wishlistspb.WishlistRemovedEvent, &wishlistspb.WishlistRemoved{
			Id: item.ID(),
		}),
	)
}
