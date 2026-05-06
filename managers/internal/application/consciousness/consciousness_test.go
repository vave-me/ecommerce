package consciousness_test

import (
	"context"
	"testing"
	"time"
	
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	
	"middleman/internal/ai"
	"middleman/internal/ddd"
	"middleman/managers/internal/application/consciousness"
)

// Mock types
type MockApp struct {
	mock.Mock
}

func (m *MockApp) ExecuteTools(ctx context.Context, toolCalls []ai.ToolCall, execCtx interface{}) ([]interface{}, error) {
	args := m.Called(ctx, toolCalls, execCtx)
	return args.Get(0).([]interface{}), args.Error(1)
}

type MockEvent struct {
	eventName string
	eventID   string
	payload   interface{}
	metadata  map[string]interface{}
}

func (e MockEvent) ID() string                      { return e.eventID }
func (e MockEvent) EventName() string               { return e.eventName }
func (e MockEvent) Payload() ddd.EventPayload       { return e.payload }
func (e MockEvent) Metadata() ddd.Metadata          { return e.metadata }
func (e MockEvent) OccurredAt() time.Time           { return time.Now() }

// Test Pattern Detection
func TestPatternDetection(t *testing.T) {
	tests := []struct {
		name           string
		event          ddd.Event
		expectedPattern *consciousness.Pattern
	}{
		{
			name: "detect cart abandonment pattern",
			event: MockEvent{
				eventName: "BasketAbandoned",
				eventID:   "test-123",
				payload: map[string]interface{}{
					"basket_id":    "basket-123",
					"user_id":      "user-456",
					"total_value":  150.0,
					"items_count":  3,
					"abandoned_at": time.Now().Add(-30 * time.Minute),
				},
			},
			expectedPattern: &consciousness.Pattern{
				Type:       "cart_abandonment",
				Confidence: 0.85,
			},
		},
		{
			name: "detect high value order pattern",
			event: MockEvent{
				eventName: "OrderCreated",
				eventID:   "test-456",
				payload: map[string]interface{}{
					"order_id":     "order-789",
					"user_id":      "user-123",
					"total_amount": 5000.0,
					"items_count":  1,
				},
			},
			expectedPattern: &consciousness.Pattern{
				Type:       "fraud_risk",
				Confidence: 0.75,
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create pattern detector
			detector := consciousness.NewPatternDetector(
				consciousness.NewActivitySurgeDetector(),
				consciousness.NewAbandonmentRiskDetector(),
				consciousness.NewFraudRiskDetector(),
				consciousness.NewSupportCrisisDetector(),
				consciousness.NewUserSurgeDetector(),
				consciousness.NewCancellationSpikeDetector(),
			)
			
			// Detect pattern
			pattern := detector.DetectPattern(context.Background(), tt.event)
			
			// Verify
			if tt.expectedPattern != nil {
				require.NotNil(t, pattern)
				assert.Equal(t, tt.expectedPattern.Type, pattern.Type)
				assert.InDelta(t, tt.expectedPattern.Confidence, pattern.Confidence, 0.1)
			}
		})
	}
}

// Test Tool Selection
func TestToolSelection(t *testing.T) {
	logger := zerolog.Nop()
	selector := consciousness.NewToolSelector(logger)
	
	tests := []struct {
		name          string
		event         ddd.Event
		expectedTools []string
	}{
		{
			name: "order created event",
			event: MockEvent{
				eventName: "OrderCreated",
				payload: map[string]interface{}{
					"order_id": "order-123",
					"user_id":  "user-456",
					"amount":   2500.0,
				},
			},
			expectedTools: []string{
				"order_get_by_id",
				"order_create",
				"user_get_by_id",
				"payment_verify_high_value",
			},
		},
		{
			name: "basket abandoned event",
			event: MockEvent{
				eventName: "BasketAbandoned",
				payload: map[string]interface{}{
					"basket_id": "basket-123",
					"user_id":   "user-789",
				},
			},
			expectedTools: []string{
				"basket_get_by_id",
				"user_get_by_id",
				"notification_send_cart_reminder",
				"offer_create_recovery",
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tools := selector.SelectToolsForEvent(context.Background(), tt.event)
			
			// Check that expected tools are included
			for _, expectedTool := range tt.expectedTools {
				assert.Contains(t, tools, expectedTool)
			}
		})
	}
}

