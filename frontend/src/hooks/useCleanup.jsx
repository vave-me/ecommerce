/**
 * PRODUCTION MEMORY LEAK PREVENTION
 * Comprehensive cleanup hook for React components
 */
import { useEffect, useRef, useCallback } from 'react';
/**
 * Central cleanup hook to prevent memory leaks
 * Automatically handles common cleanup scenarios
 */
export const useCleanup = () => {
  const cleanupFunctions = useRef([]);
  const isMountedRef = useRef(true);
  const timersRef = useRef(new Set());
  const listenersRef = useRef(new Map());
  const observersRef = useRef([]);
  const objectUrlsRef = useRef(new Set());
  const abortControllersRef = useRef(new Set());
  // Add cleanup function
  const addCleanup = useCallback((cleanupFn) => {
    if (typeof cleanupFn === 'function') {
      cleanupFunctions.current.push(cleanupFn);
    }
  }, []);
  // Timer management
  const addTimer = useCallback((timerId) => {
    timersRef.current.add(timerId);
    return timerId;
  }, []);
  const clearTimer = useCallback((timerId) => {
    clearTimeout(timerId);
    clearInterval(timerId);
    timersRef.current.delete(timerId);
  }, []);
  const safeSetTimeout = useCallback((callback, delay) => {
    const timerId = setTimeout(() => {
      if (isMountedRef.current) {
        callback();
      }
      timersRef.current.delete(timerId);
    }, delay);
    return addTimer(timerId);
  }, [addTimer]);
  const safeSetInterval = useCallback((callback, delay) => {
    const timerId = setInterval(() => {
      if (isMountedRef.current) {
        callback();
      } else {
        clearInterval(timerId);
        timersRef.current.delete(timerId);
      }
    }, delay);
    return addTimer(timerId);
  }, [addTimer]);
  // Event listener management with enhanced error handling
  const addEventListener = useCallback((target, event, handler, options) => {
    if (!target || typeof handler !== 'function') {
      return () => {};
    }
    try {
      target.addEventListener(event, handler, options);
      const key = `${target.constructor.name}-${event}`;
      if (!listenersRef.current.has(key)) {
        listenersRef.current.set(key, []);
      }
      listenersRef.current.get(key).push({ target, event, handler, options });
      // Return cleanup function
      return () => {
        try {
          target.removeEventListener(event, handler, options);
          const listeners = listenersRef.current.get(key);
          if (listeners) {
            const index = listeners.findIndex(l => 
              l.target === target && l.event === event && l.handler === handler
            );
            if (index > -1) {
              listeners.splice(index, 1);
            }
          }
        } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
      };
    } catch (error) {
      return () => {};
    }
  }, []);
  const removeEventListener = useCallback((target, event, handler) => {
    target.removeEventListener(event, handler);
    const key = `${target.constructor.name}-${event}`;
    const listeners = listenersRef.current.get(key);
    if (listeners) {
      const index = listeners.findIndex(l => 
        l.target === target && l.event === event && l.handler === handler
      );
      if (index > -1) {
        listeners.splice(index, 1);
      }
    }
  }, []);
  // Observer management
  const addObserver = useCallback((observer) => {
    observersRef.current.push(observer);
    return observer;
  }, []);
  // Object URL management
  const createObjectURL = useCallback((file) => {
    const url = URL.createObjectURL(file);
    objectUrlsRef.current.add(url);
    return url;
  }, []);
  const revokeObjectURL = useCallback((url) => {
    URL.revokeObjectURL(url);
    objectUrlsRef.current.delete(url);
  }, []);
  // AbortController management
  const createAbortController = useCallback(() => {
    const controller = new AbortController();
    abortControllersRef.current.add(controller);
    return controller;
  }, []);
  const abortController = useCallback((controller) => {
    controller.abort();
    abortControllersRef.current.delete(controller);
  }, []);
  // Safe async operations
  const safeAsync = useCallback(async (asyncFn) => {
    try {
      const result = await asyncFn();
      return isMountedRef.current ? result : null;
    } catch (error) {
      if (isMountedRef.current) {
        throw error;
      }
      return null;
    }
  }, []);
  // Component mount status
  const isMounted = useCallback(() => isMountedRef.current, []);
  // Manual cleanup trigger
  const cleanup = useCallback(() => {
    // Execute custom cleanup functions
    cleanupFunctions.current.forEach(cleanupFn => {
      try {
        cleanupFn();
      } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
    });
    cleanupFunctions.current = [];
    // Clear all timers
    timersRef.current.forEach(timerId => {
      clearTimeout(timerId);
      clearInterval(timerId);
    });
    timersRef.current.clear();
    // Remove all event listeners
    listenersRef.current.forEach(listeners => {
      listeners.forEach(({ target, event, handler }) => {
        try {
          target.removeEventListener(event, handler);
        } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
      });
    });
    listenersRef.current.clear();
    // Disconnect all observers
    observersRef.current.forEach(observer => {
      try {
        if (observer.disconnect) observer.disconnect();
        if (observer.unobserve) observer.unobserve();
      } catch (error) {
        // Initialization error - log but continue
        if (process.env.NODE_ENV === 'development') {
            console.error('Initialization error:', error);
        }
        // Continue with default behavior
    }
    });
    observersRef.current = [];
    // Revoke all object URLs
    objectUrlsRef.current.forEach(url => {
      try {
        URL.revokeObjectURL(url);
      } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
    });
    objectUrlsRef.current.clear();
    // Abort all controllers
    abortControllersRef.current.forEach(controller => {
      try {
        controller.abort();
      } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
    });
    abortControllersRef.current.clear();
  }, []);
  // Cleanup on unmount
  useEffect(() => {
    return () => {
      isMountedRef.current = false;
      cleanup();
    };
  }, [cleanup]);
  return {
    // Cleanup management
    addCleanup,
    cleanup,
    isMounted,
    // Timer management
    safeSetTimeout,
    safeSetInterval,
    addTimer,
    clearTimer,
    // Event listener management
    addEventListener,
    removeEventListener,
    // Observer management
    addObserver,
    // Object URL management
    createObjectURL,
    revokeObjectURL,
    // AbortController management
    createAbortController,
    abortController,
    // Async operations
    safeAsync
  };
};
/**
 * Hook for managing component resources with automatic cleanup
 */
