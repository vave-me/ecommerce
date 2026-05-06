package consciousness

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
	
	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
)

var (
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
	ErrQuotaExhausted    = errors.New("quota exhausted")
)

// RateLimiter provides rate limiting capabilities for the consciousness system
type RateLimiter struct {
	globalLimiter    *rate.Limiter
	componentLimiters map[string]*rate.Limiter
	userLimiters     map[string]*rate.Limiter
	quotaManager     *QuotaManager
	mu               sync.RWMutex
	logger           zerolog.Logger
	metrics          *MetricsCollector
}

// RateLimiterConfig holds configuration for rate limiting
type RateLimiterConfig struct {
	// Global limits
	GlobalRPS        int
	GlobalBurst      int
	
	// Component-specific limits
	ComponentLimits  map[string]ComponentLimit
	
	// User-specific limits
	DefaultUserRPS   int
	DefaultUserBurst int
	VIPUserRPS       int
	VIPUserBurst     int
	
	// Quota settings
	DailyQuotas      map[string]int64
	MonthlyQuotas    map[string]int64
}

// ComponentLimit defines rate limits for a specific component
type ComponentLimit struct {
	RPS   int
	Burst int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(config RateLimiterConfig, logger zerolog.Logger, metrics *MetricsCollector) *RateLimiter {
	rl := &RateLimiter{
		globalLimiter:     rate.NewLimiter(rate.Limit(config.GlobalRPS), config.GlobalBurst),
		componentLimiters: make(map[string]*rate.Limiter),
		userLimiters:      make(map[string]*rate.Limiter),
		logger:            logger,
		metrics:           metrics,
	}
	
	// Initialize component limiters
	for component, limits := range config.ComponentLimits {
		rl.componentLimiters[component] = rate.NewLimiter(
			rate.Limit(limits.RPS),
			limits.Burst,
		)
	}
	
	// Initialize quota manager
	rl.quotaManager = NewQuotaManager(config.DailyQuotas, config.MonthlyQuotas)
	
	return rl
}

// Allow checks if an operation is allowed under rate limits
func (rl *RateLimiter) Allow(ctx context.Context, component, operation, userID string) error {
	// Check global rate limit
	if !rl.globalLimiter.Allow() {
		rl.recordRateLimitExceeded("global", operation)
		return fmt.Errorf("%w: global limit", ErrRateLimitExceeded)
	}
	
	// Check component rate limit
	if err := rl.checkComponentLimit(component); err != nil {
		rl.recordRateLimitExceeded(component, operation)
		return err
	}
	
	// Check user rate limit
	if userID != "" {
		if err := rl.checkUserLimit(userID); err != nil {
			rl.recordRateLimitExceeded("user:"+userID, operation)
			return err
		}
	}
	
	// Check quotas
	if err := rl.quotaManager.CheckQuota(component, operation); err != nil {
		rl.recordQuotaExhausted(component, operation)
		return err
	}
	
	return nil
}

// Wait blocks until the operation is allowed
func (rl *RateLimiter) Wait(ctx context.Context, component, operation, userID string) error {
	// Wait for global rate limit
	if err := rl.globalLimiter.Wait(ctx); err != nil {
		return fmt.Errorf("global rate limit wait failed: %w", err)
	}
	
	// Wait for component rate limit
	if limiter := rl.getComponentLimiter(component); limiter != nil {
		if err := limiter.Wait(ctx); err != nil {
			return fmt.Errorf("component rate limit wait failed: %w", err)
		}
	}
	
	// Wait for user rate limit
	if userID != "" {
		if limiter := rl.getUserLimiter(userID); limiter != nil {
			if err := limiter.Wait(ctx); err != nil {
				return fmt.Errorf("user rate limit wait failed: %w", err)
			}
		}
	}
	
	// Check quotas (can't wait for these)
	if err := rl.quotaManager.CheckQuota(component, operation); err != nil {
		return err
	}
	
	return nil
}

// Reserve reserves tokens for future use
func (rl *RateLimiter) Reserve(component string, tokens int) *rate.Reservation {
	if limiter := rl.getComponentLimiter(component); limiter != nil {
		return limiter.ReserveN(time.Now(), tokens)
	}
	return rl.globalLimiter.ReserveN(time.Now(), tokens)
}

// checkComponentLimit checks component-specific rate limit
func (rl *RateLimiter) checkComponentLimit(component string) error {
	limiter := rl.getComponentLimiter(component)
	if limiter != nil && !limiter.Allow() {
		return fmt.Errorf("%w: component %s", ErrRateLimitExceeded, component)
	}
	return nil
}

// checkUserLimit checks user-specific rate limit
func (rl *RateLimiter) checkUserLimit(userID string) error {
	limiter := rl.getUserLimiter(userID)
	if limiter != nil && !limiter.Allow() {
		return fmt.Errorf("%w: user %s", ErrRateLimitExceeded, userID)
	}
	return nil
}

// getComponentLimiter gets or creates a component limiter
func (rl *RateLimiter) getComponentLimiter(component string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.componentLimiters[component]
	rl.mu.RUnlock()
	
	if exists {
		return limiter
	}
	
	// Create default limiter if not exists
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	// Double-check after acquiring write lock
	if limiter, exists := rl.componentLimiters[component]; exists {
		return limiter
	}
	
	// Create with default limits (10 RPS, burst of 20)
	limiter = rate.NewLimiter(10, 20)
	rl.componentLimiters[component] = limiter
	return limiter
}

// getUserLimiter gets or creates a user limiter
func (rl *RateLimiter) getUserLimiter(userID string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.userLimiters[userID]
	rl.mu.RUnlock()
	
	if exists {
		return limiter
	}
	
	// Create default limiter if not exists
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	// Double-check after acquiring write lock
	if limiter, exists := rl.userLimiters[userID]; exists {
		return limiter
	}
	
	// Create with default user limits (5 RPS, burst of 10)
	limiter = rate.NewLimiter(5, 10)
	rl.userLimiters[userID] = limiter
	return limiter
}

// UpdateUserLimits updates rate limits for a specific user
func (rl *RateLimiter) UpdateUserLimits(userID string, rps int, burst int) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	rl.userLimiters[userID] = rate.NewLimiter(rate.Limit(rps), burst)
	
	rl.logger.Info().
		Str("user_id", userID).
		Int("rps", rps).
		Int("burst", burst).
		Msg("Updated user rate limits")
}

