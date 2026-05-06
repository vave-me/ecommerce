package handlers

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"middleman/notifications/internal/constants"
	"middleman/notifications/internal/domain"

	"middleman/internal/ddd"
	"middleman/internal/di"
	"middleman/internal/errorsotel"
	"time"
)

type middlemanAlertHandlers[T ddd.Event] struct {
	catalog domain.CatalogRepository
}

var _ ddd.EventHandler[ddd.Event] = (*middlemanAlertHandlers[ddd.Event])(nil)

func NewMiddlemanHandlers(catalog domain.CatalogRepository) ddd.EventHandler[ddd.Event] {
	return middlemanAlertHandlers[ddd.Event]{
		catalog: catalog,
	}
}

func RegisterMiddlemanHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.MessageAlertAddedEvent,
		domain.CommentAlertAddedEvent,
		domain.InteractionAlertAddedEvent,
		domain.ProductAlertAddedEvent,
		domain.WishlistAlertAddedEvent,
		domain.BasketAlertAddedEvent,
		domain.OfferAlertAddedEvent,
		domain.OrderAlertAddedEvent,
		domain.SupportAlertAddedEvent,
		domain.UserAlertAddedEvent,
		domain.ReviewAlertAddedEvent,
		domain.PaymentAlertAddedEvent,
		domain.FollowingAlertAddedEvent,
		domain.AlertReadEvent,
	)
}

func RegisterMiddlemanHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		catalogHandlers := di.Get(ctx, constants.MiddlemanHandlersKey).(ddd.EventHandler[ddd.Event])

		return catalogHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.Event])

	RegisterMiddlemanHandlers(subscriber, handlers)
}

