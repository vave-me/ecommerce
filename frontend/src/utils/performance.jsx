/**
 * Consolidated performance monitoring utilities
 */
// Store metrics for Web Vitals
const metrics = {
  fcp: [], // First Contentful Paint
  lcp: [], // Largest Contentful Paint
  cls: [], // Cumulative Layout Shift
  fid: [], // First Input Delay
  ttfb: [], // Time to First Byte
  inp: [], // Interaction to Next Paint
};
// Store performance measurements
const performanceStore = {
  renderCounts: {},
  timings: {},
  measurements: new Map(),
};
// Store component render counts
const renderCounts = new Map();
// Store component-specific performance metrics
const componentMetrics = new Map();
/**
 * Initialize monitoring in the app
 * - Patches React lifecycle methods to track renders
 * - Sets up performance observers
 */
export function initializeMonitoring() {
  if (process.env.NODE_ENV !== 'development') return;
  if (process.env.NODE_ENV === 'development') {
    }
  // Set up PerformanceObserver for long tasks
  if (typeof window !== 'undefined' && 'PerformanceObserver' in window) {
    try {
      const observer = new PerformanceObserver((list) => {
        list.getEntries().forEach((entry) => {
          // Log long tasks (over 50ms)
          if (entry.duration > 50 && process.env.NODE_ENV === 'development') {
          }
        });
      });
      observer.observe({ entryTypes: ['longtask'] });
    } catch (e) {
      if (process.env.NODE_ENV === 'development') {
      }
    }
  }
  // Expose monitoring functions globally for dev tools
  if (typeof window !== 'undefined') {
    window.__PERFORMANCE_MONITOR__ = {
      generateReport: generatePerformanceReport,
      getRenderCounts: () => Object.fromEntries(renderCounts),
      resetRenderCounts,
      clearData: () => {
        performanceStore.renderCounts = {};
        performanceStore.timings = {};
        performanceStore.measurements.clear();
        renderCounts.clear();
        componentMetrics.clear();
        if (process.env.NODE_ENV === 'development') {
          }
      },
      enableVerboseLogging: () => {
        if (process.env.NODE_ENV === 'development') {
        }
        window.__PERFORMANCE_VERBOSE = true;
      },
      disableVerboseLogging: () => {
        if (process.env.NODE_ENV === 'development') {
        }
        window.__PERFORMANCE_VERBOSE = false;
      }
    };
    if (process.env.NODE_ENV === 'development') {
      }
  }
}
/**
 * Hook to track component render counts
 * @param {string} componentName - Name of the component to track
 * @returns {number} Current render count
 */
export function useRenderCount(componentName) {
  if (process.env.NODE_ENV !== 'development') {
    return 0; // No-op in production
  }
  if (!renderCounts.has(componentName)) {
    renderCounts.set(componentName, 0);
  }
  const count = renderCounts.get(componentName) + 1;
  renderCounts.set(componentName, count);
  // Log every 5 renders or on the first render
  if ((count === 1 || count % 5 === 0) && process.env.NODE_ENV === 'development') {
  }
  return count;
}
/**
 * Track component render count (alternative implementation)
 * @param {string} componentName - Name of the component
 */
export function trackRender(componentName) {
  if (process.env.NODE_ENV !== 'development') return;
  if (!performanceStore.renderCounts[componentName]) {
    performanceStore.renderCounts[componentName] = 0;
  }
  performanceStore.renderCounts[componentName]++;
  // Log every 5 renders to avoid console spam
  if (performanceStore.renderCounts[componentName] % 5 === 0 && process.env.NODE_ENV === 'development') {
    }
}
/**
 * Reset render counts for all or specific components
 * @param {string} [componentName] - Optional component name to reset
 */
export function resetRenderCounts(componentName) {
  if (componentName) {
    renderCounts.set(componentName, 0);
  } else {
    renderCounts.clear();
  }
}
/**
 * Time a function execution and log its performance
 * @param {Function} fn - Function to time
 * @param {string} label - Label for the log
 * @returns {any} Result of the function
 */
