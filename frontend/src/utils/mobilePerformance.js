/**
 * Mobile-First Performance Monitoring
 * Production-ready Core Web Vitals and mobile-specific metrics tracking
 */
// Core Web Vitals thresholds (mobile-optimized)
const MOBILE_THRESHOLDS = {
  // Largest Contentful Paint (mobile)
  LCP: {
    good: 2500,    // 2.5s for mobile
    poor: 4000     // 4.0s for mobile
  },
  // First Input Delay (mobile)
  FID: {
    good: 100,     // 100ms
    poor: 300      // 300ms
  },
  // Cumulative Layout Shift (mobile)
  CLS: {
    good: 0.1,     // 0.1
    poor: 0.25     // 0.25
  },
  // Interaction to Next Paint (mobile)
  INP: {
    good: 200,     // 200ms for mobile
    poor: 500      // 500ms for mobile
  },
  // Time to First Byte (mobile)
  TTFB: {
    good: 800,     // 800ms for mobile
    poor: 1800     // 1.8s for mobile
  },
  // First Contentful Paint (mobile)
  FCP: {
    good: 1800,    // 1.8s for mobile
    poor: 3000     // 3.0s for mobile
  }
};
// Mobile-specific metrics storage
const mobileMetrics = {
  vitals: new Map(),
  navigation: new Map(),
  resources: new Map(),
  network: new Map(),
  battery: new Map(),
  memory: new Map()
};
// Performance observer instances
let observers = new Map();
// helper for conditional logging
const perfLog = (...args) => {
  if (process.env.NEXT_PUBLIC_MOBILE_PERF_DEBUG === 'true') {
    // eslint-disable-next-line no-console
    }
};
/**
 * Initialize mobile performance monitoring
 */
export function initMobilePerformanceMonitoring() {
  if (typeof window === 'undefined') return;
  perfLog('[Mobile Perf] Initializing mobile performance monitoring...');
  // 1. Setup Core Web Vitals monitoring
  setupCoreWebVitals();
  // 2. Setup mobile-specific monitoring
  setupMobileSpecificMetrics();
  // 3. Setup network monitoring
  setupNetworkMonitoring();
  // 4. Setup battery monitoring (if available)
  setupBatteryMonitoring();
  // 5. Setup memory monitoring
  setupMemoryMonitoring();
  // 6. Setup touch interaction monitoring
  setupTouchInteractionMonitoring();
  // 7. Setup viewport monitoring
  setupViewportMonitoring();
  // 8. Setup periodic reporting
  setupPeriodicReporting();
  perfLog('[Mobile Perf] Mobile performance monitoring initialized');
}
/**
 * Setup Core Web Vitals monitoring with mobile optimizations
 */
