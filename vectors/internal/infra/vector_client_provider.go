package infra

import (
	"context"
	"fmt"
	"log"
	"middleman/vectors/internal/ports"
	"middleman/vectors/internal/vector"
	"os"
	"sync"
	"time"
)

// VectorClientProvider manages multiple vector database clients with circuit breakers
type VectorClientProvider struct {
	clients         map[string]ports.VectorDatabaseClient
	healthStatus    map[string]*VectorProviderHealth
	circuitBreakers map[string]*CircuitBreaker
	defaultName     string
	mu              sync.RWMutex
	healthTicker    *time.Ticker
	stopHealth      chan bool
}

// VectorProviderHealth tracks the health status of a vector database provider
type VectorProviderHealth struct {
	Provider         string
	Healthy          bool
	LastCheck        time.Time
	Error            string
	Latency          time.Duration
	SuccessRate      float64
	TotalRequests    int64
	FailedRequests   int64
	ConnectionStatus string
	LastOperation    time.Time
}

// CircuitBreaker implements a basic circuit breaker pattern
type CircuitBreaker struct {
	MaxFailures  int
	ResetTimeout time.Duration
	Timeout      time.Duration
	failures     int
	lastFailTime time.Time
	state        string // "closed", "open", "half-open"
	mu           sync.RWMutex
}

// NewVectorClientProvider creates a new vector client provider
func NewVectorClientProvider() *VectorClientProvider {
	provider := &VectorClientProvider{
		clients:         make(map[string]ports.VectorDatabaseClient),
		healthStatus:    make(map[string]*VectorProviderHealth),
		circuitBreakers: make(map[string]*CircuitBreaker),
		stopHealth:      make(chan bool, 1),
	}

	// Auto-discover and register Qdrant client if configuration is available
	if err := provider.autoDiscoverClients(); err != nil {
		log.Printf("Warning: Auto-discovery of vector clients failed: %v", err)
	}

	// Start health monitoring
	provider.startHealthMonitoring()

	return provider
}

// RegisterClient registers a vector database client with a given name
func (p *VectorClientProvider) RegisterClient(name string, client ports.VectorDatabaseClient) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.clients[name] = client
	p.healthStatus[name] = &VectorProviderHealth{
		Provider:         name,
		Healthy:          true,
		LastCheck:        time.Now(),
		SuccessRate:      1.0,
		ConnectionStatus: "connected",
	}

	// Create circuit breaker for this client
	p.circuitBreakers[name] = &CircuitBreaker{
		MaxFailures:  5,
		ResetTimeout: 60 * time.Second,
		Timeout:      30 * time.Second,
		state:        "closed",
	}

	if p.defaultName == "" {
		p.defaultName = name
	}

	log.Printf("Registered vector database client: %s", name)
	return nil
}

// GetClient returns a specific vector database client by name
func (p *VectorClientProvider) GetClient(ctx context.Context, provider string) (ports.VectorDatabaseClient, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	client, exists := p.clients[provider]
	if !exists {
		return nil, fmt.Errorf("vector database client not found: %s", provider)
	}

	// Check circuit breaker
	circuitBreaker := p.circuitBreakers[provider]
	if circuitBreaker != nil && circuitBreaker.IsOpen() {
		return nil, fmt.Errorf("vector database client %s is circuit-broken", provider)
	}

	// Check if client is healthy
	if health, exists := p.healthStatus[provider]; exists && !health.Healthy {
		return nil, fmt.Errorf("vector database client %s is unhealthy: %s", provider, health.Error)
	}

	return client, nil
}

// GetDefaultClient returns the default vector database client
func (p *VectorClientProvider) GetDefaultClient(ctx context.Context) (ports.VectorDatabaseClient, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.defaultName == "" {
		return nil, fmt.Errorf("no default vector database client configured")
	}

	client, exists := p.clients[p.defaultName]
	if !exists {
		return nil, fmt.Errorf("default vector database client not found: %s", p.defaultName)
	}

	return client, nil
}

