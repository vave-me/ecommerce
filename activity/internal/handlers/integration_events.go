package handlers

import (
	"context"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"middleman/activity/internal/application"
	"middleman/activity/internal/application/commands"
	"middleman/activity/internal/domain"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/internal/registry"
	"middleman/users/userspb"
	"time"
)

type integrationHandlers[T ddd.Event] struct {
	app      application.App
	users    domain.UserCacheRepository
	products domain.ProductCacheRepository
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers(reg registry.Registry, users domain.UserCacheRepository, products domain.ProductCacheRepository, app application.App, mws ...am.MessageHandlerMiddleware) am.MessageHandler {
	return am.NewEventHandler(reg, integrationHandlers[ddd.Event]{
		app:      app,
		users:    users,
		products: products,
	}, zerolog.Logger{}, mws...)
}

func RegisterIntegrationEventHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) (err error) {
	_, err = subscriber.Subscribe(userspb.UserAggregateChannel, handlers, am.MessageFilter{
		userspb.UserCreatedEvent,
		userspb.UserRenamedEvent,
		userspb.UserLoggedInEvent,
	}, am.GroupName("activity-users"))
	if err != nil {
		return err
	}

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
	case userspb.UserRenamedEvent:
		return h.onUserRenamed(ctx, event)
		//case userspb.ProductAddedEvent:
		//	return h.onProductAdded(ctx, event)
		//case userspb.ProductRebrandedEvent:
		//	return h.onProductRebranded(ctx, event)
		//case userspb.ProductRemovedEvent:
		//	return h.onProductRemoved(ctx, event)
	}

	return nil
}
func (h integrationHandlers[T]) onUserCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*userspb.UserCreated)
	userID := payload.GetId()
	id := uuid.New().String()
	cmd := commands.CreateActivity{
		ID:     id,
		UserID: userID,
	}
	return h.app.CreateActivity(ctx, cmd)
}
func (h integrationHandlers[T]) onUserRenamed(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*userspb.UserRenamed)
	return h.users.Rename(ctx, payload.GetId(), payload.GetName())
}

//func (h integrationHandlers[T]) onProductAdded(ctx context.Context, event ddd.Event) error {
//	payload := event.Payload().(*userspb.ProductAdded)
//	return h.products.Add(ctx, payload.GetId(), payload.GetName(), payload.GetDescription(), payload.GetPrice(), payload.GetUserSellerId(), payload.GetStock(), payload.GetSku(), payload.GetCategoryId(), payload.GetActive())
//}
//func (h integrationHandlers[T]) onProductRebranded(ctx context.Context, event ddd.Event) error {
//	payload := event.Payload().(*userspb.ProductRebranded)
//	return h.products.Rebrand(ctx, payload.GetId(), payload.GetName())
//}
//func (h integrationHandlers[T]) onProductPriceChanged(ctx context.Context, event ddd.Event) error {
//	payload := event.Payload().(*userspb.ProductPriceChanged)
//	return h.products.UpdatePrice(ctx, payload.GetId(), payload.GetPriceChange())
//}
//func (h integrationHandlers[T]) onProductRemoved(ctx context.Context, event ddd.Event) error {
//	payload := event.Payload().(*userspb.ProductRemoved)
//	return h.products.Remove(ctx, payload.GetId())
//}
