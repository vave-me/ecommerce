package ai

import (
	"context"
	"fmt"
)

// ReasoningRequest represents a request to reasoning models
type ReasoningRequest struct {
	Model       string               `json:"model"`
	Input       []ReasoningInputItem `json:"input"`
	Tools       []Tool               `json:"tools,omitempty"`
	Reasoning   *ReasoningConfig     `json:"reasoning,omitempty"`
	MaxTokens   *int                 `json:"max_tokens,omitempty"`
	Temperature *float64             `json:"temperature,omitempty"`
	Store       bool                 `json:"store"`
	Include     []string             `json:"include,omitempty"`
}

// ReasoningConfig configures reasoning behavior
type ReasoningConfig struct {
	Effort  string `json:"effort"`  // "low", "medium", "high"
	Summary string `json:"summary"` // "auto", "detailed", "none"
}

// ReasoningInputItem represents an input item for reasoning models
type ReasoningInputItem struct {
	Type             string `json:"type"`
	Role             string `json:"role,omitempty"`
	Content          string `json:"content,omitempty"`
	CallID           string `json:"call_id,omitempty"`
	Output           string `json:"output,omitempty"`
	ID               string `json:"id,omitempty"`
	EncryptedContent string `json:"encrypted_content,omitempty"`
	Arguments        string `json:"arguments,omitempty"`
	Name             string `json:"name,omitempty"`
	Status           string `json:"status,omitempty"`
}

// ReasoningResponse represents a response from reasoning models
type ReasoningResponse struct {
	ID      string                `json:"id"`
	Object  string                `json:"object"`
	Created int64                 `json:"created"`
	Model   string                `json:"model"`
	Output  []ReasoningOutputItem `json:"output"`
	Usage   ReasoningUsage        `json:"usage"`
	Status  string                `json:"status"`
}

// ReasoningOutputItem represents an output item from reasoning models
type ReasoningOutputItem struct {
	ID               string             `json:"id"`
	Type             string             `json:"type"`
	Summary          []ReasoningSummary `json:"summary,omitempty"`
	EncryptedContent string             `json:"encrypted_content,omitempty"`
	Content          []ReasoningContent `json:"content,omitempty"`
	Role             string             `json:"role,omitempty"`
	Status           string             `json:"status,omitempty"`
	Arguments        string             `json:"arguments,omitempty"`
	CallID           string             `json:"call_id,omitempty"`
	Name             string             `json:"name,omitempty"`
}

// ReasoningSummary represents a reasoning summary
type ReasoningSummary struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

// ReasoningContent represents reasoning content
type ReasoningContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ReasoningUsage represents token usage for reasoning models
type ReasoningUsage struct {
	InputTokens         int                    `json:"input_tokens"`
	OutputTokens        int                    `json:"output_tokens"`
	TotalTokens         int                    `json:"total_tokens"`
	InputTokensDetails  *ReasoningTokenDetails `json:"input_tokens_details,omitempty"`
	OutputTokensDetails *ReasoningTokenDetails `json:"output_tokens_details,omitempty"`
}

// ReasoningTokenDetails provides detailed token information
type ReasoningTokenDetails struct {
	CachedTokens    int `json:"cached_tokens,omitempty"`
	ReasoningTokens int `json:"reasoning_tokens,omitempty"`
}

// ReasoningClient extends OpenAIClient with reasoning model support
type ReasoningClient struct {
	*OpenAIClient
	responseAPIBase string
}

// NewReasoningClient creates a client with reasoning model support
func NewReasoningClient(apiKey, baseURL, defaultModel string) (*ReasoningClient, error) {
	baseClient, err := NewOpenAIClient(apiKey, baseURL, defaultModel)
	if err != nil {
		return nil, err
	}

	responseAPIBase := "https://api.openai.com/v1/responses"
	if baseURL != "" && baseURL != "https://api.openai.com/v1" {
		responseAPIBase = baseURL + "/responses"
	}

	return &ReasoningClient{
		OpenAIClient:    baseClient,
		responseAPIBase: responseAPIBase,
	}, nil
}

