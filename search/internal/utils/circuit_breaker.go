package utils

import (
	"context"
	"sync"
	"time"

	"github.com/stackus/errors"
)

// CircuitBreaker states
const (
	StateClosed    = "closed"
	StateOpen      = "open"
	StateHalfOpen  = "half-open"
)

// CircuitBreaker prevents cascading failures
type CircuitBreaker struct {
	mu              sync.RWMutex
	state           string
	failures        int
	lastFailureTime time.Time
	successCount    int

	// Configuration
	maxFailures      int
	timeout          time.Duration
	halfOpenRequests int
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(maxFailures int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:            StateClosed,
		maxFailures:      maxFailures,
		timeout:          timeout,
		halfOpenRequests: 3, // Allow 3 requests in half-open state
	}
}

// Call executes the function with circuit breaker protection
func (cb *CircuitBreaker) Call(ctx context.Context, fn func() error) error {
	cb.mu.RLock()
	state := cb.state
	cb.mu.RUnlock()

	switch state {
	case StateOpen:
		// Check if we should transition to half-open
		cb.mu.RLock()
		if time.Since(cb.lastFailureTime) > cb.timeout {
			cb.mu.RUnlock()
			cb.transitionToHalfOpen()
		} else {
			cb.mu.RUnlock()
			return errors.ErrUnavailable.Msg("circuit breaker is open")
		}

	case StateHalfOpen:
		// Allow limited requests through
		cb.mu.RLock()
		if cb.successCount >= cb.halfOpenRequests {
			cb.mu.RUnlock()
			cb.transitionToClosed()
		} else {
			cb.mu.RUnlock()
		}
	}

	// Execute the function
	err := fn()

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failures++
		cb.lastFailureTime = time.Now()
		cb.successCount = 0

		if cb.failures >= cb.maxFailures {
			cb.state = StateOpen
		}
		return err
	}

	// Success
	if cb.state == StateHalfOpen {
		cb.successCount++
		if cb.successCount >= cb.halfOpenRequests {
			cb.state = StateClosed
			cb.failures = 0
		}
	} else {
		cb.failures = 0
	}

	return nil
}

func (cb *CircuitBreaker) transitionToHalfOpen() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.state = StateHalfOpen
	cb.successCount = 0
}

func (cb *CircuitBreaker) transitionToClosed() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.state = StateClosed
	cb.failures = 0
	cb.successCount = 0
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() string {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}