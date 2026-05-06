import React from 'react';
import { errorHandler } from '../../utils/globalErrorHandler';
import styles from './ErrorBoundary.module.css';

class ErrorBoundary extends React.Component {
  constructor(props) {
    super(props);
    this.state = { 
      hasError: false, 
      error: null,
      errorInfo: null,
      errorId: null
    };
  }

  static getDerivedStateFromError(error) {
    // Update state so the next render will show the fallback UI
    return { hasError: true };
  }

  componentDidCatch(error, errorInfo) {
    // Generate error ID for user reference
    const errorId = `ERR_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
    
    // Log error to global error handler
    errorHandler.handleError(error, {
      context: 'React Error Boundary',
      metadata: {
        componentStack: errorInfo.componentStack,
        errorBoundary: this.props.name || 'Unknown',
        errorId
      }
    });

    // Update state with error details
    this.setState({
      error,
      errorInfo,
      errorId
    });
  }

  handleReset = () => {
    this.setState({ 
      hasError: false, 
      error: null,
      errorInfo: null,
      errorId: null
    });
  };

  render() {
    if (this.state.hasError) {
      // Custom fallback UI
      if (this.props.fallback) {
        return this.props.fallback(
          this.state.error, 
          this.state.errorInfo, 
          this.handleReset
        );
      }

      // Default fallback UI
      return (
        <div className={styles.errorBoundary}>
          <div className={styles.errorContent}>
            <div className={styles.errorIcon}>
              <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                <circle cx="12" cy="12" r="10" strokeWidth="2"/>
                <line x1="12" y1="8" x2="12" y2="12" strokeWidth="2"/>
                <line x1="12" y1="16" x2="12.01" y2="16" strokeWidth="2"/>
              </svg>
            </div>
            
            <h2 className={styles.errorTitle}>Something went wrong</h2>
            
            <p className={styles.errorMessage}>
              We're sorry, but something unexpected happened. 
              {this.props.showDetails && this.state.error?.message && (
                <span className={styles.errorDetails}>
                  {' '}Error: {this.state.error.message}
                </span>
              )}
            </p>

            {this.state.errorId && (
              <p className={styles.errorId}>
                Error ID: <code>{this.state.errorId}</code>
              </p>
            )}

            <div className={styles.errorActions}>
              <button
                onClick={this.handleReset}
                className={styles.resetButton}
              >
                Try Again
              </button>
              
              {this.props.showReload && (
                <button
                  onClick={() => window.location.reload()}
                  className={styles.reloadButton}
                >
                  Reload Page
                </button>
              )}
            </div>

            {process.env.NODE_ENV === 'development' && this.state.errorInfo && (
              <details className={styles.errorDebug}>
                <summary>Error Details (Development Only)</summary>
                <pre>{this.state.error && this.state.error.toString()}</pre>
                <pre>{this.state.errorInfo.componentStack}</pre>
              </details>
            )}
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}

// HOC for functional components
export function withErrorBoundary(Component, errorBoundaryProps) {
  const WrappedComponent = (props) => (
    <ErrorBoundary {...errorBoundaryProps}>
      <Component {...props} />
    </ErrorBoundary>
  );
  
  WrappedComponent.displayName = `withErrorBoundary(${Component.displayName || Component.name})`;
  
  return WrappedComponent;
}

export default ErrorBoundary;