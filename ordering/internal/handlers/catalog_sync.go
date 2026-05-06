package handlers

import (
	"context"
	"middleman/internal/ddd"
	"middleman/internal/di"
	"middleman/ordering/internal/domain"
)

type catalogHandlers[T ddd.Event] struct {
	catalog domain.CatalogRepository
}

func NewCatalogHandlers(catalog domain.CatalogRepository) ddd.EventHandler[ddd.Event] {
	return catalogHandlers[ddd.Event]{catalog: catalog}
}

func RegisterCatalogHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.OrderCreatedEvent,
		domain.OrderRejectedEvent,
		domain.OrderApprovedEvent,
		domain.OrderCanceledEvent,
		domain.OrderReadiedEvent,
		domain.OrderShippedEvent,
		domain.OrderDeliveredEvent,
		domain.OrderCompletedEvent,
	)
}
func RegisterCatalogHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		catalogHandlers := di.Get(ctx, "catalogHandlers").(ddd.EventHandler[ddd.Event])
		return catalogHandlers.HandleEvent(ctx, event)
	})
	subscriber := container.Get("domainDispatcher").(*ddd.EventDispatcher[ddd.Event])

	RegisterCatalogHandlers(subscriber, handlers)
}
func (h catalogHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.OrderCreatedEvent:
		return h.onOrderCreated(ctx, event)
	case domain.OrderRejectedEvent,
		domain.OrderApprovedEvent,
		domain.OrderCanceledEvent,
		domain.OrderReadiedEvent,
		domain.OrderShippedEvent,
		domain.OrderDeliveredEvent,
		domain.OrderCompletedEvent:
		return h.onOrderUpdated(ctx, event)
	}
	return nil
}

func (h catalogHandlers[T]) onOrderCreated(ctx context.Context, event ddd.Event) error {
	order := event.Payload().(*domain.Order)
	return h.catalog.AddOrder(ctx, order)
}

func (h catalogHandlers[T]) onOrderUpdated(ctx context.Context, event ddd.Event) error {
	order := event.Payload().(*domain.Order)
	return h.catalog.UpdateOrder(ctx, order)
}
