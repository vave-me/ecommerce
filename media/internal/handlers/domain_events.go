package handlers

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/media/internal/domain"
	"middleman/media/mediapb"
	"time"
)

type domainHandlers[T ddd.Event] struct {
	publisher am.EventPublisher
}

var _ ddd.EventHandler[ddd.Event] = (*domainHandlers[ddd.Event])(nil)

func NewDomainEventHandlers(publisher am.EventPublisher) ddd.EventHandler[ddd.Event] {
	return &domainHandlers[ddd.Event]{
		publisher: publisher,
	}
}

func RegisterDomainEventHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.MediaCreatedEvent,
		//domain.MediaDeletedEvent,
		domain.MediaUpdatedEvent,
		domain.ImageAddedEvent,
		domain.VideoAddedEvent,
	)
}

func (h domainHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling domain event",
				trace.WithAttributes(errorsotel.ErrAttrs(err)...),
			)
		}
		span.AddEvent("Handled domain event", trace.WithAttributes(
			attribute.Int64("TookMS", time.Since(started).Milliseconds()),
		))
	}(time.Now())

	span.AddEvent("Handling domain event", trace.WithAttributes(
		attribute.String("Event", event.EventName()),
	))

	switch event.EventName() {
	case domain.MediaCreatedEvent:
		return h.onMediaCreated(ctx, event)
	case domain.MediaUpdatedEvent:
		return h.onMediaUpdated(ctx, event)
	//case domain.MediaDeletedEvent:return h.onMediaDeleted(ctx, event)
	case domain.ImageAddedEvent:
		return h.onImageAdded(ctx, event)
	case domain.VideoAddedEvent:
		return h.onVideoAdded(ctx, event)

	}
	return nil
}

func (h domainHandlers[T]) onMediaCreated(ctx context.Context, event ddd.Event) error {
	media := event.Payload().(*domain.Media)
	return h.publisher.Publish(ctx, mediapb.MediaAggregateChannel,
		ddd.NewEvent(mediapb.MediaCreatedEvent, &mediapb.MediaCreated{
			Id:       media.ID(),
			ItemId:   media.ItemID,
			ItemType: string(media.ItemType),
			UserId:   media.UserID,
			Status:   string(media.Status),
		}),
	)
}

func (h domainHandlers[T]) onMediaUpdated(ctx context.Context, event ddd.Event) error {
	media := event.Payload().(*domain.Media)
	return h.publisher.Publish(ctx, mediapb.MediaAggregateChannel,
		ddd.NewEvent(mediapb.MediaUpdatedEvent, &mediapb.MediaUpdated{
			Id:       media.ID(),
			ItemId:   media.ItemID,
			ItemType: string(media.ItemType),
			UserId:   media.UserID,
			Status:   string(media.Status),
		}),
	)
}
func (h domainHandlers[T]) onImageAdded(ctx context.Context, event ddd.Event) error {
	image := event.Payload().(*domain.Image)
	return h.publisher.Publish(ctx, mediapb.ImageAggregateChannel,
		ddd.NewEvent(mediapb.ImageAddedEvent, &mediapb.ImageAdded{
			Id:           image.ID(),
			MediaId:      image.MediaID,
			DisplayOrder: int32(image.DisplayOrder),
			IsMain:       image.IsMain,
			Url:          image.URL,
			Metadata:     image.MetaData,
			FileType:     image.FileType,
			Thumbnail:    image.Thumbnail,
			UserId:       image.UserID,
		}),
	)
}

func (h domainHandlers[T]) onVideoAdded(ctx context.Context, event ddd.Event) error {
	video := event.Payload().(*domain.Video)
	return h.publisher.Publish(ctx, mediapb.VideoAggregateChannel,
		ddd.NewEvent(mediapb.VideoAddedEvent, &mediapb.VideoAdded{
			Id:           video.ID(),
			MediaId:      video.MediaID,
			DisplayOrder: int32(video.DisplayOrder),
			IsMain:       video.IsMain,
			Url:          video.URL,
			Metadata:     video.MetaData,
			Thumbnail:    video.Thumbnail,
			FileType:     video.FileType,
			UserId:       video.UserID,
		}),
	)
}
