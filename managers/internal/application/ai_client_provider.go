package application

import (
	"context"
	"fmt"
	"log"
	"math"
	"sync"
	"sync/atomic"
	"time"

	// Using local services package
	"middleman/managers/internal/application/services" // For services.AIClientProvider interface
	"middleman/internal/ai"                           // For ai.EnhancedAIService and other AI types

	// Using standard library errors or a common error package like github.com/pkg/errors
	// For this example, I'll use standard fmt.Errorf and check with errors.Is for context errors.
	stdErrors "errors" // Alias to avoid conflict if you have a local errors package
)

// --- CircuitBreaker Implementation ---

// CircuitBreakerState represents the state of the circuit breaker
type CircuitBreakerState int

const (
	StateClosed CircuitBreakerState = iota
	StateHalfOpen
	StateOpen
)

// CircuitBreakerConfig holds configuration for the circuit breaker
type CircuitBreakerConfig struct {
	Name             string        // Optional name for logging/identification
	MaxFailures      int           // Maximum consecutive failures before opening the circuit
	ResetTimeout     time.Duration // Time to wait in Open state before transitioning to HalfOpen
	SuccessThreshold int           // Number of consecutive successes in HalfOpen to close the circuit
	Timeout          time.Duration // Timeout for individual requests executed through the breaker
}

// CircuitBreakerMetrics holds metrics for monitoring the circuit breaker's performance
type CircuitBreakerMetrics struct {
	TotalRequests        int64
	TotalSuccesses       int64
	TotalFailures        int64
	ConsecutiveFailures  int // Renamed from 'failures' in CB struct for clarity
	ConsecutiveSuccesses int // Renamed from 'successes' in CB struct
	TotalTimeouts        int64
	StateChanges         int64
	LastStateChange      time.Time
	LastFailureTime      time.Time
}

// CircuitBreaker implements a production-ready circuit breaker pattern
type CircuitBreaker struct {
	config  CircuitBreakerConfig
	state   CircuitBreakerState
	metrics CircuitBreakerMetrics
	mutex   sync.RWMutex // RWMutex for read-heavy operations like GetState/IsHealthy
}

// RequestFunc represents a function that the circuit breaker will protect
type RequestFunc func(ctx context.Context) (interface{}, error)

// NewCircuitBreaker creates a new production-ready circuit breaker
func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	// Apply default values if not specified
	if config.MaxFailures <= 0 {
		config.MaxFailures = 3 // Default: 3 consecutive failures to open
	}
	if config.ResetTimeout <= 0 {
		config.ResetTimeout = 120 * time.Second // Default: 2 minutes in open state
	}
	if config.SuccessThreshold <= 0 {
		config.SuccessThreshold = 2 // Default: 2 successes to close from half-open
	}
	if config.Timeout <= 0 {
		config.Timeout = 600 * time.Second // Default: 10 minutes timeout for wrapped requests
	}

	cbName := config.Name
	if cbName == "" {
		cbName = "unnamed"
	}
	log.Printf("INFO: Initializing CircuitBreaker '%s' with config: MaxFailures=%d, ResetTimeout=%s, SuccessThreshold=%d, OpTimeout=%s",
		cbName, config.MaxFailures, config.ResetTimeout, config.SuccessThreshold, config.Timeout)

	return &CircuitBreaker{
		config: config,
		state:  StateClosed,
		metrics: CircuitBreakerMetrics{
			LastStateChange: time.Now(),
		},
	}
}

// Execute executes the request function with circuit breaker protection
func (cb *CircuitBreaker) Execute(ctx context.Context, fn RequestFunc) (interface{}, error) {
	cb.mutex.Lock() // Lock for state check and metrics update
	cb.metrics.TotalRequests++

	if cb.state == StateOpen {
		// Check if reset timeout has passed to transition to HalfOpen
		if time.Since(cb.metrics.LastFailureTime) > cb.config.ResetTimeout {
			cb.setStateInternal(StateHalfOpen) // Internal state change, no unlock yet
			log.Printf("INFO: CircuitBreaker '%s': Reset timeout elapsed, transitioning to HALF_OPEN", cb.config.Name)
		} else {
			// Still in Open state, request is rejected
			cb.mutex.Unlock()
			errMsg := fmt.Errorf("circuit breaker '%s' is OPEN. Last failure %s ago", cb.config.Name, time.Since(cb.metrics.LastFailureTime).Round(time.Millisecond))
			log.Printf("WARN: CircuitBreaker '%s': Request rejected. %s", cb.config.Name, errMsg.Error())
			return nil, errMsg
		}
	}
	// For Closed or HalfOpen state, proceed with the request
	cb.mutex.Unlock() // Unlock before executing the potentially long-running function

	// Create a new context with the operation timeout
	opCtx, cancel := context.WithTimeout(ctx, cb.config.Timeout)
	defer cancel()

	// Execute the function in a separate goroutine to handle its timeout vs context timeout
	var result interface{}
	var opErr error
	done := make(chan struct{})

	go func() {
		defer close(done)
		result, opErr = fn(opCtx)
	}()

	select {
	case <-done: // Operation finished (successfully or with error)
		if opErr != nil {
			// Check if the error was due to our operation context timing out
			if stdErrors.Is(opCtx.Err(), context.DeadlineExceeded) && stdErrors.Is(opErr, context.DeadlineExceeded) {
				cb.onTimeout() // Record specifically as a timeout
				log.Printf("WARN: CircuitBreaker '%s': Operation timed out after %s", cb.config.Name, cb.config.Timeout)
				return nil, fmt.Errorf("circuit breaker '%s': request timed out after %s: %w", cb.config.Name, cb.config.Timeout, opErr)
			}
			cb.onFailure(opErr)  // Record as general failure
			return result, opErr // Return the original error and result (if any)
		}
		cb.onSuccess()
		return result, nil

	case <-opCtx.Done(): // Operation context timed out
		// This case is primarily for safety if `fn` doesn't respect its context quickly enough,
		// or if the select on `done` misses due to goroutine scheduling.
		// The error from `fn` (if opCtx.Err()) should be caught in the goroutine.
		cb.onTimeout()
		timeoutErr := fmt.Errorf("circuit breaker '%s': request timed out after %s (context deadline exceeded)", cb.config.Name, cb.config.Timeout)
		log.Printf("WARN: %s", timeoutErr.Error())
		return nil, timeoutErr
	}
}

// onSuccess handles successful request completion
func (cb *CircuitBreaker) onSuccess() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.metrics.TotalSuccesses++

	switch cb.state {
	case StateHalfOpen:
		cb.metrics.ConsecutiveSuccesses++
		if cb.metrics.ConsecutiveSuccesses >= cb.config.SuccessThreshold {
			cb.setStateInternal(StateClosed)
			log.Printf("INFO: CircuitBreaker '%s': Success threshold reached in HALF_OPEN, transitioning to CLOSED", cb.config.Name)
		}
	case StateClosed:
		// Reset consecutive failures if any had occurred before success
		if cb.metrics.ConsecutiveFailures > 0 {
			log.Printf("INFO: CircuitBreaker '%s': Request successful in CLOSED state, resetting consecutive failures from %d to 0", cb.config.Name, cb.metrics.ConsecutiveFailures)
			cb.metrics.ConsecutiveFailures = 0
		}
	case StateOpen:
		// This case should ideally not be reached if canProceed works correctly,
		// but as a safeguard:
		log.Printf("WARN: CircuitBreaker '%s': onSuccess called while in OPEN state. Resetting to CLOSED.", cb.config.Name)
		cb.setStateInternal(StateClosed)
	}
}

