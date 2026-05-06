package handlers

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/streams/internal/domain"
	"middleman/streams/streamspb"

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
		domain.StreamAddedEvent,
		domain.StreamUpdatedEvent,
		domain.StreamRebrandedEvent,
		domain.StreamPriceIncreasedEvent,
		domain.StreamPriceDecreasedEvent,
		domain.StreamStockAdjustedEvent,
		domain.StreamRemovedEvent,
		domain.StreamArchivedEvent,
		domain.StreamSoldEvent,
		domain.StreamLeasedEvent,
		domain.StreamPawnedEvent,
		domain.StreamThumbnailAddedEvent,
		domain.StreamThumbnailUpdatedEvent,

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

	case domain.StreamAddedEvent:
		return h.onStreamAdded(ctx, event)
	case domain.StreamUpdatedEvent:
		return h.onStreamUpdated(ctx, event)
	case domain.StreamRebrandedEvent:
		return h.onStreamRebranded(ctx, event)
	case domain.StreamPriceIncreasedEvent:
		return h.onStreamPriceIncreased(ctx, event)
	case domain.StreamPriceDecreasedEvent:
		return h.onStreamPriceDecreased(ctx, event)
	case domain.StreamStockAdjustedEvent:
		return h.onStreamStockAdjusted(ctx, event)
	case domain.StreamArchivedEvent:
		return h.onStreamArchived(ctx, event)
	case domain.StreamRemovedEvent:
		return h.onStreamRemoved(ctx, event)
	case domain.StreamSoldEvent:
		return h.onStreamSold(ctx, event)
	case domain.StreamLeasedEvent:
		return h.onStreamLeased(ctx, event)
	case domain.StreamPawnedEvent:
		return h.onStreamPawned(ctx, event)

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

func (h domainHandlers[T]) onStreamAdded(ctx context.Context, event ddd.Event) error {
	product := event.Payload().(*domain.Stream)
	return h.publisher.Publish(ctx, streamspb.StreamAggregateChannel,
		ddd.NewEvent(streamspb.StreamAddedEvent, &streamspb.StreamAdded{
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
func (h domainHandlers[T]) onStreamUpdated(ctx context.Context, event ddd.Event) error {
	product := event.Payload().(*domain.Stream)
	return h.publisher.Publish(ctx, streamspb.StreamAggregateChannel,
		ddd.NewEvent(streamspb.StreamUpdatedEvent, &streamspb.StreamUpdated{
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

func (h domainHandlers[T]) onStreamRebranded(ctx context.Context, event ddd.Event) error {
	product := event.Payload().(*domain.Stream)
	return h.publisher.Publish(ctx, streamspb.StreamAggregateChannel,
		ddd.NewEvent(streamspb.StreamRebrandedEvent, &streamspb.StreamRebranded{
			Id:          product.ID(),
			Name:        product.Name,
			Description: product.Description,
		}),
	)
}

func (h domainHandlers[T]) onStreamPriceIncreased(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.StreamPriceIncreased)
	return h.publisher.Publish(ctx, streamspb.StreamAggregateChannel,
		ddd.NewEvent(streamspb.StreamPriceIncreasedEvent, &streamspb.StreamPriceIncreased{
			Id:       payload.StreamID,
			OldPrice: payload.OldPrice,
			NewPrice: payload.NewPrice,
		}),
	)
}

func (h domainHandlers[T]) onStreamArchived(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.StreamArchived)
	return h.publisher.Publish(ctx, streamspb.StreamAggregateChannel,
		ddd.NewEvent(streamspb.StreamArchivedEvent, &streamspb.StreamArchived{
			Id:           payload.StreamID,
			UserSellerId: payload.UserSellerID,
		}),
	)
}

func (h domainHandlers[T]) onStreamSold(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.StreamSold)
	return h.publisher.Publish(ctx, streamspb.StreamAggregateChannel,
		ddd.NewEvent(streamspb.StreamSoldEvent, &streamspb.StreamSold{
			Id:           payload.StreamID,
			UserSellerId: payload.UserSellerID,

			// BuyerId:   payload.BuyerID, if your proto has it
		}),
	)
}

func (h domainHandlers[T]) onStreamLeased(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.StreamLeased)
	return h.publisher.Publish(ctx, streamspb.StreamAggregateChannel,
		ddd.NewEvent(streamspb.StreamLeasedEvent, &streamspb.StreamLeased{
			Id:           payload.StreamID,
			UserSellerId: payload.UserSellerID,

			// LesseeId: payload.LesseeID, if needed
		}),
	)
}

func (h domainHandlers[T]) onStreamPriceDecreased(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.StreamPriceDecreased)
	return h.publisher.Publish(ctx, streamspb.StreamAggregateChannel,
		ddd.NewEvent(streamspb.StreamPriceDecreasedEvent, &streamspb.StreamPriceDecreased{
			Id:       payload.StreamID,
			OldPrice: payload.OldPrice,
			NewPrice: payload.NewPrice,
		}),
	)
}
func (h domainHandlers[T]) onStreamPawned(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.StreamPawned)
	return h.publisher.Publish(ctx, streamspb.StreamAggregateChannel,
		ddd.NewEvent(streamspb.StreamPawnedEvent, &streamspb.StreamPawned{
			Id:           payload.StreamID,
			UserSellerId: payload.UserSellerID,
		}),
	)
}

func (h domainHandlers[T]) onStreamStockAdjusted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.StreamStockAdjusted)
	return h.publisher.Publish(ctx, streamspb.StreamAggregateChannel,
		ddd.NewEvent(streamspb.StreamStockAdjustedEvent, &streamspb.StreamStockAdjusted{
			Id:       payload.StreamID,
			OldStock: payload.OldStock,
			NewStock: payload.NewStock,
		}),
	)
}
func (h domainHandlers[T]) onStreamRemoved(ctx context.Context, event ddd.Event) error {
	product := event.Payload().(*domain.Stream)
	return h.publisher.Publish(ctx, streamspb.StreamAggregateChannel,
		ddd.NewEvent(streamspb.StreamRemovedEvent, &streamspb.StreamRemoved{
			Id: product.ID(),
		}),
	)
}

