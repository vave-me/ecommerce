package handlers

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/products/internal/domain"
	"middleman/products/productspb"

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
		domain.ProductAddedEvent,
		domain.ProductUpdatedEvent,
		domain.ProductRebrandedEvent,
		domain.ProductPriceIncreasedEvent,
		domain.ProductPriceDecreasedEvent,
		domain.ProductStockAdjustedEvent,
		domain.ProductRemovedEvent,
		domain.ProductArchivedEvent,
		domain.ProductSoldEvent,
		domain.ProductLeasedEvent,
		domain.ProductPawnedEvent,
		domain.ProductThumbnailAddedEvent,
		domain.ProductThumbnailUpdatedEvent,

		// Variant events
		domain.VariantAddedEvent,
		domain.VariantPriceDecreasedEvent,
		domain.VariantPriceIncreasedEvent,
		domain.VariantStockAdjustedEvent,
		domain.VariantArchivedEvent,
		domain.VariantRemovedEvent,
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

	case domain.ProductAddedEvent:
		return h.onProductAdded(ctx, event)
	case domain.ProductUpdatedEvent:
		return h.onProductUpdated(ctx, event)
	case domain.ProductRebrandedEvent:
		return h.onProductRebranded(ctx, event)
	case domain.ProductPriceIncreasedEvent:
		return h.onProductPriceIncreased(ctx, event)
	case domain.ProductPriceDecreasedEvent:
		return h.onProductPriceDecreased(ctx, event)
	case domain.ProductStockAdjustedEvent:
		return h.onProductStockAdjusted(ctx, event)
	case domain.ProductArchivedEvent:
		return h.onProductArchived(ctx, event)
	case domain.ProductRemovedEvent:
		return h.onProductRemoved(ctx, event)
	case domain.ProductSoldEvent:
		return h.onProductSold(ctx, event)
	case domain.ProductLeasedEvent:
		return h.onProductLeased(ctx, event)
	case domain.ProductPawnedEvent:
		return h.onProductPawned(ctx, event)

	// -------------------------------------------------------------------------
	// VARIANT EVENTS
	// -------------------------------------------------------------------------
	case domain.VariantAddedEvent:
		return h.onVariantAdded(ctx, event)
	case domain.VariantPriceDecreasedEvent:
		return h.onVariantPriceDecreased(ctx, event)
	case domain.VariantPriceIncreasedEvent:
		return h.onVariantPriceIncreased(ctx, event)
	case domain.VariantStockAdjustedEvent:
		return h.onVariantStockAdjusted(ctx, event)
	case domain.VariantArchivedEvent:
		return h.onVariantArchived(ctx, event)
	case domain.VariantRemovedEvent:
		return h.onVariantRemoved(ctx, event)
	}
	return nil

}

func (h domainHandlers[T]) onProductAdded(ctx context.Context, event ddd.Event) error {
	product := event.Payload().(*domain.Product)
	return h.publisher.Publish(ctx, productspb.ProductAggregateChannel,
		ddd.NewEvent(productspb.ProductAddedEvent, &productspb.ProductAdded{
			Id:               product.ID(),
			Name:             product.Name,
			Description:      product.Description,
			BasePrice:        product.BasePrice,
			UserSellerId:     product.UserSellerID,
			CategoryId:       product.CategoryID,
			CategorySlug:     product.CategorySlug,
			Brand:            product.Brand,
			Condition:        product.Condition.String(),
			Model:            product.Model,
			Tags:             product.Tags,
			ManageStocks:     product.ManageStock,
			Stock:            product.Stock,
			Sku:              product.SKU,
			Attributes:       domainAttributesToProto(product.Attributes),
			Weight:           product.Weight,
			Height:           product.Height,
			Width:            product.Width,
			Depth:            product.Depth,
			Status:           product.Status.String(),
			Negotiable:       product.Negotiable,
			UserType:         product.UserType.String(),
			MiddlemanService: product.MiddlemanService,
			ShippingCost:     product.ShippingCost,
			HasVariants:      product.HasVariants,
			Options:          domainOptionsToProto(product.Options),
			Lat:              float32(product.Lat),
			Lng:              float32(product.Lng),
		}),
	)
}
func (h domainHandlers[T]) onProductUpdated(ctx context.Context, event ddd.Event) error {
	product := event.Payload().(*domain.Product)
	return h.publisher.Publish(ctx, productspb.ProductAggregateChannel,
		ddd.NewEvent(productspb.ProductUpdatedEvent, &productspb.ProductUpdated{
			Id:               product.ID(),
			Name:             product.Name,
			Description:      product.Description,
			BasePrice:        product.BasePrice,
			UserSellerId:     product.UserSellerID,
			CategoryId:       product.CategoryID,
			CategorySlug:     product.CategorySlug,
			Brand:            product.Brand,
			Condition:        product.Condition.String(),
			Model:            product.Model,
			Tags:             product.Tags,
			ManageStocks:     product.ManageStock,
			Stock:            product.Stock,
			Sku:              product.SKU,
			Attributes:       domainAttributesToProto(product.Attributes),
			Weight:           product.Weight,
			Height:           product.Height,
			Width:            product.Width,
			Depth:            product.Depth,
			Status:           product.Status.String(),
			Negotiable:       product.Negotiable,
			UserType:         product.UserType.String(),
			MiddlemanService: product.MiddlemanService,
			ShippingCost:     product.ShippingCost,
			HasVariants:      product.HasVariants,
			Options:          domainOptionsToProto(product.Options),
		}),
	)
}

