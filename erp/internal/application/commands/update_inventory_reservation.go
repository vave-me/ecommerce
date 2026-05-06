// DEPRECATED: This file is deprecated in favor of granular reservation commands.
// Use the following commands instead:
// - CreateInventoryReservation: Creates a new inventory reservation
// - ReleaseInventoryReservation: Releases an existing reservation
// - FulfillInventoryReservation: Marks a reservation as fulfilled (shipped)
// - TransferInventoryReservation: Transfers a reservation to another warehouse
//
// The new approach uses the InventoryReservation aggregate with proper event sourcing.

package commands

import (
	"context"
	"fmt"
	"middleman/internal/erp"
	"time"

	"middleman/erp/internal/domain"
)

// UpdateInventoryReservation command updates inventory reservations for order fulfillment
type UpdateInventoryReservation struct {
	ConnectorID  string
	OrderID      string
	Action       ReservationAction
	Reservations []InventoryReservation
}

// ReservationAction represents the type of reservation action
type ReservationAction string

const (
	ReservationActionReserve  ReservationAction = "reserve"
	ReservationActionRelease  ReservationAction = "release"
	ReservationActionFulfill  ReservationAction = "fulfill"
	ReservationActionTransfer ReservationAction = "transfer"
)

// InventoryReservation represents a single inventory reservation
type InventoryReservation struct {
	SKU                 string
	WarehouseID         string
	Quantity            int
	ReservationID       string
	ExpiresAt           *time.Time
	TransferToWarehouse string // For transfer action
}

// UpdateInventoryReservationHandler handles inventory reservation updates
type UpdateInventoryReservationHandler struct {
	registry   erp.ConnectorRegistry
	repository domain.OrderSyncRepository
}

// NewUpdateInventoryReservationHandler creates a new handler
func NewUpdateInventoryReservationHandler(
	registry erp.ConnectorRegistry,
	repository domain.OrderSyncRepository,
) UpdateInventoryReservationHandler {
	return UpdateInventoryReservationHandler{
		registry:   registry,
		repository: repository,
	}
}

// UpdateInventoryReservation updates inventory reservations (DEPRECATED)
func (h UpdateInventoryReservationHandler) UpdateInventoryReservation(ctx context.Context, cmd UpdateInventoryReservation) error {
	// Validate command
	if cmd.OrderID == "" {
		return fmt.Errorf("order ID is required")
	}
	if len(cmd.Reservations) == 0 {
		return fmt.Errorf("at least one reservation is required")
	}
	if err := validateReservationAction(cmd.Action); err != nil {
		return err
	}

	// Get connector
	connector, err := h.registry.GetConnector(cmd.ConnectorID)
	if err != nil {
		return fmt.Errorf("getting connector: %w", err)
	}

	// Create order sync record for tracking
	orderSync := &domain.OrderSync{
		ID:          generateReservationSyncID(),
		ConnectorID: cmd.ConnectorID,
		OrderID:     cmd.OrderID,
		Direction:   domain.DirectionOutbound,
		Status:      domain.OrderSyncStatusPending,
		AttemptedAt: time.Now(),
		Payload: map[string]interface{}{
			"action":       cmd.Action,
			"reservations": cmd.Reservations,
		},
	}

	if err := h.repository.Create(ctx, orderSync); err != nil {
		return fmt.Errorf("creating order sync record: %w", err)
	}

	// Process each reservation
	successCount := 0
	failedReservations := []string{}

	for _, reservation := range cmd.Reservations {
		if err := validateReservation(reservation, cmd.Action); err != nil {
			failedReservations = append(failedReservations,
				fmt.Sprintf("%s: %v", reservation.SKU, err))
			continue
		}

		// Create inventory adjustment based on action
		adjustment := createInventoryAdjustment(cmd.OrderID, reservation, cmd.Action)

		// Send to ERP
		if err := connector.UpdateInventory(ctx, []*erp.InventoryAdjustment{adjustment}); err != nil {
			failedReservations = append(failedReservations,
				fmt.Sprintf("%s: %v", reservation.SKU, err))
			continue
		}

		// Publish reservation event
		//TODO product grpc repostiory call

		// Publishing events removed - handled separately
		// if err := h.publisher.Publish(ctx, event); err != nil {
		// 	// Log but don't fail the command
		// 	failedReservations = append(failedReservations,
		// 		fmt.Sprintf("%s: failed to publish event: %v", reservation.SKU, err))
		// } else {
		successCount++
	}

	// Update order sync status
	if len(failedReservations) > 0 {
		orderSync.Status = domain.OrderSyncStatusFailed
		orderSync.Error = fmt.Sprintf("failed reservations: %v", failedReservations)
	} else {
		orderSync.Status = domain.OrderSyncStatusCompleted
	}
	orderSync.CompletedAt = ptrTime(time.Now())

	if err := h.repository.Update(ctx, orderSync); err != nil {
		return fmt.Errorf("updating order sync record: %w", err)
	}

	if len(failedReservations) > 0 {
		return fmt.Errorf("reservation update partially failed: %d succeeded, %d failed",
			successCount, len(failedReservations))
	}

	return nil
}

func validateReservationAction(action ReservationAction) error {
	switch action {
	case ReservationActionReserve, ReservationActionRelease,
		ReservationActionFulfill, ReservationActionTransfer:
		return nil
	default:
		return fmt.Errorf("invalid reservation action: %s", action)
	}
}

func validateReservation(res InventoryReservation, action ReservationAction) error {
	if res.SKU == "" {
		return fmt.Errorf("SKU is required")
	}
	if res.WarehouseID == "" {
		return fmt.Errorf("warehouse ID is required")
	}
	if res.Quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}
	if action == ReservationActionTransfer && res.TransferToWarehouse == "" {
		return fmt.Errorf("transfer destination warehouse is required")
	}
	if action == ReservationActionRelease && res.ReservationID == "" {
		return fmt.Errorf("reservation ID is required for release action")
	}
	return nil
}

func createInventoryAdjustment(orderID string, res InventoryReservation, action ReservationAction) *erp.InventoryAdjustment {
	adjustment := &erp.InventoryAdjustment{
		ReferenceID:   orderID,
		ReferenceType: "order",
		SKU:           res.SKU,
		WarehouseID:   res.WarehouseID,
		Timestamp:     time.Now(),
	}

	switch action {
	case ReservationActionReserve:
		adjustment.Type = erp.AdjustmentTypeReservation
		adjustment.ReservedDelta = res.Quantity
		adjustment.Reason = fmt.Sprintf("Reserved for order %s", orderID)

	case ReservationActionRelease:
		adjustment.Type = erp.AdjustmentTypeReservation
		adjustment.ReservedDelta = -res.Quantity
		adjustment.Reason = fmt.Sprintf("Released reservation for order %s", orderID)

	case ReservationActionFulfill:
		adjustment.Type = erp.AdjustmentTypeFulfillment
		adjustment.QuantityDelta = -res.Quantity
		adjustment.ReservedDelta = -res.Quantity
		adjustment.Reason = fmt.Sprintf("Fulfilled order %s", orderID)

	case ReservationActionTransfer:
		adjustment.Type = erp.AdjustmentTypeTransfer
		adjustment.QuantityDelta = -res.Quantity
		adjustment.TransferToWarehouse = res.TransferToWarehouse
		adjustment.Reason = fmt.Sprintf("Transfer for order %s", orderID)
	}

	return adjustment
}

func generateReservationSyncID() string {
	return fmt.Sprintf("reservation_sync_%d", time.Now().UnixNano())
}
