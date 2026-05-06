package ai

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ClientFactoryImpl implements the AIClientFactory interface
type ClientFactoryImpl struct {
	configs map[string]ProviderConfig
	mu      sync.RWMutex
}

// NewClientFactory creates a new AI client factory
func NewClientFactory() *ClientFactoryImpl {
	return &ClientFactoryImpl{
		configs: make(map[string]ProviderConfig),
	}
}

// CreateClient creates an AI client for the specified provider
func (f *ClientFactoryImpl) CreateClient(provider string, config ProviderConfig) (EnhancedAIService, error) {
	switch provider {
	case ProviderOpenAI:
		return NewOpenAIClient(config.APIKey, config.BaseURL, config.DefaultModel)
	case ProviderAnthropic:
		return NewAnthropicClient(config.APIKey, config.BaseURL, config.DefaultModel)
	case ProviderDeepSeek:
		return NewDeepSeekClient(config.APIKey, config.BaseURL, config.DefaultModel)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

// CreateMultiModalClient creates a multimodal AI client for the specified provider
func (f *ClientFactoryImpl) CreateMultiModalClient(provider string, config ProviderConfig) (MultiModalAIService, error) {
	client, err := f.CreateClient(provider, config)
	if err != nil {
		return nil, err
	}

	// Check if the client supports multimodal capabilities
	if multiModalClient, ok := client.(MultiModalAIService); ok {
		return multiModalClient, nil
	}

	return nil, fmt.Errorf("provider %s does not support multimodal capabilities", provider)
}

// SupportsMultiModal checks if a provider supports multimodal capabilities
func (f *ClientFactoryImpl) SupportsMultiModal(provider string) bool {
	switch provider {
	case ProviderOpenAI:
		return true // OpenAI supports vision, audio, and image generation
	case ProviderAnthropic:
		return false // Currently only supports text and vision (limited)
	case ProviderDeepSeek:
		return false // Currently only supports text
	default:
		return false
	}
}

// GetSupportedProviders returns the list of supported providers
func (f *ClientFactoryImpl) GetSupportedProviders() []string {
	return []string{ProviderOpenAI, ProviderAnthropic, ProviderDeepSeek}
}

// GetMultiModalProviders returns the list of providers that support multimodal capabilities
func (f *ClientFactoryImpl) GetMultiModalProviders() []string {
	providers := []string{}
	for _, provider := range f.GetSupportedProviders() {
		if f.SupportsMultiModal(provider) {
			providers = append(providers, provider)
		}
	}
	return providers
}

// RegisterProvider registers a provider configuration
func (f *ClientFactoryImpl) RegisterProvider(provider string, config ProviderConfig) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.configs[provider] = config
	return nil
}

// GetProviderConfig returns the configuration for a provider
func (f *ClientFactoryImpl) GetProviderConfig(provider string) (ProviderConfig, bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	config, exists := f.configs[provider]
	return config, exists
}

// ClientManagerImpl implements the AIClientManager interface
type ClientManagerImpl struct {
	factory           *ClientFactoryImpl
	clients           map[string]EnhancedAIService
	multiModalClients map[string]MultiModalAIService
	defaultProvider   string
	healthStatus      map[string]HealthStatus
	usageStats        map[string]UsageStats
	mu                sync.RWMutex
	healthCheckTick   *time.Ticker
	stopHealthCheck   chan struct{}
}

// NewClientManager creates a new AI client manager
func NewClientManager(factory *ClientFactoryImpl) *ClientManagerImpl {
	manager := &ClientManagerImpl{
		factory:           factory,
		clients:           make(map[string]EnhancedAIService),
		multiModalClients: make(map[string]MultiModalAIService),
		defaultProvider:   ProviderOpenAI,
		healthStatus:      make(map[string]HealthStatus),
		usageStats:        make(map[string]UsageStats),
		healthCheckTick:   time.NewTicker(5 * time.Minute),
		stopHealthCheck:   make(chan struct{}),
	}

	// Start health check routine
	go manager.healthCheckRoutine()

	return manager
}

// GetClient returns a client for the specified provider
func (m *ClientManagerImpl) GetClient(provider string) (EnhancedAIService, error) {
	m.mu.RLock()
	client, exists := m.clients[provider]
	m.mu.RUnlock()

	if exists {
		return client, nil
	}

	// Create new client
	config, exists := m.factory.GetProviderConfig(provider)
	if !exists {
		return nil, fmt.Errorf("provider %s not registered", provider)
	}

	client, err := m.factory.CreateClient(provider, config)
	if err != nil {
		return nil, err
	}

	m.mu.Lock()
	m.clients[provider] = client
	m.mu.Unlock()

	return client, nil
}

// GetMultiModalClient returns a multimodal client for the specified provider
func (m *ClientManagerImpl) GetMultiModalClient(provider string) (MultiModalAIService, error) {
	m.mu.RLock()
	client, exists := m.multiModalClients[provider]
	m.mu.RUnlock()

	if exists {
		return client, nil
	}

	// Create new multimodal client
	config, exists := m.factory.GetProviderConfig(provider)
	if !exists {
		return nil, fmt.Errorf("provider %s not registered", provider)
	}

	client, err := m.factory.CreateMultiModalClient(provider, config)
	if err != nil {
		return nil, err
	}

	m.mu.Lock()
	m.multiModalClients[provider] = client
	m.mu.Unlock()

	return client, nil
}

// GetDefaultMultiModalClient returns the default multimodal client if supported
func (m *ClientManagerImpl) GetDefaultMultiModalClient() (MultiModalAIService, error) {
	if !m.factory.SupportsMultiModal(m.defaultProvider) {
		// Try to find a provider that supports multimodal
		multiModalProviders := m.factory.GetMultiModalProviders()
		if len(multiModalProviders) == 0 {
			return nil, fmt.Errorf("no providers support multimodal capabilities")
		}
		return m.GetMultiModalClient(multiModalProviders[0])
	}
	return m.GetMultiModalClient(m.defaultProvider)
}

// GetDefaultClient returns the default AI client
func (m *ClientManagerImpl) GetDefaultClient() (EnhancedAIService, error) {
	return m.GetClient(m.defaultProvider)
}

// SetDefaultProvider sets the default provider
func (m *ClientManagerImpl) SetDefaultProvider(provider string) error {
	if !m.isProviderSupported(provider) {
		return fmt.Errorf("unsupported provider: %s", provider)
	}

	m.mu.Lock()
	m.defaultProvider = provider
	m.mu.Unlock()

	return nil
}

// CreateCompletionWithFallback creates a completion with automatic fallback
func (m *ClientManagerImpl) CreateCompletionWithFallback(ctx context.Context, request CompletionRequest) (*CompletionResponse, error) {
	providers := m.getActiveProviders()

	for _, provider := range providers {
		client, err := m.GetClient(provider)
		if err != nil {
			continue
		}

		response, err := client.CreateCompletion(ctx, request)
		if err == nil {
			return response, nil
		}

		// Log the error and try next provider
		m.recordError(provider, err)
	}

	return nil, fmt.Errorf("all providers failed")
}

// CreateCompletionWithBestProvider selects the best provider based on request requirements
func (m *ClientManagerImpl) CreateCompletionWithBestProvider(ctx context.Context, request CompletionRequest) (*CompletionResponse, error) {
	provider := m.selectBestProvider(request)

	client, err := m.GetClient(provider)
	if err != nil {
		return nil, err
	}

	return client.CreateCompletion(ctx, request)
}

// ListProviders returns all registered providers
func (m *ClientManagerImpl) ListProviders() []string {
	providers := m.factory.GetSupportedProviders()

	m.mu.RLock()
	defer m.mu.RUnlock()

	registered := make([]string, 0)
	for _, provider := range providers {
		if _, exists := m.factory.configs[provider]; exists {
			registered = append(registered, provider)
		}
	}

	return registered
}

// GetProviderHealth returns health status of all providers
func (m *ClientManagerImpl) GetProviderHealth() map[string]HealthStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	health := make(map[string]HealthStatus)
	for provider, status := range m.healthStatus {
		health[provider] = status
	}

	return health
}

