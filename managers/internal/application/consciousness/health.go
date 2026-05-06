package consciousness

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
	
	"github.com/rs/zerolog"
)

// HealthStatus represents the health status of a component
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusDegraded  HealthStatus = "degraded"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
)

// HealthChecker provides health checking capabilities
type HealthChecker struct {
	checks         map[string]Check
	mu             sync.RWMutex
	logger         zerolog.Logger
	lastCheckTime  time.Time
	cacheDuration  time.Duration
	cachedResult   *HealthCheckResult
}

// Check represents a health check function
type Check struct {
	Name        string
	CheckFunc   func(ctx context.Context) error
	Critical    bool
	Timeout     time.Duration
	SuccessRate float64 // Required success rate (0-1)
}

// HealthCheckResult represents the result of health checks
type HealthCheckResult struct {
	Status     HealthStatus              `json:"status"`
	Timestamp  time.Time                 `json:"timestamp"`
	Checks     map[string]CheckResult    `json:"checks"`
	Duration   time.Duration             `json:"duration_ms"`
	Version    string                    `json:"version"`
	Uptime     time.Duration             `json:"uptime_ms"`
}

// CheckResult represents the result of a single check
type CheckResult struct {
	Status    HealthStatus  `json:"status"`
	Message   string        `json:"message,omitempty"`
	Duration  time.Duration `json:"duration_ms"`
	Timestamp time.Time     `json:"timestamp"`
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(logger zerolog.Logger, cacheDuration time.Duration) *HealthChecker {
	return &HealthChecker{
		checks:        make(map[string]Check),
		logger:        logger,
		cacheDuration: cacheDuration,
	}
}

// RegisterCheck registers a health check
func (hc *HealthChecker) RegisterCheck(name string, check Check) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	
	check.Name = name
	if check.Timeout == 0 {
		check.Timeout = 5 * time.Second
	}
	if check.SuccessRate == 0 {
		check.SuccessRate = 1.0 // Default to 100% success required
	}
	
	hc.checks[name] = check
	
	hc.logger.Info().
		Str("check_name", name).
		Bool("critical", check.Critical).
		Dur("timeout", check.Timeout).
		Msg("Registered health check")
}

// RunHealthChecks runs all registered health checks
func (hc *HealthChecker) RunHealthChecks(ctx context.Context) *HealthCheckResult {
	// Check cache
	hc.mu.RLock()
	if hc.cachedResult != nil && time.Since(hc.lastCheckTime) < hc.cacheDuration {
		cached := hc.cachedResult
		hc.mu.RUnlock()
		return cached
	}
	hc.mu.RUnlock()
	
	// Run checks
	startTime := time.Now()
	result := &HealthCheckResult{
		Timestamp: startTime,
		Checks:    make(map[string]CheckResult),
		Version:   "1.0.0", // Should be injected from build
		Uptime:    time.Since(startTime), // Should track actual start time
	}
	
	hc.mu.RLock()
	checks := make(map[string]Check)
	for name, check := range hc.checks {
		checks[name] = check
	}
	hc.mu.RUnlock()
	
	// Run checks concurrently
	var wg sync.WaitGroup
	resultChan := make(chan struct {
		name   string
		result CheckResult
	}, len(checks))
	
	for name, check := range checks {
		wg.Add(1)
		go func(name string, check Check) {
			defer wg.Done()
			
			checkResult := hc.runSingleCheck(ctx, check)
			resultChan <- struct {
				name   string
				result CheckResult
			}{name: name, result: checkResult}
		}(name, check)
	}
	
	wg.Wait()
	close(resultChan)
	
	// Collect results
	criticalFailed := false
	degradedCount := 0
	
	for r := range resultChan {
		result.Checks[r.name] = r.result
		
		if r.result.Status == HealthStatusUnhealthy {
			if checks[r.name].Critical {
				criticalFailed = true
			} else {
				degradedCount++
			}
		} else if r.result.Status == HealthStatusDegraded {
			degradedCount++
		}
	}
	
	// Determine overall status
	if criticalFailed {
		result.Status = HealthStatusUnhealthy
	} else if degradedCount > 0 {
		result.Status = HealthStatusDegraded
	} else {
		result.Status = HealthStatusHealthy
	}
	
	result.Duration = time.Since(startTime)
	
	// Cache result
	hc.mu.Lock()
	hc.cachedResult = result
	hc.lastCheckTime = time.Now()
	hc.mu.Unlock()
	
	return result
}

