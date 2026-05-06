package handlers

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"middleman/baskets/basketspb"
	"middleman/comments/commentspb"
	"middleman/following/followingpb"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/internal/registry"
	"middleman/messages/messagespb"
	"middleman/notifications/internal/application"
	"middleman/notifications/internal/application/commands"
	"middleman/notifications/internal/domain"
	"middleman/offers/offerspb"
	"middleman/ordering/orderingpb"
	"middleman/payments/paymentspb"
	"middleman/products/productspb"
	"middleman/reviews/reviewspb"
	"middleman/support/supportpb"
	"middleman/users/userspb"
	"middleman/wishlists/wishlistspb"
	"time"
)

type integrationHandlers[T ddd.Event] struct {
	app      application.App
	catalog  domain.CatalogRepository
	users    domain.UserCacheRepository
	products domain.ProductCacheRepository
	logger   zerolog.Logger
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers(
	reg registry.Registry,
	app application.App,
	catalog domain.CatalogRepository,
	users domain.UserCacheRepository,
	products domain.ProductCacheRepository,
	logger zerolog.Logger,
	mws ...am.MessageHandlerMiddleware,
) am.MessageHandler {
	return am.NewEventHandler(reg, integrationHandlers[ddd.Event]{
		app:      app,
		catalog:  catalog,
		users:    users,
		products: products,
		logger:   logger,
	}, logger, mws...)
}

func RegisterIntegrationEventHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) (err error) {
	_, err = subscriber.Subscribe(userspb.UserAggregateChannel, handlers, am.MessageFilter{
		userspb.UserCreatedEvent,
	}, am.GroupName("notification-users"))
	if err != nil {
		return err
	}

	_, err = subscriber.Subscribe(productspb.ProductAggregateChannel, handlers, am.MessageFilter{
		productspb.ProductAddedEvent,
		productspb.ProductPriceDecreasedEvent,
		productspb.ProductPriceIncreasedEvent,
		productspb.ProductRebrandedEvent,
	}, am.GroupName("notification-products"))
	if err != nil {
		return err
	}

	_, err = subscriber.Subscribe(orderingpb.OrderAggregateChannel, handlers, am.MessageFilter{
		orderingpb.OrderCreatedEvent,
		orderingpb.OrderReadiedEvent,
		orderingpb.OrderCanceledEvent,
		orderingpb.OrderCompletedEvent,
	}, am.GroupName("notification-orders"))
	if err != nil {
		return err
	}

	_, err = subscriber.Subscribe(basketspb.BasketAggregateChannel, handlers, am.MessageFilter{
		basketspb.BasketStartedEvent,
		basketspb.BasketCanceledEvent,
		basketspb.BasketItemAddedEvent,
		basketspb.BasketItemRemovedEvent,
	}, am.GroupName("notification-baskets"))
	if err != nil {
		return err
	}

	_, err = subscriber.Subscribe(messagespb.MessageAggregateChannel, handlers, am.MessageFilter{
		messagespb.MessageSentEvent,
	}, am.GroupName("notification-messages"))
	if err != nil {
		return err
	}

	_, err = subscriber.Subscribe(commentspb.CommentAggregateChannel, handlers, am.MessageFilter{
		commentspb.CommentAddedEvent,
	}, am.GroupName("notification-comments"))
	if err != nil {
		return err
	}

	_, err = subscriber.Subscribe(reviewspb.ReviewAggregateChannel, handlers, am.MessageFilter{
		reviewspb.ReviewAddedEvent,
	}, am.GroupName("notification-reviews"))
	if err != nil {
		return err
	}

	_, err = subscriber.Subscribe(followingpb.FollowAggregateChannel, handlers, am.MessageFilter{
		followingpb.FollowAddedEvent,
	}, am.GroupName("notification-following"))
	if err != nil {
		return err
	}

	_, err = subscriber.Subscribe(offerspb.OfferAggregateChannel, handlers, am.MessageFilter{
		offerspb.OfferCreatedEvent,
	}, am.GroupName("notification-offers"))
	if err != nil {
		return err
	}

	_, err = subscriber.Subscribe(wishlistspb.WishlistItemAggregateChannel, handlers, am.MessageFilter{
		wishlistspb.WishlistItemAddedEvent,
	}, am.GroupName("notification-wishlists"))
	if err != nil {
		return err
	}

	_, err = subscriber.Subscribe(supportpb.TicketAggregateChannel, handlers, am.MessageFilter{
		supportpb.TicketCreatedEvent,
	}, am.GroupName("notification-support"))
	if err != nil {
		return err
	}

	_, err = subscriber.Subscribe(paymentspb.PaymentAggregateChannel, handlers, am.MessageFilter{
		paymentspb.PaymentConfirmedEvent,
	}, am.GroupName("notification-payments"))
	if err != nil {
		return err
	}

	return nil
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
	// products
	case productspb.ProductAddedEvent:
		return h.onProductAdded(ctx, event)
	case productspb.ProductPriceDecreasedEvent:
		return h.onProductPriceDecreased(ctx, event)
	case productspb.ProductPriceIncreasedEvent:
		return h.onProductPriceIncreased(ctx, event)

	// messages
	case messagespb.MessageSentEvent:
		return h.onMessageSent(ctx, event)

	// comments
	case commentspb.CommentAddedEvent:
		return h.onCommentAdded(ctx, event)

	// reviews
	case reviewspb.ReviewAddedEvent:
		return h.onReviewAdded(ctx, event)

	// following
	case followingpb.FollowAddedEvent:
		return h.onFollowAdded(ctx, event)

	// offers
	case offerspb.OfferCreatedEvent:
		return h.onOfferCreated(ctx, event)

	// wishlists
	case wishlistspb.WishlistItemAddedEvent:
		return h.onWishlistItemAdded(ctx, event)

	// support
	case supportpb.TicketCreatedEvent:
		return h.onTicketCreated(ctx, event)

	// ordering
	case orderingpb.OrderCreatedEvent:
		return h.onOrderCreated(ctx, event)
	case orderingpb.OrderReadiedEvent:
		return h.onOrderReadied(ctx, event)
	case orderingpb.OrderCanceledEvent:
		return h.onOrderCanceled(ctx, event)
	case orderingpb.OrderCompletedEvent:
		return h.onOrderCompleted(ctx, event)

	// payments
	case paymentspb.PaymentConfirmedEvent:
		return h.onPaymentConfirmed(ctx, event)

	// baskets
	case basketspb.BasketItemAddedEvent:
		return h.onBasketItemAdded(ctx, event)
	}
	return nil
}

