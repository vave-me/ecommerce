// This file configures the initialization of Sentry for edge features (middleware, edge routes, and so on).
// https://docs.sentry.io/platforms/javascript/guides/nextjs/

try {
  const Sentry = require('@sentry/nextjs');

Sentry.init({
  dsn: process.env.NEXT_PUBLIC_SENTRY_DSN || process.env.SENTRY_DSN,

  // Adjust this value in production, or use tracesSampler for greater control
  tracesSampleRate: process.env.NODE_ENV === 'production' ? 0.1 : 1.0,

  // Setting this option to true will print useful information to the console while you're setting up Sentry.
  debug: false,
  
  // Edge runtime specific configuration
  environment: process.env.NODE_ENV,
  
  // Edge runtime doesn't support all integrations
  integrations: [
    // Add edge-compatible integrations here
  ],
  
  beforeSend(event) {
    // Don't send events in development
    if (process.env.NODE_ENV === 'development') {
      return null;
    }
    
    return event;
  },
});
} catch (e) {
  // Sentry not installed
  console.info('Sentry not installed - edge error tracking disabled');
}