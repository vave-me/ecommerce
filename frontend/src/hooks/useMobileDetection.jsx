"use client";
import { useState, useEffect, useMemo, useRef, useCallback } from 'react';

// Unified breakpoints - consistent with CSS
export const BREAKPOINTS = {
  mobile: 768,
  tablet: 1024,
  desktop: 1280,
  wide: 1536,
};

// Global state management for performance
let globalState = {
  width: undefined,
  height: undefined,
  listeners: new Set(),
  initialized: false,
  resizeTimeoutId: null,
};

// Initialize global resize listener only once
const initializeGlobalListener = () => {
  if (typeof window === 'undefined' || globalState.initialized) return;
  
  // Set initial size
  globalState.width = window.innerWidth;
  globalState.height = window.innerHeight;
  
  // Debounced resize handler
  const handleResize = () => {
    if (globalState.resizeTimeoutId) {
      clearTimeout(globalState.resizeTimeoutId);
    }
    
    globalState.resizeTimeoutId = setTimeout(() => {
      const newWidth = window.innerWidth;
      const newHeight = window.innerHeight;
      
      // Only update if size actually changed
      if (newWidth !== globalState.width || newHeight !== globalState.height) {
        globalState.width = newWidth;
        globalState.height = newHeight;
        globalState.listeners.forEach(callback => callback({ width: newWidth, height: newHeight }));
      }
    }, 16); // ~60fps throttling
  };
  
  window.addEventListener('resize', handleResize, { passive: true });
  globalState.initialized = true;
};

/**
 * Unified hook for mobile detection and responsive design
 * Replaces: useIsMobile, useIsMobileOnly, useMobileDetection, useMediaQuery
 * 
 * @param {Object} options - Configuration options
 * @param {number} options.mobileBreakpoint - Custom mobile breakpoint (default: 768)
 * @param {boolean} options.defaultValue - Default value for SSR (default: false)
 * @returns {Object} Mobile detection and responsive utilities
 */
export function useMobileDetection(options = {}) {
  const {
    mobileBreakpoint = BREAKPOINTS.mobile,
    defaultValue = false,
  } = options;
  
  const [dimensions, setDimensions] = useState({
    width: globalState.width,
    height: globalState.height,
  });
  
  useEffect(() => {
    if (typeof window === 'undefined') return;
    
    initializeGlobalListener();
    
    // Component-specific callback
    const updateDimensions = (newDimensions) => {
      setDimensions(newDimensions);
    };
    
    globalState.listeners.add(updateDimensions);
    
    // Set initial dimensions if available
    if (globalState.width !== undefined) {
      setDimensions({ width: globalState.width, height: globalState.height });
    }
    
    return () => {
      globalState.listeners.delete(updateDimensions);
    };
  }, []);
  
  // Memoized calculations
  const result = useMemo(() => {
    const width = dimensions.width;
    const isValidWidth = typeof width === 'number';
    
    return {
      // Primary mobile detection
      isMobile: isValidWidth ? width <= mobileBreakpoint : defaultValue,
      
      // All breakpoint flags
      isTablet: isValidWidth 
        ? width > BREAKPOINTS.mobile && width <= BREAKPOINTS.tablet
        : defaultValue,
      isDesktop: isValidWidth 
        ? width > BREAKPOINTS.tablet
        : defaultValue,
      isWideScreen: isValidWidth 
        ? width >= BREAKPOINTS.wide
        : defaultValue,
      
      // Common combinations
      isMobileOrTablet: isValidWidth 
        ? width <= BREAKPOINTS.tablet
        : defaultValue,
      isDesktopOrWide: isValidWidth 
        ? width > BREAKPOINTS.tablet
        : defaultValue,
      
      // Raw dimensions
      width: width,
      height: dimensions.height,
      
      // Utility functions
      isMinWidth: (minWidth) => isValidWidth ? width >= minWidth : defaultValue,
      isMaxWidth: (maxWidth) => isValidWidth ? width <= maxWidth : defaultValue,
      
      // Breakpoint constants for reference
      breakpoints: BREAKPOINTS,
    };
  }, [dimensions, mobileBreakpoint, defaultValue]);
  
  return result;
}

/**
 * Simple hook that only returns boolean for mobile detection
 * For components that only need to know if device is mobile
 */
export function useIsMobile(customBreakpoint) {
  const { isMobile } = useMobileDetection({ 
    mobileBreakpoint: customBreakpoint 
  });
  return isMobile;
}

/**
 * Hook for full responsive design features
 * Returns all breakpoint states and utilities
 */
export function useResponsive(options = {}) {
  return useMobileDetection(options);
}

/**
 * Custom media query hook for specific queries
 * @param {string} query - Media query string
 * @returns {boolean} Whether the query matches
 */
export function useMediaQuery(query) {
  const [matches, setMatches] = useState(false);
  
  useEffect(() => {
    if (typeof window === 'undefined') return;
    
    let mediaQuery;
    let isSupported = false;
    
    try {
      mediaQuery = window.matchMedia(query);
      isSupported = true;
    } catch (e) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', e);
        }
    }
    
    const handleChange = () => {
      if (isSupported && mediaQuery) {
        setMatches(mediaQuery.matches);
      } else {
        // Fallback for common queries
        const width = window.innerWidth;
        if (query.includes('max-width: 768px')) {
          setMatches(width <= 768);
        } else if (query.includes('min-width: 769px')) {
          setMatches(width >= 769);
        } else {
          setMatches(false);
        }
      }
    };
    
    // Set initial value
    handleChange();
    
    if (isSupported && mediaQuery) {
      // Modern API
      if (mediaQuery.addEventListener) {
        mediaQuery.addEventListener('change', handleChange);
        return () => mediaQuery.removeEventListener('change', handleChange);
      } else {
        // Legacy API
        mediaQuery.addListener(handleChange);
        return () => mediaQuery.removeListener(handleChange);
      }
    } else {
      // Fallback to resize listener
      window.addEventListener('resize', handleChange, { passive: true });
      return () => window.removeEventListener('resize', handleChange);
    }
  }, [query]);
  
  return matches;
}

// Export breakpoint constants
export { BREAKPOINTS as breakpoints };

// Default export for backward compatibility
export default useMobileDetection;