// GetHealthyProvider returns the best available healthy vector database client
func (p *VectorClientProvider) GetHealthyProvider(ctx context.Context) (ports.VectorDatabaseClient, string, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var bestClient ports.VectorDatabaseClient
	var bestName string
	bestScore := -1.0

	for name, client := range p.clients {
		health, exists := p.healthStatus[name]
		if !exists || !health.Healthy {
			continue
		}

		circuitBreaker := p.circuitBreakers[name]
		if circuitBreaker != nil && circuitBreaker.IsOpen() {
			continue
		}

		// Calculate provider score based on health metrics
		score := p.calculateProviderScore(health)
		if score > bestScore {
			bestScore = score
			bestClient = client
			bestName = name
		}
	}

	if bestClient == nil {
		return nil, "", fmt.Errorf("no healthy vector database providers available")
	}

	return bestClient, bestName, nil
}

// GetOptimalProvider returns the best vector database client for the given operation type
func (p *VectorClientProvider) GetOptimalProvider(ctx context.Context, operationType string) (ports.VectorDatabaseClient, string, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var bestClient ports.VectorDatabaseClient
	var bestName string
	bestScore := -1.0

	for name, client := range p.clients {
		health, exists := p.healthStatus[name]
		if !exists || !health.Healthy {
			continue
		}

		circuitBreaker := p.circuitBreakers[name]
		if circuitBreaker != nil && circuitBreaker.IsOpen() {
			continue
		}

		// Calculate provider score with operation-specific optimizations
		score := p.calculateProviderScoreForOperation(health, operationType)
		if score > bestScore {
			bestScore = score
			bestClient = client
			bestName = name
		}
	}

	if bestClient == nil {
		return nil, "", fmt.Errorf("no healthy vector database providers available for operation: %s", operationType)
	}

	return bestClient, bestName, nil
}

// ExecuteWithCircuitBreaker executes an operation with circuit breaker protection
func (p *VectorClientProvider) ExecuteWithCircuitBreaker(ctx context.Context, provider string, operation func() error) error {
	p.mu.RLock()
	circuitBreaker := p.circuitBreakers[provider]
	health := p.healthStatus[provider]
	p.mu.RUnlock()

	if circuitBreaker == nil {
		return operation()
	}

	if circuitBreaker.IsOpen() {
		return fmt.Errorf("circuit breaker is open for provider: %s", provider)
	}

	startTime := time.Now()
	err := operation()
	duration := time.Since(startTime)

	p.mu.Lock()
	if health != nil {
		health.LastOperation = time.Now()
		health.TotalRequests++

		if err != nil {
			health.FailedRequests++
			circuitBreaker.RecordFailure()
		} else {
			circuitBreaker.RecordSuccess()
		}

		health.SuccessRate = float64(health.TotalRequests-health.FailedRequests) / float64(health.TotalRequests)
		health.Latency = duration
	}
	p.mu.Unlock()

	return err
}

// GetProviderHealth returns health status for all providers
func (p *VectorClientProvider) GetProviderHealth() map[string]*VectorProviderHealth {
	p.mu.RLock()
	defer p.mu.RUnlock()

	health := make(map[string]*VectorProviderHealth)
	for name, status := range p.healthStatus {
		health[name] = &VectorProviderHealth{
			Provider:         status.Provider,
			Healthy:          status.Healthy,
			LastCheck:        status.LastCheck,
			Error:            status.Error,
			Latency:          status.Latency,
			SuccessRate:      status.SuccessRate,
			TotalRequests:    status.TotalRequests,
			FailedRequests:   status.FailedRequests,
			ConnectionStatus: status.ConnectionStatus,
			LastOperation:    status.LastOperation,
		}
	}
	return health
}

