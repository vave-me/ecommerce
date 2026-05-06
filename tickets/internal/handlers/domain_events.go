package handlers

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/tickets/internal/domain"
	"middleman/tickets/ticketspb"

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
		domain.TicketAddedEvent,
		domain.TicketUpdatedEvent,
		domain.TicketRebrandedEvent,
		domain.TicketPriceIncreasedEvent,
		domain.TicketPriceDecreasedEvent,
		domain.TicketStockAdjustedEvent,
		domain.TicketRemovedEvent,
		domain.TicketArchivedEvent,
		domain.TicketSoldEvent,
		domain.TicketLeasedEvent,
		domain.TicketPawnedEvent,
		domain.TicketThumbnailAddedEvent,
		domain.TicketThumbnailUpdatedEvent,

		// Variant events
		domain.VariantAddedEvent,
		domain.VariantPriceDecreasedEvent,
		domain.VariantPriceIncreasedEvent,
		domain.VariantStockAdjustedEvent,
		domain.VariantArchivedEvent,
		domain.VariantRemovedEvent,
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

	case domain.TicketAddedEvent:
		return h.onTicketAdded(ctx, event)
	case domain.TicketUpdatedEvent:
		return h.onTicketUpdated(ctx, event)
	case domain.TicketRebrandedEvent:
		return h.onTicketRebranded(ctx, event)
	case domain.TicketPriceIncreasedEvent:
		return h.onTicketPriceIncreased(ctx, event)
	case domain.TicketPriceDecreasedEvent:
		return h.onTicketPriceDecreased(ctx, event)
	case domain.TicketStockAdjustedEvent:
		return h.onTicketStockAdjusted(ctx, event)
	case domain.TicketArchivedEvent:
		return h.onTicketArchived(ctx, event)
	case domain.TicketRemovedEvent:
		return h.onTicketRemoved(ctx, event)
	case domain.TicketSoldEvent:
		return h.onTicketSold(ctx, event)
	case domain.TicketLeasedEvent:
		return h.onTicketLeased(ctx, event)
	case domain.TicketPawnedEvent:
		return h.onTicketPawned(ctx, event)

	// -------------------------------------------------------------------------
	// VARIANT EVENTS
	// -------------------------------------------------------------------------
	case domain.VariantAddedEvent:
		return h.onVariantAdded(ctx, event)
	case domain.VariantPriceDecreasedEvent:
		return h.onVariantPriceDecreased(ctx, event)
	case domain.VariantPriceIncreasedEvent:
		return h.onVariantPriceIncreased(ctx, event)
	case domain.VariantStockAdjustedEvent:
		return h.onVariantStockAdjusted(ctx, event)
	case domain.VariantArchivedEvent:
		return h.onVariantArchived(ctx, event)
	case domain.VariantRemovedEvent:
		return h.onVariantRemoved(ctx, event)
	}
	return nil

}

func (h domainHandlers[T]) onTicketAdded(ctx context.Context, event ddd.Event) error {
	ticket := event.Payload().(*domain.Ticket)
	return h.publisher.Publish(ctx, ticketspb.TicketAggregateChannel,
		ddd.NewEvent(ticketspb.TicketAddedEvent, &ticketspb.TicketAdded{
			Id:               ticket.ID(),
			Name:             ticket.Name,
			Description:      ticket.Description,
			BasePrice:        ticket.BasePrice,
			UserSellerId:     ticket.UserSellerID,
			CategoryId:       ticket.CategoryID,
			CategorySlug:     ticket.CategorySlug,
			Brand:            ticket.Brand,
			Condition:        ticket.Condition.String(),
			Model:            ticket.Model,
			Tags:             ticket.Tags,
			ManageStocks:     ticket.ManageStock,
			Stock:            ticket.Stock,
			Sku:              ticket.SKU,
			Attributes:       domainAttributesToProto(ticket.Attributes),
			Weight:           ticket.Weight,
			Height:           ticket.Height,
			Width:            ticket.Width,
			Depth:            ticket.Depth,
			Status:           ticket.Status.String(),
			Negotiable:       ticket.Negotiable,
			UserType:         ticket.UserType.String(),
			MiddlemanService: ticket.MiddlemanService,
			ShippingCost:     ticket.ShippingCost,
			HasVariants:      ticket.HasVariants,
			Options:          domainOptionsToProto(ticket.Options),
			Lat:              float32(ticket.Lat),
			Lng:              float32(ticket.Lng),
		}),
	)
}
func (h domainHandlers[T]) onTicketUpdated(ctx context.Context, event ddd.Event) error {
	ticket := event.Payload().(*domain.Ticket)
	return h.publisher.Publish(ctx, ticketspb.TicketAggregateChannel,
		ddd.NewEvent(ticketspb.TicketUpdatedEvent, &ticketspb.TicketUpdated{
			Id:               ticket.ID(),
			Name:             ticket.Name,
			Description:      ticket.Description,
			BasePrice:        ticket.BasePrice,
			UserSellerId:     ticket.UserSellerID,
			CategoryId:       ticket.CategoryID,
			CategorySlug:     ticket.CategorySlug,
			Brand:            ticket.Brand,
			Condition:        ticket.Condition.String(),
			Model:            ticket.Model,
			Tags:             ticket.Tags,
			ManageStocks:     ticket.ManageStock,
			Stock:            ticket.Stock,
			Sku:              ticket.SKU,
			Attributes:       domainAttributesToProto(ticket.Attributes),
			Weight:           ticket.Weight,
			Height:           ticket.Height,
			Width:            ticket.Width,
			Depth:            ticket.Depth,
			Status:           ticket.Status.String(),
			Negotiable:       ticket.Negotiable,
			UserType:         ticket.UserType.String(),
			MiddlemanService: ticket.MiddlemanService,
			ShippingCost:     ticket.ShippingCost,
			HasVariants:      ticket.HasVariants,
			Options:          domainOptionsToProto(ticket.Options),
		}),
	)
}

