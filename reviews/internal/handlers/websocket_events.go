package handlers

import (
	"context"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/registry"
	"middleman/reviews/internal/application"
	"middleman/reviews/internal/application/commands"
	"middleman/reviews/internal/domain"
	"middleman/reviews/reviewspb"

	"time"
)

type websocketHandlers[T ddd.Websocket] struct {
	app    application.App
	logger zerolog.Logger
}

func NewWebsocketEventHandlers(reg registry.Registry, app application.App, logger zerolog.Logger, mws ...am.MessageHandlerMiddleware) am.MessageHandler {
	return am.NewWebsocketHandler(reg, websocketHandlers[ddd.Websocket]{
		app:    app,
		logger: logger,
	}, logger, mws...)
}

func RegisterWebsocketCommandHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler, logger zerolog.Logger) (err error) {

	logger.Info().Msg("Registering websocket commands")
	_, err = subscriber.Subscribe(reviewspb.WebSocketChannel, handlers, am.MessageFilter{
		reviewspb.AddReviewCommand,
	}, am.GroupName("websocket-commands"))
	if err != nil {
		logger.Error().Err(err).Msg("Failed to subscribe to websocket events")
	}

	return err
}

func (h websocketHandlers[T]) HandleWebsocket(ctx context.Context, wsCommand T) (err error) {
	span := trace.SpanFromContext(ctx)
	wsCommandName := wsCommand.WebsocketName()
	started := time.Now()

	h.logger.Info().
		Str("event", wsCommandName).
		Interface("details", wsCommand).
		Msg("Received event")

	defer func() {
		duration := time.Since(started).Milliseconds()
		span.SetAttributes(attribute.String("event.name", wsCommandName))
		span.SetAttributes(attribute.Int64("processing.duration.ms", duration))

		if err != nil {
			h.logger.Error().
				Err(err).
				Str("event", wsCommandName).
				Int64("duration_ms", duration).
				Msg("Error handling websocket event")
			span.RecordError(err)
		} else {
			h.logger.Info().
				Str("event", wsCommandName).
				Int64("duration_ms", duration).
				Msg("Handled websocket event successfully")
		}
		span.End()
	}()

	span.AddEvent("Handling websocket event", trace.WithAttributes(attribute.String("Event", wsCommandName)))

	switch wsCommandName {
	case reviewspb.AddReviewCommand:
		return h.onAddCommand(ctx, wsCommand)
	default:
		h.logger.Warn().
			Str("event", wsCommandName).
			Msg("Unhandled event")
	}

	return nil
}

func (h websocketHandlers[T]) onAddCommand(ctx context.Context, wsCommand T) (err error) {
	payload := wsCommand.Payload().(*reviewspb.AddReview)
	return h.app.AddReview(ctx, commands.AddReview{
		ID:         payload.GetId(),
		ItemID:     payload.GetItemId(),
		ItemType:   domain.ToItemType(payload.GetItemType()),
		SenderID:   payload.GetSenderId(),
		Content:    payload.GetContent(),
		CategoryID: payload.GetCategoryId(),
		ParentID:   payload.GetParentId(),
	})
}