func (h domainHandlers[T]) onStreamThumbnailAdded(ctx context.Context, event ddd.Event) error {
	product := event.Payload().(*domain.Stream)
	return h.publisher.Publish(ctx, streamspb.StreamAggregateChannel,
		ddd.NewEvent(streamspb.StreamThumbnailAddedEvent, &streamspb.StreamRemoved{
			Id: product.ID(),
		}),
	)
}

func (h domainHandlers[T]) onVariantAdded(ctx context.Context, event ddd.Event) error {
	variant := event.Payload().(*domain.Variant)
	// Publish to streamspb.VariantAggregateChannel or StreamAggregateChannel
	return h.publisher.Publish(ctx, streamspb.VariantAggregateChannel,
		ddd.NewEvent(streamspb.VariantAddedEvent, &streamspb.VariantAdded{
			Id:           variant.ID(),
			StreamId:     variant.StreamID,
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
	return h.publisher.Publish(ctx, streamspb.VariantAggregateChannel,
		ddd.NewEvent(streamspb.VariantPriceIncreasedEvent, &streamspb.VariantPriceIncreased{
			Id:       payload.VariantID,
			OldPrice: payload.OldPrice,
			NewPrice: payload.NewPrice,
		}),
	)
}

func (h domainHandlers[T]) onVariantPriceDecreased(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.VariantPriceDecreased)
	return h.publisher.Publish(ctx, streamspb.VariantAggregateChannel,
		ddd.NewEvent(streamspb.VariantPriceDecreasedEvent, &streamspb.VariantPriceDecreased{
			Id:       payload.VariantID,
			OldPrice: payload.OldPrice,
			NewPrice: payload.NewPrice,
		}),
	)
}

func (h domainHandlers[T]) onVariantStockAdjusted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.VariantStockAdjusted)
	return h.publisher.Publish(ctx, streamspb.VariantAggregateChannel,
		ddd.NewEvent(streamspb.VariantStockAdjustedEvent, &streamspb.VariantStockAdjusted{
			Id:       payload.VariantID,
			OldStock: payload.OldStock,
			NewStock: payload.NewStock,
		}),
	)
}

func (h domainHandlers[T]) onVariantArchived(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.VariantArchived)
	return h.publisher.Publish(ctx, streamspb.VariantAggregateChannel,
		ddd.NewEvent(streamspb.VariantArchivedEvent, &streamspb.VariantArchived{
			Id: payload.VariantID,
		}),
	)
}

func (h domainHandlers[T]) onVariantRemoved(ctx context.Context, event ddd.Event) error {
	variant := event.Payload().(*domain.Variant)
	return h.publisher.Publish(ctx, streamspb.VariantAggregateChannel,
		ddd.NewEvent(streamspb.VariantRemovedEvent, &streamspb.VariantRemoved{
			Id: variant.ID(),
		}),
	)
}
func domainOptionsToProto(opts []domain.Option) []*streamspb.Option {
	pbOpts := make([]*streamspb.Option, len(opts))
	for i, o := range opts {
		pbOpts[i] = &streamspb.Option{
			Name:  o.Name,
			Value: o.Value,
			// If your proto uses int64, cast the float64 or int to int64:
			Price: int64(o.Price),
		}
	}
	return pbOpts
}
func domainAttributesToProto(attrs []domain.Attribute) []*streamspb.Attribute {
	pbAttrs := make([]*streamspb.Attribute, len(attrs))
	for i, a := range attrs {
		pbAttrs[i] = &streamspb.Attribute{
			Key:   a.Key,
			Value: a.Value,
		}
	}
	return pbAttrs
}