// onFailure handles failed request completion (excluding timeouts handled by onTimeout)
func (cb *CircuitBreaker) onFailure(err error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	log.Printf("DEBUG: CircuitBreaker '%s': onFailure called with error: %v. Current state: %s, ConsecutiveFailures: %d", cb.config.Name, err, cb.GetState().String(), cb.metrics.ConsecutiveFailures)

	cb.metrics.TotalFailures++
	cb.metrics.LastFailureTime = time.Now()
	cb.metrics.ConsecutiveSuccesses = 0 // Reset consecutive successes on any failure in half-open

	switch cb.state {
	case StateClosed:
		cb.metrics.ConsecutiveFailures++
		if cb.metrics.ConsecutiveFailures >= cb.config.MaxFailures {
			cb.setStateInternal(StateOpen)
			log.Printf("WARN: CircuitBreaker '%s': Max failures (%d) reached in CLOSED state, transitioning to OPEN", cb.config.Name, cb.config.MaxFailures)
		}
	case StateHalfOpen:
		cb.setStateInternal(StateOpen) // Any failure in half-open trips it back to open
		log.Printf("WARN: CircuitBreaker '%s': Failure in HALF_OPEN state, transitioning back to OPEN", cb.config.Name)
	case StateOpen:
		// Already open, just record the failure time
		break
	}
}

// onTimeout handles request timeout specifically
func (cb *CircuitBreaker) onTimeout() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	log.Printf("DEBUG: CircuitBreaker '%s': onTimeout called. Current state: %s, ConsecutiveFailures: %d", cb.config.Name, cb.GetState().String(), cb.metrics.ConsecutiveFailures)

	cb.metrics.TotalTimeouts++
	cb.metrics.TotalFailures++ // Timeouts are considered failures
	cb.metrics.LastFailureTime = time.Now()
	cb.metrics.ConsecutiveSuccesses = 0 // Reset on timeout in half-open

	switch cb.state {
	case StateClosed:
		cb.metrics.ConsecutiveFailures++
		if cb.metrics.ConsecutiveFailures >= cb.config.MaxFailures {
			cb.setStateInternal(StateOpen)
			log.Printf("WARN: CircuitBreaker '%s': Max failures (%d, including timeout) reached in CLOSED state, transitioning to OPEN", cb.config.Name, cb.config.MaxFailures)
		}
	case StateHalfOpen:
		cb.setStateInternal(StateOpen) // Timeout in half-open trips it back to open
		log.Printf("WARN: CircuitBreaker '%s': Timeout in HALF_OPEN state, transitioning back to OPEN", cb.config.Name)
	case StateOpen:
		// Already open
		break
	}
}

// setStateInternal changes the circuit breaker state. Caller must hold the mutex.
func (cb *CircuitBreaker) setStateInternal(newState CircuitBreakerState) {
	if cb.state == newState {
		return
	}
	previousState := cb.state
	cb.state = newState
	cb.metrics.StateChanges++
	cb.metrics.LastStateChange = time.Now()

	log.Printf("INFO: CircuitBreaker '%s': State changed from %s to %s", cb.config.Name, previousState.String(), newState.String())

	// Reset counters when state changes
	if newState == StateClosed {
		cb.metrics.ConsecutiveFailures = 0
		cb.metrics.ConsecutiveSuccesses = 0
	} else if newState == StateHalfOpen {
		cb.metrics.ConsecutiveFailures = 0 // Failures leading to Open are reset when moving to HalfOpen
		cb.metrics.ConsecutiveSuccesses = 0
	} else if newState == StateOpen {
		// ConsecutiveFailures has already been incremented to MaxFailures
		cb.metrics.ConsecutiveSuccesses = 0
	}
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() CircuitBreakerState {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.state
}

// GetMetrics returns current metrics
func (cb *CircuitBreaker) GetMetrics() CircuitBreakerMetrics {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	// Return a copy to prevent race conditions if the caller modifies it
	metricsCopy := cb.metrics
	return metricsCopy
}

// IsHealthy returns true if the circuit breaker is likely to allow requests.
// Note: In HalfOpen state, it allows a limited number of requests.
func (cb *CircuitBreaker) IsHealthy() bool {
	cb.mutex.RLock() // Use RLock for read-only access
	currentState := cb.state
	lastFailTime := cb.metrics.LastFailureTime // Use metrics' lastFailureTime
	cb.mutex.RUnlock()

	if currentState == StateClosed {
		return true
	}
	if currentState == StateHalfOpen {
		return true // Allows test requests
	}
	// For StateOpen, check if it's time to transition to HalfOpen
	if currentState == StateOpen {
		// This check is also done in Execute, but IsHealthy provides a snapshot.
		// For a more accurate "can I make a call right now", rely on Execute.
		return time.Since(lastFailTime) > cb.config.ResetTimeout
	}
	return false // Should not happen with defined states
}

// Reset manually resets the circuit breaker to closed state
func (cb *CircuitBreaker) Reset() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	log.Printf("INFO: CircuitBreaker '%s': Manually reset to CLOSED state.", cb.config.Name)
	cb.setStateInternal(StateClosed)
}

// String returns a string representation of the current state
func (s CircuitBreakerState) String() string {
	switch s {
	case StateClosed:
		return "CLOSED"
	case StateHalfOpen:
		return "HALF_OPEN"
	case StateOpen:
		return "OPEN"
	default:
		return fmt.Sprintf("UNKNOWN_STATE (%d)", s)
	}
}

// --- AIClientProvider Implementation ---

// AIClientProvider provides access to different AI clients with production features.
// It implements the services.AIClientProvider interface.
type AIClientProvider struct {
	clients         map[string]ai.EnhancedAIService
	defaultProvider string
	fallbackOrder   []string // Order in which to try providers if the default fails
	circuitBreakers map[string]*CircuitBreaker
	metrics         map[string]*ProviderMetrics
	mutex           sync.RWMutex // Protects access to clients, circuitBreakers, metrics, and config fields

	// Enhanced capabilities
	performanceScores  map[string]*PerformanceScore
	requestTypeRouting map[string][]string // Request type -> preferred providers
	loadBalancer       *IntelligentLoadBalancer
	costOptimizer      *CostOptimizer
	healthMonitor      *ProviderHealthMonitor

	// Communication optimization features
	clientCache      map[string]*CachedClientInfo
	requestTracker   *RequestTracker
	adaptiveTimeouts map[string]*AdaptiveTimeout
	commOptMutex     sync.RWMutex
}

// CachedClientInfo represents cached client with optimization metadata
type CachedClientInfo struct {
	Client          ai.EnhancedAIService
	LastUsed        time.Time
	UseCount        int64
	IsHealthy       bool
	LastHealthCheck time.Time
	CreatedAt       time.Time
}

// RequestTracker monitors request patterns for optimization
type RequestTracker struct {
	activeRequests         sync.Map
	requestHistory         []RequestMetrics
	historyMutex           sync.RWMutex
	consecutiveSlowCount   int64
	lastSlowRequestTime    time.Time
	performanceDegradation bool
}

// RequestMetrics stores detailed performance data
type RequestMetrics struct {
	RequestID      string
	StartTime      time.Time
	EndTime        time.Time
	Latency        time.Duration
	Provider       string
	ToolsRequested int
	ToolsExecuted  int
	ResponseLength int
	Success        bool
	CacheHit       bool
	RequestType    string
}

// AdaptiveTimeout manages dynamic timeout adjustments
type AdaptiveTimeout struct {
	baseTimeout     time.Duration
	currentTimeout  time.Duration
	recentLatencies []time.Duration
	successRate     float64
	mutex           sync.RWMutex
	lastAdjustment  time.Time
}

