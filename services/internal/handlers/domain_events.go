package handlers

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/services/internal/domain"
	"middleman/services/servicespb"

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
		domain.ServiceAddedEvent,
		domain.ServiceUpdatedEvent,
		domain.ServiceRebrandedEvent,
		domain.ServicePriceIncreasedEvent,
		domain.ServicePriceDecreasedEvent,
		domain.ServiceStockAdjustedEvent,
		domain.ServiceRemovedEvent,
		domain.ServiceArchivedEvent,
		domain.ServiceSoldEvent,
		domain.ServiceLeasedEvent,
		domain.ServicePawnedEvent,
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

	case domain.ServiceAddedEvent:
		return h.onServiceAdded(ctx, event)
	case domain.ServiceUpdatedEvent:
		return h.onServiceUpdated(ctx, event)
	case domain.ServiceRebrandedEvent:
		return h.onServiceRebranded(ctx, event)
	case domain.ServicePriceIncreasedEvent:
		return h.onServicePriceIncreased(ctx, event)
	case domain.ServicePriceDecreasedEvent:
		return h.onServicePriceDecreased(ctx, event)
	case domain.ServiceStockAdjustedEvent:
		return h.onServiceStockAdjusted(ctx, event)
	case domain.ServiceArchivedEvent:
		return h.onServiceArchived(ctx, event)
	case domain.ServiceRemovedEvent:
		return h.onServiceRemoved(ctx, event)
	case domain.ServiceSoldEvent:
		return h.onServiceSold(ctx, event)
	case domain.ServiceLeasedEvent:
		return h.onServiceLeased(ctx, event)
	case domain.ServicePawnedEvent:
		return h.onServicePawned(ctx, event)

	}
	return nil

}

func (h domainHandlers[T]) onServiceAdded(ctx context.Context, event ddd.Event) error {
	service := event.Payload().(*domain.Service)
	return h.publisher.Publish(ctx, servicespb.ServiceAggregateChannel,
		ddd.NewEvent(servicespb.ServiceAddedEvent, &servicespb.ServiceAdded{
			Id:               service.ID(),
			Name:             service.Name,
			Description:      service.Description,
			ServiceType:      service.ServiceType,
			BasePrice:        service.BasePrice,
			Pricing:          service.Pricing,
			Availability:     service.Availability,
			ProviderName:     service.ProviderName,
			UserId:           service.UserID,
			CategoryId:       service.CategoryID,
			CategorySlug:     service.CategorySlug,
			DescriptionShort: service.DescriptionShort,
			DescriptionLong:  service.DescriptionLong,
			Qualifications:   service.Qualifications,
			Contact:          service.Contact,
			Faq:              service.Faq,
			Tags:             service.Tags,
			Status:           service.Status.String(),
			UserType:         service.UserType.String(),
			ShippingCost:     service.ShippingCost,
			Negotiable:       service.Negotiable,
			HasVariants:      service.HasVariants,
			MiddlemanService: service.MiddlemanService,
			Attributes:       domainAttributesToProto(service.Attributes),
			Options:          domainOptionsToProto(service.Options),
			Thumbnail:        service.Thumbnail,
			Lat:              float32(service.Lat),
			Lng:              float32(service.Lng),
		}),
	)
}

func (h domainHandlers[T]) onServiceUpdated(ctx context.Context, event ddd.Event) error {
	service := event.Payload().(*domain.Service)
	return h.publisher.Publish(ctx, servicespb.ServiceAggregateChannel,
		ddd.NewEvent(servicespb.ServiceUpdatedEvent, &servicespb.ServiceUpdated{
			Id:               service.ID(),
			Name:             service.Name,
			Description:      service.Description,
			ServiceType:      service.ServiceType,
			BasePrice:        service.BasePrice,
			Pricing:          service.Pricing,
			Availability:     service.Availability,
			ProviderName:     service.ProviderName,
			UserId:           service.UserID,
			CategoryId:       service.CategoryID,
			CategorySlug:     service.CategorySlug,
			DescriptionShort: service.DescriptionShort,
			DescriptionLong:  service.DescriptionLong,
			Qualifications:   service.Qualifications,
			Contact:          service.Contact,
			Faq:              service.Faq,
			Tags:             service.Tags,
			Status:           service.Status.String(),
			UserType:         service.UserType.String(),
			ShippingCost:     service.ShippingCost,
			Negotiable:       service.Negotiable,
			HasVariants:      service.HasVariants,
			MiddlemanService: service.MiddlemanService,
			Attributes:       domainAttributesToProto(service.Attributes),
			Options:          domainOptionsToProto(service.Options),
			Thumbnail:        service.Thumbnail,
			Lat:              float32(service.Lat),
			Lng:              float32(service.Lng),
		}),
	)
}
func (h domainHandlers[T]) onServiceRebranded(ctx context.Context, event ddd.Event) error {
	service := event.Payload().(*domain.Service)
	return h.publisher.Publish(ctx, servicespb.ServiceAggregateChannel,
		ddd.NewEvent(servicespb.ServiceRebrandedEvent, &servicespb.ServiceRebranded{
			Id:          service.ID(),
			Name:        service.Name,
			Description: service.Description,
		}),
	)
}

