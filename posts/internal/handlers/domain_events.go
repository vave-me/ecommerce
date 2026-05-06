package handlers

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/posts/internal/domain"
	"middleman/posts/postspb"
)

type domainHandlers[T ddd.Event] struct {
	publisher am.EventPublisher
}

// Ensure it implements the interface
var _ ddd.EventHandler[ddd.Event] = (*domainHandlers[ddd.Event])(nil)

// NewDomainEventHandlers constructs the domainHandlers with an AM publisher.
func NewDomainEventHandlers(publisher am.EventPublisher) ddd.EventHandler[ddd.Event] {
	return &domainHandlers[ddd.Event]{
		publisher: publisher,
	}
}

// RegisterDomainEventHandlers subscribes to the relevant domain events.
func RegisterDomainEventHandlers(
	subscriber ddd.EventSubscriber[ddd.Event],
	handlers ddd.EventHandler[ddd.Event],
) {
	subscriber.Subscribe(
		handlers,
		domain.PostAddedEvent,
		domain.PostUpdatedEvent,
		domain.PostArchivedEvent,
		domain.PostRemovedEvent,
		domain.PostThumbnailAddedEvent,
		domain.PostThumbnailUpdatedEvent,
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

	// Switch on the domain event name
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

// -----------------------------------------------------------------------------
// 1) onPostAdded
// -----------------------------------------------------------------------------
func (h domainHandlers[T]) onPostAdded(ctx context.Context, event ddd.Event) error {
	post := event.Payload().(*domain.Post)
	return h.publisher.Publish(
		ctx,
		postspb.PostAggregateChannel,
		ddd.NewEvent(
			postspb.PostAddedEvent,
			&postspb.PostAdded{
				Id:           post.ID(),
				Name:         post.Name,
				Description:  post.Description,
				TypeOfPost:   post.TypeOfPost.String(),
				UserId:       post.UserID,
				UserType:     post.UserType.String(),
				CategoryId:   post.CategoryID,
				CategorySlug: post.CategorySlug,
				Tags:         post.Tags,
				Status:       string(post.Status),
				Thumbnail:    post.Thumbnail,
				Lat:          float32(post.Lat),
				Lng:          float32(post.Lng),
			},
		),
	)
}

// -----------------------------------------------------------------------------
// 2) onPostUpdated
// -----------------------------------------------------------------------------
func (h domainHandlers[T]) onPostUpdated(ctx context.Context, event ddd.Event) error {
	post := event.Payload().(*domain.Post)
	return h.publisher.Publish(
		ctx,
		postspb.PostAggregateChannel,
		ddd.NewEvent(
			postspb.PostUpdatedEvent,
			&postspb.PostUpdated{
				Name:         post.Name,
				Description:  post.Description,
				TypeOfPost:   post.TypeOfPost.String(),
				CategoryId:   post.CategoryID,
				CategorySlug: post.CategorySlug,
				Tags:         post.Tags,
				Status:       string(post.Status),
				Thumbnail:    post.Thumbnail,
			},
		),
	)
}

// -----------------------------------------------------------------------------
// 3) onPostArchived
// -----------------------------------------------------------------------------
func (h domainHandlers[T]) onPostArchived(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.Post) // or *domain.Post if you store it differently
	return h.publisher.Publish(
		ctx,
		postspb.PostAggregateChannel,
		ddd.NewEvent(
			postspb.PostArchivedEvent,
			&postspb.PostArchived{
				Id:     payload.ID(),
				UserId: payload.UserID, // Adjust field name as needed
			},
		),
	)
}

// -----------------------------------------------------------------------------
// 4) onPostRemoved
// -----------------------------------------------------------------------------
func (h domainHandlers[T]) onPostRemoved(ctx context.Context, event ddd.Event) error {
	post := event.Payload().(*domain.Post)
	return h.publisher.Publish(
		ctx,
		postspb.PostAggregateChannel,
		ddd.NewEvent(
			postspb.PostRemovedEvent,
			&postspb.PostRemoved{
				Id:     post.ID(),
				UserId: post.UserID, // if desired
			},
		),
	)
}
func (h domainHandlers[T]) onPostThumbnailAdded(ctx context.Context, event ddd.Event) error {
	post := event.Payload().(*domain.PostThumbnailAdded)
	return h.publisher.Publish(
		ctx,
		postspb.PostAggregateChannel,
		ddd.NewEvent(
			postspb.PostThumbnailAddedEvent,
			&postspb.PostThumbnailAdded{
				Id:        post.PostID,
				Thumbnail: post.Thumbnail,
			},
		),
	)
}

func (h domainHandlers[T]) onPostThumbnailUpdated(ctx context.Context, event ddd.Event) error {
	post := event.Payload().(*domain.PostThumbnailUpdated)
	return h.publisher.Publish(
		ctx,
		postspb.PostAggregateChannel,
		ddd.NewEvent(
			postspb.PostThumbnailUpdatedEvent,
			&postspb.PostThumbnailUpdated{
				Id:        post.PostID,
				Thumbnail: post.Thumbnail,
			},
		),
	)
}
