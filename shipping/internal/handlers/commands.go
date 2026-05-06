package handlers

import (
	"context"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/internal/registry"
	"middleman/shipping/internal/application"
	"middleman/shipping/internal/application/commands"
	"middleman/shipping/shippingpb"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type commandHandlers struct {
	app application.App
}

func NewCommandHandlers(reg registry.Registry, app application.App, replyPublisher am.ReplyPublisher, mws ...am.MessageHandlerMiddleware) am.MessageHandler {
	return am.NewCommandHandler(reg, replyPublisher, commandHandlers{app: app}, mws...)
}

func RegisterCommandHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) error {
	_, err := subscriber.Subscribe(shippingpb.CommandChannel, handlers, am.MessageFilter{
		shippingpb.CreateShipmentCommand,
	}, am.GroupName("shipping-commands"))
	return err
}

func (h commandHandlers) HandleCommand(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	var err error
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent("error", trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		}
		span.AddEvent("done", trace.WithAttributes(attribute.Int64("ms", time.Since(started).Milliseconds())))
	}(time.Now())

	switch cmd.CommandName() {
	case shippingpb.CreateShipmentCommand:
		return h.doCreateShipment(ctx, cmd)
	}
	return nil, nil
}

func (h commandHandlers) doCreateShipment(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*shippingpb.CreateShipment)
	if err := h.app.CreateShipping(ctx, commands.CreateShipping{
		ID:       uuid.New().String(),
		OrderID:  payload.GetOrderId(),
		BasketID: payload.GetBasketId(),
	}); err != nil {
		return nil, err
	}

	return ddd.NewReply(am.SuccessReply, nil), nil
}
