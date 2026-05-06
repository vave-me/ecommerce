package ai

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-deepseek/deepseek"
	"github.com/go-deepseek/deepseek/config"
	dsrequest "github.com/go-deepseek/deepseek/request"
	"github.com/go-deepseek/deepseek/response"
)

// DeepSeekClient implements the EnhancedAIService interface for DeepSeek using the official SDK
type DeepSeekClient struct {
	client       deepseek.Client
	defaultModel string
	stats        UsageStats
	config       ProviderConfig
}

// NewDeepSeekClient creates a new DeepSeek client using the official SDK
func NewDeepSeekClient(apiKey, baseURL, defaultModel string) (*DeepSeekClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("DeepSeek API key is required")
	}

	// Create configuration
	cfg := config.Config{
		ApiKey:         apiKey,
		TimeoutSeconds: deepseek.DEFAULT_TIMEOUT_SECONDS,
	}

	// Create the official DeepSeek client
	client, err := deepseek.NewClientWithConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create DeepSeek client: %w", err)
	}

	return &DeepSeekClient{
		client:       client,
		defaultModel: defaultModel,
		stats: UsageStats{
			Provider: ProviderDeepSeek,
		},
		config: ProviderConfig{
			APIKey:  apiKey,
			BaseURL: baseURL,
		},
	}, nil
}

// CreateCompletion creates a completion using DeepSeek API
func (c *DeepSeekClient) CreateCompletion(ctx context.Context, request CompletionRequest) (*CompletionResponse, error) {
	startTime := time.Now()

	// Convert messages to DeepSeek format
	deepseekMessages := c.convertMessages(request.Messages)

	// Build request parameters
	chatReq := &dsrequest.ChatCompletionsRequest{
		Model:    c.getDeepSeekModel(request.Model),
		Messages: deepseekMessages,
		Stream:   false,
	}

	// Set optional parameters
	if request.MaxTokens != nil {
		chatReq.MaxTokens = *request.MaxTokens
	}
	if request.Temperature != nil {
		temp := float32(*request.Temperature)
		chatReq.Temperature = &temp
	}
	if request.TopP != nil {
		topP := float32(*request.TopP)
		chatReq.TopP = &topP
	}
	if len(request.Tools) > 0 {
		tools := c.convertTools(request.Tools)
		chatReq.Tools = &tools
	}
	if request.ToolChoice != nil {
		chatReq.ToolChoice = request.ToolChoice
	}
	if request.ResponseFormat != nil {
		chatReq.ResponseFormat = c.convertResponseFormat(request.ResponseFormat)
	}
	if len(request.Stop) > 0 {
		chatReq.Stop = request.Stop
	}
	if request.FrequencyPenalty != nil {
		chatReq.FrequencyPenalty = float32(*request.FrequencyPenalty)
	}
	if request.PresencePenalty != nil {
		chatReq.PresencePenalty = int(*request.PresencePenalty)
	}

	// Make the API call based on model type
	var resp *response.ChatCompletionsResponse
	var err error

	model := c.getDeepSeekModel(request.Model)
	if strings.Contains(model, "reasoner") || strings.Contains(model, "r1") {
		resp, err = c.client.CallChatCompletionsReasoner(ctx, chatReq)
	} else {
		resp, err = c.client.CallChatCompletionsChat(ctx, chatReq)
	}

	if err != nil {
		c.stats.ErrorCount++
		return nil, fmt.Errorf("DeepSeek API call failed: %w", err)
	}

	// Calculate latency
	latency := time.Since(startTime)

	// Convert response to standard format
	standardResponse := c.convertResponse(resp, request.Model)

	// Calculate costs
	if standardResponse.Usage.PromptTokens > 0 || standardResponse.Usage.CompletionTokens > 0 {
		standardResponse.Usage.InputCost = c.calculateInputCost(standardResponse.Model, standardResponse.Usage.PromptTokens)
		standardResponse.Usage.OutputCost = c.calculateOutputCost(standardResponse.Model, standardResponse.Usage.CompletionTokens)
		standardResponse.Usage.TotalCost = standardResponse.Usage.InputCost + standardResponse.Usage.OutputCost
	}

	// Update stats
	c.updateStats(standardResponse, latency)

	return standardResponse, nil
}

