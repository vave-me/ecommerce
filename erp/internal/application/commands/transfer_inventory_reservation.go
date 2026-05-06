package commands

import (
	"context"
	"fmt"

	"middleman/erp/internal/domain"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

// TransferInventoryReservation transfers a reservation to another warehouse
type TransferInventoryReservation struct {
	ReservationID string
	ToWarehouseID string
}

// TransferInventoryReservationHandler handles transferring inventory reservations
type TransferInventoryReservationHandler struct {
	reservations es.AggregateRepository[*domain.InventoryReservation]
	publisher    ddd.EventPublisher[ddd.Event]
}

// NewTransferInventoryReservationHandler creates a new handler
func NewTransferInventoryReservationHandler(
	reservations es.AggregateRepository[*domain.InventoryReservation],
	publisher ddd.EventPublisher[ddd.Event],
) TransferInventoryReservationHandler {
	return TransferInventoryReservationHandler{
		reservations: reservations,
		publisher:    publisher,
	}
}

// TransferInventoryReservation transfers an inventory reservation to another warehouse
func (h TransferInventoryReservationHandler) TransferInventoryReservation(ctx context.Context, cmd TransferInventoryReservation) error {
	// Load the reservation
	reservation, err := h.reservations.Load(ctx, cmd.ReservationID)
	if err != nil {
		return fmt.Errorf("loading reservation: %w", err)
	}
	
	// Transfer the reservation
	event, err := reservation.Transfer(cmd.ToWarehouseID)
	if err != nil {
		return fmt.Errorf("transferring reservation: %w", err)
	}
	
	// Save the reservation
	if err := h.reservations.Save(ctx, reservation); err != nil {
		return fmt.Errorf("saving reservation: %w", err)
	}
	
	// Publish domain event
	return h.publisher.Publish(ctx, event)
}