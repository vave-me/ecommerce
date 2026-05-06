package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"sync"
	"time"

	"middleman/managers/internal/application/services"
	"middleman/managers/internal/application/tools"
	"middleman/managers/internal/domain"
	ai2 "middleman/internal/ai"
)

// LLMProcessor handles LLM interactions with simplified tool execution
type LLMProcessor struct {
	aiClient       ai2.EnhancedAIService
	clientProvider services.AIClientProvider
	toolRegistry   *tools.ToolRegistry

	// Conversation context
	conversationHistory []ai2.Message
	mu                  sync.RWMutex
}

// NewLLMProcessor creates a new simplified LLM processor
func NewLLMProcessor(
	aiClient ai2.EnhancedAIService,
	clientProvider services.AIClientProvider,
	toolRegistry *tools.ToolRegistry,
) *LLMProcessor {
	return &LLMProcessor{
		aiClient:            aiClient,
		clientProvider:      clientProvider,
		toolRegistry:        toolRegistry,
		conversationHistory: make([]ai2.Message, 0),
	}
}

// ProcessRequest processes a user request with tool support
func (p *LLMProcessor) ProcessRequest(ctx context.Context, request string, assistant *domain.Assistant) (*ProcessResult, error) {
	if request == "" {
		return nil, fmt.Errorf("request cannot be empty")
	}
	if assistant == nil {
		return nil, fmt.Errorf("assistant cannot be nil")
	}

	startTime := time.Now()

	// Add user message to history
	p.addToHistory(ai2.Message{
		Role:    "user",
		Content: request,
	})

	// Prepare messages with system prompt
	messages := []ai2.Message{
		{
			Role:    "system",
			Content: assistant.SystemPrompt,
		},
	}

	// Add conversation history
	messages = append(messages, p.getRecentHistory(10)...)

	// Create completion request with tools
	temp := assistant.Temperature
	maxTok := int(assistant.MaxTokens)
	if maxTok <= 0 {
		maxTok = 2048 // Default to 2048 tokens if not set or invalid
	}
	tools := p.toolRegistry.GetToolDefinitions()

	// For assistants with tools, encourage tool usage
	toolChoice := "auto"
	if len(tools) > 0 && assistant.Type != "scheduler" {
		// Use "required" to force tool usage
		toolChoice = "required"
	}

	completionReq := &ai2.CompletionRequest{
		Messages:    messages,
		Temperature: &temp,
		MaxTokens:   &maxTok,
		Tools:       tools,
		ToolChoice:  toolChoice,
	}

	// Log before API call
	log.Printf("[LLM_PROCESSOR] Calling AI with %d messages, %d tools, tool_choice: %s", len(messages), len(tools), toolChoice)

	// Get completion from AI
	response, err := p.aiClient.CreateCompletion(ctx, *completionReq)
	if err != nil {
		log.Printf("[LLM_PROCESSOR] AI call failed: %v", err)
		return nil, fmt.Errorf("failed to get AI response: %w", err)
	}

	log.Printf("[LLM_PROCESSOR] AI response received with %d choices", len(response.Choices))

	// Extract message from response
	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no choices in AI response")
	}

	message := response.Choices[0].Message

	// Process tool calls if any
	if len(message.ToolCalls) > 0 {
		return p.processWithTools(ctx, response, assistant, startTime)
	}

	// Add assistant response to history
	p.addToHistory(message)

	responseText := message.GetContentAsString()

	return &ProcessResult{
		Response:  responseText,
		ToolsUsed: 0,
		Success:   true,
		Duration:  time.Since(startTime),
	}, nil
}