func (h domainHandlers[T]) onProductRebranded(ctx context.Context, event ddd.Event) error {
	product := event.Payload().(*domain.Product)
	return h.publisher.Publish(ctx, productspb.ProductAggregateChannel,
		ddd.NewEvent(productspb.ProductRebrandedEvent, &productspb.ProductRebranded{
			Id:          product.ID(),
			Name:        product.Name,
			Description: product.Description,
		}),
	)
}

func (h domainHandlers[T]) onProductPriceIncreased(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ProductPriceIncreased)
	return h.publisher.Publish(ctx, productspb.ProductAggregateChannel,
		ddd.NewEvent(productspb.ProductPriceIncreasedEvent, &productspb.ProductPriceIncreased{
			Id:       payload.ProductID,
			OldPrice: payload.OldPrice,
			NewPrice: payload.NewPrice,
		}),
	)
}

func (h domainHandlers[T]) onProductArchived(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.ProductArchived)
	return h.publisher.Publish(ctx, productspb.ProductAggregateChannel,
		ddd.NewEvent(productspb.ProductArchivedEvent, &productspb.ProductArchived{
			Id:           payload.ProductID,
			UserSellerId: payload.UserSellerID,
		}),
	)
}

func (h domainHandlers[T]) onProductSold(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.ProductSold)
	return h.publisher.Publish(ctx, productspb.ProductAggregateChannel,
		ddd.NewEvent(productspb.ProductSoldEvent, &productspb.ProductSold{
			Id:           payload.ProductID,
			UserSellerId: payload.UserSellerID,

			// BuyerId:   payload.BuyerID, if your proto has it
		}),
	)
}

func (h domainHandlers[T]) onProductLeased(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.ProductLeased)
	return h.publisher.Publish(ctx, productspb.ProductAggregateChannel,
		ddd.NewEvent(productspb.ProductLeasedEvent, &productspb.ProductLeased{
			Id:           payload.ProductID,
			UserSellerId: payload.UserSellerID,

			// LesseeId: payload.LesseeID, if needed
		}),
	)
}

func (h domainHandlers[T]) onProductPriceDecreased(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ProductPriceDecreased)
	return h.publisher.Publish(ctx, productspb.ProductAggregateChannel,
		ddd.NewEvent(productspb.ProductPriceDecreasedEvent, &productspb.ProductPriceDecreased{
			Id:       payload.ProductID,
			OldPrice: payload.OldPrice,
			NewPrice: payload.NewPrice,
		}),
	)
}
func (h domainHandlers[T]) onProductPawned(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.ProductPawned)
	return h.publisher.Publish(ctx, productspb.ProductAggregateChannel,
		ddd.NewEvent(productspb.ProductPawnedEvent, &productspb.ProductPawned{
			Id:           payload.ProductID,
			UserSellerId: payload.UserSellerID,
		}),
	)
}

func (h domainHandlers[T]) onProductStockAdjusted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ProductStockAdjusted)
	return h.publisher.Publish(ctx, productspb.ProductAggregateChannel,
		ddd.NewEvent(productspb.ProductStockAdjustedEvent, &productspb.ProductStockAdjusted{
			Id:       payload.ProductID,
			OldStock: payload.OldStock,
			NewStock: payload.NewStock,
		}),
	)
}
func (h domainHandlers[T]) onProductRemoved(ctx context.Context, event ddd.Event) error {
	product := event.Payload().(*domain.Product)
	return h.publisher.Publish(ctx, productspb.ProductAggregateChannel,
		ddd.NewEvent(productspb.ProductRemovedEvent, &productspb.ProductRemoved{
			Id: product.ID(),
		}),
	)
}

