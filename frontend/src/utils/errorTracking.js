/**
 * Error tracking wrapper that works with or without Sentry
 * This provides a consistent API regardless of whether Sentry is installed
 */

const ENVIRONMENT = process.env.NEXT_PUBLIC_ENV || 'development';
const isDevelopment = ENVIRONMENT === 'development';

// Fallback implementations when Sentry is not available
const fallbackImplementations = {
    captureError: (error, context = {}) => {
        if (isDevelopment) {
            console.error('Error captured:', error, context);
        }
        // In production without Sentry, you could send to your own logging endpoint
    },
    
    captureMessage: (message, level = 'info', context = {}) => {
        if (isDevelopment) {
            console.log(`[${level}] ${message}`, context);
        }
    },
    
    addBreadcrumb: (breadcrumb) => {
        if (isDevelopment) {
            console.log('Breadcrumb:', breadcrumb);
        }
    },
    
    startTransaction: (name, op = 'navigation') => {
        if (isDevelopment) {
            console.log(`Transaction started: ${name} (${op})`);
        }
        return null;
    },
    
    setUser: (user) => {
        if (isDevelopment) {
            console.log('User context set:', user);
        }
    },
    
    setContext: (key, context) => {
        if (isDevelopment) {
            console.log(`Context set: ${key}`, context);
        }
    }
};

// Export functions that will use Sentry if available, otherwise fallbacks
export const captureError = fallbackImplementations.captureError;
export const captureMessage = fallbackImplementations.captureMessage;
export const addBreadcrumb = fallbackImplementations.addBreadcrumb;
export const startTransaction = fallbackImplementations.startTransaction;
export const setUser = fallbackImplementations.setUser;
export const setContext = fallbackImplementations.setContext;

// Initialize function (no-op without Sentry)
export const initErrorTracking = () => {
    console.info('Error tracking initialized (no Sentry)');
};