/**
 * Production-safe logging utility
 * Automatically strips logging in production builds
 */

const isDevelopment = process.env.NODE_ENV === 'development';

/**
 * Development-only console logger
 * All methods are no-ops in production
 */
export const logger = {
  log: (...args) => {
    if (isDevelopment) {
      
    }
  },
  
  warn: (...args) => {
    if (isDevelopment) {
      
    }
  },
  
  error: (...args) => {
    if (isDevelopment) {
      // Error: ...args...
    }
    // Always send critical errors to tracking service
    if (typeof window !== 'undefined' && window.Sentry) {
      window.Sentry.captureMessage(args.join(' '), 'error');
    }
  },
  
  debug: (...args) => {
    if (isDevelopment) {
      
    }
  },
  
  info: (...args) => {
    if (isDevelopment) {
      
    }
  },
  
  time: (label) => {
    if (isDevelopment) {
      
    }
  },
  
  timeEnd: (label) => {
    if (isDevelopment) {
      
    }
  },
  
  table: (data) => {
    if (isDevelopment) {
      
    }
  },
  
  group: (label) => {
    if (isDevelopment) {
      
    }
  },
  
  groupEnd: () => {
    if (isDevelopment) {
      
    }
  }
};

/**
 * Performance logging for development
 */
export const perfLogger = {
  mark: (name) => {
    if (isDevelopment && typeof performance !== 'undefined') {
      performance.mark(name);
      
    }
  },
  
  measure: (name, startMark, endMark) => {
    if (isDevelopment && typeof performance !== 'undefined') {
      try {
        performance.measure(name, startMark, endMark);
        const measure = performance.getEntriesByName(name, 'measure')[0];
        console.log(`⏱️ Performance measure ${name}: ${measure.duration}ms`);
        return measure.duration;
      } catch (error) {
        // Form submission error
        if (process.env.NODE_ENV === 'development') {
            console.error('Form submission error:', error);
        }
        // Could set error state here if available
        throw error;
    }
    }
    return 0;
  },
  
  logRender: (componentName, renderTime) => {
    if (isDevelopment) {
      const emoji = renderTime > 16 ? '🐌' : renderTime > 8 ? '⚠️' : '✅';
      console.log(`${emoji} ${componentName} rendered in ${renderTime}ms`);
    }
  }
};

/**
 * API logging utility
 */
export const apiLogger = {
  request: (method, url, data) => {
    if (isDevelopment) {
      console.log(`📡 API ${method}: ${url}`, data ? { data } : '');
    }
  },
  
  response: (method, url, status, data) => {
    if (isDevelopment) {
      const emoji = status >= 400 ? '❌' : status >= 300 ? '⚠️' : '✅';
      console.log(`📡 API ${method}: ${url} ${emoji} ${status}`, { data });
    }
  },
  
  error: (method, url, error) => {
    if (isDevelopment) {
      console.error(`💥 API ${method.toUpperCase()} ERROR: ${url}`, error);
    }
    // In production, send to error tracking
    // Example: trackApiError({ method, url, error });
  }
};

/**
 * Component lifecycle logging
 */
export const componentLogger = {
  mount: (componentName) => {
    if (isDevelopment) {
      
    }
  },
  
  unmount: (componentName) => {
    if (isDevelopment) {
      
    }
  },
  
  update: (componentName, props) => {
    if (isDevelopment) {
      
    }
  },
  
  error: (componentName, error, errorInfo) => {
    if (isDevelopment) {
      // Error: `💥 ${componentName} error:`, error, errorInfo...
    }
    // In production, send to error tracking
    // Example: trackComponentError({ componentName, error, errorInfo });
  }
};

/**
 * Security logging for suspicious activities
 */
export const securityLogger = {
  suspiciousActivity: (type, details) => {
    if (isDevelopment) {
      
    }
    // In production, always log security events to monitoring service
    // Example: sendSecurityAlert({ type, details, timestamp: new Date() });
  },
  
  invalidInput: (field, value, reason) => {
    if (isDevelopment) {
      
    }
  }
};

// Default export for common usage
export default logger; 