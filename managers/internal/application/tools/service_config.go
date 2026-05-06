package tools

import "time"

// ServiceConfig holds configuration settings for tool services
type ServiceConfig struct {
	// Streaming configuration
	EnableStreaming  bool          `json:"enable_streaming"`
	StreamBufferSize int           `json:"stream_buffer_size"`
	ProgressInterval time.Duration `json:"progress_interval"`

	// Timeout configuration
	OperationTimeout  time.Duration `json:"operation_timeout"`
	ConnectionTimeout time.Duration `json:"connection_timeout"`

	// Retry configuration
	MaxRetries int           `json:"max_retries"`
	RetryDelay time.Duration `json:"retry_delay"`

	// Performance configuration
	MaxConcurrentOps int `json:"max_concurrent_ops"`
	BatchSize        int `json:"batch_size"`

	// Logging configuration
	EnableMetrics   bool   `json:"enable_metrics"`
	EnableDebugLogs bool   `json:"enable_debug_logs"`
	LogLevel        string `json:"log_level"`
}

// DefaultServiceConfig returns a default service configuration
func DefaultServiceConfig() *ServiceConfig {
	return &ServiceConfig{
		// Streaming defaults
		EnableStreaming:  true,
		StreamBufferSize: 100,
		ProgressInterval: 100 * time.Millisecond,

		// Timeout defaults
		OperationTimeout:  10 * time.Minute,
		ConnectionTimeout: 10 * time.Minute,

		// Retry defaults
		MaxRetries: 3,
		RetryDelay: 500 * time.Millisecond,

		// Performance defaults
		MaxConcurrentOps: 10,
		BatchSize:        20,

		// Logging defaults
		EnableMetrics:   true,
		EnableDebugLogs: false,
		LogLevel:        "info",
	}
}
