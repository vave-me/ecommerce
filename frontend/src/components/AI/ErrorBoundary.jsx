import React from 'react';
import { Bot, AlertCircle, RefreshCw } from '@/icons';
import styles from './ErrorBoundary.module.css';

class AIErrorBoundary extends React.Component {
    constructor(props) {
        super(props);
        this.state = { hasError: false, error: null, errorInfo: null };
    }

    static getDerivedStateFromError(error) {
        return { hasError: true };
    }

    componentDidCatch(error, errorInfo) {
        
        this.setState({
            error,
            errorInfo
        });
    }

    handleReset = () => {
        this.setState({ hasError: false, error: null, errorInfo: null });
        // Optionally reload the page
        if (this.props.onReset) {
            this.props.onReset();
        }
    };

    render() {
        if (this.state.hasError) {
            return (
                <div className={styles.errorContainer}>
                    <div className={styles.errorContent}>
                        <AlertCircle className={styles.errorIcon} />
                        <h2 className={styles.errorTitle}>AI Assistant Error</h2>
                        <p className={styles.errorMessage}>
                            {this.props.fallbackMessage || 'Something went wrong with the AI Assistant. Please try refreshing the page.'}
                        </p>
                        {process.env.NODE_ENV === 'development' && this.state.error && (
                            <details className={styles.errorDetails}>
                                <summary>Error Details</summary>
                                <pre>{this.state.error.toString()}</pre>
                                <pre>{this.state.errorInfo?.componentStack}</pre>
                            </details>
                        )}
                        <button
                            className={styles.resetButton}
                            onClick={this.handleReset}
                        >
                            <RefreshCw />
                            Try Again
                        </button>
                    </div>
                </div>
            );
        }

        return this.props.children;
    }
}

export default AIErrorBoundary;