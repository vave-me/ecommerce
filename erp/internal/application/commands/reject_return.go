package commands

import (
	"context"
	"fmt"

	"middleman/erp/internal/domain"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

// RejectReturn rejects a return
type RejectReturn struct {
	ReturnID string
	Reason   string
}

// RejectReturnHandler handles rejecting returns
type RejectReturnHandler struct {
	returns   es.AggregateRepository[*domain.Return]
	publisher ddd.EventPublisher[ddd.Event]
}

// NewRejectReturnHandler creates a new handler
func NewRejectReturnHandler(
	returns es.AggregateRepository[*domain.Return],
	publisher ddd.EventPublisher[ddd.Event],
) RejectReturnHandler {
	return RejectReturnHandler{
		returns:   returns,
		publisher: publisher,
	}
}

// RejectReturn handles the RejectReturn command
func (h RejectReturnHandler) RejectReturn(ctx context.Context, cmd RejectReturn) error {
	// Load the return
	ret, err := h.returns.Load(ctx, cmd.ReturnID)
	if err != nil {
		return fmt.Errorf("loading return: %w", err)
	}
	
	// Reject the return
	event, err := ret.RejectReturn(cmd.Reason)
	if err != nil {
		return fmt.Errorf("rejecting return: %w", err)
	}
	
	// Save the return
	if err := h.returns.Save(ctx, ret); err != nil {
		return fmt.Errorf("saving return: %w", err)
	}
	
	// Publish domain event
	return h.publisher.Publish(ctx, event)
}