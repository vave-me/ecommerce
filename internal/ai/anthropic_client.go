package ai

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// AnthropicClient implements the EnhancedAIService interface for Anthropic Claude
type AnthropicClient struct {
	client       *anthropic.Client
	defaultModel string
	stats        UsageStats
	config       ProviderConfig
}

// NewAnthropicClient creates a new Anthropic client using the official SDK
func NewAnthropicClient(apiKey, baseURL string, defaultModel string) (*AnthropicClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("Anthropic API key is required")
	}

	// Create client options
	opts := []option.RequestOption{
		option.WithAPIKey(apiKey),
	}

	// Add custom base URL if provided
	if baseURL != "" && baseURL != "https://api.anthropic.com/v1" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}

	// Create the official Anthropic client
	client := anthropic.NewClient(opts...)

	return &AnthropicClient{
		client:       &client,
		defaultModel: defaultModel,
		stats: UsageStats{
			Provider: ProviderAnthropic,
		},
		config: ProviderConfig{
			APIKey:  apiKey,
			BaseURL: baseURL,
		},
	}, nil
}

// CreateCompletion creates a completion using Anthropic API
func (c *AnthropicClient) CreateCompletion(ctx context.Context, request CompletionRequest) (*CompletionResponse, error) {
	startTime := time.Now()

	// Convert messages to Anthropic format
	anthropicMessages, systemMessage := c.convertMessages(request.Messages)

	// Build request parameters
	params := anthropic.MessageNewParams{
		Model:    c.getAnthropicModel(request.Model),
		Messages: anthropicMessages,
	}

	// Set system message if provided
	if systemMessage != "" {
		params.System = []anthropic.TextBlockParam{
			{
				Text: systemMessage,
				Type: "text",
			},
		}
	}

	// Set max tokens (required by Anthropic)
	if request.MaxTokens != nil {
		params.MaxTokens = int64(*request.MaxTokens)
	} else {
		params.MaxTokens = 1000 // Default required value
	}

	// Set optional parameters
	if request.Temperature != nil {
		params.Temperature = anthropic.Float(*request.Temperature)
	}
	if request.TopP != nil {
		params.TopP = anthropic.Float(*request.TopP)
	}
	if request.TopK != nil {
		params.TopK = anthropic.Int(int64(*request.TopK))
	}

	// Convert tools
	if len(request.Tools) > 0 {
		params.Tools = c.convertTools(request.Tools)
	}

	// Convert stop sequences
	if len(request.Stop) > 0 {
		params.StopSequences = request.Stop
	}

	// Make the API call
	response, err := c.client.Messages.New(ctx, params)
	if err != nil {
		c.stats.ErrorCount++
		return nil, fmt.Errorf("Anthropic API call failed: %w", err)
	}

	// Calculate latency
	latency := time.Since(startTime)

	// Convert response to standard format
	standardResponse := c.convertResponse(response, request.Model)

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
func (c *AnthropicClient) CreateCompletionStream(ctx context.Context, request CompletionRequest) (<-chan CompletionStreamResponse, error) {
	// Convert messages to Anthropic format
	anthropicMessages, systemMessage := c.convertMessages(request.Messages)

	// Build request parameters
	params := anthropic.MessageNewParams{
		Model:    c.getAnthropicModel(request.Model),
		Messages: anthropicMessages,
	}

	// Set system message if provided
	if systemMessage != "" {
		params.System = []anthropic.TextBlockParam{
			{
				Text: systemMessage,
				Type: "text",
			},
		}
	}

	// Set max tokens (required by Anthropic)
	if request.MaxTokens != nil {
		params.MaxTokens = int64(*request.MaxTokens)
	} else {
		params.MaxTokens = 1000
	}

	// Set optional parameters
	if request.Temperature != nil {
		params.Temperature = anthropic.Float(*request.Temperature)
	}
	if request.TopP != nil {
		params.TopP = anthropic.Float(*request.TopP)
	}
	if request.TopK != nil {
		params.TopK = anthropic.Int(int64(*request.TopK))
	}

	// Convert tools
	if len(request.Tools) > 0 {
		params.Tools = c.convertTools(request.Tools)
	}

	// Convert stop sequences
	if len(request.Stop) > 0 {
		params.StopSequences = request.Stop
	}

	// Create the streaming response
	stream := c.client.Messages.NewStreaming(ctx, params)

	// Create output channel
	resultStream := make(chan CompletionStreamResponse, 10)

	// Process streaming events in a goroutine
	go func() {
		defer close(resultStream)
		defer stream.Close()

		for stream.Next() {
			event := stream.Current()
			if streamResponse := c.convertStreamEvent(event, request.Model); streamResponse != nil {
				select {
				case resultStream <- *streamResponse:
				case <-ctx.Done():
					return
				}
			}
		}

		if err := stream.Err(); err != nil {
			c.stats.ErrorCount++
			// Create error response with required fields
			errorResponse := CompletionStreamResponse{
				ID:      fmt.Sprintf("msg_%d", time.Now().Unix()),
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
	}()

	return resultStream, nil
}

// ExecuteWithTools executes a completion with tool calling
func (c *AnthropicClient) ExecuteWithTools(ctx context.Context, messages []Message, tools []ToolDefinition) (*CompletionResponse, error) {
	convertedTools := make([]Tool, len(tools))
	for i, tool := range tools {
		convertedTools[i] = Tool{
			Type:     tool.Type,
			Function: tool.Function,
		}
	}

	request := CompletionRequest{
		Messages: messages,
		Tools:    convertedTools,
	}

	return c.CreateCompletion(ctx, request)
}

// CreateStructuredCompletion creates a completion with structured output
func (c *AnthropicClient) CreateStructuredCompletion(ctx context.Context, request CompletionRequest, schema *JSONSchemaDefinition) (*CompletionResponse, error) {
	// Anthropic doesn't have native structured output, so we'll add instructions to the system message
	if schema != nil {
		systemMsg := "Please respond with valid JSON that follows this schema: " + c.schemaToString(schema)

		// Add system message to request
		systemMessage := Message{
			Role:    RoleSystem,
			Content: systemMsg,
		}
		request.Messages = append([]Message{systemMessage}, request.Messages...)
	}

	return c.CreateCompletion(ctx, request)
}

// CountTokens estimates token count using Anthropic's approach
func (c *AnthropicClient) CountTokens(text string) (int, error) {
	// Anthropic uses approximately 3.5 characters per token
	return len(text) / 4, nil
}

// GetCapabilities returns the capabilities of this client
func (c *AnthropicClient) GetCapabilities() ClientCapabilities {
	return ClientCapabilities{
		SupportsStreaming:  true,
		SupportsTools:      true,
		SupportsStructured: false, // No native structured output
		SupportsVision:     true,
		SupportsMultiModal: true,
		SupportsReasoning:  false,  // Claude doesn't have reasoning models like o1
		SupportsThinking:   true,   // Claude supports extended thinking
		MaxContextLength:   200000, // Claude 4 context length
		MaxOutputLength:    64000,  // Claude 4 max output
		SupportedModels: []string{
			ModelClaudeOpus4,
			ModelClaudeSonnet4,
			ModelClaude37Sonnet20250219,
			ModelClaude35Sonnet20241022,
			ModelClaude35Haiku20241022,
			ModelClaude3Opus20240229,
			ModelClaudeOpus4Latest,
			ModelClaudeSonnet4Latest,
			ModelClaude37SonnetLatest,
			ModelClaude35SonnetLatest,
			ModelClaude35HaikuLatest,
		},
	}
}

// HealthCheck performs a health check
func (c *AnthropicClient) HealthCheck(ctx context.Context) error {
	// Simple health check with minimal request
	testMessages := []Message{
		{Role: RoleUser, Content: "Hello"},
	}

	request := CompletionRequest{
		Messages:  testMessages,
		MaxTokens: func(i int) *int { return &i }(10),
	}

	_, err := c.CreateCompletion(ctx, request)
	return err
}

// GetUsageStats returns usage statistics
func (c *AnthropicClient) GetUsageStats() UsageStats {
	return c.stats
}

// AnalyzeFraud analyzes content for fraud indicators
func (c *AnthropicClient) AnalyzeFraud(ctx context.Context, content string) (*SecurityAssessment, error) {
	fraudScore := 0.1
	redFlags := []string{}
	recommendations := []string{"Monitor for suspicious patterns"}

	lowerContent := strings.ToLower(content)
	if strings.Contains(lowerContent, "urgent") {
		fraudScore += 0.2
		redFlags = append(redFlags, "Contains urgency indicators")
	}
	if strings.Contains(lowerContent, "money") || strings.Contains(lowerContent, "payment") {
		fraudScore += 0.3
		redFlags = append(redFlags, "Contains financial keywords")
	}
	if strings.Contains(lowerContent, "click here") || strings.Contains(lowerContent, "verify") {
		fraudScore += 0.2
		redFlags = append(redFlags, "Contains phishing indicators")
	}

	riskLevel := "low"
	if fraudScore > 0.5 {
		riskLevel = "medium"
	}
	if fraudScore > 0.8 {
		riskLevel = "high"
	}

	return &SecurityAssessment{
		RiskLevel:       riskLevel,
		FraudScore:      fraudScore,
		Recommendations: recommendations,
		RedFlags:        redFlags,
		Confidence:      0.8,
	}, nil
}

// AssessRisk assesses risk for a completion request
func (c *AnthropicClient) AssessRisk(ctx context.Context, request CompletionRequest) (*SecurityAssessment, error) {
	var allContent strings.Builder
	for _, msg := range request.Messages {
		allContent.WriteString(msg.GetContentAsString())
		allContent.WriteString(" ")
	}

	return c.AnalyzeFraud(ctx, allContent.String())
}

// GetSecurityRecommendations provides security recommendations
func (c *AnthropicClient) GetSecurityRecommendations(ctx context.Context, content string) ([]string, error) {
	recommendations := []string{
		"Implement input validation",
		"Monitor for unusual patterns",
		"Use rate limiting",
		"Log all interactions",
	}

	lowerContent := strings.ToLower(content)
	if strings.Contains(lowerContent, "personal") {
		recommendations = append(recommendations, "Handle personal data with extra care")
	}
	if strings.Contains(lowerContent, "financial") {
		recommendations = append(recommendations, "Apply enhanced financial security measures")
	}

	return recommendations, nil
}

// GetProviderInfo returns information about this provider
func (c *AnthropicClient) GetProviderInfo() ProviderInfo {
	models := []ModelInfo{
		{
			ID:              ModelClaudeOpus4,
			Name:            "Claude Opus 4",
			Provider:        ProviderAnthropic,
			Version:         "2025-05-14",
			MaxTokens:       32000,
			ContextLength:   200000,
			MaxOutputTokens: 32000,
			InputCost:       0.015, // $15 per 1M tokens
			OutputCost:      0.075, // $75 per 1M tokens
			Capabilities:    []string{"text", "vision", "function_calling", "streaming", "extended_thinking", "computer_use"},
			Description:     "Most capable and intelligent Claude model",
			TrainingCutoff:  "2025-03",
		},
		{
			ID:              ModelClaudeSonnet4,
			Name:            "Claude Sonnet 4",
			Provider:        ProviderAnthropic,
			Version:         "2025-05-14",
			MaxTokens:       64000,
			ContextLength:   200000,
			MaxOutputTokens: 64000,
			InputCost:       0.003, // $3 per 1M tokens
			OutputCost:      0.015, // $15 per 1M tokens
			Capabilities:    []string{"text", "vision", "function_calling", "streaming", "extended_thinking", "computer_use"},
			Description:     "High-performance model with exceptional reasoning",
			TrainingCutoff:  "2025-03",
		},
		{
			ID:              ModelClaude37Sonnet20250219,
			Name:            "Claude Sonnet 3.7",
			Provider:        ProviderAnthropic,
			Version:         "2025-02-19",
			MaxTokens:       64000,
			ContextLength:   200000,
			MaxOutputTokens: 64000,
			InputCost:       0.003, // $3 per 1M tokens
			OutputCost:      0.015, // $15 per 1M tokens
			Capabilities:    []string{"text", "vision", "function_calling", "streaming", "extended_thinking"},
			Description:     "High-performance model with early extended thinking",
			TrainingCutoff:  "2024-11",
		},
		{
			ID:              ModelClaude35Sonnet20241022,
			Name:            "Claude Sonnet 3.5",
			Provider:        ProviderAnthropic,
			Version:         "2024-10-22",
			MaxTokens:       8192,
			ContextLength:   200000,
			MaxOutputTokens: 8192,
			InputCost:       0.003, // $3 per 1M tokens
			OutputCost:      0.015, // $15 per 1M tokens
			Capabilities:    []string{"text", "vision", "function_calling", "streaming", "computer_use"},
			Description:     "Previous intelligent model",
			TrainingCutoff:  "2024-04",
		},
		{
			ID:              ModelClaude35Haiku20241022,
			Name:            "Claude Haiku 3.5",
			Provider:        ProviderAnthropic,
			Version:         "2024-10-22",
			MaxTokens:       8192,
			ContextLength:   200000,
			MaxOutputTokens: 8192,
			InputCost:       0.0008, // $0.80 per 1M tokens
			OutputCost:      0.004,  // $4 per 1M tokens
			Capabilities:    []string{"text", "vision", "function_calling", "streaming"},
			Description:     "Fastest Claude model",
			TrainingCutoff:  "2024-07",
		},
	}

	return ProviderInfo{
		Name:         "Anthropic",
		Provider:     ProviderAnthropic,
		Version:      "2025.1",
		Capabilities: c.GetCapabilities(),
		Models:       models,
		Status:       "active",
	}
}

// ==================== HELPER METHODS ====================

func (c *AnthropicClient) getModel(requestModel string) string {
	if requestModel != "" {
		return requestModel
	}
	if c.defaultModel != "" {
		return c.defaultModel
	}
	return ModelClaudeSonnet4 // Default to Claude Sonnet 4
}

func (c *AnthropicClient) getAnthropicModel(requestModel string) anthropic.Model {
	model := c.getModel(requestModel)

	switch model {
	case ModelClaudeOpus4, ModelClaudeOpus4Latest:
		return anthropic.ModelClaude4Opus20250514
	case ModelClaudeSonnet4, ModelClaudeSonnet4Latest:
		return anthropic.ModelClaude4Sonnet20250514
	case ModelClaude37Sonnet20250219, ModelClaude37SonnetLatest:
		return anthropic.ModelClaude3_7Sonnet20250219
	case ModelClaude35Sonnet20241022, ModelClaude35SonnetLatest:
		return anthropic.ModelClaude3_5Sonnet20241022
	case ModelClaude35Haiku20241022, ModelClaude35HaikuLatest:
		return anthropic.ModelClaude3_5Haiku20241022
	case ModelClaude3Opus20240229:
		return anthropic.ModelClaude_3_Opus_20240229
	default:
		return anthropic.ModelClaude4Sonnet20250514 // Default fallback
	}
}

func (c *AnthropicClient) convertMessages(messages []Message) ([]anthropic.MessageParam, string) {
	var anthropicMessages []anthropic.MessageParam
	var systemMessage string

	for _, msg := range messages {
		if msg.Role == RoleSystem {
			systemMessage = msg.GetContentAsString()
			continue
		}

		// Convert role
		var role anthropic.MessageParamRole
		switch msg.Role {
		case RoleUser:
			role = anthropic.MessageParamRoleUser
		case RoleAssistant:
			role = anthropic.MessageParamRoleAssistant
		default:
			role = anthropic.MessageParamRoleUser
		}

		// Handle different content types
		content := c.convertMessageContent(msg)

		anthropicMsg := anthropic.MessageParam{
			Role:    role,
			Content: content,
		}

		anthropicMessages = append(anthropicMessages, anthropicMsg)
	}

	return anthropicMessages, systemMessage
}

func (c *AnthropicClient) convertMessageContent(msg Message) []anthropic.ContentBlockParamUnion {
	var contentBlocks []anthropic.ContentBlockParamUnion

	// Handle simple text content
	if textContent := msg.GetContentAsString(); textContent != "" {
		contentBlocks = append(contentBlocks, anthropic.NewTextBlock(textContent))
		return contentBlocks
	}

	// Handle multimodal content - check if msg.Content is a slice
	if contentSlice, ok := msg.Content.([]interface{}); ok {
		for _, content := range contentSlice {
			switch c := content.(type) {
			case map[string]interface{}:
				if contentType, ok := c["type"].(string); ok {
					switch contentType {
					case "text":
						if text, ok := c["text"].(string); ok {
							contentBlocks = append(contentBlocks, anthropic.NewTextBlock(text))
						}
					case "image_url":
						if imageURL, ok := c["image_url"].(map[string]interface{}); ok {
							if url, ok := imageURL["url"].(string); ok {
								// Handle base64 images
								if strings.HasPrefix(url, "data:image/") {
									parts := strings.Split(url, ",")
									if len(parts) == 2 {
										mediaType := "image/jpeg"
										if strings.Contains(parts[0], "image/png") {
											mediaType = "image/png"
										}
										contentBlocks = append(contentBlocks, anthropic.NewImageBlockBase64(mediaType, parts[1]))
									}
								}
							}
						}
					}
				}
			case string:
				contentBlocks = append(contentBlocks, anthropic.NewTextBlock(c))
			}
		}
	}

	// Fallback to empty text if no content
	if len(contentBlocks) == 0 {
		contentBlocks = append(contentBlocks, anthropic.NewTextBlock(""))
	}

	return contentBlocks
}

func (c *AnthropicClient) convertTools(tools []Tool) []anthropic.ToolUnionParam {
	var anthropicTools []anthropic.ToolUnionParam

	for _, tool := range tools {
		// Extract properties and required from the parameters map
		var properties interface{}
		var required []string

		if tool.Function.Parameters != nil {
			if props, ok := tool.Function.Parameters["properties"]; ok {
				properties = props
			}
			if req, ok := tool.Function.Parameters["required"]; ok {
				if reqSlice, ok := req.([]interface{}); ok {
					for _, r := range reqSlice {
						if reqStr, ok := r.(string); ok {
							required = append(required, reqStr)
						}
					}
				}
			}
		}

		toolParam := anthropic.ToolUnionParamOfTool(
			anthropic.ToolInputSchemaParam{
				Type:       "object",
				Properties: properties,
				Required:   required,
			},
			tool.Function.Name,
		)
		anthropicTools = append(anthropicTools, toolParam)
	}

	return anthropicTools
}

func (c *AnthropicClient) convertResponse(response *anthropic.Message, requestModel string) *CompletionResponse {
	// Extract content
	var content, thinking string
	var toolCalls []ToolCall

	for _, block := range response.Content {
		switch block.Type {
		case "text":
			content = block.Text
		case "thinking":
			thinking = block.Thinking
		case "tool_use":
			toolCall := ToolCall{
				ID:   block.ID,
				Type: "function",
				Function: FunctionCall{
					Name:      block.Name,
					Arguments: fmt.Sprintf("%v", block.Input),
				},
			}
			toolCalls = append(toolCalls, toolCall)
		}
	}

	// Convert stop reason
	finishReason := c.convertStopReason(&response.StopReason)

	// Extract usage
	usage := Usage{
		PromptTokens:     int(response.Usage.InputTokens),
		CompletionTokens: int(response.Usage.OutputTokens),
		TotalTokens:      int(response.Usage.InputTokens + response.Usage.OutputTokens),
	}

	// Handle cache tokens if available (these are int64, not pointers)
	if response.Usage.CacheCreationInputTokens > 0 {
		usage.CacheCreationInputTokens = int(response.Usage.CacheCreationInputTokens)
	}
	if response.Usage.CacheReadInputTokens > 0 {
		usage.CacheReadInputTokens = int(response.Usage.CacheReadInputTokens)
	}

	return &CompletionResponse{
		ID:       response.ID,
		Object:   "chat.completion",
		Created:  time.Now().Unix(),
		Model:    c.getModel(requestModel),
		Provider: ProviderAnthropic,
		Choices: []Choice{
			{
				Index: 0,
				Message: Message{
					Role:      RoleAssistant,
					Content:   content,
					Thinking:  thinking,
					ToolCalls: toolCalls,
				},
				FinishReason: finishReason,
			},
		},
		Usage: usage,
	}
}

func (c *AnthropicClient) convertStreamEvent(event anthropic.MessageStreamEventUnion, requestModel string) *CompletionStreamResponse {
	switch event.Type {
	case "content_block_start":
		startEvent := event.AsContentBlockStart()

		// Check if it's a text block
		if startEvent.ContentBlock.Type == "text" {
			return &CompletionStreamResponse{
				ID:      fmt.Sprintf("msg_%d", time.Now().Unix()),
				Object:  "chat.completion.chunk",
				Created: time.Now().Unix(),
				Model:   c.getModel(requestModel),
				Choices: []Choice{
					{
						Index: int(event.Index),
						Delta: Message{
							Role:    RoleAssistant,
							Content: "",
						},
					},
				},
			}
		}

	case "content_block_delta":
		deltaEvent := event.AsContentBlockDelta()

		// Check if it's a text delta
		if deltaEvent.Delta.Type == "text_delta" {
			return &CompletionStreamResponse{
				ID:      fmt.Sprintf("msg_%d", time.Now().Unix()),
				Object:  "chat.completion.chunk",
				Created: time.Now().Unix(),
				Model:   c.getModel(requestModel),
				Choices: []Choice{
					{
						Index: int(event.Index),
						Delta: Message{
							Role:    RoleAssistant,
							Content: deltaEvent.Delta.Text,
						},
					},
				},
			}
		}

	case "message_stop":
		return &CompletionStreamResponse{
			ID:      fmt.Sprintf("msg_%d", time.Now().Unix()),
			Object:  "chat.completion.chunk",
			Created: time.Now().Unix(),
			Model:   c.getModel(requestModel),
			Choices: []Choice{
				{
					Index:        0,
					FinishReason: "stop",
					Delta: Message{
						Role: RoleAssistant,
					},
				},
			},
		}
	}

	return nil
}

func (c *AnthropicClient) convertStopReason(stopReason *anthropic.StopReason) string {
	if stopReason == nil {
		return "unknown"
	}

	switch *stopReason {
	case "end_turn":
		return "stop"
	case "max_tokens":
		return "length"
	case "stop_sequence":
		return "stop"
	case "tool_use":
		return "tool_calls"
	default:
		return "unknown"
	}
}

func (c *AnthropicClient) schemaToString(schema *JSONSchemaDefinition) string {
	if schema == nil {
		return ""
	}

	schemaMap := map[string]interface{}{
		"type":        "object",
		"properties":  schema.Schema,
		"description": schema.Description,
	}

	// Convert to JSON string for debugging purposes
	jsonBytes := fmt.Sprintf("%+v", schemaMap)
	return jsonBytes
}

func (c *AnthropicClient) calculateInputCost(model string, tokens int) float64 {
	// Updated pricing per 1K tokens (as of 2025)
	costPer1K := map[string]float64{
		ModelClaudeOpus4:            0.015,  // $15 per 1M tokens
		ModelClaudeSonnet4:          0.003,  // $3 per 1M tokens
		ModelClaude37Sonnet20250219: 0.003,  // $3 per 1M tokens
		ModelClaude35Sonnet20241022: 0.003,  // $3 per 1M tokens
		ModelClaude35Haiku20241022:  0.0008, // $0.80 per 1M tokens
		ModelClaude3Opus20240229:    0.015,  // $15 per 1M tokens
		ModelClaudeOpus4Latest:      0.015,  // Same as Opus 4
		ModelClaudeSonnet4Latest:    0.003,  // Same as Sonnet 4
		ModelClaude37SonnetLatest:   0.003,  // Same as 3.7 Sonnet
		ModelClaude35SonnetLatest:   0.003,  // Same as 3.5 Sonnet
		ModelClaude35HaikuLatest:    0.0008, // Same as 3.5 Haiku
	}

	if cost, exists := costPer1K[model]; exists {
		return float64(tokens) / 1000.0 * cost
	}
	return float64(tokens) / 1000.0 * 0.003 // Default pricing
}

func (c *AnthropicClient) calculateOutputCost(model string, tokens int) float64 {
	// Updated pricing per 1K tokens (as of 2025)
	costPer1K := map[string]float64{
		ModelClaudeOpus4:            0.075, // $75 per 1M tokens
		ModelClaudeSonnet4:          0.015, // $15 per 1M tokens
		ModelClaude37Sonnet20250219: 0.015, // $15 per 1M tokens
		ModelClaude35Sonnet20241022: 0.015, // $15 per 1M tokens
		ModelClaude35Haiku20241022:  0.004, // $4 per 1M tokens
		ModelClaude3Opus20240229:    0.075, // $75 per 1M tokens
		ModelClaudeOpus4Latest:      0.075, // Same as Opus 4
		ModelClaudeSonnet4Latest:    0.015, // Same as Sonnet 4
		ModelClaude37SonnetLatest:   0.015, // Same as 3.7 Sonnet
		ModelClaude35SonnetLatest:   0.015, // Same as 3.5 Sonnet
		ModelClaude35HaikuLatest:    0.004, // Same as 3.5 Haiku
	}

	if cost, exists := costPer1K[model]; exists {
		return float64(tokens) / 1000.0 * cost
	}
	return float64(tokens) / 1000.0 * 0.015 // Default pricing
}

func (c *AnthropicClient) updateStats(response *CompletionResponse, latency time.Duration) {
	c.stats.RequestCount++
	c.stats.TokensUsed += int64(response.Usage.TotalTokens)
	c.stats.TotalCost += response.Usage.TotalCost
	c.stats.LastUsed = time.Now()

	// Update average latency
	if c.stats.AverageLatency == 0 {
		c.stats.AverageLatency = latency
	} else {
		c.stats.AverageLatency = (c.stats.AverageLatency + latency) / 2
	}
}
