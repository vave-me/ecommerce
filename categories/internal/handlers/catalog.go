package handlers

import (
	"context"

	"middleman/categories/internal/domain"
	"middleman/internal/ddd"
	"middleman/internal/di"
)

// catalogHandlers listens to domain events about Categories and updates the Catalog DB.
type catalogHandlers[T ddd.Event] struct {
	catalog domain.CatalogRepository
}

// compile-time interface check
var _ ddd.EventHandler[ddd.Event] = (*catalogHandlers[ddd.Event])(nil)

// NewCatalogHandlers constructs the event handler for Category-related domain events.
func NewCatalogHandlers(catalog domain.CatalogRepository) ddd.EventHandler[ddd.Event] {
	return &catalogHandlers[ddd.Event]{catalog: catalog}
}

// RegisterCatalogHandlers wires up these handlers for the domain's "Category" events.
func RegisterCatalogHandlers(
	subscriber ddd.EventSubscriber[ddd.Event],
	handlers ddd.EventHandler[ddd.Event],
) {
	subscriber.Subscribe(handlers,
		domain.CategoryAddedEvent,
		domain.CategoryUpdatedEvent, // optional if you have an 'Update' event
		domain.CategoryRebrandedEvent,
		domain.CategoryArchivedEvent,
		domain.CategoryRemovedEvent,
	)
}

// RegisterCatalogHandlersTx is a convenience function if you want to resolve
// the handler from a DI container and subscribe it with a transaction approach.
func RegisterCatalogHandlersTx(container di.Container) {
	// We wrap the actual handler in a function that gets the "catalogHandlers"
	// from the container each time
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		catalogHandlers := di.Get(ctx, "catalogHandlers").(ddd.EventHandler[ddd.Event])
		return catalogHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get("domainDispatcher").(*ddd.EventDispatcher[ddd.Event])
	RegisterCatalogHandlers(subscriber, handlers)
}

// HandleEvent routes domain events to the correct "onXYZ" method.
func (h catalogHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {

	case domain.CategoryAddedEvent:
		return h.onCategoryAdded(ctx, event)

	case domain.CategoryUpdatedEvent:
		return h.onCategoryUpdated(ctx, event)

	//case domain.CategoryRebrandedEvent:
	//	return h.onCategoryRebranded(ctx, event)

	case domain.CategoryArchivedEvent:
		return h.onCategoryArchived(ctx, event)

	case domain.CategoryRemovedEvent:
		return h.onCategoryRemoved(ctx, event)
	}

	// no-op if event is unrecognized
	return nil
}

// ---------------------------------------------------------------------------
// HANDLER IMPLEMENTATIONS
// ---------------------------------------------------------------------------

// onCategoryAdded listens to the `CategoryAdded` event payload
func (h catalogHandlers[T]) onCategoryAdded(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.CategoryAdded) // referencing the event payload struct
	return h.catalog.AddCategory(
		ctx,
		e.CategoryID,
		e.Description,
		e.ParentID,
		e.GoogleCategoryID,
		e.Tags,
		e.IsActive,
		e.Slug,
		e.SeoTitle,
		e.SeoKeywords,
		e.SeoDesc,
		e.CategoryType,
		e.Lang,
	)
}

// onCategoryUpdated handles the optional `CategoryUpdated` event
func (h catalogHandlers[T]) onCategoryUpdated(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.CategoryUpdated)
	return h.catalog.UpdateCategory(
		ctx,
		e.CategoryID,
		e.Description,
		e.ParentID,
		e.GoogleCategoryID,
		e.Tags,
		e.IsActive,
		e.Slug,
		e.SeoTitle,
		e.SeoKeywords,
		e.SeoDesc,
		e.CategoryType,
		e.Lang,
	)
}

//// onCategoryRebranded changes only the name/description, if that's your domain logic.
//func (h catalogHandlers[T]) onCategoryRebranded(ctx context.Context, event ddd.Event) error {
//	e := event.Payload().(*domain.CategoryRebranded)
//	return h.catalog.RebrandCategory(
//		ctx,
//		e.CategoryID,
//		e.NewName,
//		e.NewDesc,
//	)
//}

// onCategoryArchived marks it as archived
func (h catalogHandlers[T]) onCategoryArchived(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.CategoryArchived)
	return h.catalog.ArchiveCategory(ctx, e.CategoryID, e.CategoryID)
}

// onCategoryRemoved handles removing a category from the system
func (h catalogHandlers[T]) onCategoryRemoved(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.CategoryRemoved)
	return h.catalog.RemoveCategory(ctx, e.CategoryID, e.UserID)
}
