package application

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"strings"
	"sync"
	"time"
)

// ErrorHandlingEnhancer provides comprehensive error handling improvements
type ErrorHandlingEnhancer struct {
	errorLogger     *ErrorLogger
	retryManager    *RetryManager
	circuitBreaker  *ErrorCircuitBreaker
	errorAggregator *ErrorAggregator
	alertManager    *AlertManager
	recoveryHandler *RecoveryHandler
	metrics         *ErrorMetrics
}

// ErrorLogger handles structured error logging
type ErrorLogger struct {
	logLevel   LogLevel
	logFile    string
	structured bool
	contextual bool
	async      bool
	buffer     chan ErrorLogEntry
	shutdown   chan struct{}
	wg         sync.WaitGroup
}

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

type ErrorLogEntry struct {
	Timestamp  time.Time              `json:"timestamp"`
	Level      LogLevel               `json:"level"`
	Message    string                 `json:"message"`
	Error      error                  `json:"error"`
	Context    map[string]interface{} `json:"context"`
	StackTrace string                 `json:"stack_trace"`
	Goroutine  int                    `json:"goroutine"`
	Component  string                 `json:"component"`
	UserID     string                 `json:"user_id"`
	RequestID  string                 `json:"request_id"`
	Severity   string                 `json:"severity"`
	Category   string                 `json:"category"`
}

// RetryManager handles intelligent retry logic
type RetryManager struct {
	policies          map[string]*RetryPolicy
	backoffStrategies map[string]BackoffStrategy
	maxRetries        int
	baseDelay         time.Duration
	maxDelay          time.Duration
	jitter            bool
}

type RetryPolicy struct {
	MaxAttempts     int           `json:"max_attempts"`
	BackoffType     string        `json:"backoff_type"`
	BaseDelay       time.Duration `json:"base_delay"`
	MaxDelay        time.Duration `json:"max_delay"`
	Multiplier      float64       `json:"multiplier"`
	Jitter          bool          `json:"jitter"`
	RetryableErrors []string      `json:"retryable_errors"`
	StopConditions  []string      `json:"stop_conditions"`
}

type BackoffStrategy interface {
	NextDelay(attempt int, lastDelay time.Duration) time.Duration
}

// ErrorCircuitBreaker prevents cascading failures
type ErrorCircuitBreaker struct {
	name         string
	failureCount int
	successCount int
	lastFailure  time.Time
	state        CircuitState
	config       *ErrorCircuitBreakerConfig
	mutex        sync.RWMutex
}

type CircuitState int

const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

type ErrorCircuitBreakerConfig struct {
	FailureThreshold int           `json:"failure_threshold"`
	SuccessThreshold int           `json:"success_threshold"`
	OpenTimeout      time.Duration `json:"open_timeout"`
	HalfOpenTimeout  time.Duration `json:"half_open_timeout"`
}

// ErrorAggregator collects and groups similar errors
type ErrorAggregator struct {
	errorGroups     map[string]*ErrorGroup
	timeWindow      time.Duration
	maxGroups       int
	cleanupInterval time.Duration
	mutex           sync.RWMutex
}

type ErrorGroup struct {
	Pattern    string                 `json:"pattern"`
	Count      int                    `json:"count"`
	FirstSeen  time.Time              `json:"first_seen"`
	LastSeen   time.Time              `json:"last_seen"`
	Examples   []string               `json:"examples"`
	Context    map[string]interface{} `json:"context"`
	Severity   string                 `json:"severity"`
	Impact     string                 `json:"impact"`
	Resolution string                 `json:"resolution"`
}

// AlertManager handles error-based alerts
type AlertManager struct {
	alerts         chan Alert
	thresholds     map[string]*AlertThreshold
	subscribers    map[string][]AlertSubscriber
	suppressions   map[string]time.Time
	cooldownPeriod time.Duration
}

type Alert struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Severity   string                 `json:"severity"`
	Message    string                 `json:"message"`
	Component  string                 `json:"component"`
	Timestamp  time.Time              `json:"timestamp"`
	Context    map[string]interface{} `json:"context"`
	Actions    []string               `json:"actions"`
	Escalation string                 `json:"escalation"`
}

