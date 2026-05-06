package consciousness

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
	
	"github.com/rs/zerolog"
)

// CircuitState represents the state of a circuit breaker
type CircuitState int

const (
	StateClosed CircuitState = iota
	StateOpen
	StateHalfOpen
)

var (
	ErrCircuitOpen = errors.New("circuit breaker is open")
	ErrTooManyRequests = errors.New("too many requests in half-open state")
)

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	name               string
	maxFailures        int
	resetTimeout       time.Duration
	successThreshold   int
	timeout            time.Duration
	
	mu                 sync.RWMutex
	state              CircuitState
	failures           int
	successCount       int
	lastFailureTime    time.Time
	lastStateChange    time.Time
	
	logger             zerolog.Logger
	onStateChange      func(name string, from, to CircuitState)
}

// CircuitBreakerConfig holds configuration for circuit breaker
type CircuitBreakerConfig struct {
	Name             string
	MaxFailures      int
	ResetTimeout     time.Duration
	SuccessThreshold int
	Timeout          time.Duration
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(config CircuitBreakerConfig, logger zerolog.Logger) *CircuitBreaker {
	return &CircuitBreaker{
		name:             config.Name,
		maxFailures:      config.MaxFailures,
		resetTimeout:     config.ResetTimeout,
		successThreshold: config.SuccessThreshold,
		timeout:          config.Timeout,
		state:            StateClosed,
		lastStateChange:  time.Now(),
		logger:           logger,
	}
}

// Execute runs a function with circuit breaker protection
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func(context.Context) error) error {
	if err := cb.canExecute(); err != nil {
		return err
	}
	
	// Create timeout context if configured
	var cancel context.CancelFunc
	if cb.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, cb.timeout)
		defer cancel()
	}
	
	// Execute the function
	err := fn(ctx)
	
	// Update circuit breaker state based on result
	cb.recordResult(err)
	
	return err
}

// canExecute checks if the circuit breaker allows execution
func (cb *CircuitBreaker) canExecute() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	
	now := time.Now()
	
	switch cb.state {
	case StateOpen:
		// Check if we should transition to half-open
		if now.Sub(cb.lastFailureTime) > cb.resetTimeout {
			cb.transitionTo(StateHalfOpen)
			cb.logger.Info().
				Str("circuit", cb.name).
				Msg("Circuit breaker transitioning to half-open")
			return nil
		}
		return fmt.Errorf("%w: %s", ErrCircuitOpen, cb.name)
		
	case StateHalfOpen:
		// In half-open state, we allow limited requests
		return nil
		
	case StateClosed:
		return nil
		
	default:
		return fmt.Errorf("unknown circuit state: %v", cb.state)
	}
}

// recordResult records the result of an execution
func (cb *CircuitBreaker) recordResult(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	
	if err != nil {
		cb.recordFailure()
	} else {
		cb.recordSuccess()
	}
}

// recordFailure records a failure
func (cb *CircuitBreaker) recordFailure() {
	cb.failures++
	cb.lastFailureTime = time.Now()
	cb.successCount = 0
	
	cb.logger.Warn().
		Str("circuit", cb.name).
		Int("failures", cb.failures).
		Int("max_failures", cb.maxFailures).
		Msg("Circuit breaker recorded failure")
	
	switch cb.state {
	case StateClosed:
		if cb.failures >= cb.maxFailures {
			cb.transitionTo(StateOpen)
			cb.logger.Error().
				Str("circuit", cb.name).
				Msg("Circuit breaker opened due to excessive failures")
		}
		
	case StateHalfOpen:
		// Single failure in half-open state opens the circuit
		cb.transitionTo(StateOpen)
		cb.logger.Warn().
			Str("circuit", cb.name).
			Msg("Circuit breaker reopened due to failure in half-open state")
	}
}

