package handlers

import (
	"context"

	"middleman/internal/ddd"
	"middleman/internal/di"
	"middleman/posts/internal/domain"
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
		domain.PostAddedEvent,
		domain.PostRebrandedEvent,
		domain.PostUpdatedEvent,
		domain.PostRemovedEvent,
		domain.PostArchivedEvent,
		domain.PostThumbnailAddedEvent,
		domain.PostThumbnailUpdatedEvent,
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

	case domain.PostAddedEvent:
		return h.onPostAdded(ctx, event)
	case domain.PostUpdatedEvent:
		return h.onPostUpdated(ctx, event)

	case domain.PostArchivedEvent:
		return h.onPostArchived(ctx, event)

	case domain.PostRemovedEvent:
		return h.onPostRemoved(ctx, event)
	case domain.PostThumbnailAddedEvent:
		return h.onPostThumbnailAdded(ctx, event)
	case domain.PostThumbnailUpdatedEvent:
		return h.onPostThumbnailUpdated(ctx, event)
	}
	return nil
}

// TODO Middleman service bool or sting with couple strings as a diffrent services
func (h catalogHandlers[T]) onPostAdded(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.Post) // domain.Post after “added”
	return h.catalog.AddPost(ctx, e.ID(), e.Name, e.Description, e.TypeOfPost, e.UserID, e.UserType, e.CategoryID, e.CategorySlug, e.Tags, e.Status, e.Thumbnail, e.Lat, e.Lng)
}
func (h catalogHandlers[T]) onPostUpdated(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.Post) // domain.Post after “added”
	return h.catalog.UpdatePost(ctx, e.ID(), e.Name, e.Description, e.TypeOfPost, e.UserID, e.Tags, e.Status, e.Thumbnail)
}

func (h catalogHandlers[T]) onPostArchived(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.PostArchived)
	return h.catalog.ArchivePost(ctx, e.PostID, e.UserID)
}

func (h catalogHandlers[T]) onPostRemoved(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.PostRemoved)
	return h.catalog.RemovePost(ctx, e.PostID, e.UserID)
}
func (h catalogHandlers[T]) onPostThumbnailAdded(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.PostThumbnailAdded)
	return h.catalog.UpdatePostThumbnail(ctx, e.PostID, e.Thumbnail)
}
func (h catalogHandlers[T]) onPostThumbnailUpdated(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.PostThumbnailUpdated)
	return h.catalog.UpdatePostThumbnail(ctx, e.PostID, e.Thumbnail)
}

// a small helper if you want to store "premium"/"base" as bool. Adjust as needed.
func strToBool(s string) bool {
	return s == "true" || s == "premium"
}