func (h middlemanAlertHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
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
	case domain.MessageAlertAddedEvent:
		return h.onMessageAlertAdded(ctx, event)
	case domain.CommentAlertAddedEvent:
		return h.onCommentAlertAdded(ctx, event)
	case domain.InteractionAlertAddedEvent:
		return h.onInteractionAlertAdded(ctx, event)
	case domain.OfferAlertAddedEvent:
		return h.onOfferAlertAdded(ctx, event)
	case domain.OrderAlertAddedEvent:
		return h.onOrderAlertAdded(ctx, event)
	case domain.SupportAlertAddedEvent:
		return h.onSupportAlertAdded(ctx, event)
	case domain.ProductAlertAddedEvent:
		return h.onProductAlertAdded(ctx, event)
	case domain.UserAlertAddedEvent:
		return h.onUserAlertAdded(ctx, event)
	case domain.WishlistAlertAddedEvent:
		return h.onWishlistAlertAdded(ctx, event)
	case domain.BasketAlertAddedEvent:
		return h.onBasketAlertAdded(ctx, event)
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

func (h middlemanAlertHandlers[T]) onMessageAlertAdded(ctx context.Context, event ddd.Event) error {
	alert := event.Payload().(*domain.Alert)
	var alertID string
	if aggEvent, ok := event.(ddd.AggregateEvent); ok {
		alertID = aggEvent.AggregateID()
	}
	return h.catalog.Add(ctx, alertID, alert.UserID, alert.AlertType.String(), alert.Message, alert.Payload, alert.IsRead)
}

func (h middlemanAlertHandlers[T]) onCommentAlertAdded(ctx context.Context, event ddd.Event) error {
	alert := event.Payload().(*domain.Alert)
	var alertID string
	if aggEvent, ok := event.(ddd.AggregateEvent); ok {
		alertID = aggEvent.AggregateID()
	}
	return h.catalog.Add(ctx, alertID, alert.UserID, alert.AlertType.String(), alert.Message, alert.Payload, alert.IsRead)
}

func (h middlemanAlertHandlers[T]) onInteractionAlertAdded(ctx context.Context, event ddd.Event) error {
	alert := event.Payload().(*domain.Alert)
	var alertID string
	if aggEvent, ok := event.(ddd.AggregateEvent); ok {
		alertID = aggEvent.AggregateID()
	}
	return h.catalog.Add(ctx, alertID, alert.UserID, alert.AlertType.String(), alert.Message, alert.Payload, alert.IsRead)
}

func (h middlemanAlertHandlers[T]) onOfferAlertAdded(ctx context.Context, event ddd.Event) error {
	alert := event.Payload().(*domain.Alert)
	var alertID string
	if aggEvent, ok := event.(ddd.AggregateEvent); ok {
		alertID = aggEvent.AggregateID()
	}
	return h.catalog.Add(ctx, alertID, alert.UserID, alert.AlertType.String(), alert.Message, alert.Payload, alert.IsRead)
}

func (h middlemanAlertHandlers[T]) onOrderAlertAdded(ctx context.Context, event ddd.Event) error {
	alert := event.Payload().(*domain.Alert)
	var alertID string
	if aggEvent, ok := event.(ddd.AggregateEvent); ok {
		alertID = aggEvent.AggregateID()
	}
	return h.catalog.Add(ctx, alertID, alert.UserID, alert.AlertType.String(), alert.Message, alert.Payload, alert.IsRead)
}

func (h middlemanAlertHandlers[T]) onSupportAlertAdded(ctx context.Context, event ddd.Event) error {
	alert := event.Payload().(*domain.Alert)
	var alertID string
	if aggEvent, ok := event.(ddd.AggregateEvent); ok {
		alertID = aggEvent.AggregateID()
	}
	return h.catalog.Add(ctx, alertID, alert.UserID, alert.AlertType.String(), alert.Message, alert.Payload, alert.IsRead)
}

func (h middlemanAlertHandlers[T]) onProductAlertAdded(ctx context.Context, event ddd.Event) error {
	alert := event.Payload().(*domain.Alert)
	var alertID string
	if aggEvent, ok := event.(ddd.AggregateEvent); ok {
		alertID = aggEvent.AggregateID()
	}
	return h.catalog.Add(ctx, alertID, alert.UserID, alert.AlertType.String(), alert.Message, alert.Payload, alert.IsRead)
}

func (h middlemanAlertHandlers[T]) onUserAlertAdded(ctx context.Context, event ddd.Event) error {
	alert := event.Payload().(*domain.Alert)
	var alertID string
	if aggEvent, ok := event.(ddd.AggregateEvent); ok {
		alertID = aggEvent.AggregateID()
	}
	return h.catalog.Add(ctx, alertID, alert.UserID, alert.AlertType.String(), alert.Message, alert.Payload, alert.IsRead)
}

func (h middlemanAlertHandlers[T]) onWishlistAlertAdded(ctx context.Context, event ddd.Event) error {
	alert := event.Payload().(*domain.Alert)
	var alertID string
	if aggEvent, ok := event.(ddd.AggregateEvent); ok {
		alertID = aggEvent.AggregateID()
	}
	return h.catalog.Add(ctx, alertID, alert.UserID, alert.AlertType.String(), alert.Message, alert.Payload, alert.IsRead)
}

func (h middlemanAlertHandlers[T]) onBasketAlertAdded(ctx context.Context, event ddd.Event) error {
	alert := event.Payload().(*domain.Alert)
	var alertID string
	if aggEvent, ok := event.(ddd.AggregateEvent); ok {
		alertID = aggEvent.AggregateID()
	}
	return h.catalog.Add(ctx, alertID, alert.UserID, alert.AlertType.String(), alert.Message, alert.Payload, alert.IsRead)
}

func (h middlemanAlertHandlers[T]) onReviewAlertAdded(ctx context.Context, event ddd.Event) error {
	alert := event.Payload().(*domain.Alert)
	var alertID string
	if aggEvent, ok := event.(ddd.AggregateEvent); ok {
		alertID = aggEvent.AggregateID()
	}
	return h.catalog.Add(ctx, alertID, alert.UserID, alert.AlertType.String(), alert.Message, alert.Payload, alert.IsRead)
}

func (h middlemanAlertHandlers[T]) onPaymentAlertAdded(ctx context.Context, event ddd.Event) error {
	alert := event.Payload().(*domain.Alert)
	var alertID string
	if aggEvent, ok := event.(ddd.AggregateEvent); ok {
		alertID = aggEvent.AggregateID()
	}
	return h.catalog.Add(ctx, alertID, alert.UserID, alert.AlertType.String(), alert.Message, alert.Payload, alert.IsRead)
}

func (h middlemanAlertHandlers[T]) onFollowingAlertAdded(ctx context.Context, event ddd.Event) error {
	alert := event.Payload().(*domain.Alert)
	var alertID string
	if aggEvent, ok := event.(ddd.AggregateEvent); ok {
		alertID = aggEvent.AggregateID()
	}
	return h.catalog.Add(ctx, alertID, alert.UserID, alert.AlertType.String(), alert.Message, alert.Payload, alert.IsRead)
}

func (h middlemanAlertHandlers[T]) onAlertRead(ctx context.Context, event ddd.Event) error {
	alert := event.Payload().(*domain.AlertRead)
	return h.catalog.Read(ctx, alert.ID, alert.IsRead)
}
