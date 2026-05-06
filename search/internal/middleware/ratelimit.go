package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// RateLimitConfig defines rate limiting configuration for different methods
type RateLimitConfig struct {
	// Global limits
	GlobalRequestsPerSecond int
	GlobalBurst             int

	// Per-method limits
	MethodLimits map[string]MethodLimit

	// Default limits for unspecified methods
	DefaultAuthenticatedLimit   RateLimit
	DefaultUnauthenticatedLimit RateLimit

	// Enable/disable rate limiting
	Enabled bool
}

// MethodLimit defines rate limits for a specific method
type MethodLimit struct {
	Authenticated   RateLimit
	Unauthenticated RateLimit
	ComplexityMultiplier float64 // Multiply base rate by query complexity
}

// RateLimit defines rate limiting parameters
type RateLimit struct {
	RequestsPerMinute int
	Burst             int
}

// RateLimiter manages rate limiting for the service
type RateLimiter struct {
	config    RateLimitConfig
	global    *rate.Limiter
	userLimiters map[string]*rate.Limiter
	ipLimiters   map[string]*rate.Limiter
	mu           sync.RWMutex
	
	// Cleanup goroutine management
	cleanupInterval time.Duration
	stopCleanup     chan struct{}
}

// NewRateLimiter creates a new rate limiter with the given configuration
func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	if !config.Enabled {
		return &RateLimiter{config: config}
	}

	rl := &RateLimiter{
		config:          config,
		global:          rate.NewLimiter(rate.Limit(config.GlobalRequestsPerSecond), config.GlobalBurst),
		userLimiters:    make(map[string]*rate.Limiter),
		ipLimiters:      make(map[string]*rate.Limiter),
		cleanupInterval: 5 * time.Minute,
		stopCleanup:     make(chan struct{}),
	}

	// Start cleanup goroutine
	go rl.cleanupRoutine()

	return rl
}

// UnaryServerInterceptor creates a gRPC unary interceptor for rate limiting
func (rl *RateLimiter) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if !rl.config.Enabled {
			return handler(ctx, req)
		}

		// Check global rate limit first
		if !rl.global.Allow() {
			return nil, status.Errorf(codes.ResourceExhausted, "global rate limit exceeded")
		}

		// Extract rate limit key (user ID or IP)
		key, isAuthenticated := rl.extractRateLimitKey(ctx)
		
		// Get method-specific limits
		limit := rl.getMethodLimit(info.FullMethod, isAuthenticated)
		if limit.RequestsPerMinute == 0 {
			// Method not allowed for this auth state
			return nil, status.Errorf(codes.PermissionDenied, "method not allowed for unauthenticated users")
		}

		// Calculate query complexity if applicable
		complexity := rl.calculateQueryComplexity(req, info.FullMethod)
		adjustedLimit := rl.adjustLimitByComplexity(limit, complexity)

		// Check per-key rate limit
		limiter := rl.getLimiter(key, isAuthenticated, adjustedLimit)
		if !limiter.Allow() {
			return nil, status.Errorf(codes.ResourceExhausted, 
				"rate limit exceeded for method %s, please retry after %v", 
				info.FullMethod, rl.getRetryAfter(limiter))
		}

		return handler(ctx, req)
	}
}

// HTTPMiddleware creates an HTTP middleware for rate limiting REST endpoints
func (rl *RateLimiter) HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.config.Enabled {
			next.ServeHTTP(w, r)
			return
		}

		// Check global rate limit
		if !rl.global.Allow() {
			rl.writeRateLimitResponse(w, "global rate limit exceeded")
			return
		}

		// Extract rate limit key
		key, isAuthenticated := rl.extractHTTPRateLimitKey(r)
		
		// Map HTTP path to gRPC method (simplified - you may need to enhance this)
		method := rl.mapHTTPToGRPCMethod(r.URL.Path)
		
		// Get method-specific limits
		limit := rl.getMethodLimit(method, isAuthenticated)
		if limit.RequestsPerMinute == 0 {
			http.Error(w, "method not allowed for unauthenticated users", http.StatusForbidden)
			return
		}

		// Check per-key rate limit
		limiter := rl.getLimiter(key, isAuthenticated, limit)
		if !limiter.Allow() {
			rl.writeRateLimitResponse(w, fmt.Sprintf("rate limit exceeded, retry after %v", rl.getRetryAfter(limiter)))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// extractRateLimitKey extracts the rate limiting key from gRPC context
func (rl *RateLimiter) extractRateLimitKey(ctx context.Context) (string, bool) {
	// First try to get user ID from context (set by auth middleware)
	if userID, ok := ctx.Value("user_id").(string); ok && userID != "" {
		return "user:" + userID, true
	}

	// Fall back to IP address
	if p, ok := peer.FromContext(ctx); ok {
		return "ip:" + p.Addr.String(), false
	}

	// Last resort - use a generic key
	return "unknown", false
}

// extractHTTPRateLimitKey extracts the rate limiting key from HTTP request
func (rl *RateLimiter) extractHTTPRateLimitKey(r *http.Request) (string, bool) {
	// Check for user ID in headers (set by auth middleware)
	if userID := r.Header.Get("X-User-ID"); userID != "" {
		return "user:" + userID, true
	}

	// Fall back to IP address
	ip := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ip = strings.Split(forwarded, ",")[0]
	}
	
	return "ip:" + ip, false
}

