package commands

import (
	"context"
	"fmt"
	"middleman/internal/erp"
	"time"

	"middleman/erp/internal/domain"
)

// SendOrder command sends an order to the ERP system
type SendOrder struct {
	ConnectorID string
	Order       *erp.OrderPayload
}

// SendOrderHandler handles the SendOrder command
type SendOrderHandler struct {
	registry   erp.ConnectorRegistry
	repository domain.OrderSyncRepository
}

// NewSendOrderHandler creates a new handler
func NewSendOrderHandler(
	registry erp.ConnectorRegistry,
	repository domain.OrderSyncRepository,
) SendOrderHandler {
	return SendOrderHandler{
		registry:   registry,
		repository: repository,
	}
}

// SendOrder sends an order to the ERP system
func (h SendOrderHandler) SendOrder(ctx context.Context, cmd SendOrder) error {
	// Validate order
	if cmd.Order == nil {
		return fmt.Errorf("order is required")
	}
	if cmd.Order.OrderID == "" {
		return fmt.Errorf("order ID is required")
	}
	if cmd.Order.CustomerID == "" {
		return fmt.Errorf("customer ID is required")
	}
	if len(cmd.Order.Items) == 0 {
		return fmt.Errorf("order must have at least one item")
	}

	// Get connector
	connector, err := h.registry.GetConnector(cmd.ConnectorID)
	if err != nil {
		return fmt.Errorf("getting connector: %w", err)
	}

	// Create order sync record
	orderSync := &domain.OrderSync{
		ID:          generateOrderSyncID(),
		ConnectorID: cmd.ConnectorID,
		OrderID:     cmd.Order.OrderID,
		Direction:   domain.DirectionOutbound,
		Status:      domain.OrderSyncStatusPending,
		AttemptedAt: time.Now(),
		Payload:     cmd.Order,
	}

	if err := h.repository.Create(ctx, orderSync); err != nil {
		return fmt.Errorf("creating order sync record: %w", err)
	}

	// Send order to ERP
	if err := connector.SendOrder(ctx, cmd.Order); err != nil {
		// Update sync record with failure
		orderSync.Status = domain.OrderSyncStatusFailed
		orderSync.Error = err.Error()
		orderSync.CompletedAt = ptrTime(time.Now())
		h.repository.Update(ctx, orderSync)
		return fmt.Errorf("sending order to ERP: %w", err)
	}

	// Update sync record with success
	orderSync.Status = domain.OrderSyncStatusCompleted
	//TODO grpc order repostiory call
	orderSync.CompletedAt = ptrTime(time.Now())
	if err := h.repository.Update(ctx, orderSync); err != nil {
		return fmt.Errorf("updating order sync record: %w", err)
	}

	return nil
}
