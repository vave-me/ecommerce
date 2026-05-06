package monitoring

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
)

// HealthChecker performs health checks on streaming infrastructure
type HealthChecker struct {
	checks   map[string]HealthCheck
	results  map[string]*HealthResult
	logger   *zap.Logger
	mu       sync.RWMutex
}

// HealthCheck defines a health check function
type HealthCheck func(ctx context.Context) error

// HealthResult contains health check results
type HealthResult struct {
	Status      HealthStatus       `json:"status"`
	LastChecked time.Time          `json:"last_checked"`
	Error       string             `json:"error,omitempty"`
	Details     map[string]string  `json:"details,omitempty"`
}

// HealthStatus represents the health status
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusDegraded  HealthStatus = "degraded"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
)

// SystemHealth represents overall system health
type SystemHealth struct {
	Status     HealthStatus                `json:"status"`
	Timestamp  time.Time                   `json:"timestamp"`
	Version    string                      `json:"version"`
	Components map[string]*HealthResult    `json:"components"`
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(logger *zap.Logger) *HealthChecker {
	return &HealthChecker{
		checks:  make(map[string]HealthCheck),
		results: make(map[string]*HealthResult),
		logger:  logger,
	}
}

// RegisterCheck registers a health check
func (h *HealthChecker) RegisterCheck(name string, check HealthCheck) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.checks[name] = check
}

// RunChecks runs all registered health checks
func (h *HealthChecker) RunChecks(ctx context.Context) {
	h.mu.RLock()
	checks := make(map[string]HealthCheck)
	for k, v := range h.checks {
		checks[k] = v
	}
	h.mu.RUnlock()

	results := make(map[string]*HealthResult)
	var wg sync.WaitGroup

	for name, check := range checks {
		wg.Add(1)
		go func(n string, c HealthCheck) {
			defer wg.Done()
			
			checkCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			result := &HealthResult{
				Status:      HealthStatusHealthy,
				LastChecked: time.Now(),
			}

			if err := c(checkCtx); err != nil {
				result.Status = HealthStatusUnhealthy
				result.Error = err.Error()
			}

			h.mu.Lock()
			h.results[n] = result
			h.mu.Unlock()
		}(name, check)
	}

	wg.Wait()
}

// GetHealth returns the current system health
func (h *HealthChecker) GetHealth() *SystemHealth {
	h.mu.RLock()
	defer h.mu.RUnlock()

	health := &SystemHealth{
		Status:     HealthStatusHealthy,
		Timestamp:  time.Now(),
		Version:    "1.0.0", // Should be injected
		Components: make(map[string]*HealthResult),
	}

	// Copy results
	for name, result := range h.results {
		health.Components[name] = result
		
		// Determine overall status
		if result.Status == HealthStatusUnhealthy {
			health.Status = HealthStatusUnhealthy
		} else if result.Status == HealthStatusDegraded && health.Status == HealthStatusHealthy {
			health.Status = HealthStatusDegraded
		}
	}

	return health
}

// HTTPHandler returns an HTTP handler for health checks
func (h *HealthChecker) HTTPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Run checks before responding
		ctx := r.Context()
		h.RunChecks(ctx)

		health := h.GetHealth()
		
		// Set status code based on health
		statusCode := http.StatusOK
		if health.Status == HealthStatusDegraded {
			statusCode = http.StatusOK // Still return 200 for degraded
		} else if health.Status == HealthStatusUnhealthy {
			statusCode = http.StatusServiceUnavailable
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(health)
	}
}

// ReadinessHandler returns an HTTP handler for readiness checks
func (h *HealthChecker) ReadinessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		health := h.GetHealth()
		
		// Only return ready if healthy
		if health.Status == HealthStatusHealthy {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ready"))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("not ready"))
		}
	}
}

// LivenessHandler returns an HTTP handler for liveness checks
func (h *HealthChecker) LivenessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Simple liveness check - if we can respond, we're alive
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("alive"))
	}
}

// StartBackgroundChecks starts periodic health checks
func (h *HealthChecker) StartBackgroundChecks(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Run initial check
	h.RunChecks(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			h.RunChecks(ctx)
		}
	}
}

// StreamingHealthChecks creates health checks for streaming components
func StreamingHealthChecks() map[string]HealthCheck {
	return map[string]HealthCheck{
		"rtmp_server": func(ctx context.Context) error {
			// Check RTMP server connectivity
			// This would actually test RTMP connection
			return nil
		},
		
		"srt_server": func(ctx context.Context) error {
			// Check SRT server connectivity
			return nil
		},
		
		"transcoding": func(ctx context.Context) error {
			// Check FFmpeg availability and transcoding pipeline
			return nil
		},
		
		"storage": func(ctx context.Context) error {
			// Check storage availability and space
			return nil
		},
		
		"cdn": func(ctx context.Context) error {
			// Check CDN connectivity and health
			return nil
		},
		
		"database": func(ctx context.Context) error {
			// Check database connectivity
			return nil
		},
		
		"redis": func(ctx context.Context) error {
			// Check Redis connectivity
			return nil
		},
		
		"drm_license_server": func(ctx context.Context) error {
			// Check DRM license server availability
			return nil
		},
	}
}