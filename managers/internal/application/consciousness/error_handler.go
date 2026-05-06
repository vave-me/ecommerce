package consciousness

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"time"
	
	"github.com/rs/zerolog"
)

// ErrorHandler provides comprehensive error handling for the consciousness system
type ErrorHandler struct {
	logger         zerolog.Logger
	errorChan      chan ErrorEvent
	recoveryPolicy RecoveryPolicy
	metrics        *ErrorMetrics
}

// ErrorEvent represents an error that occurred in the system
type ErrorEvent struct {
	Timestamp   time.Time
	Component   string
	Operation   string
	Error       error
	Context     map[string]interface{}
	StackTrace  string
	Severity    ErrorSeverity
	Recoverable bool
}

// ErrorSeverity levels
type ErrorSeverity int

const (
	SeverityLow ErrorSeverity = iota
	SeverityMedium
	SeverityHigh
	SeverityCritical
)

// RecoveryPolicy defines how to handle different types of errors
type RecoveryPolicy struct {
	MaxRetries          int
	RetryDelay          time.Duration
	BackoffMultiplier   float64
	MaxBackoff          time.Duration
	CircuitBreakerLimit int
	ResetTimeout        time.Duration
}

// ErrorMetrics tracks error statistics
type ErrorMetrics struct {
	TotalErrors       int64
	ErrorsByComponent map[string]int64
	ErrorsBySeverity  map[ErrorSeverity]int64
	RecoveryAttempts  int64
	RecoverySuccesses int64
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(logger zerolog.Logger) *ErrorHandler {
	return &ErrorHandler{
		logger:    logger,
		errorChan: make(chan ErrorEvent, 1000),
		recoveryPolicy: RecoveryPolicy{
			MaxRetries:          3,
			RetryDelay:          time.Second,
			BackoffMultiplier:   2.0,
			MaxBackoff:          time.Minute,
			CircuitBreakerLimit: 10,
			ResetTimeout:        5 * time.Minute,
		},
		metrics: &ErrorMetrics{
			ErrorsByComponent: make(map[string]int64),
			ErrorsBySeverity:  make(map[ErrorSeverity]int64),
		},
	}
}

// Start begins error handling goroutine
func (eh *ErrorHandler) Start(ctx context.Context) {
	go eh.processErrors(ctx)
}

// HandleError processes an error with context
func (eh *ErrorHandler) HandleError(ctx context.Context, component, operation string, err error) {
	if err == nil {
		return
	}
	
	severity := eh.classifyError(err)
	recoverable := eh.isRecoverable(err)
	
	event := ErrorEvent{
		Timestamp:   time.Now(),
		Component:   component,
		Operation:   operation,
		Error:       err,
		Context:     eh.extractContext(ctx),
		StackTrace:  string(debug.Stack()),
		Severity:    severity,
		Recoverable: recoverable,
	}
	
	select {
	case eh.errorChan <- event:
	default:
		// Channel full, log directly
		eh.logError(event)
	}
}

// HandlePanic recovers from panics and converts them to errors
func (eh *ErrorHandler) HandlePanic(component, operation string) {
	if r := recover(); r != nil {
		err := fmt.Errorf("panic recovered: %v", r)
		eh.HandleError(context.Background(), component, operation, err)
	}
}

// processErrors handles errors from the error channel
func (eh *ErrorHandler) processErrors(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-eh.errorChan:
			eh.processError(event)
		}
	}
}

// processError handles a single error event
func (eh *ErrorHandler) processError(event ErrorEvent) {
	// Update metrics
	eh.updateMetrics(event)
	
	// Log the error
	eh.logError(event)
	
	// Attempt recovery if applicable
	if event.Recoverable {
		eh.attemptRecovery(event)
	}
	
	// Alert if critical
	if event.Severity == SeverityCritical {
		eh.sendAlert(event)
	}
}

// classifyError determines error severity
func (eh *ErrorHandler) classifyError(err error) ErrorSeverity {
	// Check for specific error types
	switch {
	case errors.Is(err, context.Canceled):
		return SeverityLow
	case errors.Is(err, context.DeadlineExceeded):
		return SeverityMedium
	case isNetworkError(err):
		return SeverityMedium
	case isDatabaseError(err):
		return SeverityHigh
	case isPanicError(err):
		return SeverityCritical
	default:
		return SeverityMedium
	}
}

