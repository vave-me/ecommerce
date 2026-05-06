/**
 * Unified Loading Component
 * Consolidates all loading states into a single, flexible component
 * 
 * Usage:
 * - Simple spinner: <LoadingSpinner />
 * - With message: <LoadingSpinner message="Loading products..." />
 * - Full screen: <LoadingSpinner fullScreen />
 * - Custom size: <LoadingSpinner size="large" />
 * - In buttons: <LoadingSpinner inline size="small" />
 */
import React, { memo } from 'react';
import { FaSpinner } from 'react-icons/fa';
import styles from './LoadingSpinner.module.css';

const LoadingSpinner = memo(({
  message = '',
  size = 'medium',
  variant = 'default',
  fullScreen = false,
  inline = false,
  className = '',
  showSpinner = true,
  ...props
}) => {
  // For inline usage in buttons/forms
  if (inline) {
    return (
      <FaSpinner 
        className={`${styles.inlineSpinner} ${styles[size]} ${className}`}
        aria-label="Loading"
        {...props}
      />
    );
  }

  const containerClass = `
    ${styles.container} 
    ${styles[size]} 
    ${styles[variant]} 
    ${fullScreen ? styles.fullScreen : ''} 
    ${className}
  `.trim();

  return (
    <div className={containerClass} role="status" aria-live="polite" {...props}>
      {showSpinner && (
        <div className={styles.spinnerContainer}>
          <div className={styles.spinner} aria-hidden="true" />
        </div>
      )}
      {message && <span className={styles.message}>{message}</span>}
    </div>
  );
});

LoadingSpinner.displayName = 'LoadingSpinner';

// Export convenience components for backward compatibility
export const LoadingState = (props) => <LoadingSpinner {...props} />;
export const LoadingPlaceholder = ({ label, ...props }) => (
  <LoadingSpinner message={label} {...props} />
);
export const Spinner = () => <LoadingSpinner variant="bounce" />;

export default LoadingSpinner;