package consciousness

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	
	"github.com/rs/zerolog"
)

// ShutdownManager handles graceful shutdown of the consciousness system
type ShutdownManager struct {
	logger          zerolog.Logger
	shutdownTimeout time.Duration
	cleanupFuncs    []CleanupFunc
	mu              sync.Mutex
	shutdownChan    chan struct{}
	wg              sync.WaitGroup
	isShuttingDown  bool
}

// CleanupFunc represents a cleanup function
type CleanupFunc struct {
	Name     string
	Func     func(context.Context) error
	Priority int // Lower number = higher priority
	Timeout  time.Duration
}

// NewShutdownManager creates a new shutdown manager
func NewShutdownManager(logger zerolog.Logger, shutdownTimeout time.Duration) *ShutdownManager {
	return &ShutdownManager{
		logger:          logger,
		shutdownTimeout: shutdownTimeout,
		cleanupFuncs:    make([]CleanupFunc, 0),
		shutdownChan:    make(chan struct{}),
	}
}

// RegisterCleanup registers a cleanup function
func (sm *ShutdownManager) RegisterCleanup(name string, priority int, timeout time.Duration, fn func(context.Context) error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	if sm.isShuttingDown {
		sm.logger.Warn().
			Str("cleanup_name", name).
			Msg("Cannot register cleanup during shutdown")
		return
	}
	
	cleanup := CleanupFunc{
		Name:     name,
		Func:     fn,
		Priority: priority,
		Timeout:  timeout,
	}
	
	// Insert in priority order
	inserted := false
	for i, existing := range sm.cleanupFuncs {
		if cleanup.Priority < existing.Priority {
			sm.cleanupFuncs = append(sm.cleanupFuncs[:i], append([]CleanupFunc{cleanup}, sm.cleanupFuncs[i:]...)...)
			inserted = true
			break
		}
	}
	
	if !inserted {
		sm.cleanupFuncs = append(sm.cleanupFuncs, cleanup)
	}
	
	sm.logger.Info().
		Str("cleanup_name", name).
		Int("priority", priority).
		Dur("timeout", timeout).
		Msg("Registered cleanup function")
}

// StartSignalHandler starts listening for shutdown signals
func (sm *ShutdownManager) StartSignalHandler(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	
	go func() {
		select {
		case sig := <-sigChan:
			sm.logger.Info().
				Str("signal", sig.String()).
				Msg("Received shutdown signal")
			sm.Shutdown()
			cancel()
		case <-ctx.Done():
			sm.logger.Info().Msg("Context cancelled")
			sm.Shutdown()
		}
	}()
	
	return ctx
}

// Shutdown initiates graceful shutdown
func (sm *ShutdownManager) Shutdown() {
	sm.mu.Lock()
	if sm.isShuttingDown {
		sm.mu.Unlock()
		sm.logger.Warn().Msg("Shutdown already in progress")
		return
	}
	sm.isShuttingDown = true
	sm.mu.Unlock()
	
	sm.logger.Info().
		Dur("timeout", sm.shutdownTimeout).
		Msg("Starting graceful shutdown")
	
	// Close shutdown channel to notify waiters
	close(sm.shutdownChan)
	
	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), sm.shutdownTimeout)
	defer cancel()
	
	// Run cleanup functions
	sm.runCleanup(ctx)
	
	// Wait for all goroutines to finish
	done := make(chan struct{})
	go func() {
		sm.wg.Wait()
		close(done)
	}()
	
	select {
	case <-done:
		sm.logger.Info().Msg("Graceful shutdown completed")
	case <-ctx.Done():
		sm.logger.Error().
			Dur("timeout", sm.shutdownTimeout).
			Msg("Shutdown timeout exceeded, forcing exit")
	}
}

// runCleanup runs all cleanup functions in priority order
func (sm *ShutdownManager) runCleanup(ctx context.Context) {
	sm.mu.Lock()
	cleanupFuncs := make([]CleanupFunc, len(sm.cleanupFuncs))
	copy(cleanupFuncs, sm.cleanupFuncs)
	sm.mu.Unlock()
	
	for _, cleanup := range cleanupFuncs {
		// Check if context is already cancelled
		select {
		case <-ctx.Done():
			sm.logger.Error().
				Str("cleanup_name", cleanup.Name).
				Msg("Shutdown context cancelled, skipping remaining cleanup")
			return
		default:
		}
		
		// Create timeout context for this cleanup
		cleanupCtx, cancel := context.WithTimeout(ctx, cleanup.Timeout)
		
		sm.logger.Info().
			Str("cleanup_name", cleanup.Name).
			Dur("timeout", cleanup.Timeout).
			Msg("Running cleanup function")
		
		startTime := time.Now()
		err := cleanup.Func(cleanupCtx)
		duration := time.Since(startTime)
		
		if err != nil {
			sm.logger.Error().
				Err(err).
				Str("cleanup_name", cleanup.Name).
				Dur("duration", duration).
				Msg("Cleanup function failed")
		} else {
			sm.logger.Info().
				Str("cleanup_name", cleanup.Name).
				Dur("duration", duration).
				Msg("Cleanup function completed")
		}
		
		cancel()
	}
}

