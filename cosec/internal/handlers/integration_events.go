package handlers

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"middleman/baskets/basketspb"
	"middleman/cosec/internal/models"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/registry"
	"middleman/internal/sec"
	"middleman/payments/paymentspb"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type integrationHandlers[T ddd.Event] struct {
	orchestrator sec.Orchestrator[*models.CheckoutData]
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers(reg registry.Registry, orchestrator sec.Orchestrator[*models.CheckoutData], logger zerolog.Logger, mws ...am.MessageHandlerMiddleware) am.MessageHandler {
	return am.NewEventHandler(reg, integrationHandlers[ddd.Event]{
		orchestrator: orchestrator,
	}, logger, mws...)
}

func RegisterIntegrationEventHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) (err error) {
	// Baskets - PRIMARY saga trigger
	if _, err = subscriber.Subscribe(basketspb.BasketAggregateChannel, handlers, am.MessageFilter{
		basketspb.BasketCheckedOutEvent,
	}, am.GroupName("cosec-baskets")); err != nil {
		return err
	}

	// Payments - for monitoring only (no saga trigger)
	if _, err = subscriber.Subscribe(paymentspb.PaymentAggregateChannel, handlers, am.MessageFilter{
		paymentspb.PaymentAuthorizedEvent,
	}, am.GroupName("cosec-payments")); err != nil {
		return err
	}
	return nil
}

func (h integrationHandlers[T]) HandleEvent(ctx context.Context, event ddd.Event) error {
	span := trace.SpanFromContext(ctx)
	logger := log.With().
		Str("service", "cosec").
		Str("handler", "integration").
		Str("event_name", event.EventName()).
		Str("event_id", event.ID()).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	startTime := time.Now()
	logger.Info().Msg("COSEC_INTEGRATION_EVENT_BEGIN: Processing integration event")

	switch event.EventName() {
	case basketspb.BasketCheckedOutEvent:
		return h.onBasketCheckedOut(ctx, event)
	case paymentspb.PaymentAuthorizedEvent:
		return h.onPaymentAuthorized(ctx, event)
	}

	logger.Debug().
		Dur("duration_ms", time.Since(startTime)).
		Msg("COSEC_INTEGRATION_EVENT_UNKNOWN: Unknown event type, skipping")

	return nil
}

func (h integrationHandlers[T]) onBasketCheckedOut(ctx context.Context, event ddd.Event) error {
	span := trace.SpanFromContext(ctx)
	logger := log.With().
		Str("service", "cosec").
		Str("handler", "integration").
		Str("operation", "onBasketCheckedOut").
		Str("event_id", event.ID()).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	startTime := time.Now()
	logger.Info().Msg("COSEC_BASKET_CHECKED_OUT_BEGIN: Processing BasketCheckedOut event - SAGA TRIGGER")

	payload := event.Payload().(*basketspb.BasketCheckedOut)

	logger.Info().
		Str("basket_id", payload.GetId()).
		Str("user_customer_id", payload.GetUserCustomerId()).
		Int("items_count", len(payload.GetItems())).
		Msg("COSEC_BASKET_CHECKED_OUT_PAYLOAD: Event payload details")

	// Log each item for detailed tracking
	for i, item := range payload.GetItems() {
		logger.Debug().
			Int("item_index", i).
			Str("product_id", item.GetProductId()).
			Str("seller_id", item.GetUserSellerId()).
			Int64("quantity", item.GetQuantity()).
			Int64("price", item.GetPrice()).
			Msg("COSEC_BASKET_CHECKED_OUT_ITEM: Item details")
	}

	// Convert to CheckoutData for saga
	checkoutData := &models.CheckoutData{
		BasketID:  payload.GetId(),
		UserID:    payload.GetUserCustomerId(),
		PaymentID: payload.GetPaymentIntentId(), // Use existing payment intent ID
		Items:     make([]models.CheckoutItem, len(payload.GetItems())),
		Total:     payload.GetTotal(), // Use the total from the event
	}

	// Log the total from event vs calculated
	var calculatedTotal int64
	for i, item := range payload.GetItems() {
		checkoutData.Items[i] = models.CheckoutItem{
			ProductID:   item.GetProductId(),
			ProductName: item.GetProductName(),
			SellerID:    item.GetUserSellerId(),
			SellerName:  item.GetUserSellerName(),
			Quantity:    item.GetQuantity(),
			Price:       item.GetPrice(),
		}
		calculatedTotal += item.GetPrice() * item.GetQuantity()
	}

	// Log any discrepancy
	if calculatedTotal != payload.GetTotal() {
		logger.Warn().
			Int64("event_total", payload.GetTotal()).
			Int64("calculated_total", calculatedTotal).
			Msg("COSEC_BASKET_CHECKED_OUT_TOTAL_MISMATCH: Total from event doesn't match calculated total")
	}

	logger.Info().
		Str("basket_id", checkoutData.BasketID).
		Str("user_id", checkoutData.UserID).
		Int("items_count", len(checkoutData.Items)).
		Int64("total_amount", checkoutData.Total).
		Msg("COSEC_BASKET_CHECKED_OUT_DATA_CREATED: CheckoutData created for saga")

	// Start the checkout saga
	sagaID := fmt.Sprintf("checkout-%s", payload.GetId())

	logger.Info().
		Str("saga_id", sagaID).
		Str("saga_name", "cosec.Checkout").
		Msg("COSEC_BASKET_CHECKED_OUT_SAGA_START: Starting checkout saga")

	if err := h.orchestrator.Start(ctx, sagaID, checkoutData); err != nil {
		logger.Error().Err(err).
			Str("saga_id", sagaID).
			Dur("duration_ms", time.Since(startTime)).
			Msg("COSEC_BASKET_CHECKED_OUT_SAGA_FAILED: Failed to start checkout saga")
		return err
	}

	logger.Info().
		Str("saga_id", sagaID).
		Str("basket_id", checkoutData.BasketID).
		Dur("duration_ms", time.Since(startTime)).
		Msg("COSEC_BASKET_CHECKED_OUT_SUCCESS: Checkout saga started successfully")

	span.AddEvent("saga_started", trace.WithAttributes(
		attribute.String("saga_id", sagaID),
		attribute.String("basket_id", checkoutData.BasketID),
		attribute.String("user_id", checkoutData.UserID),
		attribute.Int("items_count", len(checkoutData.Items)),
		attribute.Int64("total_amount", checkoutData.Total),
	))

	return nil
}

func (h integrationHandlers[T]) onPaymentAuthorized(ctx context.Context, event ddd.Event) error {
	// Payment authorized event - for monitoring/logging only
	// Saga should already be running from BasketCheckedOut
	span := trace.SpanFromContext(ctx)
	span.AddEvent("Payment authorized - saga should be running", trace.WithAttributes(
		attribute.String("PaymentID", event.Payload().(*paymentspb.PaymentAuthorized).GetPaymentIntentId()),
	))

	// NO saga start here - prevents duplicate saga corruption
	return nil
}