type AlertThreshold struct {
	ErrorRate  float64       `json:"error_rate"`
	Count      int           `json:"count"`
	TimeWindow time.Duration `json:"time_window"`
	Severity   string        `json:"severity"`
	Cooldown   time.Duration `json:"cooldown"`
}

type AlertSubscriber interface {
	HandleAlert(alert Alert) error
}

// RecoveryHandler handles panic recovery and error recovery
type RecoveryHandler struct {
	recoveryStrategies map[string]RecoveryStrategy
	panicHandlers      []PanicHandler
	errorHandlers      []ErrorHandler
	metrics            *RecoveryMetrics
}

type RecoveryStrategy interface {
	Recover(err error, context map[string]interface{}) error
}

type PanicHandler func(recovered interface{}, stackTrace string)
type ErrorHandler func(err error, context map[string]interface{}) error

type RecoveryMetrics struct {
	PanicsRecovered  int64            `json:"panics_recovered"`
	ErrorsRecovered  int64            `json:"errors_recovered"`
	FailedRecoveries int64            `json:"failed_recoveries"`
	LastRecoveryTime time.Time        `json:"last_recovery_time"`
	RecoveryMethods  map[string]int64 `json:"recovery_methods"`
}

// ErrorMetrics tracks error-related metrics
type ErrorMetrics struct {
	TotalErrors          int64              `json:"total_errors"`
	ErrorsByCategory     map[string]int64   `json:"errors_by_category"`
	ErrorsByComponent    map[string]int64   `json:"errors_by_component"`
	ErrorRate            float64            `json:"error_rate"`
	MTTR                 time.Duration      `json:"mttr"` // Mean Time To Recovery
	MTBF                 time.Duration      `json:"mtbf"` // Mean Time Between Failures
	TopErrors            []ErrorSummary     `json:"top_errors"`
	ErrorTrends          map[string][]int64 `json:"error_trends"`
	SeverityDistribution map[string]int64   `json:"severity_distribution"`

	mutex      sync.RWMutex
	lastUpdate time.Time
}

type ErrorSummary struct {
	Message  string    `json:"message"`
	Count    int64     `json:"count"`
	LastSeen time.Time `json:"last_seen"`
	Severity string    `json:"severity"`
	Impact   string    `json:"impact"`
}

// NewErrorHandlingEnhancer creates a comprehensive error handling enhancer
func NewErrorHandlingEnhancer() *ErrorHandlingEnhancer {
	ehe := &ErrorHandlingEnhancer{
		metrics: &ErrorMetrics{
			ErrorsByCategory:     make(map[string]int64),
			ErrorsByComponent:    make(map[string]int64),
			ErrorTrends:          make(map[string][]int64),
			SeverityDistribution: make(map[string]int64),
		},
	}

	// Initialize error logger
	ehe.errorLogger = &ErrorLogger{
		logLevel:   LogLevelInfo,
		structured: true,
		contextual: true,
		async:      true,
		buffer:     make(chan ErrorLogEntry, 1000),
		shutdown:   make(chan struct{}),
	}

	// Initialize retry manager
	ehe.retryManager = &RetryManager{
		policies:          make(map[string]*RetryPolicy),
		backoffStrategies: make(map[string]BackoffStrategy),
		maxRetries:        3,
		baseDelay:         time.Millisecond * 100,
		maxDelay:          time.Second * 30,
		jitter:            true,
	}

	// Initialize circuit breaker
	ehe.circuitBreaker = &ErrorCircuitBreaker{
		name:  "default",
		state: CircuitClosed,
		config: &ErrorCircuitBreakerConfig{
			FailureThreshold: 5,
			SuccessThreshold: 2,
			OpenTimeout:      time.Second * 60,
			HalfOpenTimeout:  time.Second * 30,
		},
	}

	// Initialize error aggregator
	ehe.errorAggregator = &ErrorAggregator{
		errorGroups:     make(map[string]*ErrorGroup),
		timeWindow:      time.Hour,
		maxGroups:       1000,
		cleanupInterval: time.Minute * 10,
	}

	// Initialize alert manager
	ehe.alertManager = &AlertManager{
		alerts:         make(chan Alert, 100),
		thresholds:     make(map[string]*AlertThreshold),
		subscribers:    make(map[string][]AlertSubscriber),
		suppressions:   make(map[string]time.Time),
		cooldownPeriod: time.Minute * 5,
	}

	// Initialize recovery handler
	ehe.recoveryHandler = &RecoveryHandler{
		recoveryStrategies: make(map[string]RecoveryStrategy),
		metrics: &RecoveryMetrics{
			RecoveryMethods: make(map[string]int64),
		},
	}

	// Start background workers
	go ehe.startErrorLogger()
	go ehe.startAggregator()
	go ehe.startAlertManager()
	go ehe.startMetricsCollector()

	return ehe
}

