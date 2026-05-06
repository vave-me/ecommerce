package ai

import (
	"context"
	"io"
	"time"
)

// Provider constants
const (
	ProviderOpenAI    = "openai"
	ProviderAnthropic = "anthropic"
	ProviderDeepSeek  = "deepseek"
)

// Model constants for OpenAI (updated to latest 2025 models)
const (
	ModelGPT4o         = "gpt-4o"
	ModelGPT4oMini     = "gpt-4o-mini"
	ModelGPT41Mini     = "gpt-4.1-mini"  // New GPT-4.1 mini model
	ModelGPT4Turbo     = "gpt-4-turbo"
	ModelO1Preview     = "o1-preview"
	ModelO1Mini        = "o1-mini"
	ModelGPT4oLatest   = "gpt-4o-latest"
	ModelChatGPT4oLast = "chatgpt-4o-latest"
)

// Model constants for Anthropic (updated to latest 2025 models)
const (
	ModelClaudeOpus4            = "claude-opus-4-20250514"
	ModelClaudeSonnet4          = "claude-sonnet-4-20250514"
	ModelClaude37Sonnet20250219 = "claude-3-7-sonnet-20250219"
	ModelClaude35Sonnet20241022 = "claude-3-5-sonnet-20241022"
	ModelClaude35Haiku20241022  = "claude-3-5-haiku-20241022"
	ModelClaude3Opus20240229    = "claude-3-opus-20240229"
	ModelClaudeOpus4Latest      = "claude-opus-4-0"
	ModelClaudeSonnet4Latest    = "claude-sonnet-4-0"
	ModelClaude37SonnetLatest   = "claude-3-7-sonnet-latest"
	ModelClaude35SonnetLatest   = "claude-3-5-sonnet-latest"
	ModelClaude35HaikuLatest    = "claude-3-5-haiku-latest"
)

// Model constants for DeepSeek (updated to latest 2025 models)
const (
	ModelDeepSeekV3       = "deepseek-chat"     // Points to DeepSeek-V3-0324
	ModelDeepSeekV3_0324  = "deepseek-v3-0324"  // Specific version
	ModelDeepSeekReasoner = "deepseek-reasoner" // Points to DeepSeek-R1-0528
	ModelDeepSeekR1_0528  = "deepseek-r1-0528"  // Specific version
	ModelDeepSeekCoder    = "deepseek-coder"    // Legacy for compatibility
	ModelDeepSeekChat     = "deepseek-chat"     // Alias for V3
)

// Message roles
const (
	RoleSystem    = "system"
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleTool      = "tool"
	RoleFunction  = "function"
)

// Tool types
const (
	ToolTypeFunction = "function"
	ToolTypeCustom   = "custom"
)

// Finish reasons
const (
	FinishReasonStop          = "stop"
	FinishReasonEndTurn       = "end_turn"
	FinishReasonLength        = "length"
	FinishReasonMaxTokens     = "max_tokens"
	FinishReasonToolCalls     = "tool_calls"
	FinishReasonToolUse       = "tool_use"
	FinishReasonContentFilter = "content_filter"
	FinishReasonRefusal       = "refusal"
	FinishReasonStopSequence  = "stop_sequence"
)

// ==================== MULTIMODAL DATA STRUCTURES ====================
// Note: Multimodal data structures are defined in individual client files
// to avoid circular dependencies and allow for provider-specific implementations

// ==================== EXISTING DATA STRUCTURES ====================

// Message represents a chat message
type Message struct {
	Role       string      `json:"role"`
	Content    interface{} `json:"content"` // Can be string or structured content
	Name       string      `json:"name,omitempty"`
	ToolCalls  []ToolCall  `json:"tool_calls,omitempty"`
	ToolCallID string      `json:"tool_call_id,omitempty"`
	Thinking   string      `json:"thinking,omitempty"` // For Claude extended thinking
}

