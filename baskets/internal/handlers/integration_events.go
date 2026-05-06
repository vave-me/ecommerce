package handlers

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"middleman/baskets/internal/domain"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/internal/registry"
	"middleman/products/productspb"
	"middleman/users/userspb"
)

// integrationHandlers implements ddd.EventHandler
type integrationHandlers[T ddd.Event] struct {
	users    domain.UserCacheRepository
	products domain.ProductCacheRepository
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers(
	reg registry.Registry,
	users domain.UserCacheRepository,
	products domain.ProductCacheRepository,
	mws ...am.MessageHandlerMiddleware,
) am.MessageHandler {
	return am.NewEventHandler(
		reg,
		integrationHandlers[ddd.Event]{
			users:    users,
			products: products,
		},
		zerolog.Logger{},
		mws...,
	)
}

func RegisterIntegrationEventHandlers(
	subscriber am.MessageSubscriber,
	handlers am.MessageHandler,
) (err error) {
	// Subscribe to user events
	_, err = subscriber.Subscribe(
		userspb.UserAggregateChannel,
		handlers,
		am.MessageFilter{
			userspb.UserCreatedEvent,
			userspb.UserRenamedEvent,
		},
		am.GroupName("baskets-stores"),
	)
	if err != nil {
		return err
	}

	// Subscribe to product events
	_, err = subscriber.Subscribe(
		productspb.ProductAggregateChannel,
		handlers,
		am.MessageFilter{
			productspb.ProductAddedEvent,
			productspb.ProductRebrandedEvent,
			productspb.ProductPriceIncreasedEvent,
			productspb.ProductPriceDecreasedEvent,
			productspb.ProductRemovedEvent,
		},
		am.GroupName("baskets-products"),
	)
	if err != nil {
		return err
	}

	return nil
}

// HandleEvent is the single entry point for handling events in this integration handler.
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

	// If we get here, it means no known event was handled.
	return nil
}

func (h integrationHandlers[T]) onUserCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*userspb.UserCreated)
	return h.users.Add(
		ctx,
		payload.GetId(),
		payload.GetEmail(),
		payload.GetUserName(),
		payload.GetFirstName(),
		payload.GetLastName(),
		payload.GetLocation(),
		payload.GetEnabled(),
	)
}

func (h integrationHandlers[T]) onUserRenamed(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*userspb.UserRenamed)
	return h.users.Rename(ctx, payload.GetId(), payload.GetName())
}

func protoAttrsToDomainAttrs(pbAttrs []*productspb.Attribute) []domain.Attribute {
	domainAttrs := make([]domain.Attribute, len(pbAttrs))
	for i, attr := range pbAttrs {
		domainAttrs[i] = domain.Attribute{
			Key:   attr.GetKey(),
			Value: attr.GetValue(),
		}
	}
	return domainAttrs
}

// Convert Protobuf Options to Domain Options
func protoOptsToDomainOpts(pbOpts []*productspb.Option) []domain.Option {
	domainOpts := make([]domain.Option, len(pbOpts))
	for i, opt := range pbOpts {
		domainOpts[i] = domain.Option{
			Name:  opt.GetName(),
			Value: opt.GetValue(),
			// If domain.Option.Price is a float64, cast int64 to float64
			Price: float64(opt.GetPrice()),
		}
	}
	return domainOpts
}

func (h integrationHandlers[T]) onProductAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*productspb.ProductAdded)

	// Convert the Protobuf arrays to domain arrays
	domainAttrs := protoAttrsToDomainAttrs(payload.GetAttributes())
	domainOpts := protoOptsToDomainOpts(payload.GetOptions())

	return h.products.Add(
		ctx,
		payload.GetId(),
		payload.GetName(),
		payload.GetDescription(),
		payload.GetBasePrice(),
		payload.GetUserSellerId(),
		payload.GetCategoryId(),
		payload.GetBrand(),
		domain.ToProductCondition(payload.GetCondition()),
		payload.GetModel(),
		payload.GetTags(),
		payload.GetManageStocks(),
		payload.GetStock(),
		payload.GetSku(),
		// Here we pass the domainAttrs instead of a string
		domainAttrs,
		payload.GetWeight(),
		payload.GetHeight(),
		payload.GetWidth(),
		payload.GetDepth(),
		domain.ToProductStatus(payload.GetStatus()),
		payload.GetNegotiable(),
		domain.ToUserType(payload.GetUserType()),
		payload.GetMiddlemanService(),
		payload.GetShippingCost(),
		payload.GetHasVariants(),
		// Here we pass the domainOpts instead of the raw []*productspb.Option
		domainOpts,
		payload.GetThumbnail(),
		float64(payload.GetLat()),
		float64(payload.GetLng()),
	)
}

func (h integrationHandlers[T]) onProductRebranded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*productspb.ProductRebranded)
	return h.products.Rebrand(
		ctx,
		payload.GetId(),
		payload.GetName(),
		payload.GetDescription(),
		payload.GetBasePrice(),
		payload.GetStock(),
		payload.GetSku(),
		payload.GetCategoryId(),
	)
}

func (h integrationHandlers[T]) onProductPriceIncreased(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*productspb.ProductPriceIncreased)
	return h.products.UpdatePrice(ctx, payload.GetId(), payload.GetOldPrice(), payload.GetNewPrice())
}

func (h integrationHandlers[T]) onProductPriceDecreased(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*productspb.ProductPriceDecreased)
	return h.products.UpdatePrice(ctx, payload.GetId(), payload.GetNewPrice(), payload.GetNewPrice())
}

func (h integrationHandlers[T]) onProductRemoved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*productspb.ProductRemoved)
	return h.products.Remove(ctx, payload.GetId())
}