// isRecoverable determines if an error can be recovered from
func (eh *ErrorHandler) isRecoverable(err error) bool {
	// Context errors are not recoverable
	if errors.Is(err, context.Canceled) {
		return false
	}
	
	// Network and timeout errors are often recoverable
	if isNetworkError(err) || errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	
	// Database connection errors might be recoverable
	if isDatabaseError(err) {
		return true
	}
	
	return false
}

// attemptRecovery tries to recover from an error
func (eh *ErrorHandler) attemptRecovery(event ErrorEvent) {
	eh.metrics.RecoveryAttempts++
	
	delay := eh.recoveryPolicy.RetryDelay
	
	for attempt := 1; attempt <= eh.recoveryPolicy.MaxRetries; attempt++ {
		eh.logger.Info().
			Str("component", event.Component).
			Str("operation", event.Operation).
			Int("attempt", attempt).
			Msg("Attempting recovery")
		
		// Wait before retry
		time.Sleep(delay)
		
		// Component-specific recovery logic would go here
		// For now, we just log the attempt
		
		// Exponential backoff
		delay = time.Duration(float64(delay) * eh.recoveryPolicy.BackoffMultiplier)
		if delay > eh.recoveryPolicy.MaxBackoff {
			delay = eh.recoveryPolicy.MaxBackoff
		}
	}
	
	eh.logger.Error().
		Str("component", event.Component).
		Str("operation", event.Operation).
		Msg("Recovery failed after maximum attempts")
}

// logError logs an error event with appropriate level
func (eh *ErrorHandler) logError(event ErrorEvent) {
	logger := eh.logger.With().
		Str("component", event.Component).
		Str("operation", event.Operation).
		Time("timestamp", event.Timestamp).
		Interface("context", event.Context).
		Logger()
	
	switch event.Severity {
	case SeverityLow:
		logger.Warn().Err(event.Error).Msg("Low severity error")
	case SeverityMedium:
		logger.Error().Err(event.Error).Msg("Medium severity error")
	case SeverityHigh:
		logger.Error().
			Err(event.Error).
			Str("stack_trace", event.StackTrace).
			Msg("High severity error")
	case SeverityCritical:
		logger.Fatal().
			Err(event.Error).
			Str("stack_trace", event.StackTrace).
			Msg("Critical error")
	}
}

// updateMetrics updates error metrics
func (eh *ErrorHandler) updateMetrics(event ErrorEvent) {
	eh.metrics.TotalErrors++
	eh.metrics.ErrorsByComponent[event.Component]++
	eh.metrics.ErrorsBySeverity[event.Severity]++
}

// sendAlert sends an alert for critical errors
func (eh *ErrorHandler) sendAlert(event ErrorEvent) {
	// This would integrate with alerting systems like PagerDuty, Slack, etc.
	eh.logger.Error().
		Str("component", event.Component).
		Str("operation", event.Operation).
		Err(event.Error).
		Msg("CRITICAL ALERT: System error requires immediate attention")
}

// extractContext extracts relevant context from context.Context
func (eh *ErrorHandler) extractContext(ctx context.Context) map[string]interface{} {
	contextMap := make(map[string]interface{})
	
	// Extract common context values
	if userID := ctx.Value("user_id"); userID != nil {
		contextMap["user_id"] = userID
	}
	if requestID := ctx.Value("request_id"); requestID != nil {
		contextMap["request_id"] = requestID
	}
	if traceID := ctx.Value("trace_id"); traceID != nil {
		contextMap["trace_id"] = traceID
	}
	
	return contextMap
}

// GetMetrics returns current error metrics
func (eh *ErrorHandler) GetMetrics() ErrorMetrics {
	return *eh.metrics
}

// Helper functions to classify error types
func isNetworkError(err error) bool {
	errStr := err.Error()
	return contains(errStr, "connection refused", "network", "timeout", "no such host")
}

func isDatabaseError(err error) bool {
	errStr := err.Error()
	return contains(errStr, "database", "sql", "postgres", "connection pool")
}

func isPanicError(err error) bool {
	return contains(err.Error(), "panic recovered")
}

func contains(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if contains := func(s, substr string) bool {
			return len(substr) > 0 && (s == substr || (len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || findSubstring(s, substr))))
		}(s, substr); contains {
			return true
		}
	}
	return false
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}