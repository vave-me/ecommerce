package handlers

import (
	"context"
	"middleman/media/internal/constants"
	"middleman/media/internal/domain"

	"middleman/internal/ddd"
	"middleman/internal/di"
)

// middlemanHandlers listens to domain events and updates the Middleman DB
type middlemanVideoHandlers[T ddd.Event] struct {
	middleman domain.MiddlemanVideoRepository
}

var _ ddd.EventHandler[ddd.Event] = (*middlemanVideoHandlers[ddd.Event])(nil)

func NewMiddlemanVideoHandlers(middleman domain.MiddlemanVideoRepository) ddd.EventHandler[ddd.Event] {
	return middlemanVideoHandlers[ddd.Event]{middleman: middleman}
}

func RegisterMiddlemanVideoHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.VideoAddedEvent,
		domain.VideoRemovedEvent,

		// plus anything else you'd like
	)
}

func RegisterMiddlemanVideoHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		middlemanHandlers := di.Get(ctx, constants.MiddlemanVideoHandlersKey).(ddd.EventHandler[ddd.Event])
		return middlemanHandlers.HandleEvent(ctx, event)
	})
	subscriber := container.Get("domainDispatcher").(*ddd.EventDispatcher[ddd.Event])

	RegisterMiddlemanVideoHandlers(subscriber, handlers)
}

func (h middlemanVideoHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {

	case domain.VideoAddedEvent:
		return h.onVideoAdded(ctx, event)

	case domain.VideoRemovedEvent:
		return h.onVideoRemoved(ctx, event)
	}
	return nil
}

func (h middlemanVideoHandlers[T]) onVideoAdded(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.Video) // domain.Video after “added”
	return h.middleman.AddVideo(ctx, e.ID(), e.MediaID, e.DisplayOrder, e.IsMain, e.URL, e.MetaData, e.Thumbnail, e.UserID)
}

func (h middlemanVideoHandlers[T]) onVideoRemoved(ctx context.Context, event ddd.Event) error {
	e := event.Payload().(*domain.VideoRemoved)
	return h.middleman.RemoveVideo(ctx, e.ID)
}
