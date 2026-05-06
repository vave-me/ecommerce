// This file configures the initialization of Sentry on the server side.
// The config you add here will be used whenever the server handles a request.
// https://docs.sentry.io/platforms/javascript/guides/nextjs/

try {
  const Sentry = require('@sentry/nextjs');

Sentry.init({
  dsn: process.env.NEXT_PUBLIC_SENTRY_DSN || process.env.SENTRY_DSN,

  // Adjust this value in production, or use tracesSampler for greater control
  tracesSampleRate: process.env.NODE_ENV === 'production' ? 0.1 : 1.0,

  // Setting this option to true will print useful information to the console while you're setting up Sentry.
  debug: false,
  
  // Capture unhandled promise rejections
  onUnhandledRejection: 'warn',
  
  // Environment
  environment: process.env.NODE_ENV,
  
  // Additional server-specific configuration
  autoSessionTracking: true,
  
  // Server-side specific error filtering
  beforeSend(event) {
    // Don't send events in development
    if (process.env.NODE_ENV === 'development') {
      return null;
    }
    
    // Filter out 404s and other non-critical errors
    if (event.exception?.values?.[0]?.value?.includes('ENOENT')) {
      return null;
    }
    
    return event;
  },
});
} catch (e) {
  // Sentry not installed
  console.info('Sentry not installed - server error tracking disabled');
}