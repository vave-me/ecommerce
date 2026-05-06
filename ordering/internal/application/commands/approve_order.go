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

type ApproveOrder struct {
	ID         string
	ShoppingID string
}

type ApproveOrderHandler struct {
	orders    domain.OrderRepository
	publisher ddd.EventPublisher[ddd.Event]
	logger    zerolog.Logger
}

func NewApproveOrderHandler(orders domain.OrderRepository, publisher ddd.EventPublisher[ddd.Event]) ApproveOrderHandler {
	return ApproveOrderHandler{
		orders:    orders,
		publisher: publisher,
		logger:    log.With().Str("service", "ordering").Str("handler", "ApproveOrder").Logger(),
	}
}

func (h ApproveOrderHandler) ApproveOrder(ctx context.Context, cmd ApproveOrder) error {
	span := trace.SpanFromContext(ctx)
	logger := h.logger.With().
		Str("operation", "ApproveOrder").
		Str("order_id", cmd.ID).
		Str("shopping_id", cmd.ShoppingID).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	startTime := time.Now()
	logger.Info().Msg("ORDERING_APPROVE_ORDER_BEGIN: Starting order approval - SAGA STEP 5")

	order, err := h.orders.Load(ctx, cmd.ID)
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("ORDERING_APPROVE_ORDER_LOAD_FAILED: Failed to load order for approval")
		return err
	}

	logger.Debug().
		Str("current_status", string(order.Status)).
		Int("items_count", len(order.Items)).
		Msg("ORDERING_APPROVE_ORDER_LOADED: Order loaded for approval")

	event, err := order.Approve(cmd.ShoppingID)
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("ORDERING_APPROVE_ORDER_DOMAIN_FAILED: Domain approve order operation failed")
		return errors.Wrap(err, "approve order command")
	}

	logger.Info().
		Str("event_name", event.EventName()).
		Str("new_status", string(order.Status)).
		Msg("ORDERING_APPROVE_ORDER_EVENT_CREATED: Order approved, event generated")

	if err = h.orders.Save(ctx, order); err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("ORDERING_APPROVE_ORDER_SAVE_FAILED: Failed to save approved order")
		return errors.Wrap(err, "order approval")
	}

	logger.Debug().
		Str("final_status", string(order.Status)).
		Msg("ORDERING_APPROVE_ORDER_SAVED: Approved order saved successfully")

	if err = h.publisher.Publish(ctx, event); err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Str("event_name", event.EventName()).
			Msg("ORDERING_APPROVE_ORDER_PUBLISH_FAILED: Failed to publish OrderApproved event")
		return err
	}

	logger.Info().
		Dur("duration_ms", time.Since(startTime)).
		Str("event_id", event.ID()).
		Str("event_name", event.EventName()).
		Msg("ORDERING_APPROVE_ORDER_SUCCESS: Order approved and OrderApproved event published - SAGA STEP 5 COMPLETE")

	span.AddEvent("order_approved", trace.WithAttributes(
		attribute.String("order_id", cmd.ID),
		attribute.String("shopping_id", cmd.ShoppingID),
		attribute.String("event_id", event.ID()),
		attribute.String("status", string(order.Status)),
	))

	return nil
}
