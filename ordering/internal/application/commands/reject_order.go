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

type RejectOrder struct {
	ID string
}

type RejectOrderHandler struct {
	orders    domain.OrderRepository
	publisher ddd.EventPublisher[ddd.Event]
	logger    zerolog.Logger
}

func NewRejectOrderHandler(orders domain.OrderRepository, publisher ddd.EventPublisher[ddd.Event]) RejectOrderHandler {
	return RejectOrderHandler{
		orders:    orders,
		publisher: publisher,
		logger:    log.With().Str("service", "ordering").Str("handler", "RejectOrder").Logger(),
	}
}

func (h RejectOrderHandler) RejectOrder(ctx context.Context, cmd RejectOrder) error {
	span := trace.SpanFromContext(ctx)
	logger := h.logger.With().
		Str("operation", "RejectOrder").
		Str("order_id", cmd.ID).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	startTime := time.Now()
	logger.Info().Msg("ORDERING_REJECT_ORDER_BEGIN: Starting order rejection - SAGA COMPENSATION")

	order, err := h.orders.Load(ctx, cmd.ID)
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("ORDERING_REJECT_ORDER_LOAD_FAILED: Failed to load order for rejection")
		return errors.Wrap(err, "reject order load")
	}

	logger.Debug().
		Str("current_status", string(order.Status)).
		Int("items_count", len(order.Items)).
		Msg("ORDERING_REJECT_ORDER_LOADED: Order loaded for rejection")

	event, err := order.Reject()
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("ORDERING_REJECT_ORDER_DOMAIN_FAILED: Domain reject order operation failed")
		return errors.Wrap(err, "reject order domain")
	}

	logger.Info().
		Str("event_name", event.EventName()).
		Str("new_status", string(order.Status)).
		Msg("ORDERING_REJECT_ORDER_EVENT_CREATED: Order rejected, event generated")

	if err = h.orders.Save(ctx, order); err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("ORDERING_REJECT_ORDER_SAVE_FAILED: Failed to save rejected order")
		return errors.Wrap(err, "reject order save")
	}

	logger.Debug().
		Str("final_status", string(order.Status)).
		Msg("ORDERING_REJECT_ORDER_SAVED: Rejected order saved successfully")

	if err = h.publisher.Publish(ctx, event); err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Str("event_name", event.EventName()).
			Msg("ORDERING_REJECT_ORDER_PUBLISH_FAILED: Failed to publish OrderRejected event")
		return errors.Wrap(err, "reject order publish")
	}

	logger.Info().
		Dur("duration_ms", time.Since(startTime)).
		Str("event_id", event.ID()).
		Str("event_name", event.EventName()).
		Msg("ORDERING_REJECT_ORDER_SUCCESS: Order rejected and OrderRejected event published - SAGA COMPENSATION COMPLETE")

	span.AddEvent("order_rejected", trace.WithAttributes(
		attribute.String("order_id", cmd.ID),
		attribute.String("event_id", event.ID()),
		attribute.String("status", string(order.Status)),
	))

	return nil
}