// GetContentAsString safely extracts content as string with type assertion
func (m *Message) GetContentAsString() string {
	if m.Content == nil {
		return ""
	}

	switch content := m.Content.(type) {
	case string:
		return content
	case []interface{}:
		// Handle structured content (like for vision)
		for _, item := range content {
			if itemMap, ok := item.(map[string]interface{}); ok {
				if text, exists := itemMap["text"]; exists {
					if textStr, ok := text.(string); ok {
						return textStr
					}
				}
			}
		}
		return ""
	default:
		return ""
	}
}

// ToolCall represents a function call
type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

// FunctionCall represents a function call
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// Tool represents a tool definition
type Tool struct {
	Type     string      `json:"type"`
	Function FunctionDef `json:"function"`
}

// FunctionDef represents a function definition
type FunctionDef struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// ToolDefinition represents a tool that can be called (for compatibility)
type ToolDefinition struct {
	Type     string      `json:"type"`
	Function FunctionDef `json:"function"`
}

// JSONSchemaDefinition represents a JSON schema (updated to match main codebase)
type JSONSchemaDefinition struct {
	Name        string                 `json:"name"`        // Schema name
	Description string                 `json:"description"` // Schema description
	Schema      map[string]interface{} `json:"schema"`      // JSON schema
	Strict      bool                   `json:"strict"`      // Whether to enforce strict adherence
}

// ToMap converts JSONSchemaDefinition to map[string]interface{} for compatibility
func (j *JSONSchemaDefinition) ToMap() map[string]interface{} {
	if j == nil {
		return nil
	}

	if j.Schema != nil {
		return j.Schema
	}

	return map[string]interface{}{
		"type": "object",
	}
}

// ResponseFormat represents the response format for structured output
type ResponseFormat struct {
	Type       string                 `json:"type"`
	JSONSchema map[string]interface{} `json:"json_schema,omitempty"`
	Schema     map[string]interface{} `json:"schema,omitempty"` // Alternative field name
}

// CompletionRequest represents a completion request
type CompletionRequest struct {
	Messages         []Message              `json:"messages"`
	Model            string                 `json:"model,omitempty"`
	MaxTokens        *int                   `json:"max_tokens,omitempty"`
	Temperature      *float64               `json:"temperature,omitempty"`
	TopP             *float64               `json:"top_p,omitempty"`
	TopK             *int                   `json:"top_k,omitempty"`
	Stream           bool                   `json:"stream,omitempty"`
	Tools            []Tool                 `json:"tools,omitempty"`
	ToolChoice       interface{}            `json:"tool_choice,omitempty"`
	ResponseFormat   *ResponseFormat        `json:"response_format,omitempty"`
	SystemPrompt     string                 `json:"system,omitempty"`
	Stop             []string               `json:"stop,omitempty"`
	PresencePenalty  *float64               `json:"presence_penalty,omitempty"`
	FrequencyPenalty *float64               `json:"frequency_penalty,omitempty"`
	User             string                 `json:"user,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
	Seed             *int                   `json:"seed,omitempty"`
	// Anthropic-specific
	AnthropicBeta []string        `json:"anthropic_beta,omitempty"`
	Thinking      *ThinkingConfig `json:"thinking,omitempty"`
	// DeepSeek-specific
	ReasoningEffort *string `json:"reasoning_effort,omitempty"`
}

// ThinkingConfig represents Anthropic thinking configuration
type ThinkingConfig struct {
	Type         string `json:"type"`          // "enabled"
	BudgetTokens int    `json:"budget_tokens"` // Must be >= 1024
}

// CompletionResponse represents a completion response
type CompletionResponse struct {
	ID                string     `json:"id"`
	Object            string     `json:"object"`
	Created           int64      `json:"created"`
	Model             string     `json:"model"`
	Provider          string     `json:"provider"` // Added Provider field
	SystemFingerprint string     `json:"system_fingerprint,omitempty"`
	Choices           []Choice   `json:"choices"`
	Usage             Usage      `json:"usage"`
	Container         *Container `json:"container,omitempty"` // For code execution
}

// Container represents a code execution container
type Container struct {
	ID        string `json:"id"`
	ExpiresAt string `json:"expires_at"`
}

// Choice represents a completion choice
type Choice struct {
	Index        int         `json:"index"`
	Message      Message     `json:"message"`
	Delta        Message     `json:"delta,omitempty"`
	FinishReason string      `json:"finish_reason"`
	LogProbs     interface{} `json:"logprobs,omitempty"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int     `json:"prompt_tokens"`
	CompletionTokens int     `json:"completion_tokens"`
	TotalTokens      int     `json:"total_tokens"`
	InputCost        float64 `json:"input_cost"`  // Added InputCost field
	OutputCost       float64 `json:"output_cost"` // Added OutputCost field
	TotalCost        float64 `json:"total_cost"`
	// Cache-related (for providers that support it)
	CacheCreationInputTokens int `json:"cache_creation_input_tokens,omitempty"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens,omitempty"`
	// Reasoning tokens (for DeepSeek-R1, OpenAI o1)
	CompletionTokensDetails *CompletionTokensDetails `json:"completion_tokens_details,omitempty"`
}