// runSingleCheck runs a single health check
func (hc *HealthChecker) runSingleCheck(ctx context.Context, check Check) CheckResult {
	checkCtx, cancel := context.WithTimeout(ctx, check.Timeout)
	defer cancel()
	
	startTime := time.Now()
	err := check.CheckFunc(checkCtx)
	duration := time.Since(startTime)
	
	result := CheckResult{
		Timestamp: startTime,
		Duration:  duration,
	}
	
	if err != nil {
		result.Status = HealthStatusUnhealthy
		result.Message = err.Error()
		
		hc.logger.Warn().
			Str("check_name", check.Name).
			Err(err).
			Dur("duration", duration).
			Msg("Health check failed")
	} else {
		result.Status = HealthStatusHealthy
		
		hc.logger.Debug().
			Str("check_name", check.Name).
			Dur("duration", duration).
			Msg("Health check passed")
	}
	
	return result
}

// HTTPHandler returns an HTTP handler for health checks
func (hc *HealthChecker) HTTPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		result := hc.RunHealthChecks(ctx)
		
		// Set appropriate status code
		statusCode := http.StatusOK
		switch result.Status {
		case HealthStatusUnhealthy:
			statusCode = http.StatusServiceUnavailable
		case HealthStatusDegraded:
			statusCode = http.StatusOK // Still return 200 for degraded
		}
		
		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		
		if err := json.NewEncoder(w).Encode(result); err != nil {
			hc.logger.Error().Err(err).Msg("Failed to encode health check response")
		}
	}
}

// ReadinessChecker provides readiness checking capabilities
type ReadinessChecker struct {
	checks map[string]ReadinessCheck
	mu     sync.RWMutex
	logger zerolog.Logger
}

// ReadinessCheck represents a readiness check
type ReadinessCheck struct {
	Name      string
	CheckFunc func(ctx context.Context) error
	Required  bool
}

// ReadinessResult represents the result of readiness checks
type ReadinessResult struct {
	Ready     bool                         `json:"ready"`
	Timestamp time.Time                    `json:"timestamp"`
	Checks    map[string]ReadinessCheckResult `json:"checks"`
}

// ReadinessCheckResult represents the result of a single readiness check
type ReadinessCheckResult struct {
	Ready    bool   `json:"ready"`
	Message  string `json:"message,omitempty"`
	Required bool   `json:"required"`
}

// NewReadinessChecker creates a new readiness checker
func NewReadinessChecker(logger zerolog.Logger) *ReadinessChecker {
	return &ReadinessChecker{
		checks: make(map[string]ReadinessCheck),
		logger: logger,
	}
}

// RegisterCheck registers a readiness check
func (rc *ReadinessChecker) RegisterCheck(name string, check ReadinessCheck) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	
	check.Name = name
	rc.checks[name] = check
	
	rc.logger.Info().
		Str("check_name", name).
		Bool("required", check.Required).
		Msg("Registered readiness check")
}

