import { startTransaction } from './errorTracking';

/**
 * Performance monitoring utilities
 * Tracks Core Web Vitals and custom metrics
 */

// Core Web Vitals thresholds (in milliseconds)
const THRESHOLDS = {
  FCP: { good: 1800, needsImprovement: 3000 }, // First Contentful Paint
  LCP: { good: 2500, needsImprovement: 4000 }, // Largest Contentful Paint
  FID: { good: 100, needsImprovement: 300 },   // First Input Delay
  CLS: { good: 0.1, needsImprovement: 0.25 },  // Cumulative Layout Shift
  TTFB: { good: 800, needsImprovement: 1800 }, // Time to First Byte
  INP: { good: 200, needsImprovement: 500 },   // Interaction to Next Paint
};

// Performance observer instance
let performanceObserver = null;

/**
 * Initialize performance monitoring
 */
export function initPerformanceMonitoring() {
  if (typeof window === 'undefined' || !('PerformanceObserver' in window)) {
    return;
  }

  // Observe Core Web Vitals
  observeLCP();
  observeFID();
  observeCLS();
  observeFCP();
  observeTTFB();
  observeINP();
  
  // Monitor long tasks
  observeLongTasks();
  
  // Monitor resource loading
  observeResourceTiming();
}

/**
 * Observe Largest Contentful Paint
 */
function observeLCP() {
  try {
    const observer = new PerformanceObserver((list) => {
      const entries = list.getEntries();
      const lastEntry = entries[entries.length - 1];
      const value = lastEntry.renderTime || lastEntry.loadTime;
      
      reportMetric('LCP', value, THRESHOLDS.LCP);
    });
    
    observer.observe({ entryTypes: ['largest-contentful-paint'] });
  } catch (e) {
    console.warn('LCP observer not supported');
  }
}

/**
 * Observe First Input Delay
 */
function observeFID() {
  try {
    const observer = new PerformanceObserver((list) => {
      const entries = list.getEntries();
      entries.forEach((entry) => {
        const value = entry.processingStart - entry.startTime;
        reportMetric('FID', value, THRESHOLDS.FID);
      });
    });
    
    observer.observe({ entryTypes: ['first-input'] });
  } catch (e) {
    console.warn('FID observer not supported');
  }
}

/**
 * Observe Cumulative Layout Shift
 */
function observeCLS() {
  let clsValue = 0;
  let clsEntries = [];
  
  try {
    const observer = new PerformanceObserver((list) => {
      const entries = list.getEntries();
      entries.forEach((entry) => {
        if (!entry.hadRecentInput) {
          clsValue += entry.value;
          clsEntries.push(entry);
        }
      });
    });
    
    observer.observe({ entryTypes: ['layout-shift'] });
    
    // Report CLS when page is hidden
    document.addEventListener('visibilitychange', () => {
      if (document.visibilityState === 'hidden') {
        reportMetric('CLS', clsValue, THRESHOLDS.CLS);
      }
    });
  } catch (e) {
    console.warn('CLS observer not supported');
  }
}

/**
 * Observe First Contentful Paint
 */
function observeFCP() {
  try {
    const observer = new PerformanceObserver((list) => {
      const entries = list.getEntries();
      entries.forEach((entry) => {
        if (entry.name === 'first-contentful-paint') {
          reportMetric('FCP', entry.startTime, THRESHOLDS.FCP);
        }
      });
    });
    
    observer.observe({ entryTypes: ['paint'] });
  } catch (e) {
    console.warn('FCP observer not supported');
  }
}

/**
 * Observe Time to First Byte
 */
function observeTTFB() {
  try {
    const observer = new PerformanceObserver((list) => {
      const entries = list.getEntries();
      entries.forEach((entry) => {
        const ttfb = entry.responseStart - entry.requestStart;
        reportMetric('TTFB', ttfb, THRESHOLDS.TTFB);
      });
    });
    
    observer.observe({ entryTypes: ['navigation'] });
  } catch (e) {
    // Fallback for browsers that don't support navigation timing
    window.addEventListener('load', () => {
      const navEntry = performance.getEntriesByType('navigation')[0];
      if (navEntry) {
        const ttfb = navEntry.responseStart - navEntry.requestStart;
        reportMetric('TTFB', ttfb, THRESHOLDS.TTFB);
      }
    });
  }
}

/**
 * Observe Interaction to Next Paint (INP)
 */
function observeINP() {
  let maxDuration = 0;
  
  try {
    const observer = new PerformanceObserver((list) => {
      const entries = list.getEntries();
      entries.forEach((entry) => {
        if (entry.duration > maxDuration) {
          maxDuration = entry.duration;
        }
      });
    });
    
    observer.observe({ entryTypes: ['event'] });
    
    // Report INP when page is hidden
    document.addEventListener('visibilitychange', () => {
      if (document.visibilityState === 'hidden' && maxDuration > 0) {
        reportMetric('INP', maxDuration, THRESHOLDS.INP);
      }
    });
  } catch (e) {
    console.warn('INP observer not supported');
  }
}

