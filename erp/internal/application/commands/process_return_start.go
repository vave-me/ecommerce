package commands

import (
	"context"
	"fmt"

	"middleman/erp/internal/domain"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

// ProcessReturnStart starts processing a return
type ProcessReturnStart struct {
	ReturnID    string
	ERPReturnID string // Optional ERP return ID if already created
}

// ProcessReturnStartHandler handles starting return processing
type ProcessReturnStartHandler struct {
	returns   es.AggregateRepository[*domain.Return]
	publisher ddd.EventPublisher[ddd.Event]
}

// NewProcessReturnStartHandler creates a new handler
func NewProcessReturnStartHandler(
	returns es.AggregateRepository[*domain.Return],
	publisher ddd.EventPublisher[ddd.Event],
) ProcessReturnStartHandler {
	return ProcessReturnStartHandler{
		returns:   returns,
		publisher: publisher,
	}
}

// ProcessReturnStart handles the ProcessReturnStart command
func (h ProcessReturnStartHandler) ProcessReturnStart(ctx context.Context, cmd ProcessReturnStart) error {
	// Load the return
	ret, err := h.returns.Load(ctx, cmd.ReturnID)
	if err != nil {
		return fmt.Errorf("loading return: %w", err)
	}
	
	// Start processing the return
	event, err := ret.ProcessReturn(cmd.ERPReturnID)
	if err != nil {
		return fmt.Errorf("processing return: %w", err)
	}
	
	// Save the return
	if err := h.returns.Save(ctx, ret); err != nil {
		return fmt.Errorf("saving return: %w", err)
	}
	
	// Publish domain event
	return h.publisher.Publish(ctx, event)
}