// getMethodLimit returns the rate limit for a specific method
func (rl *RateLimiter) getMethodLimit(method string, isAuthenticated bool) RateLimit {
	if methodLimit, ok := rl.config.MethodLimits[method]; ok {
		if isAuthenticated {
			return methodLimit.Authenticated
		}
		return methodLimit.Unauthenticated
	}

	// Return default limits
	if isAuthenticated {
		return rl.config.DefaultAuthenticatedLimit
	}
	return rl.config.DefaultUnauthenticatedLimit
}

// getLimiter returns or creates a rate limiter for the given key
func (rl *RateLimiter) getLimiter(key string, isAuthenticated bool, limit RateLimit) *rate.Limiter {
	rl.mu.RLock()
	limiters := rl.ipLimiters
	if isAuthenticated {
		limiters = rl.userLimiters
	}
	
	if limiter, ok := limiters[key]; ok {
		rl.mu.RUnlock()
		return limiter
	}
	rl.mu.RUnlock()

	// Create new limiter
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	// Double-check after acquiring write lock
	if limiter, ok := limiters[key]; ok {
		return limiter
	}

	limiter := rate.NewLimiter(rate.Limit(float64(limit.RequestsPerMinute)/60.0), limit.Burst)
	limiters[key] = limiter
	return limiter
}

// calculateQueryComplexity calculates the complexity score for a request
func (rl *RateLimiter) calculateQueryComplexity(req interface{}, method string) float64 {
	// This is a simplified version - you should implement based on your actual request types
	complexity := 1.0
	
	// Example complexity calculation
	switch method {
	case "/searchpb.SearchService/UnifiedSearch":
		// Add complexity based on request parameters
		// This would require type assertion to actual request types
		complexity = 2.0 // Base complexity for unified search
		
	case "/searchpb.SearchService/UnifiedFeed":
		complexity = 1.5
		
	case "/searchpb.SearchService/BatchUnifiedFeed":
		complexity = 3.0
	}
	
	return complexity
}

// adjustLimitByComplexity adjusts the rate limit based on query complexity
func (rl *RateLimiter) adjustLimitByComplexity(limit RateLimit, complexity float64) RateLimit {
	if complexity <= 1.0 {
		return limit
	}
	
	// Reduce allowed rate based on complexity
	adjusted := RateLimit{
		RequestsPerMinute: int(float64(limit.RequestsPerMinute) / complexity),
		Burst:             int(float64(limit.Burst) / complexity),
	}
	
	// Ensure minimum values
	if adjusted.RequestsPerMinute < 1 {
		adjusted.RequestsPerMinute = 1
	}
	if adjusted.Burst < 1 {
		adjusted.Burst = 1
	}
	
	return adjusted
}

// getRetryAfter calculates when the client should retry
func (rl *RateLimiter) getRetryAfter(limiter *rate.Limiter) time.Duration {
	// Get the next allowed time
	r := limiter.Reserve()
	delay := r.Delay()
	r.Cancel() // Cancel the reservation
	
	if delay == 0 {
		// If no delay, suggest 1 second
		return time.Second
	}
	return delay
}

// writeRateLimitResponse writes a rate limit error response with proper headers
func (rl *RateLimiter) writeRateLimitResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.config.GlobalRequestsPerSecond))
	w.Header().Set("X-RateLimit-Remaining", "0")
	w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Minute).Unix()))
	w.Header().Set("Retry-After", "60")
	
	w.WriteHeader(http.StatusTooManyRequests)
	fmt.Fprintf(w, `{"error": "%s"}`, message)
}

