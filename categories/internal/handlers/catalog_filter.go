package handlers

import (
	"context"
	"middleman/categories/internal/domain"
	"middleman/internal/ddd"
	"middleman/internal/di"
)

// catalogFilterHandlers listens to domain events for Filters and updates the "filters" table accordingly.
type catalogFilterHandlers[T ddd.Event] struct {
	catalog domain.CatalogFilterRepository
}

var _ ddd.EventHandler[ddd.Event] = (*catalogFilterHandlers[ddd.Event])(nil)

func NewCatalogFilterHandlers(catalog domain.CatalogFilterRepository) ddd.EventHandler[ddd.Event] {
	return catalogFilterHandlers[ddd.Event]{
		catalog: catalog,
	}
}

func RegisterCatalogFilterHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.FilterAddedEvent,
		domain.FilterArchivedEvent,
		domain.FilterRemovedEvent,
		// ... plus any others you define, e.g. "FilterRebrandedEvent"
	)
}

func RegisterCatalogFilterHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		cvh := di.Get(ctx, "catalogFilterHandlers").(ddd.EventHandler[ddd.Event])
		return cvh.HandleEvent(ctx, event)
	})

	subscriber := container.Get("domainDispatcher").(*ddd.EventDispatcher[ddd.Event])
	RegisterCatalogFilterHandlers(subscriber, handlers)
}

func (h catalogFilterHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {

	case domain.FilterAddedEvent:
		return h.onFilterAdded(ctx, event)

	case domain.FilterArchivedEvent:
		return h.onFilterArchived(ctx, event)

	case domain.FilterRemovedEvent:
		return h.onFilterRemoved(ctx, event)

		// case domain.FilterRebrandedEvent:
		//  return h.onFilterRebranded(ctx, event)

	}
	return nil
}

func (h catalogFilterHandlers[T]) onFilterAdded(ctx context.Context, event ddd.Event) error {
	agg := event.Payload().(*domain.Filter)
	return h.catalog.AddFilter(
		ctx,
		agg.ID(),
		agg.CategoryID, agg.Name, agg.FilterType, agg.Values, agg.IsActive,
	)
}

func (h catalogFilterHandlers[T]) onFilterArchived(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.FilterArchived)
	return h.catalog.ArchiveFilter(ctx, payload.FilterID)
}

func (h catalogFilterHandlers[T]) onFilterRemoved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.FilterRemoved)
	// If removing requires userSellerID, you need it in the event
	return h.catalog.RemoveFilter(ctx, payload.FilterID, "") // e.g. pass userSellerID if you have it
}

// Example if you wanted filter rebranding event
// func (h catalogFilterHandlers[T]) onFilterRebranded(ctx context.Context, event ddd.Event) error {
//   payload := event.Payload().(*domain.FilterRebranded)
//   return h.catalog.RebrandFilter(ctx, payload.FilterID, payload.NewName, payload.NewDescription)
// }
