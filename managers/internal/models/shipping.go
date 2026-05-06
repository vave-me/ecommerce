package models

import "time"

// Shipping represents a shipping record for a product
type Shipping struct {
	ID              string    `json:"id"`
	ProductID       string    `json:"product_id"`
	TrackingNumber  string    `json:"tracking_number"`
	LabelURL        string    `json:"label_url"`
	SenderName      string    `json:"sender_name"`
	SenderAddress   string    `json:"sender_address"`
	ReceiverName    string    `json:"receiver_name"`
	ReceiverAddress string    `json:"receiver_address"`
	Weight          string    `json:"weight"`
	Dimensions      string    `json:"dimensions"`
	ServiceTypes    string    `json:"service_types"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
}

// CreateShippingResponse represents the response after creating a shipping record
type CreateShippingResponse struct {
	ID string `json:"id"`
}

// TrackShippingResponse represents the response for tracking a shipment
type TrackShippingResponse struct {
	Shipping *Shipping `json:"shipping"`
}