// PerformanceScore tracks dynamic performance metrics for intelligent routing
type PerformanceScore struct {
	AverageLatency    time.Duration      `json:"average_latency"`
	SuccessRate       float64            `json:"success_rate"`
	RecentErrors      []string           `json:"recent_errors"`
	TokensPerSecond   float64            `json:"tokens_per_second"`
	CostPerToken      float64            `json:"cost_per_token"`
	QualityScore      float64            `json:"quality_score"`     // Response quality (0-1)
	ReliabilityScore  float64            `json:"reliability_score"` // Uptime and stability (0-1)
	LastUpdated       time.Time          `json:"last_updated"`
	CompositeScore    float64            `json:"composite_score"`     // Weighted overall score
	RequestTypeScores map[string]float64 `json:"request_type_scores"` // Performance per request type
}

// IntelligentLoadBalancer handles smart distribution of requests
type IntelligentLoadBalancer struct {
	strategy          LoadBalancingStrategy
	weightedProviders map[string]float64 // Provider -> weight based on performance
	requestCounts     map[string]int64   // Provider -> current load
	mutex             sync.RWMutex
}

type LoadBalancingStrategy string

const (
	StrategyRoundRobin       LoadBalancingStrategy = "round_robin"
	StrategyWeightedRandom   LoadBalancingStrategy = "weighted_random"
	StrategyPerformance      LoadBalancingStrategy = "performance"
	StrategyLeastConnections LoadBalancingStrategy = "least_connections"
)

// CostOptimizer handles intelligent cost optimization
type CostOptimizer struct {
	costThresholds map[string]float64 // Provider -> cost threshold
	budgetLimits   map[string]float64 // Daily/hourly budget limits
	costTracking   map[string]float64 // Current cost tracking
	lastReset      time.Time
}

// ProviderHealthMonitor monitors provider health in real-time
type ProviderHealthMonitor struct {
	healthChecks     map[string]*HealthCheck
	checkInterval    time.Duration
	unhealthyTimeout time.Duration
	mutex            sync.RWMutex
}

type HealthCheck struct {
	Provider     string        `json:"provider"`
	IsHealthy    bool          `json:"is_healthy"`
	LastCheck    time.Time     `json:"last_check"`
	ResponseTime time.Duration `json:"response_time"`
	ErrorRate    float64       `json:"error_rate"`
	Availability float64       `json:"availability"`
}

// Enhanced ProviderMetrics with more detailed tracking
type ProviderMetrics struct {
	TotalRequests      int64
	SuccessfulRequests int64
	FailedRequests     int64
	TotalLatency       time.Duration
	AverageLatency     time.Duration
	LastUsed           time.Time
	IsCurrentlyHealthy bool
	LastError          string
	LastSuccess        time.Time
	ProviderName       string

	// Enhanced metrics
	RequestTypeBreakdown map[string]*RequestTypeMetrics `json:"request_type_breakdown"`
	HourlyStats          map[int]*HourlyMetrics         `json:"hourly_stats"`
	TokensProcessed      int64                          `json:"tokens_processed"`
	CostAccumulated      float64                        `json:"cost_accumulated"`
	QualityRatings       []float64                      `json:"quality_ratings"`
}

type RequestTypeMetrics struct {
	Count        int64         `json:"count"`
	SuccessRate  float64       `json:"success_rate"`
	AvgLatency   time.Duration `json:"avg_latency"`
	AvgTokens    float64       `json:"avg_tokens"`
	QualityScore float64       `json:"quality_score"`
}

type HourlyMetrics struct {
	Hour         int           `json:"hour"`
	RequestCount int64         `json:"request_count"`
	SuccessRate  float64       `json:"success_rate"`
	AvgLatency   time.Duration `json:"avg_latency"`
	Cost         float64       `json:"cost"`
}

func NewAIClientProvider(
	clientMap map[string]ai.EnhancedAIService,
	defaultProviderName string,
	fallbackOrder []string,
	cbConfig CircuitBreakerConfig,
) (services.AIClientProvider, error) {
	if len(clientMap) == 0 {
		return nil, fmt.Errorf("clientMap cannot be empty")
	}

	// Validate that the defaultProviderName exists in the clientMap
	if _, exists := clientMap[defaultProviderName]; !exists {
		return nil, fmt.Errorf("default provider '%s' not found in client map", defaultProviderName)
	}

	// Validate that all fallback providers exist
	for _, provider := range fallbackOrder {
		if _, exists := clientMap[provider]; !exists {
			return nil, fmt.Errorf("fallback provider '%s' not found in client map", provider)
		}
	}

	// If fallbackOrder is empty or doesn't include the default provider, prepend it
	if len(fallbackOrder) == 0 || fallbackOrder[0] != defaultProviderName {
		fallbackOrder = append([]string{defaultProviderName}, fallbackOrder...)
	}

	circuitBreakers := make(map[string]*CircuitBreaker)
	metrics := make(map[string]*ProviderMetrics)
	performanceScores := make(map[string]*PerformanceScore)

	for providerName := range clientMap {
		// Create circuit breaker for each provider
		providerCBConfig := cbConfig
		providerCBConfig.Name = fmt.Sprintf("%s-%s", cbConfig.Name, providerName)

		circuitBreakers[providerName] = NewCircuitBreaker(providerCBConfig)

		// Initialize metrics for each provider
		metrics[providerName] = &ProviderMetrics{
			ProviderName:         providerName,
			IsCurrentlyHealthy:   true,
			RequestTypeBreakdown: make(map[string]*RequestTypeMetrics),
			HourlyStats:          make(map[int]*HourlyMetrics),
			QualityRatings:       make([]float64, 0),
		}

		// Initialize performance scores
		performanceScores[providerName] = &PerformanceScore{
			SuccessRate:       1.0,
			QualityScore:      0.8,
			ReliabilityScore:  1.0,
			CompositeScore:    0.8,
			RequestTypeScores: make(map[string]float64),
			LastUpdated:       time.Now(),
		}
	}

	// Initialize intelligent load balancer
	loadBalancer := &IntelligentLoadBalancer{
		strategy:          StrategyPerformance,
		weightedProviders: make(map[string]float64),
		requestCounts:     make(map[string]int64),
	}

	// Initialize weights based on initial performance scores
	for provider := range clientMap {
		loadBalancer.weightedProviders[provider] = 1.0 // Equal weights initially
	}

	// Initialize cost optimizer
	costOptimizer := &CostOptimizer{
		costThresholds: map[string]float64{
			"openai":    0.02,  // $0.02 per 1K tokens
			"anthropic": 0.015, // $0.015 per 1K tokens
			"deepseek":  0.001, // $0.001 per 1K tokens
		},
		budgetLimits: map[string]float64{
			"daily":  100.0, // $100 daily limit
			"hourly": 10.0,  // $10 hourly limit
		},
		costTracking: make(map[string]float64),
		lastReset:    time.Now(),
	}

	// Initialize health monitor
	healthMonitor := &ProviderHealthMonitor{
		healthChecks:     make(map[string]*HealthCheck),
		checkInterval:    time.Minute * 5,
		unhealthyTimeout: time.Minute * 15,
	}

	for provider := range clientMap {
		healthMonitor.healthChecks[provider] = &HealthCheck{
			Provider:     provider,
			IsHealthy:    true,
			LastCheck:    time.Now(),
			Availability: 1.0,
		}
	}

	// Initialize request type routing preferences
	requestTypeRouting := map[string][]string{
		"text_generation":    {"openai", "anthropic", "deepseek"},
		"code_generation":    {"openai", "deepseek", "anthropic"},
		"image_analysis":     {"openai", "anthropic"},
		"complex_reasoning":  {"anthropic", "openai", "deepseek"},
		"data_analysis":      {"openai", "anthropic", "deepseek"},
		"creative_writing":   {"anthropic", "openai", "deepseek"},
		"technical_writing":  {"openai", "anthropic", "deepseek"},
		"translation":        {"openai", "anthropic", "deepseek"},
		"summarization":      {"anthropic", "openai", "deepseek"},
		"question_answering": {"openai", "anthropic", "deepseek"},
	}

	// Initialize communication optimization features
	clientCache := make(map[string]*CachedClientInfo)
	adaptiveTimeouts := make(map[string]*AdaptiveTimeout)

	for provider := range clientMap {
		adaptiveTimeouts[provider] = &AdaptiveTimeout{
			baseTimeout:     30 * time.Second,
			currentTimeout:  30 * time.Second,
			recentLatencies: make([]time.Duration, 0, 10),
			successRate:     1.0,
			lastAdjustment:  time.Now(),
		}
	}

	requestTracker := &RequestTracker{
		requestHistory: make([]RequestMetrics, 0, 1000),
	}

	provider := &AIClientProvider{
		clients:            clientMap,
		defaultProvider:    defaultProviderName,
		fallbackOrder:      fallbackOrder,
		circuitBreakers:    circuitBreakers,
		metrics:            metrics,
		performanceScores:  performanceScores,
		requestTypeRouting: requestTypeRouting,
		loadBalancer:       loadBalancer,
		costOptimizer:      costOptimizer,
		healthMonitor:      healthMonitor,
		clientCache:        clientCache,
		requestTracker:     requestTracker,
		adaptiveTimeouts:   adaptiveTimeouts,
	}

	// Start background health monitoring
	go provider.startHealthMonitoring()
	go provider.startPerformanceScoring()

	log.Printf("INFO: AIClientProvider initialized with %d clients. Default: %s, Fallback order: %v",
		len(clientMap), defaultProviderName, fallbackOrder)

	return provider, nil
}