// CompletionTokensDetails provides detailed token breakdown
type CompletionTokensDetails struct {
	ReasoningTokens int `json:"reasoning_tokens,omitempty"`
}

// CompletionStreamResponse represents a streaming completion response
type CompletionStreamResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

// StreamCallbackFunc represents a callback function for streaming
type StreamCallbackFunc func(chunk CompletionStreamResponse) error

// ProviderConfig represents configuration for a provider
type ProviderConfig struct {
	APIKey       string   `json:"api_key"`
	BaseURL      string   `json:"base_url,omitempty"`
	DefaultModel string   `json:"default_model"`
	Enabled      bool     `json:"enabled"`
	MaxTokens    *int     `json:"max_tokens,omitempty"` // Added MaxTokens field
	Temperature  *float64 `json:"temperature,omitempty"`
	TopP         *float64 `json:"top_p,omitempty"`
	Organization string   `json:"organization,omitempty"` // For OpenAI
	Project      string   `json:"project,omitempty"`      // For OpenAI
}

// ClientCapabilities represents the capabilities of a client
type ClientCapabilities struct {
	SupportsStreaming    bool     `json:"supports_streaming"`
	SupportsTools        bool     `json:"supports_tools"`
	SupportsStructured   bool     `json:"supports_structured"`
	SupportsVision       bool     `json:"supports_vision"`
	SupportsMultiModal   bool     `json:"supports_multimodal"`
	SupportsReasoning    bool     `json:"supports_reasoning"`     // For DeepSeek-R1 type models
	SupportsThinking     bool     `json:"supports_thinking"`      // For Anthropic extended thinking
	SupportsAudio        bool     `json:"supports_audio"`         // For audio transcription/synthesis
	SupportsImageGen     bool     `json:"supports_image_gen"`     // For image generation
	SupportsToolRegistry bool     `json:"supports_tool_registry"` // For dynamic tool management
	MaxContextLength     int      `json:"max_context_length"`
	MaxOutputLength      int      `json:"max_output_length"`
	SupportedModels      []string `json:"supported_models"`
}

// ProviderInfo contains information about a provider
type ProviderInfo struct {
	Name         string             `json:"name"`
	Provider     string             `json:"provider"` // Added Provider field
	Version      string             `json:"version"`  // Added Version field
	Capabilities ClientCapabilities `json:"capabilities"`
	Models       []ModelInfo        `json:"models"`
	Status       string             `json:"status"`
}