// CheckReadiness runs all readiness checks
func (rc *ReadinessChecker) CheckReadiness(ctx context.Context) *ReadinessResult {
	result := &ReadinessResult{
		Timestamp: time.Now(),
		Checks:    make(map[string]ReadinessCheckResult),
		Ready:     true,
	}
	
	rc.mu.RLock()
	checks := make(map[string]ReadinessCheck)
	for name, check := range rc.checks {
		checks[name] = check
	}
	rc.mu.RUnlock()
	
	for name, check := range checks {
		checkResult := ReadinessCheckResult{
			Required: check.Required,
		}
		
		if err := check.CheckFunc(ctx); err != nil {
			checkResult.Ready = false
			checkResult.Message = err.Error()
			
			if check.Required {
				result.Ready = false
			}
			
			rc.logger.Warn().
				Str("check_name", name).
				Err(err).
				Bool("required", check.Required).
				Msg("Readiness check failed")
		} else {
			checkResult.Ready = true
			
			rc.logger.Debug().
				Str("check_name", name).
				Msg("Readiness check passed")
		}
		
		result.Checks[name] = checkResult
	}
	
	return result
}

// HTTPHandler returns an HTTP handler for readiness checks
func (rc *ReadinessChecker) HTTPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		result := rc.CheckReadiness(ctx)
		
		// Set appropriate status code
		statusCode := http.StatusOK
		if !result.Ready {
			statusCode = http.StatusServiceUnavailable
		}
		
		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		
		if err := json.NewEncoder(w).Encode(result); err != nil {
			rc.logger.Error().Err(err).Msg("Failed to encode readiness check response")
		}
	}
}

// ConsciousnessHealthChecks creates standard health checks for consciousness components
func ConsciousnessHealthChecks(cm *ConsciousnessManager, metrics *MetricsCollector) []Check {
	return []Check{
		{
			Name:     "memory_core",
			Critical: true,
			CheckFunc: func(ctx context.Context) error {
				// Check if memory core is responsive
				if cm.memoryCore == nil {
					return fmt.Errorf("memory core not initialized")
				}
				// Could add actual memory operations check here
				return nil
			},
		},
		{
			Name:     "pattern_detector",
			Critical: true,
			CheckFunc: func(ctx context.Context) error {
				if cm.patternDetector == nil {
					return fmt.Errorf("pattern detector not initialized")
				}
				return nil
			},
		},
		{
			Name:     "decision_maker",
			Critical: true,
			CheckFunc: func(ctx context.Context) error {
				if cm.decisionMaker == nil {
					return fmt.Errorf("decision maker not initialized")
				}
				return nil
			},
		},
		{
			Name:     "action_executor",
			Critical: true,
			CheckFunc: func(ctx context.Context) error {
				if cm.actionExecutor == nil {
					return fmt.Errorf("action executor not initialized")
				}
				return nil
			},
		},
		{
			Name:     "ai_manager",
			Critical: false,
			CheckFunc: func(ctx context.Context) error {
				if cm.aiManager == nil {
					return fmt.Errorf("AI manager not initialized")
				}
				// Check if at least one AI provider is available
				_, err := cm.aiManager.GetDefaultClient()
				return err
			},
		},
		{
			Name:     "metrics_collector",
			Critical: false,
			CheckFunc: func(ctx context.Context) error {
				if metrics == nil {
					return fmt.Errorf("metrics collector not initialized")
				}
				return nil
			},
		},
	}
}

// ConsciousnessReadinessChecks creates standard readiness checks
func ConsciousnessReadinessChecks(cm *ConsciousnessManager) []ReadinessCheck {
	return []ReadinessCheck{
		{
			Name:     "consciousness_active",
			Required: true,
			CheckFunc: func(ctx context.Context) error {
				status := cm.GetStatus()
				if !status.Active {
					return fmt.Errorf("consciousness system not active")
				}
				return nil
			},
		},
		{
			Name:     "event_processing",
			Required: true,
			CheckFunc: func(ctx context.Context) error {
				status := cm.GetStatus()
				// Check if system has processed events recently
				if time.Since(status.LastActivity) > 5*time.Minute {
					return fmt.Errorf("no recent activity detected")
				}
				return nil
			},
		},
	}
}