// GetClient returns a specific AI client by provider name, wrapped with protection.
func (p *AIClientProvider) GetClient(ctx context.Context, provider string) (ai.EnhancedAIService, error) {
	p.mutex.RLock()
	client, clientExists := p.clients[provider]
	circuitBreaker, cbExists := p.circuitBreakers[provider]
	p.mutex.RUnlock()

	if !clientExists {
		return nil, fmt.Errorf("AIClientProvider: unknown AI provider '%s'", provider)
	}
	if !cbExists {
		log.Printf("CRITICAL: AIClientProvider: Circuit breaker for provider '%s' not found. Initialization error.", provider)
		return nil, fmt.Errorf("AIClientProvider: circuit breaker critically missing for provider '%s'", provider)
	}

	// The ProtectedAIClient wrapper will use the circuit breaker's IsHealthy() before Execute.
	// No need to explicitly block here based on circuitBreaker.IsHealthy() as GetClient might be called
	// for non-Execute operations like GetCapabilities.
	return &ProtectedAIClient{
		client:         client,
		provider:       provider,
		circuitBreaker: circuitBreaker,
		clientProvider: p,
	}, nil
}

// GetDefaultClient returns the default AI client, trying fallbacks if necessary.
func (p *AIClientProvider) GetDefaultClient(ctx context.Context) (ai.EnhancedAIService, error) {
	log.Printf("INFO: AIClientProvider: Attempting to get default client for provider '%s'", p.defaultProvider)
	client, err := p.GetClient(ctx, p.defaultProvider)
	if err == nil {
		// Further check if the circuit breaker for this client is healthy *before* returning
		p.mutex.RLock()
		cb := p.circuitBreakers[p.defaultProvider]
		p.mutex.RUnlock()
		if cb != nil && cb.IsHealthy() { // Check if CB allows requests
			log.Printf("INFO: AIClientProvider: Using default provider '%s'", p.defaultProvider)
			return client, nil
		}
		if cb != nil && !cb.IsHealthy() {
			log.Printf("WARN: AIClientProvider: Default provider '%s' circuit breaker is not healthy. Attempting fallbacks.", p.defaultProvider)
		} else if cb == nil {
			log.Printf("CRITICAL: AIClientProvider: Default provider '%s' circuit breaker is nil. Attempting fallbacks.", p.defaultProvider)
		}
	}
	if err != nil { // Error from GetClient (e.g., provider not found, which shouldn't happen for default)
		log.Printf("WARN: AIClientProvider: Getting default provider '%s' failed: %v. Attempting fallbacks.", p.defaultProvider, err)
	}

	p.mutex.RLock() // Lock for reading fallbackOrder
	currentFallbackOrder := make([]string, len(p.fallbackOrder))
	copy(currentFallbackOrder, p.fallbackOrder)
	p.mutex.RUnlock()

	for _, fallbackProviderName := range currentFallbackOrder {
		if fallbackProviderName == p.defaultProvider {
			continue
		}
		log.Printf("INFO: AIClientProvider: Attempting fallback to provider '%s'", fallbackProviderName)
		fallbackClient, fbErr := p.GetClient(ctx, fallbackProviderName)
		if fbErr == nil {
			p.mutex.RLock()
			cb := p.circuitBreakers[fallbackProviderName]
			p.mutex.RUnlock()
			if cb != nil && cb.IsHealthy() {
				log.Printf("INFO: AIClientProvider: Successfully fell back to provider '%s'", fallbackProviderName)
				return fallbackClient, nil
			}
			if cb != nil && !cb.IsHealthy() {
				log.Printf("WARN: AIClientProvider: Fallback provider '%s' circuit breaker is not healthy.", fallbackProviderName)
			} else if cb == nil {
				log.Printf("CRITICAL: AIClientProvider: Fallback provider '%s' circuit breaker is nil.", fallbackProviderName)
			}
		}
		if fbErr != nil {
			log.Printf("WARN: AIClientProvider: Getting fallback provider '%s' failed: %v", fallbackProviderName, fbErr)
		}
	}

	log.Printf("ERROR: AIClientProvider: All AI providers (default and fallbacks) are unavailable or non-operational.")
	return nil, fmt.Errorf("AIClientProvider: all AI providers are unavailable")
}

// GetHealthyProvider returns the first healthy provider based on fallback order.
func (p *AIClientProvider) GetHealthyProvider(ctx context.Context) (ai.EnhancedAIService, string, error) {
	p.mutex.RLock()
	currentFallbackOrder := make([]string, len(p.fallbackOrder))
	copy(currentFallbackOrder, p.fallbackOrder)
	p.mutex.RUnlock()

	for _, providerName := range currentFallbackOrder {
		p.mutex.RLock()
		cb, cbExists := p.circuitBreakers[providerName]
		p.mutex.RUnlock()

		if cbExists && cb.IsHealthy() {
			client, err := p.GetClient(ctx, providerName) // GetClient will return the ProtectedAIClient
			if err == nil {
				log.Printf("INFO: AIClientProvider: Found healthy provider '%s' via GetHealthyProvider.", providerName)
				return client, providerName, nil
			}
			log.Printf("WARN: AIClientProvider: Error getting client instance for healthy provider '%s': %v", providerName, err)
		}
	}

	log.Println("ERROR: AIClientProvider: No healthy AI providers available via GetHealthyProvider.")
	return nil, "", fmt.Errorf("AIClientProvider: no healthy AI providers available")
}

