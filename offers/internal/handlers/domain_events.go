package handlers

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/offers/internal/domain"
	"middleman/offers/offerspb"
	"time"
)

// domainHandlers is a generic handler that processes domain events (T ddd.Event)
// and publishes integration events via an am.EventPublisher
type domainHandlers[T ddd.Event] struct {
	publisher am.EventPublisher
}

var _ ddd.EventHandler[ddd.Event] = (*domainHandlers[ddd.Event])(nil)

// NewDomainEventHandlers creates a new domain handler that
// listens to domain events and publishes them to an integration channel
func NewDomainEventHandlers(publisher am.EventPublisher) ddd.EventHandler[ddd.Event] {
	return &domainHandlers[ddd.Event]{
		publisher: publisher,
	}
}

// RegisterDomainEventHandlers subscribes to all domain events we intend to handle
func RegisterDomainEventHandlers(
	subscriber ddd.EventSubscriber[ddd.Event],
	handlers ddd.EventHandler[ddd.Event],
) {
	subscriber.Subscribe(handlers,
		// Offer events
		domain.OfferCreatedEvent,
		domain.OfferActivatedEvent,
		domain.OfferClosedEvent,
		domain.OfferAcceptedEvent,

		domain.BuyNowCreatedEvent,
		domain.BuyNowConfirmedEvent,

		domain.LeaseCreatedEvent, // might be same or different from "LeaseAdded"
		domain.LeaseDefaultedEvent,
		domain.LeaseEndedEvent,

		// BuyBack events
		domain.BuyBackCreatedEvent,
		domain.BuyBackCanceledEvent,
		domain.BuyBackExpiredEvent,
		domain.BuyBackRedeemedEvent,

		domain.ReservationCreatedEvent,
		domain.ReservationCanceledEvent,
		domain.ReservationExpiredEvent,
		domain.ReservationRedeemedEvent,
	)
}

// HandleEvent routes domain events to the appropriate onXxx handler
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

	// -------------------
	// Offer events
	// -------------------
	case domain.OfferCreatedEvent:
		return h.onOfferCreated(ctx, event)

	case domain.OfferActivatedEvent:
		return h.onOfferActivated(ctx, event)

	case domain.OfferClosedEvent:
		return h.onOfferClosed(ctx, event)

	case domain.OfferAcceptedEvent:
		return h.onOfferAccepted(ctx, event)

	// -------------------
	// BuyNow events
	// -------------------
	case domain.BuyNowCreatedEvent:
		return h.onBuyNowCreated(ctx, event)

	case domain.BuyNowConfirmedEvent:
		return h.onBuyNowConfirmed(ctx, event)

	// -------------------
	// Lease events
	// -------------------

	case domain.LeaseCreatedEvent:
		return h.onLeaseCreated(ctx, event)

	case domain.LeaseDefaultedEvent:
		return h.onLeaseDefaulted(ctx, event)

	case domain.LeaseEndedEvent:
		return h.onLeaseEnded(ctx, event)

	// -------------------
	// BuyBack events
	// -------------------
	case domain.BuyBackCreatedEvent:
		return h.onBuyBackCreated(ctx, event)

	case domain.BuyBackCanceledEvent:
		return h.onBuyBackCanceled(ctx, event)

	case domain.BuyBackExpiredEvent:
		return h.onBuyBackExpired(ctx, event)

	case domain.BuyBackRedeemedEvent:
		return h.onBuyBackRedeemed(ctx, event)

	case domain.ReservationCreatedEvent:
		return h.onReservationCreated(ctx, event)

	case domain.ReservationCanceledEvent:
		return h.onReservationCanceled(ctx, event)

	case domain.ReservationExpiredEvent:
		return h.onReservationExpired(ctx, event)

	case domain.ReservationRedeemedEvent:
		return h.onReservationRedeemed(ctx, event)

	}
	return nil
}

// ----------------------------------------------
// Offer Event Handlers
// ----------------------------------------------

func (h domainHandlers[T]) onOfferCreated(ctx context.Context, event ddd.Event) error {
	offer := event.Payload().(*domain.Offer)
	// Publish integration event
	return h.publisher.Publish(ctx, offerspb.OfferAggregateChannel,
		ddd.NewEvent(offerspb.OfferCreatedEvent, &offerspb.OfferCreated{
			Id:             offer.ID(),
			UserSellerId:   offer.UserSellerID,
			UserCustomerId: offer.UserCustomerID,
			ProductId:      offer.ProductID,
		}),
	)
}

