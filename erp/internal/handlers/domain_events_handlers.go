package handlers

import (
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"middleman/erp/erppb"
	"middleman/erp/internal/domain"
	"middleman/internal/ddd"
	"time"
)

// Invoice event handlers
func (h domainHandlers[T]) onInvoiceCreated(ctx context.Context, event ddd.Event) error {
	invoice := event.Payload().(*domain.Invoice)
	return h.publisher.Publish(ctx, erppb.InvoiceAggregateChannel,
		ddd.NewEvent(erppb.InvoiceCreatedEvent, &erppb.InvoiceCreated{
			Id:            invoice.ID(),
			InvoiceNumber: invoice.InvoiceNumber,
			OrderId:       invoice.OrderID,
			CustomerId:    invoice.CustomerID,
			CustomerName:  "", // Address doesn't have name field
			CustomerEmail: "", // Address doesn't have email field
			Type:          string(invoice.Type),
			TotalAmount:   int64(invoice.TotalAmount * 100), // Convert to cents
			Currency:      invoice.Currency,
		}),
	)
}

func (h domainHandlers[T]) onInvoiceApproved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.InvoiceApproved)
	// Get invoice ID from event metadata
	var invoiceID string
	if aggEvent, ok := event.(ddd.AggregateEvent); ok {
		invoiceID = aggEvent.AggregateID()
	}
	return h.publisher.Publish(ctx, erppb.InvoiceAggregateChannel,
		ddd.NewEvent(erppb.InvoiceApprovedEvent, &erppb.InvoiceApproved{
			Id:         invoiceID,
			ApprovedBy: payload.ApprovedBy,
			ApprovedAt: timestamppb.New(payload.ApprovedAt),
		}),
	)
}

func (h domainHandlers[T]) onInvoiceSent(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.InvoiceSent)
	// Get invoice ID from event metadata
	var invoiceID string
	if aggEvent, ok := event.(ddd.AggregateEvent); ok {
		invoiceID = aggEvent.AggregateID()
	}
	return h.publisher.Publish(ctx, erppb.InvoiceAggregateChannel,
		ddd.NewEvent(erppb.InvoiceSentEvent, &erppb.InvoiceSent{
			Id:             invoiceID,
			Method:         "email", // Default method
			RecipientEmail: "",      // Would need to get from customer data
			SentAt:         timestamppb.New(payload.SentAt),
		}),
	)
}

func (h domainHandlers[T]) onInvoiceVoided(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.InvoiceVoided)
	// Get invoice ID from event metadata
	var invoiceID string
	if aggEvent, ok := event.(ddd.AggregateEvent); ok {
		invoiceID = aggEvent.AggregateID()
	}
	return h.publisher.Publish(ctx, erppb.InvoiceAggregateChannel,
		ddd.NewEvent(erppb.InvoiceVoidedEvent, &erppb.InvoiceVoided{
			Id:       invoiceID,
			VoidedBy: payload.VoidedBy,
			VoidedAt: timestamppb.New(payload.VoidedAt),
			Reason:   payload.Reason,
		}),
	)
}

func (h domainHandlers[T]) onInvoicePaymentReceived(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.InvoicePaymentReceived)
	// Get invoice ID from event metadata
	var invoiceID string
	if aggEvent, ok := event.(ddd.AggregateEvent); ok {
		invoiceID = aggEvent.AggregateID()
	}
	return h.publisher.Publish(ctx, erppb.InvoiceAggregateChannel,
		ddd.NewEvent(erppb.InvoicePaymentReceivedEvent, &erppb.InvoicePaymentReceived{
			Id:            invoiceID,
			Amount:        int64(payload.Amount * 100), // Convert to cents
			PaymentMethod: payload.PaymentMethod,
			TransactionId: payload.TransactionID,
			PaymentDate:   timestamppb.New(payload.PaymentDate),
		}),
	)
}

func (h domainHandlers[T]) onInvoiceLinkedToERP(ctx context.Context, event ddd.Event) error {
	// This is an internal event, might not need to publish externally
	return nil
}

// Return event handlers
func (h domainHandlers[T]) onReturnCreated(ctx context.Context, event ddd.Event) error {
	// Return created is an internal event, we don't publish it externally
	return nil
}

func (h domainHandlers[T]) onReturnApproved(ctx context.Context, event ddd.Event) error {
	// Return approved is an internal event
	return nil
}

func (h domainHandlers[T]) onReturnProcessed(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ReturnProcessed)
	// Get return ID from event metadata
	var returnID string
	if aggEvent, ok := event.(ddd.AggregateEvent); ok {
		returnID = aggEvent.AggregateID()
	}
	return h.publisher.Publish(ctx, erppb.SyncAggregateChannel,
		ddd.NewEvent(erppb.ReturnProcessedEvent, &erppb.ReturnProcessed{
			Id:          returnID,
			ProcessedAt: timestamppb.New(payload.ProcessedAt),
		}),
	)
}

