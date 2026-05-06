"use client";
import { useEffect, useCallback, useRef } from 'react';
/**
 * Advanced performance monitoring hook for Core Web Vitals and component performance
 * Tracks: LCP, FID, CLS, TTFB, FCP, and custom metrics
 */
export function usePerformanceMonitoring(componentName = 'Unknown') {
  const renderStartTime = useRef(performance.now());
  const metricsRef = useRef({});
  // Track component render performance
  const trackRender = useCallback((phase = 'mount') => {
    const renderTime = performance.now() - renderStartTime.current;
    if (process.env.NODE_ENV === 'development') {
    }
    // Store metrics for analysis
    metricsRef.current[`${componentName}_${phase}`] = renderTime;
    // Reset timer for next measurement
    renderStartTime.current = performance.now();
  }, [componentName]);
  // Track Core Web Vitals
  const trackWebVitals = useCallback(() => {
    if (typeof window === 'undefined') return;
    // Largest Contentful Paint (LCP)
    const observeLCP = () => {
      const observer = new PerformanceObserver((list) => {
        const entries = list.getEntries();
        const lastEntry = entries[entries.length - 1];
        if (lastEntry) {
          const lcp = lastEntry.startTime;
          metricsRef.current.lcp = lcp;
          if (process.env.NODE_ENV === 'development') {
          }
        }
      });
      try {
        observer.observe({ entryTypes: ['largest-contentful-paint'] });
      } catch (e) {
        if (process.env.NODE_ENV === 'development') {
        }
      }
    };
    // First Input Delay (FID)
    const observeFID = () => {
      const observer = new PerformanceObserver((list) => {
        const entries = list.getEntries();
        entries.forEach((entry) => {
          const fid = entry.processingStart - entry.startTime;
          metricsRef.current.fid = fid;
          if (process.env.NODE_ENV === 'development') {
          }
        });
      });
      try {
        observer.observe({ entryTypes: ['first-input'] });
      } catch (e) {
        if (process.env.NODE_ENV === 'development') {
        }
      }
    };
    // Cumulative Layout Shift (CLS)
    const observeCLS = () => {
      let clsValue = 0;
      const observer = new PerformanceObserver((list) => {
        const entries = list.getEntries();
        entries.forEach((entry) => {
          if (!entry.hadRecentInput) {
            clsValue += entry.value;
            metricsRef.current.cls = clsValue;
          }
        });
      });
      try {
        observer.observe({ entryTypes: ['layout-shift'] });
      } catch (e) {
        if (process.env.NODE_ENV === 'development') {
        }
      }
    };
    // Time to First Byte (TTFB)
    const observeTTFB = () => {
      const observer = new PerformanceObserver((list) => {
        const entries = list.getEntries();
        entries.forEach((entry) => {
          if (entry.entryType === 'navigation') {
            const ttfb = entry.responseStart - entry.requestStart;
            metricsRef.current.ttfb = ttfb;
          }
        });
      });
      try {
        observer.observe({ entryTypes: ['navigation'] });
      } catch (e) {
        if (process.env.NODE_ENV === 'development') {
        }
      }
    };
    // Initialize all observers
    observeLCP();
    observeFID();
    observeCLS();
    observeTTFB();
  }, []);
  // Memory usage tracking
  const trackMemoryUsage = useCallback(() => {
    if (typeof window === 'undefined' || !performance.memory) return;
    const memory = performance.memory;
    const memoryInfo = {
      usedJSHeapSize: (memory.usedJSHeapSize / 1048576).toFixed(2), // MB
      totalJSHeapSize: (memory.totalJSHeapSize / 1048576).toFixed(2), // MB
      jsHeapSizeLimit: (memory.jsHeapSizeLimit / 1048576).toFixed(2), // MB
    };
    metricsRef.current.memory = memoryInfo;
  }, []);
  // Get performance insights
  const getPerformanceInsights = useCallback(() => {
    const metrics = metricsRef.current;
    const insights = [];
    // LCP insights
    if (metrics.lcp) {
      if (metrics.lcp > 4000) {
        insights.push({
          metric: 'LCP',
          severity: 'high',
          message: 'Largest Contentful Paint is poor. Consider optimizing images and reducing server response times.',
          value: `${metrics.lcp.toFixed(2)}ms`
        });
      } else if (metrics.lcp > 2500) {
        insights.push({
          metric: 'LCP',
          severity: 'medium',
          message: 'Largest Contentful Paint needs improvement. Optimize critical resources.',
          value: `${metrics.lcp.toFixed(2)}ms`
        });
      }
    }
    // FID insights
    if (metrics.fid) {
      if (metrics.fid > 300) {
        insights.push({
          metric: 'FID',
          severity: 'high',
          message: 'First Input Delay is poor. Reduce JavaScript execution time.',
          value: `${metrics.fid.toFixed(2)}ms`
        });
      } else if (metrics.fid > 100) {
        insights.push({
          metric: 'FID',
          severity: 'medium',
          message: 'First Input Delay needs improvement. Consider code splitting.',
          value: `${metrics.fid.toFixed(2)}ms`
        });
      }
    }
    // CLS insights
    if (metrics.cls) {
      if (metrics.cls > 0.25) {
        insights.push({
          metric: 'CLS',
          severity: 'high',
          message: 'Cumulative Layout Shift is poor. Reserve space for dynamic content.',
          value: metrics.cls.toFixed(4)
        });
      } else if (metrics.cls > 0.1) {
        insights.push({
          metric: 'CLS',
          severity: 'medium',
          message: 'Cumulative Layout Shift needs improvement. Avoid layout shifts.',
          value: metrics.cls.toFixed(4)
        });
      }
    }
    return insights;
  }, []);
  // Initialize monitoring
  useEffect(() => {
    trackWebVitals();
    trackMemoryUsage();
    // Track initial render
    trackRender('mount');
    // Set up periodic memory tracking (optimized for production)
    const memoryInterval = setInterval(trackMemoryUsage, 60000); // Every 60 seconds (optimized)
    return () => {
      clearInterval(memoryInterval);
      trackRender('unmount');
    };
  }, [trackWebVitals, trackMemoryUsage, trackRender]);
  return {
    trackRender,
    getMetrics: () => metricsRef.current,
    getPerformanceInsights,
    trackMemoryUsage
  };
}
/**
 * HOC for automatic performance tracking
 */
export function withPerformanceTracking(WrappedComponent, componentName) {
  return function PerformanceTrackedComponent(props) {
    const { trackRender } = usePerformanceMonitoring(componentName);
    useEffect(() => {
      trackRender('render');
    });
    return <WrappedComponent {...props} />;
  };
} 