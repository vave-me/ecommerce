package handlers

import (
	"context"

	"middleman/internal/ddd"
	"middleman/internal/di"
	"middleman/streams/internal/domain"
)

// catalogHandlers listens to domain events and updates the Catalog DB
type catalogHandlers[T ddd.Event] struct {
	catalog domain.CatalogRepository
}

var _ ddd.EventHandler[ddd.Event] = (*catalogHandlers[ddd.Event])(nil)

func NewCatalogHandlers(catalog domain.CatalogRepository) ddd.EventHandler[ddd.Event] {
	return catalogHandlers[ddd.Event]{catalog: catalog}
}

func RegisterCatalogHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.StreamAddedEvent,
		domain.StreamUpdatedEvent,
		domain.StreamRebrandedEvent,
		domain.StreamPriceIncreasedEvent,
		domain.StreamPriceDecreasedEvent,
		domain.StreamRemovedEvent,
		domain.StreamArchivedEvent,
		domain.StreamStockAdjustedEvent,
		domain.StreamSoldEvent,
		domain.StreamLeasedEvent,
		domain.StreamPawnedEvent,
		domain.StreamNegotiableToggledEvent,
		domain.StreamThumbnailAddedEvent,
		domain.StreamThumbnailRemovedEvent,
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

	case domain.StreamAddedEvent:
		return h.onStreamAdded(ctx, event)
	case domain.StreamRebrandedEvent:
		return h.onStreamRebranded(ctx, event)
	case domain.StreamUpdatedEvent:
		return h.onStreamUpdated(ctx, event)

	case domain.StreamPriceIncreasedEvent:
		return h.onStreamPriceIncreased(ctx, event)

	case domain.StreamPriceDecreasedEvent:
		return h.onStreamPriceDecreased(ctx, event)

	case domain.StreamStockAdjustedEvent:
		return h.onStreamStockAdjusted(ctx, event)

	case domain.StreamArchivedEvent:
		return h.onStreamArchived(ctx, event)

	case domain.StreamSoldEvent:
		return h.onStreamSold(ctx, event)

	case domain.StreamLeasedEvent:
		return h.onStreamLeased(ctx, event)

	case domain.StreamPawnedEvent:
		return h.onStreamPawned(ctx, event)

	case domain.StreamNegotiableToggledEvent:
		return h.onStreamNegotiableToggled(ctx, event)

	case domain.StreamRemovedEvent:
		return h.onStreamRemoved(ctx, event)
	case domain.StreamThumbnailAddedEvent:
		return h.onStreamThumbnailAdded(ctx, event)
	case domain.StreamThumbnailUpdatedEvent:
		return h.onStreamThumbnailUpdated(ctx, event)

	}

	return nil
}

// ---------------------------------------------------------------------------
// HANDLER IMPLEMENTATIONS
// ---------------------------------------------------------------------------
// TODO Middleman service bool or sting with couple strings as a diffrent services
func (h catalogHandlers[T]) onStreamAdded(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.Stream) // domain.Stream after “added”
	return h.catalog.AddStream(ctx,
		e.ID(),
		e.Name,
		e.Description,
		e.BasePrice,
		e.UserSellerID,
		e.CategoryID,
		e.CategorySlug,
		e.Brand,
		e.Condition,
		e.Model,
		e.Tags,
		e.ManageStock,
		e.Stock,
		e.SKU,
		e.Attributes,
		e.Weight,
		e.Height,
		e.Width,
		e.Depth,
		e.Status,
		e.Negotiable,
		e.UserType,
		e.MiddlemanService,
		e.ShippingCost,
		e.HasVariants,
		e.Options,
		e.Thumbnail,
		e.Lat,
		e.Lng,
	)
}
func (h catalogHandlers[T]) onStreamUpdated(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.Stream) // domain.Stream after “added”
	return h.catalog.UpdateStream(ctx,
		e.ID(),
		e.Name,
		e.Description,
		e.BasePrice,
		e.CategoryID,
		e.CategorySlug,
		e.Brand,
		e.Condition,
		e.Model,
		e.Tags,
		e.ManageStock,
		e.Stock,
		e.SKU,
		e.Attributes,
		e.Weight,
		e.Height,
		e.Width,
		e.Depth,
		e.Status,
		e.Negotiable,
		e.UserType,
		e.MiddlemanService,
		e.ShippingCost,
		e.HasVariants,
		e.Options,
		e.Thumbnail,
		e.Lat,
		e.Lng,
	)
}
func (h catalogHandlers[T]) onStreamRebranded(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.Stream)
	return h.catalog.RebrandStream(
		ctx,
		e.ID(),
		e.Name,
		e.Description,
		e.CategoryID,
		e.CategorySlug,
		e.Brand,
		e.Model,
		string(e.Condition),
		e.Tags,
	)
}

func (h catalogHandlers[T]) onStreamPriceIncreased(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.StreamPriceIncreased)
	// If your domain event has OldPrice/NewPrice, do:
	return h.catalog.UpdatePrice(ctx, e.StreamID, e.OldPrice, e.NewPrice)
}

func (h catalogHandlers[T]) onStreamPriceDecreased(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.StreamPriceDecreased)
	return h.catalog.UpdatePrice(ctx, e.StreamID, e.OldPrice, e.NewPrice)
}

func (h catalogHandlers[T]) onStreamStockAdjusted(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.StreamStockAdjusted)
	return h.catalog.AdjustStock(ctx, e.StreamID, "", e.OldStock, e.NewStock)
}

func (h catalogHandlers[T]) onStreamArchived(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.StreamArchived)
	return h.catalog.ArchiveStream(ctx, e.StreamID, e.UserSellerID)
}

func (h catalogHandlers[T]) onStreamSold(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.StreamSold)
	return h.catalog.MarkStreamSold(ctx, e.StreamID, e.UserSellerID, e.FinalPrice)
}

func (h catalogHandlers[T]) onStreamLeased(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.StreamLeased)
	// store monthlyPrice, leaseTerm if you'd like
	return h.catalog.MarkStreamLeased(ctx, e.StreamID, e.UserSellerID)
}

func (h catalogHandlers[T]) onStreamPawned(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.StreamPawned)
	return h.catalog.MarkStreamPawned(ctx, e.StreamID, e.UserSellerID)
}

func (h catalogHandlers[T]) onStreamNegotiableToggled(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.StreamNegotiableToggled)
	return h.catalog.ToggleNegotiable(ctx, e.StreamID, "", e.OldValue) // "userSellerID" left out if not used
}

func (h catalogHandlers[T]) onStreamRemoved(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.StreamRemoved)
	return h.catalog.RemoveStream(ctx, e.StreamID, e.UserSellerID)
}

func (h catalogHandlers[T]) onStreamThumbnailAdded(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.StreamThumbnailAdded)
	return h.catalog.UpdateThumbnail(ctx, e.StreamID, e.Thumbnail)
}
func (h catalogHandlers[T]) onStreamThumbnailUpdated(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.StreamThumbnailUpdated)
	return h.catalog.UpdateThumbnail(ctx, e.StreamID, e.Thumbnail)
}

// a small helper if you want to store "premium"/"base" as bool. Adjust as needed.
func strToBool(s string) bool {
	return s == "true" || s == "premium"
}
