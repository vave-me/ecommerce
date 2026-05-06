package domain

import "errors"

// Stream errors
var (
	ErrStreamNotFound      = errors.New("stream not found")
	ErrStreamAlreadyActive = errors.New("stream already active")
	ErrStreamNotActive     = errors.New("stream not active")
	ErrMaxViewersReached   = errors.New("maximum viewers reached")
	ErrInvalidStreamState  = errors.New("invalid stream state")
	ErrStreamExpired       = errors.New("stream has expired")
)

// Webhook errors
var (
	ErrWebhookNotFound       = errors.New("webhook subscription not found")
	ErrWebhookInvalidURL     = errors.New("invalid webhook URL")
	ErrWebhookDeliveryFailed = errors.New("webhook delivery failed")
	ErrWebhookMaxRetries     = errors.New("webhook max retries exceeded")
)

// Auth errors
var (
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrInvalidToken      = errors.New("invalid token")
	ErrTokenExpired      = errors.New("token expired")
	ErrInsufficientScope = errors.New("insufficient scope")
)

// Validation errors
var (
	ErrInvalidInput      = errors.New("invalid input")
	ErrMissingField      = errors.New("required field missing")
	ErrInvalidFormat     = errors.New("invalid format")
	ErrValueOutOfRange   = errors.New("value out of range")
	ErrDuplicateResource = errors.New("duplicate resource")
)

// Infrastructure errors
var (
	ErrDatabaseConnection = errors.New("database connection failed")
	ErrCacheUnavailable   = errors.New("cache unavailable")
	ErrExternalService    = errors.New("external service error")
	ErrTimeout            = errors.New("operation timeout")
)