func (h integrationHandlers[T]) onProductAdded(ctx context.Context, event T) error {
	payload := event.Payload().(*productspb.ProductAdded)
	alertId := uuid.New().String()

	h.logger.Info().
		Str("alert_id", alertId).
		Str("product_id", payload.GetId()).
		Str("user_id", payload.GetUserSellerId()).
		Msg("Creating product added alert")

	// Add product to cache for future lookups
	if err := h.products.Add(ctx, payload.GetId(), payload.GetUserSellerId(),
		payload.GetName(), payload.GetDescription(), payload.GetBasePrice()); err != nil {
		h.logger.Warn().Err(err).Msg("Failed to cache product")
	}

	return h.app.AddProductAlert(ctx, commands.AddProductAlert{
		ID:        alertId,
		ProductID: payload.GetId(),
		UserID:    payload.GetUserSellerId(),
		Message:   "Your product has been successfully listed",
	})
}

func (h integrationHandlers[T]) onProductPriceDecreased(ctx context.Context, event T) error {
	payload := event.Payload().(*productspb.ProductPriceDecreased)
	alertId := uuid.New().String()

	product, err := h.products.Find(ctx, payload.GetId())
	if err != nil {
		h.logger.Error().Err(err).Str("product_id", payload.GetId()).
			Msg("Failed to find product for price decrease notification")
		return err
	}

	h.logger.Info().
		Str("alert_id", alertId).
		Str("product_id", payload.GetId()).
		Float64("old_price", float64(payload.GetOldPrice())).
		Float64("new_price", float64(payload.GetNewPrice())).
		Msg("Creating product price decreased alert")

	return h.app.AddProductAlert(ctx, commands.AddProductAlert{
		ID:        alertId,
		ProductID: payload.GetId(),
		UserID:    product.UserSellerID,
		Message: fmt.Sprintf("Your product price has been decreased from $%.2f to $%.2f",
			payload.GetOldPrice(), payload.GetNewPrice()),
	})
}

