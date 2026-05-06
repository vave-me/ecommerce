package handlers

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"middleman/activity/activitypb"
	"middleman/baskets/basketspb"
	"middleman/comments/commentspb"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/internal/registry"
	"middleman/messages/messagespb"
	"middleman/metrics/internal/application"
	"middleman/metrics/internal/models"
	"middleman/posts/postspb"
	"middleman/products/productspb"
	"middleman/reviews/reviewspb"
	"middleman/services/servicespb"
	"middleman/users/userspb"
	"middleman/wishlists/wishlistspb"
	"time"
)

type integrationHandlers[T ddd.Event] struct {
	itemMetrics application.ItemMetricRepository
	userMetrics application.UserMetricRepository
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers(reg registry.Registry, itemMetrics application.ItemMetricRepository, userMetrics application.UserMetricRepository,

	mws ...am.MessageHandlerMiddleware) am.MessageHandler {
	return am.NewEventHandler(reg, integrationHandlers[ddd.Event]{
		itemMetrics: itemMetrics,
		userMetrics: userMetrics,
	}, zerolog.Logger{}, mws...)
}
func RegisterIntegrationEventHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) (err error) {

	if _, err = subscriber.Subscribe(activitypb.InteractionAggregateChannel, handlers, am.MessageFilter{
		activitypb.InteractionAddedEvent,
		activitypb.InteractionRemovedEvent,
	}, am.GroupName("metrics-activity")); err != nil {
		return
	}

	if _, err = subscriber.Subscribe(wishlistspb.WishlistAggregateChannel, handlers, am.MessageFilter{
		wishlistspb.WishlistItemAddedEvent,
		wishlistspb.WishlistItemRemovedEvent,
	}, am.GroupName("metrics-wishlists")); err != nil {
		return
	}

	if _, err = subscriber.Subscribe(commentspb.CommentAggregateChannel, handlers, am.MessageFilter{
		commentspb.CommentAddedEvent,
	}, am.GroupName("metrics-comments")); err != nil {
		return
	}

	if _, err = subscriber.Subscribe(productspb.ProductAggregateChannel, handlers, am.MessageFilter{
		productspb.ProductAddedEvent,
		productspb.ProductUpdatedEvent,
		productspb.ProductRemovedEvent,
	}, am.GroupName("metrics-product")); err != nil {
		return
	}

	if _, err = subscriber.Subscribe(postspb.PostAggregateChannel, handlers, am.MessageFilter{
		postspb.PostAddedEvent,
		postspb.PostUpdatedEvent,
		postspb.PostRemovedEvent,
	}, am.GroupName("metrics-posts")); err != nil {
		return
	}

	if _, err = subscriber.Subscribe(messagespb.MessageAggregateChannel, handlers, am.MessageFilter{
		messagespb.MessageSentEvent,
	}, am.GroupName("metrics-vehicles")); err != nil {
		return
	}

	if _, err = subscriber.Subscribe(basketspb.BasketAggregateChannel, handlers, am.MessageFilter{
		basketspb.BasketItemAddedEvent,
		basketspb.BasketItemRemovedEvent,
	}, am.GroupName("metrics-basket")); err != nil {
		return
	}
	if _, err = subscriber.Subscribe(servicespb.ServiceAggregateChannel, handlers, am.MessageFilter{
		servicespb.ServiceAddedEvent,
		servicespb.ServiceUpdatedEvent,
		servicespb.ServiceRemovedEvent,
	}, am.GroupName("metrics-services")); err != nil {
		return
	}
	if _, err = subscriber.Subscribe(userspb.UserAggregateChannel, handlers, am.MessageFilter{
		userspb.UserCreatedEvent,
		userspb.UserRenamedEvent,
	}, am.GroupName("metrics-users")); err != nil {
		return
	}

	return
}
func (h integrationHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling integration event",
				trace.WithAttributes(errorsotel.ErrAttrs(err)...),
			)
		}
		span.AddEvent("Handled integration event", trace.WithAttributes(
			attribute.Int64("TookMS", time.Since(started).Milliseconds()),
		))
	}(time.Now())

	span.AddEvent("Handling integration event", trace.WithAttributes(
		attribute.String("Event", event.EventName()),
	))

	switch event.EventName() {
	case userspb.UserCreatedEvent:
		return h.onUserCreated(ctx, event)
	case productspb.ProductAddedEvent:
		return h.onProductAdded(ctx, event)
	case productspb.ProductRemovedEvent:
		return h.onProductRemoved(ctx, event)
	case postspb.PostAddedEvent:
		return h.onPostAdded(ctx, event)
	case postspb.PostRemovedEvent:
		return h.onPostRemoved(ctx, event)
	case servicespb.ServiceAddedEvent:
		return h.onServiceAdded(ctx, event)
	case servicespb.ServiceRemovedEvent:
		return h.onServiceRemoved(ctx, event)
	case messagespb.MessageSentEvent:
		return h.onMessageSent(ctx, event)
	case wishlistspb.WishlistItemAddedEvent:
		return h.onWishlistItemAdded(ctx, event)
	case wishlistspb.WishlistItemRemovedEvent:
		return h.onWishlistItemRemoved(ctx, event)
	case commentspb.CommentAddedEvent:
		return h.onCommentAdded(ctx, event)
	case activitypb.InteractionAddedEvent:
		return h.onInteractionAdded(ctx, event)
	case activitypb.InteractionRemovedEvent:
		return h.onInteractionRemoved(ctx, event)
	}

	return nil
}