// HandleError handles an error with comprehensive processing
func (ehe *ErrorHandlingEnhancer) HandleError(err error, context map[string]interface{}) error {
	if err == nil {
		return nil
	}

	// Log the error
	ehe.logError(err, context)

	// Aggregate similar errors
	ehe.aggregateError(err, context)

	// Check circuit breaker
	if ehe.circuitBreaker.ShouldBreak(err) {
		return fmt.Errorf("circuit breaker open: %w", err)
	}

	// Update metrics
	ehe.updateMetrics(err, context)

	// Check for alerts
	ehe.checkAlerts(err, context)

	// Attempt recovery
	if recoveredErr := ehe.attemptRecovery(err, context); recoveredErr != nil {
		return recoveredErr
	}

	return err
}

// WithRetry executes a function with retry logic
func (ehe *ErrorHandlingEnhancer) WithRetry(ctx context.Context, operation func() error, policyName string) error {
	policy, exists := ehe.retryManager.policies[policyName]
	if !exists {
		policy = &RetryPolicy{
			MaxAttempts: ehe.retryManager.maxRetries,
			BackoffType: "exponential",
			BaseDelay:   ehe.retryManager.baseDelay,
			MaxDelay:    ehe.retryManager.maxDelay,
			Multiplier:  2.0,
			Jitter:      ehe.retryManager.jitter,
		}
	}

	for attempt := 0; attempt < policy.MaxAttempts; attempt++ {
		err := operation()
		if err == nil {
			return nil
		}

		// Check if error is retryable
		if !ehe.isRetryableError(err, policy) {
			return err
		}

		// Check context cancellation
		if ctx.Err() != nil {
			return ctx.Err()
		}

		// Calculate delay
		delay := ehe.calculateBackoff(attempt, policy)

		// Wait before retry
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}

		log.Printf("Retrying operation (attempt %d/%d) after %v", attempt+1, policy.MaxAttempts, delay)
	}

	return fmt.Errorf("operation failed after %d attempts", policy.MaxAttempts)
}

// SafeExecute executes a function with panic recovery
func (ehe *ErrorHandlingEnhancer) SafeExecute(name string, fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// Get stack trace
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			stackTrace := string(buf[:n])

			// Handle panic
			for _, handler := range ehe.recoveryHandler.panicHandlers {
				handler(r, stackTrace)
			}

			// Convert panic to error
			err = fmt.Errorf("panic in %s: %v", name, r)

			// Log panic
			ehe.logError(err, map[string]interface{}{
				"panic_value": r,
				"stack_trace": stackTrace,
				"component":   name,
			})

			// Update metrics
			ehe.recoveryHandler.metrics.PanicsRecovered++
			ehe.recoveryHandler.metrics.LastRecoveryTime = time.Now()
		}
	}()

	return fn()
}

