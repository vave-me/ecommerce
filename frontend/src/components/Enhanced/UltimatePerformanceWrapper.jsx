/**
 * ULTIMATE PERFORMANCE WRAPPER
 * Production-ready component optimization with real-time monitoring
 * Combines memoization, virtualization, lazy loading, and resource optimization
 */
import React, { 
  memo, 
  useMemo, 
  useCallback, 
  useRef, 
  useEffect, 
  useState,
  Suspense,
  lazy,
  startTransition
} from 'react';
import { useIntersectionObserver } from '../optimized/PerformanceOptimizer';
import webVitalsOptimizer from '../../utils/coreWebVitalsOptimizer';
// Performance monitoring hook
const usePerformanceMonitoring = (componentName) => {
  const renderStart = useRef(performance.now());
  const renderCount = useRef(0);
  useEffect(() => {
    renderCount.current++;
    const renderTime = performance.now() - renderStart.current;
    if (renderTime > 16) { // Slower than 60fps
    }
    renderStart.current = performance.now();
  });
  return { renderCount: renderCount.current };
};
// Ultimate memoization wrapper with deep comparison
const withUltimateMemo = (Component, compareProps) => {
  return memo(Component, compareProps || ((prevProps, nextProps) => {
    // Smart comparison logic
    const prevKeys = Object.keys(prevProps);
    const nextKeys = Object.keys(nextProps);
    if (prevKeys.length !== nextKeys.length) return false;
    return prevKeys.every(key => {
      const prev = prevProps[key];
      const next = nextProps[key];
      // Handle functions (should be memoized by parent)
      if (typeof prev === 'function' && typeof next === 'function') {
        return prev === next;
      }
      // Handle arrays with shallow comparison
      if (Array.isArray(prev) && Array.isArray(next)) {
        return prev.length === next.length && 
               prev.every((item, idx) => item === next[idx]);
      }
      // Handle objects with JSON comparison (expensive but thorough)
      if (typeof prev === 'object' && typeof next === 'object' && prev !== null && next !== null) {
        return JSON.stringify(prev) === JSON.stringify(next);
      }
      return prev === next;
    });
  }));
};
// Lazy loading with priority levels
const createLazyComponent = (importFn, priority = 'normal') => {
  const LazyComponent = lazy(() => {
    // Prioritize loading based on importance
    if (priority === 'high') {
      return importFn();
    } else if (priority === 'low') {
      return new Promise(resolve => {
        // Delay low priority components
        setTimeout(() => resolve(importFn()), 100);
      });
    }
    return importFn();
  });
  return LazyComponent;
};
// Virtual scrolling for large lists
const VirtualizedList = memo(({ 
  items, 
  renderItem, 
  itemHeight = 100, 
  containerHeight = 400,
  overscan = 5 
}) => {
  const [scrollTop, setScrollTop] = useState(0);
  const containerRef = useRef(null);
  const visibleRange = useMemo(() => {
    const start = Math.max(0, Math.floor(scrollTop / itemHeight) - overscan);
    const end = Math.min(
      items.length,
      Math.ceil((scrollTop + containerHeight) / itemHeight) + overscan
    );
    return { start, end };
  }, [scrollTop, itemHeight, containerHeight, items.length, overscan]);
  const visibleItems = useMemo(() => {
    return items.slice(visibleRange.start, visibleRange.end).map((item, index) => ({
      item,
      index: visibleRange.start + index
    }));
  }, [items, visibleRange]);
  const handleScroll = useCallback((e) => {
    setScrollTop(e.target.scrollTop);
  }, []);
  const totalHeight = items.length * itemHeight;
  const offsetY = visibleRange.start * itemHeight;
  return (
    <div
      ref={containerRef}
      style={{ 
        height: containerHeight, 
        overflow: 'auto',
        willChange: 'scroll-position'
      }}
      onScroll={handleScroll}
    >
      <div style={{ height: totalHeight, position: 'relative' }}>
        <div 
          style={{ 
            transform: `translateY(${offsetY}px)`,
            willChange: 'transform'
          }}
        >
          {visibleItems.map(({ item, index }) => (
            <div 
              key={index} 
              style={{ 
                height: itemHeight,
                contain: 'layout style paint'
              }}
            >
              {renderItem(item, index)}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
});
VirtualizedList.displayName = 'VirtualizedList';
// Optimized image component with progressive loading
const OptimizedImage = memo(({ 
  src, 
  alt, 
  width, 
  height, 
  className,
  priority = false,
  quality = 85,
  ...props 
}) => {
  const [isLoaded, setIsLoaded] = useState(false);
  const [isInView, setIsInView] = useState(priority);
  const imgRef = useRef(null);
  const { elementRef } = useIntersectionObserver({
    threshold: 0.1,
    rootMargin: '50px'
  });
  useEffect(() => {
    if (elementRef.current) {
      const observer = new IntersectionObserver(([entry]) => {
        if (entry.isIntersecting) {
          setIsInView(true);
          observer.disconnect();
        }
      }, { threshold: 0.1, rootMargin: '50px' });
      observer.observe(elementRef.current);
      return () => observer.disconnect();
    }
  }, [elementRef]);
  const handleLoad = useCallback(() => {
    startTransition(() => {
      setIsLoaded(true);
    });
  }, []);
  return (
    <div 
      ref={elementRef}
      className={`image-container ${className || ''}`}
      style={{ 
        aspectRatio: width && height ? `${width}/${height}` : undefined,
        backgroundColor: '#f0f0f0'
      }}
    >
      {isInView && (
        <img
          ref={imgRef}
          src={src}
          alt={alt}
          width={width}
          height={height}
          loading={priority ? 'eager' : 'lazy'}
          decoding="async"
          onLoad={handleLoad}
          style={{
            width: '100%',
            height: 'auto',
            opacity: isLoaded ? 1 : 0,
            transition: 'opacity 0.3s ease',
            contain: 'layout style paint'
          }}
          {...props}
        />
      )}
      {!isLoaded && (
        <div 
          className="image-placeholder"
          style={{
            width: '100%',
            height: '100%',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            backgroundColor: '#f0f0f0',
            color: '#666'
          }}
        >
          Loading...
        </div>
      )}
    </div>
  );
});
OptimizedImage.displayName = 'OptimizedImage';
// Advanced performance wrapper component
const UltimatePerformanceWrapper = ({
  children,
  componentName = 'WrappedComponent',
  enableMonitoring = process.env.NODE_ENV === 'development',
  enableVirtualization = false,
  virtualizationProps = {},
  lazy = false,
  lazyProps = {},
  memoization = true,
  customMemoComparison = null
}) => {
  const { renderCount } = usePerformanceMonitoring(componentName);
  const [hasError, setHasError] = useState(false);
  // Error boundary logic
  useEffect(() => {
    const handleError = (error, errorInfo) => {
      setHasError(true);
      // Report to monitoring service
      if (window.gtag) {
        window.gtag('event', 'exception', {
          description: `${componentName}: ${error.message}`,
          fatal: false
        });
      }
    };
    window.addEventListener('error', handleError);
    return () => window.removeEventListener('error', handleError);
  }, [componentName]);
  // Performance budget monitoring
  useEffect(() => {
    if (enableMonitoring && renderCount > 50) {
    }
  }, [renderCount, enableMonitoring, componentName]);
  // Memoized component wrapper
  const MemoizedComponent = useMemo(() => {
    if (!memoization) return children;
    if (React.isValidElement(children)) {
      const OriginalComponent = children.type;
      return withUltimateMemo(OriginalComponent, customMemoComparison);
    }
    return children;
  }, [children, memoization, customMemoComparison]);
  // Error fallback
  if (hasError) {
    return (
      <div className="error-boundary" style={{ padding: '20px', color: 'red' }}>
        <h3>Something went wrong with {componentName}</h3>
        <button onClick={() => setHasError(false)}>
          Try Again
        </button>
      </div>
    );
  }
  // Virtualization wrapper
  if (enableVirtualization) {
    return <VirtualizedList {...virtualizationProps} />;
  }
  // Lazy loading wrapper
  if (lazy) {
    const LazyComponent = createLazyComponent(
      () => Promise.resolve({ default: () => children }),
      lazyProps.priority
    );
    return (
      <Suspense fallback={<div>Loading {componentName}...</div>}>
        <LazyComponent />
      </Suspense>
    );
  }
  // Development monitoring overlay
  const DevMonitoringOverlay = enableMonitoring ? () => (
    <div 
      style={{
        position: 'fixed',
        top: 10,
        right: 10,
        background: 'rgba(0,0,0,0.8)',
        color: 'white',
        padding: '5px 10px',
        fontSize: '12px',
        zIndex: 9999,
        borderRadius: '4px'
      }}
    >
      {componentName}: {renderCount} renders
    </div>
  ) : null;
  return (
    <>
      {memoization ? MemoizedComponent : children}
      {DevMonitoringOverlay && <DevMonitoringOverlay />}
    </>
  );
};
// HOC version for easy wrapping
const withUltimatePerformance = (options = {}) => (Component) => {
  const WrappedComponent = (props) => (
    <UltimatePerformanceWrapper 
      componentName={Component.displayName || Component.name}
      {...options}
    >
      <Component {...props} />
    </UltimatePerformanceWrapper>
  );
  WrappedComponent.displayName = `UltimatePerformance(${Component.displayName || Component.name})`;
  return WrappedComponent;
};
// Resource preloading utility
const useResourcePreloader = (resources = []) => {
  useEffect(() => {
    resources.forEach(resource => {
      const link = document.createElement('link');
      link.rel = resource.rel || 'prefetch';
      link.href = resource.href;
      if (resource.as) link.as = resource.as;
      if (resource.type) link.type = resource.type;
      document.head.appendChild(link);
    });
  }, [resources]);
};
// Advanced state optimization hook
const useOptimizedState = (initialState, options = {}) => {
  const { 
    debounceMs = 0, 
    throttleMs = 0,
    enableHistory = false,
    maxHistorySize = 10 
  } = options;
  const [state, setState] = useState(initialState);
  const [history, setHistory] = useState(enableHistory ? [initialState] : []);
  const timeoutRef = useRef(null);
  const lastUpdateRef = useRef(Date.now());
  const optimizedSetState = useCallback((newState) => {
    const now = Date.now();
    // Throttling
    if (throttleMs > 0 && now - lastUpdateRef.current < throttleMs) {
      return;
    }
    // Debouncing
    if (debounceMs > 0) {
      clearTimeout(timeoutRef.current);
      timeoutRef.current = setTimeout(() => {
        setState(newState);
        lastUpdateRef.current = Date.now();
        if (enableHistory) {
          setHistory(prev => {
            const newHistory = [...prev, newState];
            return newHistory.length > maxHistorySize 
              ? newHistory.slice(-maxHistorySize)
              : newHistory;
          });
        }
      }, debounceMs);
      return;
    }
    // Immediate update
    setState(newState);
    lastUpdateRef.current = now;
    if (enableHistory) {
      setHistory(prev => {
        const newHistory = [...prev, newState];
        return newHistory.length > maxHistorySize 
          ? newHistory.slice(-maxHistorySize)
          : newHistory;
      });
    }
  }, [debounceMs, throttleMs, enableHistory, maxHistorySize]);
  return [state, optimizedSetState, history];
};
// Performance context for global optimization
const PerformanceContext = React.createContext({
  metrics: {},
  optimizations: new Set(),
  enableDebug: false
});
const PerformanceProvider = ({ children, enableDebug = false }) => {
  const [metrics, setMetrics] = useState({});
  const [optimizations] = useState(new Set());
  useEffect(() => {
    if (enableDebug) {
      // Monitor web vitals
      const updateMetrics = (metric) => {
        setMetrics(prev => ({
          ...prev,
          [metric.name]: metric.value
        }));
      };
      webVitalsOptimizer.getMetrics && updateMetrics(webVitalsOptimizer.getMetrics());
    }
  }, [enableDebug]);
  const value = useMemo(() => ({
    metrics,
    optimizations,
    enableDebug
  }), [metrics, optimizations, enableDebug]);
  return (
    <PerformanceContext.Provider value={value}>
      {children}
    </PerformanceContext.Provider>
  );
};
export {
  UltimatePerformanceWrapper,
  withUltimatePerformance,
  VirtualizedList,
  OptimizedImage,
  useResourcePreloader,
  useOptimizedState,
  PerformanceProvider,
  PerformanceContext,
  withUltimateMemo,
  createLazyComponent
};
export default UltimatePerformanceWrapper; 