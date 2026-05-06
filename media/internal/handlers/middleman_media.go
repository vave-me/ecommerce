package handlers

import (
	"context"
	"middleman/media/internal/constants"
	"middleman/media/internal/domain"

	"middleman/internal/ddd"
	"middleman/internal/di"
)

// middlemanHandlers listens to domain events and updates the Middleman DB
type middlemanMediaHandlers[T ddd.Event] struct {
	middleman domain.MiddlemanMediaRepository
}

var _ ddd.EventHandler[ddd.Event] = (*middlemanMediaHandlers[ddd.Event])(nil)

func NewMiddlemanMediaHandlers(middleman domain.MiddlemanMediaRepository) ddd.EventHandler[ddd.Event] {
	return middlemanMediaHandlers[ddd.Event]{middleman: middleman}
}

func RegisterMiddlemanMediaHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.MediaCreatedEvent,
		domain.MediaUpdatedEvent,
		domain.ImageAddedEvent,
		domain.VideoAddedEvent,
		domain.ImageRemovedEvent,

		// plus anything else you'd like
	)
}

func RegisterMiddlemanMediaHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		middlemanHandlers := di.Get(ctx, constants.MiddlemanMediaHandlersKey).(ddd.EventHandler[ddd.Event])
		return middlemanHandlers.HandleEvent(ctx, event)
	})
	subscriber := container.Get("domainDispatcher").(*ddd.EventDispatcher[ddd.Event])

	RegisterMiddlemanMediaHandlers(subscriber, handlers)
}

func (h middlemanMediaHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {

	case domain.MediaCreatedEvent:
		return h.onMediaCreated(ctx, event)
	case domain.MediaUpdatedEvent:
		return h.onMediaUpdated(ctx, event)
	case domain.ImageRemovedEvent:
		return h.onMediaRemoved(ctx, event)
	case domain.ImageAddedEvent:
		return h.onImageAdded(ctx, event)
	case domain.VideoAddedEvent:
		return h.onVideoAdded(ctx, event)

	}
	return nil
}

func (h middlemanMediaHandlers[T]) onMediaCreated(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.Media) // domain.Image after “added”
	return h.middleman.AddMedia(ctx, e.ID(), e.ItemID, e.ItemType, e.UserID, e.Status)
}
func (h middlemanMediaHandlers[T]) onMediaUpdated(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.Media) // domain.Image after “added”
	return h.middleman.UpdateMedia(ctx, e.ID(), e.ItemID, e.ItemType, e.UserID, e.Status)
}
func (h middlemanMediaHandlers[T]) onMediaRemoved(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.MediaDeleted)
	return h.middleman.RemoveMedia(ctx, e.ID)
}

func (h middlemanMediaHandlers[T]) onImageAdded(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.Image)
	return h.middleman.AddMediaItemOrder(ctx, e.MediaID, e.ID(), e.URL, e.DisplayOrder)
}

func (h middlemanMediaHandlers[T]) onVideoAdded(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.Video)
	return h.middleman.AddMediaItemOrder(ctx, e.MediaID, e.ID(), e.URL, e.DisplayOrder)
}
