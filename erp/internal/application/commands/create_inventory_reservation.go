package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"middleman/erp/internal/domain"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

// CreateInventoryReservation creates a new inventory reservation
type CreateInventoryReservation struct {
	ReservationID string
	OrderID       string
	SKU           string
	WarehouseID   string
	Quantity      int
	ExpiresAt     *time.Time
	ConnectorID   string
}

// CreateInventoryReservationHandler handles creating inventory reservations
type CreateInventoryReservationHandler struct {
	reservations es.AggregateRepository[*domain.InventoryReservation]
	publisher    ddd.EventPublisher[ddd.Event]
}

// NewCreateInventoryReservationHandler creates a new handler
func NewCreateInventoryReservationHandler(
	reservations es.AggregateRepository[*domain.InventoryReservation],
	publisher ddd.EventPublisher[ddd.Event],
) CreateInventoryReservationHandler {
	return CreateInventoryReservationHandler{
		reservations: reservations,
		publisher:    publisher,
	}
}

// CreateInventoryReservation creates a new inventory reservation
func (h CreateInventoryReservationHandler) CreateInventoryReservation(ctx context.Context, cmd CreateInventoryReservation) error {
	// Generate reservation ID if not provided
	if cmd.ReservationID == "" {
		cmd.ReservationID = uuid.New().String()
	}
	
	// Load the reservation aggregate
	reservation, err := h.reservations.Load(ctx, cmd.ReservationID)
	if err != nil {
		return fmt.Errorf("loading reservation: %w", err)
	}
	
	// Create the reservation
	event, err := reservation.CreateReservation(
		cmd.OrderID,
		cmd.SKU,
		cmd.WarehouseID,
		cmd.Quantity,
		cmd.ExpiresAt,
		cmd.ConnectorID,
	)
	if err != nil {
		return fmt.Errorf("creating reservation: %w", err)
	}
	
	// Save the reservation
	if err := h.reservations.Save(ctx, reservation); err != nil {
		return fmt.Errorf("saving reservation: %w", err)
	}
	
	// Publish domain event
	return h.publisher.Publish(ctx, event)
}