// ModelInfo provides detailed information about a model
type ModelInfo struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Provider        string   `json:"provider"`   // Added Provider field
	Version         string   `json:"version"`    // Added Version field
	MaxTokens       int      `json:"max_tokens"` // Added MaxTokens field
	ContextLength   int      `json:"context_length"`
	MaxOutputTokens int      `json:"max_output_tokens"`
	InputCost       float64  `json:"input_cost"`  // Added InputCost field
	OutputCost      float64  `json:"output_cost"` // Added OutputCost field
	Capabilities    []string `json:"capabilities"`
	Description     string   `json:"description,omitempty"`
	TrainingCutoff  string   `json:"training_cutoff,omitempty"`
}

// HealthStatus represents the health status of a provider
type HealthStatus struct {
	Provider  string        `json:"provider"`
	Healthy   bool          `json:"healthy"`
	LastCheck time.Time     `json:"last_check"`
	Error     string        `json:"error,omitempty"`
	Latency   time.Duration `json:"latency"`
}

// UsageStats represents usage statistics
type UsageStats struct {
	Provider       string        `json:"provider"`
	RequestCount   int64         `json:"request_count"`
	TokensUsed     int64         `json:"tokens_used"`
	ErrorCount     int64         `json:"error_count"`
	AverageLatency time.Duration `json:"average_latency"`
	LastUsed       time.Time     `json:"last_used"`
	TotalCost      float64       `json:"total_cost"`
}

// SecurityAssessment represents a security assessment result
type SecurityAssessment struct {
	RiskLevel       string   `json:"risk_level"`
	FraudScore      float64  `json:"fraud_score"`
	Recommendations []string `json:"recommendations"`
	RedFlags        []string `json:"red_flags"`
	Confidence      float64  `json:"confidence"`
}

// ==================== ENHANCED AI SERVICE INTERFACE ====================

// EnhancedAIService defines the main interface for AI providers with multimodal capabilities
type EnhancedAIService interface {
	// Core completion methods
	CreateCompletion(ctx context.Context, request CompletionRequest) (*CompletionResponse, error)
	CreateCompletionStream(ctx context.Context, request CompletionRequest) (<-chan CompletionStreamResponse, error)

	// Tool/function calling
	ExecuteWithTools(ctx context.Context, messages []Message, tools []ToolDefinition) (*CompletionResponse, error)

	// Structured output
	CreateStructuredCompletion(ctx context.Context, request CompletionRequest, schema *JSONSchemaDefinition) (*CompletionResponse, error)

	// Utility methods
	CountTokens(text string) (int, error)
	GetCapabilities() ClientCapabilities
	HealthCheck(ctx context.Context) error
	GetUsageStats() UsageStats

	// Security features
	AnalyzeFraud(ctx context.Context, content string) (*SecurityAssessment, error)
	AssessRisk(ctx context.Context, request CompletionRequest) (*SecurityAssessment, error)
	GetSecurityRecommendations(ctx context.Context, content string) ([]string, error)

	// Provider info
	GetProviderInfo() ProviderInfo
}

// MultiModalAIService defines the extended interface for multimodal AI capabilities
type MultiModalAIService interface {
	EnhancedAIService

	// Vision capabilities
	AnalyzeImage(ctx context.Context, request interface{}) (*CompletionResponse, error)
	AnalyzeImages(ctx context.Context, imagePaths []string, prompt string) (*CompletionResponse, error)

	// Audio capabilities
	TranscribeAudio(ctx context.Context, request interface{}) (interface{}, error)
	TranslateAudio(ctx context.Context, audioPath string) (interface{}, error)
	SynthesizeSpeech(ctx context.Context, request interface{}) (io.Reader, error)

	// Image generation capabilities
	GenerateImage(ctx context.Context, request interface{}) (interface{}, error)
	EditImage(ctx context.Context, imagePath, maskPath, prompt string) (interface{}, error)
	CreateImageVariation(ctx context.Context, imagePath string, n int) (interface{}, error)

	// Tool registry capabilities
	RegisterTool(tool ToolDefinition) error
	UnregisterTool(name string) error
	ListRegisteredTools() []ToolDefinition
	GetRegisteredTool(name string) (ToolDefinition, bool)
	ExecuteWithRegisteredTools(ctx context.Context, messages []Message, toolNames []string) (*CompletionResponse, error)
}

