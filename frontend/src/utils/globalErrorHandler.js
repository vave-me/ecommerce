/**
 * Global Error Handler
 * Centralized error handling for production readiness
 * 
 * Features:
 * - Error categorization and standardization
 * - User-friendly error messages
 * - Error tracking integration (Sentry)
 * - Development vs production handling
 * - Rate limiting for error reporting
 */

import { logger } from './logger';
import { captureError, captureMessage, addBreadcrumb } from './errorTracking';

// Error types
export const ErrorTypes = {
  NETWORK: 'NETWORK_ERROR',
  AUTHENTICATION: 'AUTH_ERROR',
  VALIDATION: 'VALIDATION_ERROR',
  PERMISSION: 'PERMISSION_ERROR',
  NOT_FOUND: 'NOT_FOUND_ERROR',
  SERVER: 'SERVER_ERROR',
  CLIENT: 'CLIENT_ERROR',
  UNKNOWN: 'UNKNOWN_ERROR'
};

// Error severity levels
export const ErrorSeverity = {
  LOW: 'low',
  MEDIUM: 'medium',
  HIGH: 'high',
  CRITICAL: 'critical'
};

// Rate limiting for error reporting
const errorReportingCache = new Map();
const ERROR_CACHE_DURATION = 60000; // 1 minute
const MAX_ERRORS_PER_TYPE = 10;

class GlobalErrorHandler {
  constructor() {
    this.errorHandlers = new Map();
    this.defaultHandler = this.handleGenericError;
    this.setupGlobalHandlers();
  }

  /**
   * Setup global error handlers
   */
  setupGlobalHandlers() {
    if (typeof window !== 'undefined') {
      // Handle unhandled promise rejections
      window.addEventListener('unhandledrejection', (event) => {
        this.handleError(event.reason, {
          type: ErrorTypes.UNKNOWN,
          severity: ErrorSeverity.HIGH,
          context: 'Unhandled Promise Rejection'
        });
        event.preventDefault();
      });

      // Handle global errors
      window.addEventListener('error', (event) => {
        this.handleError(event.error || event, {
          type: ErrorTypes.UNKNOWN,
          severity: ErrorSeverity.HIGH,
          context: 'Global Error Handler'
        });
        event.preventDefault();
      });
    }
  }

  /**
   * Register custom error handler for specific error types
   */
  registerHandler(errorType, handler) {
    this.errorHandlers.set(errorType, handler);
  }

  /**
   * Main error handling method
   */
  handleError(error, options = {}) {
    const {
      type = this.categorizeError(error),
      severity = this.determineSeverity(error),
      context = '',
      userId = null,
      metadata = {}
    } = options;

    // Create standardized error object
    const errorInfo = {
      type,
      severity,
      message: this.extractErrorMessage(error),
      userMessage: this.getUserFriendlyMessage(error, type),
      stack: error?.stack,
      context,
      userId,
      metadata: {
        ...metadata,
        timestamp: new Date().toISOString(),
        userAgent: typeof navigator !== 'undefined' ? navigator.userAgent : '',
        url: typeof window !== 'undefined' ? window.location.href : ''
      }
    };

    // Log error in development
    if (process.env.NODE_ENV === 'development') {
      logger.error('Error Details:', errorInfo);
    }

    // Check rate limiting
    if (this.shouldReportError(errorInfo)) {
      // Report to error tracking service
      this.reportToErrorTracking(errorInfo);
    }

    // Execute specific handler or default
    const handler = this.errorHandlers.get(type) || this.defaultHandler;
    return handler(errorInfo);
  }

  /**
   * Categorize error based on its properties
   */
  categorizeError(error) {
    // Network errors
    if (error?.code === 'ECONNABORTED' || 
        error?.message?.includes('Network') ||
        error?.message?.includes('fetch')) {
      return ErrorTypes.NETWORK;
    }

    // Authentication errors
    if (error?.response?.status === 401 || 
        error?.response?.status === 403) {
      return ErrorTypes.AUTHENTICATION;
    }

    // Not found errors
    if (error?.response?.status === 404) {
      return ErrorTypes.NOT_FOUND;
    }

    // Validation errors
    if (error?.response?.status === 400 || 
        error?.response?.status === 422) {
      return ErrorTypes.VALIDATION;
    }

    // Server errors
    if (error?.response?.status >= 500) {
      return ErrorTypes.SERVER;
    }

    // Default to unknown
    return ErrorTypes.UNKNOWN;
  }