func (h domainHandlers[T]) onOfferActivated(ctx context.Context, event ddd.Event) error {
	// domain.OfferActivated is presumably a struct, cast & handle
	// (assuming event.Payload() is domain.Offer or domain.OfferActivated).
	// For demonstration, let's assume it's domain.OfferActivated struct:
	//   type OfferActivated struct { OfferID string }
	activated := event.Payload().(*domain.OfferActivated)
	// Publish integration event
	return h.publisher.Publish(ctx, offerspb.OfferAggregateChannel,
		ddd.NewEvent(offerspb.OfferActivatedEvent, &offerspb.OfferActivated{
			OfferId: activated.ProductID,
		}),
	)
}

func (h domainHandlers[T]) onOfferClosed(ctx context.Context, event ddd.Event) error {
	closed := event.Payload().(*domain.OfferClosed)
	return h.publisher.Publish(ctx, offerspb.OfferAggregateChannel,
		ddd.NewEvent(offerspb.OfferClosedEvent, &offerspb.OfferClosed{
			OfferId: closed.OfferID,
			Reason:  closed.Reason,
		}),
	)
}

func (h domainHandlers[T]) onOfferAccepted(ctx context.Context, event ddd.Event) error {
	accepted := event.Payload().(*domain.OfferAccepted)
	return h.publisher.Publish(ctx, offerspb.OfferAggregateChannel,
		ddd.NewEvent(offerspb.OfferAcceptedEvent, &offerspb.OfferAccepted{
			OfferId:        accepted.OfferID,
			UserCustomerId: accepted.UserCustomerID,
		}),
	)
}

// ----------------------------------------------
// BuyNow Event Handlers
// ----------------------------------------------

func (h domainHandlers[T]) onBuyNowCreated(ctx context.Context, event ddd.Event) error {
	// domain.BuyNowCreated
	created := event.Payload().(*domain.BuyNow)
	// Then build an integration event
	return h.publisher.Publish(ctx, offerspb.BuyNowAggregateChannel,
		ddd.NewEvent(offerspb.BuyNowCreatedEvent, &offerspb.BuyNowCreated{
			Id:         created.ID(),
			OfferId:    created.OfferID,
			FinalPrice: created.FinalPrice,
			// etc...
		}),
	)
}

func (h domainHandlers[T]) onBuyNowConfirmed(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.BuyNow)
	return h.publisher.Publish(ctx, offerspb.BuyNowAggregateChannel,
		ddd.NewEvent(offerspb.BuyNowConfirmedEvent, &offerspb.BuyNowConfirmed{
			BuyNowId: payload.ID(),
			// maybe confirm date/time
		}),
	)
}

// ----------------------------------------------
// Lease Event Handlers
// ----------------------------------------------

func (h domainHandlers[T]) onLeaseAdded(ctx context.Context, event ddd.Event) error {
	added := event.Payload().(*domain.Lease)
	return h.publisher.Publish(ctx, offerspb.LeaseAggregateChannel,
		ddd.NewEvent(offerspb.LeaseAddedEvent, &offerspb.LeaseCreated{
			Id:              added.ID(),
			OfferId:         added.OfferID,
			MonthlyPrice:    added.MonthlyPrice,
			LeaseTermMonths: added.LeaseTermMonths,
			// etc...
		}),
	)
}

func (h domainHandlers[T]) onLeaseCreated(ctx context.Context, event ddd.Event) error {
	created := event.Payload().(*domain.Lease)
	return h.publisher.Publish(ctx, offerspb.LeaseAggregateChannel,
		ddd.NewEvent(offerspb.LeaseCreatedEvent, &offerspb.LeaseCreated{
			Id:              created.ID(),
			OfferId:         created.OfferID,
			MonthlyPrice:    created.MonthlyPrice,
			LeaseTermMonths: created.LeaseTermMonths,
			// etc...
		}),
	)
}

func (h domainHandlers[T]) onLeaseDefaulted(ctx context.Context, event ddd.Event) error {
	def := event.Payload().(*domain.Lease)
	return h.publisher.Publish(ctx, offerspb.LeaseAggregateChannel,
		ddd.NewEvent(offerspb.LeaseDefaultedEvent, &offerspb.LeaseDefaulted{
			LeaseId: def.ID(),
		}),
	)
}

