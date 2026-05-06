package handlers

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/notifications/internal/domain"
	"middleman/notifications/notificationspb"
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
		domain.BasketAlertAddedEvent,
		domain.ProductAlertAddedEvent,
		domain.UserAlertAddedEvent,
		domain.CommentAlertAddedEvent,
		domain.InteractionAlertAddedEvent,
		domain.SupportAlertAddedEvent,
		domain.WishlistAlertAddedEvent,
		domain.MessageAlertAddedEvent,
		domain.OfferAlertAddedEvent,
		domain.OrderAlertAddedEvent,
		domain.ReviewAlertAddedEvent,
		domain.PaymentAlertAddedEvent,
		domain.FollowingAlertAddedEvent,
		domain.AlertReadEvent,
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
	case domain.MessageAlertAddedEvent:
		return h.onMessageAlertAdded(ctx, event)
	case domain.InteractionAlertAddedEvent:
		return h.onInteractionAlertAdded(ctx, event)
	case domain.OrderAlertAddedEvent:
		return h.onOrderAlertAdded(ctx, event)
	case domain.WishlistAlertAddedEvent:
		return h.onWishlistAlertAdded(ctx, event)
	case domain.SupportAlertAddedEvent:
		return h.onSupportAlertAdded(ctx, event)
	case domain.OfferAlertAddedEvent:
		return h.onOfferAlertAdded(ctx, event)
	case domain.ProductAlertAddedEvent:
		return h.onProductAlertAdded(ctx, event)
	case domain.UserAlertAddedEvent:
		return h.onUserAlertAdded(ctx, event)
	case domain.BasketAlertAddedEvent:
		return h.onBasketAlertAdded(ctx, event)
	case domain.CommentAlertAddedEvent:
		return h.onCommentAlertAdded(ctx, event)
	case domain.ReviewAlertAddedEvent:
		return h.onReviewAlertAdded(ctx, event)
	case domain.PaymentAlertAddedEvent:
		return h.onPaymentAlertAdded(ctx, event)
	case domain.FollowingAlertAddedEvent:
		return h.onFollowingAlertAdded(ctx, event)
	case domain.AlertReadEvent:
		return h.onAlertRead(ctx, event)
	}
	return nil
}

func (h domainHandlers[T]) onCommentAlertAdded(ctx context.Context, event ddd.Event) error {
	alert := event.Payload().(*domain.Alert)
	return h.publisher.Publish(ctx, notificationspb.AlertAggregateChannel,
		ddd.NewEvent(notificationspb.CommentAlertAddedEvent, &notificationspb.CommentAlertAdded{
			UserId:  alert.UserID,
			Message: alert.Message,
		}),
	)
}
func (h domainHandlers[T]) onBasketAlertAdded(ctx context.Context, event ddd.Event) error {
	alert := event.Payload().(*domain.Alert)
	return h.publisher.Publish(ctx, notificationspb.AlertAggregateChannel,
		ddd.NewEvent(notificationspb.BasketAlertAddedEvent, &notificationspb.BasketAlertAdded{
			UserId:  alert.UserID,
			Message: alert.Message,
		}),
	)
}
func (h domainHandlers[T]) onUserAlertAdded(ctx context.Context, event ddd.Event) error {
	alert := event.Payload().(*domain.Alert)
	return h.publisher.Publish(ctx, notificationspb.AlertAggregateChannel,
		ddd.NewEvent(notificationspb.UserAlertAddedEvent, &notificationspb.UserAlertAdded{
			UserId:  alert.UserID,
			Message: alert.Message,
		}),
	)
}

