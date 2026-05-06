package handlers

import (
	"context"
	"middleman/media/internal/constants"
	"middleman/media/internal/domain"

	"middleman/internal/ddd"
	"middleman/internal/di"
)

// middlemanHandlers listens to domain events and updates the Middleman DB
type middlemanImageHandlers[T ddd.Event] struct {
	middleman domain.MiddlemanImageRepository
}

var _ ddd.EventHandler[ddd.Event] = (*middlemanImageHandlers[ddd.Event])(nil)

func NewMiddlemanImageHandlers(middleman domain.MiddlemanImageRepository) ddd.EventHandler[ddd.Event] {
	return middlemanImageHandlers[ddd.Event]{middleman: middleman}
}

func RegisterMiddlemanImageHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.ImageAddedEvent,
		domain.ImageRemovedEvent,

		// plus anything else you'd like
	)
}

func RegisterMiddlemanImageHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		middlemanHandlers := di.Get(ctx, constants.MiddlemanImageHandlersKey).(ddd.EventHandler[ddd.Event])
		return middlemanHandlers.HandleEvent(ctx, event)
	})
	subscriber := container.Get("domainDispatcher").(*ddd.EventDispatcher[ddd.Event])

	RegisterMiddlemanImageHandlers(subscriber, handlers)
}

func (h middlemanImageHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {

	case domain.ImageAddedEvent:
		return h.onImageAdded(ctx, event)

	case domain.ImageRemovedEvent:
		return h.onImageRemoved(ctx, event)
	}
	return nil
}

func (h middlemanImageHandlers[T]) onImageAdded(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.Image) // domain.Image after “added”
	return h.middleman.AddImage(ctx, e.ID(), e.MediaID, e.DisplayOrder, e.IsMain, e.URL, e.MetaData, e.Thumbnail, e.UserID)
}

func (h middlemanImageHandlers[T]) onImageRemoved(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.ImageRemoved)
	return h.middleman.RemoveImage(ctx, e.ID)
}
