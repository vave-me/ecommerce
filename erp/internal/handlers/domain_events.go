package handlers

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"middleman/erp/internal/domain"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
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
		// Invoice domain events
		domain.InvoiceCreatedEvent,
		domain.InvoiceApprovedEvent,
		domain.InvoiceSentEvent,
		domain.InvoiceVoidedEvent,
		domain.InvoicePaymentReceivedEvent,
		domain.InvoiceLinkedToERPEvent,

		// Return domain events
		domain.ReturnCreatedEvent,
		domain.ReturnApprovedEvent,
		domain.ReturnProcessedEvent,
		domain.ReturnCompletedEvent,
		domain.ReturnRejectedEvent,
		domain.ReturnItemsRestockedEvent,
		domain.ReturnLinkedToERPEvent,
		domain.ReturnSyncFailedEvent,

		// Inventory reservation domain events
		domain.ReservationCreatedEvent,
		domain.ReservationReleasedEvent,
		domain.ReservationFulfilledEvent,
		domain.ReservationTransferredEvent,
		domain.ReservationExpiredEvent,
		
		// Sync events
		domain.ProductsSyncCompletedEvent,
		domain.StockSyncCompletedEvent,
		domain.PricesSyncCompletedEvent,
		domain.CustomersSyncCompletedEvent,
		
		// Order events
		domain.OrderSentToERPEvent,
		domain.OrderSyncedFromERPEvent,
		
		// Webhook events
		domain.WebhookReceivedEvent,
		domain.WebhookProcessedEvent,
		domain.WebhookFailedEvent,
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
	// Invoice events
	case domain.InvoiceCreatedEvent:
		return h.onInvoiceCreated(ctx, event)
	case domain.InvoiceApprovedEvent:
		return h.onInvoiceApproved(ctx, event)
	case domain.InvoiceSentEvent:
		return h.onInvoiceSent(ctx, event)
	case domain.InvoiceVoidedEvent:
		return h.onInvoiceVoided(ctx, event)
	case domain.InvoicePaymentReceivedEvent:
		return h.onInvoicePaymentReceived(ctx, event)
	case domain.InvoiceLinkedToERPEvent:
		return h.onInvoiceLinkedToERP(ctx, event)

	// Return events
	case domain.ReturnCreatedEvent:
		return h.onReturnCreated(ctx, event)
	case domain.ReturnApprovedEvent:
		return h.onReturnApproved(ctx, event)
	case domain.ReturnProcessedEvent:
		return h.onReturnProcessed(ctx, event)
	case domain.ReturnCompletedEvent:
		return h.onReturnCompleted(ctx, event)
	case domain.ReturnRejectedEvent:
		return h.onReturnRejected(ctx, event)
	case domain.ReturnItemsRestockedEvent:
		return h.onReturnItemsRestocked(ctx, event)
	case domain.ReturnLinkedToERPEvent:
		return h.onReturnLinkedToERP(ctx, event)
	case domain.ReturnSyncFailedEvent:
		return h.onReturnSyncFailed(ctx, event)

	// Inventory reservation events
	case domain.ReservationCreatedEvent:
		return h.onReservationCreated(ctx, event)
	case domain.ReservationReleasedEvent:
		return h.onReservationReleased(ctx, event)
	case domain.ReservationFulfilledEvent:
		return h.onReservationFulfilled(ctx, event)
	case domain.ReservationTransferredEvent:
		return h.onReservationTransferred(ctx, event)
	case domain.ReservationExpiredEvent:
		return h.onReservationExpired(ctx, event)
		
	// Sync events
	case domain.ProductsSyncCompletedEvent:
		return h.onProductsSyncCompleted(ctx, event)
	case domain.StockSyncCompletedEvent:
		return h.onStockSyncCompleted(ctx, event)
	case domain.PricesSyncCompletedEvent:
		return h.onPricesSyncCompleted(ctx, event)
	case domain.CustomersSyncCompletedEvent:
		return h.onCustomersSyncCompleted(ctx, event)
		
	// Order events
	case domain.OrderSentToERPEvent:
		return h.onOrderSentToERP(ctx, event)
	case domain.OrderSyncedFromERPEvent:
		return h.onOrderSyncedFromERP(ctx, event)
		
	// Webhook events
	case domain.WebhookReceivedEvent:
		return h.onWebhookReceived(ctx, event)
	case domain.WebhookProcessedEvent:
		return h.onWebhookProcessed(ctx, event)
	case domain.WebhookFailedEvent:
		return h.onWebhookFailed(ctx, event)

	}

	return nil
}
