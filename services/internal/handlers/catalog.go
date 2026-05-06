package handlers

import (
	"context"

	"middleman/internal/ddd"
	"middleman/internal/di"
	"middleman/services/internal/domain"
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
		domain.ServiceAddedEvent,
		domain.ServiceUpdatedEvent,
		domain.ServiceRebrandedEvent,
		domain.ServicePriceIncreasedEvent,
		domain.ServicePriceDecreasedEvent,
		domain.ServiceRemovedEvent,
		domain.ServiceArchivedEvent,
		domain.ServiceStockAdjustedEvent,
		domain.ServiceSoldEvent,
		domain.ServiceLeasedEvent,
		domain.ServicePawnedEvent,
		domain.ServiceNegotiableToggledEvent,
		// plus anything else you'd like
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

	case domain.ServiceAddedEvent:
		return h.onServiceAdded(ctx, event)
	case domain.ServiceUpdatedEvent:
		return h.onServiceUpdated(ctx, event)

	case domain.ServiceRebrandedEvent:
		return h.onServiceRebranded(ctx, event)

	case domain.ServicePriceIncreasedEvent:
		return h.onServicePriceIncreased(ctx, event)

	case domain.ServicePriceDecreasedEvent:
		return h.onServicePriceDecreased(ctx, event)

	case domain.ServiceArchivedEvent:
		return h.onServiceArchived(ctx, event)

	case domain.ServiceSoldEvent:
		return h.onServiceSold(ctx, event)

	case domain.ServiceLeasedEvent:
		return h.onServiceLeased(ctx, event)

	case domain.ServiceNegotiableToggledEvent:
		return h.onServiceNegotiableToggled(ctx, event)

	case domain.ServiceRemovedEvent:
		return h.onServiceRemoved(ctx, event)
	}
	return nil
}

// ---------------------------------------------------------------------------
// HANDLER IMPLEMENTATIONS
// ---------------------------------------------------------------------------
// TODO Middleman service bool or sting with couple strings as a diffrent services
func (h catalogHandlers[T]) onServiceAdded(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.Service) // domain.Service after “added”
	return h.catalog.AddService(ctx,
		e.ID(),
		e.Name,
		e.Description,
		e.ServiceType,
		e.BasePrice,
		e.Pricing,
		e.Availability,
		e.ProviderName,
		e.UserID,
		e.CategoryID,
		e.CategorySlug,
		e.DescriptionShort,
		e.DescriptionLong,
		e.Qualifications,
		e.Contact,
		e.Faq,
		e.Tags,
		e.Status,
		e.UserType,
		e.ShippingCost,
		e.Negotiable,
		e.HasVariants,
		e.MiddlemanService,
		e.Attributes,
		e.Options,
		e.Thumbnail,
		e.Lat,
		e.Lng,
	)
}

func (h catalogHandlers[T]) onServiceUpdated(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.Service) // domain.Service after “added”
	return h.catalog.UpdateService(ctx,
		e.ID(),
		e.Name,
		e.Description,
		e.ServiceType,
		e.BasePrice,
		e.Pricing,
		e.Availability,
		e.ProviderName,
		e.UserID,
		e.CategoryID,
		e.CategorySlug,
		e.DescriptionShort,
		e.DescriptionLong,
		e.Qualifications,
		e.Contact,
		e.Faq,
		e.Tags,
		e.Status,
		e.UserType,
		e.ShippingCost,
		e.Negotiable,
		e.HasVariants,
		e.MiddlemanService,
		e.Attributes,
		e.Options,
		e.Thumbnail,
		e.Lat,
		e.Lng,
	)
}

func (h catalogHandlers[T]) onServiceRebranded(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.Service)
	return h.catalog.RebrandService(
		ctx,
		e.ID(),
		e.Name,
		e.Description,
		e.Tags,
		e.Qualifications, e.Faq,
	)
}

func (h catalogHandlers[T]) onServicePriceIncreased(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.ServicePriceIncreased)
	// If your domain event has OldPrice/NewPrice, do:
	return h.catalog.UpdatePrice(ctx, e.ServiceID, e.OldPrice, e.NewPrice)
}

func (h catalogHandlers[T]) onServicePriceDecreased(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.ServicePriceDecreased)
	return h.catalog.UpdatePrice(ctx, e.ServiceID, e.OldPrice, e.NewPrice)
}

func (h catalogHandlers[T]) onServiceArchived(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.ServiceArchived)
	return h.catalog.ArchiveService(ctx, e.ServiceID, e.UserID)
}

func (h catalogHandlers[T]) onServiceSold(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.ServiceSold)
	return h.catalog.MarkServiceSold(ctx, e.ServiceID, e.UserID, e.FinalPrice)
}

func (h catalogHandlers[T]) onServiceLeased(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.ServiceLeased)
	// store monthlyPrice, leaseTerm if you'd like
	return h.catalog.MarkServiceLeased(ctx, e.ServiceID, e.UserID)
}

func (h catalogHandlers[T]) onServiceNegotiableToggled(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.ServiceNegotiableToggled)
	return h.catalog.ToggleNegotiable(ctx, e.ServiceID, "", e.OldValue) // "userSellerID" left out if not used
}

func (h catalogHandlers[T]) onServiceRemoved(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.ServiceRemoved)
	return h.catalog.RemoveService(ctx, e.ServiceID, e.UserID)
}

func strToBool(s string) bool {
	return s == "true" || s == "premium"
}
