package handlers

import (
	"context"
	"middleman/internal/ddd"
	"middleman/internal/di"
	"middleman/tickets/internal/domain"
)

// catalogVariantHandlers listens to domain events for Variants and updates the "variants" table accordingly.
type catalogVariantHandlers[T ddd.Event] struct {
	catalog domain.CatalogVariantRepository
}

var _ ddd.EventHandler[ddd.Event] = (*catalogVariantHandlers[ddd.Event])(nil)

func NewCatalogVariantHandlers(catalog domain.CatalogVariantRepository) ddd.EventHandler[ddd.Event] {
	return catalogVariantHandlers[ddd.Event]{
		catalog: catalog,
	}
}

func RegisterCatalogVariantHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.VariantAddedEvent,
		domain.VariantPriceIncreasedEvent,
		domain.VariantPriceDecreasedEvent,
		domain.VariantStockAdjustedEvent,
		domain.VariantArchivedEvent,
		domain.VariantRemovedEvent,
		// ... plus any others you define, e.g. "VariantRebrandedEvent"
	)
}

func RegisterCatalogVariantHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		cvh := di.Get(ctx, "catalogVariantHandlers").(ddd.EventHandler[ddd.Event])
		return cvh.HandleEvent(ctx, event)
	})

	subscriber := container.Get("domainDispatcher").(*ddd.EventDispatcher[ddd.Event])
	RegisterCatalogVariantHandlers(subscriber, handlers)
}

func (h catalogVariantHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {

	case domain.VariantAddedEvent:
		return h.onVariantAdded(ctx, event)

	case domain.VariantPriceIncreasedEvent:
		return h.onVariantPriceIncreased(ctx, event)

	case domain.VariantPriceDecreasedEvent:
		return h.onVariantPriceDecreased(ctx, event)

	case domain.VariantStockAdjustedEvent:
		return h.onVariantStockAdjusted(ctx, event)

	case domain.VariantArchivedEvent:
		return h.onVariantArchived(ctx, event)

	case domain.VariantRemovedEvent:
		return h.onVariantRemoved(ctx, event)

		// case domain.VariantRebrandedEvent:
		//  return h.onVariantRebranded(ctx, event)

	}
	return nil
}

func (h catalogVariantHandlers[T]) onVariantAdded(ctx context.Context, event ddd.Event) error {
	agg := event.Payload().(*domain.Variant)
	return h.catalog.AddVariant(
		ctx,
		agg.ID(),
		agg.TicketID,
		agg.Status,
		agg.SKU,
		agg.Barcode,
		agg.Condition,
		agg.VariantPrice,
		agg.CurrencyCode,
		agg.Stock,
		agg.Weight,
		agg.Height,
		agg.Width,
		agg.Depth,
		agg.Attributes,
		agg.IsAvailable,
		agg.HasOptions,
		agg.Options,
	)
}

func (h catalogVariantHandlers[T]) onVariantPriceIncreased(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.VariantPriceIncreased)
	return h.catalog.UpdateVariantPrice(ctx, payload.VariantID, payload.OldPrice, payload.NewPrice)
}

func (h catalogVariantHandlers[T]) onVariantPriceDecreased(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.VariantPriceDecreased)
	return h.catalog.UpdateVariantPrice(ctx, payload.VariantID, payload.OldPrice, payload.NewPrice)
}

func (h catalogVariantHandlers[T]) onVariantStockAdjusted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.VariantStockAdjusted)
	return h.catalog.AdjustVariantStock(ctx, payload.VariantID, payload.OldStock, payload.NewStock)
}

func (h catalogVariantHandlers[T]) onVariantArchived(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.VariantArchived)
	return h.catalog.ArchiveVariant(ctx, payload.VariantID)
}

func (h catalogVariantHandlers[T]) onVariantRemoved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.VariantRemoved)
	// If removing requires userSellerID, you need it in the event
	return h.catalog.RemoveVariant(ctx, payload.VariantID, "") // e.g. pass userSellerID if you have it
}

// Example if you wanted variant rebranding event
// func (h catalogVariantHandlers[T]) onVariantRebranded(ctx context.Context, event ddd.Event) error {
//   payload := event.Payload().(*domain.VariantRebranded)
//   return h.catalog.RebrandVariant(ctx, payload.VariantID, payload.NewName, payload.NewDescription)
// }
