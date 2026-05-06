package handlers

import (
	"context"

	"middleman/internal/ddd"
	"middleman/internal/di"
	"middleman/tickets/internal/domain"
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
		domain.TicketAddedEvent,
		domain.TicketUpdatedEvent,
		domain.TicketRebrandedEvent,
		domain.TicketPriceIncreasedEvent,
		domain.TicketPriceDecreasedEvent,
		domain.TicketRemovedEvent,
		domain.TicketArchivedEvent,
		domain.TicketStockAdjustedEvent,
		domain.TicketSoldEvent,
		domain.TicketLeasedEvent,
		domain.TicketPawnedEvent,
		domain.TicketNegotiableToggledEvent,
		domain.TicketThumbnailAddedEvent,
		domain.TicketThumbnailRemovedEvent,
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

	case domain.TicketAddedEvent:
		return h.onTicketAdded(ctx, event)
	case domain.TicketRebrandedEvent:
		return h.onTicketRebranded(ctx, event)
	case domain.TicketUpdatedEvent:
		return h.onTicketUpdated(ctx, event)

	case domain.TicketPriceIncreasedEvent:
		return h.onTicketPriceIncreased(ctx, event)

	case domain.TicketPriceDecreasedEvent:
		return h.onTicketPriceDecreased(ctx, event)

	case domain.TicketStockAdjustedEvent:
		return h.onTicketStockAdjusted(ctx, event)

	case domain.TicketArchivedEvent:
		return h.onTicketArchived(ctx, event)

	case domain.TicketSoldEvent:
		return h.onTicketSold(ctx, event)

	case domain.TicketLeasedEvent:
		return h.onTicketLeased(ctx, event)

	case domain.TicketPawnedEvent:
		return h.onTicketPawned(ctx, event)

	case domain.TicketNegotiableToggledEvent:
		return h.onTicketNegotiableToggled(ctx, event)

	case domain.TicketRemovedEvent:
		return h.onTicketRemoved(ctx, event)
	case domain.TicketThumbnailAddedEvent:
		return h.onTicketThumbnailAdded(ctx, event)
	case domain.TicketThumbnailUpdatedEvent:
		return h.onTicketThumbnailUpdated(ctx, event)

	}

	return nil
}

// ---------------------------------------------------------------------------
// HANDLER IMPLEMENTATIONS
// ---------------------------------------------------------------------------
// TODO Middleman service bool or sting with couple strings as a diffrent services
func (h catalogHandlers[T]) onTicketAdded(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.Ticket) // domain.Ticket after “added”
	return h.catalog.AddTicket(ctx,
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
func (h catalogHandlers[T]) onTicketUpdated(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.Ticket) // domain.Ticket after “added”
	return h.catalog.UpdateTicket(ctx,
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
func (h catalogHandlers[T]) onTicketRebranded(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.Ticket)
	return h.catalog.RebrandTicket(
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

func (h catalogHandlers[T]) onTicketPriceIncreased(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.TicketPriceIncreased)
	// If your domain event has OldPrice/NewPrice, do:
	return h.catalog.UpdatePrice(ctx, e.TicketID, e.OldPrice, e.NewPrice)
}

func (h catalogHandlers[T]) onTicketPriceDecreased(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.TicketPriceDecreased)
	return h.catalog.UpdatePrice(ctx, e.TicketID, e.OldPrice, e.NewPrice)
}

func (h catalogHandlers[T]) onTicketStockAdjusted(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.TicketStockAdjusted)
	return h.catalog.AdjustStock(ctx, e.TicketID, "", e.OldStock, e.NewStock)
}

func (h catalogHandlers[T]) onTicketArchived(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.TicketArchived)
	return h.catalog.ArchiveTicket(ctx, e.TicketID, e.UserSellerID)
}

func (h catalogHandlers[T]) onTicketSold(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.TicketSold)
	return h.catalog.MarkTicketSold(ctx, e.TicketID, e.UserSellerID, e.FinalPrice)
}

func (h catalogHandlers[T]) onTicketLeased(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.TicketLeased)
	// store monthlyPrice, leaseTerm if you'd like
	return h.catalog.MarkTicketLeased(ctx, e.TicketID, e.UserSellerID)
}

func (h catalogHandlers[T]) onTicketPawned(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.TicketPawned)
	return h.catalog.MarkTicketPawned(ctx, e.TicketID, e.UserSellerID)
}

func (h catalogHandlers[T]) onTicketNegotiableToggled(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.TicketNegotiableToggled)
	return h.catalog.ToggleNegotiable(ctx, e.TicketID, "", e.OldValue) // "userSellerID" left out if not used
}

func (h catalogHandlers[T]) onTicketRemoved(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.TicketRemoved)
	return h.catalog.RemoveTicket(ctx, e.TicketID, e.UserSellerID)
}

func (h catalogHandlers[T]) onTicketThumbnailAdded(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.TicketThumbnailAdded)
	return h.catalog.UpdateThumbnail(ctx, e.TicketID, e.Thumbnail)
}
func (h catalogHandlers[T]) onTicketThumbnailUpdated(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.TicketThumbnailUpdated)
	return h.catalog.UpdateThumbnail(ctx, e.TicketID, e.Thumbnail)
}

// a small helper if you want to store "premium"/"base" as bool. Adjust as needed.
func strToBool(s string) bool {
	return s == "true" || s == "premium"
}