func (h domainHandlers[T]) onServicePriceIncreased(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ServicePriceIncreased)
	return h.publisher.Publish(ctx, servicespb.ServiceAggregateChannel,
		ddd.NewEvent(servicespb.ServicePriceIncreasedEvent, &servicespb.ServicePriceIncreased{
			Id:       payload.ServiceID,
			OldPrice: payload.OldPrice,
			NewPrice: payload.NewPrice,
		}),
	)
}

func (h domainHandlers[T]) onServiceArchived(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.ServiceArchived)
	return h.publisher.Publish(ctx, servicespb.ServiceAggregateChannel,
		ddd.NewEvent(servicespb.ServiceArchivedEvent, &servicespb.ServiceArchived{
			Id:     payload.ServiceID,
			UserId: payload.UserID,
		}),
	)
}

func (h domainHandlers[T]) onServiceSold(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.ServiceSold)
	return h.publisher.Publish(ctx, servicespb.ServiceAggregateChannel,
		ddd.NewEvent(servicespb.ServiceSoldEvent, &servicespb.ServiceSold{
			Id:     payload.ServiceID,
			UserId: payload.UserID,

			// BuyerId:   payload.BuyerID, if your proto has it
		}),
	)
}

func (h domainHandlers[T]) onServiceLeased(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.ServiceLeased)
	return h.publisher.Publish(ctx, servicespb.ServiceAggregateChannel,
		ddd.NewEvent(servicespb.ServiceLeasedEvent, &servicespb.ServiceLeased{
			Id:     payload.ServiceID,
			UserId: payload.UserID,

			// LesseeId: payload.LesseeID, if needed
		}),
	)
}

func (h domainHandlers[T]) onServicePriceDecreased(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ServicePriceDecreased)
	return h.publisher.Publish(ctx, servicespb.ServiceAggregateChannel,
		ddd.NewEvent(servicespb.ServicePriceDecreasedEvent, &servicespb.ServicePriceDecreased{
			Id:       payload.ServiceID,
			OldPrice: payload.OldPrice,
			NewPrice: payload.NewPrice,
		}),
	)
}
func (h domainHandlers[T]) onServicePawned(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.ServicePawned)
	return h.publisher.Publish(ctx, servicespb.ServiceAggregateChannel,
		ddd.NewEvent(servicespb.ServicePawnedEvent, &servicespb.ServicePawned{
			Id:     payload.ServiceID,
			UserId: payload.UserID,
		}),
	)
}

func (h domainHandlers[T]) onServiceStockAdjusted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ServiceStockAdjusted)
	return h.publisher.Publish(ctx, servicespb.ServiceAggregateChannel,
		ddd.NewEvent(servicespb.ServiceStockAdjustedEvent, &servicespb.ServiceStockAdjusted{
			Id:       payload.ServiceID,
			OldStock: payload.OldStock,
			NewStock: payload.NewStock,
		}),
	)
}
func (h domainHandlers[T]) onServiceRemoved(ctx context.Context, event ddd.Event) error {
	service := event.Payload().(*domain.Service)
	return h.publisher.Publish(ctx, servicespb.ServiceAggregateChannel,
		ddd.NewEvent(servicespb.ServiceRemovedEvent, &servicespb.ServiceRemoved{
			Id: service.ID(),
		}),
	)
}

func domainOptionsToProto(opts []domain.Option) []*servicespb.Option {
	pbOpts := make([]*servicespb.Option, len(opts))
	for i, o := range opts {
		pbOpts[i] = &servicespb.Option{
			Name:  o.Name,
			Value: o.Value,
			// If your proto uses int64, cast the float64 or int to int64:
			Price: int64(o.Price),
		}
	}
	return pbOpts
}
func domainAttributesToProto(attrs []domain.Attribute) []*servicespb.Attribute {
	pbAttrs := make([]*servicespb.Attribute, len(attrs))
	for i, a := range attrs {
		pbAttrs[i] = &servicespb.Attribute{
			Key:   a.Key,
			Value: a.Value,
		}
	}
	return pbAttrs
}