func (h domainHandlers[T]) onLeaseEnded(ctx context.Context, event ddd.Event) error {
	ended := event.Payload().(*domain.Lease)
	return h.publisher.Publish(ctx, offerspb.LeaseAggregateChannel,
		ddd.NewEvent(offerspb.LeaseEndedEvent, &offerspb.LeaseEnded{
			LeaseId: ended.ID(),
		}),
	)
}

// ----------------------------------------------
// BuyBack Event Handlers
// ----------------------------------------------

func (h domainHandlers[T]) onBuyBackCreated(ctx context.Context, event ddd.Event) error {
	created := event.Payload().(*domain.BuyBack)
	return h.publisher.Publish(ctx, offerspb.BuyBackAggregateChannel,
		ddd.NewEvent(offerspb.BuyBackCreatedEvent, &offerspb.BuyBackCreated{
			Id:              created.ID(),
			OfferId:         created.OfferID,
			LockedPrice:     created.LockedPrice,
			RedemptionPrice: created.RedemptionFee,
			// ...
		}),
	)
}

func (h domainHandlers[T]) onBuyBackCanceled(ctx context.Context, event ddd.Event) error {
	canceled := event.Payload().(*domain.BuyBack)
	return h.publisher.Publish(ctx, offerspb.BuyBackAggregateChannel,
		ddd.NewEvent(offerspb.BuyBackCanceledEvent, &offerspb.BuyBackCanceled{
			BuyBackId: canceled.ID(),
		}),
	)
}

func (h domainHandlers[T]) onBuyBackExpired(ctx context.Context, event ddd.Event) error {
	expired := event.Payload().(*domain.BuyBack)
	return h.publisher.Publish(ctx, offerspb.BuyBackAggregateChannel,
		ddd.NewEvent(offerspb.BuyBackExpiredEvent, &offerspb.BuyBackExpired{
			BuyBackId: expired.ID(),
		}),
	)
}

func (h domainHandlers[T]) onBuyBackRedeemed(ctx context.Context, event ddd.Event) error {
	redeemed := event.Payload().(*domain.BuyBack)
	return h.publisher.Publish(ctx, offerspb.BuyBackAggregateChannel,
		ddd.NewEvent(offerspb.BuyBackRedeemedEvent, &offerspb.BuyBackRedeemed{
			BuyBackId: redeemed.ID(),
		}),
	)
}

func (h domainHandlers[T]) onReservationCreated(ctx context.Context, event ddd.Event) error {
	created := event.Payload().(*domain.Reservation)
	return h.publisher.Publish(ctx, offerspb.ReservationAggregateChannel,
		ddd.NewEvent(offerspb.ReservationCreatedEvent, &offerspb.ReservationCreated{
			Id:             created.ID(),
			OfferId:        created.OfferID,
			LockedPrice:    created.LockedPrice,
			ReservationFee: created.ReservationFee,
			// ...
		}),
	)
}

func (h domainHandlers[T]) onReservationCanceled(ctx context.Context, event ddd.Event) error {
	canceled := event.Payload().(*domain.Reservation)
	return h.publisher.Publish(ctx, offerspb.ReservationAggregateChannel,
		ddd.NewEvent(offerspb.ReservationCanceledEvent, &offerspb.ReservationCanceled{
			ReservationId: canceled.ID(),
		}),
	)
}

func (h domainHandlers[T]) onReservationExpired(ctx context.Context, event ddd.Event) error {
	expired := event.Payload().(*domain.Reservation)
	return h.publisher.Publish(ctx, offerspb.ReservationAggregateChannel,
		ddd.NewEvent(offerspb.ReservationExpiredEvent, &offerspb.ReservationExpired{
			ReservationId: expired.ID(),
		}),
	)
}

func (h domainHandlers[T]) onReservationRedeemed(ctx context.Context, event ddd.Event) error {
	redeemed := event.Payload().(*domain.Reservation)
	return h.publisher.Publish(ctx, offerspb.ReservationAggregateChannel,
		ddd.NewEvent(offerspb.ReservationRedeemedEvent, &offerspb.ReservationRedeemed{
			ReservationId: redeemed.ID(),
		}),
	)
}