// processWithTools handles responses that require tool execution
func (p *LLMProcessor) processWithTools(ctx context.Context, response *ai2.CompletionResponse, assistant *domain.Assistant, startTime time.Time) (*ProcessResult, error) {
	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no choices in AI response")
	}

	message := response.Choices[0].Message

	// Execute all tool calls in parallel
	toolResults := p.executeToolsParallel(ctx, message.ToolCalls)

	// Add assistant message with tool calls to history
	p.addToHistory(message)

	// Add tool results to history
	for _, result := range toolResults {
		p.addToHistory(ai2.Message{
			Role:       "tool",
			Content:    result.Content,
			ToolCallID: result.ToolCallID,
		})
	}

	// Get final response from AI with tool results
	messages := []ai2.Message{
		{
			Role:    "system",
			Content: assistant.SystemPrompt,
		},
	}
	messages = append(messages, p.getRecentHistory(20)...)

	// Create pointers for temperature and max tokens
	temp := assistant.Temperature
	maxTok := int(assistant.MaxTokens)
	if maxTok <= 0 {
		maxTok = 2048 // Default to 2048 tokens if not set or invalid
	}

	finalReq := &ai2.CompletionRequest{
		Messages:    messages,
		Temperature: &temp,
		MaxTokens:   &maxTok,
		Tools:       p.toolRegistry.GetToolDefinitions(),
	}

	finalResponse, err := p.aiClient.CreateCompletion(ctx, *finalReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get final response: %w", err)
	}

	// Extract message from response
	if len(finalResponse.Choices) == 0 {
		return nil, fmt.Errorf("no choices in final AI response")
	}

	finalMessage := finalResponse.Choices[0].Message

	// If more tools are requested, process them recursively (with depth limit)
	if len(finalMessage.ToolCalls) > 0 {
		return p.processWithTools(ctx, finalResponse, assistant, startTime)
	}

	// Add final response to history
	p.addToHistory(finalMessage)

	// Store tool call information in metadata
	toolCallInfo := make([]interface{}, 0, len(message.ToolCalls))
	for i, tc := range message.ToolCalls {
		var params map[string]interface{}
		json.Unmarshal([]byte(tc.Function.Arguments), &params)
		toolCallInfo = append(toolCallInfo, map[string]interface{}{
			"name":       tc.Function.Name,
			"parameters": params,
			"result":     toolResults[i].Content,
		})
	}

	return &ProcessResult{
		Response:  finalMessage.GetContentAsString(),
		ToolsUsed: len(message.ToolCalls),
		Success:   true,
		Duration:  time.Since(startTime),
		Metadata: map[string]interface{}{
			"tool_calls": toolCallInfo,
		},
	}, nil
}

// executeToolsParallel executes multiple tools in parallel
func (p *LLMProcessor) executeToolsParallel(ctx context.Context, toolCalls []ai2.ToolCall) []ToolResult {
	results := make([]ToolResult, len(toolCalls))
	wg := sync.WaitGroup{}

	for i, toolCall := range toolCalls {
		wg.Add(1)
		go func(index int, tc ai2.ToolCall) {
			defer wg.Done()

			// Parse arguments
			var params map[string]interface{}
			if err := json.Unmarshal([]byte(tc.Function.Arguments), &params); err != nil {
				results[index] = ToolResult{
					ToolCallID: tc.ID,
					Content:    fmt.Sprintf(`{"error": "Failed to parse arguments: %v"}`, err),
				}
				return
			}

			// Execute tool
			result, err := p.toolRegistry.ExecuteTool(ctx, tc.Function.Name, params)
			if err != nil {
				results[index] = ToolResult{
					ToolCallID: tc.ID,
					Content:    fmt.Sprintf(`{"error": "%v"}`, err),
				}
				return
			}

			// Marshal result
			resultJSON, err := json.Marshal(result)
			if err != nil {
				results[index] = ToolResult{
					ToolCallID: tc.ID,
					Content:    fmt.Sprintf(`{"error": "Failed to marshal result: %v"}`, err),
				}
				return
			}

			results[index] = ToolResult{
				ToolCallID: tc.ID,
				Content:    string(resultJSON),
			}
		}(i, toolCall)
	}

	wg.Wait()
	return results
}

