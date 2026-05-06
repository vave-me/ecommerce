"use client";
import React from "react";
import { useTranslations } from "next-intl";
import { AlertCircle, RefreshCw } from "@/icons";
import styles from "./WishlistErrorBoundary.module.css";
/**
 * Error Boundary for Wishlist components
 * Catches and handles errors gracefully
 */
class WishlistErrorBoundary extends React.Component {
    constructor(props) {
        super(props);
        this.state = { hasError: false, error: null };
    }
    static getDerivedStateFromError(error) {
        // Update state so the next render will show the fallback UI
        return { hasError: true, error };
    }
    componentDidCatch(error, errorInfo) {
        // Log error to console in development
        if (process.env.NODE_ENV === 'development') {
        }
        // You can also log the error to an error reporting service here
    }
    handleReset = () => {
        this.setState({ hasError: false, error: null });
        // Optionally reload the page
        if (this.props.onReset) {
            this.props.onReset();
        }
    };
    render() {
        if (this.state.hasError) {
            return (
                <ErrorFallback 
                    error={this.state.error} 
                    onReset={this.handleReset}
                />
            );
        }
        return this.props.children;
    }
}
/**
 * Fallback component shown when an error occurs
 */
function ErrorFallback({ error, onReset }) {
    const t = useTranslations('Wishlist');
    return (
        <div className={styles.errorContainer}>
            <div className={styles.errorContent}>
                <AlertCircle size={48} className={styles.errorIcon} />
                <h2 className={styles.errorTitle}>{t('errorBoundaryTitle', 'Something went wrong')}</h2>
                <p className={styles.errorMessage}>
                    {t('errorBoundaryMessage', 'We encountered an error while loading your wishlist.')}
                </p>
                {process.env.NODE_ENV === 'development' && error && (
                    <details className={styles.errorDetails}>
                        <summary>{t('errorDetails', 'Error details')}</summary>
                        <pre>{error.toString()}</pre>
                    </details>
                )}
                <button 
                    onClick={onReset}
                    className={styles.resetButton}
                >
                    <RefreshCw size={16} />
                    <span>{t('tryAgain', 'Try Again')}</span>
                </button>
            </div>
        </div>
    );
}
export default WishlistErrorBoundary;