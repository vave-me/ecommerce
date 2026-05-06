package handlers

import (
	"context"
	"time"

	"middleman/internal/ddd"
	"middleman/internal/di"
	"middleman/shipping/internal/domain"
)

// catalogHandlers listens to domain events and updates the Catalog DB
type catalogHandlers[T ddd.Event] struct {
	catalog domain.ShippingCatalogRepository
}

var _ ddd.EventHandler[ddd.Event] = (*catalogHandlers[ddd.Event])(nil)

func NewCatalogHandlers(catalog domain.ShippingCatalogRepository) ddd.EventHandler[ddd.Event] {
	return catalogHandlers[ddd.Event]{catalog: catalog}
}

func RegisterCatalogHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.ShipmentCreatedEvent,
		domain.CarrierAssignedEvent,
		domain.ShipmentStartedEvent,
		domain.ShipmentStatusUpdatedEvent,
		domain.ShipmentCancelledEvent,
		domain.ShipmentDeliveredEvent,
		domain.PickupScheduledEvent,
		domain.ShipmentReturnedEvent,
	)
}

func RegisterCatalogHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		catalogHandlers := di.Get(ctx, "catalogHandlers").(ddd.EventHandler[ddd.Event])
		return catalogHandlers.HandleEvent(ctx, event)
	})
	subscriber := container.Get("domainDispatcher").(*ddd.EventDispatcher[ddd.Event])

	RegisterCatalogHandlers(subscriber, handlers)
}

func (h catalogHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.ShipmentCreatedEvent:
		return h.onShipmentCreated(ctx, event)
	case domain.CarrierAssignedEvent:
		return h.onCarrierAssigned(ctx, event)
	case domain.ShipmentStartedEvent:
		return h.onShipmentStarted(ctx, event)
	case domain.ShipmentStatusUpdatedEvent:
		return h.onShipmentStatusUpdated(ctx, event)
	case domain.ShipmentCancelledEvent:
		return h.onShipmentCancelled(ctx, event)
	case domain.ShipmentDeliveredEvent:
		return h.onShipmentDelivered(ctx, event)
	case domain.PickupScheduledEvent:
		return h.onPickupScheduled(ctx, event)
	case domain.ShipmentReturnedEvent:
		return h.onShipmentReturned(ctx, event)
	}
	return nil
}

func (h catalogHandlers[T]) onShipmentCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ShipmentCreated)
	shipment := &domain.CatalogShipment{
		ID:              payload.ID,
		ProductID:       payload.ProductID,
		OrderID:         payload.OrderID,
		BasketID:        payload.BasketID,
		TrackingNumber:  payload.TrackingNumber,
		LabelURL:        payload.LabelUrl,
		SenderName:      payload.SenderName,
		SenderAddress:   payload.SenderAddress,
		ReceiverName:    payload.ReceiverName,
		ReceiverAddress: payload.ReceiverAddress,
		Weight:          payload.Weight,
		Dimensions:      payload.Dimensions,
		ServiceType:     domain.ServiceType(payload.ServiceType),
		Status:          domain.ShipmentStatusCreated,
		CreatedAt:       event.OccurredAt(),
		UpdatedAt:       event.OccurredAt(),
	}
	return h.catalog.AddShipment(ctx, shipment)
}

func (h catalogHandlers[T]) onCarrierAssigned(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.CarrierAssigned)
	return h.catalog.UpdateShipmentStatus(ctx, payload.ShipmentID, domain.ShipmentStatusAssigned, event.OccurredAt())
}

func (h catalogHandlers[T]) onShipmentStarted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ShipmentStarted)
	return h.catalog.UpdateShipmentStatus(ctx, payload.ShipmentID, domain.ShipmentStatusInTransit, event.OccurredAt())
}

func (h catalogHandlers[T]) onShipmentStatusUpdated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ShipmentStatusUpdated)
	return h.catalog.UpdateShipmentStatus(ctx, payload.ShipmentID, domain.ShipmentStatus(payload.Status), event.OccurredAt())
}

func (h catalogHandlers[T]) onShipmentCancelled(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ShipmentCancelled)
	return h.catalog.CancelShipment(ctx, payload.ShipmentID, event.OccurredAt())
}

func (h catalogHandlers[T]) onShipmentDelivered(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ShipmentDelivered)
	return h.catalog.UpdateDeliveryInfo(ctx, payload.ShipmentID, event.OccurredAt())
}

func (h catalogHandlers[T]) onPickupScheduled(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.PickupScheduled)
	// Parse the pickup time
	pickupTime, err := time.Parse(time.RFC3339, payload.PickupTime)
	if err != nil {
		return err
	}
	return h.catalog.UpdatePickupInfo(ctx, payload.ShipmentID, pickupTime)
}

func (h catalogHandlers[T]) onShipmentReturned(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ShipmentReturned)
	return h.catalog.UpdateShipmentStatus(ctx, payload.ShipmentID, domain.ShipmentStatusReturned, event.OccurredAt())
}