export function timeExecution(fn, label) {
  if (process.env.NODE_ENV !== 'development') {
    return fn();
  }
  const start = performance.now();
  const result = fn();
  const end = performance.now();
  return result;
}
/**
 * Measure execution time of a function with console.time
 * @param {Function} fn - Function to measure
 * @param {string} label - Label for the measurement
 * @returns {any} - Result of the function
 */
export function measureExecutionTime(fn, label) {
  if (process.env.NODE_ENV !== 'development') return fn();
  if (process.env.NODE_ENV === 'development') {
    }
  const result = fn();
  if (process.env.NODE_ENV === 'development') {
    }
  return result;
}
/**
 * Start measuring a performance timing
 * @param {string} label - Label for the measurement
 */
export function startMeasure(label) {
  if (process.env.NODE_ENV !== 'development') return;
  performanceStore.timings[label] = performance.now();
}
/**
 * End measuring a performance timing and log the result
 * @param {string} label - Label for the measurement
 * @returns {number} - Duration in milliseconds
 */
export function endMeasure(label) {
  if (process.env.NODE_ENV !== 'development') return 0;
  const startTime = performanceStore.timings[label];
  if (!startTime) {
    if (process.env.NODE_ENV === 'development') {
    }
    return 0;
  }
  const duration = performance.now() - startTime;
  if (process.env.NODE_ENV === 'development') {
  }
  // Store the measurement for reporting
  if (!performanceStore.measurements.has(label)) {
    performanceStore.measurements.set(label, []);
  }
  performanceStore.measurements.get(label).push(duration);
  return duration;
}
/**
 * Create a performance report for all tracked metrics
 * @returns {Object} - Performance report
 */
export function generatePerformanceReport() {
  if (process.env.NODE_ENV !== 'development') return {};
  const report = {
    timestamp: new Date().toISOString(),
    componentRenders: { 
      ...performanceStore.renderCounts,
      ...Object.fromEntries(renderCounts)
    },
    measurements: {},
    webVitals: getAverageMetrics(),
    componentMetrics: Object.fromEntries(componentMetrics)
  };
  // Calculate stats for each measurement
  performanceStore.measurements.forEach((durations, label) => {
    if (durations.length === 0) return;
    const sum = durations.reduce((acc, val) => acc + val, 0);
    const avg = sum / durations.length;
    const min = Math.min(...durations);
    const max = Math.max(...durations);
    report.measurements[label] = {
      average: avg.toFixed(2),
      min: min.toFixed(2),
      max: max.toFixed(2),
      samples: durations.length,
    };
  });
  if (process.env.NODE_ENV === 'development') {
    }
  return report;
}
/**
 * Create a performance wrapper for a component for debugging purposes
 * @param {React.Component} Component - Component to wrap
 * @param {string} name - Display name for the component
 * @returns {React.Component} Wrapped component with performance logging
 */
export function withPerformanceTracking(Component, name) {
  if (process.env.NODE_ENV !== 'development') {
    return Component; // No-op in production
  }
  // Get display name for the component
  const displayName = name || Component.displayName || Component.name || 'Unknown';
  // Create wrapper component
  const WrappedComponent = (props) => {
    const start = performance.now();
    const result = <Component {...props} />;
    const end = performance.now();
    // Only log if render takes more than 5ms
    if (end - start > 5) {
    }
    return result;
  };
  WrappedComponent.displayName = `Performance(${displayName})`;
  return WrappedComponent;
}
/**
 * Track a specific user interaction for performance monitoring
 * @param {string} name - Name of the interaction
 * @param {Function} callback - The interaction callback to measure
 * @returns {Function} - Wrapped function that measures performance
 */
export function trackInteraction(name, callback) {
  return (...args) => {
    const start = performance.now();
    try {
      return callback(...args);
    } finally {
      const duration = performance.now() - start;
      // Log long interactions (over 100ms)
      if (duration > 100) {
      }
      // Record to PerformanceObserver
      performance.mark(`${name}-end`);
    }
  };
}
/**
 * React hook for component performance monitoring
 * @param {string} componentName - Name of the component to monitor
 * @returns {Object} - Performance monitoring utilities
 */