// recordRateLimitExceeded records rate limit exceeded metrics
func (rl *RateLimiter) recordRateLimitExceeded(limiterType, operation string) {
	if rl.metrics != nil {
		rl.metrics.RecordCustomMetric(
			fmt.Sprintf("rate_limit_exceeded_%s_%s", limiterType, operation),
			1,
		)
	}
	
	rl.logger.Warn().
		Str("limiter_type", limiterType).
		Str("operation", operation).
		Msg("Rate limit exceeded")
}

// recordQuotaExhausted records quota exhausted metrics
func (rl *RateLimiter) recordQuotaExhausted(component, operation string) {
	if rl.metrics != nil {
		rl.metrics.RecordCustomMetric(
			fmt.Sprintf("quota_exhausted_%s_%s", component, operation),
			1,
		)
	}
	
	rl.logger.Warn().
		Str("component", component).
		Str("operation", operation).
		Msg("Quota exhausted")
}

// QuotaManager manages usage quotas
type QuotaManager struct {
	dailyQuotas   map[string]int64
	monthlyQuotas map[string]int64
	usage         map[string]*UsageTracker
	mu            sync.RWMutex
}

// UsageTracker tracks usage over time periods
type UsageTracker struct {
	Daily   *TimeWindowCounter
	Monthly *TimeWindowCounter
}