func (h integrationHandlers[T]) onProductPriceIncreased(ctx context.Context, event T) error {
	payload := event.Payload().(*productspb.ProductPriceIncreased)
	alertId := uuid.New().String()

	product, err := h.products.Find(ctx, payload.GetId())
	if err != nil {
		h.logger.Error().Err(err).Str("product_id", payload.GetId()).
			Msg("Failed to find product for price increase notification")
		return err
	}

	h.logger.Info().
		Str("alert_id", alertId).
		Str("product_id", payload.GetId()).
		Float64("old_price", float64(payload.GetOldPrice())).
		Float64("new_price", float64(payload.GetNewPrice())).
		Msg("Creating product price increased alert")

	return h.app.AddProductAlert(ctx, commands.AddProductAlert{
		ID:        alertId,
		ProductID: payload.GetId(),
		UserID:    product.UserSellerID,
		Message: fmt.Sprintf("Your product price has been increased from $%.2f to $%.2f",
			payload.GetOldPrice(), payload.GetNewPrice()),
	})
}

func (h integrationHandlers[T]) onMessageSent(ctx context.Context, event T) error {
	payload := event.Payload().(*messagespb.MessageSent)
	alertId := uuid.New().String()

	h.logger.Info().
		Str("alert_id", alertId).
		Str("message_id", payload.GetId()).
		Str("recipient_id", payload.GetRecipientId()).
		Str("sender_id", payload.GetSenderId()).
		Msg("Creating message alert")

	return h.app.AddMessageAlert(ctx, commands.AddMessageAlert{
		ID:              alertId,
		UserID:          payload.GetRecipientId(),
		ProductID:       payload.GetItemId(),
		MessageID:       payload.GetId(),
		MessageSenderID: payload.GetSenderId(),
		Message:         "You have a new message",
	})
}

func (h integrationHandlers[T]) onCommentAdded(ctx context.Context, event T) error {
	payload := event.Payload().(*commentspb.CommentAdded)
	alertId := uuid.New().String()

	// Notify the product owner about the new comment
	product, err := h.products.Find(ctx, payload.GetItemId())
	if err != nil {
		h.logger.Error().Err(err).Str("product_id", payload.GetItemId()).
			Msg("Failed to find product for comment notification")
		return err
	}

	h.logger.Info().
		Str("alert_id", alertId).
		Str("comment_id", payload.GetId()).
		Str("product_id", payload.GetItemId()).
		Str("user_id", product.UserSellerID).
		Msg("Creating comment alert")

	return h.app.AddCommentAlert(ctx, commands.AddCommentAlert{
		ID:          alertId,
		UserID:      product.UserSellerID,
		CommentID:   payload.GetId(),
		UserAddedID: payload.GetSenderId(),
		ProductID:   payload.GetItemId(),
		Message:     "New comment on your product",
	})
}

func (h integrationHandlers[T]) onReviewAdded(ctx context.Context, event T) error {
	payload := event.Payload().(*reviewspb.ReviewAdded)
	alertId := uuid.New().String()

	// Notify the product owner about the new review
	product, err := h.products.Find(ctx, payload.GetItemId())
	if err != nil {
		h.logger.Error().Err(err).Str("product_id", payload.GetItemId()).
			Msg("Failed to find product for review notification")
		return err
	}

	h.logger.Info().
		Str("alert_id", alertId).
		Str("review_id", payload.GetId()).
		Str("product_id", payload.GetItemId()).
		Str("user_id", product.UserSellerID).
		Msg("Creating review alert")

	return h.app.AddReviewAlert(ctx, commands.AddReviewAlert{
		ID:        alertId,
		UserID:    product.UserSellerID,
		ReviewID:  payload.GetId(),
		ProductID: payload.GetItemId(),
		Message:   fmt.Sprintf("New %d-star review on your product", payload.GetContent()),
	})
}

