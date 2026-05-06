"use client";
/**
 * HYDRATION-SAFE UTILITIES
 * Client-side utilities for preventing hydration mismatches
 * Server-safe fallbacks for dynamic content
 */
import { useEffect, useState } from 'react';
/**
 * Hook to safely handle dynamic content that might cause hydration mismatches
 */
export const useSafeDynamicContent = (getDynamicContent, dependencies = []) => {
  const [content, setContent] = useState(null);
  const [isReady, setIsReady] = useState(false);
  useEffect(() => {
    // Only generate dynamic content on client
    setContent(getDynamicContent());
    setIsReady(true);
  }, dependencies);
  return { content, isReady };
};
/**
 * Safe component for rendering dynamic/random content
 */
export const SafeDynamicContent = ({ 
  children, 
  fallback = null, 
  suppressHydrationWarning = true 
}) => {
  const [mounted, setMounted] = useState(false);
  useEffect(() => {
    setMounted(true);
  }, []);
  if (!mounted) {
    return <div suppressHydrationWarning={suppressHydrationWarning}>{fallback}</div>;
  }
  return <div suppressHydrationWarning={suppressHydrationWarning}>{children}</div>;
};
/**
 * Time-based content that's hydration safe
 */
export const SafeTimeDisplay = ({ 
  date, 
  format = 'relative',
  fallback = ''
}) => {
  const [timeString, setTimeString] = useState(fallback);
  useEffect(() => {
    if (!date) return;
    const updateTime = () => {
      try {
        const dateObj = new Date(date);
        const now = new Date();
        if (format === 'relative') {
          const diffMs = now - dateObj;
          const diffMins = Math.floor(diffMs / (1000 * 60));
          if (diffMins < 1) {
            setTimeString('just now');
          } else if (diffMins < 60) {
            setTimeString(`${diffMins}m ago`);
          } else {
            const diffHours = Math.floor(diffMins / 60);
            if (diffHours < 24) {
              setTimeString(`${diffHours}h ago`);
            } else {
              const diffDays = Math.floor(diffHours / 24);
              setTimeString(`${diffDays}d ago`);
            }
          }
        } else {
          setTimeString(dateObj.toISOString().split('T')[0]);
        }
      } catch (error) {
        setTimeString(fallback);
      }
    };
    updateTime();
    // Update relative times every minute
    if (format === 'relative') {
      const interval = setInterval(updateTime, 60000);
      return () => clearInterval(interval);
    }
  }, [date, format, fallback]);
  return <SafeDynamicContent fallback={fallback}>{timeString}</SafeDynamicContent>;
};
/**
 * Safe hydration wrapper component
 * Prevents hydration mismatches for dynamic content
 */
export const HydrationBoundary = ({ children, fallback = null }) => {
  const [isHydrated, setIsHydrated] = useState(false);
  useEffect(() => {
    setIsHydrated(true);
  }, []);
  if (!isHydrated) {
    return fallback;
  }
  return children;
};
/**
 * Hook for hydration-safe styles
 */
export const useHydrationSafeStyles = (styles = {}) => {
  const [isHydrated, setIsHydrated] = useState(false);
  const [safeStyles, setSafeStyles] = useState({});
  useEffect(() => {
    setIsHydrated(true);
    setSafeStyles(styles);
  }, [styles]);
  // Return empty styles during SSR to prevent mismatch
  return isHydrated ? safeStyles : {};
};
/**
 * Hook for client-only operations
 */
export const useClientOnly = (callback, dependencies = []) => {
  const [isClient, setIsClient] = useState(false);
  useEffect(() => {
    setIsClient(true);
    if (callback) {
      callback();
    }
  }, dependencies);
  return isClient;
};
/**
 * HYDRATION-SAFE UTILITIES - PHASE 8
 * Prevents SSR/CSR mismatches and improves hydration performance
 * 
 * FIX 101: Hydration Safety & Performance Optimization
 */
import { useState, useEffect, useRef, useCallback } from 'react';
/**
 * Hydration-safe state hook
 * Prevents SSR/CSR mismatches by deferring state initialization
 */
export const useHydrationSafeState = (clientValue, serverValue = null) => {
  const [isHydrated, setIsHydrated] = useState(false);
  const [value, setValue] = useState(serverValue);
  useEffect(() => {
    setIsHydrated(true);
    setValue(clientValue);
  }, [clientValue]);
  return [isHydrated ? value : serverValue, setValue, isHydrated];
};
/**
 * Hydration-safe media query hook
 * Prevents layout shifts during hydration
 */