func (h domainHandlers[T]) onProductAlertAdded(ctx context.Context, event ddd.Event) error {
	alert := event.Payload().(*domain.Alert)
	return h.publisher.Publish(ctx, notificationspb.AlertAggregateChannel,
		ddd.NewEvent(notificationspb.ProductAlertAddedEvent, &notificationspb.ProductAlertAdded{
			UserId:  alert.UserID,
			Message: alert.Message,
		}),
	)
}
func (h domainHandlers[T]) onOfferAlertAdded(ctx context.Context, event ddd.Event) error {
	alert := event.Payload().(*domain.Alert)
	return h.publisher.Publish(ctx, notificationspb.AlertAggregateChannel,
		ddd.NewEvent(notificationspb.OfferAlertAddedEvent, &notificationspb.OfferAlertAdded{
			UserId:  alert.UserID,
			Message: alert.Message,
		}),
	)
}
func (h domainHandlers[T]) onSupportAlertAdded(ctx context.Context, event ddd.Event) error {
	alert := event.Payload().(*domain.Alert)
	return h.publisher.Publish(ctx, notificationspb.AlertAggregateChannel,
		ddd.NewEvent(notificationspb.SupportAlertAddedEvent, &notificationspb.SupportAlertAdded{
			UserId:  alert.UserID,
			Message: alert.Message,
		}),
	)
}
func (h domainHandlers[T]) onWishlistAlertAdded(ctx context.Context, event ddd.Event) error {
	alert := event.Payload().(*domain.Alert)
	return h.publisher.Publish(ctx, notificationspb.AlertAggregateChannel,
		ddd.NewEvent(notificationspb.WishlistAlertAddedEvent, &notificationspb.WishlistAlertAdded{
			UserId:  alert.UserID,
			Message: alert.Message,
		}),
	)
}
func (h domainHandlers[T]) onMessageAlertAdded(ctx context.Context, event ddd.Event) error {
	alert := event.Payload().(*domain.Alert)
	return h.publisher.Publish(ctx, notificationspb.AlertAggregateChannel,
		ddd.NewEvent(notificationspb.MessageAlertAddedEvent, &notificationspb.MessageAlertAdded{
			UserId:  alert.UserID,
			Message: alert.Message,
		}),
	)
}
func (h domainHandlers[T]) onInteractionAlertAdded(ctx context.Context, event ddd.Event) error {
	alert := event.Payload().(*domain.Alert)
	return h.publisher.Publish(ctx, notificationspb.AlertAggregateChannel,
		ddd.NewEvent(notificationspb.InteractionAlertAddedEvent, &notificationspb.InteractionAlertAdded{
			UserId:  alert.UserID,
			Message: alert.Message,
		}),
	)
}
func (h domainHandlers[T]) onOrderAlertAdded(ctx context.Context, event ddd.Event) error {
	alert := event.Payload().(*domain.Alert)
	return h.publisher.Publish(ctx, notificationspb.AlertAggregateChannel,
		ddd.NewEvent(notificationspb.OrderAlertAddedEvent, &notificationspb.OrderAlertAdded{
			UserId:  alert.UserID,
			Message: alert.Message,
		}),
	)
}

func (h domainHandlers[T]) onReviewAlertAdded(ctx context.Context, event ddd.Event) error {
	alert := event.Payload().(*domain.Alert)
	return h.publisher.Publish(ctx, notificationspb.AlertAggregateChannel,
		ddd.NewEvent(notificationspb.ReviewAlertAddedEvent, &notificationspb.ReviewAlertAdded{
			UserId:  alert.UserID,
			Message: alert.Message,
		}),
	)
}

func (h domainHandlers[T]) onPaymentAlertAdded(ctx context.Context, event ddd.Event) error {
	alert := event.Payload().(*domain.Alert)
	return h.publisher.Publish(ctx, notificationspb.AlertAggregateChannel,
		ddd.NewEvent(notificationspb.PaymentAlertAddedEvent, &notificationspb.PaymentAlertAdded{
			UserId:  alert.UserID,
			Message: alert.Message,
		}),
	)
}

func (h domainHandlers[T]) onFollowingAlertAdded(ctx context.Context, event ddd.Event) error {
	alert := event.Payload().(*domain.Alert)
	// Extract FollowerID from payload
	followerID := ""
	if fid, ok := alert.Payload["FollowerID"].(string); ok {
		followerID = fid
	}
	return h.publisher.Publish(ctx, notificationspb.AlertAggregateChannel,
		ddd.NewEvent(notificationspb.FollowingAlertAddedEvent, &notificationspb.FollowingAlertAdded{
			UserId:     alert.UserID,
			Message:    alert.Message,
			FollowerId: followerID,
		}),
	)
}

func (h domainHandlers[T]) onAlertRead(ctx context.Context, event ddd.Event) error {
	alert := event.Payload().(*domain.AlertRead)
	return h.publisher.Publish(ctx, notificationspb.AlertAggregateChannel,
		ddd.NewEvent(notificationspb.AlertReadEvent, &notificationspb.AlertRead{
			Id: alert.ID,
		}),
	)
}
