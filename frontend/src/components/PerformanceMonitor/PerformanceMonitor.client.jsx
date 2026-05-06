"use client";
import { useEffect, useState, memo } from 'react';
import { initMobilePerformanceMonitoring, cleanupMobilePerformanceMonitoring } from '../../utils/mobilePerformance';
/**
 * Client component that initializes performance monitoring in development
 * 
 * OPTIMIZED: Memoized to prevent unnecessary re-initialization
 */
const PerformanceMonitor = memo(function PerformanceMonitor() {
  const [isInitialized, setIsInitialized] = useState(false);
  const [error, setError] = useState(null);
  useEffect(() => {
    // Only initialize in browser environment
    if (typeof window === 'undefined') return;
    try {
      // Initialize mobile performance monitoring
      initMobilePerformanceMonitoring();
      setIsInitialized(true);
      // Report initialization success
      if (process.env.NODE_ENV === 'development') {
        }
    } catch (err) {
      setError(err);
    }
    // Cleanup function
    return () => {
      try {
        cleanupMobilePerformanceMonitoring();
        if (process.env.NODE_ENV === 'development') {
          }
      } catch (err) {
        // Form submission error
        if (process.env.NODE_ENV === 'development') {
            console.error('Form submission error:', err);
        }
        // Could set error state here if available
        throw err;
    }
    };
  }, []);
  // Global error handler for performance monitoring
  useEffect(() => {
    const handleError = (event) => {
      // Only log performance-related errors
      if (event.error?.message?.includes('Performance') || 
          event.error?.message?.includes('Observer')) {
        setError(event.error);
      }
    };
    const handleUnhandledRejection = (event) => {
      // Only log performance-related rejections
      if (event.reason?.message?.includes('Performance') || 
          event.reason?.message?.includes('Observer')) {
        setError(event.reason);
      }
    };
    window.addEventListener('error', handleError);
    window.addEventListener('unhandledrejection', handleUnhandledRejection);
    return () => {
      window.removeEventListener('error', handleError);
      window.removeEventListener('unhandledrejection', handleUnhandledRejection);
    };
  }, []);
  // Don't render anything in production to avoid visual clutter
  if (process.env.NODE_ENV === 'production') {
    return null;
  }
  // Development indicator
  return (
    <div 
      style={{
        position: 'fixed',
        top: 0,
        right: 0,
        zIndex: 9999,
        padding: '2px 6px',
        fontSize: '10px',
        backgroundColor: isInitialized ? '#10b981' : error ? '#ef4444' : '#f59e0b',
        color: 'white',
        borderRadius: '0 0 0 4px',
        fontFamily: 'monospace',
        opacity: 0.7,
      }}
      title={
        error 
          ? `Performance monitoring error: ${error.message}`
          : isInitialized 
            ? 'Performance monitoring active'
            : 'Performance monitoring initializing...'
      }
    >
      {error ? '❌' : isInitialized ? '📊' : '⏳'}
    </div>
  );
});
export default PerformanceMonitor; 