// logError logs an error with context
func (ehe *ErrorHandlingEnhancer) logError(err error, context map[string]interface{}) {
	entry := ErrorLogEntry{
		Timestamp: time.Now(),
		Level:     LogLevelError,
		Message:   err.Error(),
		Error:     err,
		Context:   context,
		Goroutine: runtime.NumGoroutine(),
		Severity:  ehe.determineSeverity(err),
		Category:  ehe.categorizeError(err),
	}

	// Add request/user context if available
	if context != nil {
		if userID, ok := context["user_id"]; ok {
			entry.UserID = fmt.Sprintf("%v", userID)
		}
		if requestID, ok := context["request_id"]; ok {
			entry.RequestID = fmt.Sprintf("%v", requestID)
		}
		if component, ok := context["component"]; ok {
			entry.Component = fmt.Sprintf("%v", component)
		}
	}

	// Add stack trace for severe errors
	if entry.Severity == "high" || entry.Severity == "critical" {
		buf := make([]byte, 4096)
		n := runtime.Stack(buf, false)
		entry.StackTrace = string(buf[:n])
	}

	// Send to logger
	select {
	case ehe.errorLogger.buffer <- entry:
	default:
		log.Printf("Error logger buffer full, dropping error: %v", err)
	}
}

// aggregateError groups similar errors
func (ehe *ErrorHandlingEnhancer) aggregateError(err error, context map[string]interface{}) {
	pattern := ehe.extractErrorPattern(err)

	ehe.errorAggregator.mutex.Lock()
	defer ehe.errorAggregator.mutex.Unlock()

	group, exists := ehe.errorAggregator.errorGroups[pattern]
	if !exists {
		group = &ErrorGroup{
			Pattern:   pattern,
			Count:     0,
			FirstSeen: time.Now(),
			Examples:  []string{},
			Context:   make(map[string]interface{}),
			Severity:  ehe.determineSeverity(err),
		}
		ehe.errorAggregator.errorGroups[pattern] = group
	}

	group.Count++
	group.LastSeen = time.Now()

	// Add example if we have less than 5
	if len(group.Examples) < 5 {
		group.Examples = append(group.Examples, err.Error())
	}

	// Merge context
	for k, v := range context {
		group.Context[k] = v
	}
}

// updateMetrics updates error metrics
func (ehe *ErrorHandlingEnhancer) updateMetrics(err error, context map[string]interface{}) {
	ehe.metrics.mutex.Lock()
	defer ehe.metrics.mutex.Unlock()

	ehe.metrics.TotalErrors++

	// Update by category
	category := ehe.categorizeError(err)
	ehe.metrics.ErrorsByCategory[category]++

	// Update by component
	if context != nil {
		if component, ok := context["component"]; ok {
			componentStr := fmt.Sprintf("%v", component)
			ehe.metrics.ErrorsByComponent[componentStr]++
		}
	}

	// Update severity distribution
	severity := ehe.determineSeverity(err)
	ehe.metrics.SeverityDistribution[severity]++

	ehe.metrics.lastUpdate = time.Now()
}

// checkAlerts checks if alerts should be triggered
func (ehe *ErrorHandlingEnhancer) checkAlerts(err error, context map[string]interface{}) {
	category := ehe.categorizeError(err)
	threshold, exists := ehe.alertManager.thresholds[category]
	if !exists {
		return
	}

	// Check if alert is suppressed
	suppressUntil, suppressed := ehe.alertManager.suppressions[category]
	if suppressed && time.Now().Before(suppressUntil) {
		return
	}

	// Check threshold
	errorCount := ehe.metrics.ErrorsByCategory[category]
	if errorCount >= int64(threshold.Count) {
		alert := Alert{
			ID:        fmt.Sprintf("%s-%d", category, time.Now().Unix()),
			Type:      "error_threshold",
			Severity:  threshold.Severity,
			Message:   fmt.Sprintf("Error threshold exceeded for %s: %d errors", category, errorCount),
			Component: category,
			Timestamp: time.Now(),
			Context:   context,
			Actions:   []string{"investigate", "scale_up", "notify_team"},
		}

		select {
		case ehe.alertManager.alerts <- alert:
			// Suppress future alerts for cooldown period
			ehe.alertManager.suppressions[category] = time.Now().Add(threshold.Cooldown)
		default:
			log.Printf("Alert queue full, dropping alert: %s", alert.Message)
		}
	}
}

