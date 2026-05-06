package commands

import (
	"context"
	"fmt"

	"middleman/erp/internal/domain"
	"middleman/internal/ddd"
	"middleman/internal/erp"
	"middleman/internal/es"
)

// SyncReturnToERP syncs a return to the ERP system
type SyncReturnToERP struct {
	ReturnID    string
	ConnectorID string
}

// SyncReturnToERPHandler handles syncing returns to ERP
type SyncReturnToERPHandler struct {
	returns   es.AggregateRepository[*domain.Return]
	registry  erp.ConnectorRegistry
	publisher ddd.EventPublisher[ddd.Event]
}

// NewSyncReturnToERPHandler creates a new handler
func NewSyncReturnToERPHandler(
	returns es.AggregateRepository[*domain.Return],
	registry erp.ConnectorRegistry,
	publisher ddd.EventPublisher[ddd.Event],
) SyncReturnToERPHandler {
	return SyncReturnToERPHandler{
		returns:   returns,
		registry:  registry,
		publisher: publisher,
	}
}

// SyncReturnToERP syncs a return to ERP
func (h SyncReturnToERPHandler) SyncReturnToERP(ctx context.Context, cmd SyncReturnToERP) error {
	// Load the return
	ret, err := h.returns.Load(ctx, cmd.ReturnID)
	if err != nil {
		return fmt.Errorf("loading return: %w", err)
	}
	
	// Get the connector
	connector, err := h.registry.GetConnector(cmd.ConnectorID)
	if err != nil {
		return fmt.Errorf("getting connector: %w", err)
	}
	
	// Convert to ERP return format
	erpReturn := &erp.ReturnPayload{
		ReturnID:        ret.ID(),
		OriginalOrderID: ret.OriginalOrderID,
		CustomerID:      ret.CustomerID,
		Reason:          string(ret.Reason),
		Status:          string(ret.Status),
		CreatedAt:       ret.CreatedAt,
		Items:           convertReturnItems(ret.Items),
		TotalRefund:     ret.RefundDetails.Amount,
		RefundMethod:    string(ret.RefundDetails.Method),
		WarehouseID:     ret.WarehouseID,
		Notes:           ret.Notes,
	}
	
	// Process return in ERP
	processedReturn, err := connector.ProcessReturn(ctx, erpReturn)
	if err != nil {
		// Record sync failure
		event, recordErr := ret.RecordSyncFailure(err)
		if recordErr != nil {
			return fmt.Errorf("recording sync failure: %w", recordErr)
		}
		if err := h.returns.Save(ctx, ret); err != nil {
			return fmt.Errorf("saving return after sync failure: %w", err)
		}
		if err := h.publisher.Publish(ctx, event); err != nil {
			return fmt.Errorf("publishing sync failure event: %w", err)
		}
		return fmt.Errorf("processing return in ERP: %w", err)
	}
	
	// Link to ERP
	if processedReturn != nil && processedReturn.ExternalID != "" {
		event, err := ret.LinkToERP(processedReturn.ExternalID)
		if err != nil {
			return fmt.Errorf("linking to ERP: %w", err)
		}
		// Save the return with the link event
		if err := h.returns.Save(ctx, ret); err != nil {
			return fmt.Errorf("saving return after linking: %w", err)
		}
		// Publish the link event
		if err := h.publisher.Publish(ctx, event); err != nil {
			return fmt.Errorf("publishing link event: %w", err)
		}
		return nil
	}
	
	// Update inventory for returned items
	for _, item := range ret.Items {
		// Only restock items in acceptable condition
		if ret.ShouldRestock(item.Condition) {
			adjustment := &erp.InventoryAdjustment{
				ReferenceID:   ret.ID(),
				ReferenceType: "return",
				SKU:           item.SKU,
				WarehouseID:   ret.WarehouseID,
				Type:          erp.AdjustmentTypeReturn,
				QuantityDelta: item.Quantity,
				Reason:        fmt.Sprintf("Return %s - %s", ret.ReturnNumber, item.Condition),
				Timestamp:     ret.CreatedAt,
			}
			
			if err := connector.UpdateInventory(ctx, []*erp.InventoryAdjustment{adjustment}); err != nil {
				// Log but don't fail the return sync
				// Could add this to events or logs
			}
		}
	}
	
	// Save the return
	if err := h.returns.Save(ctx, ret); err != nil {
		return fmt.Errorf("saving return: %w", err)
	}
	
	return nil
}

func convertReturnItems(items []domain.ReturnItem) []erp.ReturnItem {
	result := make([]erp.ReturnItem, len(items))
	for i, item := range items {
		result[i] = erp.ReturnItem{
			SKU:           item.SKU,
			Quantity:      item.Quantity,
			Reason:        string(item.Condition),
			RefundAmount:  item.RefundAmount,
			RestockingFee: item.RestockingFee,
			Notes:         item.InspectionNotes,
		}
	}
	return result
}