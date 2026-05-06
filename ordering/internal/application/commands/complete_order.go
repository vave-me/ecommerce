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

type CompleteOrder struct {
	ID        string
	InvoiceID string
}

type CompleteOrderHandler struct {
	orders    domain.OrderRepository
	publisher ddd.EventPublisher[ddd.Event]
	logger    zerolog.Logger
}

func NewCompleteOrderHandler(orders domain.OrderRepository, publisher ddd.EventPublisher[ddd.Event]) CompleteOrderHandler {
	return CompleteOrderHandler{
		orders:    orders,
		publisher: publisher,
		logger:    log.With().Str("service", "ordering").Str("handler", "CompleteOrder").Logger(),
	}
}

func (h CompleteOrderHandler) CompleteOrder(ctx context.Context, cmd CompleteOrder) error {
	span := trace.SpanFromContext(ctx)
	logger := h.logger.With().
		Str("operation", "CompleteOrder").
		Str("order_id", cmd.ID).
		Str("invoice_id", cmd.InvoiceID).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	startTime := time.Now()
	logger.Info().Msg("ORDERING_COMPLETE_ORDER_BEGIN: Starting order completion")

	order, err := h.orders.Load(ctx, cmd.ID)
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("ORDERING_COMPLETE_ORDER_LOAD_FAILED: Failed to load order for completion")
		return errors.Wrap(err, "complete order load")
	}

	logger.Debug().
		Str("current_status", string(order.Status)).
		Int("items_count", len(order.Items)).
		Msg("ORDERING_COMPLETE_ORDER_LOADED: Order loaded for completion")

	event, err := order.Complete(cmd.InvoiceID)
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("ORDERING_COMPLETE_ORDER_DOMAIN_FAILED: Domain complete order operation failed")
		return errors.Wrap(err, "complete order domain")
	}

	logger.Info().
		Str("event_name", event.EventName()).
		Str("new_status", string(order.Status)).
		Msg("ORDERING_COMPLETE_ORDER_EVENT_CREATED: Order completed, event generated")

	if err = h.orders.Save(ctx, order); err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("ORDERING_COMPLETE_ORDER_SAVE_FAILED: Failed to save completed order")
		return errors.Wrap(err, "complete order save")
	}

	logger.Debug().
		Str("final_status", string(order.Status)).
		Msg("ORDERING_COMPLETE_ORDER_SAVED: Completed order saved successfully")

	if err = h.publisher.Publish(ctx, event); err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Str("event_name", event.EventName()).
			Msg("ORDERING_COMPLETE_ORDER_PUBLISH_FAILED: Failed to publish OrderCompleted event")
		return errors.Wrap(err, "complete order publish")
	}

	logger.Info().
		Dur("duration_ms", time.Since(startTime)).
		Str("event_id", event.ID()).
		Str("event_name", event.EventName()).
		Msg("ORDERING_COMPLETE_ORDER_SUCCESS: Order completed and OrderCompleted event published")

	span.AddEvent("order_completed", trace.WithAttributes(
		attribute.String("order_id", cmd.ID),
		attribute.String("invoice_id", cmd.InvoiceID),
		attribute.String("event_id", event.ID()),
		attribute.String("status", string(order.Status)),
	))

	return nil
}