func (h domainHandlers[T]) onTicketRebranded(ctx context.Context, event ddd.Event) error {
	ticket := event.Payload().(*domain.Ticket)
	return h.publisher.Publish(ctx, ticketspb.TicketAggregateChannel,
		ddd.NewEvent(ticketspb.TicketRebrandedEvent, &ticketspb.TicketRebranded{
			Id:          ticket.ID(),
			Name:        ticket.Name,
			Description: ticket.Description,
		}),
	)
}

func (h domainHandlers[T]) onTicketPriceIncreased(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TicketPriceIncreased)
	return h.publisher.Publish(ctx, ticketspb.TicketAggregateChannel,
		ddd.NewEvent(ticketspb.TicketPriceIncreasedEvent, &ticketspb.TicketPriceIncreased{
			Id:       payload.TicketID,
			OldPrice: payload.OldPrice,
			NewPrice: payload.NewPrice,
		}),
	)
}

func (h domainHandlers[T]) onTicketArchived(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.TicketArchived)
	return h.publisher.Publish(ctx, ticketspb.TicketAggregateChannel,
		ddd.NewEvent(ticketspb.TicketArchivedEvent, &ticketspb.TicketArchived{
			Id:           payload.TicketID,
			UserSellerId: payload.UserSellerID,
		}),
	)
}

func (h domainHandlers[T]) onTicketSold(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.TicketSold)
	return h.publisher.Publish(ctx, ticketspb.TicketAggregateChannel,
		ddd.NewEvent(ticketspb.TicketSoldEvent, &ticketspb.TicketSold{
			Id:           payload.TicketID,
			UserSellerId: payload.UserSellerID,

			// BuyerId:   payload.BuyerID, if your proto has it
		}),
	)
}

func (h domainHandlers[T]) onTicketLeased(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.TicketLeased)
	return h.publisher.Publish(ctx, ticketspb.TicketAggregateChannel,
		ddd.NewEvent(ticketspb.TicketLeasedEvent, &ticketspb.TicketLeased{
			Id:           payload.TicketID,
			UserSellerId: payload.UserSellerID,

			// LesseeId: payload.LesseeID, if needed
		}),
	)
}

func (h domainHandlers[T]) onTicketPriceDecreased(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TicketPriceDecreased)
	return h.publisher.Publish(ctx, ticketspb.TicketAggregateChannel,
		ddd.NewEvent(ticketspb.TicketPriceDecreasedEvent, &ticketspb.TicketPriceDecreased{
			Id:       payload.TicketID,
			OldPrice: payload.OldPrice,
			NewPrice: payload.NewPrice,
		}),
	)
}
func (h domainHandlers[T]) onTicketPawned(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.TicketPawned)
	return h.publisher.Publish(ctx, ticketspb.TicketAggregateChannel,
		ddd.NewEvent(ticketspb.TicketPawnedEvent, &ticketspb.TicketPawned{
			Id:           payload.TicketID,
			UserSellerId: payload.UserSellerID,
		}),
	)
}

