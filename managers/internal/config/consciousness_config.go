package config

import (
	"os"
	"strconv"
	"time"
)

// ConsciousnessConfig holds configuration for the consciousness system
type ConsciousnessConfig struct {
	// Core Settings
	Enabled               bool
	DecisionDelayMS       int
	ConfidenceThreshold   float64
	MaxActionsPerMinute   int
	LearningEnabled       bool
	
	// AI Provider Settings
	DefaultProvider     string
	OpenAIAPIKey        string
	AnthropicAPIKey     string
	DeepSeekAPIKey      string
	FallbackEnabled     bool
	
	// AI Models
	OpenAIModel         string
	AnthropicModel      string
	DeepSeekModel       string
	
	// Tool Execution
	ToolTimeout         time.Duration
	ToolRetryCount      int
	MaxConcurrentTools  int
	
	// Rules Configuration
	RulesConfigPath     string
	RulesReloadInterval time.Duration
	
	// Performance
	MaxEventProcessingTime time.Duration
	EventBatchSize         int
	EventBufferSize        int
	
	// Feature Flags
	AIDecisionsEnabled      bool
	AutonomousActionsEnabled bool
	LearningModeEnabled     bool
	DryRunMode              bool
	
	// Rate Limiting
	RateLimitEnabled        bool
	RateLimitPerMinute      int
	
	// Circuit Breaker
	CircuitBreakerEnabled    bool
	CircuitBreakerMaxFailures int
	CircuitBreakerResetTimeout time.Duration
	CircuitBreakerSuccessThreshold int
}

// LoadConsciousnessConfig loads configuration from environment variables
func LoadConsciousnessConfig() *ConsciousnessConfig {
	return &ConsciousnessConfig{
		// Core Settings
		Enabled:               getEnvBool("MANAGER_CONSCIOUSNESS_ENABLED", true),
		DecisionDelayMS:       getEnvInt("MANAGER_DECISION_DELAY_MS", 1000),
		ConfidenceThreshold:   getEnvFloat("MANAGER_CONFIDENCE_THRESHOLD", 0.8),
		MaxActionsPerMinute:   getEnvInt("MANAGER_MAX_ACTIONS_PER_MINUTE", 10),
		LearningEnabled:       getEnvBool("MANAGER_LEARNING_ENABLED", true),
		
		// AI Provider Settings
		DefaultProvider:     getEnvString("AI_PROVIDER_DEFAULT", "deepseek"),
		OpenAIAPIKey:        getEnvString("AI_PROVIDER_OPENAI_API_KEY", ""),
		AnthropicAPIKey:     getEnvString("AI_PROVIDER_ANTHROPIC_API_KEY", ""),
		DeepSeekAPIKey:      getEnvString("AI_PROVIDER_DEEPSEEK_API_KEY", ""),
		FallbackEnabled:     getEnvBool("AI_PROVIDER_FALLBACK_ENABLED", true),
		
		// AI Models
		OpenAIModel:         getEnvString("AI_MODEL_OPENAI", "gpt-4o-mini"),
		AnthropicModel:      getEnvString("AI_MODEL_ANTHROPIC", "claude-3-5-sonnet-latest"),
		DeepSeekModel:       getEnvString("AI_MODEL_DEEPSEEK", "deepseek-chat"),
		
		// Tool Execution
		ToolTimeout:         getEnvDuration("MANAGER_TOOL_TIMEOUT", 30*time.Second),
		ToolRetryCount:      getEnvInt("MANAGER_TOOL_RETRY_COUNT", 3),
		MaxConcurrentTools:  getEnvInt("MANAGER_TOOL_MAX_CONCURRENT", 10),
		
		// Rules Configuration
		RulesConfigPath:     getEnvString("RULES_CONFIG_PATH", "/etc/managers/rules.yaml"),
		RulesReloadInterval: getEnvDuration("RULES_RELOAD_INTERVAL", 5*time.Minute),
		
		// Performance
		MaxEventProcessingTime: getEnvDuration("MAX_EVENT_PROCESSING_TIME", 5*time.Second),
		EventBatchSize:         getEnvInt("EVENT_BATCH_SIZE", 100),
		EventBufferSize:        getEnvInt("EVENT_BUFFER_SIZE", 1000),
		
		// Feature Flags
		AIDecisionsEnabled:      getEnvBool("FEATURE_AI_DECISIONS", true),
		AutonomousActionsEnabled: getEnvBool("FEATURE_AUTONOMOUS_ACTIONS", true),
		LearningModeEnabled:     getEnvBool("FEATURE_LEARNING_MODE", true),
		DryRunMode:              getEnvBool("FEATURE_DRY_RUN_MODE", false),
		
		// Rate Limiting
		RateLimitEnabled:   getEnvBool("RATE_LIMIT_ENABLED", true),
		RateLimitPerMinute: getEnvInt("RATE_LIMIT_REQUESTS_PER_MINUTE", 100),
		
		// Circuit Breaker
		CircuitBreakerEnabled:    getEnvBool("CIRCUIT_BREAKER_ENABLED", true),
		CircuitBreakerMaxFailures: getEnvInt("CIRCUIT_BREAKER_MAX_FAILURES", 5),
		CircuitBreakerResetTimeout: getEnvDuration("CIRCUIT_BREAKER_RESET_TIMEOUT", 2*time.Minute),
		CircuitBreakerSuccessThreshold: getEnvInt("CIRCUIT_BREAKER_SUCCESS_THRESHOLD", 3),
	}
}

// Helper functions
func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}