package tools

import (
	"context"
	"fmt"
	"log"
	"time"

	"middleman/managers/internal/domain"
)

// ShippingToolService handles shipping and logistics operations
type ShippingToolService struct {
	shipping domain.ShippingRepository
}

// NewShippingToolService creates a new shipping tool service
func NewShippingToolService(shippingRepo domain.ShippingRepository) *ShippingToolService {
	return &ShippingToolService{
		shipping: shippingRepo,
	}
}

// ExecuteOperation executes a shipping operation with streaming progress
func (s *ShippingToolService) ExecuteOperation(
	ctx context.Context,
	operation string,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	// Send initial progress
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "shipping_operation",
		Status:   "progress",
		Progress: 25.0,
		Metadata: map[string]interface{}{
			"step":      "initializing_shipping_operation",
			"operation": operation,
		},
		Timestamp: time.Now().Unix(),
	}

	switch operation {
	case "create_shipping":
		return s.createShipping(ctx, parameters, streamChan, toolID)
	case "track_shipping":
		return s.trackShipping(ctx, parameters, streamChan, toolID)
	default:
		return s.handleUnsupportedOperation(operation, streamChan, toolID)
	}
}

func (s *ShippingToolService) createShipping(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	productID := getStringParam(parameters, "product_id", "")
	trackingNumber := getStringParam(parameters, "tracking_number", "")
	labelURL := getStringParam(parameters, "label_url", "")
	senderName := getStringParam(parameters, "sender_name", "")
	senderAddress := getStringParam(parameters, "sender_address", "")
	receiverName := getStringParam(parameters, "receiver_name", "")
	receiverAddress := getStringParam(parameters, "receiver_address", "")
	weight := getStringParam(parameters, "weight", "")
	dimensions := getStringParam(parameters, "dimensions", "")
	serviceTypes := getStringParam(parameters, "service_types", "")

	if productID == "" || trackingNumber == "" {
		return nil, fmt.Errorf("product_id and tracking_number are required for create_shipping operation")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "shipping_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":            "creating_shipping",
			"product_id":      productID,
			"tracking_number": trackingNumber,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("ShippingToolService: Creating shipping for product: %s", productID)
	response, err := s.shipping.CreateShipping(ctx, productID, trackingNumber, labelURL, senderName, senderAddress, receiverName, receiverAddress, weight, dimensions, serviceTypes)
	if err != nil {
		return nil, fmt.Errorf("create shipping failed: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "shipping_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"shipping":        response,
			"product_id":      productID,
			"tracking_number": trackingNumber,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "shipping",
		"operation":   "create_shipping",
		"result":      response,
		"product_id":  productID,
	}, nil
}

func (s *ShippingToolService) trackShipping(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	trackingNumber := getStringParam(parameters, "tracking_number", "")
	if trackingNumber == "" {
		trackingNumber = getStringParam(parameters, "id", "")
	}

	if trackingNumber == "" {
		return nil, fmt.Errorf("tracking_number or id parameter required for track_shipping operation")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "shipping_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":            "tracking_shipping",
			"tracking_number": trackingNumber,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("ShippingToolService: Tracking shipping with number: %s", trackingNumber)
	response, err := s.shipping.TrackShipping(ctx, trackingNumber)
	if err != nil {
		return nil, fmt.Errorf("track shipping failed: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "shipping_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"tracking_info":   response,
			"tracking_number": trackingNumber,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type":     "shipping",
		"operation":       "track_shipping",
		"result":          response,
		"tracking_number": trackingNumber,
	}, nil
}

func (s *ShippingToolService) handleUnsupportedOperation(
	operation string,
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	streamChan <- ToolExecutionStream{
		ID:        toolID,
		ToolName:  "shipping_operation",
		Status:    "error",
		Progress:  100.0,
		Error:     fmt.Sprintf("Shipping operation '%s' not implemented", operation),
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "shipping",
		"operation":   operation,
		"message":     fmt.Sprintf("Shipping operation '%s' not implemented yet", operation),
	}, nil
}