func (h domainHandlers[T]) onReturnCompleted(ctx context.Context, event ddd.Event) error {
	// Return completed is an internal event
	return nil
}

func (h domainHandlers[T]) onReturnRejected(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ReturnRejected)
	// Get return ID from event metadata
	var returnID string
	if aggEvent, ok := event.(ddd.AggregateEvent); ok {
		returnID = aggEvent.AggregateID()
	}
	return h.publisher.Publish(ctx, erppb.SyncAggregateChannel,
		ddd.NewEvent(erppb.ReturnFailedEvent, &erppb.ReturnFailed{
			Id:     returnID,
			Reason: payload.Reason,
		}),
	)
}

func (h domainHandlers[T]) onReturnItemsRestocked(ctx context.Context, event ddd.Event) error {
	// Internal event, might not need external publishing
	return nil
}

func (h domainHandlers[T]) onReturnLinkedToERP(ctx context.Context, event ddd.Event) error {
	// Internal event, might not need external publishing
	return nil
}

func (h domainHandlers[T]) onReturnSyncFailed(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ReturnSyncFailed)
	// Get return ID from event metadata
	var returnID string
	if aggEvent, ok := event.(ddd.AggregateEvent); ok {
		returnID = aggEvent.AggregateID()
	}
	return h.publisher.Publish(ctx, erppb.SyncAggregateChannel,
		ddd.NewEvent(erppb.ReturnFailedEvent, &erppb.ReturnFailed{
			Id:     returnID,
			Reason: payload.Error,
		}),
	)
}

// Inventory reservation event handlers
func (h domainHandlers[T]) onReservationCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ReservationCreated)
	// Get reservation ID from event metadata
	var reservationID string
	if aggEvent, ok := event.(ddd.AggregateEvent); ok {
		reservationID = aggEvent.AggregateID()
	}
	return h.publisher.Publish(ctx, erppb.SyncAggregateChannel,
		ddd.NewEvent(erppb.InventoryReservedEvent, &erppb.InventoryReserved{
			ReservationId: reservationID,
			ConnectorId:   payload.ConnectorID,
			OrderId:       payload.OrderID,
			Items: []*erppb.ReservedItem{{
				Sku:      payload.SKU,
				Quantity: int64(payload.Quantity),
				Location: payload.WarehouseID,
			}},
			ReservedAt: timestamppb.New(payload.CreatedAt),
		}),
	)
}

func (h domainHandlers[T]) onReservationReleased(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ReservationReleased)
	// Get reservation ID from event metadata
	var reservationID string
	if aggEvent, ok := event.(ddd.AggregateEvent); ok {
		reservationID = aggEvent.AggregateID()
	}
	return h.publisher.Publish(ctx, erppb.SyncAggregateChannel,
		ddd.NewEvent(erppb.InventoryReleasedEvent, &erppb.InventoryReleased{
			ReservationId: reservationID,
			ReleasedAt:    timestamppb.New(payload.ReleasedAt),
		}),
	)
}

func (h domainHandlers[T]) onReservationFulfilled(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ReservationFulfilled)
	// Get reservation ID from event metadata
	var reservationID string
	if aggEvent, ok := event.(ddd.AggregateEvent); ok {
		reservationID = aggEvent.AggregateID()
	}
	return h.publisher.Publish(ctx, erppb.SyncAggregateChannel,
		ddd.NewEvent(erppb.InventoryConfirmedEvent, &erppb.InventoryConfirmed{
			ReservationId: reservationID,
			ConfirmedAt:   timestamppb.New(payload.FulfilledAt),
		}),
	)
}

func (h domainHandlers[T]) onReservationTransferred(ctx context.Context, event ddd.Event) error {
	// Internal event, might not need external publishing
	return nil
}

func (h domainHandlers[T]) onReservationExpired(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ReservationExpired)
	// Get reservation ID from event metadata
	var reservationID string
	if aggEvent, ok := event.(ddd.AggregateEvent); ok {
		reservationID = aggEvent.AggregateID()
	}
	return h.publisher.Publish(ctx, erppb.SyncAggregateChannel,
		ddd.NewEvent(erppb.InventoryReleasedEvent, &erppb.InventoryReleased{
			ReservationId: reservationID,
			ReleasedAt:    timestamppb.New(payload.ExpiredAt),
		}),
	)
}

