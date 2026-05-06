package tools

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"

	"sync"
	"time"
)

// ToolHandler is a function that handles a specific tool execution
type ToolHandler func(ctx context.Context, r *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error)

// ComprehensiveToolRegistry implements ALL repository methods as tools
type ComprehensiveToolRegistry struct {
	*ToolRegistry
	handlers map[string]ToolHandler
	mu       sync.RWMutex
	metrics  *ToolMetrics
}

// ToolMetrics tracks tool execution metrics
type ToolMetrics struct {
	mu         sync.RWMutex
	executions map[string]*ExecutionMetrics
}

// ExecutionMetrics tracks metrics for a specific tool
type ExecutionMetrics struct {
	TotalCalls      int64
	SuccessfulCalls int64
	FailedCalls     int64
	TotalDuration   time.Duration
	LastExecution   time.Time
}

// NewComprehensiveToolRegistry creates a registry with all repository methods exposed
func NewComprehensiveToolRegistry(registry *ToolRegistry) *ComprehensiveToolRegistry {
	r := &ComprehensiveToolRegistry{
		ToolRegistry: registry,
		handlers:     make(map[string]ToolHandler),
		metrics: &ToolMetrics{
			executions: make(map[string]*ExecutionMetrics),
		},
	}

	// Initialize all tool handlers
	r.initializeActivityHandlers()
	r.initializeBasketHandlers()
	r.initializeCategoryHandlers()
	r.initializeCommentHandlers()
	r.initializeFollowingHandlers()
	r.initializeGeocodingHandlers()
	r.initializeMailerHandlers()
	r.initializeMediaHandlers()
	r.initializeMessageHandlers()
	r.initializeMetricHandlers()
	r.initializeNewsletterHandlers()
	r.initializeNotificationHandlers()
	r.initializeOfferHandlers()
	r.initializeOrderHandlers()
	r.initializePaymentHandlers()
	r.initializePostHandlers()
	r.initializeProductHandlers()
	r.initializeReviewHandlers()
	r.initializeServiceHandlers()
	r.initializeShippingHandlers()
	r.initializeSupportHandlers()
	r.initializeUserHandlers()
	r.initializeVectorHandlers()
	r.initializeWishlistHandlers()

	return r
}

// MainExecute is the main entry point for tool execution with comprehensive error handling
func (r *ComprehensiveToolRegistry) MainExecute(ctx context.Context, toolName string, params map[string]interface{}) (result interface{}, err error) {
	startTime := time.Now()

	// Update metrics
	r.updateMetrics(toolName, func(m *ExecutionMetrics) {
		m.TotalCalls++
		m.LastExecution = startTime
	})

	// Defer panic recovery and metrics update
	defer func() {
		duration := time.Since(startTime)

		if panicVal := recover(); panicVal != nil {
			err = fmt.Errorf("panic in tool %s: %v", toolName, panicVal)

			// Update failure metrics
			r.updateMetrics(toolName, func(m *ExecutionMetrics) {
				m.FailedCalls++
				m.TotalDuration += duration
			})
		} else if err != nil {
			// Update failure metrics
			r.updateMetrics(toolName, func(m *ExecutionMetrics) {
				m.FailedCalls++
				m.TotalDuration += duration
			})
		} else {
			// Update success metrics
			r.updateMetrics(toolName, func(m *ExecutionMetrics) {
				m.SuccessfulCalls++
				m.TotalDuration += duration
			})
		}
	}()

	// Validate context
	if ctx == nil {
		return nil, fmt.Errorf("context is nil for tool %s", toolName)
	}

	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context cancelled before executing tool %s: %w", toolName, ctx.Err())
	default:
	}

	// Validate parameters
	if params == nil {
		params = make(map[string]interface{})
	}

	// Get handler
	r.mu.RLock()
	handler, exists := r.handlers[toolName]
	r.mu.RUnlock()

	if !exists {
		availableTools := r.GetAvailableTools()
		return nil, fmt.Errorf("tool '%s' not found in registry. Available tools: %v", toolName, availableTools[:min(10, len(availableTools))])
	}

	// Execute handler
	result, err = handler(ctx, r, params)

	// Check context again after execution
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context cancelled after executing tool %s: %w", toolName, ctx.Err())
	default:
	}

	return result, err
}

// GetAvailableTools returns a sorted list of all available tool names
func (r *ComprehensiveToolRegistry) GetAvailableTools() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tools := make([]string, 0, len(r.handlers))
	for name := range r.handlers {
		tools = append(tools, name)
	}
	return tools
}

// HasTool checks if a tool exists in the registry
func (r *ComprehensiveToolRegistry) HasTool(toolName string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.handlers[toolName]
	return exists
}

// GetToolMetrics returns metrics for a specific tool
func (r *ComprehensiveToolRegistry) GetToolMetrics(toolName string) *ExecutionMetrics {
	r.metrics.mu.RLock()
	defer r.metrics.mu.RUnlock()

	if metrics, exists := r.metrics.executions[toolName]; exists {
		// Return a copy to prevent modification
		copy := *metrics
		return &copy
	}
	return nil
}

// GetAllMetrics returns metrics for all tools
func (r *ComprehensiveToolRegistry) GetAllMetrics() map[string]ExecutionMetrics {
	r.metrics.mu.RLock()
	defer r.metrics.mu.RUnlock()

	result := make(map[string]ExecutionMetrics)
	for name, metrics := range r.metrics.executions {
		result[name] = *metrics
	}
	return result
}

// RegisterTool allows dynamic registration of new tools
func (r *ComprehensiveToolRegistry) RegisterTool(name string, handler ToolHandler) error {
	if name == "" {
		return fmt.Errorf("tool name cannot be empty")
	}
	if handler == nil {
		return fmt.Errorf("tool handler cannot be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.handlers[name]; exists {
		return fmt.Errorf("tool '%s' already registered", name)
	}

	r.handlers[name] = handler
	log.Printf("[COMPREHENSIVE_TOOL_REGISTRY] Registered new tool: %s", name)
	return nil
}

// UnregisterTool removes a tool from the registry
func (r *ComprehensiveToolRegistry) UnregisterTool(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.handlers[name]; !exists {
		return fmt.Errorf("tool '%s' not found", name)
	}

	delete(r.handlers, name)
	log.Printf("[COMPREHENSIVE_TOOL_REGISTRY] Unregistered tool: %s", name)
	return nil
}

// updateMetrics safely updates metrics for a tool
func (r *ComprehensiveToolRegistry) updateMetrics(toolName string, update func(*ExecutionMetrics)) {
	r.metrics.mu.Lock()
	defer r.metrics.mu.Unlock()

	if _, exists := r.metrics.executions[toolName]; !exists {
		r.metrics.executions[toolName] = &ExecutionMetrics{}
	}

	update(r.metrics.executions[toolName])
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