/**
 * Observe long tasks (blocking main thread)
 */
function observeLongTasks() {
  try {
    const observer = new PerformanceObserver((list) => {
      const entries = list.getEntries();
      entries.forEach((entry) => {
        if (entry.duration > 50) { // Tasks longer than 50ms
          reportLongTask(entry);
        }
      });
    });
    
    observer.observe({ entryTypes: ['longtask'] });
  } catch (e) {
    console.warn('Long task observer not supported');
  }
}

/**
 * Observe resource timing
 */
function observeResourceTiming() {
  try {
    const observer = new PerformanceObserver((list) => {
      const entries = list.getEntries();
      entries.forEach((entry) => {
        if (entry.duration > 1000) { // Resources taking more than 1s
          reportSlowResource(entry);
        }
      });
    });
    
    observer.observe({ entryTypes: ['resource'] });
  } catch (e) {
    console.warn('Resource timing observer not supported');
  }
}

/**
 * Report metric to monitoring service
 */
function reportMetric(name, value, threshold) {
  const rating = getMetricRating(value, threshold);
  
  // Log in development
  if (process.env.NODE_ENV === 'development') {
    console.log(`[Performance] ${name}: ${value.toFixed(2)}ms (${rating})`);
  }
  
  // Send to monitoring service
  if (typeof window !== 'undefined' && window.gtag) {
    window.gtag('event', name, {
      event_category: 'Web Vitals',
      event_label: rating,
      value: Math.round(value),
      non_interaction: true,
    });
  }
  
  // Send to Sentry
  const transaction = startTransaction(`web-vital-${name}`, 'measure');
  if (transaction) {
    transaction.setData('value', value);
    transaction.setData('rating', rating);
    transaction.finish();
  }
}

/**
 * Get metric rating based on thresholds
 */
function getMetricRating(value, threshold) {
  if (value <= threshold.good) return 'good';
  if (value <= threshold.needsImprovement) return 'needs-improvement';
  return 'poor';
}

/**
 * Report long task
 */
function reportLongTask(entry) {
  if (process.env.NODE_ENV === 'development') {
    console.warn(`[Performance] Long task detected: ${entry.duration.toFixed(2)}ms`);
  }
  
  // You can send this to your monitoring service
}

/**
 * Report slow resource
 */
function reportSlowResource(entry) {
  if (process.env.NODE_ENV === 'development') {
    console.warn(`[Performance] Slow resource: ${entry.name} (${entry.duration.toFixed(2)}ms)`);
  }
  
  // You can send this to your monitoring service
}

/**
 * Custom performance marks and measures
 */
export const performanceMark = {
  start: (name) => {
    if (typeof window !== 'undefined' && window.performance) {
      window.performance.mark(`${name}-start`);
    }
  },
  
  end: (name) => {
    if (typeof window !== 'undefined' && window.performance) {
      window.performance.mark(`${name}-end`);
      
      try {
        window.performance.measure(name, `${name}-start`, `${name}-end`);
        const measure = window.performance.getEntriesByName(name, 'measure')[0];
        
        if (measure) {
          if (process.env.NODE_ENV === 'development') {
            console.log(`[Performance] ${name}: ${measure.duration.toFixed(2)}ms`);
          }
          
          return measure.duration;
        }
      } catch (e) {
        // Ignore if mark doesn't exist
      }
    }
    
    return null;
  },
  
  clear: (name) => {
    if (typeof window !== 'undefined' && window.performance) {
      window.performance.clearMarks(`${name}-start`);
      window.performance.clearMarks(`${name}-end`);
      window.performance.clearMeasures(name);
    }
  }
};

/**
 * Track component render performance
 */
export function trackComponentPerformance(componentName, renderTime) {
  if (renderTime > 16) { // More than one frame (16ms)
    if (process.env.NODE_ENV === 'development') {
      console.warn(`[Performance] Slow render: ${componentName} (${renderTime.toFixed(2)}ms)`);
    }
    
    // Send to monitoring
    const transaction = startTransaction(`component-render-${componentName}`, 'measure');
    if (transaction) {
      transaction.setData('renderTime', renderTime);
      transaction.setData('slow', renderTime > 16);
      transaction.finish();
    }
  }
}

// Auto-initialize on import
if (typeof window !== 'undefined') {
  // Wait for page load to ensure all resources are available
  if (document.readyState === 'complete') {
    initPerformanceMonitoring();
  } else {
    window.addEventListener('load', initPerformanceMonitoring);
  }
}