function setupCoreWebVitals() {
  // LCP (Largest Contentful Paint)
  if ('PerformanceObserver' in window) {
    try {
      const lcpObserver = new PerformanceObserver((list) => {
        const entries = list.getEntries();
        const lastEntry = entries[entries.length - 1];
        const metric = {
          name: 'LCP',
          value: lastEntry.startTime,
          rating: getRating(lastEntry.startTime, MOBILE_THRESHOLDS.LCP),
          timestamp: Date.now(),
          url: window.location.href,
          isMobile: isMobileDevice()
        };
        mobileMetrics.vitals.set('LCP', metric);
        reportMetric(metric);
      });
      lcpObserver.observe({ entryTypes: ['largest-contentful-paint'] });
      observers.set('LCP', lcpObserver);
    } catch (e) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', e);
        }
    }
  }
  // FID (First Input Delay)
  if ('PerformanceObserver' in window) {
    try {
      const fidObserver = new PerformanceObserver((list) => {
        list.getEntries().forEach((entry) => {
          const metric = {
            name: 'FID',
            value: entry.processingStart - entry.startTime,
            rating: getRating(entry.processingStart - entry.startTime, MOBILE_THRESHOLDS.FID),
            timestamp: Date.now(),
            url: window.location.href,
            isMobile: isMobileDevice()
          };
          mobileMetrics.vitals.set('FID', metric);
          reportMetric(metric);
        });
      });
      fidObserver.observe({ entryTypes: ['first-input'] });
      observers.set('FID', fidObserver);
    } catch (e) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', e);
        }
    }
  }
  // CLS (Cumulative Layout Shift)
  if ('PerformanceObserver' in window) {
    try {
      let clsValue = 0;
      let clsEntries = [];
      const clsObserver = new PerformanceObserver((list) => {
        list.getEntries().forEach((entry) => {
          if (!entry.hadRecentInput) {
            clsValue += entry.value;
            clsEntries.push(entry);
          }
        });
        const metric = {
          name: 'CLS',
          value: clsValue,
          rating: getRating(clsValue, MOBILE_THRESHOLDS.CLS),
          timestamp: Date.now(),
          url: window.location.href,
          isMobile: isMobileDevice(),
          entries: clsEntries.length
        };
        mobileMetrics.vitals.set('CLS', metric);
        reportMetric(metric);
      });
      clsObserver.observe({ entryTypes: ['layout-shift'] });
      observers.set('CLS', clsObserver);
    } catch (e) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', e);
        }
    }
  }
  // INP (Interaction to Next Paint) - newer metric
  if ('PerformanceObserver' in window) {
    try {
      const inpObserver = new PerformanceObserver((list) => {
        list.getEntries().forEach((entry) => {
          const metric = {
            name: 'INP',
            value: entry.duration,
            rating: getRating(entry.duration, MOBILE_THRESHOLDS.INP),
            timestamp: Date.now(),
            url: window.location.href,
            isMobile: isMobileDevice(),
            interactionType: entry.name
          };
          mobileMetrics.vitals.set('INP', metric);
          reportMetric(metric);
        });
      });
      inpObserver.observe({ entryTypes: ['event'] });
      observers.set('INP', inpObserver);
    } catch (e) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', e);
        }
    }
  }
}
/**
 * Setup mobile-specific performance metrics
 */
function setupMobileSpecificMetrics() {
  // Navigation timing
  if ('performance' in window && performance.timing) {
    const navigation = performance.timing;
    const metrics = {
      dns: navigation.domainLookupEnd - navigation.domainLookupStart,
      tcp: navigation.connectEnd - navigation.connectStart,
      ssl: navigation.secureConnectionStart ? navigation.connectEnd - navigation.secureConnectionStart : 0,
      ttfb: navigation.responseStart - navigation.navigationStart,
      download: navigation.responseEnd - navigation.responseStart,
      domParse: navigation.domContentLoadedEventStart - navigation.responseEnd,
      domReady: navigation.domContentLoadedEventEnd - navigation.navigationStart,
      loadComplete: navigation.loadEventEnd - navigation.navigationStart
    };
    mobileMetrics.navigation.set('timing', {
      ...metrics,
      timestamp: Date.now(),
      isMobile: isMobileDevice()
    });
  }
  // Resource timing for mobile optimization
  if ('PerformanceObserver' in window) {
    try {
      const resourceObserver = new PerformanceObserver((list) => {
        list.getEntries().forEach((entry) => {
          // Focus on critical resources for mobile
          if (entry.initiatorType === 'img' || 
              entry.initiatorType === 'css' || 
              entry.initiatorType === 'script') {
            const metric = {
              name: entry.name,
              type: entry.initiatorType,
              duration: entry.duration,
              size: entry.transferSize || 0,
              cached: entry.transferSize === 0 && entry.decodedBodySize > 0,
              timestamp: Date.now(),
              isMobile: isMobileDevice()
            };
            mobileMetrics.resources.set(entry.name, metric);
            // Report slow resources on mobile
            if (entry.duration > 1000) {
              reportMetric({
                name: 'SLOW_RESOURCE',
                value: entry.duration,
                resource: entry.name,
                type: entry.initiatorType,
                ...metric
              });
            }
          }
        });
      });
      resourceObserver.observe({ entryTypes: ['resource'] });
      observers.set('resource', resourceObserver);
    } catch (e) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', e);
        }
    }
  }
}
/**
 * Setup network monitoring for mobile
 */
