package handlers

import (
	"context"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/internal/registry"
	"middleman/products/productspb"
	"middleman/users/userspb"
	"middleman/wishlists/internal/application"
	"middleman/wishlists/internal/application/commands"
	"middleman/wishlists/internal/domain"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type integrationHandlers[T ddd.Event] struct {
	app      application.App
	products domain.ProductCacheRepository
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers(reg registry.Registry, products domain.ProductCacheRepository, app application.App, mws ...am.MessageHandlerMiddleware) am.MessageHandler {
	return am.NewEventHandler(reg, integrationHandlers[ddd.Event]{
		app:      app,
		products: products,
	}, zerolog.Logger{}, mws...)
}

func RegisterIntegrationEventHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) (err error) {

	_, err = subscriber.Subscribe(userspb.UserAggregateChannel, handlers, am.MessageFilter{
		userspb.UserCreatedEvent,
	}, am.GroupName("wishlist-users"))

	_, err = subscriber.Subscribe(productspb.ProductAggregateChannel, handlers, am.MessageFilter{
		productspb.ProductAddedEvent,
		productspb.ProductRebrandedEvent,
		productspb.ProductPriceIncreasedEvent,
		productspb.ProductPriceDecreasedEvent,
		productspb.ProductRemovedEvent,
	}, am.GroupName("wishlist-products"))

	return err
}

func (h integrationHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling integration event",
				trace.WithAttributes(errorsotel.ErrAttrs(err)...),
			)
		}
		span.AddEvent("Handled integration event", trace.WithAttributes(
			attribute.Int64("TookMS", time.Since(started).Milliseconds()),
		))
	}(time.Now())

	span.AddEvent("Handling integration event", trace.WithAttributes(
		attribute.String("Event", event.EventName()),
	))

	switch event.EventName() {
	case userspb.UserCreatedEvent:
		return h.onUserCreated(ctx, event)
	case productspb.ProductAddedEvent:
		return h.onProductAdded(ctx, event)
	case productspb.ProductRebrandedEvent:
		return h.onProductRebranded(ctx, event)
	case productspb.ProductPriceIncreasedEvent:
		return h.onProductPriceIncreased(ctx, event)
	case productspb.ProductPriceDecreasedEvent:
		return h.onProductPriceDecreased(ctx, event)
	case productspb.ProductRemovedEvent:
		return h.onProductRemoved(ctx, event)
	}

	return nil
}
func (h integrationHandlers[T]) onUserCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*userspb.UserCreated)
	userID := payload.GetId()
	id := uuid.New().String()
	cmd := commands.CreateWishlist{
		ID:     id,
		Name:   "Default",
		UserID: userID,
	}
	return h.app.CreateWishlist(ctx, cmd)
}

func (h integrationHandlers[T]) onProductAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*productspb.ProductAdded)
	return h.products.Add(ctx, payload.GetId(), payload.GetName(), payload.GetDescription(), payload.GetBasePrice(), payload.GetUserSellerId(), payload.GetStock(), payload.GetSku(), payload.GetCategoryId())
}

func (h integrationHandlers[T]) onProductRebranded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*productspb.ProductRebranded)
	return h.products.Rebrand(ctx, payload.GetId(), payload.GetName())
}

func (h integrationHandlers[T]) onProductPriceIncreased(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*productspb.ProductPriceIncreased)
	return h.products.UpdatePrice(ctx, payload.GetId(), payload.GetNewPrice())
}
func (h integrationHandlers[T]) onProductPriceDecreased(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*productspb.ProductPriceDecreased)
	return h.products.UpdatePrice(ctx, payload.GetId(), payload.GetNewPrice())
}
func (h integrationHandlers[T]) onProductRemoved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*productspb.ProductRemoved)
	return h.products.Remove(ctx, payload.GetId())
}
