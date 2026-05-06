/**
 * PRODUCTION ERROR STATE SYSTEM
 * Comprehensive error handling with improved UX
 */
import React, { memo } from 'react';
import { AlertTriangle, RefreshCw, Home, ArrowLeft } from 'lucide-react';
import styles from './ErrorState.module.css';
const ErrorState = memo(({
  title = 'Something went wrong',
  message = 'We encountered an unexpected error. Please try again.',
  error = null,
  onRetry = null,
  onGoBack = null,
  onGoHome = null,
  showDetails = false,
  variant = 'default',
  size = 'medium',
  className = '',
  ...props
}) => {
  const [showErrorDetails, setShowErrorDetails] = React.useState(false);
  const [isRetrying, setIsRetrying] = React.useState(false);
  const handleRetry = async () => {
    if (!onRetry || isRetrying) return;
    setIsRetrying(true);
    try {
      await onRetry();
    } catch (err) {
        // Handle error silently for better UX
        if (process.env.NODE_ENV === 'development') {
            console.error('Event handler error:', err);
        }
    } finally {
      setIsRetrying(false);
    }
  };
  const containerClass = `
    ${styles.container} 
    ${styles[size]} 
    ${styles[variant]} 
    ${className}
  `.trim();
  const errorDetails = error?.message || error?.toString?.() || '';
  const shouldShowDetails = showDetails && errorDetails && process.env.NODE_ENV === 'development';
  return (
    <div className={containerClass} role="alert" {...props}>
      <div className={styles.iconContainer}>
        <AlertTriangle className={styles.errorIcon} size={size === 'small' ? 32 : size === 'large' ? 48 : 40} />
      </div>
      <div className={styles.content}>
        <h3 className={styles.title}>{title}</h3>
        <p className={styles.message}>{message}</p>
        {shouldShowDetails && (
          <div className={styles.errorDetails}>
            <button
              className={styles.detailsToggle}
              onClick={() => setShowErrorDetails(!showErrorDetails)}
              type="button"
            >
              {showErrorDetails ? 'Hide' : 'Show'} technical details
            </button>
            {showErrorDetails && (
              <pre className={styles.errorText}>
                {errorDetails}
              </pre>
            )}
          </div>
        )}
      </div>
      <div className={styles.actions}>
        {onRetry && (
          <button
            className={`${styles.button} ${styles.primary}`}
            onClick={handleRetry}
            disabled={isRetrying}
            type="button"
          >
            <RefreshCw 
              className={`${styles.buttonIcon} ${isRetrying ? styles.spinning : ''}`} 
              size={16} 
            />
            {isRetrying ? 'Retrying...' : 'Try Again'}
          </button>
        )}
        {onGoBack && (
          <button
            className={`${styles.button} ${styles.secondary}`}
            onClick={onGoBack}
            type="button"
          >
            <ArrowLeft className={styles.buttonIcon} size={16} />
            Go Back
          </button>
        )}
        {onGoHome && (
          <button
            className={`${styles.button} ${styles.secondary}`}
            onClick={onGoHome}
            type="button"
          >
            <Home className={styles.buttonIcon} size={16} />
            Go Home
          </button>
        )}
      </div>
    </div>
  );
});
ErrorState.displayName = 'ErrorState';
export default ErrorState; 