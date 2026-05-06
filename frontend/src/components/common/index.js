/**
 * Common Components Export
 * Centralized export for all unified/common components
 */

// Loading States
export { default as LoadingSpinner, LoadingState, LoadingPlaceholder, Spinner } from './LoadingSpinner';

// Error Handling
export { default as ErrorBoundary } from './ErrorBoundary';

// Empty States
export { default as EmptyState, EmptyPlaceholder } from './EmptyState';

// User Components
export { default as UserAvatar } from './UserAvatar';

// Re-export specific variants for convenience
export const InlineSpinner = (props) => <LoadingSpinner inline {...props} />;
export const FullScreenLoader = (props) => <LoadingSpinner fullScreen {...props} />;