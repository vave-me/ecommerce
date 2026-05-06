/**
 * PRODUCTION LOADING STATE SYSTEM
 * Lightweight, accessible, and consistent loading states
 * Replaces heavy animations with optimized alternatives
 */
import React, { memo } from 'react';
import styles from './LoadingState.module.css';
const LoadingState = memo(({
  message = 'Loading...',
  size = 'medium',
  variant = 'default',
  fullScreen = false,
  showSpinner = true,
  className = '',
  ...props
}) => {
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
      <span className={styles.message}>{message}</span>
    </div>
  );
});
LoadingState.displayName = 'LoadingState';
export default LoadingState; 