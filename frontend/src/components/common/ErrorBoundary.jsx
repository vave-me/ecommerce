"use client";
/**
 * Unified Error Boundary Component
 * Consolidates all error boundary implementations into a single, flexible component
 * 
 * Features:
 * - Error logging to external services
 * - Toast notifications
 * - Customizable fallback UI
 * - Development-mode error details
 * - Recovery actions
 */
import React from 'react';
import PropTypes from 'prop-types';
import { toast } from 'react-toastify';
import styles from './ErrorBoundary.module.css';

class ErrorBoundary extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      hasError: false,
      error: null,
      errorInfo: null
    };
  }

  static getDerivedStateFromError(error) {
    return { hasError: true };
  }

  componentDidCatch(error, errorInfo) {
    this.setState({
      error: error,
      errorInfo: errorInfo
    });

    // Call custom error handler if provided
    if (this.props.onError) {
      this.props.onError(error, errorInfo);
    }

    // Show toast notification unless disabled
    if (!this.props.silent) {
      toast.error(this.props.errorMessage || 'Something went wrong. The team has been notified.');
    }

    // Log to external service in production
    if (typeof window !== 'undefined' && window.REPORT_ERROR_API) {
      try {
        window.REPORT_ERROR_API.capture({
          error: error.toString(),
          componentStack: errorInfo.componentStack,
          url: window.location.href,
          timestamp: new Date().toISOString(),
          component: this.props.name || 'Unknown'
        });
      } catch (loggingError) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', loggingError);
        }
    }
    }
  }

  handleReset = () => {
    this.setState({
      hasError: false,
      error: null,
      errorInfo: null
    });
    
    if (this.props.resetError) {
      this.props.resetError();
    }
  };

  render() {
    if (this.state.hasError) {
      // Use custom fallback if provided
      if (this.props.fallback) {
        return this.props.fallback(this.state.error, this.handleReset);
      }

      // Default fallback UI
      return (
        <div className={styles.fallbackContainer}>
          <div className={styles.errorIcon}>⚠️</div>
          <h2 className={styles.errorTitle}>
            {this.props.title || 'Something went wrong'}
          </h2>
          <p className={styles.errorMessage}>
            {this.props.message || 'We\'ve been notified and are working on a fix. Please try refreshing the page.'}
          </p>
          
          {(this.props.showDetails || process.env.NODE_ENV === 'development') && this.state.error && (
            <details className={styles.errorDetails}>
              <summary className={styles.errorSummary}>Error Details</summary>
              <div className={styles.errorStack}>
                <strong>Error:</strong> {this.state.error.toString()}
                {this.state.errorInfo && (
                  <>
                    <br /><br />
                    <strong>Component Stack:</strong>
                    <pre>{this.state.errorInfo.componentStack}</pre>
                  </>
                )}
              </div>
            </details>
          )}
          
          <div className={styles.errorActions}>
            <button 
              onClick={() => window.location.reload()}
              className={styles.primaryButton}
            >
              Reload Page
            </button>
            {this.props.resetError && (
              <button
                onClick={this.handleReset}
                className={styles.secondaryButton}
              >
                Try Again
              </button>
            )}
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}

ErrorBoundary.propTypes = {
  children: PropTypes.node.isRequired,
  fallback: PropTypes.func,
  onError: PropTypes.func,
  resetError: PropTypes.func,
  showDetails: PropTypes.bool,
  silent: PropTypes.bool,
  errorMessage: PropTypes.string,
  title: PropTypes.string,
  message: PropTypes.string,
  name: PropTypes.string
};

ErrorBoundary.defaultProps = {
  showDetails: false,
  silent: false
};

export default ErrorBoundary;