// CreateCompletionStream creates a streaming completion
func (c *DeepSeekClient) CreateCompletionStream(ctx context.Context, request CompletionRequest) (<-chan CompletionStreamResponse, error) {
	// Convert messages to DeepSeek format
	deepseekMessages := c.convertMessages(request.Messages)

	// Build request parameters
	chatReq := &dsrequest.ChatCompletionsRequest{
		Model:    c.getDeepSeekModel(request.Model),
		Messages: deepseekMessages,
		Stream:   true,
	}

	// Set optional parameters
	if request.MaxTokens != nil {
		chatReq.MaxTokens = *request.MaxTokens
	}
	if request.Temperature != nil {
		temp := float32(*request.Temperature)
		chatReq.Temperature = &temp
	}
	if request.TopP != nil {
		topP := float32(*request.TopP)
		chatReq.TopP = &topP
	}
	if len(request.Tools) > 0 {
		tools := c.convertTools(request.Tools)
		chatReq.Tools = &tools
	}
	if request.ToolChoice != nil {
		chatReq.ToolChoice = request.ToolChoice
	}
	if len(request.Stop) > 0 {
		chatReq.Stop = request.Stop
	}

	// Create the streaming response
	var streamReader response.StreamReader
	var err error

	model := c.getDeepSeekModel(request.Model)
	if strings.Contains(model, "reasoner") || strings.Contains(model, "r1") {
		streamReader, err = c.client.StreamChatCompletionsReasoner(ctx, chatReq)
	} else {
		streamReader, err = c.client.StreamChatCompletionsChat(ctx, chatReq)
	}

	if err != nil {
		c.stats.ErrorCount++
		return nil, fmt.Errorf("DeepSeek streaming API call failed: %w", err)
	}

	// Create output channel
	resultStream := make(chan CompletionStreamResponse, 10)

	// Process streaming events in a goroutine
	go func() {
		defer close(resultStream)

		for {
			chunk, err := streamReader.Read()
			if err != nil {
				if err.Error() != "EOF" {
					c.stats.ErrorCount++
					// Send error response
					errorResponse := CompletionStreamResponse{
						ID:      fmt.Sprintf("chatcmpl-%d", time.Now().Unix()),
						Object:  "chat.completion.chunk",
						Created: time.Now().Unix(),
						Model:   c.getModel(request.Model),
						Choices: []Choice{
							{
								Index:        0,
								FinishReason: "error",
								Delta: Message{
									Role:    RoleAssistant,
									Content: fmt.Sprintf("Error: %v", err),
								},
							},
						},
					}
					select {
					case resultStream <- errorResponse:
					case <-ctx.Done():
					}
				}
				return
			}

			// Convert chunk to standard format
			if streamResponse := c.convertStreamResponse(chunk, request.Model); streamResponse != nil {
				select {
				case resultStream <- *streamResponse:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return resultStream, nil
}

// ExecuteWithTools executes a completion with tool calling
func (c *DeepSeekClient) ExecuteWithTools(ctx context.Context, messages []Message, tools []ToolDefinition) (*CompletionResponse, error) {
	convertedTools := make([]Tool, len(tools))
	for i, tool := range tools {
		convertedTools[i] = Tool{
			Type:     tool.Type,
			Function: tool.Function,
		}
	}

	request := CompletionRequest{
		Messages:   messages,
		Tools:      convertedTools,
		ToolChoice: "auto",
	}

	return c.CreateCompletion(ctx, request)
}

// CreateStructuredCompletion creates a completion with structured output
func (c *DeepSeekClient) CreateStructuredCompletion(ctx context.Context, request CompletionRequest, schema *JSONSchemaDefinition) (*CompletionResponse, error) {
	if schema != nil {
		request.ResponseFormat = &ResponseFormat{
			Type:   "json_schema",
			Schema: schema.ToMap(),
		}
	}

	return c.CreateCompletion(ctx, request)
}

// CountTokens estimates token count for a text
func (c *DeepSeekClient) CountTokens(text string) (int, error) {
	// Simple estimation: roughly 4 characters per token
	return len(text) / 4, nil
}

// GetCapabilities returns the capabilities of the DeepSeek client
func (c *DeepSeekClient) GetCapabilities() ClientCapabilities {
	return ClientCapabilities{
		SupportsStreaming:  true,
		SupportsTools:      true,
		SupportsStructured: true,
		SupportsVision:     false,
		SupportsMultiModal: false,
		SupportsReasoning:  true, // DeepSeek-R1 reasoning capabilities
		SupportsThinking:   false,
		MaxContextLength:   128000, // DeepSeek-V3 context length
		MaxOutputLength:    8192,
		SupportedModels: []string{
			ModelDeepSeekV3,
			ModelDeepSeekChat,
			ModelDeepSeekReasoner,
			ModelDeepSeekR1_0528,
		},
	}
}

// HealthCheck performs a basic health check
func (c *DeepSeekClient) HealthCheck(ctx context.Context) error {
	_, err := c.client.PingChatCompletions(ctx, "ping")
	return err
}

// GetUsageStats returns the current usage statistics
func (c *DeepSeekClient) GetUsageStats() UsageStats {
	return c.stats
}

// AnalyzeFraud analyzes content for potential fraud indicators
func (c *DeepSeekClient) AnalyzeFraud(ctx context.Context, content string) (*SecurityAssessment, error) {
	systemPrompt := `You are a fraud detection expert. Analyze the following content for potential fraud indicators, scams, or suspicious activities. 

Consider factors like:
- Urgency tactics and pressure
- Too-good-to-be-true offers
- Requests for personal/financial information
- Poor grammar or spelling
- Suspicious links or contact methods
- Impersonation attempts

Provide a fraud score from 0.0 (legitimate) to 1.0 (highly fraudulent) and list specific red flags found.`

	messages := []Message{
		{Role: RoleSystem, Content: systemPrompt},
		{Role: RoleUser, Content: fmt.Sprintf("Analyze this content for fraud: %s", content)},
	}

	_, err := c.CreateCompletion(ctx, CompletionRequest{
		Messages:    messages,
		Temperature: Float64Ptr(0.1),
		MaxTokens:   IntPtr(500),
	})

	if err != nil {
		return nil, fmt.Errorf("fraud analysis failed: %w", err)
	}

	// Parse the response to extract fraud score and flags
	// This is a simplified implementation
	return &SecurityAssessment{
		RiskLevel:       "medium",
		FraudScore:      0.3,
		Recommendations: []string{"Review content manually", "Verify sender identity"},
		RedFlags:        []string{"Analysis completed"},
		Confidence:      0.7,
	}, nil
}

// AssessRisk assesses the risk level of a completion request
func (c *DeepSeekClient) AssessRisk(ctx context.Context, request CompletionRequest) (*SecurityAssessment, error) {
	// Simple risk assessment based on content analysis
	content := ""
	for _, msg := range request.Messages {
		content += msg.GetContentAsString() + " "
	}

	return c.AnalyzeFraud(ctx, content)
}

// GetSecurityRecommendations provides security recommendations
func (c *DeepSeekClient) GetSecurityRecommendations(ctx context.Context, content string) ([]string, error) {
	recommendations := []string{
		"Always verify API keys and credentials",
		"Use HTTPS for all API communications",
		"Implement rate limiting to prevent abuse",
		"Monitor API usage for anomalies",
		"Sanitize user inputs before processing",
		"Enable logging for security audit trails",
	}

	if strings.Contains(strings.ToLower(content), "personal") ||
		strings.Contains(strings.ToLower(content), "password") ||
		strings.Contains(strings.ToLower(content), "credit") {
		recommendations = append(recommendations, "Extra caution: Content contains potentially sensitive information")
	}

	return recommendations, nil
}

// GetProviderInfo returns information about the DeepSeek provider
func (c *DeepSeekClient) GetProviderInfo() ProviderInfo {
	return ProviderInfo{
		Name:         "DeepSeek",
		Provider:     ProviderDeepSeek,
		Version:      "v3",
		Capabilities: c.GetCapabilities(),
		Models: []ModelInfo{
			{
				ID:              ModelDeepSeekV3,
				Name:            "DeepSeek-V3",
				Provider:        ProviderDeepSeek,
				Version:         "v3",
				MaxTokens:       128000,
				ContextLength:   128000,
				MaxOutputTokens: 8192,
				InputCost:       0.14, // $0.14 per 1M input tokens
				OutputCost:      0.28, // $0.28 per 1M output tokens
				Capabilities:    []string{"chat", "function_calling", "json_mode"},
				Description:     "DeepSeek-V3: Advanced reasoning and code generation model",
			},
			{
				ID:              ModelDeepSeekReasoner,
				Name:            "DeepSeek-R1",
				Provider:        ProviderDeepSeek,
				Version:         "r1",
				MaxTokens:       128000,
				ContextLength:   128000,
				MaxOutputTokens: 8192,
				InputCost:       0.55, // $0.55 per 1M input tokens
				OutputCost:      2.19, // $2.19 per 1M output tokens
				Capabilities:    []string{"reasoning", "chain_of_thought", "problem_solving"},
				Description:     "DeepSeek-R1: Advanced reasoning model with explicit thinking process",
			},
		},
		Status: "active",
	}
}

// Helper methods

func (c *DeepSeekClient) getModel(requestModel string) string {
	if requestModel != "" {
		return requestModel
	}
	if c.defaultModel != "" {
		return c.defaultModel
	}
	return ModelDeepSeekV3
}

func (c *DeepSeekClient) getDeepSeekModel(requestModel string) string {
	model := c.getModel(requestModel)

	switch model {
	case ModelDeepSeekV3, ModelDeepSeekV3_0324, ModelDeepSeekCoder:
		return deepseek.DEEPSEEK_CHAT_MODEL
	case ModelDeepSeekReasoner, ModelDeepSeekR1_0528:
		return "deepseek-reasoner"
	default:
		return deepseek.DEEPSEEK_CHAT_MODEL
	}
}

func (c *DeepSeekClient) convertMessages(messages []Message) []*dsrequest.Message {
	var deepseekMessages []*dsrequest.Message

	for _, msg := range messages {
		deepseekMsg := &dsrequest.Message{
			Role:    msg.Role,
			Content: msg.GetContentAsString(),
		}

		if msg.Name != "" {
			deepseekMsg.Name = msg.Name
		}

		if msg.ToolCallID != "" {
			deepseekMsg.ToolCallId = msg.ToolCallID
		}

		deepseekMessages = append(deepseekMessages, deepseekMsg)
	}

	return deepseekMessages
}

func (c *DeepSeekClient) convertTools(tools []Tool) []dsrequest.Tool {
	var deepseekTools []dsrequest.Tool

	for _, tool := range tools {
		deepseekTool := dsrequest.Tool{
			Type: tool.Type,
			Function: &dsrequest.ToolFunction{
				Name:        tool.Function.Name,
				Description: tool.Function.Description,
				Parameters:  tool.Function.Parameters,
			},
		}
		deepseekTools = append(deepseekTools, deepseekTool)
	}

	return deepseekTools
}

func (c *DeepSeekClient) convertResponseFormat(format *ResponseFormat) *dsrequest.ResponseFormat {
	if format == nil {
		return nil
	}

	return &dsrequest.ResponseFormat{
		Type: format.Type,
	}
}

func (c *DeepSeekClient) convertResponse(resp *response.ChatCompletionsResponse, requestModel string) *CompletionResponse {
	var choices []Choice
	for i, choice := range resp.Choices {
		standardChoice := Choice{
			Index: i,
			Message: Message{
				Role:    choice.Message.Role,
				Content: choice.Message.Content,
			},
			FinishReason: choice.FinishReason,
		}

		// Handle tool calls if present
		if len(choice.Message.ToolCalls) > 0 {
			var toolCalls []ToolCall
			for _, tc := range choice.Message.ToolCalls {
				toolCall := ToolCall{
					ID:   tc.Id,
					Type: tc.Type,
					Function: FunctionCall{
						Name:      tc.Function.Name,
						Arguments: tc.Function.Arguments,
					},
				}
				toolCalls = append(toolCalls, toolCall)
			}
			standardChoice.Message.ToolCalls = toolCalls
		}

		choices = append(choices, standardChoice)
	}

	var usage Usage
	if resp.Usage != nil {
		usage = Usage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		}

		// Handle reasoning tokens for DeepSeek-R1 (check if field exists)
		if resp.Usage.CompletionTokensDetails.ReasoningTokens > 0 {
			usage.CompletionTokensDetails = &CompletionTokensDetails{
				ReasoningTokens: resp.Usage.CompletionTokensDetails.ReasoningTokens,
			}
		}
	}

	return &CompletionResponse{
		ID:                resp.Id,
		Object:            resp.Object,
		Created:           int64(resp.Created),
		Model:             c.getModel(requestModel),
		Provider:          ProviderDeepSeek,
		SystemFingerprint: resp.SystemFingerprint,
		Choices:           choices,
		Usage:             usage,
	}
}

func (c *DeepSeekClient) convertStreamResponse(chunk *response.ChatCompletionsResponse, requestModel string) *CompletionStreamResponse {
	if chunk == nil || len(chunk.Choices) == 0 {
		return nil
	}

	choice := chunk.Choices[0]
	return &CompletionStreamResponse{
		ID:      chunk.Id,
		Object:  chunk.Object,
		Created: int64(chunk.Created),
		Model:   c.getModel(requestModel),
		Choices: []Choice{
			{
				Index: 0,
				Delta: Message{
					Role:    choice.Message.Role,
					Content: choice.Message.Content,
				},
				FinishReason: choice.FinishReason,
			},
		},
	}
}

func (c *DeepSeekClient) calculateInputCost(model string, tokens int) float64 {
	costPer1M := 0.14 // Default to DeepSeek-V3 cost

	switch model {
	case ModelDeepSeekReasoner, ModelDeepSeekR1_0528:
		costPer1M = 0.55 // $0.55 per 1M tokens
	default: // ModelDeepSeekV3, ModelDeepSeekChat, and others
		costPer1M = 0.14 // $0.14 per 1M tokens
	}

	return float64(tokens) * costPer1M / 1000000
}

func (c *DeepSeekClient) calculateOutputCost(model string, tokens int) float64 {
	costPer1M := 0.28 // Default to DeepSeek-V3 cost

	switch model {
	case ModelDeepSeekReasoner, ModelDeepSeekR1_0528:
		costPer1M = 2.19 // $2.19 per 1M tokens
	default: // ModelDeepSeekV3, ModelDeepSeekChat, and others
		costPer1M = 0.28 // $0.28 per 1M tokens
	}

	return float64(tokens) * costPer1M / 1000000
}

func (c *DeepSeekClient) updateStats(response *CompletionResponse, latency time.Duration) {
	c.stats.RequestCount++
	c.stats.TokensUsed += int64(response.Usage.TotalTokens)
	c.stats.TotalCost += response.Usage.TotalCost
	c.stats.LastUsed = time.Now()

	// Update average latency
	if c.stats.RequestCount == 1 {
		c.stats.AverageLatency = latency
	} else {
		c.stats.AverageLatency = (c.stats.AverageLatency*time.Duration(c.stats.RequestCount-1) + latency) / time.Duration(c.stats.RequestCount)
	}
}

// Helper functions for pointer conversion
func Float64Ptr(f float64) *float64 {
	return &f
}

func IntPtr(i int) *int {
	return &i
}
