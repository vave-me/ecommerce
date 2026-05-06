package consciousness

import (
	"context"
	"sync"
	"time"
	
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// MetricsCollector collects and exposes consciousness system metrics
type MetricsCollector struct {
	// Event processing metrics
	eventsProcessed   *prometheus.CounterVec
	eventDuration     *prometheus.HistogramVec
	
	// Pattern detection metrics
	patternsDetected  *prometheus.CounterVec
	patternConfidence *prometheus.HistogramVec
	
	// Decision making metrics
	decisionsMade     *prometheus.CounterVec
	decisionDuration  *prometheus.HistogramVec
	decisionQuality   *prometheus.GaugeVec
	
	// Action execution metrics
	actionsExecuted   *prometheus.CounterVec
	actionDuration    *prometheus.HistogramVec
	actionSuccess     *prometheus.CounterVec
	
	// Tool usage metrics
	toolsUsed         *prometheus.CounterVec
	toolDuration      *prometheus.HistogramVec
	toolErrors        *prometheus.CounterVec
	
	// System health metrics
	systemHealth      *prometheus.GaugeVec
	errorRate         *prometheus.CounterVec
	recoveryAttempts  *prometheus.CounterVec
	
	// Circuit breaker metrics
	circuitState      *prometheus.GaugeVec
	circuitFailures   *prometheus.CounterVec
	
	// Performance metrics
	memoryUsage       prometheus.Gauge
	goroutines        prometheus.Gauge
	cpuUsage          prometheus.Gauge
	
	// Custom metrics storage
	customMetrics     map[string]float64
	mu                sync.RWMutex
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector(namespace string) *MetricsCollector {
	mc := &MetricsCollector{
		customMetrics: make(map[string]float64),
	}
	
	// Initialize event processing metrics
	mc.eventsProcessed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "consciousness",
			Name:      "events_processed_total",
			Help:      "Total number of events processed",
		},
		[]string{"event_type", "status"},
	)
	
	mc.eventDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "consciousness",
			Name:      "event_processing_duration_seconds",
			Help:      "Duration of event processing",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"event_type"},
	)
	
	// Initialize pattern detection metrics
	mc.patternsDetected = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "consciousness",
			Name:      "patterns_detected_total",
			Help:      "Total number of patterns detected",
		},
		[]string{"pattern_type"},
	)
	
	mc.patternConfidence = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "consciousness",
			Name:      "pattern_confidence",
			Help:      "Confidence scores of detected patterns",
			Buckets:   []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0},
		},
		[]string{"pattern_type"},
	)
	
	// Initialize decision making metrics
	mc.decisionsMade = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "consciousness",
			Name:      "decisions_made_total",
			Help:      "Total number of decisions made",
		},
		[]string{"decision_type", "source"},
	)
	
	mc.decisionDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "consciousness",
			Name:      "decision_duration_seconds",
			Help:      "Duration of decision making",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"decision_type"},
	)
	
	mc.decisionQuality = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "consciousness",
			Name:      "decision_quality_score",
			Help:      "Quality score of decisions (0-1)",
		},
		[]string{"decision_type"},
	)
	
	// Initialize action execution metrics
	mc.actionsExecuted = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "consciousness",
			Name:      "actions_executed_total",
			Help:      "Total number of actions executed",
		},
		[]string{"action_type", "status"},
	)
	
	mc.actionDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "consciousness",
			Name:      "action_duration_seconds",
			Help:      "Duration of action execution",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"action_type"},
	)
	
	mc.actionSuccess = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "consciousness",
			Name:      "action_success_total",
			Help:      "Total number of successful actions",
		},
		[]string{"action_type"},
	)
	
	// Initialize tool usage metrics
	mc.toolsUsed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "consciousness",
			Name:      "tools_used_total",
			Help:      "Total number of tool invocations",
		},
		[]string{"tool_name", "status"},
	)
	
	mc.toolDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "consciousness",
			Name:      "tool_duration_seconds",
			Help:      "Duration of tool execution",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"tool_name"},
	)
	
	mc.toolErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "consciousness",
			Name:      "tool_errors_total",
			Help:      "Total number of tool errors",
		},
		[]string{"tool_name", "error_type"},
	)
	
	// Initialize system health metrics
	mc.systemHealth = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "consciousness",
			Name:      "system_health",
			Help:      "Overall system health (0=unhealthy, 1=healthy)",
		},
		[]string{"component"},
	)
	
	mc.errorRate = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "consciousness",
			Name:      "errors_total",
			Help:      "Total number of errors",
		},
		[]string{"component", "severity"},
	)
	
	mc.recoveryAttempts = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "consciousness",
			Name:      "recovery_attempts_total",
			Help:      "Total number of recovery attempts",
		},
		[]string{"component", "status"},
	)
	
	// Initialize circuit breaker metrics
	mc.circuitState = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "consciousness",
			Name:      "circuit_breaker_state",
			Help:      "Circuit breaker state (0=closed, 1=open, 2=half-open)",
		},
		[]string{"circuit_name"},
	)
	
	mc.circuitFailures = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "consciousness",
			Name:      "circuit_breaker_failures_total",
			Help:      "Total number of circuit breaker failures",
		},
		[]string{"circuit_name"},
	)
	
	// Initialize performance metrics
	mc.memoryUsage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "consciousness",
			Name:      "memory_usage_bytes",
			Help:      "Current memory usage in bytes",
		},
	)
	
	mc.goroutines = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "consciousness",
			Name:      "goroutines_count",
			Help:      "Current number of goroutines",
		},
	)
	
	mc.cpuUsage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "consciousness",
			Name:      "cpu_usage_percent",
			Help:      "Current CPU usage percentage",
		},
	)
	
	return mc
}

