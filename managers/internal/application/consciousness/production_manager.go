package consciousness

import (
	"context"
	"fmt"
	"net/http"
	"time"
	
	"github.com/rs/zerolog"
	"middleman/internal/ai"
	"middleman/internal/ddd"
	"middleman/managers/internal/config"
)

// ProductionConsciousnessManager extends ConsciousnessManager with production features
type ProductionConsciousnessManager struct {
	*ConsciousnessManager
	
	// Production components
	errorHandler      *ErrorHandler
	circuitBreakers   *CircuitBreakerManager
	rateLimiter       *RateLimiter
	metricsCollector  *MetricsCollector
	tracingManager    *TracingManager
	healthChecker     *HealthChecker
	readinessChecker  *ReadinessChecker
	shutdownManager   *ShutdownManager
	enhancedLogger    *EnhancedLogger
	
	// Configuration
	config            *config.ConsciousnessConfig
}

// NewProductionConsciousnessManager creates a production-ready consciousness manager
func NewProductionConsciousnessManager(
	app App,
	memoryCore *MemoryCore,
	patternDetector *PatternDetector,
	learningProcessor *LearningProcessor,
	decisionOrchestrator *DecisionOrchestrator,
	actionExecutor *ActionExecutor,
	aiManager ai.AIClientManager,
	config *config.ConsciousnessConfig,
	logger zerolog.Logger,
) (*ProductionConsciousnessManager, error) {
	// Create base consciousness manager
	baseCM := NewConsciousnessManager(
		app,
		memoryCore,
		patternDetector,
		learningProcessor,
		decisionOrchestrator,
		actionExecutor,
		aiManager,
		logger,
	)
	
	// Create production components
	errorHandler := NewErrorHandler(logger)
	circuitBreakers := NewCircuitBreakerManager(logger)
	metricsCollector := NewMetricsCollector("managers")
	tracingManager := NewTracingManager("managers.consciousness", logger, metricsCollector)
	enhancedLogger := NewEnhancedLogger(logger, tracingManager)
	
	// Create rate limiter
	rateLimiterConfig := RateLimiterConfig{
		GlobalRPS:        config.MaxActionsPerMinute / 60,
		GlobalBurst:      config.MaxActionsPerMinute / 30,
		ComponentLimits: map[string]ComponentLimit{
			"event_processing": {RPS: 100, Burst: 200},
			"pattern_detection": {RPS: 50, Burst: 100},
			"decision_making":   {RPS: 30, Burst: 60},
			"action_execution":  {RPS: 20, Burst: 40},
		},
		DefaultUserRPS:   10,
		DefaultUserBurst: 20,
		VIPUserRPS:       50,
		VIPUserBurst:     100,
		DailyQuotas: map[string]int64{
			"ai_calls:total":      10000,
			"tool_executions:total": 50000,
		},
		MonthlyQuotas: map[string]int64{
			"ai_calls:total":      250000,
			"tool_executions:total": 1500000,
		},
	}
	rateLimiter := NewRateLimiter(rateLimiterConfig, logger, metricsCollector)
	
	// Create health and readiness checkers
	healthChecker := NewHealthChecker(logger, 30*time.Second)
	readinessChecker := NewReadinessChecker(logger)
	
	// Create shutdown manager
	shutdownManager := NewShutdownManager(logger, 30*time.Second)
	
	pcm := &ProductionConsciousnessManager{
		ConsciousnessManager: baseCM,
		errorHandler:         errorHandler,
		circuitBreakers:      circuitBreakers,
		rateLimiter:          rateLimiter,
		metricsCollector:     metricsCollector,
		tracingManager:       tracingManager,
		healthChecker:        healthChecker,
		readinessChecker:     readinessChecker,
		shutdownManager:      shutdownManager,
		enhancedLogger:       enhancedLogger,
		config:               config,
	}
	
	// Register health checks
	pcm.registerHealthChecks()
	
	// Register readiness checks
	pcm.registerReadinessChecks()
	
	// Register shutdown handlers
	pcm.registerShutdownHandlers()
	
	// Start background processes
	pcm.startBackgroundProcesses(context.Background())
	
	return pcm, nil
}

// ProcessEvent processes incoming platform events with production features
func (pcm *ProductionConsciousnessManager) ProcessEvent(ctx context.Context, event ddd.Event) error {
	// Add shutdown check
	if pcm.shutdownManager.IsShuttingDown() {
		return fmt.Errorf("consciousness system is shutting down")
	}
	
	// Start tracing
	ctx, span := pcm.tracingManager.StartEventProcessingSpan(ctx, event.EventName(), event.ID())
	defer func() {
		pcm.tracingManager.EndSpanWithStatus(span, nil)
	}()
	
	// Apply rate limiting
	if err := pcm.rateLimiter.Allow(ctx, "event_processing", event.EventName(), ""); err != nil {
		pcm.metricsCollector.RecordError("rate_limiter", "medium")
		pcm.tracingManager.RecordError(ctx, err, "Rate limit exceeded")
		return err
	}
	
	// Wrap with circuit breaker
	cbConfig := CircuitBreakerConfig{
		Name:             fmt.Sprintf("event_%s", event.EventName()),
		MaxFailures:      5,
		ResetTimeout:     2 * time.Minute,
		SuccessThreshold: 3,
		Timeout:          pcm.config.MaxEventProcessingTime,
	}
	
	cb := pcm.circuitBreakers.GetOrCreate(cbConfig)
	
	// Execute with circuit breaker
	err := cb.Execute(ctx, func(cbCtx context.Context) error {
		// Add panic recovery
		defer pcm.errorHandler.HandlePanic("consciousness_manager", "process_event")
		
		// Record start time
		startTime := time.Now()
		
		// Log with enhanced logger
		pcm.enhancedLogger.LogEventProcessing(ctx, event, "started")
		
		// Process event with base manager
		err := pcm.ConsciousnessManager.ProcessEvent(cbCtx, event)
		
		// Record metrics
		duration := time.Since(startTime)
		pcm.metricsCollector.RecordEventProcessed(event.EventName(), duration, err == nil)
		
		// Handle error if any
		if err != nil {
			pcm.errorHandler.HandleError(cbCtx, "consciousness_manager", "process_event", err)
			pcm.enhancedLogger.Error(ctx).
				Err(err).
				Str("event_type", event.EventName()).
				Msg("Event processing failed")
		} else {
			pcm.enhancedLogger.LogEventProcessing(ctx, event, "completed")
		}
		
		return err
	})
	
	// Update circuit breaker metrics
	if err != nil {
		pcm.metricsCollector.RecordCircuitBreakerFailure(cb.name)
	}
	pcm.metricsCollector.UpdateCircuitBreakerState(cb.name, cb.GetState())
	
	return err
}