func (h integrationHandlers[T]) onFollowAdded(ctx context.Context, event T) error {
	payload := event.Payload().(*followingpb.FollowAdded)
	alertId := uuid.New().String()

	h.logger.Info().
		Str("alert_id", alertId).
		Str("user_id", payload.GetUserId()).
		Str("follower_id", payload.GetFollowedUserId()).
		Msg("Creating following alert")

	return h.app.AddFollowingAlert(ctx, commands.AddFollowingAlert{
		ID:         alertId,
		UserID:     payload.GetUserId(),
		FollowerID: payload.GetFollowedUserId(),
		Message:    "You have a new follower",
	})
}

func (h integrationHandlers[T]) onOfferCreated(ctx context.Context, event T) error {
	payload := event.Payload().(*offerspb.OfferCreated)
	alertId := uuid.New().String()

	// Notify the product owner about the new offer
	product, err := h.products.Find(ctx, payload.GetProductId())
	if err != nil {
		h.logger.Error().Err(err).Str("product_id", payload.GetProductId()).
			Msg("Failed to find product for offer notification")
		return err
	}

	h.logger.Info().
		Str("alert_id", alertId).
		Str("offer_id", payload.GetId()).
		Str("product_id", payload.GetProductId()).
		Str("user_id", product.UserSellerID).
		Msg("Creating offer alert")

	return h.app.AddOfferAlert(ctx, commands.AddOfferAlert{
		ID:        alertId,
		UserID:    product.UserSellerID,
		OfferID:   payload.GetId(),
		ProductID: payload.GetProductId(),
		Message:   "New offer on your product",
	})
}

func (h integrationHandlers[T]) onWishlistItemAdded(ctx context.Context, event T) error {
	payload := event.Payload().(*wishlistspb.WishlistItemAdded)
	alertId := uuid.New().String()

	// Notify the product owner that their product was added to a wishlist
	product, err := h.products.Find(ctx, payload.GetItemId())
	if err != nil {
		h.logger.Error().Err(err).Str("product_id", payload.GetItemId()).
			Msg("Failed to find product for wishlist notification")
		return err
	}

	h.logger.Info().
		Str("alert_id", alertId).
		Str("wishlist_id", payload.GetId()).
		Str("product_id", payload.GetItemId()).
		Str("user_id", product.UserSellerID).
		Msg("Creating wishlist alert")

	return h.app.AddWishlistAlert(ctx, commands.AddWishlistAlert{
		ID:         alertId,
		UserID:     product.UserSellerID,
		WishlistID: payload.GetId(),
		ProductID:  payload.GetItemId(),
		Message:    "Your product was added to a wishlist",
	})
}

func (h integrationHandlers[T]) onTicketCreated(ctx context.Context, event T) error {
	payload := event.Payload().(*supportpb.TicketCreated)
	alertId := uuid.New().String()

	h.logger.Info().
		Str("alert_id", alertId).
		Str("ticket_id", payload.GetId()).
		Str("user_id", payload.GetCreatedBy()).
		Msg("Creating support ticket alert")

	return h.app.AddSupportAlert(ctx, commands.AddSupportAlert{
		ID:       alertId,
		UserID:   payload.GetCreatedBy(),
		TicketID: payload.GetId(),
		Message:  "Your support ticket has been created",
	})
}