  /**
   * Determine error severity
   */
  determineSeverity(error) {
    const status = error?.response?.status;

    if (status >= 500 || error?.code === 'ECONNREFUSED') {
      return ErrorSeverity.CRITICAL;
    }
    if (status === 401 || status === 403) {
      return ErrorSeverity.HIGH;
    }
    if (status === 404 || status === 400) {
      return ErrorSeverity.MEDIUM;
    }
    return ErrorSeverity.LOW;
  }

  /**
   * Extract error message
   */
  extractErrorMessage(error) {
    if (error?.response?.data?.message) {
      return error.response.data.message;
    }
    if (error?.message) {
      return error.message;
    }
    if (typeof error === 'string') {
      return error;
    }
    return 'An unexpected error occurred';
  }

  /**
   * Get user-friendly error message
   */
  getUserFriendlyMessage(error, type) {
    const messages = {
      [ErrorTypes.NETWORK]: 'Unable to connect. Please check your internet connection.',
      [ErrorTypes.AUTHENTICATION]: 'Please log in to continue.',
      [ErrorTypes.PERMISSION]: 'You don\'t have permission to perform this action.',
      [ErrorTypes.NOT_FOUND]: 'The requested resource was not found.',
      [ErrorTypes.VALIDATION]: 'Please check your input and try again.',
      [ErrorTypes.SERVER]: 'Something went wrong on our end. Please try again later.',
      [ErrorTypes.UNKNOWN]: 'An unexpected error occurred. Please try again.'
    };

    return messages[type] || messages[ErrorTypes.UNKNOWN];
  }

  /**
   * Check if error should be reported (rate limiting)
   */
  shouldReportError(errorInfo) {
    const key = `${errorInfo.type}-${errorInfo.message}`;
    const now = Date.now();
    
    // Clean old entries
    for (const [cacheKey, data] of errorReportingCache.entries()) {
      if (now - data.firstSeen > ERROR_CACHE_DURATION) {
        errorReportingCache.delete(cacheKey);
      }
    }

    // Check if we should report
    const cached = errorReportingCache.get(key);
    if (cached) {
      cached.count++;
      if (cached.count > MAX_ERRORS_PER_TYPE) {
        return false;
      }
    } else {
      errorReportingCache.set(key, {
        firstSeen: now,
        count: 1
      });
    }

    return true;
  }

  /**
   * Report error to tracking service (Sentry, etc.)
   */
  reportToErrorTracking(errorInfo) {
    // Use our Sentry integration
    captureError(errorInfo.error || new Error(errorInfo.message), {
      level: errorInfo.severity,
      tags: {
        errorType: errorInfo.type,
        context: errorInfo.context,
        category: errorInfo.category
      },
      extra: {
        ...errorInfo.metadata,
        timestamp: errorInfo.timestamp,
        userMessage: errorInfo.userMessage
      },
      user: errorInfo.userId ? { id: errorInfo.userId } : undefined
    });
    
    // Add breadcrumb for better debugging
    addBreadcrumb({
      category: 'error',
      message: errorInfo.message,
      level: errorInfo.severity,
      data: {
        type: errorInfo.type,
        context: errorInfo.context
      }
    });
  }

  /**
   * Default error handler
   */
  handleGenericError(errorInfo) {
    // In production, you might want to show a toast notification
    if (process.env.NODE_ENV === 'development') {
      console.error('[GlobalErrorHandler]', errorInfo);
    }
    
    return {
      success: false,
      error: errorInfo.userMessage,
      errorCode: errorInfo.type
    };
  }

  /**
   * Handle async errors in a standardized way
   */
  async handleAsync(asyncFn, options = {}) {
    try {
      return await asyncFn();
    } catch (error) {
      return this.handleError(error, options);
    }
  }
}

// Create singleton instance
export const errorHandler = new GlobalErrorHandler();

// Register specific handlers
errorHandler.registerHandler(ErrorTypes.AUTHENTICATION, (errorInfo) => {
  // Clear auth tokens and redirect to login
  if (typeof window !== 'undefined') {
    // Import dynamically to avoid circular dependencies
    import('./secureTokenStorage').then(({ secureTokenStorage }) => {
      secureTokenStorage.clearTokens();
      window.location.href = '/login';
    });
  }
  return { success: false, error: errorInfo.userMessage };
});

errorHandler.registerHandler(ErrorTypes.NETWORK, (errorInfo) => {
  // Could implement retry logic here
  return {
    success: false,
    error: errorInfo.userMessage,
    retry: true
  };
});

// Export utilities
export const handleError = errorHandler.handleError.bind(errorHandler);
export const handleAsync = errorHandler.handleAsync.bind(errorHandler);
export const registerErrorHandler = errorHandler.registerHandler.bind(errorHandler);

export default errorHandler;