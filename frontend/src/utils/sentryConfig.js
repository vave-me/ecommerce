// Sentry configuration - only loads if Sentry is installed
let Sentry = null;
let ExtraErrorData = null;

try {
    Sentry = require('@sentry/nextjs');
    const integrations = require('@sentry/integrations');
    ExtraErrorData = integrations.ExtraErrorData;
} catch (e) {
    // Sentry not installed - will use fallback implementations
    console.info('Sentry not installed - error tracking disabled');
}

/**
 * Sentry configuration for production error tracking
 */

const SENTRY_DSN = process.env.NEXT_PUBLIC_SENTRY_DSN;
const ENVIRONMENT = process.env.NEXT_PUBLIC_ENV || 'development';
const RELEASE = process.env.NEXT_PUBLIC_APP_VERSION || 'unknown';

// Only initialize in production or staging
export const initSentry = () => {
    if (!Sentry || !SENTRY_DSN || ENVIRONMENT === 'development') {
        console.info('Sentry is disabled');
        return;
    }

    Sentry.init({
        dsn: SENTRY_DSN,
        environment: ENVIRONMENT,
        release: RELEASE,
        
        // Performance Monitoring
        tracesSampleRate: ENVIRONMENT === 'production' ? 0.1 : 1.0,
        
        // Session Replay
        replaysSessionSampleRate: 0.1,
        replaysOnErrorSampleRate: 1.0,
        
        // Integrations
        integrations: [
            new ExtraErrorData({ depth: 10 }),
            new Sentry.BrowserTracing({
                // Set sampling rate for performance monitoring
                tracingOrigins: ['localhost', /^https:\/\/yourapp\.com\/api/],
                // Capture interactions
                routingInstrumentation: Sentry.nextRouterInstrumentation,
            }),
            new Sentry.Replay({
                // Mask sensitive content
                maskAllText: false,
                maskAllInputs: true,
                blockAllMedia: false,
                // Capture console logs in replays
                beforeAddRecordingEvent: (event) => {
                    // Don't record sensitive events
                    if (event.data?.tag === 'password') {
                        return null;
                    }
                    return event;
                },
            }),
        ],
        
        // Filtering
        ignoreErrors: [
            // Browser extensions
            'top.GLOBALS',
            'ResizeObserver loop limit exceeded',
            'Non-Error promise rejection captured',
            // Network errors that are expected
            /NetworkError/i,
            /fetch.*failed/i,
            // User-caused errors
            'User denied permission',
            'User cancelled',
        ],
        
        denyUrls: [
            // Chrome extensions
            /extensions\//i,
            /^chrome:\/\//i,
            /^moz-extension:\/\//i,
            // Other browser extensions
            /^safari-extension:\/\//i,
        ],
        
        // Data scrubbing
        beforeSend(event, hint) {
            // Filter out non-errors
            if (event.exception && !event.exception.values?.[0]?.value) {
                return null;
            }
            
            // Add user context if available
            const user = getUserContext();
            if (user) {
                event.user = {
                    id: user.id,
                    username: user.username,
                    // Don't send email or other PII
                };
            }
            
            // Add custom context
            event.contexts = {
                ...event.contexts,
                app: {
                    build_time: process.env.NEXT_PUBLIC_BUILD_TIME,
                    deployment: process.env.NEXT_PUBLIC_DEPLOYMENT_ID,
                },
            };
            
            // Sanitize sensitive data
            if (event.request?.cookies) {
                event.request.cookies = '[Filtered]';
            }
            
            if (event.extra) {
                event.extra = sanitizeData(event.extra);
            }
            
            return event;
        },
        
        // Breadcrumbs configuration
        beforeBreadcrumb(breadcrumb, hint) {
            // Filter out noisy breadcrumbs
            if (breadcrumb.category === 'console' && breadcrumb.level === 'debug') {
                return null;
            }
            
            // Don't log navigation to sensitive pages
            if (breadcrumb.category === 'navigation') {
                const sensitiveRoutes = ['/payment', '/checkout', '/admin'];
                if (sensitiveRoutes.some(route => breadcrumb.data?.to?.includes(route))) {
                    breadcrumb.data = { ...breadcrumb.data, to: '[Filtered]' };
                }
            }
            
            return breadcrumb;
        },
    });
};