// GetAllUsageStats returns usage statistics for all providers
func (m *ClientManagerImpl) GetAllUsageStats() map[string]UsageStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := make(map[string]UsageStats)
	for provider, client := range m.clients {
		stats[provider] = client.GetUsageStats()
	}

	return stats
}

// EnableProvider enables a provider
func (m *ClientManagerImpl) EnableProvider(provider string) error {
	config, exists := m.factory.GetProviderConfig(provider)
	if !exists {
		return fmt.Errorf("provider %s not found", provider)
	}

	config.Enabled = true
	return m.factory.RegisterProvider(provider, config)
}

// DisableProvider disables a provider
func (m *ClientManagerImpl) DisableProvider(provider string) error {
	config, exists := m.factory.GetProviderConfig(provider)
	if !exists {
		return fmt.Errorf("provider %s not found", provider)
	}

	config.Enabled = false
	return m.factory.RegisterProvider(provider, config)
}

// UpdateProviderConfig updates the configuration for a provider
func (m *ClientManagerImpl) UpdateProviderConfig(provider string, config ProviderConfig) error {
	// Remove existing client to force recreation with new config
	m.mu.Lock()
	delete(m.clients, provider)
	m.mu.Unlock()

	return m.factory.RegisterProvider(provider, config)
}

