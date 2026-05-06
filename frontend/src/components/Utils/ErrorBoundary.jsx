// src/components/ErrorBoundary.jsx
import React from 'react';
import PropTypes from 'prop-types';
import styles from './ErrorBoundary.module.css';
class ErrorBoundary extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            hasError: false,
            error: null
        };
    }
    static getDerivedStateFromError(error) {
        // Update state to render fallback UI
        return {
            hasError: true,
            error
        };
    }
    componentDidCatch(error, errorInfo) {
        // You can log the error to an error reporting service here
        // Call onError if provided
        if (this.props.onError) {
            this.props.onError(error);
        }
    }
    render() {
        if (this.state.hasError) {
            return (
                <div className={styles.fallbackContainer}>
                    <h1>Something went wrong.</h1>
                    <p>Please try refreshing the page or contact support if the problem persists.</p>
                    {this.props.debug && (
                        <div className={styles.errorDetails}>
                            <h2>Error Details:</h2>
                            <p>{this.state.error && this.state.error.message}</p>
                        </div>
                    )}
                </div>
            );
        }
        return this.props.children;
    }
}
ErrorBoundary.propTypes = {
    children: PropTypes.node.isRequired,
    debug: PropTypes.bool,
    onError: PropTypes.func
};
ErrorBoundary.defaultProps = {
    debug: false
};
export default ErrorBoundary;
