package handlers

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"middleman/internal/ddd"
	"middleman/internal/di"
	"middleman/internal/errorsotel"
	"middleman/reviews/internal/constants"
	"middleman/reviews/internal/domain"
	"time"
)

type middlemanHandlers[T ddd.Event] struct {
	middleman domain.MiddlemanCacheRepository
}

var _ ddd.EventHandler[ddd.Event] = (*middlemanHandlers[ddd.Event])(nil)

func NewMiddlemanHandlers(middleman domain.MiddlemanCacheRepository) ddd.EventHandler[ddd.Event] {
	return middlemanHandlers[ddd.Event]{
		middleman: middleman,
	}
}

func RegisterMiddlemanHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.ReviewAddedEvent,
		domain.ReviewEditedEvent,
	)
}

func RegisterMiddlemanHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		middlemanHandlers := di.Get(ctx, constants.MiddlemanHandlersKey).(ddd.EventHandler[ddd.Event])

		return middlemanHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.Event])

	RegisterMiddlemanHandlers(subscriber, handlers)
}

func (h middlemanHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling middleman event",
				trace.WithAttributes(errorsotel.ErrAttrs(err)...),
			)
		}
		span.AddEvent("Handled middleman event", trace.WithAttributes(
			attribute.Int64("TookMS", time.Since(started).Milliseconds()),
		))
	}(time.Now())

	span.AddEvent("Handling middleman event", trace.WithAttributes(
		attribute.String("Event", event.EventName()),
	))

	switch event.EventName() {
	case domain.ReviewAddedEvent:
		return h.onReviewAdded(ctx, event)
	case domain.ReviewEditedEvent:
		return h.onReviewAdded(ctx, event)
	}
	return nil
}

func (h middlemanHandlers[T]) onReviewAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Review)
	return h.middleman.Add(
		ctx,
		payload.ID(),
		payload.SenderID,
		payload.ItemID,
		domain.ToItemType(payload.ItemType),
		payload.Content,
		payload.CategoryID,
		payload.ParentID)
}

//func (h middlemanHandlers[T]) onReviewEdited(ctx context.Context, event ddd.Event) error {
//	payload := event.Payload().(*domain.Review)
//	return h.middleman.Edit(ctx, payload.ID(), payload.ItemID, payload.Content)
//}
