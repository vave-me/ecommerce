package commands

import (
	"context"
	"fmt"

	"middleman/erp/internal/domain"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

// FulfillInventoryReservation fulfills an inventory reservation (marks as shipped)
type FulfillInventoryReservation struct {
	ReservationID string
}

// FulfillInventoryReservationHandler handles fulfilling inventory reservations
type FulfillInventoryReservationHandler struct {
	reservations es.AggregateRepository[*domain.InventoryReservation]
	publisher    ddd.EventPublisher[ddd.Event]
}

// NewFulfillInventoryReservationHandler creates a new handler
func NewFulfillInventoryReservationHandler(
	reservations es.AggregateRepository[*domain.InventoryReservation],
	publisher ddd.EventPublisher[ddd.Event],
) FulfillInventoryReservationHandler {
	return FulfillInventoryReservationHandler{
		reservations: reservations,
		publisher:    publisher,
	}
}

// FulfillInventoryReservation fulfills an inventory reservation
func (h FulfillInventoryReservationHandler) FulfillInventoryReservation(ctx context.Context, cmd FulfillInventoryReservation) error {
	// Load the reservation
	reservation, err := h.reservations.Load(ctx, cmd.ReservationID)
	if err != nil {
		return fmt.Errorf("loading reservation: %w", err)
	}
	
	// Fulfill the reservation
	event, err := reservation.Fulfill()
	if err != nil {
		return fmt.Errorf("fulfilling reservation: %w", err)
	}
	
	// Save the reservation
	if err := h.reservations.Save(ctx, reservation); err != nil {
		return fmt.Errorf("saving reservation: %w", err)
	}
	
	// Publish domain event
	return h.publisher.Publish(ctx, event)
}