// Close stops health monitoring and cleans up resources
func (p *VectorClientProvider) Close() error {
	if p.healthTicker != nil {
		p.healthTicker.Stop()
	}
	select {
	case p.stopHealth <- true:
	default:
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	for name, client := range p.clients {
		if err := client.Close(); err != nil {
			log.Printf("Error closing vector client %s: %v", name, err)
		}
	}

	return nil
}

// autoDiscoverClients automatically discovers and registers available vector database clients
func (p *VectorClientProvider) autoDiscoverClients() error {
	// Try to create Qdrant client if configuration is available
	qdrantHost := os.Getenv("QDRANT_HOST")
	qdrantPort := os.Getenv("QDRANT_PORT")

	if qdrantHost == "" {
		qdrantHost = "localhost"
	}
	if qdrantPort == "" {
		qdrantPort = "6334"
	}

	config := vector.Config{
		QdrantHost:     qdrantHost,
		QdrantPort:     qdrantPort,
		CollectionName: "vectors",
		VectorSize:     1536, // OpenAI text-embedding-3-small dimensions
	}

	vectorService, err := vector.NewVectorService(config)
	if err != nil {
		log.Printf("Failed to create Qdrant vector service: %v", err)
		return err
	}

	vectorClient := NewQdrantVectorClient(vectorService)
	if err := p.RegisterClient("qdrant", vectorClient); err != nil {
		log.Printf("Failed to register Qdrant vector client: %v", err)
	} else {
		log.Printf("Auto-registered Qdrant vector client at %s:%s", qdrantHost, qdrantPort)
	}

	return nil
}

// startHealthMonitoring starts background health monitoring
func (p *VectorClientProvider) startHealthMonitoring() {
	p.healthTicker = time.NewTicker(2 * time.Minute) // Check every 2 minutes

	go func() {
		for {
			select {
			case <-p.healthTicker.C:
				p.performHealthChecks()
			case <-p.stopHealth:
				return
			}
		}
	}()
}

// performHealthChecks checks the health of all registered clients
func (p *VectorClientProvider) performHealthChecks() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for name, client := range p.clients {
		startTime := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

		err := client.HealthCheck(ctx)
		latency := time.Since(startTime)
		cancel()

		health := p.healthStatus[name]
		health.LastCheck = time.Now()
		health.Latency = latency

		if err != nil {
			health.Healthy = false
			health.Error = err.Error()
			health.ConnectionStatus = "disconnected"
			log.Printf("Health check failed for vector provider %s: %v", name, err)
		} else {
			health.Healthy = true
			health.Error = ""
			health.ConnectionStatus = "connected"
		}
	}
}

// calculateProviderScore calculates a score for provider selection
func (p *VectorClientProvider) calculateProviderScore(health *VectorProviderHealth) float64 {
	if !health.Healthy {
		return -1
	}

	// Base score from success rate
	score := health.SuccessRate * 100

	// Adjust for latency (lower is better)
	latencyScore := 100.0 - (float64(health.Latency.Milliseconds()) / 10.0)
	if latencyScore < 0 {
		latencyScore = 0
	}
	score += latencyScore * 0.4

	// Bonus for recent activity
	if time.Since(health.LastOperation) < 5*time.Minute {
		score += 5
	}

	return score
}

// calculateProviderScoreForOperation calculates a score for specific operations
func (p *VectorClientProvider) calculateProviderScoreForOperation(health *VectorProviderHealth, operationType string) float64 {
	baseScore := p.calculateProviderScore(health)

	if baseScore < 0 {
		return baseScore
	}

	// Operation-specific optimizations
	switch operationType {
	case "search", "query":
		// For search operations, prioritize low latency
		baseScore += 15
	case "index", "batch_index":
		// For indexing operations, prioritize stability
		baseScore += 10
	case "delete", "batch_delete":
		// For delete operations, prioritize reliability
		baseScore += 5
	}

	return baseScore
}

// Circuit Breaker implementation

func NewCircuitBreaker(maxFailures int, resetTimeout, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		MaxFailures:  maxFailures,
		ResetTimeout: resetTimeout,
		Timeout:      timeout,
		state:        "closed",
	}
}

func (cb *CircuitBreaker) IsOpen() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	if cb.state == "open" {
		if time.Since(cb.lastFailTime) > cb.ResetTimeout {
			cb.state = "half-open"
			return false
		}
		return true
	}
	return false
}

func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures = 0
	if cb.state == "half-open" {
		cb.state = "closed"
	}
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures++
	cb.lastFailTime = time.Now()

	if cb.failures >= cb.MaxFailures {
		cb.state = "open"
	}
}