// TODO on all onUserCreated methods need to be rethinked way of caching what need to be cached
func (h integrationHandlers[T]) onUserCreated(ctx context.Context, event T) error {
	payload := event.Payload().(*userspb.UserCreated)
	return h.userMetrics.AddUserMetric(ctx, payload.GetId(), "private")
}

func (h integrationHandlers[T]) onProductAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*productspb.ProductAdded)

	// Update item metrics only - we don't know the correct user ID field
	return h.itemMetrics.AddMetric(ctx, payload.GetId(), models.ProductType.String(), payload.GetCategoryId(), payload.GetBasePrice(), float64(payload.GetLat()), float64(payload.GetLng()))
}

func (h integrationHandlers[T]) onProductRemoved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*productspb.ProductRemoved)
	return h.itemMetrics.RemoveItemMetric(ctx, payload.GetId())
}

func (h integrationHandlers[T]) onPostAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*postspb.PostAdded)

	// Update item metrics only - we don't know the correct user ID field
	return h.itemMetrics.AddMetric(ctx, payload.GetId(), models.PostType.String(), payload.GetCategoryId(), 0, float64(payload.GetLat()), float64(payload.GetLng()))
}

func (h integrationHandlers[T]) onPostRemoved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*postspb.PostRemoved)
	return h.itemMetrics.RemoveItemMetric(ctx, payload.GetId())
}

func (h integrationHandlers[T]) onMessageSent(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*messagespb.MessageSent)
	return h.itemMetrics.UpdateItemMetric(ctx, payload.GetId(), models.MetricTypeCountMessage.String(), models.MetricTypeActionAdd.String())
}
func (h integrationHandlers[T]) onServiceAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*servicespb.ServiceAdded)

	// Update item metrics only - we don't know the correct user ID field
	return h.itemMetrics.AddMetric(ctx, payload.GetId(), models.ServiceType.String(), payload.GetCategoryId(), payload.GetBasePrice(), float64(payload.GetLat()), float64(payload.GetLng()))
}

func (h integrationHandlers[T]) onServiceRemoved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*servicespb.ServiceRemoved)
	return h.itemMetrics.RemoveItemMetric(ctx, payload.GetId())
}
func (h integrationHandlers[T]) onWishlistItemAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*wishlistspb.WishlistItemAdded)
	return h.itemMetrics.UpdateItemMetric(ctx, payload.GetId(), models.MetricTypeCountWishlist.String(), models.MetricTypeActionAdd.String())
}
func (h integrationHandlers[T]) onWishlistItemRemoved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*wishlistspb.WishlistItemRemoved)
	return h.itemMetrics.UpdateItemMetric(ctx, payload.GetId(), models.MetricTypeCountWishlist.String(), models.MetricTypeActionRemove.String())
}
func (h integrationHandlers[T]) onCommentAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*commentspb.CommentAdded)

	// Update item metric only - we don't know the correct user ID field
	return h.itemMetrics.UpdateItemMetric(ctx, payload.GetItemId(), models.MetricTypeCountComment.String(), models.MetricTypeActionAdd.String())
}
func (h integrationHandlers[T]) onBasketItemAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*basketspb.BasketItemAdded)
	return h.itemMetrics.UpdateItemMetric(ctx, payload.GetProductId(), models.MetricTypeCountBasket.String(), models.MetricTypeActionAdd.String())
}
func (h integrationHandlers[T]) onBasketItemRemoved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*basketspb.BasketItemRemoved)
	return h.itemMetrics.UpdateItemMetric(ctx, payload.GetProductId(), models.MetricTypeCountBasket.String(), models.MetricTypeActionRemove.String())
}
func (h integrationHandlers[T]) onInteractionAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*activitypb.InteractionAdded)
	fmt.Println(payload.GetId())
	fmt.Println(payload.GetItemId())
	fmt.Println(payload.GetItemType())
	fmt.Println(payload.GetActionType())
	return h.itemMetrics.UpdateItemMetric(ctx, payload.GetItemId(), payload.GetActionType(), models.MetricTypeActionAdd.String())
}
func (h integrationHandlers[T]) onInteractionRemoved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*activitypb.InteractionRemoved)
	fmt.Println(payload.GetId())
	fmt.Println(payload.GetItemId())
	fmt.Println(payload.GetItemType())
	fmt.Println(payload.GetActionType())
	return h.itemMetrics.UpdateItemMetric(ctx, payload.GetItemId(), payload.GetActionType(), models.MetricTypeActionRemove.String())
}
func (h integrationHandlers[T]) onReviewAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*reviewspb.ReviewAdded)
	return h.itemMetrics.UpdateItemMetric(ctx, payload.GetItemId(), models.MetricTypeCountReview.String(), models.MetricTypeActionAdd.String())
}
func (h integrationHandlers[T]) onReviewRemoved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*reviewspb.ReviewRemoved)
	return h.itemMetrics.UpdateItemMetric(ctx, payload.GetId(), models.MetricTypeCountReview.String(), models.MetricTypeActionRemove.String())
}
