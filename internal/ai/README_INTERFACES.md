# AI Interfaces Documentation

## Overview

The `assistants/internal/ai` package provides a comprehensive interface system for AI providers with support for both standard and multimodal capabilities. This document outlines the interfaces, their capabilities, and how to use them.

## Interface Hierarchy

### 1. EnhancedAIService
The base interface that all AI providers must implement. Provides core text-based AI capabilities:

```go
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
```

### 2. MultiModalAIService
Extended interface for providers that support multimodal capabilities (vision, audio, image generation):

```go
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
```

## Supported Providers

### OpenAI (Full MultiModal Support) ✅
- **Text**: GPT-4o, GPT-4-Turbo, O1-Preview, O1-Mini
- **Vision**: Image analysis with GPT-4o
- **Audio**: Whisper transcription, text-to-speech synthesis
- **Image Generation**: DALL-E 2 & 3
- **Tools**: Function calling with tool registry
- **Features**: Streaming, structured output, reasoning tokens

**Capabilities**: 
```go
ClientCapabilities{
    SupportsStreaming:     true,
    SupportsTools:         true,
    SupportsStructured:    true,
    SupportsVision:        true,
    SupportsMultiModal:    true,
    SupportsReasoning:     true,
    SupportsThinking:      false,
    SupportsAudio:         true,
    SupportsImageGen:      true,
    SupportsToolRegistry:  true,
}
```

### Anthropic (Text + Limited Vision) ✅
- **Text**: Claude-4 Opus/Sonnet, Claude-3.7-Sonnet, Claude-3.5-Sonnet/Haiku
- **Vision**: Limited image analysis
- **Features**: Extended thinking, beta features, streaming
- **Limitations**: No audio or image generation

**Capabilities**:
```go
ClientCapabilities{
    SupportsStreaming:     true,
    SupportsTools:         true,
    SupportsStructured:    true,
    SupportsVision:        true,
    SupportsMultiModal:    false,
    SupportsReasoning:     false,
    SupportsThinking:      true,
    SupportsAudio:         false,
    SupportsImageGen:      false,
    SupportsToolRegistry:  false,
}
```

### DeepSeek (Text + Reasoning) ✅
- **Text**: DeepSeek-V3, DeepSeek-Chat
- **Reasoning**: DeepSeek-R1, DeepSeek-Reasoner
- **Features**: Streaming, function calling, reasoning tokens
- **Limitations**: No vision, audio, or image generation

**Capabilities**:
```go
ClientCapabilities{
    SupportsStreaming:     true,
    SupportsTools:         true,
    SupportsStructured:    true,
    SupportsVision:        false,
    SupportsMultiModal:    false,
    SupportsReasoning:     true,
    SupportsThinking:      false,
    SupportsAudio:         false,
    SupportsImageGen:      false,
    SupportsToolRegistry:  false,
}
```

## Factory and Manager Interfaces

### AIClientFactory
Creates and manages AI client instances:

```go
type AIClientFactory interface {
    CreateClient(provider string, config ProviderConfig) (EnhancedAIService, error)
    CreateMultiModalClient(provider string, config ProviderConfig) (MultiModalAIService, error)
    GetSupportedProviders() []string
    GetMultiModalProviders() []string
    SupportsMultiModal(provider string) bool
    RegisterProvider(provider string, config ProviderConfig) error
    GetProviderConfig(provider string) (ProviderConfig, bool)
}
```

### AIClientManager
High-level management with fallback and load balancing:

```go
type AIClientManager interface {
    // Standard client management
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
```

## Usage Examples

### Basic Text Completion
```go
factory := NewClientFactory()
config := ProviderConfig{
    APIKey: "your-api-key",
    DefaultModel: "gpt-4o",
}
factory.RegisterProvider(ProviderOpenAI, config)

client, err := factory.CreateClient(ProviderOpenAI, config)
if err != nil {
    return err
}

response, err := client.CreateCompletion(ctx, CompletionRequest{
    Messages: []Message{
        {Role: RoleUser, Content: "Hello, world!"},
    },
})
```

### Multimodal Operations (OpenAI Only)
```go
// Get multimodal client
multiModalClient, err := factory.CreateMultiModalClient(ProviderOpenAI, config)
if err != nil {
    return err
}

// Vision analysis
visionRequest := VisionAnalysisRequest{
    ImagePath: "path/to/image.jpg",
    Questions: []string{"What do you see in this image?"},
    MaxTokens: 300,
}
response, err := multiModalClient.AnalyzeImage(ctx, visionRequest)

// Audio transcription
audioRequest := AudioTranscriptionRequest{
    AudioPath: "path/to/audio.wav",
    Language: "en",
}
transcription, err := multiModalClient.TranscribeAudio(ctx, audioRequest)

// Image generation
imageRequest := ImageGenerationRequest{
    Prompt: "A beautiful sunset over mountains",
    Model: "dall-e-3",
    Size: "1024x1024",
}
images, err := multiModalClient.GenerateImage(ctx, imageRequest)
```

### Client Manager with Fallback
```go
manager := NewClientManager(factory)

// Register multiple providers
manager.factory.RegisterProvider(ProviderOpenAI, openaiConfig)
manager.factory.RegisterProvider(ProviderAnthropic, anthropicConfig)

// Automatic fallback if primary provider fails
response, err := manager.CreateCompletionWithFallback(ctx, request)

// Choose best provider based on request requirements
response, err := manager.CreateCompletionWithBestProvider(ctx, request)
```

## Key Features

### Enhanced Capabilities
- **Streaming Support**: Real-time response streaming for all providers
- **Tool/Function Calling**: Dynamic function execution with schema validation
- **Structured Output**: JSON schema-enforced responses
- **Security Features**: Fraud detection, risk assessment, content filtering
- **Cost Tracking**: Detailed usage statistics and cost calculation
- **Health Monitoring**: Provider health checks and latency tracking

### Multimodal Features (OpenAI)
- **Vision**: Image analysis, multi-image processing
- **Audio**: Speech-to-text transcription, translation, text-to-speech synthesis
- **Image Generation**: DALL-E integration with editing and variations
- **Tool Registry**: Dynamic tool registration and management

### Provider-Specific Enhancements
- **OpenAI**: Full multimodal support, latest models (GPT-4o, O1, DALL-E 3)
- **Anthropic**: Extended thinking, beta features, Claude-4 models
- **DeepSeek**: Reasoning models (R1), cost-effective chat models

## Data Structures

The package includes comprehensive data structures for:
- Request/Response objects for all operations
- Provider configuration and capabilities
- Usage statistics and health monitoring
- Security assessments and recommendations
- Multimodal content handling

## Thread Safety

All client implementations and managers are thread-safe with proper mutex protection for concurrent operations.

## Migration Notes

The interfaces maintain backward compatibility while adding new multimodal capabilities. Existing code using `EnhancedAIService` will continue to work without changes. To access multimodal features, upgrade to `MultiModalAIService` for supported providers. 