// mapHTTPToGRPCMethod maps HTTP paths to gRPC method names
func (rl *RateLimiter) mapHTTPToGRPCMethod(path string) string {
	// This is a simplified mapping - enhance based on your actual routes
	mappings := map[string]string{
		"/v1/search/unified":        "/searchpb.SearchService/UnifiedSearch",
		"/v1/feed/unified":          "/searchpb.SearchService/UnifiedFeed",
		"/v1/products/search":       "/searchpb.SearchService/SearchProductsWithFilters",
		"/v1/posts/search":          "/searchpb.SearchService/SearchPostsWithFilters",
	}
	
	for prefix, method := range mappings {
		if strings.HasPrefix(path, prefix) {
			return method
		}
	}
	
	return path
}

// cleanupRoutine periodically cleans up old rate limiters
func (rl *RateLimiter) cleanupRoutine() {
	ticker := time.NewTicker(rl.cleanupInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			rl.cleanup()
		case <-rl.stopCleanup:
			return
		}
	}
}

// cleanup removes inactive rate limiters
func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	// In a production system, you'd want to track last access time
	// and remove limiters that haven't been used recently
	// For now, we'll just limit the total number
	
	maxLimiters := 10000
	
	if len(rl.userLimiters) > maxLimiters {
		// Remove oldest entries (this is simplified - use LRU in production)
		count := 0
		for key := range rl.userLimiters {
			delete(rl.userLimiters, key)
			count++
			if len(rl.userLimiters) <= maxLimiters/2 {
				break
			}
		}
	}
	
	if len(rl.ipLimiters) > maxLimiters {
		count := 0
		for key := range rl.ipLimiters {
			delete(rl.ipLimiters, key)
			count++
			if len(rl.ipLimiters) <= maxLimiters/2 {
				break
			}
		}
	}
}

// Stop gracefully stops the rate limiter
func (rl *RateLimiter) Stop() {
	close(rl.stopCleanup)
}

// GetDefaultConfig returns a default rate limiting configuration
func GetDefaultConfig() RateLimitConfig {
	return RateLimitConfig{
		Enabled:                 true,
		GlobalRequestsPerSecond: 1000,
		GlobalBurst:             100,
		
		MethodLimits: map[string]MethodLimit{
			// Critical - Resource intensive operations
			"/searchpb.SearchService/UnifiedSearch": {
				Authenticated:   RateLimit{RequestsPerMinute: 10, Burst: 2},
				Unauthenticated: RateLimit{RequestsPerMinute: 5, Burst: 1},
				ComplexityMultiplier: 2.0,
			},
			"/searchpb.SearchService/UnifiedFeed": {
				Authenticated:   RateLimit{RequestsPerMinute: 20, Burst: 5},
				Unauthenticated: RateLimit{RequestsPerMinute: 10, Burst: 2},
				ComplexityMultiplier: 1.5,
			},
			"/searchpb.SearchService/BatchUnifiedFeed": {
				Authenticated:   RateLimit{RequestsPerMinute: 5, Burst: 1},
				Unauthenticated: RateLimit{RequestsPerMinute: 0, Burst: 0}, // Disabled
				ComplexityMultiplier: 3.0,
			},
			
			// Moderate - Standard searches
			"/searchpb.SearchService/SearchProductsWithFilters": {
				Authenticated:   RateLimit{RequestsPerMinute: 30, Burst: 10},
				Unauthenticated: RateLimit{RequestsPerMinute: 15, Burst: 5},
				ComplexityMultiplier: 1.0,
			},
			"/searchpb.SearchService/SearchPostsWithFilters": {
				Authenticated:   RateLimit{RequestsPerMinute: 30, Burst: 10},
				Unauthenticated: RateLimit{RequestsPerMinute: 15, Burst: 5},
				ComplexityMultiplier: 1.0,
			},
			
			// Light - Simple lookups
			"/searchpb.SearchService/GetProduct": {
				Authenticated:   RateLimit{RequestsPerMinute: 100, Burst: 20},
				Unauthenticated: RateLimit{RequestsPerMinute: 50, Burst: 10},
				ComplexityMultiplier: 0.5,
			},
			"/searchpb.SearchService/GetPost": {
				Authenticated:   RateLimit{RequestsPerMinute: 100, Burst: 20},
				Unauthenticated: RateLimit{RequestsPerMinute: 50, Burst: 10},
				ComplexityMultiplier: 0.5,
			},
		},
		
		// Default limits for unspecified methods
		DefaultAuthenticatedLimit:   RateLimit{RequestsPerMinute: 60, Burst: 10},
		DefaultUnauthenticatedLimit: RateLimit{RequestsPerMinute: 30, Burst: 5},
	}
}