// TimeWindowCounter counts events within a time window
type TimeWindowCounter struct {
	count      int64
	windowStart time.Time
	windowDuration time.Duration
	mu         sync.Mutex
}

// NewQuotaManager creates a new quota manager
func NewQuotaManager(dailyQuotas, monthlyQuotas map[string]int64) *QuotaManager {
	qm := &QuotaManager{
		dailyQuotas:   dailyQuotas,
		monthlyQuotas: monthlyQuotas,
		usage:         make(map[string]*UsageTracker),
	}
	
	// Initialize usage trackers
	for key := range dailyQuotas {
		qm.usage[key] = &UsageTracker{
			Daily:   NewTimeWindowCounter(24 * time.Hour),
			Monthly: NewTimeWindowCounter(30 * 24 * time.Hour),
		}
	}
	
	return qm
}

// CheckQuota checks if an operation is within quota limits
func (qm *QuotaManager) CheckQuota(component, operation string) error {
	key := fmt.Sprintf("%s:%s", component, operation)
	
	qm.mu.RLock()
	tracker, exists := qm.usage[key]
	dailyLimit, hasDailyLimit := qm.dailyQuotas[key]
	monthlyLimit, hasMonthlyLimit := qm.monthlyQuotas[key]
	qm.mu.RUnlock()
	
	if !exists {
		// No quota tracking for this operation
		return nil
	}
	
	// Check daily quota
	if hasDailyLimit {
		dailyUsage := tracker.Daily.Count()
		if dailyUsage >= dailyLimit {
			return fmt.Errorf("%w: daily quota for %s", ErrQuotaExhausted, key)
		}
	}
	
	// Check monthly quota
	if hasMonthlyLimit {
		monthlyUsage := tracker.Monthly.Count()
		if monthlyUsage >= monthlyLimit {
			return fmt.Errorf("%w: monthly quota for %s", ErrQuotaExhausted, key)
		}
	}
	
	// Increment usage
	tracker.Daily.Increment()
	tracker.Monthly.Increment()
	
	return nil
}

// NewTimeWindowCounter creates a new time window counter
func NewTimeWindowCounter(windowDuration time.Duration) *TimeWindowCounter {
	return &TimeWindowCounter{
		windowStart:    time.Now(),
		windowDuration: windowDuration,
	}
}

// Count returns the current count, resetting if window has expired
func (twc *TimeWindowCounter) Count() int64 {
	twc.mu.Lock()
	defer twc.mu.Unlock()
	
	now := time.Now()
	if now.Sub(twc.windowStart) > twc.windowDuration {
		// Window has expired, reset
		twc.count = 0
		twc.windowStart = now
	}
	
	return twc.count
}

// Increment increments the counter
func (twc *TimeWindowCounter) Increment() {
	twc.mu.Lock()
	defer twc.mu.Unlock()
	
	now := time.Now()
	if now.Sub(twc.windowStart) > twc.windowDuration {
		// Window has expired, reset
		twc.count = 1
		twc.windowStart = now
	} else {
		twc.count++
	}
}

// GetUsageStats returns current usage statistics
func (qm *QuotaManager) GetUsageStats() map[string]UsageStats {
	qm.mu.RLock()
	defer qm.mu.RUnlock()
	
	stats := make(map[string]UsageStats)
	
	for key, tracker := range qm.usage {
		dailyLimit, _ := qm.dailyQuotas[key]
		monthlyLimit, _ := qm.monthlyQuotas[key]
		
		stats[key] = UsageStats{
			DailyUsage:    tracker.Daily.Count(),
			DailyLimit:    dailyLimit,
			MonthlyUsage:  tracker.Monthly.Count(),
			MonthlyLimit:  monthlyLimit,
		}
	}
	
	return stats
}

// UsageStats represents usage statistics
type UsageStats struct {
	DailyUsage   int64
	DailyLimit   int64
	MonthlyUsage int64
	MonthlyLimit int64
}