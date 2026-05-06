package commands

import (
	"context"
	"fmt"

	"middleman/erp/internal/domain"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

// RestockReturnItems records that return items have been restocked
type RestockReturnItems struct {
	ReturnID string
	Items    []domain.RestockedItem
}

// RestockReturnItemsHandler handles restocking return items
type RestockReturnItemsHandler struct {
	returns   es.AggregateRepository[*domain.Return]
	publisher ddd.EventPublisher[ddd.Event]
}

// NewRestockReturnItemsHandler creates a new handler
func NewRestockReturnItemsHandler(
	returns es.AggregateRepository[*domain.Return],
	publisher ddd.EventPublisher[ddd.Event],
) RestockReturnItemsHandler {
	return RestockReturnItemsHandler{
		returns:   returns,
		publisher: publisher,
	}
}

// RestockReturnItems handles the RestockReturnItems command
func (h RestockReturnItemsHandler) RestockReturnItems(ctx context.Context, cmd RestockReturnItems) error {
	// Load the return
	ret, err := h.returns.Load(ctx, cmd.ReturnID)
	if err != nil {
		return fmt.Errorf("loading return: %w", err)
	}
	
	// Record items restocked
	event, err := ret.RestockItems(cmd.Items)
	if err != nil {
		return fmt.Errorf("restocking items: %w", err)
	}
	
	// Save the return
	if err := h.returns.Save(ctx, ret); err != nil {
		return fmt.Errorf("saving return: %w", err)
	}
	
	// Publish domain event
	return h.publisher.Publish(ctx, event)
}