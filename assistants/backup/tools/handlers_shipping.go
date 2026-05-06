package tools

import (
	"context"
	"fmt"
	"strings"
)

// ==================== SHIPPING HANDLERS ====================
func (r *ComprehensiveToolRegistry) initializeShippingHandlers() {
	r.handlers["shipping_create"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		productID := getStringParam(params, "product_id")
		trackingNumber := getStringParam(params, "tracking_number")
		labelURL := getStringParam(params, "label_url")
		senderName := getStringParam(params, "sender_name")
		senderAddress := getStringParam(params, "sender_address")
		receiverName := getStringParam(params, "receiver_name")
		receiverAddress := getStringParam(params, "receiver_address")
		dimensions := getStringParam(params, "dimensions")

		// Validate required parameters
		if err := ValidateIDParam("product_id", productID); err != nil {
			return nil, fmt.Errorf("invalid product_id: %w", err)
		}
		if trackingNumber == "" {
			return nil, fmt.Errorf("tracking_number is required")
		}
		if senderName == "" {
			return nil, fmt.Errorf("sender_name is required")
		}
		if senderAddress == "" {
			return nil, fmt.Errorf("sender_address is required")
		}
		if receiverName == "" {
			return nil, fmt.Errorf("receiver_name is required")
		}
		if receiverAddress == "" {
			return nil, fmt.Errorf("receiver_address is required")
		}

		// Validate weight
		weight := getFloat64Param(params, "weight", 0)
		if weight <= 0 {
			return nil, fmt.Errorf("weight must be greater than zero")
		}
		weightStr := fmt.Sprintf("%.2f", weight)
		
		// Convert service types array to comma-separated string
		serviceTypesArray := getStringArrayParam(params, "service_types")
		if len(serviceTypesArray) == 0 {
			return nil, fmt.Errorf("at least one service type is required")
		}
		serviceTypesStr := strings.Join(serviceTypesArray, ",")

		// Sanitize string inputs
		senderName = SanitizeString(senderName)
		senderAddress = SanitizeString(senderAddress)
		receiverName = SanitizeString(receiverName)
		receiverAddress = SanitizeString(receiverAddress)
		dimensions = SanitizeString(dimensions)
		
		return reg.shippingRepo.CreateNewShipmentWithDetails(ctx,
			productID,
			trackingNumber,
			labelURL,
			senderName,
			senderAddress,
			receiverName,
			receiverAddress,
			weightStr,
			dimensions,
			serviceTypesStr)
	}

	r.handlers["shipping_track"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		trackingNumber := getStringParam(params, "tracking_number")
		if trackingNumber == "" {
			return nil, fmt.Errorf("tracking_number is required")
		}
		return reg.shippingRepo.TrackShipmentByTrackingNumber(ctx, trackingNumber)
	}
}