function setupNetworkMonitoring() {
  if ('connection' in navigator) {
    const connection = navigator.connection || navigator.mozConnection || navigator.webkitConnection;
    if (connection) {
      const networkMetric = {
        effectiveType: connection.effectiveType,
        downlink: connection.downlink,
        rtt: connection.rtt,
        saveData: connection.saveData,
        timestamp: Date.now()
      };
      mobileMetrics.network.set('initial', networkMetric);
      // Monitor network changes
      connection.addEventListener('change', () => {
        const updatedMetric = {
          effectiveType: connection.effectiveType,
          downlink: connection.downlink,
          rtt: connection.rtt,
          saveData: connection.saveData,
          timestamp: Date.now()
        };
        mobileMetrics.network.set('current', updatedMetric);
        reportMetric({
          name: 'NETWORK_CHANGE',
          ...updatedMetric
        });
      });
    }
  }
}
/**
 * Setup battery monitoring for mobile devices
 */
function setupBatteryMonitoring() {
  if ('getBattery' in navigator) {
    navigator.getBattery().then((battery) => {
      const batteryMetric = {
        level: battery.level,
        charging: battery.charging,
        chargingTime: battery.chargingTime,
        dischargingTime: battery.dischargingTime,
        timestamp: Date.now()
      };
      mobileMetrics.battery.set('initial', batteryMetric);
      // Monitor battery changes
      battery.addEventListener('levelchange', () => {
        const updatedMetric = {
          level: battery.level,
          charging: battery.charging,
          timestamp: Date.now()
        };
        mobileMetrics.battery.set('current', updatedMetric);
        // Report low battery impact on performance
        if (battery.level < 0.2) {
          reportMetric({
            name: 'LOW_BATTERY',
            level: battery.level,
            timestamp: Date.now()
          });
        }
      });
    }).catch(() => {
      // Battery API not available
    });
  }
}
/**
 * Setup memory monitoring
 */
function setupMemoryMonitoring() {
  if ('memory' in performance) {
    const memoryMetric = {
      used: performance.memory.usedJSHeapSize,
      total: performance.memory.totalJSHeapSize,
      limit: performance.memory.jsHeapSizeLimit,
      timestamp: Date.now()
    };
    mobileMetrics.memory.set('initial', memoryMetric);
    // Monitor memory usage periodically
    setInterval(() => {
      const currentMemory = {
        used: performance.memory.usedJSHeapSize,
        total: performance.memory.totalJSHeapSize,
        limit: performance.memory.jsHeapSizeLimit,
        timestamp: Date.now()
      };
      mobileMetrics.memory.set('current', currentMemory);
      // Report high memory usage
      const usagePercent = (currentMemory.used / currentMemory.limit) * 100;
      if (usagePercent > 80) {
        reportMetric({
          name: 'HIGH_MEMORY_USAGE',
          percentage: usagePercent,
          ...currentMemory
        });
      }
    }, 30000); // Check every 30 seconds
  }
}
/**
 * Setup touch interaction monitoring
 */
function setupTouchInteractionMonitoring() {
  let touchStartTime = 0;
  let touchInteractions = [];
  document.addEventListener('touchstart', (event) => {
    touchStartTime = performance.now();
  }, { passive: true });
  document.addEventListener('touchend', (event) => {
    if (touchStartTime > 0) {
      const duration = performance.now() - touchStartTime;
      touchInteractions.push({
        duration,
        timestamp: Date.now(),
        touches: event.changedTouches.length
      });
      // Report slow touch interactions
      if (duration > 100) {
        reportMetric({
          name: 'SLOW_TOUCH_INTERACTION',
          value: duration,
          timestamp: Date.now()
        });
      }
      touchStartTime = 0;
    }
  }, { passive: true });
}
/**
 * Setup viewport monitoring for mobile
 */
