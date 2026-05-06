package handlers

import (
	"context"
	"time"
	
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/newsletters/internal/domain"
	"middleman/newsletters/newsletterspb"
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
		domain.NewsletterCreatedEvent,
		domain.NewsletterUpdatedEvent,
		domain.NewsletterActivatedEvent,
		domain.NewsletterDeactivatedEvent,
		domain.NewsletterDeletedEvent,
		domain.SubscribedEvent,
		domain.UnsubscribedEvent,
		domain.PreferencesUpdatedEvent,
		domain.SubscriptionPausedEvent,
		domain.SubscriptionResumedEvent,
		domain.EditionCreatedEvent,
		domain.EditionUpdatedEvent,
		domain.EditionScheduledEvent,
		domain.EditionSendingEvent,
		domain.EditionSentEvent,
		domain.TemplateCreatedEvent,
		domain.TemplateUpdatedEvent,
		domain.TemplateDeletedEvent,
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
	case domain.NewsletterCreatedEvent:
		return h.onNewsletterCreated(ctx, event)
	case domain.SubscribedEvent:
		return h.onSubscribed(ctx, event)
	case domain.EditionSentEvent:
		return h.onEditionSent(ctx, event)
	}
	
	return nil
}

func (h domainHandlers[T]) onNewsletterCreated(ctx context.Context, event ddd.Event) error {
	newsletter := event.Payload().(*domain.Newsletter)
	return h.publisher.Publish(ctx, newsletterspb.NewsletterAggregateChannel,
		ddd.NewEvent(newsletterspb.NewsletterCreatedEvent, &newsletterspb.NewsletterCreated{
			NewsletterId: newsletter.ID(),
			UserId:       newsletter.UserID,
			Name:         newsletter.Name,
			Description:  newsletter.Description,
			Frequency:    newsletter.Frequency.String(),
			Category:     newsletter.Category,
			TemplateId:   newsletter.TemplateID,
		}),
	)
}

func (h domainHandlers[T]) onSubscribed(ctx context.Context, event ddd.Event) error {
	subscription := event.Payload().(*domain.Subscription)
	return h.publisher.Publish(ctx, newsletterspb.SubscriptionAggregateChannel,
		ddd.NewEvent(newsletterspb.SubscribedEvent, &newsletterspb.UserSubscribed{
			SubscriptionId: subscription.ID(),
			UserId:         subscription.UserID,
			NewsletterId:   subscription.NewsletterID,
		}),
	)
}

func (h domainHandlers[T]) onEditionSent(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.EditionSent)
	return h.publisher.Publish(ctx, newsletterspb.EditionAggregateChannel,
		ddd.NewEvent(newsletterspb.EditionSentEvent, &newsletterspb.EditionSent{
			EditionId:       payload.EditionID,
			RecipientsCount: int32(payload.RecipientCount),
		}),
	)
}