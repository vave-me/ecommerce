import { jest } from '@jest/globals';
import { errorHandler, ErrorTypes } from '../../src/utils/globalErrorHandler';
import logger from '../../src/utils/logger';

// Mock logger
jest.mock('../../src/utils/logger', () => ({
    default: {
        error: jest.fn(),
        warn: jest.fn()
    }
}));

// Mock window.Sentry
global.window = {
    Sentry: {
        captureException: jest.fn(),
        withScope: jest.fn((callback) => {
            callback({
                setContext: jest.fn(),
                setTag: jest.fn(),
                setLevel: jest.fn()
            });
        })
    }
};

describe('errorHandler', () => {
    beforeEach(() => {
        jest.clearAllMocks();
        // Reset error handler state
        errorHandler.clearHandlers();
    });

    describe('categorizeError', () => {
        it('should categorize network errors correctly', () => {
            const networkError = new Error('Network request failed');
            networkError.code = 'ECONNREFUSED';
            
            const category = errorHandler.categorizeError(networkError);
            expect(category).toBe(ErrorTypes.NETWORK);
        });

        it('should categorize authentication errors (401)', () => {
            const authError = new Error('Unauthorized');
            authError.response = { status: 401 };
            
            const category = errorHandler.categorizeError(authError);
            expect(category).toBe(ErrorTypes.AUTHENTICATION);
        });

        it('should categorize authentication errors (403)', () => {
            const authError = new Error('Forbidden');
            authError.response = { status: 403 };
            
            const category = errorHandler.categorizeError(authError);
            expect(category).toBe(ErrorTypes.AUTHENTICATION);
        });

        it('should categorize validation errors (400)', () => {
            const validationError = new Error('Bad Request');
            validationError.response = { status: 400 };
            
            const category = errorHandler.categorizeError(validationError);
            expect(category).toBe(ErrorTypes.VALIDATION);
        });

        it('should categorize not found errors (404)', () => {
            const notFoundError = new Error('Not Found');
            notFoundError.response = { status: 404 };
            
            const category = errorHandler.categorizeError(notFoundError);
            expect(category).toBe(ErrorTypes.NOT_FOUND);
        });

        it('should categorize server errors (500+)', () => {
            const serverError = new Error('Internal Server Error');
            serverError.response = { status: 500 };
            
            const category = errorHandler.categorizeError(serverError);
            expect(category).toBe(ErrorTypes.SERVER);
        });

        it('should categorize permission errors', () => {
            const permissionError = new Error('Permission denied to access resource');
            
            const category = errorHandler.categorizeError(permissionError);
            expect(category).toBe(ErrorTypes.PERMISSION);
        });

        it('should categorize timeout errors', () => {
            const timeoutError = new Error('Request timeout');
            
            const category = errorHandler.categorizeError(timeoutError);
            expect(category).toBe(ErrorTypes.TIMEOUT);
        });

        it('should return UNKNOWN for unrecognized errors', () => {
            const unknownError = new Error('Something went wrong');
            
            const category = errorHandler.categorizeError(unknownError);
            expect(category).toBe(ErrorTypes.UNKNOWN);
        });
    });

    describe('getUserMessage', () => {
        it('should return user-friendly message for network errors', () => {
            const message = errorHandler.getUserMessage(ErrorTypes.NETWORK);
            expect(message).toBe('Unable to connect. Please check your internet connection.');
        });

        it('should return user-friendly message for authentication errors', () => {
            const message = errorHandler.getUserMessage(ErrorTypes.AUTHENTICATION);
            expect(message).toBe('Please log in to continue.');
        });

        it('should return user-friendly message for validation errors', () => {
            const message = errorHandler.getUserMessage(ErrorTypes.VALIDATION);
            expect(message).toBe('Please check your input and try again.');
        });

        it('should use custom message if provided', () => {
            const customMessage = 'Custom error message';
            const message = errorHandler.getUserMessage(ErrorTypes.NETWORK, customMessage);
            expect(message).toBe(customMessage);
        });

        it('should return default message for unknown error type', () => {
            const message = errorHandler.getUserMessage('INVALID_TYPE');
            expect(message).toBe('An unexpected error occurred. Please try again.');
        });
    });

    describe('handleError', () => {
        it('should log error with logger', () => {
            const error = new Error('Test error');
            const context = { context: 'Test Context' };
            
            errorHandler.handleError(error, context);
            
            expect(logger.error).toHaveBeenCalledWith('Test error', expect.objectContaining({
                error,
                ...context,
                category: ErrorTypes.UNKNOWN,
                timestamp: expect.any(Number)
            }));
        });

        it('should send error to Sentry in production', () => {
            const originalEnv = process.env.NODE_ENV;
            process.env.NODE_ENV = 'production';
            
            const error = new Error('Production error');
            errorHandler.handleError(error);
            
            expect(window.Sentry.captureException).toHaveBeenCalledWith(error);
            
            process.env.NODE_ENV = originalEnv;
        });

        it('should not send error to Sentry in development', () => {
            process.env.NODE_ENV = 'development';
            
            const error = new Error('Dev error');
            errorHandler.handleError(error);
            
            expect(window.Sentry.captureException).not.toHaveBeenCalled();
        });

        it('should handle errors without Sentry gracefully', () => {
            const tempSentry = window.Sentry;
            delete window.Sentry;
            
            const error = new Error('Test error');
            
            // Should not throw
            expect(() => {
                errorHandler.handleError(error);
            }).not.toThrow();
            
            window.Sentry = tempSentry;
        });
    });

    describe('addHandler', () => {
        it('should add custom error handler', () => {
            const handler = jest.fn();
            
            errorHandler.addHandler(ErrorTypes.NETWORK, handler);
            
            const error = new Error('Network error');
            error.code = 'ECONNREFUSED';
            
            errorHandler.handleError(error);
            
            expect(handler).toHaveBeenCalledWith(error, expect.any(Object));
        });

        it('should call multiple handlers for the same error type', () => {
            const handler1 = jest.fn();
            const handler2 = jest.fn();
            
            errorHandler.addHandler(ErrorTypes.NETWORK, handler1);
            errorHandler.addHandler(ErrorTypes.NETWORK, handler2);
            
            const error = new Error('Network error');
            error.code = 'ECONNREFUSED';
            
            errorHandler.handleError(error);
            
            expect(handler1).toHaveBeenCalled();
            expect(handler2).toHaveBeenCalled();
        });

        it('should not call handler for different error type', () => {
            const handler = jest.fn();
            
            errorHandler.addHandler(ErrorTypes.NETWORK, handler);
            
            const error = new Error('Validation error');
            error.response = { status: 400 };
            
            errorHandler.handleError(error);
            
            expect(handler).not.toHaveBeenCalled();
        });
    });

    describe('clearHandlers', () => {
        it('should remove all custom handlers', () => {
            const handler = jest.fn();
            
            errorHandler.addHandler(ErrorTypes.NETWORK, handler);
            errorHandler.clearHandlers();
            
            const error = new Error('Network error');
            error.code = 'ECONNREFUSED';
            
            errorHandler.handleError(error);
            
            expect(handler).not.toHaveBeenCalled();
        });
    });

    describe('error context', () => {
        it('should include metadata in error context', () => {
            const error = new Error('Test error');
            const metadata = {
                userId: '123',
                action: 'fetchData'
            };
            
            errorHandler.handleError(error, {
                context: 'API Call',
                metadata
            });
            
            expect(logger.error).toHaveBeenCalledWith(
                'Test error',
                expect.objectContaining({
                    metadata
                })
            );
        });

        it('should include timestamp in error context', () => {
            const error = new Error('Test error');
            
            errorHandler.handleError(error);
            
            expect(logger.error).toHaveBeenCalledWith(
                'Test error',
                expect.objectContaining({
                    timestamp: expect.any(Number)
                })
            );
        });
    });

    describe('error recovery', () => {
        it('should handle handler errors gracefully', () => {
            const faultyHandler = jest.fn().mockImplementation(() => {
                throw new Error('Handler error');
            });
            
            errorHandler.addHandler(ErrorTypes.NETWORK, faultyHandler);
            
            const error = new Error('Network error');
            error.code = 'ECONNREFUSED';
            
            // Should not throw
            expect(() => {
                errorHandler.handleError(error);
            }).not.toThrow();
            
            // Should still log the original error
            expect(logger.error).toHaveBeenCalled();
        });
    });
});