package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	ai2 "middleman/internal/ai"
)

// ToolExecutor provides simple tool execution
type ToolExecutor struct {
	registry *ToolRegistry
	mu       sync.RWMutex
	metrics  map[string]*SimpleToolMetrics
}

// SimpleToolMetrics tracks tool execution metrics for the simple executor
type SimpleToolMetrics struct {
	ExecutionCount int64
	SuccessCount   int64
	ErrorCount     int64
	TotalDuration  int64 // in milliseconds
}

// NewToolExecutor creates a new tool executor
func NewToolExecutor(registry *ToolRegistry) *ToolExecutor {
	return &ToolExecutor{
		registry: registry,
		metrics:  make(map[string]*SimpleToolMetrics),
	}
}

// ExecuteToolCall executes a single tool call
func (e *ToolExecutor) ExecuteToolCall(ctx context.Context, toolCall ai2.ToolCall) (*ToolResult, error) {
	log.Printf("[TOOL_EXECUTOR] Executing tool: %s (ID: %s)", toolCall.Function.Name, toolCall.ID)

	// Parse arguments
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &params); err != nil {
		e.recordError(toolCall.Function.Name)
		return &ToolResult{
			ToolCallID: toolCall.ID,
			Success:    false,
			Error:      fmt.Sprintf("Failed to parse arguments: %v", err),
		}, err
	}

	// Execute tool
	startTime := currentTimeMillis()
	result, err := e.registry.ExecuteTool(ctx, toolCall.Function.Name, params)
	duration := currentTimeMillis() - startTime

	if err != nil {
		e.recordError(toolCall.Function.Name)
		log.Printf("[TOOL_EXECUTOR] Tool %s failed: %v", toolCall.Function.Name, err)
		return &ToolResult{
			ToolCallID: toolCall.ID,
			ToolName:   toolCall.Function.Name,
			Success:    false,
			Error:      err.Error(),
			Duration:   duration,
		}, err
	}

	// Record success
	e.recordSuccess(toolCall.Function.Name, duration)

	// Marshal result
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return &ToolResult{
			ToolCallID: toolCall.ID,
			ToolName:   toolCall.Function.Name,
			Success:    false,
			Error:      fmt.Sprintf("Failed to marshal result: %v", err),
			Duration:   duration,
		}, err
	}

	log.Printf("[TOOL_EXECUTOR] Tool %s completed successfully in %dms", toolCall.Function.Name, duration)

	return &ToolResult{
		ToolCallID: toolCall.ID,
		ToolName:   toolCall.Function.Name,
		Success:    true,
		Result:     result,
		ResultJSON: string(resultJSON),
		Duration:   duration,
	}, nil
}

// ExecuteParallel executes multiple tool calls in parallel
func (e *ToolExecutor) ExecuteParallel(ctx context.Context, toolCalls []ai2.ToolCall) []*ToolResult {
	results := make([]*ToolResult, len(toolCalls))
	wg := sync.WaitGroup{}

	log.Printf("[TOOL_EXECUTOR] Executing %d tools in parallel", len(toolCalls))

	for i, toolCall := range toolCalls {
		wg.Add(1)
		go func(index int, tc ai2.ToolCall) {
			defer wg.Done()
			result, _ := e.ExecuteToolCall(ctx, tc)
			results[index] = result
		}(i, toolCall)
	}

	wg.Wait()

	// Log summary
	successCount := 0
	for _, result := range results {
		if result.Success {
			successCount++
		}
	}
	log.Printf("[TOOL_EXECUTOR] Parallel execution complete: %d/%d succeeded", successCount, len(toolCalls))

	return results
}

// ExecuteSequential executes tool calls one by one
func (e *ToolExecutor) ExecuteSequential(ctx context.Context, toolCalls []ai2.ToolCall) ([]*ToolResult, error) {
	results := make([]*ToolResult, 0, len(toolCalls))

	for _, toolCall := range toolCalls {
		result, err := e.ExecuteToolCall(ctx, toolCall)
		results = append(results, result)

		// Stop on first error if critical
		if err != nil && isCriticalError(err) {
			return results, fmt.Errorf("critical error in tool %s: %w", toolCall.Function.Name, err)
		}
	}

	return results, nil
}

// GetMetrics returns execution metrics for all tools
func (e *ToolExecutor) GetMetrics() map[string]*SimpleToolMetrics {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// Create a copy to avoid race conditions
	metricsCopy := make(map[string]*SimpleToolMetrics)
	for name, metrics := range e.metrics {
		metricsCopy[name] = &SimpleToolMetrics{
			ExecutionCount: metrics.ExecutionCount,
			SuccessCount:   metrics.SuccessCount,
			ErrorCount:     metrics.ErrorCount,
			TotalDuration:  metrics.TotalDuration,
		}
	}

	return metricsCopy
}

// GetToolDefinitions returns all available tool definitions
func (e *ToolExecutor) GetToolDefinitions() []ai2.Tool {
	return e.registry.GetToolDefinitions()
}

// recordSuccess records a successful tool execution
func (e *ToolExecutor) recordSuccess(toolName string, duration int64) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.metrics[toolName]; !exists {
		e.metrics[toolName] = &SimpleToolMetrics{}
	}

	e.metrics[toolName].ExecutionCount++
	e.metrics[toolName].SuccessCount++
	e.metrics[toolName].TotalDuration += duration
}

// recordError records a failed tool execution
func (e *ToolExecutor) recordError(toolName string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.metrics[toolName]; !exists {
		e.metrics[toolName] = &SimpleToolMetrics{}
	}

	e.metrics[toolName].ExecutionCount++
	e.metrics[toolName].ErrorCount++
}

// ToolResult represents the result of a tool execution
type ToolResult struct {
	ToolCallID string      `json:"tool_call_id"`
	ToolName   string      `json:"tool_name"`
	Success    bool        `json:"success"`
	Result     interface{} `json:"result,omitempty"`
	ResultJSON string      `json:"result_json,omitempty"`
	Error      string      `json:"error,omitempty"`
	Duration   int64       `json:"duration_ms"`
}

// ToJSON converts the result to JSON string for LLM consumption
func (r *ToolResult) ToJSON() string {
	if r.Success {
		return r.ResultJSON
	}
	return fmt.Sprintf(`{"error": "%s"}`, r.Error)
}

// currentTimeMillis returns the current time in milliseconds
func currentTimeMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func isCriticalError(err error) bool {
	// Define which errors should stop sequential execution
	// For now, no errors are considered critical
	return false
}