// attemptRecovery attempts to recover from an error
func (ehe *ErrorHandlingEnhancer) attemptRecovery(err error, context map[string]interface{}) error {
	errorType := ehe.categorizeError(err)

	strategy, exists := ehe.recoveryHandler.recoveryStrategies[errorType]
	if !exists {
		return err
	}

	recoveredErr := strategy.Recover(err, context)
	if recoveredErr == nil {
		ehe.recoveryHandler.metrics.ErrorsRecovered++
		ehe.recoveryHandler.metrics.RecoveryMethods[errorType]++
	} else {
		ehe.recoveryHandler.metrics.FailedRecoveries++
	}

	return recoveredErr
}

// Helper methods

func (ehe *ErrorHandlingEnhancer) determineSeverity(err error) string {
	errStr := err.Error()

	// Critical errors
	if contains(errStr, []string{"panic", "fatal", "critical", "corruption"}) {
		return "critical"
	}

	// High severity errors
	if contains(errStr, []string{"timeout", "connection", "database", "auth"}) {
		return "high"
	}

	// Medium severity errors
	if contains(errStr, []string{"validation", "permission", "not found"}) {
		return "medium"
	}

	return "low"
}

func (ehe *ErrorHandlingEnhancer) categorizeError(err error) string {
	errStr := err.Error()

	if contains(errStr, []string{"sql", "database", "db"}) {
		return "database"
	}
	if contains(errStr, []string{"http", "request", "response"}) {
		return "network"
	}
	if contains(errStr, []string{"auth", "token", "permission"}) {
		return "authentication"
	}
	if contains(errStr, []string{"validation", "invalid", "malformed"}) {
		return "validation"
	}
	if contains(errStr, []string{"timeout", "deadline"}) {
		return "timeout"
	}

	return "unknown"
}

func (ehe *ErrorHandlingEnhancer) extractErrorPattern(err error) string {
	// Simplified pattern extraction - in production, this would be more sophisticated
	errStr := err.Error()

	// Replace numbers and UUIDs with placeholders
	// This is a basic implementation
	return errStr
}

func (ehe *ErrorHandlingEnhancer) isRetryableError(err error, policy *RetryPolicy) bool {
	errStr := err.Error()

	// Check stop conditions
	for _, stopCondition := range policy.StopConditions {
		if contains(errStr, []string{stopCondition}) {
			return false
		}
	}

	// Check retryable errors
	if len(policy.RetryableErrors) > 0 {
		return contains(errStr, policy.RetryableErrors)
	}

	// Default retryable errors
	retryablePatterns := []string{"timeout", "connection", "temporary", "rate limit"}
	return contains(errStr, retryablePatterns)
}

func (ehe *ErrorHandlingEnhancer) calculateBackoff(attempt int, policy *RetryPolicy) time.Duration {
	delay := policy.BaseDelay

	switch policy.BackoffType {
	case "exponential":
		delay = time.Duration(float64(delay) * (policy.Multiplier * float64(attempt)))
	case "linear":
		delay = time.Duration(int64(delay) * int64(attempt+1))
	case "fixed":
		// delay remains the same
	}

	// Cap at max delay
	if delay > policy.MaxDelay {
		delay = policy.MaxDelay
	}

	// Add jitter if enabled
	if policy.Jitter {
		jitterAmount := time.Duration(float64(delay) * 0.1) // 10% jitter
		delay += time.Duration(float64(jitterAmount) * (2*rand.Float64() - 1))
	}

	return delay
}

// Background workers

func (ehe *ErrorHandlingEnhancer) startErrorLogger() {
	ehe.errorLogger.wg.Add(1)
	defer ehe.errorLogger.wg.Done()

	for {
		select {
		case entry := <-ehe.errorLogger.buffer:
			ehe.writeLogEntry(entry)
		case <-ehe.errorLogger.shutdown:
			return
		}
	}
}

