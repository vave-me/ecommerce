package commands

import (
	"context"
	"fmt"

	"middleman/erp/internal/domain"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

// ReleaseInventoryReservation releases an existing inventory reservation
type ReleaseInventoryReservation struct {
	ReservationID string
	Reason        string
}

// ReleaseInventoryReservationHandler handles releasing inventory reservations
type ReleaseInventoryReservationHandler struct {
	reservations es.AggregateRepository[*domain.InventoryReservation]
	publisher    ddd.EventPublisher[ddd.Event]
}

// NewReleaseInventoryReservationHandler creates a new handler
func NewReleaseInventoryReservationHandler(
	reservations es.AggregateRepository[*domain.InventoryReservation],
	publisher ddd.EventPublisher[ddd.Event],
) ReleaseInventoryReservationHandler {
	return ReleaseInventoryReservationHandler{
		reservations: reservations,
		publisher:    publisher,
	}
}

// ReleaseInventoryReservation releases an inventory reservation
func (h ReleaseInventoryReservationHandler) ReleaseInventoryReservation(ctx context.Context, cmd ReleaseInventoryReservation) error {
	// Load the reservation
	reservation, err := h.reservations.Load(ctx, cmd.ReservationID)
	if err != nil {
		return fmt.Errorf("loading reservation: %w", err)
	}
	
	// Release the reservation
	event, err := reservation.Release(cmd.Reason)
	if err != nil {
		return fmt.Errorf("releasing reservation: %w", err)
	}
	
	// Save the reservation
	if err := h.reservations.Save(ctx, reservation); err != nil {
		return fmt.Errorf("saving reservation: %w", err)
	}
	
	// Publish domain event
	return h.publisher.Publish(ctx, event)
}