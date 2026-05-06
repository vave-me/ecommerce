/**
 * Performance Optimization Utilities
 * Provides memoization helpers and performance monitoring for React components
 */
import React, { useMemo, useCallback, useRef, useEffect, useState, memo } from 'react';
// ============================================================================
// MEMOIZATION HELPERS
// ============================================================================
/**
 * Creates a stable object reference that only changes when dependencies change
 * Prevents unnecessary re-renders caused by object recreation
 */
export const useStableObject = (obj, deps) => {
  return useMemo(() => obj, deps);
};
/**
 * Creates a stable array reference that only changes when dependencies change
 * Prevents unnecessary re-renders caused by array recreation
 */
export const useStableArray = (arr, deps) => {
  return useMemo(() => arr, deps);
};
/**
 * Memoizes expensive computations with dependency tracking
 */
export const useExpensiveComputation = (computeFn, deps) => {
  return useMemo(() => {
    const startTime = performance.now();
    const result = computeFn();
    const endTime = performance.now();
    if (process.env.NODE_ENV === 'development' && endTime - startTime > 16) {
    }
    return result;
  }, deps);
};
/**
 * Creates a debounced callback that prevents excessive function calls
 */
export const useDebouncedCallback = (callback, delay, deps) => {
  const timeoutRef = useRef(null);
  return useCallback((...args) => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }
    timeoutRef.current = setTimeout(() => {
      callback(...args);
    }, delay);
  }, [callback, delay, ...deps]);
};
/**
 * Creates a throttled callback that limits function execution frequency
 */
export const useThrottledCallback = (callback, delay, deps) => {
  const lastCallRef = useRef(0);
  return useCallback((...args) => {
    const now = Date.now();
    if (now - lastCallRef.current >= delay) {
      lastCallRef.current = now;
      callback(...args);
    }
  }, [callback, delay, ...deps]);
};
// ============================================================================
// COMPONENT OPTIMIZATION HELPERS
// ============================================================================
/**
 * Higher-order component that adds performance monitoring
 */
export const withPerformanceMonitoring = (WrappedComponent, componentName) => {
  const MemoizedComponent = memo(WrappedComponent);
  return React.forwardRef((props, ref) => {
    const renderStartTime = useRef(performance.now());
    useEffect(() => {
      const renderEndTime = performance.now();
      const renderTime = renderEndTime - renderStartTime.current;
      if (process.env.NODE_ENV === 'development' && renderTime > 16) {
      }
    });
    return <MemoizedComponent {...props} ref={ref} />;
  });
};
/**
 * Hook for lazy loading components with error boundaries
 */
export const useLazyComponent = (importFn) => {
  const [Component, setComponent] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  useEffect(() => {
    importFn()
      .then(module => {
        setComponent(() => module.default || module);
        setLoading(false);
      })
      .catch(err => {
        setError(err);
        setLoading(false);
      });
  }, []);
  return { Component, loading, error };
};
// ============================================================================
// INTERSECTION OBSERVER OPTIMIZATION
// ============================================================================
/**
 * Optimized intersection observer hook with performance monitoring
 */
export const useIntersectionObserver = (options = {}) => {
  const [isIntersecting, setIsIntersecting] = useState(false);
  const targetRef = useRef(null);
  const observerRef = useRef(null);
  const { threshold = 0.1, rootMargin = '0px', root = null } = options;
  useEffect(() => {
    const target = targetRef.current;
    if (!target) return;
    const observer = new IntersectionObserver(
      ([entry]) => {
        setIsIntersecting(entry.isIntersecting);
      },
      { threshold, rootMargin, root }
    );
    observer.observe(target);
    observerRef.current = observer;
    return () => {
      if (observerRef.current) {
        observerRef.current.disconnect();
      }
    };
  }, [threshold, rootMargin, root]);
  return [targetRef, isIntersecting];
};
// ============================================================================
// SCROLL OPTIMIZATION
// ============================================================================
/**
 * Optimized scroll position hook with throttling
 */
export const useOptimizedScrollPosition = (throttleMs = 16) => {
  const [scrollPosition, setScrollPosition] = useState({ x: 0, y: 0 });
  const updateScrollPosition = useThrottledCallback(() => {
    setScrollPosition({
      x: window.scrollX,
      y: window.scrollY
    });
  }, throttleMs, []);
  useEffect(() => {
    window.addEventListener('scroll', updateScrollPosition, { passive: true });
    return () => window.removeEventListener('scroll', updateScrollPosition);
  }, [updateScrollPosition]);
  return scrollPosition;
};
// ============================================================================
// VIRTUAL SCROLLING UTILITIES
// ============================================================================
/**
 * Basic virtual scrolling implementation for large lists
 */
export const useVirtualScrolling = (items, itemHeight, containerHeight) => {
  const [scrollTop, setScrollTop] = useState(0);
  const visibleCount = Math.ceil(containerHeight / itemHeight) + 2; // Buffer
  const startIndex = Math.floor(scrollTop / itemHeight);
  const endIndex = Math.min(startIndex + visibleCount, items.length);
  const visibleItems = useMemo(() => {
    return items.slice(startIndex, endIndex).map((item, index) => ({
      ...item,
      index: startIndex + index,
      top: (startIndex + index) * itemHeight
    }));
  }, [items, startIndex, endIndex, itemHeight]);
  const totalHeight = items.length * itemHeight;
  const onScroll = useCallback((e) => {
    setScrollTop(e.target.scrollTop);
  }, []);
  return {
    visibleItems,
    totalHeight,
    onScroll,
    containerStyle: {
      height: containerHeight,
      overflow: 'auto'
    }
  };
};
// ============================================================================
// BUNDLE SIZE OPTIMIZATION
// ============================================================================
/**
 * Dynamic import with loading state management
 */
export const useDynamicImport = (importPath, fallback = null) => {
  const [module, setModule] = useState(fallback);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  useEffect(() => {
    const loadModule = async () => {
      try {
        const imported = await import(importPath);
        setModule(imported.default || imported);
      } catch (err) {
        setError(err);
      } finally {
        setLoading(false);
      }
    };
    loadModule();
  }, [importPath]);
  return { module, loading, error };
};
// ============================================================================
// PERFORMANCE MONITORING
// ============================================================================
/**
 * Performance metrics collection
 */
export const usePerformanceMetrics = (componentName) => {
  const mountTime = useRef(performance.now());
  const renderCount = useRef(0);
  useEffect(() => {
    renderCount.current += 1;
  });
  useEffect(() => {
    return () => {
      const totalTime = performance.now() - mountTime.current;
      if (process.env.NODE_ENV === 'development') {
        }
    };
  }, [componentName]);
};
/**
 * Memory usage monitoring (development only)
 */
export const useMemoryMonitoring = (componentName) => {
  useEffect(() => {
    if (process.env.NODE_ENV === 'development' && performance.memory) {
      const logMemory = () => {
      };
      const interval = setInterval(logMemory, 5000);
      return () => clearInterval(interval);
    }
  }, [componentName]);
};
// ============================================================================
// EXPORTS
// ============================================================================
export default {
  useStableObject,
  useStableArray,
  useExpensiveComputation,
  useDebouncedCallback,
  useThrottledCallback,
  withPerformanceMonitoring,
  useLazyComponent,
  useIntersectionObserver,
  useOptimizedScrollPosition,
  useVirtualScrolling,
  useDynamicImport,
  usePerformanceMetrics,
  useMemoryMonitoring
}; 