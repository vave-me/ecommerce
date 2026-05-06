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

// CreateReturn creates a new return
type CreateReturn struct {
	ReturnID        string
	ReturnNumber    string
	OriginalOrderID string
	CustomerID      string
	CustomerName    string
	CustomerEmail   string
	Reason          domain.ReturnReason
	Items           []domain.ReturnItem
	RefundMethod    domain.RefundMethod
	RefundAmount    float64
	WarehouseID     string
	Notes           string
	ConnectorID     string
}

// CreateReturnHandler handles creating returns
type CreateReturnHandler struct {
	returns   es.AggregateRepository[*domain.Return]
	publisher ddd.EventPublisher[ddd.Event]
}

// NewCreateReturnHandler creates a new handler
func NewCreateReturnHandler(
	returns es.AggregateRepository[*domain.Return],
	publisher ddd.EventPublisher[ddd.Event],
) CreateReturnHandler {
	return CreateReturnHandler{
		returns:   returns,
		publisher: publisher,
	}
}

// CreateReturn handles the CreateReturn command
func (h CreateReturnHandler) CreateReturn(ctx context.Context, cmd CreateReturn) error {
	// Generate return ID if not provided
	if cmd.ReturnID == "" {
		cmd.ReturnID = uuid.New().String()
	}
	
	// Generate return number if not provided
	if cmd.ReturnNumber == "" {
		cmd.ReturnNumber = generateReturnNumber()
	}
	
	// Load the return aggregate (will create new if doesn't exist)
	ret, err := h.returns.Load(ctx, cmd.ReturnID)
	if err != nil {
		return fmt.Errorf("loading return: %w", err)
	}
	
	// Create the return
	event, err := ret.CreateReturn(
		cmd.ReturnNumber,
		cmd.OriginalOrderID,
		cmd.CustomerID,
		cmd.CustomerName,
		cmd.CustomerEmail,
		cmd.Reason,
		cmd.Items,
		cmd.RefundMethod,
		cmd.RefundAmount,
		cmd.WarehouseID,
		cmd.Notes,
		cmd.ConnectorID,
	)
	if err != nil {
		return fmt.Errorf("creating return: %w", err)
	}
	
	// Save the return
	if err := h.returns.Save(ctx, ret); err != nil {
		return fmt.Errorf("saving return: %w", err)
	}
	
	// Publish domain event
	return h.publisher.Publish(ctx, event)
}

func generateReturnNumber() string {
	return fmt.Sprintf("RET-%d", time.Now().Unix())
}