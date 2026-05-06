package infra

import (
	"context"
	"fmt"
	"log"
	ai2 "middleman/internal/ai"
	"middleman/vectors/internal/ports"
	"os"
	"sync"
	"time"
)

// EmbeddingClientProvider manages multiple embedding clients with health monitoring
type EmbeddingClientProvider struct {
	clients      map[string]ports.EmbeddingClient
	healthStatus map[string]*ProviderHealth
	defaultName  string
	mu           sync.RWMutex
	healthTicker *time.Ticker
	stopHealth   chan bool
}

// ProviderHealth tracks the health status of an embedding provider
type ProviderHealth struct {
	Provider       string
	Healthy        bool
	LastCheck      time.Time
	Error          string
	Latency        time.Duration
	SuccessRate    float64
	TotalRequests  int64
	FailedRequests int64
}

// NewEmbeddingClientProvider creates a new embedding client provider
func NewEmbeddingClientProvider() *EmbeddingClientProvider {
	provider := &EmbeddingClientProvider{
		clients:      make(map[string]ports.EmbeddingClient),
		healthStatus: make(map[string]*ProviderHealth),
		stopHealth:   make(chan bool, 1),
	}

	// Auto-discover and register OpenAI client if API key is available
	if err := provider.autoDiscoverClients(); err != nil {
		log.Printf("Warning: Auto-discovery of embedding clients failed: %v", err)
	}

	// Start health monitoring
	provider.startHealthMonitoring()

	return provider
}

// RegisterClient registers an embedding client with a given name
func (p *EmbeddingClientProvider) RegisterClient(name string, client ports.EmbeddingClient) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.clients[name] = client
	p.healthStatus[name] = &ProviderHealth{
		Provider:    name,
		Healthy:     true,
		LastCheck:   time.Now(),
		SuccessRate: 1.0,
	}

	if p.defaultName == "" {
		p.defaultName = name
	}

	log.Printf("Registered embedding client: %s", name)
	return nil
}

// GetClient returns a specific embedding client by name
func (p *EmbeddingClientProvider) GetClient(ctx context.Context, provider string) (ports.EmbeddingClient, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	client, exists := p.clients[provider]
	if !exists {
		return nil, fmt.Errorf("embedding client not found: %s", provider)
	}

	// Check if client is healthy
	if health, exists := p.healthStatus[provider]; exists && !health.Healthy {
		return nil, fmt.Errorf("embedding client %s is unhealthy: %s", provider, health.Error)
	}

	return client, nil
}

// GetDefaultClient returns the default embedding client
func (p *EmbeddingClientProvider) GetDefaultClient(ctx context.Context) (ports.EmbeddingClient, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.defaultName == "" {
		return nil, fmt.Errorf("no default embedding client configured")
	}

	client, exists := p.clients[p.defaultName]
	if !exists {
		return nil, fmt.Errorf("default embedding client not found: %s", p.defaultName)
	}

	return client, nil
}

// GetOptimalProvider returns the best available embedding client for the given entity type
func (p *EmbeddingClientProvider) GetOptimalProvider(ctx context.Context, entityType string) (ports.EmbeddingClient, string, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var bestClient ports.EmbeddingClient
	var bestName string
	bestScore := -1.0

	for name, client := range p.clients {
		health, exists := p.healthStatus[name]
		if !exists || !health.Healthy {
			continue
		}

		// Calculate provider score based on health metrics
		score := p.calculateProviderScore(health, entityType)
		if score > bestScore {
			bestScore = score
			bestClient = client
			bestName = name
		}
	}

	if bestClient == nil {
		return nil, "", fmt.Errorf("no healthy embedding providers available")
	}

	return bestClient, bestName, nil
}

// GetProviderHealth returns health status for all providers
func (p *EmbeddingClientProvider) GetProviderHealth() map[string]*ProviderHealth {
	p.mu.RLock()
	defer p.mu.RUnlock()

	health := make(map[string]*ProviderHealth)
	for name, status := range p.healthStatus {
		health[name] = &ProviderHealth{
			Provider:       status.Provider,
			Healthy:        status.Healthy,
			LastCheck:      status.LastCheck,
			Error:          status.Error,
			Latency:        status.Latency,
			SuccessRate:    status.SuccessRate,
			TotalRequests:  status.TotalRequests,
			FailedRequests: status.FailedRequests,
		}
	}
	return health
}

// Close stops health monitoring and cleans up resources
func (p *EmbeddingClientProvider) Close() error {
	if p.healthTicker != nil {
		p.healthTicker.Stop()
	}
	select {
	case p.stopHealth <- true:
	default:
	}
	return nil
}

// autoDiscoverClients automatically discovers and registers available embedding clients
func (p *EmbeddingClientProvider) autoDiscoverClients() error {
	// Try to create OpenAI client if API key is available
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey != "" {
		baseURL := os.Getenv("OPENAI_BASE_URL")
		model := os.Getenv("OPENAI_EMBEDDING_MODEL")
		if model == "" {
			model = "text-embedding-3-small"
		}

		openAIClient, err := ai2.NewOpenAIClient(apiKey, baseURL, model)
		if err != nil {
			log.Printf("Failed to create OpenAI client: %v", err)
		} else {
			embeddingClient := NewOpenAIEmbeddingClient(openAIClient, EmbeddingClientConfig{
				Model:         model,
				Dimensions:    1536,
				PromptEnabled: true,
			})

			if err := p.RegisterClient("openai", embeddingClient); err != nil {
				log.Printf("Failed to register OpenAI embedding client: %v", err)
			} else {
				log.Printf("Auto-registered OpenAI embedding client with model: %s", model)
			}
		}
	}

	// Could add other providers here (Anthropic, DeepSeek, etc.)
	return nil
}

// startHealthMonitoring starts background health monitoring
func (p *EmbeddingClientProvider) startHealthMonitoring() {
	p.healthTicker = time.NewTicker(5 * time.Minute) // Check every 5 minutes

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
func (p *EmbeddingClientProvider) performHealthChecks() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for name, client := range p.clients {
		startTime := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		err := client.HealthCheck(ctx)
		latency := time.Since(startTime)
		cancel()

		health := p.healthStatus[name]
		health.LastCheck = time.Now()
		health.Latency = latency

		if err != nil {
			health.Healthy = false
			health.Error = err.Error()
			health.FailedRequests++
			log.Printf("Health check failed for embedding provider %s: %v", name, err)
		} else {
			health.Healthy = true
			health.Error = ""
		}

		health.TotalRequests++
		health.SuccessRate = float64(health.TotalRequests-health.FailedRequests) / float64(health.TotalRequests)
	}
}

// calculateProviderScore calculates a score for provider selection
func (p *EmbeddingClientProvider) calculateProviderScore(health *ProviderHealth, entityType string) float64 {
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
	score += latencyScore * 0.3

	// Entity-type specific optimizations
	switch entityType {
	case "product", "post":
		// For content-heavy entities, prefer higher-dimension models
		score += 10
	case "user":
		// For user entities, prefer faster models
		score += 5
	}

	return score
}