function setupViewportMonitoring() {
  let resizeTimeout;
  const reportViewport = () => {
    const viewport = {
      width: window.innerWidth,
      height: window.innerHeight,
      devicePixelRatio: window.devicePixelRatio || 1,
      orientation: screen.orientation ? screen.orientation.angle : 0,
      timestamp: Date.now()
    };
    reportMetric({
      name: 'VIEWPORT_CHANGE',
      ...viewport
    });
  };
  window.addEventListener('resize', () => {
    clearTimeout(resizeTimeout);
    resizeTimeout = setTimeout(reportViewport, 250);
  });
  window.addEventListener('orientationchange', () => {
    setTimeout(reportViewport, 500); // Delay for orientation change
  });
}
/**
 * Setup periodic reporting
 */
function setupPeriodicReporting() {
  // Report metrics every 30 seconds
  setInterval(() => {
    const report = generateMobilePerformanceReport();
    if (report.metrics.length > 0) {
      reportMetric({
        name: 'PERIODIC_REPORT',
        report,
        timestamp: Date.now()
      });
    }
  }, 30000);
  // Report on page visibility change
  document.addEventListener('visibilitychange', () => {
    if (document.visibilityState === 'hidden') {
      const report = generateMobilePerformanceReport();
      reportMetric({
        name: 'PAGE_HIDDEN_REPORT',
        report,
        timestamp: Date.now()
      });
    }
  });
}
/**
 * Generate comprehensive mobile performance report
 */
export function generateMobilePerformanceReport() {
  if (typeof window === 'undefined' || typeof navigator === 'undefined') {
    return {
      timestamp: Date.now(),
      url: '',
      userAgent: '',
      isMobile: false,
      metrics: [],
      summary: {
        totalMetrics: 0,
        vitalsCount: 0,
        performanceScore: 0,
        recommendations: []
      }
    };
  }
  const report = {
    timestamp: Date.now(),
    url: window.location.href,
    userAgent: navigator.userAgent,
    isMobile: isMobileDevice(),
    metrics: [],
    summary: {}
  };
  // Add Core Web Vitals
  mobileMetrics.vitals.forEach((metric, name) => {
    report.metrics.push({
      category: 'vitals',
      name,
      ...metric
    });
  });
  // Add navigation metrics
  mobileMetrics.navigation.forEach((metric, name) => {
    report.metrics.push({
      category: 'navigation',
      name,
      ...metric
    });
  });
  // Add network metrics
  mobileMetrics.network.forEach((metric, name) => {
    report.metrics.push({
      category: 'network',
      name,
      ...metric
    });
  });
  // Add battery metrics
  mobileMetrics.battery.forEach((metric, name) => {
    report.metrics.push({
      category: 'battery',
      name,
      ...metric
    });
  });
  // Add memory metrics
  mobileMetrics.memory.forEach((metric, name) => {
    report.metrics.push({
      category: 'memory',
      name,
      ...metric
    });
  });
  // Generate summary
  report.summary = {
    totalMetrics: report.metrics.length,
    vitalsCount: Array.from(mobileMetrics.vitals.keys()).length,
    performanceScore: calculateMobilePerformanceScore(),
    recommendations: generateMobileRecommendations()
  };
  return report;
}
/**
 * Calculate mobile performance score (0-100)
 */