// WaitForShutdown returns a channel that's closed when shutdown begins
func (sm *ShutdownManager) WaitForShutdown() <-chan struct{} {
	return sm.shutdownChan
}

// AddGoroutine increments the wait group for tracking goroutines
func (sm *ShutdownManager) AddGoroutine() {
	sm.wg.Add(1)
}

// DoneGoroutine decrements the wait group
func (sm *ShutdownManager) DoneGoroutine() {
	sm.wg.Done()
}

// IsShuttingDown returns true if shutdown is in progress
func (sm *ShutdownManager) IsShuttingDown() bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	return sm.isShuttingDown
}

// ConsciousnessShutdownHandler creates shutdown handlers for consciousness components
type ConsciousnessShutdownHandler struct {
	manager         *ConsciousnessManager
	errorHandler    *ErrorHandler
	metricsCollector *MetricsCollector
	rateLimiter     *RateLimiter
	circuitBreakers *CircuitBreakerManager
	logger          zerolog.Logger
}

// NewConsciousnessShutdownHandler creates a new consciousness shutdown handler
func NewConsciousnessShutdownHandler(
	manager *ConsciousnessManager,
	errorHandler *ErrorHandler,
	metrics *MetricsCollector,
	rateLimiter *RateLimiter,
	circuitBreakers *CircuitBreakerManager,
	logger zerolog.Logger,
) *ConsciousnessShutdownHandler {
	return &ConsciousnessShutdownHandler{
		manager:         manager,
		errorHandler:    errorHandler,
		metricsCollector: metrics,
		rateLimiter:     rateLimiter,
		circuitBreakers: circuitBreakers,
		logger:          logger,
	}
}

// RegisterShutdownHandlers registers all consciousness shutdown handlers
func (csh *ConsciousnessShutdownHandler) RegisterShutdownHandlers(sm *ShutdownManager) {
	// Priority 1: Stop accepting new events
	sm.RegisterCleanup("stop_event_processing", 1, 5*time.Second, func(ctx context.Context) error {
		csh.logger.Info().Msg("Stopping event processing")
		// Implementation would stop event listeners
		return nil
	})
	
	// Priority 2: Complete in-flight operations
	sm.RegisterCleanup("complete_operations", 2, 30*time.Second, func(ctx context.Context) error {
		csh.logger.Info().Msg("Waiting for in-flight operations to complete")
		// Wait for current decisions and actions to complete
		status := csh.manager.GetStatus()
		csh.logger.Info().
			Int64("events_processed", status.EventsProcessed).
			Int64("decisions_made", status.DecisionsMade).
			Msg("Final consciousness statistics")
		return nil
	})
	
	// Priority 3: Flush metrics
	sm.RegisterCleanup("flush_metrics", 3, 5*time.Second, func(ctx context.Context) error {
		csh.logger.Info().Msg("Flushing metrics")
		// Flush any buffered metrics
		return nil
	})
	
	// Priority 4: Close circuit breakers
	sm.RegisterCleanup("close_circuit_breakers", 4, 5*time.Second, func(ctx context.Context) error {
		csh.logger.Info().Msg("Closing circuit breakers")
		metrics := csh.circuitBreakers.GetAllMetrics()
		for _, metric := range metrics {
			csh.logger.Info().
				Str("circuit", metric.Name).
				Str("state", metric.State).
				Int("failures", metric.Failures).
				Msg("Circuit breaker final state")
		}
		return nil
	})
	
	// Priority 5: Save state
	sm.RegisterCleanup("save_state", 5, 10*time.Second, func(ctx context.Context) error {
		csh.logger.Info().Msg("Saving consciousness state")
		// Save any persistent state
		return nil
	})
	
	// Priority 6: Close connections
	sm.RegisterCleanup("close_connections", 6, 5*time.Second, func(ctx context.Context) error {
		csh.logger.Info().Msg("Closing connections")
		// Close database connections, message queues, etc.
		return nil
	})
}

// GracefulShutdownMiddleware creates middleware that respects shutdown
func GracefulShutdownMiddleware(sm *ShutdownManager) func(next func(context.Context) error) func(context.Context) error {
	return func(next func(context.Context) error) func(context.Context) error {
		return func(ctx context.Context) error {
			// Check if we're shutting down
			if sm.IsShuttingDown() {
				return fmt.Errorf("service is shutting down")
			}
			
			// Add to wait group
			sm.AddGoroutine()
			defer sm.DoneGoroutine()
			
			// Create context that's cancelled on shutdown
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			
			// Watch for shutdown
			go func() {
				select {
				case <-sm.WaitForShutdown():
					cancel()
				case <-ctx.Done():
				}
			}()
			
			return next(ctx)
		}
	}
}