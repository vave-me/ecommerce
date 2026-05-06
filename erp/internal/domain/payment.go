package domain

import "time"

// Payment represents a payment entity
type Payment struct {
	ID              string    `json:"id"`
	UserCustomerID  string    `json:"user_customer_id"`
	Amount          int64     `json:"amount"`
	PaymentMethodID string    `json:"payment_method_id"`
	Status          string    `json:"status"`
	ClientSecret    string    `json:"client_secret,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// AuthorizePaymentResponse represents the response for payment authorization
type AuthorizePaymentResponse struct {
	ID           string `json:"id"`
	ClientSecret string `json:"client_secret"`
}

// ConfirmPaymentResponse represents the response for payment confirmation
type ConfirmPaymentResponse struct {
	ID            string `json:"id"`
	PaymentStatus string `json:"payment_status"`
}

// CreateInvoiceResponse represents the response for invoice creation
type CreateInvoiceResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

// AdjustInvoiceResponse represents the response for invoice adjustment
type AdjustInvoiceResponse struct {
	ID         string    `json:"id"`
	AdjustedAt time.Time `json:"adjusted_at"`
	NewAmount  int64     `json:"new_amount"`
}

// PayInvoiceResponse represents the response for invoice payment
type PayInvoiceResponse struct {
	ID            string    `json:"id"`
	PaidAt        time.Time `json:"paid_at"`
	PaymentStatus string    `json:"payment_status"`
}

// CapturePaymentResponse represents the response for payment capture
type CapturePaymentResponse struct {
	PaymentID      string `json:"payment_id"`
	PaymentStatus  string `json:"payment_status"`
	CapturedAmount int64  `json:"captured_amount"`
}

// HandleWebhookResponse represents the response for webhook handling
type HandleWebhookResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Payment status constants
const (
	PaymentStatusPending    = "pending"
	PaymentStatusAuthorized = "authorized"
	PaymentStatusConfirmed  = "confirmed"
	PaymentStatusCaptured   = "captured"
	PaymentStatusFailed     = "failed"
	PaymentStatusCanceled   = "canceled"
)
