package commands

import (
	"context"
	"fmt"

	"middleman/erp/internal/domain"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

// CompleteReturn completes a return
type CompleteReturn struct {
	ReturnID            string
	RefundTransactionID string
}

// CompleteReturnHandler handles completing returns
type CompleteReturnHandler struct {
	returns   es.AggregateRepository[*domain.Return]
	publisher ddd.EventPublisher[ddd.Event]
}

// NewCompleteReturnHandler creates a new handler
func NewCompleteReturnHandler(
	returns es.AggregateRepository[*domain.Return],
	publisher ddd.EventPublisher[ddd.Event],
) CompleteReturnHandler {
	return CompleteReturnHandler{
		returns:   returns,
		publisher: publisher,
	}
}

// CompleteReturn handles the CompleteReturn command
func (h CompleteReturnHandler) CompleteReturn(ctx context.Context, cmd CompleteReturn) error {
	// Load the return
	ret, err := h.returns.Load(ctx, cmd.ReturnID)
	if err != nil {
		return fmt.Errorf("loading return: %w", err)
	}

	// Complete the return
	event, err := ret.CompleteReturn(cmd.RefundTransactionID)
	if err != nil {
		return fmt.Errorf("completing return: %w", err)
	}

	// Save the return
	if err := h.returns.Save(ctx, ret); err != nil {
		return fmt.Errorf("saving return: %w", err)
	}

	// Publish domain event
	return h.publisher.Publish(ctx, event)
}