func (h integrationHandlers[T]) onOrderCreated(ctx context.Context, event T) error {
	payload := event.Payload().(*orderingpb.OrderCreated)
	alertId := uuid.New().String()

	h.logger.Info().
		Str("alert_id", alertId).
		Str("order_id", payload.GetId()).
		Str("customer_id", payload.GetUserCustomerId()).
		Msg("Creating order created alert")

	// Get first product ID from order items for notification
	productID := ""
	if items := payload.GetItems(); len(items) > 0 {
		productID = items[0].GetProductId()
	}

	return h.app.AddOrderAlert(ctx, commands.AddOrderAlert{
		ID:             alertId,
		OrderID:        payload.GetId(),
		UserID:         payload.GetUserCustomerId(),
		ProductID:      productID,
		UserCustomerID: payload.GetUserCustomerId(),
		Message:        "Your order has been created",
	})
}

func (h integrationHandlers[T]) onOrderReadied(ctx context.Context, event T) error {
	payload := event.Payload().(*orderingpb.OrderReadied)
	alertId := uuid.New().String()

	h.logger.Info().
		Str("alert_id", alertId).
		Str("order_id", payload.GetId()).
		Msg("Creating order readied alert")

	// For readied/canceled/completed events, we need order repository to get customer info
	// For now, we'll skip the notification if we can't determine the customer
	h.logger.Warn().Msg("Order readied event missing customer information - notification skipped")
	return nil
}

func (h integrationHandlers[T]) onOrderCanceled(ctx context.Context, event T) error {
	payload := event.Payload().(*orderingpb.OrderCanceled)
	alertId := uuid.New().String()

	h.logger.Info().
		Str("alert_id", alertId).
		Str("order_id", payload.GetId()).
		Msg("Creating order canceled alert")

	// For readied/canceled/completed events, we need order repository to get customer info
	// For now, we'll skip the notification if we can't determine the customer
	h.logger.Warn().Msg("Order canceled event missing customer information - notification skipped")
	return nil
}

func (h integrationHandlers[T]) onOrderCompleted(ctx context.Context, event T) error {
	payload := event.Payload().(*orderingpb.OrderCompleted)
	alertId := uuid.New().String()

	h.logger.Info().
		Str("alert_id", alertId).
		Str("order_id", payload.GetId()).
		Msg("Creating order completed alert")

	// For readied/canceled/completed events, we need order repository to get customer info
	// For now, we'll skip the notification if we can't determine the customer
	h.logger.Warn().Msg("Order completed event missing customer information - notification skipped")
	return nil
}

func (h integrationHandlers[T]) onPaymentConfirmed(ctx context.Context, event T) error {
	payload := event.Payload().(*paymentspb.PaymentConfirmed)
	alertId := uuid.New().String()

	h.logger.Info().
		Str("alert_id", alertId).
		Str("payment_id", payload.GetPaymentId()).
		Str("order_id", payload.GetOrderId()).
		Msg("Creating payment confirmed alert")

	// Payment events should include user ID to send notifications
	// For now, we create the alert with order context
	return h.app.AddPaymentAlert(ctx, commands.AddPaymentAlert{
		ID:        alertId,
		UserID:    "", // Will need to be resolved from order
		PaymentID: payload.GetPaymentId(),
		OrderID:   payload.GetOrderId(),
		Message:   "Payment confirmed for your order",
	})
}

func (h integrationHandlers[T]) onBasketItemAdded(ctx context.Context, event T) error {
	payload := event.Payload().(*basketspb.BasketItemAdded)
	alertId := uuid.New().String()

	// Notify the product owner that their product was added to a basket
	product, err := h.products.Find(ctx, payload.GetProductId())
	if err != nil {
		h.logger.Error().Err(err).Str("product_id", payload.GetProductId()).
			Msg("Failed to find product for basket notification")
		return err
	}

	// Extract basket ID from aggregate event
	var basketID string

	h.logger.Info().
		Str("alert_id", alertId).
		Str("basket_id", basketID).
		Str("product_id", payload.GetProductId()).
		Str("user_id", product.UserSellerID).
		Msg("Creating basket item added alert")

	return h.app.AddBasketAlert(ctx, commands.AddBasketAlert{
		ID:             alertId,
		BasketID:       basketID,
		ProductID:      payload.GetProductId(),
		UserCustomerID: product.UserSellerID,
		Message:        "Your product was added to a basket",
	})
}
