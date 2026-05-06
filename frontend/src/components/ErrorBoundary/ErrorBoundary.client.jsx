"use client";
import React from 'react';
import PropTypes from 'prop-types';
import { toast } from 'react-toastify';
class ErrorBoundary extends React.Component {
  constructor(props) {
    super(props);
    this.state = { hasError: false, error: null, errorInfo: null };
  }
  static getDerivedStateFromError(error) {
    return { hasError: true };
  }
  componentDidCatch(error, errorInfo) {
    this.setState({
      error: error,
      errorInfo: errorInfo
    });
    // Log the error to console
    // Notify the user with toast
    toast.error('Something went wrong. The team has been notified.');
    // Here you would typically log to an error reporting service
    // e.g., Sentry, LogRocket, etc.
    if (typeof window !== 'undefined' && window.REPORT_ERROR_API) {
      try {
        window.REPORT_ERROR_API.capture({
          error: error.toString(),
          componentStack: errorInfo.componentStack,
          url: window.location.href,
          timestamp: new Date().toISOString()
        });
      } catch (loggingError) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', loggingError);
        }
    }
    }
  }
  render() {
    if (this.state.hasError) {
      // Check for custom fallback
      if (this.props.fallback) {
        return this.props.fallback(this.state.error);
      }
      // Default fallback UI
      return (
        <div style={{
          padding: '20px',
          margin: '20px 0',
          borderRadius: '8px',
          backgroundColor: '#fff8f8',
          border: '1px solid #ffcdd2',
          boxShadow: '0 2px 4px rgba(0,0,0,0.05)'
        }}>
          <h2 style={{ color: '#d32f2f' }}>Something went wrong</h2>
          <p>We've been notified and are working on a fix.</p>
          {this.props.showDetails && (
            <details style={{ whiteSpace: 'pre-wrap', margin: '10px 0' }}>
              <summary style={{ cursor: 'pointer' }}>Error Details</summary>
              <div style={{ 
                padding: '10px',
                marginTop: '10px',
                backgroundColor: '#f5f5f5', 
                border: '1px solid #e0e0e0',
                borderRadius: '4px',
                fontSize: '12px',
                fontFamily: 'monospace'
              }}>
                {this.state.error && this.state.error.toString()}
                <br />
                {this.state.errorInfo && this.state.errorInfo.componentStack}
              </div>
            </details>
          )}
          <div style={{ marginTop: '20px' }}>
            <button 
              onClick={() => window.location.reload()}
              style={{
                padding: '8px 16px',
                backgroundColor: '#2196f3',
                color: 'white',
                border: 'none',
                borderRadius: '4px',
                cursor: 'pointer'
              }}
            >
              Reload Page
            </button>
            {this.props.resetError && (
              <button
                onClick={this.props.resetError}
                style={{
                  padding: '8px 16px',
                  backgroundColor: 'transparent',
                  color: '#2196f3',
                  border: '1px solid #2196f3',
                  borderRadius: '4px',
                  marginLeft: '10px',
                  cursor: 'pointer'
                }}
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
  resetError: PropTypes.func,
  showDetails: PropTypes.bool
};
ErrorBoundary.defaultProps = {
  showDetails: process.env.NODE_ENV === 'development'
};
export default ErrorBoundary; 