// recordSuccess records a success
func (cb *CircuitBreaker) recordSuccess() {
	cb.successCount++
	
	switch cb.state {
	case StateHalfOpen:
		cb.logger.Info().
			Str("circuit", cb.name).
			Int("successes", cb.successCount).
			Int("threshold", cb.successThreshold).
			Msg("Circuit breaker recorded success in half-open state")
			
		if cb.successCount >= cb.successThreshold {
			cb.failures = 0
			cb.transitionTo(StateClosed)
			cb.logger.Info().
				Str("circuit", cb.name).
				Msg("Circuit breaker closed after successful recovery")
		}
		
	case StateClosed:
		// Reset failure count on success in closed state
		cb.failures = 0
	}
}

// transitionTo changes the circuit breaker state
func (cb *CircuitBreaker) transitionTo(newState CircuitState) {
	if cb.state == newState {
		return
	}
	
	oldState := cb.state
	cb.state = newState
	cb.lastStateChange = time.Now()
	
	if cb.onStateChange != nil {
		cb.onStateChange(cb.name, oldState, newState)
	}
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// GetMetrics returns circuit breaker metrics
func (cb *CircuitBreaker) GetMetrics() CircuitBreakerMetrics {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	
	return CircuitBreakerMetrics{
		Name:            cb.name,
		State:           cb.state.String(),
		Failures:        cb.failures,
		SuccessCount:    cb.successCount,
		LastFailureTime: cb.lastFailureTime,
		LastStateChange: cb.lastStateChange,
	}
}

// CircuitBreakerMetrics contains circuit breaker statistics
type CircuitBreakerMetrics struct {
	Name            string
	State           string
	Failures        int
	SuccessCount    int
	LastFailureTime time.Time
	LastStateChange time.Time
}

// String returns string representation of circuit state
func (cs CircuitState) String() string {
	switch cs {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// CircuitBreakerManager manages multiple circuit breakers
type CircuitBreakerManager struct {
	breakers map[string]*CircuitBreaker
	mu       sync.RWMutex
	logger   zerolog.Logger
}

// NewCircuitBreakerManager creates a new circuit breaker manager
func NewCircuitBreakerManager(logger zerolog.Logger) *CircuitBreakerManager {
	return &CircuitBreakerManager{
		breakers: make(map[string]*CircuitBreaker),
		logger:   logger,
	}
}

// GetOrCreate gets an existing circuit breaker or creates a new one
func (cbm *CircuitBreakerManager) GetOrCreate(config CircuitBreakerConfig) *CircuitBreaker {
	cbm.mu.Lock()
	defer cbm.mu.Unlock()
	
	if cb, exists := cbm.breakers[config.Name]; exists {
		return cb
	}
	
	cb := NewCircuitBreaker(config, cbm.logger)
	cb.onStateChange = cbm.onStateChange
	cbm.breakers[config.Name] = cb
	
	return cb
}

// Get returns a circuit breaker by name
func (cbm *CircuitBreakerManager) Get(name string) (*CircuitBreaker, bool) {
	cbm.mu.RLock()
	defer cbm.mu.RUnlock()
	
	cb, exists := cbm.breakers[name]
	return cb, exists
}

// GetAllMetrics returns metrics for all circuit breakers
func (cbm *CircuitBreakerManager) GetAllMetrics() []CircuitBreakerMetrics {
	cbm.mu.RLock()
	defer cbm.mu.RUnlock()
	
	metrics := make([]CircuitBreakerMetrics, 0, len(cbm.breakers))
	for _, cb := range cbm.breakers {
		metrics = append(metrics, cb.GetMetrics())
	}
	
	return metrics
}

// onStateChange handles circuit breaker state changes
func (cbm *CircuitBreakerManager) onStateChange(name string, from, to CircuitState) {
	cbm.logger.Info().
		Str("circuit", name).
		Str("from_state", from.String()).
		Str("to_state", to.String()).
		Msg("Circuit breaker state changed")
}