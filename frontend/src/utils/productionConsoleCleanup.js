/**
 * PRODUCTION CONSOLE CLEANUP
 * FIX 69: Eliminates ALL console statements that are causing performance issues
 * 
 * Removes debug logs, warnings, and unnecessary console output in production
 */

/**
 * Initialize production console cleanup
 * Removes all debug/development console statements
 */
export const initProductionConsoleCleanup = () => {
  if (process.env.NODE_ENV === 'production') {
    // Store original methods for critical errors only
    const originalError = console.error;
    const originalWarn = console.warn;
    
    // Override all console methods with no-ops
    const noop = () => {};
    
    console.log = noop;
    console.debug = noop;
    console.info = noop;
    console.trace = noop;
    console.table = noop;
    console.time = noop;
    console.timeEnd = noop;
    console.group = noop;
    console.groupEnd = noop;
    console.groupCollapsed = noop;
    console.count = noop;
    console.countReset = noop;
    console.profile = noop;
    console.profileEnd = noop;
    console.clear = noop;
    console.dir = noop;
    console.dirxml = noop;
    console.timeLog = noop;
    
    // Filter error messages to only show critical ones
    console.error = (...args) => {
      const message = args[0];
      if (typeof message === 'string') {
        // Only log truly critical errors
        const criticalPatterns = [
          'Authentication failed',
          'Network request failed',
          'Fatal error',
          'Security violation',
          'Payment error',
          'Database error'
        ];
        
        const isCritical = criticalPatterns.some(pattern => 
          message.toLowerCase().includes(pattern.toLowerCase())
        );
        
        if (isCritical) {
          originalError(...args);
        }
      }
    };
    
    // Filter warnings to only show security/performance related
    console.warn = (...args) => {
      const message = args[0];
      if (typeof message === 'string') {
        const importantPatterns = [
          'Security',
          'Performance critical',
          'Memory leak',
          'Deprecated API will be removed'
        ];
        
        const isImportant = importantPatterns.some(pattern => 
          message.includes(pattern)
        );
        
        if (isImportant) {
          originalWarn(...args);
        }
      }
    };
  }
};

/**
 * Development-only logger that's automatically stripped in production
 */
export const devLog = (...args) => {
  if (process.env.NODE_ENV === 'development') {
    }
};

/**
 * Performance logger for development debugging
 */
export const perfLog = (operation, startTime) => {
  if (process.env.NODE_ENV === 'development') {
    const duration = performance.now() - startTime;
    if (duration > 16) { // Slower than 60fps
      // Performance warning: operation took too long
    }
  }
};

/**
 * Safe error logger that reports to analytics in production
 */
export const safeErrorLog = (error, context = '') => {
  if (process.env.NODE_ENV === 'development') {
    // Error: '[ERROR]', context, error...
  } else {
    // In production, send to analytics instead of console
    if (typeof window !== 'undefined' && window.gtag) {
      try {
        window.gtag('event', 'exception', {
          description: `${context}: ${error.message || error}`,
          fatal: false
        });
      } catch (e) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', e);
        }
    }
    }
  }
};

/**
 * API logger that's production-safe
 */
export const apiLog = {
  request: (method, url, data) => {
    if (process.env.NODE_ENV === 'development') {
      // API request logged for debugging
    }
  },
  
  response: (method, url, status, data) => {
    if (process.env.NODE_ENV === 'development') {
      const emoji = status >= 400 ? '❌' : status >= 300 ? '⚠️' : '✅';
      // API response logged for debugging
    }
  },
  
  error: (method, url, error) => {
    if (process.env.NODE_ENV === 'development') {
      // Error: `💥 API ${method.toUpperCase(...} ERROR: ${url}`, error);
    } else {
      // Report API errors to analytics in production
      safeErrorLog(error, `API ${method} ${url}`);
    }
  }
};

/**
 * Component lifecycle logger for development
 */
export const componentLog = {
  mount: (componentName) => {
    if (process.env.NODE_ENV === 'development') {
      }
  },
  
  unmount: (componentName) => {
    if (process.env.NODE_ENV === 'development') {
      }
  },
  
  render: (componentName, props) => {
    if (process.env.NODE_ENV === 'development') {
      }
  },
  
  error: (componentName, error, errorInfo) => {
    safeErrorLog(error, `Component ${componentName}`);
  }
};

/**
 * Feed/state logger that won't spam production
 */
export const feedLog = {
  update: (message, data) => {
    if (process.env.NODE_ENV === 'development') {
      }
  },
  
  items: (message, items) => {
    if (process.env.NODE_ENV === 'development') {
      }
  },
  
  scroll: (message, conditions) => {
    if (process.env.NODE_ENV === 'development') {
      }
  }
};

// Initialize cleanup immediately
initProductionConsoleCleanup();

export default {
  initProductionConsoleCleanup,
  devLog,
  perfLog,
  safeErrorLog,
  apiLog,
  componentLog,
  feedLog
}; 