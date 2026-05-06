package handlers

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"middleman/comments/commentspb"
	"middleman/comments/internal/domain"
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
		domain.CommentAddedEvent,
		domain.CommentRemovedEvent,
		domain.CommentFlaggedEvent,
		domain.CommentRejectedEvent,
		domain.CommentFlaggedEvent,
		domain.CommentEditedEvent,
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
	case domain.CommentAddedEvent:
		return h.onCommentAdded(ctx, event)
		//case domain.CommentApprovedEvent:
		//	return h.onCommentApproved(ctx, event)
		//case domain.CommentRemovedEvent:
		//	return h.onommentRemoved(ctx, event)
		//case domain.CommentEditedEvent:
		//	return h.onCommentEdited(ctx, event)
		//case domain.CommentFlaggedEvent:return h.onCommentFlagged(ctx, event)
		//case domain.CommentRejectedEvent:return h.onCommentRejected(ctx, event)

	}
	return nil
}

//	func (h domainHandlers[T]) onCommentRejected(ctx context.Context, event ddd.Event) error {
//		comment := event.Payload().(*domain.Comment)
//		return h.publisher.Publish(ctx, commentspb.CommentAggregateChannel,
//			ddd.NewEvent(commentspb.CommentAddedEvent, &commentspb.CommentAdded{
//				Id:        comment.ID(),
//				SenderId:  comment.SenderID,
//				ItemId: comment.ItemID,
//				Content:   comment.Content,
//			}),
//		)
//	}
//
//	func (h domainHandlers[T]) onCommentFlagged(ctx context.Context, event ddd.Event) error {
//		comment := event.Payload().(*domain.Comment)
//		return h.publisher.Publish(ctx, commentspb.CommentAggregateChannel,
//			ddd.NewEvent(commentspb.CommentAddedEvent, &commentspb.CommentAdded{
//				Id:        comment.ID(),
//				SenderId:  comment.SenderID,
//				ItemId: comment.ItemID,
//				Content:   comment.Content,
//			}),
//		)
//	}
//
//	func (h domainHandlers[T]) onCommentEdited(ctx context.Context, event ddd.Event) error {
//		comment := event.Payload().(*domain.Comment)
//		return h.publisher.Publish(ctx, commentspb.CommentAggregateChannel,
//			ddd.NewEvent(commentspb.CommentAddedEvent, &commentspb.CommentAdded{
//				Id:        comment.ID(),
//				SenderId:  comment.SenderID,
//				ItemId: comment.ItemID,
//				Content:   comment.Content,
//			}),
//		)
//	}
func (h domainHandlers[T]) onCommentAdded(ctx context.Context, event ddd.Event) error {
	comment := event.Payload().(*domain.Comment)
	return h.publisher.Publish(ctx, commentspb.CommentAggregateChannel,
		ddd.NewEvent(commentspb.CommentAddedEvent, &commentspb.CommentAdded{
			Id:         comment.ID(),
			SenderId:   comment.SenderID,
			ItemId:     comment.ItemID,
			ItemType:   comment.ItemType,
			Content:    comment.Content,
			CategoryId: comment.CategoryID,
			ParentId:   comment.ParentID,
		}),
	)
}

//func (h domainHandlers[T]) onCommentApproved(ctx context.Context, event ddd.Event) error {
//	comment := event.Payload().(*domain.Comment)
//	return h.publisher.Publish(ctx, commentspb.CommentAggregateChannel,
//		ddd.NewEvent(commentspb.CommentAddedEvent, &commentspb.CommentAdded{
//			Id:        comment.ID(),
//			SenderId:  comment.SenderID,
//			ItemId: comment.ItemID,
//			Content:   comment.Content,
//		}),
//	)
//}
//func (h domainHandlers[T]) onommentRemoved(ctx context.Context, event ddd.Event) error {
//	comment := event.Payload().(*domain.Comment)
//	return h.publisher.Publish(ctx, commentspb.CommentAggregateChannel,
//		ddd.NewEvent(commentspb.CommentAddedEvent, &commentspb.CommentAdded{
//			Id:        comment.ID(),
//			SenderId:  comment.SenderID,
//			ItemId: comment.ItemID,
//			Content:   comment.Content,
//		}),
//	)
//}
