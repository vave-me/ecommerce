import React, {memo} from 'react';
import {AlertCircle, RefreshCw} from '@/icons';
import styles from './ErrorPlaceholder.module.css';
/**
 * ErrorPlaceholder - Atomic Design Component
 * Displays an error state with message and retry option
 *
 * @param {Object} props - Component props
 * @param {string} props.message - Error message
 * @param {Function} props.onRetry - Retry function
 * @param {string} props.title - Error title
 * @param {string} props.size - Size variant (sm, md, lg)
 * @returns {JSX.Element} Rendered error placeholder
 */
const ErrorPlaceholder = memo(({
                                   message = 'Something went wrong. Please try again.',
                                   onRetry = null,
                                   title = 'Error',
                                   size = 'md'
                               }) => {
    return (
        <div className={styles.container}>
            <div className={`${styles.content} ${styles[size]}`}>
                <div className={styles.iconWrapper}>
                    <AlertCircle
                        className={`${styles.icon} ${styles[`icon${size.charAt(0).toUpperCase() + size.slice(1)}`]}`}/>
                </div>
                <div className={styles.textContent}>
                    <h3 className={styles.title}>{title}</h3>
                    <p className={styles.message}>{message}</p>
                </div>
                {onRetry && (
                    <button
                        className={styles.retryButton}
                        onClick={onRetry}
                        aria-label="Retry"
                    >
                        <RefreshCw className={styles.retryIcon} size={16}/>
                        <span>Try Again</span>
                    </button>
                )}
            </div>
        </div>
    );
});
ErrorPlaceholder.displayName = 'ErrorPlaceholder';
export default ErrorPlaceholder; 