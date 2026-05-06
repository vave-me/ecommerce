package commands

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"middleman/internal/ddd"
	"middleman/ordering/internal/domain"

	"github.com/stackus/errors"
)

type CreateOrder struct {
	ID             string
	UserCustomerID string
	BasketID       string
	Items          []domain.Item
	PaymentIntent  string
}

type CreateOrderHandler struct {
	orders    domain.OrderRepository
	publisher ddd.EventPublisher[ddd.Event]
	logger    zerolog.Logger
}

func NewCreateOrderHandler(orders domain.OrderRepository, publisher ddd.EventPublisher[ddd.Event]) CreateOrderHandler {
	return CreateOrderHandler{
		orders:    orders,
		publisher: publisher,
		logger:    log.With().Str("service", "ordering").Str("handler", "CreateOrder").Logger(),
	}
}

func (h CreateOrderHandler) CreateOrder(ctx context.Context, cmd CreateOrder) error {
	span := trace.SpanFromContext(ctx)
	logger := h.logger.With().
		Str("operation", "CreateOrder").
		Str("order_id", cmd.ID).
		Str("user_customer_id", cmd.UserCustomerID).
		Str("basket_id", cmd.BasketID).
		Str("payment_intent", cmd.PaymentIntent).
		Int("items_count", len(cmd.Items)).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	startTime := time.Now()
	logger.Info().Msg("ORDERING_CREATE_ORDER_BEGIN: Starting order creation - SAGA STEP 1")

	// Log command details for debugging
	logger.Debug().
		Interface("items", cmd.Items).
		Msg("ORDERING_CREATE_ORDER_COMMAND_DETAILS: Full command details")

	// Log each item for detailed tracking
	for i, item := range cmd.Items {
		logger.Debug().
			Int("item_index", i).
			Str("product_id", item.ProductID).
			Str("seller_id", item.UserSellerID).
			Str("seller_name", item.UserSellerName).
			Str("product_name", item.ProductName).
			Int64("quantity", item.Quantity).
			Int64("price", item.Price).
			Msg("ORDERING_CREATE_ORDER_ITEM: Item details")
	}

	// Load or create the order aggregator
	order, err := h.orders.Load(ctx, cmd.ID)
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("ORDERING_CREATE_ORDER_LOAD_FAILED: Failed to load order aggregator")
		return err
	}

	logger.Debug().
		Str("current_status", string(order.Status)).
		Int("existing_items", len(order.Items)).
		Msg("ORDERING_CREATE_ORDER_LOADED: Order aggregator loaded successfully")

	// Attempt to create the order in the aggregator
	event, err := order.CreateOrder(cmd.ID, cmd.UserCustomerID, cmd.BasketID, cmd.PaymentIntent, cmd.Items)
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("ORDERING_CREATE_ORDER_DOMAIN_FAILED: Domain create order operation failed")
		return errors.Wrap(err, "create order command")
	}

	logger.Info().
		Str("event_name", event.EventName()).
		Str("new_status", string(order.Status)).
		Msg("ORDERING_CREATE_ORDER_EVENT_CREATED: Order created, event generated")

	// Save the aggregator state
	if err = h.orders.Save(ctx, order); err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("ORDERING_CREATE_ORDER_SAVE_FAILED: Failed to save order aggregator")
		return errors.Wrap(err, "order creation")
	}

	logger.Debug().
		Str("final_status", string(order.Status)).
		Int("final_items_count", len(order.Items)).
		Msg("ORDERING_CREATE_ORDER_SAVED: Order aggregator saved successfully")

	// Publish the domain event
	if pubErr := h.publisher.Publish(ctx, event); pubErr != nil {
		logger.Error().Err(pubErr).
			Dur("duration_ms", time.Since(startTime)).
			Str("event_name", event.EventName()).
			Msg("ORDERING_CREATE_ORDER_PUBLISH_FAILED: Failed to publish OrderCreated event")
		return pubErr
	}

	logger.Info().
		Dur("duration_ms", time.Since(startTime)).
		Str("event_id", event.ID()).
		Str("event_name", event.EventName()).
		Msg("ORDERING_CREATE_ORDER_SUCCESS: Order created and OrderCreated event published - SAGA STEP 1 COMPLETE")

	span.AddEvent("order_created", trace.WithAttributes(
		attribute.String("order_id", cmd.ID),
		attribute.String("user_customer_id", cmd.UserCustomerID),
		attribute.String("basket_id", cmd.BasketID),
		attribute.Int("items_count", len(cmd.Items)),
		attribute.String("event_id", event.ID()),
		attribute.String("status", string(order.Status)),
	))

	return nil
}
