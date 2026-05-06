package queries

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"middleman/ordering/internal/domain"

	"github.com/stackus/errors"
)

type GetOrder struct {
	ID string
}

type GetOrderHandler struct {
	repo   domain.OrderRepository
	logger zerolog.Logger
}

func NewGetOrderHandler(repo domain.OrderRepository) GetOrderHandler {
	return GetOrderHandler{
		repo:   repo,
		logger: log.With().Str("service", "ordering").Str("handler", "GetOrder").Logger(),
	}
}

func (h GetOrderHandler) GetOrder(ctx context.Context, query GetOrder) (*domain.Order, error) {
	span := trace.SpanFromContext(ctx)
	logger := h.logger.With().
		Str("operation", "GetOrder").
		Str("order_id", query.ID).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	startTime := time.Now()
	logger.Debug().Msg("ORDERING_GET_ORDER_BEGIN: Retrieving order")

	order, err := h.repo.Load(ctx, query.ID)
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("ORDERING_GET_ORDER_FAILED: Failed to retrieve order")
		return nil, errors.Wrap(err, "get order query")
	}

	logger.Debug().
		Str("order_status", string(order.Status)).
		Int("items_count", len(order.Items)).
		Dur("duration_ms", time.Since(startTime)).
		Msg("ORDERING_GET_ORDER_SUCCESS: Order retrieved successfully")

	span.AddEvent("order_retrieved", trace.WithAttributes(
		attribute.String("order_id", query.ID),
		attribute.String("status", string(order.Status)),
		attribute.Int("items_count", len(order.Items)),
	))

	return order, nil
}