export function usePerformance(componentName) {
  if (process.env.NODE_ENV !== 'development') {
    // Return no-op functions in production
    return {
      trackApiCall: async (apiName, apiCall) => await apiCall(),
      trackError: () => {},
      getMetrics: () => ({})
    };
  }
  // Create a new metrics object if one doesn't exist
  if (!componentMetrics.has(componentName)) {
    componentMetrics.set(componentName, {
      renderCount: 0,
      totalRenderTime: 0,
      apiCalls: {},
      errors: []
    });
  }
  // Record start time for component render
  const startTime = performance.now();
  // Update metrics when the component rerenders
  const metrics = componentMetrics.get(componentName);
  metrics.renderCount++;
  // Calculate render time and update metrics on next tick
  setTimeout(() => {
    const renderTime = performance.now() - startTime;
    metrics.totalRenderTime += renderTime;
  }, 0);
  // Track API calls
  const trackApiCall = async (apiName, apiCall) => {
    const start = performance.now();
    try {
      const result = await apiCall();
      const duration = performance.now() - start;
      // Initialize array if needed
      if (!metrics.apiCalls[apiName]) {
        metrics.apiCalls[apiName] = [];
      }
      metrics.apiCalls[apiName].push(duration);
      return result;
    } catch (error) {
      metrics.errors.push({
        type: 'api',
        apiName,
        error: error.message,
        timestamp: new Date().toISOString()
      });
      throw error;
    }
  };
  // Track errors
  const trackError = (error) => {
    metrics.errors.push({
      type: 'component',
      error: error.message,
      stack: error.stack,
      timestamp: new Date().toISOString()
    });
  };
  return {
    trackApiCall,
    trackError,
    getMetrics: () => ({
      renderCount: metrics.renderCount,
      averageRenderTime: metrics.totalRenderTime / metrics.renderCount,
      apiCalls: {...metrics.apiCalls},
      errors: [...metrics.errors]
    })
  };
}
/**
 * Capture and process Web Vitals metrics
 * @param {Object} metric - Web Vitals metric object
 */
export function captureWebVitals(metric) {
  // Store metric for local analysis
  if (metrics[metric.name]) {
    metrics[metric.name].push(metric.value);
  }
  // Log to console in development
  if (process.env.NODE_ENV === 'development') {
    }
  // Send to analytics service
  try {
    // Use fetch with keepalive to ensure data is sent even during page transitions
    fetch('/api/vitals', {
      method: 'POST',
      body: JSON.stringify(metric),
      keepalive: true,
      headers: {
        'Content-Type': 'application/json'
      }
    }).catch(() => {
      // Silently fail - don't interrupt user experience for analytics
    });
  } catch (error) {
    // No-op in production, log in development
    if (process.env.NODE_ENV === 'development') {
    }
  }
}
/**
 * Hook to monitor Web Vitals (can be used with Next.js reportWebVitals)
 * Alias for captureWebVitals for naming consistency with Next.js
 * @param {Object} metric - Web Vitals metric
 */
export function reportWebVitals(metric) {
  // Only run in development or when specifically enabled
  if (process.env.NODE_ENV !== 'development' && 
      !process.env.NEXT_PUBLIC_TRACK_VITALS) return;
  captureWebVitals(metric);
}
/**
 * Get average values for each metric
 * @returns {Object} Average values for stored metrics
 */
export function getAverageMetrics() {
  const result = {};
  Object.keys(metrics).forEach(key => {
    const values = metrics[key];
    if (values.length === 0) return;
    const sum = values.reduce((a, b) => a + b, 0);
    result[key] = sum / values.length;
  });
  return result;
}
/**
 * Clear stored web vitals metrics
 */
export function clearMetrics() {
  Object.keys(metrics).forEach(key => {
    metrics[key] = [];
  });
}
// Initialize monitoring if we're in the browser
if (typeof window !== 'undefined' && process.env.NODE_ENV === 'development') {
  // Auto-initialize when imported
  setTimeout(() => {
    initializeMonitoring();
  }, 0);
} 