// RecordEventProcessed records an event processing metric
func (mc *MetricsCollector) RecordEventProcessed(eventType string, duration time.Duration, success bool) {
	status := "success"
	if !success {
		status = "failure"
	}
	
	mc.eventsProcessed.WithLabelValues(eventType, status).Inc()
	mc.eventDuration.WithLabelValues(eventType).Observe(duration.Seconds())
}

// RecordPatternDetected records a pattern detection metric
func (mc *MetricsCollector) RecordPatternDetected(patternType string, confidence float64) {
	mc.patternsDetected.WithLabelValues(patternType).Inc()
	mc.patternConfidence.WithLabelValues(patternType).Observe(confidence)
}

// RecordDecisionMade records a decision making metric
func (mc *MetricsCollector) RecordDecisionMade(decisionType, source string, duration time.Duration) {
	mc.decisionsMade.WithLabelValues(decisionType, source).Inc()
	mc.decisionDuration.WithLabelValues(decisionType).Observe(duration.Seconds())
}

// RecordActionExecuted records an action execution metric
func (mc *MetricsCollector) RecordActionExecuted(actionType string, duration time.Duration, success bool) {
	status := "success"
	if !success {
		status = "failure"
	}
	
	mc.actionsExecuted.WithLabelValues(actionType, status).Inc()
	mc.actionDuration.WithLabelValues(actionType).Observe(duration.Seconds())
	
	if success {
		mc.actionSuccess.WithLabelValues(actionType).Inc()
	}
}

// RecordToolUsage records tool usage metrics
func (mc *MetricsCollector) RecordToolUsage(toolName string, duration time.Duration, success bool, errorType string) {
	status := "success"
	if !success {
		status = "failure"
	}
	
	mc.toolsUsed.WithLabelValues(toolName, status).Inc()
	mc.toolDuration.WithLabelValues(toolName).Observe(duration.Seconds())
	
	if !success && errorType != "" {
		mc.toolErrors.WithLabelValues(toolName, errorType).Inc()
	}
}

// RecordError records an error metric
func (mc *MetricsCollector) RecordError(component string, severity string) {
	mc.errorRate.WithLabelValues(component, severity).Inc()
}

// RecordRecoveryAttempt records a recovery attempt
func (mc *MetricsCollector) RecordRecoveryAttempt(component string, success bool) {
	status := "success"
	if !success {
		status = "failure"
	}
	mc.recoveryAttempts.WithLabelValues(component, status).Inc()
}

// UpdateSystemHealth updates the system health metric
func (mc *MetricsCollector) UpdateSystemHealth(component string, health float64) {
	mc.systemHealth.WithLabelValues(component).Set(health)
}

// UpdateCircuitBreakerState updates circuit breaker metrics
func (mc *MetricsCollector) UpdateCircuitBreakerState(circuitName string, state CircuitState) {
	stateValue := float64(state)
	mc.circuitState.WithLabelValues(circuitName).Set(stateValue)
}

// RecordCircuitBreakerFailure records a circuit breaker failure
func (mc *MetricsCollector) RecordCircuitBreakerFailure(circuitName string) {
	mc.circuitFailures.WithLabelValues(circuitName).Inc()
}

// UpdatePerformanceMetrics updates performance metrics
func (mc *MetricsCollector) UpdatePerformanceMetrics(memoryBytes float64, goroutines float64, cpuPercent float64) {
	mc.memoryUsage.Set(memoryBytes)
	mc.goroutines.Set(goroutines)
	mc.cpuUsage.Set(cpuPercent)
}

// SetDecisionQuality sets the decision quality score
func (mc *MetricsCollector) SetDecisionQuality(decisionType string, quality float64) {
	mc.decisionQuality.WithLabelValues(decisionType).Set(quality)
}

// RecordCustomMetric records a custom metric
func (mc *MetricsCollector) RecordCustomMetric(name string, value float64) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.customMetrics[name] = value
}

// GetCustomMetric retrieves a custom metric
func (mc *MetricsCollector) GetCustomMetric(name string) (float64, bool) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	value, exists := mc.customMetrics[name]
	return value, exists
}

// StartPerformanceMonitoring starts a goroutine to monitor performance metrics
func (mc *MetricsCollector) StartPerformanceMonitoring(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				mc.collectPerformanceMetrics()
			}
		}
	}()
}

// collectPerformanceMetrics collects current performance metrics
func (mc *MetricsCollector) collectPerformanceMetrics() {
	// This would collect actual metrics from runtime
	// For now, using placeholder values
	// In production, you would use runtime.MemStats, runtime.NumGoroutine(), etc.
	
	// Example implementation:
	// var m runtime.MemStats
	// runtime.ReadMemStats(&m)
	// mc.UpdatePerformanceMetrics(float64(m.Alloc), float64(runtime.NumGoroutine()), getCPUPercent())
}