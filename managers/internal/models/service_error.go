package models

// ServiceError provides a consistent, JSON-friendly error payload consumed by LLM clients.
// Code follows a kebab-case convention (e.g., "validation-error", "not-found", "internal-error").
// Message is human-readable; Details can carry any additional machine-readable context.
type ServiceError struct {
	Code    string      `json:"code"`              // Stable, short error identifier
	Message string      `json:"message"`           // Human-readable explanation
	Details interface{} `json:"details,omitempty"` // Optional structured data
}

// Error implements the error interface, returning the human-readable message
func (se *ServiceError) Error() string {
	return se.Message
}