// recordCallMetric is called by ProtectedAIClient to update metrics after a call.
func (p *AIClientProvider) recordCallMetric(provider string, success bool, latency time.Duration, callErr error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	metric, exists := p.metrics[provider]
	if !exists {
		log.Printf("ERROR: AIClientProvider: Metrics for provider '%s' not found during recordCallMetric. Initializing.", provider)
		metric = &ProviderMetrics{ProviderName: provider}
		p.metrics[provider] = metric
	}

	metric.TotalRequests++
	metric.LastUsed = time.Now()
	metric.TotalLatency += latency

	cb, cbExists := p.circuitBreakers[provider]

	if success {
		metric.SuccessfulRequests++
		metric.LastError = ""
		metric.LastSuccess = time.Now()
		if cbExists {
			metric.IsCurrentlyHealthy = cb.IsHealthy()
		} else {
			metric.IsCurrentlyHealthy = true
		} // Assume healthy if CB missing (bad state)
	} else {
		metric.FailedRequests++
		if callErr != nil {
			metric.LastError = callErr.Error()
		} else {
			metric.LastError = "Unknown failure"
		}
		if cbExists {
			metric.IsCurrentlyHealthy = cb.IsHealthy()
		} else {
			metric.IsCurrentlyHealthy = false
		} // Assume unhealthy if CB missing on failure
	}

	if metric.TotalRequests > 0 {
		metric.AverageLatency = time.Duration(metric.TotalLatency.Nanoseconds() / metric.TotalRequests)
	}

	log.Printf("INFO: AIClientProvider: Metrics updated for '%s': Success: %t, Latency: %s, TotalReq: %d, AvgLat: %s, HealthyNow: %t, LastErr: %s",
		provider, success, latency, metric.TotalRequests, metric.AverageLatency, metric.IsCurrentlyHealthy, metric.LastError)
}

// GetMetrics returns a copy of metrics for all providers.
func (p *AIClientProvider) GetMetrics() map[string]*ProviderMetrics {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	result := make(map[string]*ProviderMetrics, len(p.metrics))
	for providerName, m := range p.metrics {
		metricCopy := *m // Create a copy
		result[providerName] = &metricCopy
	}
	return result
}

// ProtectedAIClient wraps an AI client, adding circuit breaker protection and metrics reporting.
type ProtectedAIClient struct {
	client         ai.EnhancedAIService
	provider       string
	circuitBreaker *CircuitBreaker
	clientProvider *AIClientProvider // Changed from pointer to interface to pointer to concrete type for recordCallMetric
}

// --- Implementations for ai.EnhancedAIService for ProtectedAIClient ---

// CreateCompletion wraps the underlying client's method.
func (p *ProtectedAIClient) CreateCompletion(ctx context.Context, request ai.CompletionRequest) (resp *ai.CompletionResponse, err error) {
	start := time.Now()
	// Defer metrics recording to ensure it always runs for this call
	defer func(callStart time.Time) {
		// `err` will have its final value when defer executes
		p.clientProvider.recordCallMetric(p.provider, err == nil, time.Since(callStart), err)
		if r := recover(); r != nil { // Catch panics from underlying client or CB Execute
			err = fmt.Errorf("panic in ProtectedAIClient.CreateCompletion for provider %s: %v", p.provider, r)
			log.Printf("CRITICAL: %v", err)
			// Ensure metrics are updated for panic as well
			p.clientProvider.recordCallMetric(p.provider, false, time.Since(callStart), err)
			// Re-throw panic if needed, or return error. For robustness, return error.
			// panic(r)
		}
	}(start)

	var result interface{}
	result, err = p.circuitBreaker.Execute(ctx, func(c context.Context) (interface{}, error) {
		return p.client.CreateCompletion(c, request)
	})

	if err != nil {
		return nil, err // Error already recorded by defer or from CB
	}
	typedResult, ok := result.(*ai.CompletionResponse)
	if !ok {
		err = fmt.Errorf("ProtectedAIClient: CreateCompletion received unexpected type from circuit breaker for provider %s: %T", p.provider, result)
		return nil, err // Error will be recorded by defer
	}
	return typedResult, nil
}

// CreateCompletionStream wraps the underlying client's method.
func (p *ProtectedAIClient) CreateCompletionStream(ctx context.Context, request ai.CompletionRequest) (stream <-chan ai.CompletionStreamResponse, err error) {
	start := time.Now()
	defer func(callStart time.Time) {
		p.clientProvider.recordCallMetric(p.provider, err == nil, time.Since(callStart), err)
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in ProtectedAIClient.CreateCompletionStream for provider %s: %v", p.provider, r)
			log.Printf("CRITICAL: %v", err)
			p.clientProvider.recordCallMetric(p.provider, false, time.Since(callStart), err)
		}
	}(start)

	var result interface{}
	result, err = p.circuitBreaker.Execute(ctx, func(c context.Context) (interface{}, error) {
		return p.client.CreateCompletionStream(c, request)
	})
	if err != nil {
		return nil, err
	}
	typedResult, ok := result.(<-chan ai.CompletionStreamResponse)
	if !ok {
		err = fmt.Errorf("ProtectedAIClient: CreateCompletionStream received unexpected type for provider %s: %T", p.provider, result)
		return nil, err
	}
	return typedResult, nil
}

// ExecuteWithTools wraps the underlying client's method.
func (p *ProtectedAIClient) ExecuteWithTools(ctx context.Context, messages []ai.Message, tools []ai.ToolDefinition) (resp *ai.CompletionResponse, err error) {
	start := time.Now()
	defer func(callStart time.Time) {
		p.clientProvider.recordCallMetric(p.provider, err == nil, time.Since(callStart), err)
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in ProtectedAIClient.ExecuteWithTools for provider %s: %v", p.provider, r)
			log.Printf("CRITICAL: %v", err)
			p.clientProvider.recordCallMetric(p.provider, false, time.Since(callStart), err)
		}
	}(start)

	var result interface{}
	result, err = p.circuitBreaker.Execute(ctx, func(c context.Context) (interface{}, error) {
		return p.client.ExecuteWithTools(c, messages, tools)
	})
	if err != nil {
		return nil, err
	}
	typedResult, ok := result.(*ai.CompletionResponse)
	if !ok {
		err = fmt.Errorf("ProtectedAIClient: ExecuteWithTools received unexpected type for provider %s: %T", p.provider, result)
		return nil, err
	}
	return typedResult, nil
}

// CreateStructuredCompletion wraps the underlying client's method.
func (p *ProtectedAIClient) CreateStructuredCompletion(ctx context.Context, request ai.CompletionRequest, schema *ai.JSONSchemaDefinition) (resp *ai.CompletionResponse, err error) {
	start := time.Now()
	defer func(callStart time.Time) {
		p.clientProvider.recordCallMetric(p.provider, err == nil, time.Since(callStart), err)
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in ProtectedAIClient.CreateStructuredCompletion for provider %s: %v", p.provider, r)
			log.Printf("CRITICAL: %v", err)
			p.clientProvider.recordCallMetric(p.provider, false, time.Since(callStart), err)
		}
	}(start)

	var result interface{}
	result, err = p.circuitBreaker.Execute(ctx, func(c context.Context) (interface{}, error) {
		return p.client.CreateStructuredCompletion(c, request, schema)
	})
	if err != nil {
		return nil, err
	}
	typedResult, ok := result.(*ai.CompletionResponse)
	if !ok {
		err = fmt.Errorf("ProtectedAIClient: CreateStructuredCompletion received unexpected type for provider %s: %T", p.provider, result)
		return nil, err
	}
	return typedResult, nil
}

// CountTokens is a direct passthrough.
func (p *ProtectedAIClient) CountTokens(text string) (int, error) {
	return p.client.CountTokens(text)
}

// GetCapabilities is a direct passthrough.
func (p *ProtectedAIClient) GetCapabilities() ai.ClientCapabilities {
	return p.client.GetCapabilities()
}

