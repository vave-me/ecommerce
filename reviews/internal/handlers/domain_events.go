package handlers

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/reviews/internal/domain"
	"middleman/reviews/reviewspb"
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
		domain.ReviewAddedEvent,
		domain.ReviewRemovedEvent,
		domain.ReviewFlaggedEvent,
		domain.ReviewRejectedEvent,
		domain.ReviewFlaggedEvent,
		domain.ReviewEditedEvent,
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
	case domain.ReviewAddedEvent:
		return h.onReviewAdded(ctx, event)
		//case domain.ReviewApprovedEvent:
		//	return h.onReviewApproved(ctx, event)
		//case domain.ReviewRemovedEvent:
		//	return h.onommentRemoved(ctx, event)
		//case domain.ReviewEditedEvent:
		//	return h.onReviewEdited(ctx, event)
		//case domain.ReviewFlaggedEvent:return h.onReviewFlagged(ctx, event)
		//case domain.ReviewRejectedEvent:return h.onReviewRejected(ctx, event)

	}
	return nil
}

//	func (h domainHandlers[T]) onReviewRejected(ctx context.Context, event ddd.Event) error {
//		review := event.Payload().(*domain.Review)
//		return h.publisher.Publish(ctx, reviewspb.ReviewAggregateChannel,
//			ddd.NewEvent(reviewspb.ReviewAddedEvent, &reviewspb.ReviewAdded{
//				Id:        review.ID(),
//				SenderId:  review.SenderID,
//				ItemId: review.ItemID,
//				Content:   review.Content,
//			}),
//		)
//	}
//
//	func (h domainHandlers[T]) onReviewFlagged(ctx context.Context, event ddd.Event) error {
//		review := event.Payload().(*domain.Review)
//		return h.publisher.Publish(ctx, reviewspb.ReviewAggregateChannel,
//			ddd.NewEvent(reviewspb.ReviewAddedEvent, &reviewspb.ReviewAdded{
//				Id:        review.ID(),
//				SenderId:  review.SenderID,
//				ItemId: review.ItemID,
//				Content:   review.Content,
//			}),
//		)
//	}
//
//	func (h domainHandlers[T]) onReviewEdited(ctx context.Context, event ddd.Event) error {
//		review := event.Payload().(*domain.Review)
//		return h.publisher.Publish(ctx, reviewspb.ReviewAggregateChannel,
//			ddd.NewEvent(reviewspb.ReviewAddedEvent, &reviewspb.ReviewAdded{
//				Id:        review.ID(),
//				SenderId:  review.SenderID,
//				ItemId: review.ItemID,
//				Content:   review.Content,
//			}),
//		)
//	}
func (h domainHandlers[T]) onReviewAdded(ctx context.Context, event ddd.Event) error {
	review := event.Payload().(*domain.Review)
	return h.publisher.Publish(ctx, reviewspb.ReviewAggregateChannel,
		ddd.NewEvent(reviewspb.ReviewAddedEvent, &reviewspb.ReviewAdded{
			Id:         review.ID(),
			SenderId:   review.SenderID,
			ItemId:     review.ItemID,
			ItemType:   review.ItemType,
			Content:    review.Content,
			CategoryId: review.CategoryID,
			ParentId:   review.ParentID,
		}),
	)
}

//func (h domainHandlers[T]) onReviewApproved(ctx context.Context, event ddd.Event) error {
//	review := event.Payload().(*domain.Review)
//	return h.publisher.Publish(ctx, reviewspb.ReviewAggregateChannel,
//		ddd.NewEvent(reviewspb.ReviewAddedEvent, &reviewspb.ReviewAdded{
//			Id:        review.ID(),
//			SenderId:  review.SenderID,
//			ItemId: review.ItemID,
//			Content:   review.Content,
//		}),
//	)
//}
//func (h domainHandlers[T]) onommentRemoved(ctx context.Context, event ddd.Event) error {
//	review := event.Payload().(*domain.Review)
//	return h.publisher.Publish(ctx, reviewspb.ReviewAggregateChannel,
//		ddd.NewEvent(reviewspb.ReviewAddedEvent, &reviewspb.ReviewAdded{
//			Id:        review.ID(),
//			SenderId:  review.SenderID,
//			ItemId: review.ItemID,
//			Content:   review.Content,
//		}),
//	)
//}
