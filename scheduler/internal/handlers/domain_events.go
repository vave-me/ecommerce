package handlers

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/scheduler/internal/domain"
	"middleman/scheduler/schedulerspb"
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
		domain.SchedulerCreatedEvent,
		domain.ActionAddedEvent,
		domain.ActionRemovedEvent,
		domain.ActionUpdatedEvent,
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
	case domain.SchedulerCreatedEvent:
		return h.onSchedulerCreated(ctx, event)
	case domain.ActionAddedEvent:
		return h.onActionAdded(ctx, event)
	case domain.ActionRemovedEvent:
		return h.onActionRemoved(ctx, event)
	case domain.ActionUpdatedEvent:
		return h.onActionUpdated(ctx, event)
	}
	return nil
}

func (h domainHandlers[T]) onActionAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ActionAdded)
	return h.publisher.Publish(ctx, schedulerspb.ActionAggregateChannel,
		ddd.NewEvent(schedulerspb.ActionAddedEvent, &schedulerspb.ActionAdded{
			Id:                  payload.ID,
			SchedulerId:         payload.SchedulerID,
			NaturalLanguageTask: payload.NaturalLanguageTask,
			ExecutionTime:       timestamppb.New(payload.ExecutionTime),
			Status:              payload.Status,
			CreatedAt:           timestamppb.New(payload.CreatedAt),
		}),
	)
}
func (h domainHandlers[T]) onSchedulerCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.SchedulerCreated)
	return h.publisher.Publish(ctx, schedulerspb.SchedulerAggregateChannel,
		ddd.NewEvent(schedulerspb.SchedulerCreatedEvent, &schedulerspb.SchedulerCreated{
			Id:     payload.ID,
			UserId: payload.UserID,
		}),
	)
}

func (h domainHandlers[T]) onActionRemoved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ActionRemoved)
	return h.publisher.Publish(ctx, schedulerspb.ActionAggregateChannel,
		ddd.NewEvent(schedulerspb.ActionRemovedEvent, &schedulerspb.ActionRemoved{
			Id:          payload.ActionID,
			SchedulerId: payload.SchedulerID,
		}),
	)
}

func (h domainHandlers[T]) onActionUpdated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ActionUpdated)
	pbEvent := &schedulerspb.ActionUpdated{
		Id:          payload.ActionID,
		SchedulerId: payload.SchedulerID,
		Status:      payload.Status,
		Result:      payload.Result,
		ErrorMessage: payload.ErrorMessage,
	}
	if payload.ExecutedAt != nil {
		pbEvent.ExecutedAt = timestamppb.New(*payload.ExecutedAt)
	}
	return h.publisher.Publish(ctx, schedulerspb.ActionAggregateChannel,
		ddd.NewEvent(schedulerspb.ActionUpdatedEvent, pbEvent),
	)
}