// HealthCheck is wrapped.
func (p *ProtectedAIClient) HealthCheck(ctx context.Context) (err error) {
	start := time.Now()
	healthCheckTimeout := 600 * time.Second // Specific timeout for health checks
	healthCheckCtx, cancel := context.WithTimeout(ctx, healthCheckTimeout)
	defer cancel()

	defer func(callStart time.Time) {
		// err will be set if Execute fails or underlying HealthCheck fails
		p.clientProvider.recordCallMetric(p.provider, err == nil, time.Since(callStart), err)
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in ProtectedAIClient.HealthCheck for provider %s: %v", p.provider, r)
			log.Printf("CRITICAL: %v", err)
			p.clientProvider.recordCallMetric(p.provider, false, time.Since(callStart), err) // Record panic as failure
		}
	}(start)

	_, err = p.circuitBreaker.Execute(healthCheckCtx, func(c context.Context) (interface{}, error) {
		return nil, p.client.HealthCheck(c)
	})
	return err // err will be from Execute or nil
}

// GetUsageStats is a direct passthrough.
func (p *ProtectedAIClient) GetUsageStats() ai.UsageStats {
	return p.client.GetUsageStats()
}

// GetProviderInfo is a direct passthrough.
func (p *ProtectedAIClient) GetProviderInfo() ai.ProviderInfo {
	return p.client.GetProviderInfo()
}

// --- Security feature wrappers ---
// These are wrapped similarly to CreateCompletion if they involve external calls.

func (p *ProtectedAIClient) AnalyzeFraud(ctx context.Context, data string) (resp *ai.SecurityAssessment, err error) {
	start := time.Now()
	defer func(cs time.Time) {
		p.clientProvider.recordCallMetric(p.provider, err == nil, time.Since(cs), err)
		if r := recover(); r != nil {
			err = fmt.Errorf("panic_analyze_fraud_%s: %v", p.provider, r)
			log.Printf("CRITICAL: %v", err)
			p.clientProvider.recordCallMetric(p.provider, false, time.Since(cs), err)
		}
	}(start)
	var res interface{}
	res, err = p.circuitBreaker.Execute(ctx, func(c context.Context) (interface{}, error) { return p.client.AnalyzeFraud(c, data) })
	if err != nil {
		return nil, err
	}
	if r, ok := res.(*ai.SecurityAssessment); ok {
		return r, nil
	}
	err = fmt.Errorf("unexpected type from AnalyzeFraud for provider %s: %T", p.provider, res)
	return nil, err
}

func (p *ProtectedAIClient) AssessRisk(ctx context.Context, request ai.CompletionRequest) (resp *ai.SecurityAssessment, err error) {
	start := time.Now()
	defer func(cs time.Time) {
		p.clientProvider.recordCallMetric(p.provider, err == nil, time.Since(cs), err)
		if r := recover(); r != nil {
			err = fmt.Errorf("panic_assess_risk_%s: %v", p.provider, r)
			log.Printf("CRITICAL: %v", err)
			p.clientProvider.recordCallMetric(p.provider, false, time.Since(cs), err)
		}
	}(start)
	var res interface{}
	res, err = p.circuitBreaker.Execute(ctx, func(c context.Context) (interface{}, error) { return p.client.AssessRisk(c, request) })
	if err != nil {
		return nil, err
	}
	if r, ok := res.(*ai.SecurityAssessment); ok {
		return r, nil
	}
	err = fmt.Errorf("unexpected type from AssessRisk for provider %s: %T", p.provider, res)
	return nil, err
}

func (p *ProtectedAIClient) GetSecurityRecommendations(ctx context.Context, content string) (resp []string, err error) {
	start := time.Now()
	defer func(cs time.Time) {
		p.clientProvider.recordCallMetric(p.provider, err == nil, time.Since(cs), err)
		if r := recover(); r != nil {
			err = fmt.Errorf("panic_get_sec_reco_%s: %v", p.provider, r)
			log.Printf("CRITICAL: %v", err)
			p.clientProvider.recordCallMetric(p.provider, false, time.Since(cs), err)
		}
	}(start)
	var res interface{}
	res, err = p.circuitBreaker.Execute(ctx, func(c context.Context) (interface{}, error) { return p.client.GetSecurityRecommendations(c, content) })
	if err != nil {
		return nil, err
	}
	if r, ok := res.([]string); ok {
		return r, nil
	}
	err = fmt.Errorf("unexpected type from GetSecurityRecommendations for provider %s: %T", p.provider, res)
	return nil, err
}

// startHealthMonitoring runs background health checks on all providers
func (p *AIClientProvider) startHealthMonitoring() {
	ticker := time.NewTicker(p.healthMonitor.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.performHealthChecks()
		}
	}
}

// performHealthChecks checks the health of all providers
func (p *AIClientProvider) performHealthChecks() {
	p.healthMonitor.mutex.Lock()
	defer p.healthMonitor.mutex.Unlock()

	for providerName, client := range p.clients {
		startTime := time.Now()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		err := client.HealthCheck(ctx)
		cancel()

		responseTime := time.Since(startTime)
		check := p.healthMonitor.healthChecks[providerName]

		if err != nil {
			check.IsHealthy = false
			check.ErrorRate = math.Min(1.0, check.ErrorRate+0.1)
			log.Printf("WARN: Health check failed for provider %s: %v", providerName, err)
		} else {
			check.IsHealthy = true
			check.ErrorRate = math.Max(0.0, check.ErrorRate-0.05)
		}

		check.LastCheck = time.Now()
		check.ResponseTime = responseTime

		// Update availability (exponential moving average)
		if err == nil {
			check.Availability = 0.9*check.Availability + 0.1*1.0
		} else {
			check.Availability = 0.9*check.Availability + 0.1*0.0
		}

		// Update provider metrics health status
		p.mutex.Lock()
		if metrics, exists := p.metrics[providerName]; exists {
			metrics.IsCurrentlyHealthy = check.IsHealthy
		}
		p.mutex.Unlock()
	}
}

// startPerformanceScoring continuously updates performance scores
func (p *AIClientProvider) startPerformanceScoring() {
	ticker := time.NewTicker(time.Minute * 2) // Update every 2 minutes
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.updatePerformanceScores()
		}
	}
}

// updatePerformanceScores recalculates performance scores for all providers
func (p *AIClientProvider) updatePerformanceScores() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for providerName, metrics := range p.metrics {
		score := p.performanceScores[providerName]

		// Calculate success rate
		if metrics.TotalRequests > 0 {
			score.SuccessRate = float64(metrics.SuccessfulRequests) / float64(metrics.TotalRequests)
		}

		// Update average latency
		score.AverageLatency = metrics.AverageLatency

		// Calculate tokens per second (rough estimate)
		if metrics.TotalLatency > 0 && metrics.TokensProcessed > 0 {
			score.TokensPerSecond = float64(metrics.TokensProcessed) / metrics.TotalLatency.Seconds()
		}

		// Get health metrics
		p.healthMonitor.mutex.RLock()
		if healthCheck, exists := p.healthMonitor.healthChecks[providerName]; exists {
			score.ReliabilityScore = healthCheck.Availability
		}
		p.healthMonitor.mutex.RUnlock()

		// Calculate quality score from recent ratings
		if len(metrics.QualityRatings) > 0 {
			sum := 0.0
			count := math.Min(10, float64(len(metrics.QualityRatings))) // Last 10 ratings
			for i := len(metrics.QualityRatings) - int(count); i < len(metrics.QualityRatings); i++ {
				sum += metrics.QualityRatings[i]
			}
			score.QualityScore = sum / count
		}

		// Calculate composite score (weighted)
		score.CompositeScore = 0.3*score.SuccessRate +
			0.2*score.ReliabilityScore +
			0.2*score.QualityScore +
			0.15*(1.0-math.Min(1.0, score.AverageLatency.Seconds()/10.0)) + // Latency penalty
			0.15*math.Min(1.0, score.TokensPerSecond/100.0) // Throughput bonus

		score.LastUpdated = time.Now()

		// Update load balancer weights
		p.loadBalancer.mutex.Lock()
		p.loadBalancer.weightedProviders[providerName] = score.CompositeScore
		p.loadBalancer.mutex.Unlock()
	}
}