// Helper methods

func (m *ClientManagerImpl) isProviderSupported(provider string) bool {
	for _, p := range m.factory.GetSupportedProviders() {
		if p == provider {
			return true
		}
	}
	return false
}

func (m *ClientManagerImpl) getActiveProviders() []string {
	providers := make([]string, 0)

	for _, provider := range m.factory.GetSupportedProviders() {
		config, exists := m.factory.GetProviderConfig(provider)
		if exists && config.Enabled {
			providers = append(providers, provider)
		}
	}

	return providers
}

func (m *ClientManagerImpl) selectBestProvider(request CompletionRequest) string {
	// Simple selection logic - can be enhanced based on requirements
	// For now, prioritize based on model availability and cost

	if request.Model != "" {
		// Check which provider supports this model
		if isOpenAIModel(request.Model) {
			return ProviderOpenAI
		}
		if isAnthropicModel(request.Model) {
			return ProviderAnthropic
		}
		if isDeepSeekModel(request.Model) {
			return ProviderDeepSeek
		}
	}

	// Default to the cheapest for general requests
	return ProviderDeepSeek
}

func (m *ClientManagerImpl) recordError(provider string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	status, exists := m.healthStatus[provider]
	if !exists {
		status = HealthStatus{
			Provider: provider,
			Healthy:  false,
		}
	}

	status.Healthy = false
	status.LastCheck = time.Now()
	status.Error = err.Error()

	m.healthStatus[provider] = status
}

func (m *ClientManagerImpl) healthCheckRoutine() {
	for {
		select {
		case <-m.healthCheckTick.C:
			m.performHealthChecks()
		case <-m.stopHealthCheck:
			return
		}
	}
}

func (m *ClientManagerImpl) performHealthChecks() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, provider := range m.getActiveProviders() {
		client, err := m.GetClient(provider)
		if err != nil {
			m.recordError(provider, err)
			continue
		}

		start := time.Now()
		err = client.HealthCheck(ctx)
		latency := time.Since(start)

		m.mu.Lock()
		if err != nil {
			m.healthStatus[provider] = HealthStatus{
				Provider:  provider,
				Healthy:   false,
				LastCheck: time.Now(),
				Error:     err.Error(),
				Latency:   latency,
			}
		} else {
			m.healthStatus[provider] = HealthStatus{
				Provider:  provider,
				Healthy:   true,
				LastCheck: time.Now(),
				Latency:   latency,
			}
		}
		m.mu.Unlock()
	}
}

// Stop stops the health check routine
func (m *ClientManagerImpl) Stop() {
	close(m.stopHealthCheck)
	m.healthCheckTick.Stop()
}

// Helper functions to identify model providers
func isOpenAIModel(model string) bool {
	openAIModels := []string{
		ModelGPT4o, ModelGPT4oMini, ModelGPT4Turbo,
		ModelO1Preview, ModelO1Mini, ModelGPT4oLatest,
		ModelChatGPT4oLast,
	}

	for _, m := range openAIModels {
		if m == model {
			return true
		}
	}
	return false
}

func isAnthropicModel(model string) bool {
	anthropicModels := []string{
		ModelClaudeOpus4, ModelClaudeSonnet4,
		ModelClaude37Sonnet20250219, ModelClaude35Sonnet20241022,
		ModelClaude35Haiku20241022, ModelClaude3Opus20240229,
		ModelClaudeOpus4Latest, ModelClaudeSonnet4Latest,
		ModelClaude37SonnetLatest, ModelClaude35SonnetLatest,
		ModelClaude35HaikuLatest,
	}

	for _, m := range anthropicModels {
		if m == model {
			return true
		}
	}
	return false
}

func isDeepSeekModel(model string) bool {
	deepseekModels := []string{
		ModelDeepSeekV3, ModelDeepSeekV3_0324,
		ModelDeepSeekReasoner, ModelDeepSeekR1_0528,
		ModelDeepSeekCoder, ModelDeepSeekChat,
	}

	for _, m := range deepseekModels {
		if m == model {
			return true
		}
	}
	return false
}