/**
 * Get user context for Sentry
 */
function getUserContext() {
    // This should be implemented based on your auth system
    // Example:
    if (typeof window !== 'undefined') {
        const userStr = localStorage.getItem('user');
        if (userStr) {
            try {
                const user = JSON.parse(userStr);
                return {
                    id: user.userId,
                    username: user.username,
                };
            } catch {
                return null;
            }
        }
    }
    return null;
}

/**
 * Sanitize sensitive data from objects
 */
function sanitizeData(data) {
    const sensitiveKeys = [
        'password', 'token', 'secret', 'api_key', 'apiKey',
        'credit_card', 'creditCard', 'ssn', 'social_security'
    ];
    
    const sanitized = { ...data };
    
    Object.keys(sanitized).forEach(key => {
        const lowerKey = key.toLowerCase();
        if (sensitiveKeys.some(sensitive => lowerKey.includes(sensitive))) {
            sanitized[key] = '[Filtered]';
        } else if (typeof sanitized[key] === 'object' && sanitized[key] !== null) {
            sanitized[key] = sanitizeData(sanitized[key]);
        }
    });
    
    return sanitized;
}

/**
 * Custom error capture with additional context
 */
export function captureError(error, context = {}) {
    if (!Sentry || !SENTRY_DSN || ENVIRONMENT === 'development') {
        console.error('Error captured:', error, context);
        return;
    }
    
    Sentry.withScope((scope) => {
        // Set error level based on error type
        if (context.level) {
            scope.setLevel(context.level);
        } else if (error.status >= 500) {
            scope.setLevel('error');
        } else if (error.status >= 400) {
            scope.setLevel('warning');
        }
        
        // Add context
        if (context.tags) {
            Object.entries(context.tags).forEach(([key, value]) => {
                scope.setTag(key, value);
            });
        }
        
        if (context.user) {
            scope.setUser(context.user);
        }
        
        if (context.extra) {
            Object.entries(context.extra).forEach(([key, value]) => {
                scope.setExtra(key, value);
            });
        }
        
        // Capture the error
        Sentry.captureException(error);
    });
}

/**
 * Capture custom messages
 */
export function captureMessage(message, level = 'info', context = {}) {
    if (!Sentry || !SENTRY_DSN || ENVIRONMENT === 'development') {
        console.log(`[${level}] ${message}`, context);
        return;
    }
    
    Sentry.withScope((scope) => {
        scope.setLevel(level);
        
        if (context.tags) {
            Object.entries(context.tags).forEach(([key, value]) => {
                scope.setTag(key, value);
            });
        }
        
        if (context.extra) {
            Object.entries(context.extra).forEach(([key, value]) => {
                scope.setExtra(key, value);
            });
        }
        
        Sentry.captureMessage(message, level);
    });
}

/**
 * Performance monitoring
 */
export function startTransaction(name, op = 'navigation') {
    if (!Sentry || !SENTRY_DSN || ENVIRONMENT === 'development') {
        return null;
    }
    
    return Sentry.startTransaction({
        name,
        op,
    });
}

/**
 * Add breadcrumb for user actions
 */
export function addBreadcrumb(breadcrumb) {
    if (!Sentry || !SENTRY_DSN || ENVIRONMENT === 'development') {
        return;
    }
    
    Sentry.addBreadcrumb({
        timestamp: Date.now(),
        ...breadcrumb,
    });
}

// Auto-initialize on import
if (typeof window !== 'undefined') {
    initSentry();
}