// ProcessWithHistory processes a request with conversation history
func (p *LLMProcessor) ProcessWithHistory(ctx context.Context, assistant *domain.Assistant, currentMessage string, history []domain.ConversationMessage, contextData map[string]interface{}) (string, []domain.AssistantAction, float64, error) {
	// Clear history and rebuild it properly
	p.ClearHistory()

	// Add conversation history
	for _, msg := range history {
		p.addToHistory(ai2.Message{
			Role:    string(msg.Role),
			Content: msg.Content,
		})
	}

	// Process the current message
	result, err := p.ProcessRequest(ctx, currentMessage, assistant)
	if err != nil {
		return "", nil, 0.0, err
	}

	// Convert tool usage into actions
	actions := []domain.AssistantAction{}
	if result.ToolsUsed > 0 {
		// Extract tool calls from metadata if available
		if result.Metadata != nil {
			if toolCalls, ok := result.Metadata["tool_calls"].([]interface{}); ok {
				for _, tc := range toolCalls {
					if toolCall, ok := tc.(map[string]interface{}); ok {
						action := domain.AssistantAction{
							Type:        "tool_call",
							Endpoint:    toolCall["name"].(string),
							Method:      "execute",
							Parameters:  toolCall["parameters"].(map[string]interface{}),
							Description: fmt.Sprintf("Tool call: %s", toolCall["name"].(string)),
						}
						actions = append(actions, action)
					}
				}
			}
		}
	}

	// Calculate confidence based on success and tool usage
	confidence := 0.9
	if !result.Success {
		confidence = 0.5
	} else if result.ToolsUsed > 0 {
		confidence = 0.95 // Higher confidence when tools are used
	}

	return result.Response, actions, confidence, nil
}

// ShouldUseAI determines if AI should be used for the assistant
func (p *LLMProcessor) ShouldUseAI(assistant interface{}) bool {
	// Always use AI for now
	return true
}

// ProcessStreamingRequest processes a request with streaming response
// Note: Streaming is not currently supported by the AI service
func (p *LLMProcessor) ProcessStreamingRequest(ctx context.Context, request string, assistant *domain.Assistant, streamChan chan<- StreamChunk) error {
	defer close(streamChan)

	// For now, we'll use non-streaming and simulate streaming
	result, err := p.ProcessRequest(ctx, request, assistant)
	if err != nil {
		streamChan <- StreamChunk{
			Type:    "error",
			Content: fmt.Sprintf("Error processing request: %v", err),
		}
		return err
	}

	// Send the response as a single chunk
	streamChan <- StreamChunk{
		Type:    "content",
		Content: result.Response,
	}

	streamChan <- StreamChunk{
		Type:    "complete",
		Content: "Processing complete",
	}

	return nil
}

// ClearHistory clears the conversation history
func (p *LLMProcessor) ClearHistory() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.conversationHistory = make([]ai2.Message, 0)
}

// addToHistory adds a message to conversation history
func (p *LLMProcessor) addToHistory(msg ai2.Message) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.conversationHistory = append(p.conversationHistory, msg)

	// Keep only last 100 messages
	if len(p.conversationHistory) > 100 {
		p.conversationHistory = p.conversationHistory[len(p.conversationHistory)-100:]
	}
}

// getRecentHistory returns the most recent messages from history
func (p *LLMProcessor) getRecentHistory(count int) []ai2.Message {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if len(p.conversationHistory) <= count {
		return p.conversationHistory
	}

	return p.conversationHistory[len(p.conversationHistory)-count:]
}

// ProcessResult represents the result of processing a request
type ProcessResult struct {
	Response  string
	ToolsUsed int
	Success   bool
	Duration  time.Duration
	Metadata  map[string]interface{}
}

// ToolResult represents a tool execution result
type ToolResult struct {
	ToolCallID string
	Content    string
}

// StreamChunk represents a streaming response chunk
type StreamChunk struct {
	Content string
	Type    string // "content", "tool_start", "tool_result", "complete", "error"
}
