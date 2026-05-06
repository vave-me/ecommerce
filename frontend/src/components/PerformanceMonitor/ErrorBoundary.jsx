'use client';
import React from 'react';
import { performanceMetrics } from './performanceMetrics';
export class ErrorBoundary extends React.Component {
  constructor(props) {
    super(props);
    this.state = { hasError: false, error: null };
  }
  static getDerivedStateFromError(error) {
    return { hasError: true, error };
  }
  componentDidCatch(error, errorInfo) {
    performanceMetrics.trackError(this.props.componentName || 'Unknown', error);
    // Log to console in development
    if (process.env.NODE_ENV === 'development') {
    }
  }
  render() {
    if (this.state.hasError) {
      return this.props.fallback || (
        <div className="error-boundary">
          <h2>Something went wrong</h2>
          <button onClick={() => this.setState({ hasError: false })}>
            Try again
          </button>
        </div>
      );
    }
    return this.props.children;
  }
} 