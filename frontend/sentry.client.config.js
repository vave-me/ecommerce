// This file configures the initialization of Sentry on the client side.
// The config you add here will be used whenever a user loads a page in their browser.
// https://docs.sentry.io/platforms/javascript/guides/nextjs/

// Only initialize if Sentry is installed
try {
  const Sentry = require('@sentry/nextjs');

Sentry.init({
  dsn: process.env.NEXT_PUBLIC_SENTRY_DSN,
  
  // Adjust this value in production, or use tracesSampler for greater control
  tracesSampleRate: process.env.NODE_ENV === 'production' ? 0.1 : 1.0,
  
  // Setting this option to true will print useful information to the console while you're setting up Sentry.
  debug: false,

  replaysOnErrorSampleRate: 1.0,
  replaysSessionSampleRate: process.env.NODE_ENV === 'production' ? 0.1 : 0.0,

  integrations: [
    new Sentry.Replay({
      maskAllText: false,
      maskAllInputs: true,
      blockAllMedia: false,
    }),
  ],
  
  // Filter out noisy errors
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
    // Next.js specific
    'NEXT_NOT_FOUND',
  ],
  
  beforeSend(event) {
    // Don't send events in development
    if (process.env.NODE_ENV === 'development') {
      return null;
    }
    
    // Filter out non-actionable errors
    if (event.exception?.values?.[0]?.value?.includes('ChunkLoadError')) {
      return null;
    }
    
    return event;
  },
});
} catch (e) {
  // Sentry not installed
  console.info('Sentry not installed - client error tracking disabled');
}