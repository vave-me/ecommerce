/**
 * ADVANCED HOOK OPTIMIZATION UTILITIES
 * Production-ready performance patterns for React hooks
 * - Stable reference helpers
 * - Memory leak prevention
 * - Optimized state management
 * - Event listener optimization
 */
import { useRef, useCallback, useMemo, useEffect, useState } from 'react';
/**
 * Creates a stable reference that persists across renders
 * Prevents unnecessary re-renders caused by object recreation
 */
export const useStableRef = (value) => {
  const ref = useRef(value);
  ref.current = value;
  return ref;
};
/**
 * Optimized previous value hook
 * Returns the previous value of a state/prop
 */
export const usePrevious = (value) => {
  const ref = useRef();
  useEffect(() => {
    ref.current = value;
  });
  return ref.current;
};
/**
 * Advanced state setter that only updates when value actually changes
 * Prevents unnecessary re-renders from duplicate state updates
 */
export const useOptimizedState = (initialValue) => {
  const [state, setState] = useState(initialValue);
  const setOptimizedState = useCallback((newValue) => {
    setState(prevValue => {
      const actualNewValue = typeof newValue === 'function' ? newValue(prevValue) : newValue;
      // Deep comparison for objects and arrays
      if (typeof actualNewValue === 'object' && actualNewValue !== null) {
        return JSON.stringify(actualNewValue) !== JSON.stringify(prevValue) 
          ? actualNewValue 
          : prevValue;
      }
      // Simple comparison for primitives
      return actualNewValue !== prevValue ? actualNewValue : prevValue;
    });
  }, []);
  return [state, setOptimizedState];
};
/**
 * Debounced state hook with optimized performance
 * Delays state updates until after the specified delay
 */
export const useDebouncedState = (initialValue, delay = 300) => {
  const [value, setValue] = useState(initialValue);
  const [debouncedValue, setDebouncedValue] = useState(initialValue);
  const timeoutRef = useRef(null);
  const setValueDebounced = useCallback((newValue) => {
    setValue(newValue);
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }
    timeoutRef.current = setTimeout(() => {
      setDebouncedValue(newValue);
    }, delay);
  }, [delay]);
  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);
  return [value, debouncedValue, setValueDebounced];
};
/**
 * Throttled callback hook
 * Limits function execution to once per specified interval
 */
export const useThrottledCallback = (callback, delay) => {
  const lastRunRef = useRef(0);
  const timeoutRef = useRef(null);
  return useCallback((...args) => {
    const now = Date.now();
    if (now - lastRunRef.current >= delay) {
      lastRunRef.current = now;
      callback(...args);
    } else {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
      timeoutRef.current = setTimeout(() => {
        lastRunRef.current = Date.now();
        callback(...args);
      }, delay - (now - lastRunRef.current));
    }
  }, [callback, delay]);
};
/**
 * Optimized event listener hook with automatic cleanup
 * Prevents memory leaks and provides stable event handling
 */
export const useOptimizedEventListener = (
  target, 
  event, 
  handler, 
  options = {}
) => {
  const savedHandler = useRef();
  const { passive = true, capture = false, once = false } = options;
  // Update ref when handler changes
  useEffect(() => {
    savedHandler.current = handler;
  }, [handler]);
  useEffect(() => {
    if (!target?.addEventListener) return;
    const eventListener = (event) => savedHandler.current?.(event);
    const opts = { passive, capture, once };
    target.addEventListener(event, eventListener, opts);
    return () => {
      target.removeEventListener(event, eventListener, opts);
    };
  }, [target, event, passive, capture, once]);
};
/**
 * Local storage state hook with optimized performance
 * Syncs state with localStorage and handles JSON serialization
 */
export const useLocalStorageState = (key, initialValue) => {
  const [storedValue, setStoredValue] = useState(() => {
    try {
      if (typeof window === 'undefined') return initialValue;
      const item = window.localStorage.getItem(key);
      return item ? JSON.parse(item) : initialValue;
    } catch (error) {
      return initialValue;
    }
  });
  const setValue = useCallback((value) => {
    try {
      setStoredValue(value);
      if (typeof window !== 'undefined') {
        window.localStorage.setItem(key, JSON.stringify(value));
      }
    } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
  }, [key]);
  return [storedValue, setValue];
};
/**
 * Async state hook with loading and error handling
 * Optimized for API calls and async operations
 */
