package tools

import (
	"context"
	"fmt"
	"log"
	"time"

	"middleman/managers/internal/domain"
)

// SupportToolService handles customer support and ticket management operations
type SupportToolService struct {
	support domain.SupportRepository
}

// NewSupportToolService creates a new support tool service
func NewSupportToolService(supportRepo domain.SupportRepository) *SupportToolService {
	return &SupportToolService{
		support: supportRepo,
	}
}

// ExecuteOperation executes a support operation with streaming progress
func (s *SupportToolService) ExecuteOperation(
	ctx context.Context,
	operation string,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	// Send initial progress
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "support_operation",
		Status:   "progress",
		Progress: 25.0,
		Metadata: map[string]interface{}{
			"step":      "initializing_support_operation",
			"operation": operation,
		},
		Timestamp: time.Now().Unix(),
	}

	switch operation {
	case "start_support":
		return s.startSupport(ctx, parameters, streamChan, toolID)
	case "create_ticket":
		return s.createTicket(ctx, parameters, streamChan, toolID)
	case "list_tickets":
		return s.listTickets(ctx, parameters, streamChan, toolID)
	case "get_ticket":
		return s.getTicket(ctx, parameters, streamChan, toolID)
	case "update_ticket":
		return s.updateTicket(ctx, parameters, streamChan, toolID)
	case "delete_ticket", "close_ticket":
		return s.deleteTicket(ctx, parameters, streamChan, toolID)
	default:
		return s.handleUnsupportedOperation(operation, streamChan, toolID)
	}
}

func (s *SupportToolService) startSupport(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	userID := getStringParam(parameters, "user_id", "")
	if userID == "" {
		return nil, fmt.Errorf("user_id parameter required for start_support operation")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "support_operation",
		Status:   "progress",
		Progress: 50.0,
		Metadata: map[string]interface{}{
			"step":    "starting_support",
			"user_id": userID,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("SupportToolService: Starting support for user: %s", userID)
	response, err := s.support.StartSupport(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("start support failed: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "support_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"support_session": response,
			"user_id":         userID,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "support",
		"operation":   "start_support",
		"result":      response,
		"user_id":     userID,
	}, nil
}

func (s *SupportToolService) createTicket(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	supportID := getStringParam(parameters, "support_id", "")
	title := getStringParam(parameters, "title", "")
	description := getStringParam(parameters, "description", "")

	if supportID == "" || title == "" || description == "" {
		return nil, fmt.Errorf("support_id, title, and description are required for create_ticket operation")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "support_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":       "creating_ticket",
			"support_id": supportID,
			"title":      title,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("SupportToolService: Creating ticket for support: %s", supportID)
	response, err := s.support.CreateTicket(ctx, supportID, title, description)
	if err != nil {
		return nil, fmt.Errorf("create ticket failed: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "support_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"ticket":     response,
			"support_id": supportID,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "ticket",
		"operation":   "create_ticket",
		"result":      response,
		"support_id":  supportID,
	}, nil
}

func (s *SupportToolService) listTickets(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	supportID := getStringParam(parameters, "support_id", "")
	if supportID == "" {
		return nil, fmt.Errorf("support_id parameter required for list_tickets operation")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "support_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":       "listing_tickets",
			"support_id": supportID,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("SupportToolService: Listing tickets for support: %s", supportID)
	response, err := s.support.GetTickets(ctx, supportID)
	if err != nil {
		return nil, fmt.Errorf("list tickets failed: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "support_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"tickets":    response,
			"support_id": supportID,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "ticket",
		"operation":   "list_tickets",
		"result":      response,
		"support_id":  supportID,
	}, nil
}

func (s *SupportToolService) getTicket(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	supportID := getStringParam(parameters, "support_id", "")
	ticketID := getStringParam(parameters, "id", "")
	if ticketID == "" {
		ticketID = getStringParam(parameters, "ticket_id", "")
	}

	if supportID == "" || ticketID == "" {
		return nil, fmt.Errorf("support_id and ticket_id are required for get_ticket operation")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "support_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":       "getting_ticket",
			"support_id": supportID,
			"ticket_id":  ticketID,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("SupportToolService: Getting ticket: %s for support: %s", ticketID, supportID)
	response, err := s.support.GetTicket(ctx, ticketID)
	if err != nil {
		return nil, fmt.Errorf("get ticket failed: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "support_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"ticket":     response,
			"support_id": supportID,
			"ticket_id":  ticketID,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "ticket",
		"operation":   "get_ticket",
		"result":      response,
		"support_id":  supportID,
		"ticket_id":   ticketID,
	}, nil
}

func (s *SupportToolService) updateTicket(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	supportID := getStringParam(parameters, "support_id", "")
	ticketID := getStringParam(parameters, "id", "")
	if ticketID == "" {
		ticketID = getStringParam(parameters, "ticket_id", "")
	}
	status := getStringParam(parameters, "status", "")

	if supportID == "" || ticketID == "" {
		return nil, fmt.Errorf("support_id and ticket_id are required for update_ticket operation")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "support_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":       "updating_ticket",
			"support_id": supportID,
			"ticket_id":  ticketID,
			"status":     status,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("SupportToolService: Updating ticket: %s for support: %s", ticketID, supportID)
	response, err := s.support.UpdateTicket(ctx, ticketID, status, "")
	if err != nil {
		return nil, fmt.Errorf("update ticket failed: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "support_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"ticket":     response,
			"support_id": supportID,
			"ticket_id":  ticketID,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "ticket",
		"operation":   "update_ticket",
		"result":      response,
		"support_id":  supportID,
		"ticket_id":   ticketID,
	}, nil
}

func (s *SupportToolService) deleteTicket(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	supportID := getStringParam(parameters, "support_id", "")
	ticketID := getStringParam(parameters, "id", "")
	if ticketID == "" {
		ticketID = getStringParam(parameters, "ticket_id", "")
	}

	if supportID == "" || ticketID == "" {
		return nil, fmt.Errorf("support_id and ticket_id are required for delete_ticket operation")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "support_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":       "closing_ticket",
			"support_id": supportID,
			"ticket_id":  ticketID,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("SupportToolService: Closing ticket: %s for support: %s", ticketID, supportID)
	response, err := s.support.CloseTicket(ctx, ticketID, "closed_via_api")
	if err != nil {
		return nil, fmt.Errorf("delete ticket failed: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "support_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"ticket":     response,
			"support_id": supportID,
			"ticket_id":  ticketID,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "ticket",
		"operation":   "delete_ticket",
		"result":      response,
		"support_id":  supportID,
		"ticket_id":   ticketID,
	}, nil
}

func (s *SupportToolService) handleUnsupportedOperation(
	operation string,
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	streamChan <- ToolExecutionStream{
		ID:        toolID,
		ToolName:  "support_operation",
		Status:    "error",
		Progress:  100.0,
		Error:     fmt.Sprintf("Support operation '%s' not implemented", operation),
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "support",
		"operation":   operation,
		"message":     fmt.Sprintf("Support operation '%s' not implemented yet", operation),
	}, nil
}
