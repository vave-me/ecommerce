package handlers

import (
	"context"

	"middleman/internal/ddd"
	"middleman/internal/di"
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
		domain.ProductAddedEvent,
		domain.ProductUpdatedEvent,
		domain.ProductRebrandedEvent,
		domain.ProductPriceIncreasedEvent,
		domain.ProductPriceDecreasedEvent,
		domain.ProductRemovedEvent,
		domain.ProductArchivedEvent,
		domain.ProductStockAdjustedEvent,
		domain.ProductSoldEvent,
		domain.ProductLeasedEvent,
		domain.ProductPawnedEvent,
		domain.ProductNegotiableToggledEvent,
		domain.ProductThumbnailAddedEvent,
		domain.ProductThumbnailRemovedEvent,
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

	case domain.ProductAddedEvent:
		return h.onProductAdded(ctx, event)
	case domain.ProductRebrandedEvent:
		return h.onProductRebranded(ctx, event)
	case domain.ProductUpdatedEvent:
		return h.onProductUpdated(ctx, event)

	case domain.ProductPriceIncreasedEvent:
		return h.onProductPriceIncreased(ctx, event)

	case domain.ProductPriceDecreasedEvent:
		return h.onProductPriceDecreased(ctx, event)

	case domain.ProductStockAdjustedEvent:
		return h.onProductStockAdjusted(ctx, event)

	case domain.ProductArchivedEvent:
		return h.onProductArchived(ctx, event)

	case domain.ProductSoldEvent:
		return h.onProductSold(ctx, event)

	case domain.ProductLeasedEvent:
		return h.onProductLeased(ctx, event)

	case domain.ProductPawnedEvent:
		return h.onProductPawned(ctx, event)

	case domain.ProductNegotiableToggledEvent:
		return h.onProductNegotiableToggled(ctx, event)

	case domain.ProductRemovedEvent:
		return h.onProductRemoved(ctx, event)
	case domain.ProductThumbnailAddedEvent:
		return h.onProductThumbnailAdded(ctx, event)
	case domain.ProductThumbnailUpdatedEvent:
		return h.onProductThumbnailUpdated(ctx, event)

	}

	return nil
}

// ---------------------------------------------------------------------------
// HANDLER IMPLEMENTATIONS
// ---------------------------------------------------------------------------
// TODO Middleman service bool or sting with couple strings as a diffrent services
func (h catalogHandlers[T]) onProductAdded(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.Product) // domain.Product after “added”
	return h.catalog.AddProduct(ctx,
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
func (h catalogHandlers[T]) onProductUpdated(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.Product) // domain.Product after “added”
	return h.catalog.UpdateProduct(ctx,
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
func (h catalogHandlers[T]) onProductRebranded(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.Product)
	return h.catalog.RebrandProduct(
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

func (h catalogHandlers[T]) onProductPriceIncreased(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.ProductPriceIncreased)
	// If your domain event has OldPrice/NewPrice, do:
	return h.catalog.UpdatePrice(ctx, e.ProductID, e.OldPrice, e.NewPrice)
}

func (h catalogHandlers[T]) onProductPriceDecreased(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.ProductPriceDecreased)
	return h.catalog.UpdatePrice(ctx, e.ProductID, e.OldPrice, e.NewPrice)
}

func (h catalogHandlers[T]) onProductStockAdjusted(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.ProductStockAdjusted)
	return h.catalog.AdjustStock(ctx, e.ProductID, "", e.OldStock, e.NewStock)
}

func (h catalogHandlers[T]) onProductArchived(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.ProductArchived)
	return h.catalog.ArchiveProduct(ctx, e.ProductID, e.UserSellerID)
}

func (h catalogHandlers[T]) onProductSold(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.ProductSold)
	return h.catalog.MarkProductSold(ctx, e.ProductID, e.UserSellerID, e.FinalPrice)
}

func (h catalogHandlers[T]) onProductLeased(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.ProductLeased)
	// store monthlyPrice, leaseTerm if you'd like
	return h.catalog.MarkProductLeased(ctx, e.ProductID, e.UserSellerID)
}

func (h catalogHandlers[T]) onProductPawned(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.ProductPawned)
	return h.catalog.MarkProductPawned(ctx, e.ProductID, e.UserSellerID)
}

func (h catalogHandlers[T]) onProductNegotiableToggled(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.ProductNegotiableToggled)
	return h.catalog.ToggleNegotiable(ctx, e.ProductID, "", e.OldValue) // "userSellerID" left out if not used
}

func (h catalogHandlers[T]) onProductRemoved(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.ProductRemoved)
	return h.catalog.RemoveProduct(ctx, e.ProductID, e.UserSellerID)
}

func (h catalogHandlers[T]) onProductThumbnailAdded(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.ProductThumbnailAdded)
	return h.catalog.UpdateThumbnail(ctx, e.ProductID, e.Thumbnail)
}
func (h catalogHandlers[T]) onProductThumbnailUpdated(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.ProductThumbnailUpdated)
	return h.catalog.UpdateThumbnail(ctx, e.ProductID, e.Thumbnail)
}

// a small helper if you want to store "premium"/"base" as bool. Adjust as needed.
func strToBool(s string) bool {
	return s == "true" || s == "premium"
}