func (h domainHandlers[T]) onTicketStockAdjusted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TicketStockAdjusted)
	return h.publisher.Publish(ctx, ticketspb.TicketAggregateChannel,
		ddd.NewEvent(ticketspb.TicketStockAdjustedEvent, &ticketspb.TicketStockAdjusted{
			Id:       payload.TicketID,
			OldStock: payload.OldStock,
			NewStock: payload.NewStock,
		}),
	)
}
func (h domainHandlers[T]) onTicketRemoved(ctx context.Context, event ddd.Event) error {
	ticket := event.Payload().(*domain.Ticket)
	return h.publisher.Publish(ctx, ticketspb.TicketAggregateChannel,
		ddd.NewEvent(ticketspb.TicketRemovedEvent, &ticketspb.TicketRemoved{
			Id: ticket.ID(),
		}),
	)
}

func (h domainHandlers[T]) onTicketThumbnailAdded(ctx context.Context, event ddd.Event) error {
	ticket := event.Payload().(*domain.Ticket)
	return h.publisher.Publish(ctx, ticketspb.TicketAggregateChannel,
		ddd.NewEvent(ticketspb.TicketThumbnailAddedEvent, &ticketspb.TicketRemoved{
			Id: ticket.ID(),
		}),
	)
}

func (h domainHandlers[T]) onVariantAdded(ctx context.Context, event ddd.Event) error {
	variant := event.Payload().(*domain.Variant)
	// Publish to ticketspb.VariantAggregateChannel or TicketAggregateChannel
	return h.publisher.Publish(ctx, ticketspb.VariantAggregateChannel,
		ddd.NewEvent(ticketspb.VariantAddedEvent, &ticketspb.VariantAdded{
			Id:           variant.ID(),
			TicketId:     variant.TicketID,
			Sku:          variant.SKU,
			Barcode:      variant.Barcode,
			VariantPrice: variant.VariantPrice,
			Stock:        variant.Stock,
			Weight:       variant.Weight,
			Height:       variant.Height,
			Width:        variant.Width,
			Depth:        variant.Depth,
			IsAvailable:  variant.IsAvailable,
			// Add other fields as needed, e.g. currencyCode, attributes, etc.
		}),
	)
}

func (h domainHandlers[T]) onVariantPriceIncreased(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.VariantPriceIncreased)
	return h.publisher.Publish(ctx, ticketspb.VariantAggregateChannel,
		ddd.NewEvent(ticketspb.VariantPriceIncreasedEvent, &ticketspb.VariantPriceIncreased{
			Id:       payload.VariantID,
			OldPrice: payload.OldPrice,
			NewPrice: payload.NewPrice,
		}),
	)
}

func (h domainHandlers[T]) onVariantPriceDecreased(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.VariantPriceDecreased)
	return h.publisher.Publish(ctx, ticketspb.VariantAggregateChannel,
		ddd.NewEvent(ticketspb.VariantPriceDecreasedEvent, &ticketspb.VariantPriceDecreased{
			Id:       payload.VariantID,
			OldPrice: payload.OldPrice,
			NewPrice: payload.NewPrice,
		}),
	)
}

func (h domainHandlers[T]) onVariantStockAdjusted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.VariantStockAdjusted)
	return h.publisher.Publish(ctx, ticketspb.VariantAggregateChannel,
		ddd.NewEvent(ticketspb.VariantStockAdjustedEvent, &ticketspb.VariantStockAdjusted{
			Id:       payload.VariantID,
			OldStock: payload.OldStock,
			NewStock: payload.NewStock,
		}),
	)
}

func (h domainHandlers[T]) onVariantArchived(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.VariantArchived)
	return h.publisher.Publish(ctx, ticketspb.VariantAggregateChannel,
		ddd.NewEvent(ticketspb.VariantArchivedEvent, &ticketspb.VariantArchived{
			Id: payload.VariantID,
		}),
	)
}

func (h domainHandlers[T]) onVariantRemoved(ctx context.Context, event ddd.Event) error {
	variant := event.Payload().(*domain.Variant)
	return h.publisher.Publish(ctx, ticketspb.VariantAggregateChannel,
		ddd.NewEvent(ticketspb.VariantRemovedEvent, &ticketspb.VariantRemoved{
			Id: variant.ID(),
		}),
	)
}
func domainOptionsToProto(opts []domain.Option) []*ticketspb.Option {
	pbOpts := make([]*ticketspb.Option, len(opts))
	for i, o := range opts {
		pbOpts[i] = &ticketspb.Option{
			Name:  o.Name,
			Value: o.Value,
			// If your proto uses int64, cast the float64 or int to int64:
			Price: int64(o.Price),
		}
	}
	return pbOpts
}
func domainAttributesToProto(attrs []domain.Attribute) []*ticketspb.Attribute {
	pbAttrs := make([]*ticketspb.Attribute, len(attrs))
	for i, a := range attrs {
		pbAttrs[i] = &ticketspb.Attribute{
			Key:   a.Key,
			Value: a.Value,
		}
	}
	return pbAttrs
}
