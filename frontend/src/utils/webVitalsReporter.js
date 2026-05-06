/**
 * WEB VITALS REPORTER
 * Performance monitoring and analytics integration
 * 
 * Tracks Core Web Vitals for performance optimization
 */
/**
 * Web Vitals Reporter Utility
 * Reports Core Web Vitals metrics for performance monitoring
 */
// Import web-vitals library if available
let webVitals;
try {
    webVitals = require('web-vitals');
} catch (error) {
        // Form submission error
        if (process.env.NODE_ENV === 'development') {
            console.error('Form submission error:', error);
        }
        // Could set error state here if available
        throw error;
    }
/**
 * Report Web Vitals metrics
 * @param {Function} onPerfEntry - Callback function to handle metrics
 */
const reportWebVitals = (onPerfEntry) => {
    if (onPerfEntry && onPerfEntry instanceof Function && webVitals) {
        // Report Core Web Vitals
        webVitals.getCLS(onPerfEntry);  // Cumulative Layout Shift
        webVitals.getFID(onPerfEntry);  // First Input Delay
        webVitals.getFCP(onPerfEntry);  // First Contentful Paint
        webVitals.getLCP(onPerfEntry);  // Largest Contentful Paint
        webVitals.getTTFB(onPerfEntry); // Time to First Byte
    } else if (!webVitals) {
        // Fallback: Use Performance API directly
        if (typeof window !== 'undefined' && 'performance' in window) {
            const observer = new PerformanceObserver((list) => {
                list.getEntries().forEach((entry) => {
                    if (onPerfEntry && onPerfEntry instanceof Function) {
                        onPerfEntry({
                            name: entry.name,
                            value: entry.duration || entry.value,
                            id: entry.entryType,
                            delta: entry.duration || entry.value
                        });
                    }
                });
            });
            try {
                observer.observe({ entryTypes: ['navigation', 'paint', 'largest-contentful-paint'] });
            } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
        }
    }
};
/**
 * Default performance entry handler
 * Logs metrics to console in development
 */
const defaultHandler = (metric) => {
    if (process.env.NODE_ENV === 'development') {
        }
    // In production, you might want to send metrics to analytics service
    // Example: analytics.track('web-vital', metric);
};
// Auto-report with default handler if this file is imported directly
if (typeof window !== 'undefined') {
    reportWebVitals(defaultHandler);
}
export { reportWebVitals, defaultHandler };
/**
 * Initialize Web Vitals reporting
 * Dynamically imports web-vitals library only when needed
 */
export const initWebVitals = async () => {
  try {
    // Dynamic import to reduce initial bundle size
    const { getCLS, getFID, getFCP, getLCP, getTTFB } = await import('web-vitals');
    // Report all Core Web Vitals
    getCLS(reportWebVitals);
    getFID(reportWebVitals);
    getFCP(reportWebVitals);
    getLCP(reportWebVitals);
    getTTFB(reportWebVitals);
    if (process.env.NODE_ENV === 'development') {
      }
  } catch (error) {
    // Gracefully handle if web-vitals is not installed
    if (process.env.NODE_ENV === 'development') {
    }
  }
};
/**
 * Performance observer for custom metrics
 */
export const createPerformanceObserver = () => {
  if (typeof window === 'undefined' || !window.PerformanceObserver) {
    return null;
  }
  const observer = new PerformanceObserver((list) => {
    list.getEntries().forEach((entry) => {
      // Report custom performance metrics
      if (process.env.NODE_ENV === 'development') {
        }
      // Send to analytics in production
      if (process.env.NODE_ENV === 'production' && window.gtag) {
        window.gtag('event', 'custom_performance', {
          event_category: 'Performance',
          event_label: entry.name,
          value: Math.round(entry.duration || entry.processingStart || 0),
          custom_map: {
            entry_type: entry.entryType,
            start_time: entry.startTime
          }
        });
      }
    });
  });
  try {
    // Observe various performance entry types
    observer.observe({ 
      entryTypes: ['navigation', 'paint', 'largest-contentful-paint', 'first-input', 'layout-shift'] 
    });
    return observer;
  } catch (error) {
    if (process.env.NODE_ENV === 'development') {
    }
    return null;
  }
};
export default {
  reportWebVitals,
  initWebVitals,
  createPerformanceObserver
}; 