package models

// Alert represents a notification alert
type Alert struct {
	ID      string            `json:"id"`
	UserID  string            `json:"user_id"`
	Type    string            `json:"type"`
	Message string            `json:"message"`
	Payload map[string]string `json:"payload"`
	IsRead  bool              `json:"is_read"`
}

// ListAlertsResponse represents the response for listing alerts
type ListAlertsResponse struct {
	Alerts []Alert `json:"alerts"`
	Total  int64   `json:"total"`
}

// GetAlertsByTypeResponse represents the response for getting alerts by type
type GetAlertsByTypeResponse struct {
	Alerts []Alert `json:"alerts"`
}

// Alert type constants
const (
	AlertTypeInfo     = "info"
	AlertTypeWarning  = "warning"
	AlertTypeError    = "error"
	AlertTypeSuccess  = "success"
	AlertTypeSystem   = "system"
	AlertTypeUser     = "user"
	AlertTypeOrder    = "order"
	AlertTypePayment  = "payment"
	AlertTypeShipping = "shipping"
	AlertTypeOffer    = "offer"
)