// Sync event handlers
func (h domainHandlers[T]) onProductsSyncCompleted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ProductsSyncCompleted)
	return h.publisher.Publish(ctx, erppb.SyncAggregateChannel,
		ddd.NewEvent(erppb.ProductsSyncCompletedEvent, &erppb.ProductsSyncCompleted{
			SyncId:         "sync-" + time.Now().Format("20060102150405"),
			ConnectorId:    payload.ERPType,
			ProductsSynced: payload.SuccessCount,
			ProductsFailed: payload.FailedCount,
			CompletedAt:    timestamppb.New(payload.CompletedAt),
		}),
	)
}

func (h domainHandlers[T]) onStockSyncCompleted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.StockSyncCompleted)
	return h.publisher.Publish(ctx, erppb.SyncAggregateChannel,
		ddd.NewEvent(erppb.StockSyncCompletedEvent, &erppb.StockSyncCompleted{
			SyncId:      "sync-" + time.Now().Format("20060102150405"),
			ConnectorId: payload.ERPType,
			StockSynced: payload.SuccessCount,
			StockFailed: payload.FailedCount,
			CompletedAt: timestamppb.New(payload.CompletedAt),
		}),
	)
}

func (h domainHandlers[T]) onPricesSyncCompleted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.PricesSyncCompleted)
	return h.publisher.Publish(ctx, erppb.SyncAggregateChannel,
		ddd.NewEvent(erppb.PricesSyncCompletedEvent, &erppb.PricesSyncCompleted{
			SyncId:       "sync-" + time.Now().Format("20060102150405"),
			ConnectorId:  payload.ERPType,
			PricesSynced: payload.SuccessCount,
			PricesFailed: payload.FailedCount,
			CompletedAt:  timestamppb.New(payload.CompletedAt),
		}),
	)
}

func (h domainHandlers[T]) onCustomersSyncCompleted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.CustomersSyncCompleted)
	return h.publisher.Publish(ctx, erppb.SyncAggregateChannel,
		ddd.NewEvent(erppb.CustomersSyncCompletedEvent, &erppb.CustomersSyncCompleted{
			SyncId:          "sync-" + time.Now().Format("20060102150405"),
			ConnectorId:     payload.ERPType,
			CustomersSynced: payload.SuccessCount,
			CustomersFailed: payload.FailedCount,
			CompletedAt:     timestamppb.New(payload.CompletedAt),
		}),
	)
}

// Order event handlers
func (h domainHandlers[T]) onOrderSentToERP(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.OrderSentToERP)
	return h.publisher.Publish(ctx, erppb.OrderAggregateChannel,
		ddd.NewEvent(erppb.OrderSentToERPEvent, &erppb.OrderSentToERP{
			Id:          "order-" + payload.OrderID,
			OrderId:     payload.OrderID,
			ExternalId:  payload.ERPOrderID,
			ConnectorId: payload.ConnectorID,
			SentAt:      timestamppb.New(payload.SentAt),
		}),
	)
}

func (h domainHandlers[T]) onOrderSyncedFromERP(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.OrderSyncedFromERP)
	return h.publisher.Publish(ctx, erppb.OrderAggregateChannel,
		ddd.NewEvent(erppb.OrderSyncedFromERPEvent, &erppb.OrderSyncedFromERP{
			Id:          "order-" + payload.OrderID,
			OrderId:     payload.OrderID,
			ExternalId:  payload.ERPOrderID,
			ConnectorId: payload.ERPType,
			SyncedAt:    timestamppb.New(payload.SyncedAt),
		}),
	)
}

// Webhook event handlers
func (h domainHandlers[T]) onWebhookReceived(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.WebhookReceived)
	return h.publisher.Publish(ctx, erppb.WebhookAggregateChannel,
		ddd.NewEvent(erppb.WebhookReceivedEvent, &erppb.WebhookReceived{
			Id:          payload.WebhookID,
			ConnectorId: payload.ERPType,
			EventType:   payload.EventType,
			ReceivedAt:  timestamppb.New(payload.ReceivedAt),
		}),
	)
}

func (h domainHandlers[T]) onWebhookProcessed(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.WebhookProcessed)
	return h.publisher.Publish(ctx, erppb.WebhookAggregateChannel,
		ddd.NewEvent(erppb.WebhookProcessedEvent, &erppb.WebhookProcessed{
			Id:          payload.WebhookID,
			ProcessedAt: timestamppb.New(payload.ProcessedAt),
		}),
	)
}

func (h domainHandlers[T]) onWebhookFailed(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.WebhookFailed)
	return h.publisher.Publish(ctx, erppb.WebhookAggregateChannel,
		ddd.NewEvent(erppb.WebhookFailedEvent, &erppb.WebhookFailed{
			Id:       payload.WebhookID,
			Error:    payload.Error,
			FailedAt: timestamppb.New(payload.FailedAt),
		}),
	)
}
