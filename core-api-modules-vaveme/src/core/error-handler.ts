import { AxiosError } from 'axios';

export enum ErrorSeverity {
  LOW = 'low',
  MEDIUM = 'medium',
  HIGH = 'high',
  CRITICAL = 'critical'
}

export interface ApiError {
  success: false;
  error: string;
  userMessage: string;
  statusCode?: number;
  severity: ErrorSeverity;
  timestamp: string;
  endpoint?: string;
  operation?: string;
  details?: any;
}

export class ApiErrorHandler {
  static handle(
    error: unknown,
    endpoint?: string,
    operation?: string
  ): ApiError {
    const timestamp = new Date().toISOString();

    if (error instanceof AxiosError) {
      const statusCode = error.response?.status;
      const data = error.response?.data;

      return {
        success: false,
        error: data?.error || error.message,
        userMessage: this.getUserMessage(statusCode, data),
        statusCode,
        severity: this.getSeverity(statusCode),
        timestamp,
        endpoint,
        operation,
        details: data,
      };
    }

    // Generic error handling
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Unknown error',
      userMessage: 'An unexpected error occurred. Please try again.',
      severity: ErrorSeverity.MEDIUM,
      timestamp,
      endpoint,
      operation,
    };
  }

  private static getUserMessage(statusCode?: number, data?: any): string {
    // Check for custom error message from API
    if (data?.userMessage) return data.userMessage;
    if (data?.message) return data.message;

    // Status code based messages
    switch (statusCode) {
      case 400:
        return 'Invalid request. Please check your input and try again.';
      case 401:
        return 'You need to be logged in to perform this action.';
      case 403:
        return 'You do not have permission to perform this action.';
      case 404:
        return 'The requested resource was not found.';
      case 409:
        return 'This action conflicts with existing data.';
      case 422:
        return 'The provided data is invalid. Please check and try again.';
      case 429:
        return 'Too many requests. Please wait a moment and try again.';
      case 500:
        return 'Server error. Our team has been notified.';
      case 502:
      case 503:
        return 'Service temporarily unavailable. Please try again later.';
      default:
        return 'An error occurred. Please try again.';
    }
  }

  private static getSeverity(statusCode?: number): ErrorSeverity {
    if (!statusCode) return ErrorSeverity.MEDIUM;

    if (statusCode >= 500) return ErrorSeverity.HIGH;
    if (statusCode === 401 || statusCode === 403) return ErrorSeverity.MEDIUM;
    if (statusCode === 404) return ErrorSeverity.LOW;
    if (statusCode >= 400) return ErrorSeverity.MEDIUM;
    
    return ErrorSeverity.LOW;
  }

  static isRetryable(error: unknown): boolean {
    if (!(error instanceof AxiosError)) return false;

    const statusCode = error.response?.status;
    if (!statusCode) return true; // Network errors are retryable

    // Retryable status codes
    const retryableStatuses = [408, 429, 500, 502, 503, 504];
    return retryableStatuses.includes(statusCode);
  }

  static getRetryDelay(attempt: number, baseDelay: number = 1000): number {
    // Exponential backoff with jitter
    const exponentialDelay = baseDelay * Math.pow(2, attempt - 1);
    const jitter = Math.random() * 0.1 * exponentialDelay;
    return Math.min(exponentialDelay + jitter, 30000); // Cap at 30 seconds
  }
}