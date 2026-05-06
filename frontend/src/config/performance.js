/**
 * PRODUCTION PERFORMANCE CONFIGURATION
 * 
 * FIX 63: Centralized performance settings
 * Optimizes application performance based on environment
 */
export const PERFORMANCE_CONFIG = {
  // Console logging (production override)
  enableConsoleLogging: process.env.NODE_ENV === 'development',
  // Component monitoring
  enablePerformanceMonitoring: process.env.NODE_ENV === 'development',
  // Memory management
  maxCacheSize: process.env.NODE_ENV === 'production' ? 50 : 100,
  // Rendering optimizations
  enableVirtualization: true,
  virtualListThreshold: 20, // Items before virtualization kicks in
  // Network optimizations
  enableRequestDebouncing: true,
  debounceDelay: 300,
  // Bundle optimizations
  enableCodeSplitting: true,
  enableLazyLoading: true,
  // Image optimizations
  enableImageOptimization: true,
  imageQuality: process.env.NODE_ENV === 'production' ? 85 : 95,
  // Animation optimizations
  respectReducedMotion: true,
  enableAnimations: !window?.matchMedia?.('(prefers-reduced-motion: reduce)')?.matches,
  // Memory thresholds
  memoryWarningThreshold: 50 * 1024 * 1024, // 50MB
  memoryCriticalThreshold: 100 * 1024 * 1024, // 100MB
};
// Environment-specific overrides
if (process.env.NODE_ENV === 'production') {
  // Production-specific optimizations
  Object.assign(PERFORMANCE_CONFIG, {
    enablePerformanceMonitoring: false,
    enableConsoleLogging: false,
    maxCacheSize: 25, // Smaller cache in production
  });
}
export default PERFORMANCE_CONFIG;
