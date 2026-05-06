import React, {memo} from 'react';
import {Loader2} from '@/icons';
import styles from './LoadingPlaceholder.module.css';
/**
 * LoadingPlaceholder - Atomic Design Component
 * Displays a loading state with spinner and message
 *
 * @param {Object} props - Component props
 * @param {string} props.label - Loading message
 * @param {string} props.size - Size variant (sm, md, lg)
 * @param {boolean} props.fullScreen - Whether to show full screen loading
 * @returns {JSX.Element} Rendered loading placeholder
 */
const LoadingPlaceholder = memo(({
                                     label = 'Loading...',
                                     size = 'md',
                                     fullScreen = false
                                 }) => {
    const containerClass = fullScreen
        ? `${styles.container} ${styles.fullScreen}`
        : styles.container;
    return (
        <div className={containerClass}>
            <div className={`${styles.content} ${styles[size]}`}>
                <div className={styles.spinnerWrapper}>
                    <Loader2
                        className={`${styles.spinner} ${styles[`spinner${size.charAt(0).toUpperCase() + size.slice(1)}`]}`}/>
                </div>
                <p className={styles.label}>{label}</p>
            </div>
        </div>
    );
});
LoadingPlaceholder.displayName = 'LoadingPlaceholder';
export default LoadingPlaceholder; 