export const useHydrationSafeMediaQuery = (query, serverDefault = false) => {
  const [matches, setMatches] = useState(serverDefault);
  const [isHydrated, setIsHydrated] = useState(false);
  useEffect(() => {
    setIsHydrated(true);
    if (typeof window !== 'undefined') {
      const mediaQuery = window.matchMedia(query);
      setMatches(mediaQuery.matches);
      const handler = (e) => setMatches(e.matches);
      mediaQuery.addEventListener('change', handler);
      return () => mediaQuery.removeEventListener('change', handler);
    }
  }, [query]);
  return [isHydrated ? matches : serverDefault, isHydrated];
};
/**
 * Hydration-safe localStorage hook
 * Prevents SSR errors and improves performance
 */
export const useHydrationSafeLocalStorage = (key, defaultValue) => {
  const [isHydrated, setIsHydrated] = useState(false);
  const [value, setValue] = useState(defaultValue);
  useEffect(() => {
    setIsHydrated(true);
    try {
      const item = localStorage.getItem(key);
      if (item !== null) {
        setValue(JSON.parse(item));
      }
    } catch (error) {
        // Initialization error - log but continue
        if (process.env.NODE_ENV === 'development') {
            console.error('Initialization error:', error);
        }
        // Continue with default behavior
    }
  }, [key]);
  const setStoredValue = useCallback((newValue) => {
    try {
      setValue(newValue);
      if (typeof window !== 'undefined') {
        localStorage.setItem(key, JSON.stringify(newValue));
      }
    } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
  }, [key]);
  return [isHydrated ? value : defaultValue, setStoredValue, isHydrated];
};
/**
 * Hydration-safe viewport detection
 * Prevents layout shifts and improves mobile experience
 */
export const useHydrationSafeViewport = () => {
  const [viewport, setViewport] = useState({
    width: 1024, // Safe server default
    height: 768,
    isMobile: false,
    isTablet: false,
    isDesktop: true
  });
  const [isHydrated, setIsHydrated] = useState(false);
  useEffect(() => {
    setIsHydrated(true);
    if (typeof window !== 'undefined') {
      const updateViewport = () => {
        const width = window.innerWidth;
        const height = window.innerHeight;
        setViewport({
          width,
          height,
          isMobile: width < 768,
          isTablet: width >= 768 && width < 1024,
          isDesktop: width >= 1024
        });
      };
      updateViewport();
      window.addEventListener('resize', updateViewport);
      return () => window.removeEventListener('resize', updateViewport);
    }
  }, []);
  return [viewport, isHydrated];
};
/**
 * Hydration-safe theme detection
 * Prevents flash of incorrect theme
 */
export const useHydrationSafeTheme = () => {
  const [theme, setTheme] = useState('light'); // Safe server default
  const [isHydrated, setIsHydrated] = useState(false);
  useEffect(() => {
    setIsHydrated(true);
    if (typeof window !== 'undefined') {
      // Check localStorage first
      try {
        const savedTheme = localStorage.getItem('theme');
        if (savedTheme) {
          setTheme(savedTheme);
          return;
        }
      } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
      // Fall back to system preference
      const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
      setTheme(mediaQuery.matches ? 'dark' : 'light');
      const handler = (e) => setTheme(e.matches ? 'dark' : 'light');
      mediaQuery.addEventListener('change', handler);
      return () => mediaQuery.removeEventListener('change', handler);
    }
  }, []);
  const updateTheme = useCallback((newTheme) => {
    setTheme(newTheme);
    try {
      localStorage.setItem('theme', newTheme);
    } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
  }, []);
  return [isHydrated ? theme : 'light', updateTheme, isHydrated];
};
/**
 * Hydration-safe intersection observer
 * Improves performance for lazy loading
 */
export const useHydrationSafeIntersection = (options = {}) => {
  const ref = useRef(null);
  const [isIntersecting, setIsIntersecting] = useState(false);
  const [isHydrated, setIsHydrated] = useState(false);
  useEffect(() => {
    setIsHydrated(true);
    if (typeof window !== 'undefined' && ref.current) {
      const observer = new IntersectionObserver(
        ([entry]) => setIsIntersecting(entry.isIntersecting),
        {
          threshold: 0.1,
          rootMargin: '50px',
          ...options
        }
      );
      observer.observe(ref.current);
      return () => observer.disconnect();
    }
  }, [options]);
  return [ref, isHydrated ? isIntersecting : false, isHydrated];
};
/**
 * Hydration-safe component wrapper
 * Prevents hydration mismatches for client-only components
 */
