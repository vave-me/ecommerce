package handlers

import (
	"context"
	"time"

	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/internal/registry"
	"middleman/products/internal/application"
	"middleman/products/internal/application/commands"
	"middleman/products/productspb"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// commandHandlers routes incoming NATS/JetStream commands to the application.
// Only ReserveProduct is handled for now.

type commandHandlers struct {
	app application.App
}

var _ ddd.CommandHandler[ddd.Command] = (*commandHandlers)(nil)

// NewCommandHandlers wires registry-based serialization, a reply publisher and middlewares.
func NewCommandHandlers(reg registry.Registry, app application.App, replyPublisher am.ReplyPublisher, mws ...am.MessageHandlerMiddleware) am.MessageHandler {
	return am.NewCommandHandler(reg, replyPublisher, commandHandlers{app: app}, mws...)
}

// RegisterCommandHandlers subscribes the given subscriber to the products command channel.
func RegisterCommandHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) error {
	_, err := subscriber.Subscribe(productspb.CommandChannel, handlers, am.MessageFilter{
		productspb.ReserveProductCommand,
		productspb.ReleaseProductCommand,
		productspb.ReserveProductsCommand,
		productspb.ReleaseProductsCommand,
	}, am.GroupName("products-commands"))
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
	case productspb.ReserveProductCommand:
		return h.doReserveProduct(ctx, cmd)
	case productspb.ReleaseProductCommand:
		return h.doReleaseProduct(ctx, cmd)
	case productspb.ReserveProductsCommand:
		return h.doReserveProducts(ctx, cmd)
	case productspb.ReleaseProductsCommand:
		return h.doReleaseProducts(ctx, cmd)
	}
	return nil, nil
}

func (h commandHandlers) doReserveProduct(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*productspb.ReserveProduct)

	if err := h.app.ReserveProduct(ctx, commands.ReserveProduct{
		ProductID: payload.GetProductId(),
		Quantity:  payload.GetQuantity(),
	}); err != nil {
		return nil, err
	}

	// publish success reply so saga can continue
	return ddd.NewReply(am.SuccessReply, nil), nil
}

func (h commandHandlers) doReleaseProduct(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*productspb.ReleaseProduct)

	if err := h.app.ReleaseProduct(ctx, commands.ReleaseProduct{
		ProductID: payload.GetProductId(),
		Quantity:  payload.GetQuantity(),
	}); err != nil {
		return nil, err
	}

	return ddd.NewReply(am.SuccessReply, nil), nil
}

func (h commandHandlers) doReserveProducts(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*productspb.ReserveProducts)

	// Process each item in the batch
	for _, item := range payload.GetItems() {
		if err := h.app.ReserveProduct(ctx, commands.ReserveProduct{
			ProductID: item.GetProductId(),
			Quantity:  item.GetQuantity(),
		}); err != nil {
			// If any reservation fails, return error
			// The saga will handle compensation
			return nil, err
		}
	}

	// All reservations successful
	return ddd.NewReply(am.SuccessReply, nil), nil
}

func (h commandHandlers) doReleaseProducts(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*productspb.ReleaseProducts)

	// Process each item in the batch
	for _, item := range payload.GetItems() {
		if err := h.app.ReleaseProduct(ctx, commands.ReleaseProduct{
			ProductID: item.GetProductId(),
			Quantity:  item.GetQuantity(),
		}); err != nil {
			// Continue with other releases even if one fails
			// This is a compensation action, so we try to release as much as possible
			continue
		}
	}

	return ddd.NewReply(am.SuccessReply, nil), nil
}