func (h domainHandlers[T]) onProductThumbnailAdded(ctx context.Context, event ddd.Event) error {
	product := event.Payload().(*domain.Product)
	return h.publisher.Publish(ctx, productspb.ProductAggregateChannel,
		ddd.NewEvent(productspb.ProductThumbnailAddedEvent, &productspb.ProductRemoved{
			Id: product.ID(),
		}),
	)
}

func (h domainHandlers[T]) onVariantAdded(ctx context.Context, event ddd.Event) error {
	variant := event.Payload().(*domain.Variant)
	// Publish to productspb.VariantAggregateChannel or ProductAggregateChannel
	return h.publisher.Publish(ctx, productspb.VariantAggregateChannel,
		ddd.NewEvent(productspb.VariantAddedEvent, &productspb.VariantAdded{
			Id:           variant.ID(),
			ProductId:    variant.ProductID,
			Sku:          variant.SKU,
			Barcode:      variant.Barcode,
			VariantPrice: variant.VariantPrice,
			Stock:        variant.Stock,
			Weight:       variant.Weight,
			Height:       variant.Height,
			Width:        variant.Width,
			Depth:        variant.Depth,
			IsAvailable:  variant.IsAvailable,
			// Add other fields as needed, e.g. currencyCode, attributes, etc.
		}),
	)
}

func (h domainHandlers[T]) onVariantPriceIncreased(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.VariantPriceIncreased)
	return h.publisher.Publish(ctx, productspb.VariantAggregateChannel,
		ddd.NewEvent(productspb.VariantPriceIncreasedEvent, &productspb.VariantPriceIncreased{
			Id:       payload.VariantID,
			OldPrice: payload.OldPrice,
			NewPrice: payload.NewPrice,
		}),
	)
}

func (h domainHandlers[T]) onVariantPriceDecreased(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.VariantPriceDecreased)
	return h.publisher.Publish(ctx, productspb.VariantAggregateChannel,
		ddd.NewEvent(productspb.VariantPriceDecreasedEvent, &productspb.VariantPriceDecreased{
			Id:       payload.VariantID,
			OldPrice: payload.OldPrice,
			NewPrice: payload.NewPrice,
		}),
	)
}

func (h domainHandlers[T]) onVariantStockAdjusted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.VariantStockAdjusted)
	return h.publisher.Publish(ctx, productspb.VariantAggregateChannel,
		ddd.NewEvent(productspb.VariantStockAdjustedEvent, &productspb.VariantStockAdjusted{
			Id:       payload.VariantID,
			OldStock: payload.OldStock,
			NewStock: payload.NewStock,
		}),
	)
}

func (h domainHandlers[T]) onVariantArchived(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.VariantArchived)
	return h.publisher.Publish(ctx, productspb.VariantAggregateChannel,
		ddd.NewEvent(productspb.VariantArchivedEvent, &productspb.VariantArchived{
			Id: payload.VariantID,
		}),
	)
}

func (h domainHandlers[T]) onVariantRemoved(ctx context.Context, event ddd.Event) error {
	variant := event.Payload().(*domain.Variant)
	return h.publisher.Publish(ctx, productspb.VariantAggregateChannel,
		ddd.NewEvent(productspb.VariantRemovedEvent, &productspb.VariantRemoved{
			Id: variant.ID(),
		}),
	)
}
func domainOptionsToProto(opts []domain.Option) []*productspb.Option {
	pbOpts := make([]*productspb.Option, len(opts))
	for i, o := range opts {
		pbOpts[i] = &productspb.Option{
			Name:  o.Name,
			Value: o.Value,
			// If your proto uses int64, cast the float64 or int to int64:
			Price: int64(o.Price),
		}
	}
	return pbOpts
}
func domainAttributesToProto(attrs []domain.Attribute) []*productspb.Attribute {
	pbAttrs := make([]*productspb.Attribute, len(attrs))
	for i, a := range attrs {
		pbAttrs[i] = &productspb.Attribute{
			Key:   a.Key,
			Value: a.Value,
		}
	}
	return pbAttrs
}