function calculateMobilePerformanceScore() {
  let score = 100;
  let penalties = 0;
  // Check Core Web Vitals
  mobileMetrics.vitals.forEach((metric) => {
    if (metric.rating === 'poor') {
      penalties += 20;
    } else if (metric.rating === 'needs-improvement') {
      penalties += 10;
    }
  });
  // Check network conditions
  const currentNetwork = mobileMetrics.network.get('current');
  if (currentNetwork && currentNetwork.effectiveType === 'slow-2g') {
    penalties += 15;
  } else if (currentNetwork && currentNetwork.effectiveType === '2g') {
    penalties += 10;
  }
  // Check battery level
  const currentBattery = mobileMetrics.battery.get('current');
  if (currentBattery && currentBattery.level < 0.2) {
    penalties += 5;
  }
  return Math.max(0, score - penalties);
}
/**
 * Generate mobile-specific performance recommendations
 */
function generateMobileRecommendations() {
  const recommendations = [];
  // Check LCP
  const lcp = mobileMetrics.vitals.get('LCP');
  if (lcp && lcp.rating === 'poor') {
    recommendations.push('Optimize Largest Contentful Paint: Consider image optimization, preloading critical resources, or reducing server response times');
  }
  // Check CLS
  const cls = mobileMetrics.vitals.get('CLS');
  if (cls && cls.rating === 'poor') {
    recommendations.push('Reduce Cumulative Layout Shift: Set dimensions for images and ads, avoid inserting content above existing content');
  }
  // Check network
  const network = mobileMetrics.network.get('current');
  if (network && (network.effectiveType === 'slow-2g' || network.effectiveType === '2g')) {
    recommendations.push('Optimize for slow networks: Enable compression, reduce bundle sizes, implement progressive loading');
  }
  // Check memory usage
  const memory = mobileMetrics.memory.get('current');
  if (memory) {
    const usagePercent = (memory.used / memory.limit) * 100;
    if (usagePercent > 80) {
      recommendations.push('Reduce memory usage: Optimize images, implement lazy loading, clean up unused resources');
    }
  }
  return recommendations;
}
/**
 * Report metric to analytics/monitoring service
 */
function reportMetric(metric) {
  // In development, log to console
  perfLog('[Mobile Perf]', metric);
  // In production, send to analytics service
  if (process.env.NODE_ENV === 'production') {
    // Example: Send to Google Analytics, Sentry, or custom endpoint
    try {
      // gtag('event', 'mobile_performance', metric);
      // or
      // fetch('/api/analytics/performance', {
      //   method: 'POST',
      //   headers: { 'Content-Type': 'application/json' },
      //   body: JSON.stringify(metric)
      // });
    } catch (e) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', e);
        }
        throw e; // Re-throw for caller to handle
    }
  }
}
/**
 * Get performance rating based on thresholds
 */
function getRating(value, thresholds) {
  if (value <= thresholds.good) return 'good';
  if (value <= thresholds.poor) return 'needs-improvement';
  return 'poor';
}
/**
 * Check if device is mobile
 */
function isMobileDevice() {
  if (typeof window === 'undefined' || typeof navigator === 'undefined') {
    return false;
  }
  return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent) ||
         window.innerWidth <= 768;
}
/**
 * Cleanup performance monitoring
 */
export function cleanupMobilePerformanceMonitoring() {
  observers.forEach((observer) => {
    observer.disconnect();
  });
  observers.clear();
  // Clear metrics
  Object.values(mobileMetrics).forEach(map => map.clear());
  perfLog('[Mobile Perf] Performance monitoring cleaned up');
}
/**
 * Get current mobile metrics
 */
export function getCurrentMobileMetrics() {
  if (typeof window === 'undefined') {
    return {
      vitals: {},
      navigation: {},
      network: {},
      battery: {},
      memory: {}
    };
  }
  return {
    vitals: Object.fromEntries(mobileMetrics.vitals),
    navigation: Object.fromEntries(mobileMetrics.navigation),
    network: Object.fromEntries(mobileMetrics.network),
    battery: Object.fromEntries(mobileMetrics.battery),
    memory: Object.fromEntries(mobileMetrics.memory)
  };
} 