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

type CancelOrder struct {
	ID string
}

type CancelOrderHandler struct {
	orders    domain.OrderRepository
	publisher ddd.EventPublisher[ddd.Event]
	logger    zerolog.Logger
}

func NewCancelOrderHandler(orders domain.OrderRepository, publisher ddd.EventPublisher[ddd.Event]) CancelOrderHandler {
	return CancelOrderHandler{
		orders:    orders,
		publisher: publisher,
		logger:    log.With().Str("service", "ordering").Str("handler", "CancelOrder").Logger(),
	}
}

func (h CancelOrderHandler) CancelOrder(ctx context.Context, cmd CancelOrder) error {
	span := trace.SpanFromContext(ctx)
	logger := h.logger.With().
		Str("operation", "CancelOrder").
		Str("order_id", cmd.ID).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	startTime := time.Now()
	logger.Info().Msg("ORDERING_CANCEL_ORDER_BEGIN: Starting order cancellation")

	order, err := h.orders.Load(ctx, cmd.ID)
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("ORDERING_CANCEL_ORDER_LOAD_FAILED: Failed to load order for cancellation")
		return errors.Wrap(err, "cancel order load")
	}

	logger.Debug().
		Str("current_status", string(order.Status)).
		Int("items_count", len(order.Items)).
		Msg("ORDERING_CANCEL_ORDER_LOADED: Order loaded for cancellation")

	event, err := order.Cancel()
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("ORDERING_CANCEL_ORDER_DOMAIN_FAILED: Domain cancel order operation failed")
		return errors.Wrap(err, "cancel order domain")
	}

	logger.Info().
		Str("event_name", event.EventName()).
		Str("new_status", string(order.Status)).
		Msg("ORDERING_CANCEL_ORDER_EVENT_CREATED: Order cancelled, event generated")

	// // TODO CH8 remove; handled in the saga
	// if err = h.shopping.Cancel(ctx, order.ShoppingID); err != nil {
	// 	return err
	// }

	if err = h.orders.Save(ctx, order); err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("ORDERING_CANCEL_ORDER_SAVE_FAILED: Failed to save cancelled order")
		return errors.Wrap(err, "cancel order save")
	}

	logger.Debug().
		Str("final_status", string(order.Status)).
		Msg("ORDERING_CANCEL_ORDER_SAVED: Cancelled order saved successfully")

	if err = h.publisher.Publish(ctx, event); err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Str("event_name", event.EventName()).
			Msg("ORDERING_CANCEL_ORDER_PUBLISH_FAILED: Failed to publish OrderCancelled event")
		return errors.Wrap(err, "cancel order publish")
	}

	logger.Info().
		Dur("duration_ms", time.Since(startTime)).
		Str("event_id", event.ID()).
		Str("event_name", event.EventName()).
		Msg("ORDERING_CANCEL_ORDER_SUCCESS: Order cancelled and OrderCancelled event published")

	span.AddEvent("order_cancelled", trace.WithAttributes(
		attribute.String("order_id", cmd.ID),
		attribute.String("event_id", event.ID()),
		attribute.String("status", string(order.Status)),
	))

	return nil
}