// GetOptimalProvider returns the best provider for a specific request type
func (p *AIClientProvider) GetOptimalProvider(ctx context.Context, requestType string) (ai.EnhancedAIService, string, error) {
	// Check if we have specific routing for this request type
	if preferredProviders, exists := p.requestTypeRouting[requestType]; exists {
		for _, providerName := range preferredProviders {
			if client, err := p.tryProvider(ctx, providerName); err == nil {
				return client, providerName, nil
			}
		}
	}

	// Fall back to performance-based selection
	return p.getPerformanceBasedProvider(ctx)
}

// tryProvider attempts to get a specific provider if it's healthy
func (p *AIClientProvider) tryProvider(ctx context.Context, providerName string) (ai.EnhancedAIService, error) {
	p.mutex.RLock()
	client, exists := p.clients[providerName]
	circuitBreaker := p.circuitBreakers[providerName]
	p.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("provider %s not found", providerName)
	}

	if !circuitBreaker.IsHealthy() {
		return nil, fmt.Errorf("provider %s is unhealthy", providerName)
	}

	// Check cost limits
	if !p.checkCostLimits(providerName) {
		return nil, fmt.Errorf("provider %s exceeds cost limits", providerName)
	}

	return &ProtectedAIClient{
		client:         client,
		provider:       providerName,
		circuitBreaker: circuitBreaker,
		clientProvider: p,
	}, nil
}

// getPerformanceBasedProvider selects provider based on current performance scores
func (p *AIClientProvider) getPerformanceBasedProvider(ctx context.Context) (ai.EnhancedAIService, string, error) {
	p.loadBalancer.mutex.RLock()
	defer p.loadBalancer.mutex.RUnlock()

	var bestProvider string
	var bestScore float64

	for providerName, weight := range p.loadBalancer.weightedProviders {
		if weight > bestScore {
			// Check if provider is healthy
			if _, err := p.tryProvider(ctx, providerName); err == nil {
				bestProvider = providerName
				bestScore = weight
			}
		}
	}

	if bestProvider == "" {
		return nil, "", fmt.Errorf("no healthy providers available")
	}

	client, err := p.tryProvider(ctx, bestProvider)
	return client, bestProvider, err
}

// checkCostLimits verifies if provider is within cost constraints
func (p *AIClientProvider) checkCostLimits(providerName string) bool {
	// Check if cost tracking needs reset (daily)
	if time.Since(p.costOptimizer.lastReset) > time.Hour*24 {
		p.costOptimizer.costTracking = make(map[string]float64)
		p.costOptimizer.lastReset = time.Now()
	}

	// Check provider-specific cost threshold
	if threshold, exists := p.costOptimizer.costThresholds[providerName]; exists {
		if currentCost := p.costOptimizer.costTracking[providerName]; currentCost > threshold {
			return false
		}
	}

	// Check daily budget limit
	totalDailyCost := 0.0
	for _, cost := range p.costOptimizer.costTracking {
		totalDailyCost += cost
	}

	if dailyLimit := p.costOptimizer.budgetLimits["daily"]; totalDailyCost > dailyLimit {
		return false
	}

	return true
}

// RecordCost tracks cost for a provider
func (p *AIClientProvider) RecordCost(providerName string, cost float64) {
	p.costOptimizer.costTracking[providerName] += cost

	// Update provider metrics
	p.mutex.Lock()
	if metrics, exists := p.metrics[providerName]; exists {
		metrics.CostAccumulated += cost
	}
	p.mutex.Unlock()
}

// RecordQualityRating records a quality rating for a provider's response
func (p *AIClientProvider) RecordQualityRating(providerName string, rating float64) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if metrics, exists := p.metrics[providerName]; exists {
		metrics.QualityRatings = append(metrics.QualityRatings, rating)

		// Keep only last 100 ratings
		if len(metrics.QualityRatings) > 100 {
			metrics.QualityRatings = metrics.QualityRatings[len(metrics.QualityRatings)-100:]
		}
	}
}

// GetPerformanceScores returns current performance scores for all providers
func (p *AIClientProvider) GetPerformanceScores() map[string]*PerformanceScore {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	scores := make(map[string]*PerformanceScore)
	for name, score := range p.performanceScores {
		// Create a copy to avoid race conditions
		scoreCopy := *score
		scores[name] = &scoreCopy
	}

	return scores
}

// GetHealthStatus returns current health status for all providers
func (p *AIClientProvider) GetHealthStatus() map[string]*HealthCheck {
	p.healthMonitor.mutex.RLock()
	defer p.healthMonitor.mutex.RUnlock()

	status := make(map[string]*HealthCheck)
	for name, check := range p.healthMonitor.healthChecks {
		// Create a copy to avoid race conditions
		checkCopy := *check
		status[name] = &checkCopy
	}

	return status
}

// UpdateRequestTypeRouting allows dynamic updating of request type preferences
func (p *AIClientProvider) UpdateRequestTypeRouting(requestType string, providers []string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.requestTypeRouting[requestType] = providers
	log.Printf("INFO: Updated routing for request type '%s' to providers: %v", requestType, providers)
}

// === COMMUNICATION OPTIMIZATION METHODS ===

// GetCachedClient returns a cached client if available and healthy
func (p *AIClientProvider) GetCachedClient(ctx context.Context, provider string) (ai.EnhancedAIService, bool) {
	p.commOptMutex.RLock()
	cached, exists := p.clientCache[provider]
	p.commOptMutex.RUnlock()

	if !exists {
		return nil, false
	}

	// Check if cache entry is still valid (not expired)
	if time.Since(cached.LastUsed) > 5*time.Minute || !cached.IsHealthy {
		p.commOptMutex.Lock()
		delete(p.clientCache, provider)
		p.commOptMutex.Unlock()
		return nil, false
	}

	// Update usage stats
	atomic.AddInt64(&cached.UseCount, 1)
	cached.LastUsed = time.Now()

	log.Printf("DEBUG: CommunicationOptimizer: Cache hit for provider '%s' (uses: %d)", provider, cached.UseCount)
	return cached.Client, true
}

// CacheClient stores a client in the cache with optimization metadata
func (p *AIClientProvider) CacheClient(provider string, client ai.EnhancedAIService) {
	p.commOptMutex.Lock()
	defer p.commOptMutex.Unlock()

	p.clientCache[provider] = &CachedClientInfo{
		Client:          client,
		LastUsed:        time.Now(),
		UseCount:        1,
		IsHealthy:       true,
		LastHealthCheck: time.Now(),
		CreatedAt:       time.Now(),
	}

	log.Printf("DEBUG: CommunicationOptimizer: Cached client for provider '%s'", provider)
}

// StartRequestTracking begins tracking a request for performance analysis
func (p *AIClientProvider) StartRequestTracking(requestID, provider, requestType string, toolsRequested int) {
	metrics := RequestMetrics{
		RequestID:      requestID,
		StartTime:      time.Now(),
		Provider:       provider,
		ToolsRequested: toolsRequested,
		RequestType:    requestType,
	}

	p.requestTracker.activeRequests.Store(requestID, metrics)
	log.Printf("DEBUG: RequestTracker: Started tracking request %s for provider %s", requestID, provider)
}

