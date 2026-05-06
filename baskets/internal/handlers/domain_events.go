package handlers

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"middleman/baskets/basketspb"
	"middleman/baskets/internal/domain"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"

	"time"
)

type domainHandlers[T ddd.Event] struct {
	publisher am.EventPublisher
}

var _ ddd.EventHandler[ddd.Event] = (*domainHandlers[ddd.Event])(nil)

func NewDomainEventHandlers(publisher am.EventPublisher) ddd.EventHandler[ddd.Event] {
	return &domainHandlers[ddd.Event]{
		publisher: publisher,
	}
}

func RegisterDomainEventHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.BasketStartedEvent,
		domain.BasketItemAddedEvent,
		domain.BasketItemRemovedEvent,
		domain.BasketCanceledEvent,
		domain.BasketCheckedOutEvent,
	)
}
func (h domainHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling domain event",
				trace.WithAttributes(errorsotel.ErrAttrs(err)...),
			)
		}
		span.AddEvent("Handled domain event", trace.WithAttributes(
			attribute.Int64("TookMS", time.Since(started).Milliseconds()),
		))
	}(time.Now())

	span.AddEvent("Handling domain event", trace.WithAttributes(
		attribute.String("Event", event.EventName()),
	))

	switch event.EventName() {
	case domain.BasketStartedEvent:
		return h.onBasketStarted(ctx, event)
	case domain.BasketItemAddedEvent:
		return h.onItemAdded(ctx, event)
	case domain.BasketItemRemovedEvent:
		return h.onItemRemoved(ctx, event)
	case domain.BasketCanceledEvent:
		return h.onBasketCanceled(ctx, event)
	case domain.BasketCheckedOutEvent:
		return h.onBasketCheckedOut(ctx, event)
	}
	return nil
}
func (h domainHandlers[T]) onBasketStarted(ctx context.Context, event ddd.Event) error {
	basket := event.Payload().(*domain.Basket)

	return h.publisher.Publish(ctx, basketspb.BasketAggregateChannel,
		ddd.NewEvent(basketspb.BasketStartedEvent, &basketspb.BasketStarted{
			Id: basket.ID(),
		}),
	)
}

func (h domainHandlers[T]) onItemAdded(ctx context.Context, event ddd.Event) error {
	basket := event.Payload().(*domain.Basket)

	items := make([]*basketspb.Item, 0, len(basket.Items))
	for _, item := range basket.Items {
		items = append(items, &basketspb.Item{
			UserSellerId:   item.UserSellerID,
			ProductId:      item.ProductID,
			UserSellerName: item.UserSellerName,
			ProductName:    item.ProductName,
			ProductPrice:   item.ProductPrice,
			Quantity:       item.Quantity,
		})
		err := h.publisher.Publish(ctx, basketspb.BasketAggregateChannel, ddd.NewEvent(basketspb.BasketItemAddedEvent, &basketspb.BasketItemAdded{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
		}))
		if err != nil {
			return err
		}
	}
	return nil
}
func (h domainHandlers[T]) onItemRemoved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Basket)
	return h.publisher.Publish(ctx, basketspb.BasketAggregateChannel,
		ddd.NewEvent(basketspb.BasketItemRemovedEvent, &basketspb.BasketItemRemoved{
			ProductId: payload.ID(),
		}),
	)
}
func (h domainHandlers[T]) onBasketCanceled(ctx context.Context, event ddd.Event) error {
	basket := event.Payload().(*domain.Basket)
	return h.publisher.Publish(ctx, basketspb.BasketAggregateChannel,
		ddd.NewEvent(basketspb.BasketCanceledEvent, &basketspb.BasketCanceled{
			Id: basket.ID(),
		}),
	)
}

func (h domainHandlers[T]) onBasketCheckedOut(ctx context.Context, event ddd.Event) error {
	basket := event.Payload().(*domain.Basket)

	items := make([]*basketspb.BasketCheckedOut_Item, 0, len(basket.Items))
	var total int64
	for _, it := range basket.Items {
		lineTotal := it.ProductPrice * it.Quantity
		total += lineTotal

		items = append(items, &basketspb.BasketCheckedOut_Item{
			UserSellerId:   it.UserSellerID,
			ProductId:      it.ProductID,
			UserSellerName: it.UserSellerName,
			ProductName:    it.ProductName,
			Price:          it.ProductPrice,
			Quantity:       it.Quantity,
		})
	}

	// Now include user_customer_id, payment_method_id, total, items, etc.
	return h.publisher.Publish(ctx, basketspb.BasketAggregateChannel,
		ddd.NewEvent(basketspb.BasketCheckedOutEvent, &basketspb.BasketCheckedOut{
			Id:               basket.ID(),
			UserCustomerId:   basket.UserCustomerID,
			Total:            total,
			Items:            items,
			PaymentIntentId:  basket.PaymentIntentID, // Include the payment intent ID from the basket
		}),
	)
}
