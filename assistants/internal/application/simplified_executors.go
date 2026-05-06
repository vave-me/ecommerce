package application

import (
	"context"
	"fmt"
	"middleman/managers/internal/application/tools"
	ai2 "middleman/internal/ai"
	"time"
)

// SimplifiedToolExecutor handles tool execution without complex abstractions
type SimplifiedToolExecutor struct {
	registry     *tools.ToolRegistry
	toolExecutor *tools.ToolExecutor
	config       *ToolConfig
}

// NewSimplifiedToolExecutor creates a new simplified tool executor
func NewSimplifiedToolExecutor(registry *tools.ToolRegistry, config *ToolConfig) *SimplifiedToolExecutor {
	if config == nil {
		config = DefaultToolConfig()
	}
	return &SimplifiedToolExecutor{
		registry:     registry,
		toolExecutor: tools.NewToolExecutor(registry),
		config:       config,
	}
}

// ExecuteTools executes multiple tools and returns results
func (e *SimplifiedToolExecutor) ExecuteTools(ctx context.Context, toolCalls []ai2.ToolCall, execCtx *tools.ToolExecutionContext) ([]tools.ToolExecutionResult, error) {
	// Execute tools in parallel
	toolResults := e.toolExecutor.ExecuteParallel(ctx, toolCalls)

	// Convert to ToolExecutionResult format
	results := make([]tools.ToolExecutionResult, len(toolResults))
	for i, result := range toolResults {
		results[i] = tools.ToolExecutionResult{
			ToolName: result.ToolName,
			Status:   getToolStatus(result.Success, result.Error),
			Result:   result.Result,
			Error:    result.Error,
			Duration: result.Duration,
			Metadata: map[string]interface{}{
				"tool_call_id":      result.ToolCallID,
				"execution_context": execCtx,
			},
		}
	}

	return results, nil
}

// SimplifiedStreamingExecutor handles streaming tool execution
type SimplifiedStreamingExecutor struct {
	registry *tools.ToolRegistry
	config   *StreamingConfig
}

// NewSimplifiedStreamingExecutor creates a new streaming executor
func NewSimplifiedStreamingExecutor(registry *tools.ToolRegistry, config *StreamingConfig) *SimplifiedStreamingExecutor {
	if config == nil {
		config = DefaultStreamingConfig()
	}
	return &SimplifiedStreamingExecutor{
		registry: registry,
		config:   config,
	}
}

// ExecuteToolsStream executes tools with streaming support
func (e *SimplifiedStreamingExecutor) ExecuteToolsStream(ctx context.Context, toolCalls []ai2.ToolCall, execCtx *tools.ToolExecutionContext) (<-chan tools.ToolExecutionStream, error) {
	streamChan := make(chan tools.ToolExecutionStream, e.config.StreamBufferSize)

	go func() {
		defer close(streamChan)
		e.executeWithStreaming(ctx, toolCalls, execCtx, streamChan)
	}()

	return streamChan, nil
}

func (e *SimplifiedStreamingExecutor) executeWithStreaming(ctx context.Context, toolCalls []ai2.ToolCall, execCtx *tools.ToolExecutionContext, streamChan chan<- tools.ToolExecutionStream) {
	toolExecutor := tools.NewToolExecutor(e.registry)

	for i, toolCall := range toolCalls {
		// Send start event
		select {
		case streamChan <- tools.ToolExecutionStream{
			ID:       fmt.Sprintf("%s-%d", execCtx.RequestID, i),
			ToolName: toolCall.Function.Name,
			Status:   "started",
			Progress: 0,
			Metadata: map[string]interface{}{
				"call_index":  i,
				"total_tools": len(toolCalls),
			},
			Timestamp: time.Now().Unix(),
		}:
		case <-ctx.Done():
			return
		}

		// Execute tool
		result, err := toolExecutor.ExecuteToolCall(ctx, toolCall)

		// Send completion event
		status := "completed"
		errorMsg := ""
		if err != nil {
			status = "error"
			errorMsg = err.Error()
		}

		select {
		case streamChan <- tools.ToolExecutionStream{
			ID:       fmt.Sprintf("%s-%d", execCtx.RequestID, i),
			ToolName: toolCall.Function.Name,
			Status:   status,
			Progress: 100,
			Result:   result.Result,
			Error:    errorMsg,
			Metadata: map[string]interface{}{
				"call_index":  i,
				"total_tools": len(toolCalls),
				"duration_ms": result.Duration,
			},
			Timestamp: time.Now().Unix(),
		}:
		case <-ctx.Done():
			return
		}
	}
}

// Helper function to get tool status
func getToolStatus(success bool, errorMsg string) string {
	if success {
		return "completed"
	}
	if errorMsg != "" {
		return "error"
	}
	return "failed"
}
