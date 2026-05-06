package handlers

import (
	"context"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/internal/registry"
	"middleman/ordering/internal/application"
	"middleman/ordering/internal/application/commands"
	"middleman/ordering/internal/domain"
	"middleman/ordering/orderingpb"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type commandHandlers struct {
	app application.App
}

func NewCommandHandlers(reg registry.Registry, app application.App, replyPublisher am.ReplyPublisher, mws ...am.MessageHandlerMiddleware) am.MessageHandler {
	return am.NewCommandHandler(reg, replyPublisher, commandHandlers{
		app: app,
	}, mws...)
}

func RegisterCommandHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) error {
	_, err := subscriber.Subscribe(orderingpb.CommandChannel, handlers, am.MessageFilter{
		orderingpb.CreateOrderCommand,
		orderingpb.RejectOrderCommand,
		orderingpb.ApproveOrderCommand,
	}, am.GroupName("ordering-commands"))
	return err
}

func (h commandHandlers) HandleCommand(ctx context.Context, cmd ddd.Command) (reply ddd.Reply, err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling command",
				trace.WithAttributes(errorsotel.ErrAttrs(err)...),
			)
		}
		span.AddEvent("Handled command", trace.WithAttributes(
			attribute.Int64("TookMS", time.Since(started).Milliseconds()),
		))
	}(time.Now())

	span.AddEvent("Handling command", trace.WithAttributes(
		attribute.String("Command", cmd.CommandName()),
	))

	switch cmd.CommandName() {
	case orderingpb.CreateOrderCommand:
		return h.doCreateOrder(ctx, cmd)
	case orderingpb.RejectOrderCommand:
		return h.doRejectOrder(ctx, cmd)
	case orderingpb.ApproveOrderCommand:
		return h.doApproveOrder(ctx, cmd)
	}

	return nil, nil
}

func (h commandHandlers) doRejectOrder(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*orderingpb.RejectOrder)

	return nil, h.app.RejectOrder(ctx, commands.RejectOrder{ID: payload.GetId()})
}

func (h commandHandlers) doApproveOrder(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*orderingpb.ApproveOrder)

	return nil, h.app.ApproveOrder(ctx, commands.ApproveOrder{
		ID:         payload.GetId(),
		ShoppingID: payload.GetShoppingId(),
	})
}

func (h commandHandlers) doCreateOrder(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*orderingpb.CreateOrderRequest)

	// Convert protobuf Items to domain Items
	items := make([]domain.Item, len(payload.GetItems()))
	for i, pbItem := range payload.GetItems() {
		items[i] = domain.Item{
			UserSellerID:   pbItem.GetUserSellerId(),
			ProductID:      pbItem.GetProductId(),
			UserSellerName: pbItem.GetUserSellerName(),
			ProductName:    pbItem.GetProductName(),
			Price:          pbItem.GetPrice(),
			Quantity:       pbItem.GetQuantity(),
		}
	}

	return nil, h.app.CreateOrder(ctx, commands.CreateOrder{
		ID:             payload.GetId(),
		UserCustomerID: payload.GetUserCustomerId(),
		BasketID:       payload.GetBasketId(),
		Items:          items,
		PaymentIntent:  payload.GetPaymentIntent(),
	})
}