// Test Circuit Breaker
func TestCircuitBreaker(t *testing.T) {
	logger := zerolog.Nop()
	
	config := consciousness.CircuitBreakerConfig{
		Name:             "test-circuit",
		MaxFailures:      3,
		ResetTimeout:     100 * time.Millisecond,
		SuccessThreshold: 2,
		Timeout:          50 * time.Millisecond,
	}
	
	cb := consciousness.NewCircuitBreaker(config, logger)
	
	// Test circuit opens after failures
	failCount := 0
	for i := 0; i < 5; i++ {
		err := cb.Execute(context.Background(), func(ctx context.Context) error {
			failCount++
			if failCount <= 3 {
				return assert.AnError
			}
			return nil
		})
		
		if i < 3 {
			assert.Error(t, err)
		} else {
			// Circuit should be open
			assert.ErrorIs(t, err, consciousness.ErrCircuitOpen)
		}
	}
	
	// Wait for reset timeout
	time.Sleep(150 * time.Millisecond)
	
	// Circuit should be half-open, allowing one request
	err := cb.Execute(context.Background(), func(ctx context.Context) error {
		return nil
	})
	assert.NoError(t, err)
	
	// Need one more success to close circuit
	err = cb.Execute(context.Background(), func(ctx context.Context) error {
		return nil
	})
	assert.NoError(t, err)
	
	// Circuit should be closed now
	assert.Equal(t, consciousness.StateClosed, cb.GetState())
}

// Test Rate Limiter
func TestRateLimiter(t *testing.T) {
	logger := zerolog.Nop()
	metrics := consciousness.NewMetricsCollector("test")
	
	config := consciousness.RateLimiterConfig{
		GlobalRPS:        2,
		GlobalBurst:      4,
		ComponentLimits: map[string]consciousness.ComponentLimit{
			"test_component": {RPS: 1, Burst: 2},
		},
	}
	
	limiter := consciousness.NewRateLimiter(config, logger, metrics)
	
	// Test rate limiting
	successCount := 0
	for i := 0; i < 10; i++ {
		err := limiter.Allow(context.Background(), "test_component", "test_op", "user-123")
		if err == nil {
			successCount++
		}
	}
	
	// Should allow burst initially
	assert.LessOrEqual(t, successCount, config.GlobalBurst)
}

// Test Health Checks
func TestHealthChecks(t *testing.T) {
	logger := zerolog.Nop()
	checker := consciousness.NewHealthChecker(logger, 1*time.Second)
	
	// Register a healthy check
	checker.RegisterCheck("test_healthy", consciousness.Check{
		Name:     "test_healthy",
		Critical: true,
		CheckFunc: func(ctx context.Context) error {
			return nil
		},
	})
	
	// Register an unhealthy check
	checker.RegisterCheck("test_unhealthy", consciousness.Check{
		Name:     "test_unhealthy",
		Critical: false,
		CheckFunc: func(ctx context.Context) error {
			return assert.AnError
		},
	})
	
	// Run health checks
	result := checker.RunHealthChecks(context.Background())
	
	assert.Equal(t, consciousness.HealthStatusDegraded, result.Status)
	assert.Equal(t, consciousness.HealthStatusHealthy, result.Checks["test_healthy"].Status)
	assert.Equal(t, consciousness.HealthStatusUnhealthy, result.Checks["test_unhealthy"].Status)
}

// Test Error Handler
func TestErrorHandler(t *testing.T) {
	logger := zerolog.Nop()
	handler := consciousness.NewErrorHandler(logger)
	
	ctx := context.Background()
	handler.Start(ctx)
	
	// Test error handling
	handler.HandleError(ctx, "test_component", "test_operation", assert.AnError)
	
	// Test panic recovery
	func() {
		defer handler.HandlePanic("test_component", "panic_operation")
		panic("test panic")
	}()
	
	// Give some time for async processing
	time.Sleep(100 * time.Millisecond)
	
	// Check metrics
	metrics := handler.GetMetrics()
	assert.Greater(t, metrics.TotalErrors, int64(0))
}

// Test Shutdown Manager
func TestShutdownManager(t *testing.T) {
	logger := zerolog.Nop()
	sm := consciousness.NewShutdownManager(logger, 1*time.Second)
	
	cleanupCalled := false
	sm.RegisterCleanup("test_cleanup", 1, 100*time.Millisecond, func(ctx context.Context) error {
		cleanupCalled = true
		return nil
	})
	
	// Start shutdown
	go sm.Shutdown()
	
	// Wait for shutdown signal
	<-sm.WaitForShutdown()
	
	// Give time for cleanup
	time.Sleep(200 * time.Millisecond)
	
	assert.True(t, cleanupCalled)
	assert.True(t, sm.IsShuttingDown())
}

// Integration test for full consciousness flow
func TestConsciousnessIntegration(t *testing.T) {
	// This would be a more comprehensive test with real components
	// For now, just a placeholder
	t.Skip("Full integration test requires complete setup")
}