func (ehe *ErrorHandlingEnhancer) writeLogEntry(entry ErrorLogEntry) {
	// In production, this would write to files, databases, or external services
	log.Printf("[%s] %s: %s", entry.Level.String(), entry.Component, entry.Message)

	if entry.Context != nil && len(entry.Context) > 0 {
		log.Printf("Context: %v", entry.Context)
	}

	if entry.StackTrace != "" {
		log.Printf("Stack trace: %s", entry.StackTrace)
	}
}

func (ehe *ErrorHandlingEnhancer) startAggregator() {
	ticker := time.NewTicker(ehe.errorAggregator.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ehe.cleanupErrorGroups()
		}
	}
}

func (ehe *ErrorHandlingEnhancer) cleanupErrorGroups() {
	ehe.errorAggregator.mutex.Lock()
	defer ehe.errorAggregator.mutex.Unlock()

	cutoff := time.Now().Add(-ehe.errorAggregator.timeWindow)

	for pattern, group := range ehe.errorAggregator.errorGroups {
		if group.LastSeen.Before(cutoff) {
			delete(ehe.errorAggregator.errorGroups, pattern)
		}
	}
}

func (ehe *ErrorHandlingEnhancer) startAlertManager() {
	for {
		select {
		case alert := <-ehe.alertManager.alerts:
			ehe.processAlert(alert)
		}
	}
}

func (ehe *ErrorHandlingEnhancer) processAlert(alert Alert) {
	log.Printf("ALERT [%s]: %s", alert.Severity, alert.Message)

	// In production, this would send to monitoring systems, Slack, email, etc.
	subscribers, exists := ehe.alertManager.subscribers[alert.Type]
	if exists {
		for _, subscriber := range subscribers {
			if err := subscriber.HandleAlert(alert); err != nil {
				log.Printf("Failed to send alert to subscriber: %v", err)
			}
		}
	}
}

func (ehe *ErrorHandlingEnhancer) startMetricsCollector() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ehe.updateErrorRate()
		}
	}
}

func (ehe *ErrorHandlingEnhancer) updateErrorRate() {
	ehe.metrics.mutex.Lock()
	defer ehe.metrics.mutex.Unlock()

	// Calculate error rate (simplified)
	totalRequests := int64(1000) // This would come from request metrics
	if totalRequests > 0 {
		ehe.metrics.ErrorRate = float64(ehe.metrics.TotalErrors) / float64(totalRequests)
	}
}

// Utility functions

func contains(str string, patterns []string) bool {
	for _, pattern := range patterns {
		if strings.Contains(strings.ToLower(str), strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}

func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	case LogLevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Circuit breaker methods

func (cb *ErrorCircuitBreaker) ShouldBreak(err error) bool {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	now := time.Now()

	switch cb.state {
	case CircuitClosed:
		cb.failureCount++
		cb.lastFailure = now

		if cb.failureCount >= cb.config.FailureThreshold {
			cb.state = CircuitOpen
			log.Printf("Circuit breaker '%s' opened after %d failures", cb.name, cb.failureCount)
			return true
		}

	case CircuitOpen:
		if now.Sub(cb.lastFailure) > cb.config.OpenTimeout {
			cb.state = CircuitHalfOpen
			cb.successCount = 0
			log.Printf("Circuit breaker '%s' transitioning to half-open", cb.name)
		} else {
			return true
		}

	case CircuitHalfOpen:
		cb.failureCount++
		cb.lastFailure = now
		cb.state = CircuitOpen
		log.Printf("Circuit breaker '%s' failed in half-open state, reopening", cb.name)
		return true
	}

	return false
}

// GetMetrics returns current error metrics
func (ehe *ErrorHandlingEnhancer) GetMetrics() *ErrorMetrics {
	ehe.metrics.mutex.RLock()
	defer ehe.metrics.mutex.RUnlock()

	metricsCopy := *ehe.metrics
	return &metricsCopy
}

// Shutdown gracefully shuts down the error handling enhancer
func (ehe *ErrorHandlingEnhancer) Shutdown(ctx context.Context) error {
	log.Println("Shutting down error handling enhancer...")

	close(ehe.errorLogger.shutdown)
	ehe.errorLogger.wg.Wait()

	log.Println("Error handling enhancer shutdown completed")
	return nil
}