export const HydrationSafeComponent = ({ 
  children, 
  fallback = null,
  className = ''
}) => {
  const [isHydrated, setIsHydrated] = useState(false);
  useEffect(() => {
    setIsHydrated(true);
  }, []);
  if (!isHydrated) {
    return fallback ? (
      <div className={className}>{fallback}</div>
    ) : null;
  }
  return children;
};
/**
 * Hydration-safe user agent detection
 * Prevents SSR/CSR mismatches for device-specific features
 */
export const useHydrationSafeUserAgent = () => {
  const [userAgent, setUserAgent] = useState({
    isSafari: false,
    isChrome: false,
    isFirefox: false,
    isEdge: false,
    isMobile: false,
    isIOS: false,
    isAndroid: false
  });
  const [isHydrated, setIsHydrated] = useState(false);
  useEffect(() => {
    setIsHydrated(true);
    if (typeof window !== 'undefined' && navigator.userAgent) {
      const ua = navigator.userAgent;
      setUserAgent({
        isSafari: /Safari/.test(ua) && !/Chrome/.test(ua),
        isChrome: /Chrome/.test(ua) && !/Edge/.test(ua),
        isFirefox: /Firefox/.test(ua),
        isEdge: /Edge/.test(ua),
        isMobile: /Mobile|Android|iPhone|iPad/.test(ua),
        isIOS: /iPhone|iPad|iPod/.test(ua),
        isAndroid: /Android/.test(ua)
      });
    }
  }, []);
  return [userAgent, isHydrated];
};
/**
 * Performance-optimized scroll detection
 * Throttled and hydration-safe
 */
export const useHydrationSafeScroll = (throttleMs = 16) => {
  const [scrollY, setScrollY] = useState(0);
  const [isScrolling, setIsScrolling] = useState(false);
  const [isHydrated, setIsHydrated] = useState(false);
  const timeoutRef = useRef(null);
  useEffect(() => {
    setIsHydrated(true);
    if (typeof window !== 'undefined') {
      let ticking = false;
      const updateScrollY = () => {
        setScrollY(window.scrollY);
        setIsScrolling(true);
        // Clear existing timeout
        if (timeoutRef.current) {
          clearTimeout(timeoutRef.current);
        }
        // Set scrolling to false after throttle period
        timeoutRef.current = setTimeout(() => {
          setIsScrolling(false);
        }, throttleMs * 2);
        ticking = false;
      };
      const onScroll = () => {
        if (!ticking) {
          requestAnimationFrame(updateScrollY);
          ticking = true;
        }
      };
      window.addEventListener('scroll', onScroll, { passive: true });
      return () => {
        window.removeEventListener('scroll', onScroll);
        if (timeoutRef.current) {
          clearTimeout(timeoutRef.current);
        }
      };
    }
  }, [throttleMs]);
  return [scrollY, isScrolling, isHydrated];
};
/**
 * Hydration-safe feature detection
 * Prevents errors for unsupported browser features
 */
export const useHydrationSafeFeatures = () => {
  const [features, setFeatures] = useState({
    hasIntersectionObserver: false,
    hasResizeObserver: false,
    hasWebP: false,
    hasTouch: false,
    hasHover: false,
    supportsPassive: false
  });
  const [isHydrated, setIsHydrated] = useState(false);
  useEffect(() => {
    setIsHydrated(true);
    if (typeof window !== 'undefined') {
      // Test for passive event support
      let supportsPassive = false;
      try {
        const opts = Object.defineProperty({}, 'passive', {
          get() {
            supportsPassive = true;
            return true;
          }
        });
        window.addEventListener('testPassive', null, opts);
        window.removeEventListener('testPassive', null, opts);
      } catch (e) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', e);
        }
    }
      setFeatures({
        hasIntersectionObserver: 'IntersectionObserver' in window,
        hasResizeObserver: 'ResizeObserver' in window,
        hasWebP: (() => {
          const canvas = document.createElement('canvas');
          return canvas.toDataURL('image/webp').indexOf('webp') > -1;
        })(),
        hasTouch: 'ontouchstart' in window,
        hasHover: window.matchMedia('(hover: hover)').matches,
        supportsPassive
      });
    }
  }, []);
  return [features, isHydrated];
};
export default {
  useSafeDynamicContent,
  SafeDynamicContent,
  SafeTimeDisplay,
  HydrationBoundary,
  useHydrationSafeStyles,
  useClientOnly,
  useHydrationSafeState,
  useHydrationSafeMediaQuery,
  useHydrationSafeLocalStorage,
  useHydrationSafeViewport,
  useHydrationSafeTheme,
  useHydrationSafeIntersection,
  useHydrationSafeUserAgent,
  useHydrationSafeScroll,
  useHydrationSafeFeatures,
  HydrationSafeComponent
}; 