// CreateReasoningCompletion creates a completion using reasoning models
func (c *ReasoningClient) CreateReasoningCompletion(ctx context.Context, request ReasoningRequest) (*ReasoningResponse, error) {
	// Implement direct HTTP calls to Responses API since it's not in the Go SDK yet
	// This would use the responseAPIBase URL and proper request formatting

	// For now, return an error indicating this needs implementation
	return nil, fmt.Errorf("reasoning models require Responses API implementation")
}

// ExecuteSequentialTools executes tools sequentially for reasoning models
func (c *ReasoningClient) ExecuteSequentialTools(
	ctx context.Context,
	input []ReasoningInputItem,
	tools []Tool,
	maxIterations int,
) (*ReasoningResponse, error) {
	conversation := input
	iteration := 0

	for iteration < maxIterations {
		request := ReasoningRequest{
			Model: c.defaultModel,
			Input: conversation,
			Tools: tools,
			Store: true,
			Reasoning: &ReasoningConfig{
				Effort:  "medium",
				Summary: "auto",
			},
		}

		response, err := c.CreateReasoningCompletion(ctx, request)
		if err != nil {
			return nil, fmt.Errorf("reasoning completion failed at iteration %d: %w", iteration, err)
		}

		// Check if we have function calls to execute
		functionCalls := []ReasoningOutputItem{}
		for _, item := range response.Output {
			if item.Type == "function_call" {
				functionCalls = append(functionCalls, item)
			}
		}

		if len(functionCalls) == 0 {
			// No more function calls, we're done
			return response, nil
		}

		// Add all output items to conversation (including reasoning)
		for _, item := range response.Output {
			conversation = append(conversation, ReasoningInputItem{
				Type:             item.Type,
				ID:               item.ID,
				EncryptedContent: item.EncryptedContent,
				Arguments:        item.Arguments,
				Name:             item.Name,
				CallID:           item.CallID,
				Status:           item.Status,
			})
		}

		// Execute function calls sequentially
		for _, funcCall := range functionCalls {
			// This would integrate with the existing tool execution system
			result := c.executeToolCall(ctx, funcCall)

			conversation = append(conversation, ReasoningInputItem{
				Type:   "function_call_output",
				CallID: funcCall.CallID,
				Output: result,
			})
		}

		iteration++
	}

	return nil, fmt.Errorf("max iterations (%d) reached without completion", maxIterations)
}

// executeToolCall executes a single tool call
func (c *ReasoningClient) executeToolCall(ctx context.Context, funcCall ReasoningOutputItem) string {
	// This would integrate with the existing tool execution system
	// For now, return a placeholder
	return fmt.Sprintf("Tool %s executed with args: %s", funcCall.Name, funcCall.Arguments)
}

// IsReasoningModel checks if the model is a reasoning model
func IsReasoningModel(model string) bool {
	reasoningModels := []string{
		"o1-preview", "o1-mini", "o3", "o4-mini",
		"o3-mini", "gpt-4-1-mini",
	}

	for _, rm := range reasoningModels {
		if model == rm {
			return true
		}
	}
	return false
}

// GetReasoningCapabilities returns capabilities specific to reasoning models
func (c *ReasoningClient) GetReasoningCapabilities() ReasoningCapabilities {
	return ReasoningCapabilities{
		SupportsSequentialTools:  true,
		SupportsReasoningSummary: true,
		SupportsEncryptedContent: true,
		MaxIterations:            10,
		SupportedEffortLevels:    []string{"low", "medium", "high"},
		SupportedSummaryTypes:    []string{"auto", "detailed", "none"},
	}
}

// ReasoningCapabilities defines reasoning model capabilities
type ReasoningCapabilities struct {
	SupportsSequentialTools  bool     `json:"supports_sequential_tools"`
	SupportsReasoningSummary bool     `json:"supports_reasoning_summary"`
	SupportsEncryptedContent bool     `json:"supports_encrypted_content"`
	MaxIterations            int      `json:"max_iterations"`
	SupportedEffortLevels    []string `json:"supported_effort_levels"`
	SupportedSummaryTypes    []string `json:"supported_summary_types"`
}
