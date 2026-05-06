import React, { memo, useMemo, useCallback, Suspense, lazy } from 'react';
import { ErrorBoundary } from '../ErrorBoundary/ErrorBoundary';
/**
 * ADAN Performance Optimizer Utility
 * High-impact performance patterns for immediate optimization
 * 
 * Features:
 * - Auto-memoization wrapper
 * - Lazy loading helpers
 * - Stable callback factory
 * - Image optimization
 * - Bundle splitting utilities
 */
/**
 * Enhanced React.memo wrapper with deep comparison support
 * Use for components with complex props
 */
export const withDeepMemo = (Component, customCompare) => {
  const MemoizedComponent = memo(Component, customCompare || ((prevProps, nextProps) => {
    // Custom deep comparison for arrays and objects
    const keys = Object.keys(prevProps);
    if (keys.length !== Object.keys(nextProps).length) return false;
    return keys.every(key => {
      const prev = prevProps[key];
      const next = nextProps[key];
      // Handle arrays
      if (Array.isArray(prev) && Array.isArray(next)) {
        return prev.length === next.length && 
               prev.every((item, index) => item === next[index]);
      }
      // Handle objects
      if (typeof prev === 'object' && typeof next === 'object' && prev !== null && next !== null) {
        return JSON.stringify(prev) === JSON.stringify(next);
      }
      return prev === next;
    });
  }));
  MemoizedComponent.displayName = `DeepMemo(${Component.displayName || Component.name})`;
  return MemoizedComponent;
};
/**
 * Lazy component factory with error boundary and loading states
 * Automatic code splitting for better performance
 */
export const createLazyComponent = (importFn, fallback = null) => {
  const LazyComponent = lazy(importFn);
  return React.forwardRef((props, ref) => (
    <ErrorBoundary fallback={<div>Error loading component</div>}>
      <Suspense fallback={fallback || <div>Loading...</div>}>
        <LazyComponent {...props} ref={ref} />
      </Suspense>
    </ErrorBoundary>
  ));
};
/**
 * Stable callback factory to prevent unnecessary re-renders
 * Automatically memoizes callbacks with dependencies
 */
export const useStableCallback = (callback, deps = []) => {
  return useCallback(callback, deps);
};
/**
 * Optimized image component with automatic format detection
 * Implements Next.js Image with performance best practices
 */
import Image from 'next/image';
export const OptimizedImage = memo(({
  src,
  alt,
  width,
  height,
  priority = false,
  quality = 85,
  placeholder = 'blur',
  blurDataURL,
  sizes,
  fill = false,
  className,
  ...props
}) => {
  // Auto-generate blur placeholder for better UX
  const autoBlurDataURL = useMemo(() => {
    if (blurDataURL) return blurDataURL;
    // Generate a simple blur placeholder
    return `data:image/svg+xml;base64,${btoa(
      `<svg width="${width || 400}" height="${height || 300}" xmlns="http://www.w3.org/2000/svg">
        <rect width="100%" height="100%" fill="#e5e7eb"/>
        <rect width="60%" height="20%" x="20%" y="40%" rx="4" fill="#d1d5db"/>
        <rect width="40%" height="15%" x="20%" y="65%" rx="4" fill="#d1d5db"/>
      </svg>`
    )}`;
  }, [width, height, blurDataURL]);
  // Responsive sizes for different breakpoints
  const responsiveSizes = useMemo(() => {
    if (sizes) return sizes;
    return '(max-width: 768px) 100vw, (max-width: 1400px) 50vw, 33vw';
  }, [sizes]);
  return (
    <Image
      src={src}
      alt={alt}
      width={width}
      height={height}
      fill={fill}
      priority={priority}
      quality={quality}
      placeholder={placeholder}
      blurDataURL={autoBlurDataURL}
      sizes={responsiveSizes}
      className={className}
      {...props}
    />
  );
});
OptimizedImage.displayName = 'OptimizedImage';
/**
 * Performance monitoring hook
 * Tracks component render performance
 */
export const usePerformanceMonitor = (componentName) => {
  React.useEffect(() => {
    const startTime = performance.now();
    return () => {
      const endTime = performance.now();
      const renderTime = endTime - startTime;
      if (renderTime > 16) { // Slower than 60fps
      }
    };
  });
};
/**
 * Bundle optimization helper
 * Dynamic imports with preloading
 */
export const withPreload = (importFn) => {
  let componentPromise = null;
  const preload = () => {
    if (!componentPromise) {
      componentPromise = importFn();
    }
    return componentPromise;
  };
  const LazyComponent = lazy(() => preload());
  LazyComponent.preload = preload;
  return LazyComponent;
};
/**
 * Virtual list optimization for large datasets
 * Renders only visible items
 */
export const VirtualList = memo(({
  items,
  renderItem,
  itemHeight = 100,
  containerHeight = 400,
  overscan = 5
}) => {
  const [scrollTop, setScrollTop] = React.useState(0);
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
  const totalHeight = items.length * itemHeight;
  const offsetY = visibleRange.start * itemHeight;
  return (
    <div
      style={{ height: containerHeight, overflow: 'auto' }}
      onScroll={(e) => setScrollTop(e.target.scrollTop)}
    >
      <div style={{ height: totalHeight, position: 'relative' }}>
        <div style={{ transform: `translateY(${offsetY}px)` }}>
          {visibleItems.map(({ item, index }) => (
            <div key={index} style={{ height: itemHeight }}>
              {renderItem(item, index)}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
});
VirtualList.displayName = 'VirtualList';
/**
 * Intersection Observer hook for lazy loading
 * Optimizes resource loading based on visibility
 */
export const useIntersectionObserver = (options = {}) => {
  const [isVisible, setIsVisible] = React.useState(false);
  const [hasBeenVisible, setHasBeenVisible] = React.useState(false);
  const elementRef = React.useRef(null);
  React.useEffect(() => {
    const element = elementRef.current;
    if (!element) return;
    const observer = new IntersectionObserver(
      ([entry]) => {
        const visible = entry.isIntersecting;
        setIsVisible(visible);
        if (visible && !hasBeenVisible) {
          setHasBeenVisible(true);
        }
      },
      {
        threshold: 0.1,
        rootMargin: '50px',
        ...options
      }
    );
    observer.observe(element);
    return () => observer.disconnect();
  }, [hasBeenVisible, options]);
  return { elementRef, isVisible, hasBeenVisible };
};
/**
 * Debounced value hook for performance optimization
 * Prevents excessive API calls and re-renders
 */
export const useDebouncedValue = (value, delay = 300) => {
  const [debouncedValue, setDebouncedValue] = React.useState(value);
  React.useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);
    return () => clearTimeout(handler);
  }, [value, delay]);
  return debouncedValue;
};
/**
 * Component performance wrapper
 * Automatically applies performance optimizations
 */
export const withPerformanceOptimizations = (Component) => {
  const OptimizedComponent = memo(React.forwardRef((props, ref) => {
    usePerformanceMonitor(Component.displayName || Component.name);
    return <Component {...props} ref={ref} />;
  }));
  OptimizedComponent.displayName = `Optimized(${Component.displayName || Component.name})`;
  return OptimizedComponent;
};
export default {
  withDeepMemo,
  createLazyComponent,
  useStableCallback,
  OptimizedImage,
  usePerformanceMonitor,
  withPreload,
  VirtualList,
  useIntersectionObserver,
  useDebouncedValue,
  withPerformanceOptimizations
}; 