export const useAsyncState = (asyncFunction, dependencies = []) => {
  const [state, setState] = useState({
    data: null,
    loading: false,
    error: null
  });
  const abortControllerRef = useRef(null);
  const execute = useCallback(async (...args) => {
    // Cancel previous request
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }
    abortControllerRef.current = new AbortController();
    setState(prev => ({ ...prev, loading: true, error: null }));
    try {
      const result = await asyncFunction(...args);
      if (!abortControllerRef.current.signal.aborted) {
        setState({ data: result, loading: false, error: null });
      }
    } catch (error) {
      if (!abortControllerRef.current.signal.aborted) {
        setState(prev => ({ ...prev, loading: false, error }));
      }
    }
  }, [asyncFunction, ...dependencies]);
  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }
    };
  }, []);
  return [state, execute];
};
/**
 * Optimized component render tracking hook
 * Helps identify performance bottlenecks in development
 */
export const useRenderTracking = (componentName) => {
  const renderCount = useRef(0);
  const renderTimes = useRef([]);
  useEffect(() => {
    if (process.env.NODE_ENV === 'development') {
      renderCount.current += 1;
      const renderTime = performance.now();
      renderTimes.current.push(renderTime);
      // Keep only last 10 render times
      if (renderTimes.current.length > 10) {
        renderTimes.current.shift();
      }
      // Log excessive renders
      if (renderCount.current > 20 && renderCount.current % 5 === 0) {
      }
    }
  });
  return process.env.NODE_ENV === 'development' 
    ? {
        renderCount: renderCount.current,
        avgRenderTime: renderTimes.current.length > 1 
          ? (renderTimes.current[renderTimes.current.length - 1] - renderTimes.current[0]) / renderTimes.current.length
          : 0
      }
    : null;
};
/**
 * Window size hook with optimized performance
 * Uses passive event listeners and throttling
 */
export const useOptimizedWindowSize = () => {
  const [windowSize, setWindowSize] = useState({
    width: typeof window !== 'undefined' ? window.innerWidth : 0,
    height: typeof window !== 'undefined' ? window.innerHeight : 0,
  });
  const updateSize = useThrottledCallback(() => {
    setWindowSize({
      width: window.innerWidth,
      height: window.innerHeight,
    });
  }, 100);
  useOptimizedEventListener(typeof window !== 'undefined' ? window : null, 'resize', updateSize, {
    passive: true
  });
  return windowSize;
};
/**
 * Intersection observer hook with performance optimizations
 * Lazy loading and visibility detection
 */
export const useOptimizedIntersectionObserver = (options = {}) => {
  const [isIntersecting, setIsIntersecting] = useState(false);
  const [hasIntersected, setHasIntersected] = useState(false);
  const elementRef = useRef(null);
  const observerRef = useRef(null);
  useEffect(() => {
    const element = elementRef.current;
    if (!element) return;
    const observerOptions = {
      threshold: 0.1,
      rootMargin: '50px',
      ...options
    };
    observerRef.current = new IntersectionObserver(([entry]) => {
      const isVisible = entry.isIntersecting;
      setIsIntersecting(isVisible);
      if (isVisible && !hasIntersected) {
        setHasIntersected(true);
      }
    }, observerOptions);
    observerRef.current.observe(element);
    return () => {
      if (observerRef.current) {
        observerRef.current.disconnect();
      }
    };
  }, [hasIntersected, options]);
  return {
    elementRef,
    isIntersecting,
    hasIntersected
  };
};
/**
 * Media query hook with performance optimizations
 * Responsive design patterns
 */
export const useOptimizedMediaQuery = (query) => {
  const [matches, setMatches] = useState(false);
  useEffect(() => {
    if (typeof window === 'undefined') return;
    const mediaQuery = window.matchMedia(query);
    setMatches(mediaQuery.matches);
    const handler = (event) => setMatches(event.matches);
    // Use addListener for older browsers, addEventListener for newer ones
    if (mediaQuery.addListener) {
      mediaQuery.addListener(handler);
      return () => mediaQuery.removeListener(handler);
    } else {
      mediaQuery.addEventListener('change', handler);
      return () => mediaQuery.removeEventListener('change', handler);
    }
  }, [query]);
  return matches;
};
export default {
  useStableRef,
  usePrevious,
  useOptimizedState,
  useDebouncedState,
  useThrottledCallback,
  useOptimizedEventListener,
  useLocalStorageState,
  useAsyncState,
  useRenderTracking,
  useOptimizedWindowSize,
  useOptimizedIntersectionObserver,
  useOptimizedMediaQuery
}; 