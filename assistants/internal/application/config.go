package application

import (
	"time"

	"middleman/assistants/internal/application/services"
	"middleman/assistants/internal/domain"
)

// Config contains configuration for the application
type Config struct {
	LLMProcessor      services.LLMProcessor
	SpeechProcessor   services.SpeechProcessor
	VisionProcessor   services.VisionProcessor
	DocumentProcessor services.DocumentProcessor
	DataProcessor     services.DataProcessor
	PromptProvider    domain.SystemPromptProvider
	ToolConfig        *ToolConfig
	StreamingConfig   *StreamingConfig
}

// ToolConfig configures tool execution behavior
type ToolConfig struct {
	MaxConcurrentTools   int
	ToolExecutionTimeout time.Duration
	EnableMetrics        bool
}

// StreamingConfig configures streaming behavior
type StreamingConfig struct {
	MaxConcurrentTools    int
	ToolExecutionTimeout  time.Duration
	StreamBufferSize      int
	EnableProgressUpdates bool
	ChunkSize             int
}

// DefaultToolConfig returns the default tool configuration
func DefaultToolConfig() *ToolConfig {
	return &ToolConfig{
		MaxConcurrentTools:   20,
		ToolExecutionTimeout: 10 * time.Minute,
		EnableMetrics:        true,
	}
}

// DefaultStreamingConfig returns the default streaming configuration
func DefaultStreamingConfig() *StreamingConfig {
	return &StreamingConfig{
		MaxConcurrentTools:    20,
		ToolExecutionTimeout:  10 * time.Minute,
		StreamBufferSize:      100,
		EnableProgressUpdates: true,
		ChunkSize:             1024,
	}
}

// ValidateConfig validates and sets defaults for the configuration
func ValidateConfig(config *Config) *Config {
	if config == nil {
		panic("config cannot be nil")
	}

	// Validate required processors
	if config.LLMProcessor == nil {
		panic("LLMProcessor is required")
	}

	// Set default configs if not provided
	if config.ToolConfig == nil {
		config.ToolConfig = DefaultToolConfig()
	}
	if config.StreamingConfig == nil {
		config.StreamingConfig = DefaultStreamingConfig()
	}

	// Validate tool config
	if config.ToolConfig.MaxConcurrentTools <= 0 {
		config.ToolConfig.MaxConcurrentTools = 20
	}
	if config.ToolConfig.ToolExecutionTimeout <= 0 {
		config.ToolConfig.ToolExecutionTimeout = 10 * time.Minute
	}

	// Validate streaming config
	if config.StreamingConfig.MaxConcurrentTools <= 0 {
		config.StreamingConfig.MaxConcurrentTools = 20
	}
	if config.StreamingConfig.ToolExecutionTimeout <= 0 {
		config.StreamingConfig.ToolExecutionTimeout = 10 * time.Minute
	}
	if config.StreamingConfig.StreamBufferSize <= 0 {
		config.StreamingConfig.StreamBufferSize = 100
	}
	if config.StreamingConfig.ChunkSize <= 0 {
		config.StreamingConfig.ChunkSize = 1024
	}

	return config
}