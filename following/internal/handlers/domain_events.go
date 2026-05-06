package handlers

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"middleman/following/followingpb"
	"middleman/following/internal/domain"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"time"
)

type domainHandlers[T ddd.Event] struct {
	publisher am.EventPublisher
}

func NewDomainEventHandlers(publisher am.EventPublisher) ddd.EventHandler[ddd.Event] {
	return &domainHandlers[ddd.Event]{
		publisher: publisher,
	}
}

func RegisterDomainEventHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.FollowAddedEvent,
		domain.FollowRemovedEvent,
		domain.FollowFlaggedEvent,
		domain.FollowRejectedEvent,
		domain.FollowFlaggedEvent,
		domain.FollowEditedEvent,
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
	case domain.FollowAddedEvent:
		return h.onFollowAdded(ctx, event)
		//case domain.FollowApprovedEvent:
		//	return h.onFollowApproved(ctx, event)
		//case domain.FollowRemovedEvent:
		//	return h.onommentRemoved(ctx, event)
		//case domain.FollowEditedEvent:
		//	return h.onFollowEdited(ctx, event)
		//case domain.FollowFlaggedEvent:return h.onFollowFlagged(ctx, event)
		//case domain.FollowRejectedEvent:return h.onFollowRejected(ctx, event)

	}
	return nil
}

//	func (h domainHandlers[T]) onFollowRejected(ctx context.Context, event ddd.Event) error {
//		follow := event.Payload().(*domain.Follow)
//		return h.publisher.Publish(ctx, followingpb.FollowAggregateChannel,
//			ddd.NewEvent(followingpb.FollowAddedEvent, &followingpb.FollowAdded{
//				Id:        follow.ID(),
//				UserId:  follow.UserID,
//				FollowedUserId: follow.FollowedUserID,
//				Content:   follow.Content,
//			}),
//		)
//	}
//
//	func (h domainHandlers[T]) onFollowFlagged(ctx context.Context, event ddd.Event) error {
//		follow := event.Payload().(*domain.Follow)
//		return h.publisher.Publish(ctx, followingpb.FollowAggregateChannel,
//			ddd.NewEvent(followingpb.FollowAddedEvent, &followingpb.FollowAdded{
//				Id:        follow.ID(),
//				UserId:  follow.UserID,
//				FollowedUserId: follow.FollowedUserID,
//				Content:   follow.Content,
//			}),
//		)
//	}
//
//	func (h domainHandlers[T]) onFollowEdited(ctx context.Context, event ddd.Event) error {
//		follow := event.Payload().(*domain.Follow)
//		return h.publisher.Publish(ctx, followingpb.FollowAggregateChannel,
//			ddd.NewEvent(followingpb.FollowAddedEvent, &followingpb.FollowAdded{
//				Id:        follow.ID(),
//				UserId:  follow.UserID,
//				FollowedUserId: follow.FollowedUserID,
//				Content:   follow.Content,
//			}),
//		)
//	}
func (h domainHandlers[T]) onFollowAdded(ctx context.Context, event ddd.Event) error {
	follow := event.Payload().(*domain.Follow)
	return h.publisher.Publish(ctx, followingpb.FollowAggregateChannel,
		ddd.NewEvent(followingpb.FollowAddedEvent, &followingpb.FollowAdded{
			Id:               follow.ID(),
			UserId:           follow.UserID,
			FollowedUserId:   follow.FollowedUserID,
			FollowedUserType: follow.FollowedUserType,
			Content:          follow.Content,
			CategoryId:       follow.CategoryID,
			ParentId:         follow.ParentID,
		}),
	)
}

//func (h domainHandlers[T]) onFollowApproved(ctx context.Context, event ddd.Event) error {
//	follow := event.Payload().(*domain.Follow)
//	return h.publisher.Publish(ctx, followingpb.FollowAggregateChannel,
//		ddd.NewEvent(followingpb.FollowAddedEvent, &followingpb.FollowAdded{
//			Id:        follow.ID(),
//			UserId:  follow.UserID,
//			FollowedUserId: follow.FollowedUserID,
//			Content:   follow.Content,
//		}),
//	)
//}
//func (h domainHandlers[T]) onommentRemoved(ctx context.Context, event ddd.Event) error {
//	follow := event.Payload().(*domain.Follow)
//	return h.publisher.Publish(ctx, followingpb.FollowAggregateChannel,
//		ddd.NewEvent(followingpb.FollowAddedEvent, &followingpb.FollowAdded{
//			Id:        follow.ID(),
//			UserId:  follow.UserID,
//			FollowedUserId: follow.FollowedUserID,
//			Content:   follow.Content,
//		}),
//	)
//}
