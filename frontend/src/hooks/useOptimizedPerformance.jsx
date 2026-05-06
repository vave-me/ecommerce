/**
 * OPTIMIZED PERFORMANCE MONITORING
 * Lightweight production-ready performance tracking
 * Replaces heavy performance monitoring with minimal overhead
 */
import { useEffect, useRef, useCallback } from 'react';
/**
 * Lightweight performance monitoring hook
 * Only tracks essential metrics with minimal overhead
 */
export const useOptimizedPerformance = (componentName) => {
  const mountTimeRef = useRef(null);
  const renderCountRef = useRef(0);
  const errorCountRef = useRef(0);
  const lastErrorRef = useRef(null);
  // Track component lifecycle
  useEffect(() => {
    mountTimeRef.current = performance.now();
    renderCountRef.current = 0;
    // Only enable detailed tracking in development
    if (process.env.NODE_ENV === 'development') {
      }
    return () => {
      if (process.env.NODE_ENV === 'development') {
        const totalTime = performance.now() - mountTimeRef.current;
      }
    };
  }, [componentName]);
  // Track renders (lightweight)
  useEffect(() => {
    renderCountRef.current += 1;
    // Only warn about excessive renders in development
    if (process.env.NODE_ENV === 'development' && renderCountRef.current > 20) {
    }
  });
  // Lightweight error tracking
  const trackError = useCallback((error) => {
    errorCountRef.current += 1;
    lastErrorRef.current = {
      message: error.message,
      timestamp: Date.now(),
      count: errorCountRef.current
    };
    if (process.env.NODE_ENV === 'development') {
    }
  }, [componentName]);
  // Performance metrics (minimal)
  const getMetrics = useCallback(() => {
    const uptime = mountTimeRef.current ? performance.now() - mountTimeRef.current : 0;
    return {
      componentName,
      uptime: Math.round(uptime),
      renderCount: renderCountRef.current,
      errorCount: errorCountRef.current,
      lastError: lastErrorRef.current,
      avgRenderTime: renderCountRef.current > 0 ? Math.round(uptime / renderCountRef.current) : 0
    };
  }, [componentName]);
  // Only return essential functions
  return {
    trackError,
    getMetrics
  };
};
/**
 * Lightweight performance wrapper for components
 * Minimal overhead alternative to heavy performance tracking
 */
export const withOptimizedPerformance = (Component, componentName) => {
  const WrappedComponent = (props) => {
    const { trackError } = useOptimizedPerformance(componentName);
    // Wrap component in error boundary
    try {
      return <Component {...props} />;
    } catch (error) {
      trackError(error);
      // In development, show error details
      if (process.env.NODE_ENV === 'development') {
        return (
          <div style={{ 
            padding: '20px', 
            border: '2px solid red', 
            background: '#ffe6e6',
            color: '#cc0000',
            fontFamily: 'monospace'
          }}>
            <h4>Component Error: {componentName}</h4>
            <p>{error.message}</p>
          </div>
        );
      }
      // In production, return null or fallback
      return null;
    }
  };
  WrappedComponent.displayName = `OptimizedPerformance(${componentName})`;
  return WrappedComponent;
};
/**
 * Global performance utilities (production-safe)
 */
export const performanceUtils = {
  // Mark performance milestones
  mark: (name) => {
    if (typeof performance !== 'undefined' && performance.mark) {
      performance.mark(name);
    }
  },
  // Measure between marks
  measure: (name, start, end) => {
    if (typeof performance !== 'undefined' && performance.measure) {
      try {
        performance.measure(name, start, end);
        const measure = performance.getEntriesByName(name, 'measure')[0];
        if (process.env.NODE_ENV === 'development' && measure) {
        }
        return measure?.duration || 0;
      } catch (error) {
        // Silently handle errors in production
        if (process.env.NODE_ENV === 'development') {
        }
        return 0;
      }
    }
    return 0;
  },
  // Clear performance entries (cleanup)
  clearEntries: () => {
    if (typeof performance !== 'undefined') {
      try {
        performance.clearMarks();
        performance.clearMeasures();
      } catch (error) {
        // Form submission error
        if (process.env.NODE_ENV === 'development') {
            console.error('Form submission error:', error);
        }
        // Could set error state here if available
        throw error;
    }
    }
  }
};
export default useOptimizedPerformance; 