export const useResourceManager = () => {
  const cleanup = useCleanup();
  const resourcesRef = useRef({
    streams: new Set(),
    connections: new Set(),
    subscriptions: new Set(),
    workers: new Set()
  });
  // Media stream management
  const addStream = useCallback((stream) => {
    resourcesRef.current.streams.add(stream);
    cleanup.addCleanup(() => {
      if (stream && stream.getTracks) {
        stream.getTracks().forEach(track => track.stop());
      }
    });
    return stream;
  }, [cleanup]);
  // WebSocket connection management
  const addConnection = useCallback((connection) => {
    resourcesRef.current.connections.add(connection);
    cleanup.addCleanup(() => {
      if (connection && connection.close) {
        connection.close();
      }
    });
    return connection;
  }, [cleanup]);
  // Subscription management (e.g., Redux, RxJS)
  const addSubscription = useCallback((subscription) => {
    resourcesRef.current.subscriptions.add(subscription);
    cleanup.addCleanup(() => {
      if (subscription && subscription.unsubscribe) {
        subscription.unsubscribe();
      } else if (typeof subscription === 'function') {
        subscription();
      }
    });
    return subscription;
  }, [cleanup]);
  // Web Worker management
  const addWorker = useCallback((worker) => {
    resourcesRef.current.workers.add(worker);
    cleanup.addCleanup(() => {
      if (worker && worker.terminate) {
        worker.terminate();
      }
    });
    return worker;
  }, [cleanup]);
  return {
    ...cleanup,
    addStream,
    addConnection,
    addSubscription,
    addWorker
  };
};
/**
 * Hook for performance monitoring and memory usage tracking
 */
export const usePerformanceMonitor = (componentName) => {
  const cleanup = useCleanup();
  const startTime = useRef(performance.now());
  const renderCount = useRef(0);
  const memoryUsage = useRef([]);
  useEffect(() => {
    renderCount.current += 1;
    // Memory monitoring (if available)
    if (performance.memory) {
      memoryUsage.current.push({
        used: performance.memory.usedJSHeapSize,
        total: performance.memory.totalJSHeapSize,
        timestamp: Date.now()
      });
      // Keep only last 10 measurements
      if (memoryUsage.current.length > 10) {
        memoryUsage.current.shift();
      }
    }
  });
  useEffect(() => {
    cleanup.addCleanup(() => {
      const endTime = performance.now();
      const totalTime = endTime - startTime.current;
      if (process.env.NODE_ENV === 'development') {
        if (memoryUsage.current.length > 1) {
          const memStart = memoryUsage.current[0];
          const memEnd = memoryUsage.current[memoryUsage.current.length - 1];
          const memDiff = memEnd.used - memStart.used;
        }
        }
    });
  }, [cleanup, componentName]);
  return {
    renderCount: renderCount.current,
    uptime: performance.now() - startTime.current,
    memoryUsage: memoryUsage.current
  };
}; 