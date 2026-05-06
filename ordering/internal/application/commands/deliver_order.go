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

type DeliverOrder struct {
	ID string
}

type DeliverOrderHandler struct {
	orders    domain.OrderRepository
	publisher ddd.EventPublisher[ddd.Event]
	logger    zerolog.Logger
}

func NewDeliverOrderHandler(orders domain.OrderRepository, publisher ddd.EventPublisher[ddd.Event]) DeliverOrderHandler {
	return DeliverOrderHandler{
		orders:    orders,
		publisher: publisher,
		logger:    log.With().Str("service", "ordering").Str("handler", "DeliverOrder").Logger(),
	}
}

func (h DeliverOrderHandler) DeliverOrder(ctx context.Context, cmd DeliverOrder) error {
	span := trace.SpanFromContext(ctx)
	logger := h.logger.With().
		Str("operation", "DeliverOrder").
		Str("order_id", cmd.ID).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	startTime := time.Now()
	logger.Info().Msg("ORDERING_DELIVER_ORDER_BEGIN: Starting order delivery")

	order, err := h.orders.Load(ctx, cmd.ID)
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("ORDERING_DELIVER_ORDER_LOAD_FAILED: Failed to load order for delivery")
		return errors.Wrap(err, "deliver order load")
	}

	logger.Debug().
		Str("current_status", string(order.Status)).
		Int("items_count", len(order.Items)).
		Msg("ORDERING_DELIVER_ORDER_LOADED: Order loaded for delivery")

	event, err := order.Deliver()
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("ORDERING_DELIVER_ORDER_DOMAIN_FAILED: Domain deliver order operation failed")
		return errors.Wrap(err, "deliver order domain")
	}

	logger.Info().
		Str("event_name", event.EventName()).
		Str("new_status", string(order.Status)).
		Msg("ORDERING_DELIVER_ORDER_EVENT_CREATED: Order delivered, event generated")

	if err = h.orders.Save(ctx, order); err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("ORDERING_DELIVER_ORDER_SAVE_FAILED: Failed to save delivered order")
		return errors.Wrap(err, "deliver order save")
	}

	logger.Debug().
		Str("final_status", string(order.Status)).
		Msg("ORDERING_DELIVER_ORDER_SAVED: Delivered order saved successfully")

	if err = h.publisher.Publish(ctx, event); err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Str("event_name", event.EventName()).
			Msg("ORDERING_DELIVER_ORDER_PUBLISH_FAILED: Failed to publish OrderDelivered event")
		return errors.Wrap(err, "deliver order publish")
	}

	logger.Info().
		Dur("duration_ms", time.Since(startTime)).
		Str("event_id", event.ID()).
		Str("event_name", event.EventName()).
		Msg("ORDERING_DELIVER_ORDER_SUCCESS: Order delivered and OrderDelivered event published")

	span.AddEvent("order_delivered", trace.WithAttributes(
		attribute.String("order_id", cmd.ID),
		attribute.String("event_id", event.ID()),
		attribute.String("status", string(order.Status)),
	))

	return nil
}
