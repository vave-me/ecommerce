package domain

import (
	"context"
	"middleman/assistants/internal/models"
)

type ShippingRepository interface {
	// Shipping management
	CreateNewShipmentWithDetails(ctx context.Context, productID, trackingNumber, labelURL, senderName, senderAddress, receiverName, receiverAddress, weight, dimensions, serviceTypes string) (*models.CreateShippingResponse, error)
	TrackShipmentByTrackingNumber(ctx context.Context, trackingNumber string) (*models.TrackShippingResponse, error)
}
