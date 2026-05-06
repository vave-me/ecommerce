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

type ShipOrder struct {
	ID string
}

type ShipOrderHandler struct {
	orders    domain.OrderRepository
	publisher ddd.EventPublisher[ddd.Event]
	logger    zerolog.Logger
}

func NewShipOrderHandler(orders domain.OrderRepository, publisher ddd.EventPublisher[ddd.Event]) ShipOrderHandler {
	return ShipOrderHandler{
		orders:    orders,
		publisher: publisher,
		logger:    log.With().Str("service", "ordering").Str("handler", "ShipOrder").Logger(),
	}
}

func (h ShipOrderHandler) ShipOrder(ctx context.Context, cmd ShipOrder) error {
	span := trace.SpanFromContext(ctx)
	logger := h.logger.With().
		Str("operation", "ShipOrder").
		Str("order_id", cmd.ID).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	startTime := time.Now()
	logger.Info().Msg("ORDERING_SHIP_ORDER_BEGIN: Starting order shipment")

	order, err := h.orders.Load(ctx, cmd.ID)
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("ORDERING_SHIP_ORDER_LOAD_FAILED: Failed to load order for shipment")
		return errors.Wrap(err, "ship order load")
	}

	logger.Debug().
		Str("current_status", string(order.Status)).
		Int("items_count", len(order.Items)).
		Msg("ORDERING_SHIP_ORDER_LOADED: Order loaded for shipment")

	event, err := order.Ship()
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("ORDERING_SHIP_ORDER_DOMAIN_FAILED: Domain ship order operation failed")
		return errors.Wrap(err, "ship order domain")
	}

	logger.Info().
		Str("event_name", event.EventName()).
		Str("new_status", string(order.Status)).
		Msg("ORDERING_SHIP_ORDER_EVENT_CREATED: Order shipped, event generated")

	if err = h.orders.Save(ctx, order); err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("ORDERING_SHIP_ORDER_SAVE_FAILED: Failed to save shipped order")
		return errors.Wrap(err, "ship order save")
	}

	logger.Debug().
		Str("final_status", string(order.Status)).
		Msg("ORDERING_SHIP_ORDER_SAVED: Shipped order saved successfully")

	if err = h.publisher.Publish(ctx, event); err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Str("event_name", event.EventName()).
			Msg("ORDERING_SHIP_ORDER_PUBLISH_FAILED: Failed to publish OrderShipped event")
		return errors.Wrap(err, "ship order publish")
	}

	logger.Info().
		Dur("duration_ms", time.Since(startTime)).
		Str("event_id", event.ID()).
		Str("event_name", event.EventName()).
		Msg("ORDERING_SHIP_ORDER_SUCCESS: Order shipped and OrderShipped event published")

	span.AddEvent("order_shipped", trace.WithAttributes(
		attribute.String("order_id", cmd.ID),
		attribute.String("event_id", event.ID()),
		attribute.String("status", string(order.Status)),
	))

	return nil
}
