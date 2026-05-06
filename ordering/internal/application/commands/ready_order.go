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

type ReadyOrder struct {
	ID string
}

type ReadyOrderHandler struct {
	orders    domain.OrderRepository
	publisher ddd.EventPublisher[ddd.Event]
	logger    zerolog.Logger
}

func NewReadyOrderHandler(orders domain.OrderRepository, publisher ddd.EventPublisher[ddd.Event]) ReadyOrderHandler {
	return ReadyOrderHandler{
		orders:    orders,
		publisher: publisher,
		logger:    log.With().Str("service", "ordering").Str("handler", "ReadyOrder").Logger(),
	}
}

func (h ReadyOrderHandler) ReadyOrder(ctx context.Context, cmd ReadyOrder) error {
	span := trace.SpanFromContext(ctx)
	logger := h.logger.With().
		Str("operation", "ReadyOrder").
		Str("order_id", cmd.ID).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	startTime := time.Now()
	logger.Info().Msg("ORDERING_READY_ORDER_BEGIN: Starting order ready status update")

	order, err := h.orders.Load(ctx, cmd.ID)
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("ORDERING_READY_ORDER_LOAD_FAILED: Failed to load order for ready status")
		return errors.Wrap(err, "ready order load")
	}

	logger.Debug().
		Str("current_status", string(order.Status)).
		Int("items_count", len(order.Items)).
		Msg("ORDERING_READY_ORDER_LOADED: Order loaded for ready status update")

	event, err := order.Ready()
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("ORDERING_READY_ORDER_DOMAIN_FAILED: Domain ready order operation failed")
		return errors.Wrap(err, "ready order domain")
	}

	logger.Info().
		Str("event_name", event.EventName()).
		Str("new_status", string(order.Status)).
		Msg("ORDERING_READY_ORDER_EVENT_CREATED: Order marked as ready, event generated")

	if err = h.orders.Save(ctx, order); err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("ORDERING_READY_ORDER_SAVE_FAILED: Failed to save ready order")
		return errors.Wrap(err, "ready order save")
	}

	logger.Debug().
		Str("final_status", string(order.Status)).
		Msg("ORDERING_READY_ORDER_SAVED: Ready order saved successfully")

	if err = h.publisher.Publish(ctx, event); err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Str("event_name", event.EventName()).
			Msg("ORDERING_READY_ORDER_PUBLISH_FAILED: Failed to publish OrderReady event")
		return errors.Wrap(err, "ready order publish")
	}

	logger.Info().
		Dur("duration_ms", time.Since(startTime)).
		Str("event_id", event.ID()).
		Str("event_name", event.EventName()).
		Msg("ORDERING_READY_ORDER_SUCCESS: Order marked as ready and OrderReady event published")

	span.AddEvent("order_ready", trace.WithAttributes(
		attribute.String("order_id", cmd.ID),
		attribute.String("event_id", event.ID()),
		attribute.String("status", string(order.Status)),
	))

	return nil
}