// CompleteRequestTracking finishes tracking a request and updates performance metrics
func (p *AIClientProvider) CompleteRequestTracking(requestID string, success bool, toolsExecuted, responseLength int, cacheHit bool) {
	if activeMetrics, exists := p.requestTracker.activeRequests.LoadAndDelete(requestID); exists {
		metrics := activeMetrics.(RequestMetrics)
		metrics.EndTime = time.Now()
		metrics.Latency = metrics.EndTime.Sub(metrics.StartTime)
		metrics.Success = success
		metrics.ToolsExecuted = toolsExecuted
		metrics.ResponseLength = responseLength
		metrics.CacheHit = cacheHit

		// Add to history
		p.requestTracker.historyMutex.Lock()
		p.requestTracker.requestHistory = append(p.requestTracker.requestHistory, metrics)

		// Keep only last 1000 requests
		if len(p.requestTracker.requestHistory) > 1000 {
			p.requestTracker.requestHistory = p.requestTracker.requestHistory[100:]
		}

		// Check for performance degradation
		if metrics.Latency > 10*time.Second {
			atomic.AddInt64(&p.requestTracker.consecutiveSlowCount, 1)
			p.requestTracker.lastSlowRequestTime = time.Now()

			if p.requestTracker.consecutiveSlowCount > 3 {
				p.requestTracker.performanceDegradation = true
				log.Printf("WARN: RequestTracker: Performance degradation detected for provider %s", metrics.Provider)
			}
		} else {
			atomic.StoreInt64(&p.requestTracker.consecutiveSlowCount, 0)
			p.requestTracker.performanceDegradation = false
		}

		p.requestTracker.historyMutex.Unlock()

		// Update adaptive timeout
		p.UpdateAdaptiveTimeout(metrics.Provider, metrics.Latency, success)

		log.Printf("DEBUG: RequestTracker: Completed tracking request %s - Latency: %v, Success: %t",
			requestID, metrics.Latency, success)
	}
}

// UpdateAdaptiveTimeout adjusts timeout based on recent performance
func (p *AIClientProvider) UpdateAdaptiveTimeout(provider string, latency time.Duration, success bool) {
	p.commOptMutex.Lock()
	defer p.commOptMutex.Unlock()

	timeout, exists := p.adaptiveTimeouts[provider]
	if !exists {
		return
	}

	timeout.mutex.Lock()
	defer timeout.mutex.Unlock()

	// Add to recent latencies
	timeout.recentLatencies = append(timeout.recentLatencies, latency)
	if len(timeout.recentLatencies) > 10 {
		timeout.recentLatencies = timeout.recentLatencies[1:]
	}

	// Update success rate
	currentSuccessRate := timeout.successRate
	if success {
		timeout.successRate = 0.9*currentSuccessRate + 0.1*1.0
	} else {
		timeout.successRate = 0.9*currentSuccessRate + 0.1*0.0
	}

	// Adjust timeout if enough time has passed since last adjustment
	if time.Since(timeout.lastAdjustment) > time.Minute {
		// Calculate average latency
		if len(timeout.recentLatencies) > 0 {
			var totalLatency time.Duration
			for _, lat := range timeout.recentLatencies {
				totalLatency += lat
			}
			avgLatency := totalLatency / time.Duration(len(timeout.recentLatencies))

			// Adjust timeout based on performance
			if timeout.successRate > 0.95 && avgLatency < timeout.currentTimeout/2 {
				// Performing well, can reduce timeout
				newTimeout := timeout.currentTimeout - time.Second*5
				if newTimeout < timeout.baseTimeout/2 {
					newTimeout = timeout.baseTimeout / 2
				}
				timeout.currentTimeout = newTimeout
				log.Printf("INFO: AdaptiveTimeout: Reduced timeout for provider %s to %v", provider, newTimeout)
			} else if timeout.successRate < 0.8 || avgLatency > time.Duration(float64(timeout.currentTimeout)*0.8) {
				// Performing poorly, increase timeout
				newTimeout := timeout.currentTimeout + time.Second*10
				if newTimeout > timeout.baseTimeout*3 {
					newTimeout = timeout.baseTimeout * 3
				}
				timeout.currentTimeout = newTimeout
				log.Printf("INFO: AdaptiveTimeout: Increased timeout for provider %s to %v", provider, newTimeout)
			}

			timeout.lastAdjustment = time.Now()
		}
	}
}

// GetAdaptiveTimeout returns the current adaptive timeout for a provider
func (p *AIClientProvider) GetAdaptiveTimeout(provider string) time.Duration {
	p.commOptMutex.RLock()
	defer p.commOptMutex.RUnlock()

	if timeout, exists := p.adaptiveTimeouts[provider]; exists {
		timeout.mutex.RLock()
		defer timeout.mutex.RUnlock()
		return timeout.currentTimeout
	}

	return 30 * time.Second // Default timeout
}

// GetRequestHistory returns recent request history for analysis
func (p *AIClientProvider) GetRequestHistory(limit int) []RequestMetrics {
	p.requestTracker.historyMutex.RLock()
	defer p.requestTracker.historyMutex.RUnlock()

	history := p.requestTracker.requestHistory
	if limit > 0 && len(history) > limit {
		history = history[len(history)-limit:]
	}

	// Return copy to avoid race conditions
	result := make([]RequestMetrics, len(history))
	copy(result, history)
	return result
}

// IsPerformanceDegraded returns true if performance degradation is detected
func (p *AIClientProvider) IsPerformanceDegraded() bool {
	p.requestTracker.historyMutex.RLock()
	defer p.requestTracker.historyMutex.RUnlock()
	return p.requestTracker.performanceDegradation
}

// GetCommunicationStats returns detailed communication optimization statistics
func (p *AIClientProvider) GetCommunicationStats() map[string]interface{} {
	p.commOptMutex.RLock()
	defer p.commOptMutex.RUnlock()

	stats := make(map[string]interface{})

	// Cache statistics
	cacheStats := make(map[string]interface{})
	totalCacheHits := int64(0)
	totalCacheEntries := len(p.clientCache)

	for provider, cached := range p.clientCache {
		totalCacheHits += cached.UseCount
		cacheStats[provider] = map[string]interface{}{
			"use_count":         cached.UseCount,
			"last_used":         cached.LastUsed,
			"is_healthy":        cached.IsHealthy,
			"age":               time.Since(cached.CreatedAt),
			"last_health_check": cached.LastHealthCheck,
		}
	}

	stats["cache"] = map[string]interface{}{
		"total_entries": totalCacheEntries,
		"total_hits":    totalCacheHits,
		"providers":     cacheStats,
	}

	// Adaptive timeout statistics
	timeoutStats := make(map[string]interface{})
	for provider, timeout := range p.adaptiveTimeouts {
		timeout.mutex.RLock()
		timeoutStats[provider] = map[string]interface{}{
			"base_timeout":    timeout.baseTimeout,
			"current_timeout": timeout.currentTimeout,
			"success_rate":    timeout.successRate,
			"last_adjustment": timeout.lastAdjustment,
		}
		timeout.mutex.RUnlock()
	}

	stats["adaptive_timeouts"] = timeoutStats

	// Request tracking statistics
	p.requestTracker.historyMutex.RLock()
	recentRequests := len(p.requestTracker.requestHistory)
	avgLatency := time.Duration(0)
	successfulRequests := 0

	if recentRequests > 0 {
		totalLatency := time.Duration(0)
		for _, req := range p.requestTracker.requestHistory {
			totalLatency += req.Latency
			if req.Success {
				successfulRequests++
			}
		}
		avgLatency = totalLatency / time.Duration(recentRequests)
	}

	stats["request_tracking"] = map[string]interface{}{
		"recent_requests":      recentRequests,
		"avg_latency":          avgLatency,
		"success_rate":         float64(successfulRequests) / float64(recentRequests),
		"consecutive_slow":     p.requestTracker.consecutiveSlowCount,
		"performance_degraded": p.requestTracker.performanceDegradation,
		"last_slow_request":    p.requestTracker.lastSlowRequestTime,
	}
	p.requestTracker.historyMutex.RUnlock()

	return stats
}
