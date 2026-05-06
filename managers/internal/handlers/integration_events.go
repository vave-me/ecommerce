package handlers

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"middleman/activity/activitypb"
	"middleman/baskets/basketspb"
	"middleman/comments/commentspb"
	"middleman/following/followingpb"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/internal/registry"
	"middleman/managers/internal/application"
	"middleman/managers/internal/application/consciousness"
	"middleman/messages/messagespb"
	"middleman/offers/offerspb"
	"middleman/ordering/orderingpb"
	"middleman/payments/paymentspb"
	"middleman/posts/postspb"
	"middleman/products/productspb"
	"middleman/reviews/reviewspb"
	"middleman/services/servicespb"
	"middleman/support/supportpb"
	"middleman/users/userspb"
	"middleman/wishlists/wishlistspb"
)

type integrationHandlers[T ddd.Event] struct {
	app                  application.App
	consciousnessManager interface{ ProcessEvent(context.Context, ddd.Event) error }
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers(reg registry.Registry, app application.App, mws ...am.MessageHandlerMiddleware) am.MessageHandler {
	return am.NewEventHandler(reg, integrationHandlers[ddd.Event]{
		app: app,
	}, zerolog.Logger{}, mws...)
}

// NewIntegrationEventHandlersWithConsciousness creates handlers with consciousness manager
func NewIntegrationEventHandlersWithConsciousness(reg registry.Registry, app application.App, consciousnessManager interface{ ProcessEvent(context.Context, ddd.Event) error }, mws ...am.MessageHandlerMiddleware) am.MessageHandler {
	return am.NewEventHandler(reg, integrationHandlers[ddd.Event]{
		app:                  app,
		consciousnessManager: consciousnessManager,
	}, zerolog.Logger{}, mws...)
}

func RegisterIntegrationEventHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) (err error) {
	// User events
	if _, err = subscriber.Subscribe(userspb.UserAggregateChannel, handlers, am.MessageFilter{
		userspb.UserCreatedEvent,
		userspb.UserRenamedEvent,
		userspb.UserEnabledToggledEvent,
		userspb.UserUpdatedEvent,
		userspb.UserLoggedInEvent,
		userspb.UserLoggedOutEvent,
	}, am.GroupName("managers-users")); err != nil {
		return err
	}

	// Product events
	if _, err = subscriber.Subscribe(productspb.ProductAggregateChannel, handlers, am.MessageFilter{
		productspb.ProductAddedEvent,
		productspb.ProductUpdatedEvent,
		productspb.ProductRemovedEvent,
		productspb.ProductPriceIncreasedEvent,
		productspb.ProductPriceDecreasedEvent,
		productspb.ProductSoldEvent,
		productspb.ProductStockAdjustedEvent,
		productspb.ProductArchivedEvent,
	}, am.GroupName("managers-products")); err != nil {
		return err
	}

	// Order events
	if _, err = subscriber.Subscribe(orderingpb.OrderAggregateChannel, handlers, am.MessageFilter{
		orderingpb.OrderCreatedEvent,
		orderingpb.OrderReadiedEvent,
		orderingpb.OrderCanceledEvent,
		orderingpb.OrderCompletedEvent,
	}, am.GroupName("managers-orders")); err != nil {
		return err
	}

	// Payment events
	if _, err = subscriber.Subscribe(paymentspb.PaymentAggregateChannel, handlers, am.MessageFilter{
		paymentspb.PaymentAuthorizedEvent,
		paymentspb.PaymentConfirmedEvent,
	}, am.GroupName("managers-payments")); err != nil {
		return err
	}

	// Basket events
	if _, err = subscriber.Subscribe(basketspb.BasketAggregateChannel, handlers, am.MessageFilter{
		basketspb.BasketStartedEvent,
		basketspb.BasketItemAddedEvent,
		basketspb.BasketItemRemovedEvent,
		basketspb.BasketCanceledEvent,
		basketspb.BasketCheckedOutEvent,
	}, am.GroupName("managers-baskets")); err != nil {
		return err
	}

	// Review events
	if _, err = subscriber.Subscribe(reviewspb.ReviewAggregateChannel, handlers, am.MessageFilter{
		reviewspb.ReviewAddedEvent,
	}, am.GroupName("managers-reviews")); err != nil {
		return err
	}

	// Support events
	if _, err = subscriber.Subscribe(supportpb.TicketAggregateChannel, handlers, am.MessageFilter{
		supportpb.TicketCreatedEvent,
		supportpb.TicketUpdatedEvent,
		supportpb.TicketEscalatedEvent,
		supportpb.TicketClosedEvent,
		supportpb.TicketReplyAddedEvent,
	}, am.GroupName("managers-support")); err != nil {
		return err
	}

	// Wishlist events
	if _, err = subscriber.Subscribe(wishlistspb.WishlistAggregateChannel, handlers, am.MessageFilter{
		wishlistspb.WishlistCreatedEvent,
		wishlistspb.WishlistRemovedEvent,
	}, am.GroupName("managers-wishlists")); err != nil {
		return err
	}

	// Wishlist item events
	if _, err = subscriber.Subscribe(wishlistspb.WishlistItemAggregateChannel, handlers, am.MessageFilter{
		wishlistspb.WishlistItemAddedEvent,
		wishlistspb.WishlistItemRemovedEvent,
	}, am.GroupName("managers-wishlist-items")); err != nil {
		return err
	}

	// Comment events
	if _, err = subscriber.Subscribe(commentspb.CommentAggregateChannel, handlers, am.MessageFilter{
		commentspb.CommentAddedEvent,
	}, am.GroupName("managers-comments")); err != nil {
		return err
	}

	// Activity events
	if _, err = subscriber.Subscribe(activitypb.ActivityAggregateChannel, handlers, am.MessageFilter{
		activitypb.ActivityCreatedEvent,
	}, am.GroupName("managers-activity")); err != nil {
		return err
	}

	// Interaction events
	if _, err = subscriber.Subscribe(activitypb.InteractionAggregateChannel, handlers, am.MessageFilter{
		activitypb.InteractionAddedEvent,
		activitypb.InteractionUpdatedEvent,
		activitypb.InteractionRemovedEvent,
	}, am.GroupName("managers-interactions")); err != nil {
		return err
	}

	// Following events
	if _, err = subscriber.Subscribe(followingpb.FollowAggregateChannel, handlers, am.MessageFilter{
		followingpb.FollowAddedEvent,
	}, am.GroupName("managers-following")); err != nil {
		return err
	}

	// Messages events
	if _, err = subscriber.Subscribe(messagespb.MessageAggregateChannel, handlers, am.MessageFilter{
		messagespb.MessageSentEvent,
		messagespb.MessageDeletedEvent,
	}, am.GroupName("managers-messages")); err != nil {
		return err
	}

	// Conversation events
	if _, err = subscriber.Subscribe(messagespb.ConversationAggregateChannel, handlers, am.MessageFilter{
		messagespb.ConversationStartedEvent,
		messagespb.ConversationActiveToggledEvent,
	}, am.GroupName("managers-conversations")); err != nil {
		return err
	}

	// Offers events
	if _, err = subscriber.Subscribe(offerspb.OfferAggregateChannel, handlers, am.MessageFilter{
		offerspb.OfferCreatedEvent,
		offerspb.OfferAcceptedEvent,
		offerspb.OfferActivatedEvent,
		offerspb.OfferClosedEvent,
	}, am.GroupName("managers-offers")); err != nil {
		return err
	}

	// Posts events
	if _, err = subscriber.Subscribe(postspb.PostAggregateChannel, handlers, am.MessageFilter{
		postspb.PostAddedEvent,
		postspb.PostUpdatedEvent,
		postspb.PostRemovedEvent,
		postspb.PostArchivedEvent,
		postspb.PostThumbnailAddedEvent,
		postspb.PostThumbnailUpdatedEvent,
	}, am.GroupName("managers-posts")); err != nil {
		return err
	}

	// Services events
	if _, err = subscriber.Subscribe(servicespb.ServiceAggregateChannel, handlers, am.MessageFilter{
		servicespb.ServiceAddedEvent,
		servicespb.ServiceUpdatedEvent,
		servicespb.ServiceRemovedEvent,
	}, am.GroupName("managers-services")); err != nil {
		return err
	}

	return nil
}
func (h integrationHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	// Add logging to see what events are being received
	zerolog.Ctx(ctx).Info().
		Str("event_type", event.EventName()).
		Str("event_id", event.ID()).
		Msg("Integration handler received event")
	
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

	// Route to consciousness manager if available
	if h.consciousnessManager != nil {
		zerolog.Ctx(ctx).Info().Msg("Routing event to consciousness manager")
		err = h.consciousnessManager.ProcessEvent(ctx, event)
		if err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to process event in consciousness manager")
		} else {
			zerolog.Ctx(ctx).Info().Msg("Successfully processed event in consciousness manager")
		}
		return err
	}
	
	// Fallback to standard processing
	zerolog.Ctx(ctx).Info().Msg("Calling app.ProcessPlatformEvent")
	err = h.app.ProcessPlatformEvent(ctx, event)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to process platform event")
	} else {
		zerolog.Ctx(ctx).Info().Msg("Successfully processed platform event")
	}
	return err
}

