import React, { memo } from 'react';
import styles from './LoadingStates.module.css';

/**
 * Spinner loading component
 */
export const LoadingSpinner = memo(function LoadingSpinner({ 
    size = 'medium', 
    color = 'primary',
    text = 'Loading...',
    showText = true
}) {
    return (
        <div className={`${styles.spinnerContainer} ${styles[size]}`}>
            <div className={`${styles.spinner} ${styles[color]}`} />
            {showText && <p className={styles.loadingText}>{text}</p>}
        </div>
    );
});

/**
 * Skeleton loader for content placeholders
 */
export const SkeletonLoader = memo(function SkeletonLoader({ 
    count = 1, 
    height = 20,
    width = '100%',
    className = ''
}) {
    return (
        <div className={`${styles.skeletonContainer} ${className}`}>
            {Array.from({ length: count }).map((_, index) => (
                <div
                    key={index}
                    className={styles.skeleton}
                    style={{ height: `${height}px`, width }}
                />
            ))}
        </div>
    );
});

/**
 * Card skeleton loader
 */
export const CardSkeleton = memo(function CardSkeleton({ count = 1 }) {
    return (
        <div className={styles.cardSkeletonContainer}>
            {Array.from({ length: count }).map((_, index) => (
                <div key={index} className={styles.cardSkeleton}>
                    <div className={styles.imageSkeleton} />
                    <div className={styles.contentSkeleton}>
                        <SkeletonLoader height={24} width="70%" />
                        <SkeletonLoader height={16} width="100%" count={2} />
                        <SkeletonLoader height={20} width="30%" />
                    </div>
                </div>
            ))}
        </div>
    );
});

/**
 * Full page loader
 */
export const PageLoader = memo(function PageLoader({ text = 'Loading page...' }) {
    return (
        <div className={styles.pageLoader}>
            <LoadingSpinner size="large" text={text} />
        </div>
    );
});

/**
 * Inline loader for buttons/forms
 */
export const InlineLoader = memo(function InlineLoader({ className = '' }) {
    return (
        <span className={`${styles.inlineLoader} ${className}`}>
            <span className={styles.dot} />
            <span className={styles.dot} />
            <span className={styles.dot} />
        </span>
    );
});

/**
 * Progress bar loader
 */
export const ProgressLoader = memo(function ProgressLoader({ 
    progress = 0, 
    showPercentage = true,
    color = 'primary' 
}) {
    const clampedProgress = Math.min(100, Math.max(0, progress));
    
    return (
        <div className={styles.progressContainer}>
            <div className={styles.progressBar}>
                <div 
                    className={`${styles.progressFill} ${styles[color]}`}
                    style={{ width: `${clampedProgress}%` }}
                />
            </div>
            {showPercentage && (
                <span className={styles.progressText}>{clampedProgress}%</span>
            )}
        </div>
    );
});