// registerHealthChecks registers all health checks
func (pcm *ProductionConsciousnessManager) registerHealthChecks() {
	// Register consciousness health checks
	checks := ConsciousnessHealthChecks(pcm.ConsciousnessManager, pcm.metricsCollector)
	for _, check := range checks {
		pcm.healthChecker.RegisterCheck(check.Name, check)
	}
	
	// Register additional production health checks
	pcm.healthChecker.RegisterCheck("error_rate", Check{
		Name:     "error_rate",
		Critical: false,
		CheckFunc: func(ctx context.Context) error {
			metrics := pcm.errorHandler.GetMetrics()
			if metrics.TotalErrors > 1000 {
				return fmt.Errorf("high error rate: %d total errors", metrics.TotalErrors)
			}
			return nil
		},
	})
	
	pcm.healthChecker.RegisterCheck("circuit_breakers", Check{
		Name:     "circuit_breakers",
		Critical: false,
		CheckFunc: func(ctx context.Context) error {
			openCount := 0
			for _, cb := range pcm.circuitBreakers.GetAllMetrics() {
				if cb.State == "open" {
					openCount++
				}
			}
			if openCount > 5 {
				return fmt.Errorf("too many open circuit breakers: %d", openCount)
			}
			return nil
		},
	})
}

// registerReadinessChecks registers all readiness checks
func (pcm *ProductionConsciousnessManager) registerReadinessChecks() {
	// Register consciousness readiness checks
	checks := ConsciousnessReadinessChecks(pcm.ConsciousnessManager)
	for _, check := range checks {
		pcm.readinessChecker.RegisterCheck(check.Name, check)
	}
	
	// Register additional production readiness checks
	pcm.readinessChecker.RegisterCheck("ai_available", ReadinessCheck{
		Name:     "ai_available",
		Required: false,
		CheckFunc: func(ctx context.Context) error {
			_, err := pcm.aiManager.GetDefaultClient()
			return err
		},
	})
}

// registerShutdownHandlers registers shutdown handlers
func (pcm *ProductionConsciousnessManager) registerShutdownHandlers() {
	handler := NewConsciousnessShutdownHandler(
		pcm.ConsciousnessManager,
		pcm.errorHandler,
		pcm.metricsCollector,
		pcm.rateLimiter,
		pcm.circuitBreakers,
		pcm.logger,
	)
	handler.RegisterShutdownHandlers(pcm.shutdownManager)
}

// startBackgroundProcesses starts background monitoring and maintenance
func (pcm *ProductionConsciousnessManager) startBackgroundProcesses(ctx context.Context) {
	// Start error handler
	pcm.errorHandler.Start(ctx)
	
	// Start performance monitoring
	pcm.metricsCollector.StartPerformanceMonitoring(ctx, 30*time.Second)
	
	// Start signal handler for graceful shutdown
	ctx = pcm.shutdownManager.StartSignalHandler(ctx)
	
	// Periodic health check logging
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				result := pcm.healthChecker.RunHealthChecks(ctx)
				pcm.logger.Info().
					Str("status", string(result.Status)).
					Interface("checks", result.Checks).
					Msg("Periodic health check")
			}
		}
	}()
}

// GetHealthHandler returns HTTP handler for health checks
func (pcm *ProductionConsciousnessManager) GetHealthHandler() func(w http.ResponseWriter, r *http.Request) {
	return pcm.healthChecker.HTTPHandler()
}

// GetReadinessHandler returns HTTP handler for readiness checks
func (pcm *ProductionConsciousnessManager) GetReadinessHandler() func(w http.ResponseWriter, r *http.Request) {
	return pcm.readinessChecker.HTTPHandler()
}

// GetMetrics returns current metrics
func (pcm *ProductionConsciousnessManager) GetMetrics() interface{} {
	return struct {
		ErrorMetrics         ErrorMetrics                `json:"errors"`
		CircuitBreakerMetrics []CircuitBreakerMetrics    `json:"circuit_breakers"`
		RateLimiterStats     map[string]UsageStats       `json:"rate_limits"`
		ConsciousnessStatus  ConsciousnessStatus         `json:"consciousness"`
	}{
		ErrorMetrics:         pcm.errorHandler.GetMetrics(),
		CircuitBreakerMetrics: pcm.circuitBreakers.GetAllMetrics(),
		RateLimiterStats:     pcm.rateLimiter.quotaManager.GetUsageStats(),
		ConsciousnessStatus:  pcm.GetStatus(),
	}
}

// Shutdown initiates graceful shutdown
func (pcm *ProductionConsciousnessManager) Shutdown() {
	pcm.shutdownManager.Shutdown()
}