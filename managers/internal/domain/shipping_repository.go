package domain

import (
	"context"
	"middleman/managers/internal/models"
)

type ShippingRepository interface {
	// Shipping management
	CreateShipping(ctx context.Context, productID, trackingNumber, labelURL, senderName, senderAddress, receiverName, receiverAddress, weight, dimensions, serviceTypes string) (*models.CreateShippingResponse, error)
	TrackShipping(ctx context.Context, trackingNumber string) (*models.TrackShippingResponse, error)
}
