package handlers

import (
	"context"
	"time"

	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/internal/registry"
	"middleman/streams/internal/application"
	"middleman/streams/internal/application/commands"
	"middleman/streams/streamspb"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// commandHandlers routes incoming NATS/JetStream commands to the application.
// Only ReserveStream is handled for now.

type commandHandlers struct {
	app application.App
}

var _ ddd.CommandHandler[ddd.Command] = (*commandHandlers)(nil)

// NewCommandHandlers wires registry-based serialization, a reply publisher and middlewares.
func NewCommandHandlers(reg registry.Registry, app application.App, replyPublisher am.ReplyPublisher, mws ...am.MessageHandlerMiddleware) am.MessageHandler {
	return am.NewCommandHandler(reg, replyPublisher, commandHandlers{app: app}, mws...)
}

// RegisterCommandHandlers subscribes the given subscriber to the streams command channel.
func RegisterCommandHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) error {
	_, err := subscriber.Subscribe(streamspb.CommandChannel, handlers, am.MessageFilter{
		streamspb.ReserveStreamCommand,
		streamspb.ReleaseStreamCommand,
		streamspb.ReserveStreamsCommand,
		streamspb.ReleaseStreamsCommand,
	}, am.GroupName("streams-commands"))
	return err
}

// HandleCommand dispatches based on CommandName.
func (h commandHandlers) HandleCommand(ctx context.Context, cmd ddd.Command) (reply ddd.Reply, err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent("error", trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		}
		span.AddEvent("done", trace.WithAttributes(attribute.Int64("ms", time.Since(started).Milliseconds())))
	}(time.Now())

	switch cmd.CommandName() {
	case streamspb.ReserveStreamCommand:
		return h.doReserveStream(ctx, cmd)
	case streamspb.ReleaseStreamCommand:
		return h.doReleaseStream(ctx, cmd)
	case streamspb.ReserveStreamsCommand:
		return h.doReserveStreams(ctx, cmd)
	case streamspb.ReleaseStreamsCommand:
		return h.doReleaseStreams(ctx, cmd)
	}
	return nil, nil
}

func (h commandHandlers) doReserveStream(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*streamspb.ReserveStream)

	if err := h.app.ReserveStream(ctx, commands.ReserveStream{
		StreamID: payload.GetStreamId(),
		Quantity: payload.GetQuantity(),
	}); err != nil {
		return nil, err
	}

	// publish success reply so saga can continue
	return ddd.NewReply(am.SuccessReply, nil), nil
}

func (h commandHandlers) doReleaseStream(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*streamspb.ReleaseStream)

	if err := h.app.ReleaseStream(ctx, commands.ReleaseStream{
		StreamID: payload.GetStreamId(),
		Quantity: payload.GetQuantity(),
	}); err != nil {
		return nil, err
	}

	return ddd.NewReply(am.SuccessReply, nil), nil
}

func (h commandHandlers) doReserveStreams(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*streamspb.ReserveStreams)

	// Process each item in the batch
	for _, item := range payload.GetItems() {
		if err := h.app.ReserveStream(ctx, commands.ReserveStream{
			StreamID: item.GetStreamId(),
			Quantity: item.GetQuantity(),
		}); err != nil {
			// If any reservation fails, return error
			// The saga will handle compensation
			return nil, err
		}
	}

	// All reservations successful
	return ddd.NewReply(am.SuccessReply, nil), nil
}

func (h commandHandlers) doReleaseStreams(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*streamspb.ReleaseStreams)

	// Process each item in the batch
	for _, item := range payload.GetItems() {
		if err := h.app.ReleaseStream(ctx, commands.ReleaseStream{
			StreamID: item.GetStreamId(),
			Quantity: item.GetQuantity(),
		}); err != nil {
			// Continue with other releases even if one fails
			// This is a compensation action, so we try to release as much as possible
			continue
		}
	}

	return ddd.NewReply(am.SuccessReply, nil), nil
}
