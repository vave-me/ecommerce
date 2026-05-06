package commands

import (
	"context"
	"fmt"

	"middleman/erp/internal/domain"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

// ApproveReturn approves a return
type ApproveReturn struct {
	ReturnID   string
	ApprovedBy string
}

// ApproveReturnHandler handles approving returns
type ApproveReturnHandler struct {
	returns   es.AggregateRepository[*domain.Return]
	publisher ddd.EventPublisher[ddd.Event]
}

// NewApproveReturnHandler creates a new handler
func NewApproveReturnHandler(
	returns es.AggregateRepository[*domain.Return],
	publisher ddd.EventPublisher[ddd.Event],
) ApproveReturnHandler {
	return ApproveReturnHandler{
		returns:   returns,
		publisher: publisher,
	}
}

// ApproveReturn handles the ApproveReturn command
func (h ApproveReturnHandler) ApproveReturn(ctx context.Context, cmd ApproveReturn) error {
	// Load the return
	ret, err := h.returns.Load(ctx, cmd.ReturnID)
	if err != nil {
		return fmt.Errorf("loading return: %w", err)
	}

	// Approve the return
	event, err := ret.ApproveReturn(cmd.ApprovedBy)
	if err != nil {
		return fmt.Errorf("approving return: %w", err)
	}

	// Save the return
	if err := h.returns.Save(ctx, ret); err != nil {
		return fmt.Errorf("saving return: %w", err)
	}

	// Publish domain event
	return h.publisher.Publish(ctx, event)
}