// AIClientFactory defines the interface for creating AI clients
type AIClientFactory interface {
	CreateClient(provider string, config ProviderConfig) (EnhancedAIService, error)
	CreateMultiModalClient(provider string, config ProviderConfig) (MultiModalAIService, error)
	GetSupportedProviders() []string
	GetMultiModalProviders() []string
	SupportsMultiModal(provider string) bool
	RegisterProvider(provider string, config ProviderConfig) error
	GetProviderConfig(provider string) (ProviderConfig, bool)
}

// AIClientManager defines the interface for managing multiple AI clients
type AIClientManager interface {
	// Client management
	GetClient(provider string) (EnhancedAIService, error)
	GetDefaultClient() (EnhancedAIService, error)
	SetDefaultProvider(provider string) error

	// Multimodal client management
	GetMultiModalClient(provider string) (MultiModalAIService, error)
	GetDefaultMultiModalClient() (MultiModalAIService, error)

	// Enhanced operations with fallback
	CreateCompletionWithFallback(ctx context.Context, request CompletionRequest) (*CompletionResponse, error)
	CreateCompletionWithBestProvider(ctx context.Context, request CompletionRequest) (*CompletionResponse, error)

	// Provider management
	ListProviders() []string
	GetProviderHealth() map[string]HealthStatus
	GetAllUsageStats() map[string]UsageStats

	// Configuration
	EnableProvider(provider string) error
	DisableProvider(provider string) error
	UpdateProviderConfig(provider string, config ProviderConfig) error
}

// AIError represents an AI-specific error
type AIError struct {
	Provider string                 `json:"provider"`
	Type     string                 `json:"type"`
	Message  string                 `json:"message"`
	Code     string                 `json:"code,omitempty"`
	Details  map[string]interface{} `json:"details,omitempty"`
}

func (e *AIError) Error() string {
	return e.Message
}

// StreamReader defines the interface for streaming responses
type StreamReader interface {
	io.Closer
	ReadChunk() (*CompletionStreamResponse, error)
}

// Additional types needed for assistants service compilation

// APIError represents a standard API error
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

func (e *APIError) Error() string {
	return e.Message
}

// ContentFilterResult represents content filtering results
type ContentFilterResult struct {
	Filtered bool   `json:"filtered"`
	Category string `json:"category"`
	Severity string `json:"severity"`
	Reason   string `json:"reason"`
}

// ContentFilterError represents content filtering errors
type ContentFilterError struct {
	APIError
	Results map[string]ContentFilterResult `json:"results"`
}

// RateLimitInfo contains rate limit information
type RateLimitInfo struct {
	RequestsRemaining int           `json:"requests_remaining"`
	TokensRemaining   int           `json:"tokens_remaining"`
	ResetTime         time.Time     `json:"reset_time"`
	RetryAfter        time.Duration `json:"retry_after"`
}

// RateLimitError represents rate limiting errors
type RateLimitError struct {
	APIError
	RateLimitInfo RateLimitInfo `json:"rate_limit_info"`
}

// StreamChunk represents a chunk in streaming response (alias for compatibility)
type StreamChunk = CompletionStreamResponse

// StreamOptions represents streaming options
type StreamOptions struct {
	IncludeUsage bool `json:"include_usage"`
}

// RequireContext validates that a context is present and valid
func RequireContext(ctx context.Context) error {
	if ctx == nil {
		return &APIError{
			Code:    "missing_context",
			Message: "Context is required",
			Type:    "invalid_request_error",
		}
	}
	return nil
}
