/**
 * CORE WEB VITALS OPTIMIZER - PRODUCTION READY
 * Real-time performance monitoring and optimization
 * Targets LCP < 2.5s, FID < 100ms, CLS < 0.1, TTFB < 800ms
 */
import { getCLS, getFID, getFCP, getLCP, getTTFB } from 'web-vitals';
class CoreWebVitalsOptimizer {
  constructor() {
    this.metrics = new Map();
    this.thresholds = {
      LCP: 2500, // ms
      FID: 100,  // ms
      CLS: 0.1,  // score
      FCP: 1800, // ms
      TTFB: 800  // ms
    };
    this.optimizations = new Set();
    this.init();
  }
  init() {
    if (typeof window === 'undefined') return;
    // Track all Core Web Vitals
    getCLS(this.handleMetric.bind(this));
    getFID(this.handleMetric.bind(this));
    getFCP(this.handleMetric.bind(this));
    getLCP(this.handleMetric.bind(this));
    getTTFB(this.handleMetric.bind(this));
    // Implement real-time optimizations
    this.implementLCPOptimizations();
    this.implementCLSOptimizations();
    this.implementFIDOptimizations();
  }
  handleMetric(metric) {
    this.metrics.set(metric.name, metric);
    // Auto-optimize if threshold exceeded
    if (metric.value > this.thresholds[metric.name]) {
      this.optimizeMetric(metric);
    }
    // Send to analytics in production
    if (process.env.NODE_ENV === 'production' && window.gtag) {
      window.gtag('event', metric.name, {
        event_category: 'Web Vitals',
        value: Math.round(metric.name === 'CLS' ? metric.value * 1000 : metric.value),
        non_interaction: true,
      });
    }
  }
  // LCP (Largest Contentful Paint) Optimizations
  implementLCPOptimizations() {
    // Preload critical resources
    this.preloadCriticalResources();
    // Optimize images with intersection observer
    this.optimizeImageLoading();
    // Implement resource hints
    this.addResourceHints();
  }
  preloadCriticalResources() {
    if (this.optimizations.has('preload')) return;
    const criticalResources = [
      { href: '/fonts/inter.woff2', as: 'font', type: 'font/woff2', crossorigin: 'anonymous' },
      { href: '/images/logo-small-120-high-quality.webp', as: 'image' }
    ];
    criticalResources.forEach(resource => {
      const link = document.createElement('link');
      link.rel = 'preload';
      Object.assign(link, resource);
      document.head.appendChild(link);
    });
    this.optimizations.add('preload');
  }
  optimizeImageLoading() {
    if (this.optimizations.has('images')) return;
    const imageObserver = new IntersectionObserver((entries) => {
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          const img = entry.target;
          // Implement progressive loading
          if (img.dataset.src) {
            img.src = img.dataset.src;
            img.classList.add('loaded');
            imageObserver.unobserve(img);
          }
        }
      });
    }, { rootMargin: '50px' });
    // Observe all lazy images
    document.querySelectorAll('img[data-src]').forEach(img => {
      imageObserver.observe(img);
    });
    this.optimizations.add('images');
  }
  addResourceHints() {
    if (this.optimizations.has('hints')) return;
    const hints = [
      { rel: 'dns-prefetch', href: '//fonts.googleapis.com' },
      { rel: 'dns-prefetch', href: '//www.google-analytics.com' },
      { rel: 'preconnect', href: 'https://api.openai.com' }
    ];
    hints.forEach(hint => {
      const link = document.createElement('link');
      Object.assign(link, hint);
      document.head.appendChild(link);
    });
    this.optimizations.add('hints');
  }
  // CLS (Cumulative Layout Shift) Optimizations
  implementCLSOptimizations() {
    this.fixLayoutShifts();
    this.reserveSpaceForDynamicContent();
  }
  fixLayoutShifts() {
    if (this.optimizations.has('cls')) return;
    // Add CSS for layout stability
    const style = document.createElement('style');
    style.textContent = `
      /* Prevent layout shifts */
      img, video { 
        aspect-ratio: attr(width) / attr(height);
        height: auto;
      }
      .feed-item-loading {
        min-height: 200px;
        background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
        background-size: 200% 100%;
        animation: loading 1.5s infinite;
      }
      @keyframes loading {
        0% { background-position: 200% 0; }
        100% { background-position: -200% 0; }
      }
      /* Reserve space for ads */
      .ad-container {
        min-height: 250px;
        position: relative;
      }
    `;
    document.head.appendChild(style);
    this.optimizations.add('cls');
  }
  reserveSpaceForDynamicContent() {
    // Observe dynamic content insertion
    const contentObserver = new MutationObserver((mutations) => {
      mutations.forEach(mutation => {
        if (mutation.type === 'childList') {
          mutation.addedNodes.forEach(node => {
            if (node.nodeType === 1 && node.classList.contains('dynamic-content')) {
              // Reserve space before content loads
              node.style.minHeight = '200px';
            }
          });
        }
      });
    });
    contentObserver.observe(document.body, {
      childList: true,
      subtree: true
    });
  }
  // FID (First Input Delay) Optimizations
  implementFIDOptimizations() {
    this.optimizeJavaScriptExecution();
    this.implementInputResponseOptimization();
  }
  optimizeJavaScriptExecution() {
    if (this.optimizations.has('fid')) return;
    // Break up long tasks
    this.breakUpLongTasks();
    // Use scheduler when available
    if ('scheduler' in window && 'postTask' in window.scheduler) {
      this.useSchedulerAPI();
    }
    this.optimizations.add('fid');
  }
  breakUpLongTasks() {
    // Monkey patch setTimeout for task yielding
    const originalSetTimeout = window.setTimeout;
    window.yieldToMain = () => {
      return new Promise(resolve => {
        originalSetTimeout(resolve, 0);
      });
    };
  }
  useSchedulerAPI() {
    // Prioritize user interactions
    window.scheduleUserInteraction = (callback) => {
      return window.scheduler.postTask(callback, { priority: 'user-blocking' });
    };
    window.scheduleBackgroundTask = (callback) => {
      return window.scheduler.postTask(callback, { priority: 'background' });
    };
  }
  implementInputResponseOptimization() {
    // Debounce frequent inputs
    const debouncedHandlers = new Map();
    document.addEventListener('input', (e) => {
      if (e.target.matches('input, textarea')) {
        const handler = debouncedHandlers.get(e.target);
        if (handler) clearTimeout(handler);
        debouncedHandlers.set(e.target, setTimeout(() => {
          // Process input with lower priority
          if (window.scheduleBackgroundTask) {
            window.scheduleBackgroundTask(() => {
              e.target.dispatchEvent(new Event('debouncedInput'));
            });
          }
        }, 150));
      }
    }, { passive: true });
  }
  // Performance budget monitoring
  checkPerformanceBudget() {
    const budget = {
      javascript: 200 * 1024, // 200KB
      images: 500 * 1024,     // 500KB
      total: 1000 * 1024      // 1MB
    };
    if ('performance' in window && 'getEntriesByType' in performance) {
      const resources = performance.getEntriesByType('resource');
      const usage = resources.reduce((acc, resource) => {
        const size = resource.transferSize || 0;
        if (resource.name.includes('.js')) {
          acc.javascript += size;
        } else if (resource.name.includes('.jpg') || resource.name.includes('.png') || resource.name.includes('.webp')) {
          acc.images += size;
        }
        acc.total += size;
        return acc;
      }, { javascript: 0, images: 0, total: 0 });
      // Warn if budget exceeded
      Object.entries(budget).forEach(([type, limit]) => {
        if (usage[type] > limit) {
        }
      });
      return { usage, budget, withinBudget: usage.total <= budget.total };
    }
  }
  // Real-time metric optimization
  optimizeMetric(metric) {
    switch (metric.name) {
      case 'LCP':
        this.optimizeLCP();
        break;
      case 'FID':
        this.optimizeFID();
        break;
      case 'CLS':
        this.optimizeCLS();
        break;
    }
  }
  optimizeLCP() {
    // Emergency LCP optimization
    const images = document.querySelectorAll('img');
    images.forEach(img => {
      if (!img.loading) {
        img.loading = 'lazy';
      }
      if (!img.decoding) {
        img.decoding = 'async';
      }
    });
  }
  optimizeFID() {
    // Defer non-critical JavaScript
    const scripts = document.querySelectorAll('script:not([async]):not([defer])');
    scripts.forEach(script => {
      if (!script.src.includes('critical')) {
        script.defer = true;
      }
    });
  }
  optimizeCLS() {
    // Add aspect ratios to images without them
    const images = document.querySelectorAll('img:not([width]):not([height])');
    images.forEach(img => {
      img.style.aspectRatio = '16/9'; // Default fallback
    });
  }
  getMetrics() {
    const metricsObj = {};
    this.metrics.forEach((value, key) => {
      metricsObj[key] = {
        value: value.value,
        rating: this.getRating(key, value.value),
        threshold: this.thresholds[key]
      };
    });
    return metricsObj;
  }
  getRating(metricName, value) {
    const threshold = this.thresholds[metricName];
    if (value <= threshold) return 'good';
    if (value <= threshold * 1.5) return 'needs-improvement';
    return 'poor';
  }
}
// Create singleton instance
const webVitalsOptimizer = new CoreWebVitalsOptimizer();
export { webVitalsOptimizer, CoreWebVitalsOptimizer };
export default webVitalsOptimizer; 