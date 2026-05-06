// src/components/Skeletons/MediaSkeleton.jsx
import React, { memo } from 'react';
import styles from './MediaSkeleton.module.css';
/**
 * Skeleton loader for media content while it's being loaded
 * Displays a pulsing animation in the same dimensions as the actual media
 */
const MediaSkeleton = memo(() => {
    return (
        <div className={styles.skeletonContainer}>
            <div className={styles.thumbnailColumnSkeleton}>
                {/* Generate multiple thumbnail skeletons */}
                {[...Array(4)].map((_, index) => (
                    <div
                        key={`thumb-skeleton-${index}`}
                        className={styles.thumbnailSkeleton}
                        aria-hidden="true"
                    />
                ))}
            </div>
            <div className={styles.mainMediaSkeleton}>
                <div className={styles.shimmer} />
            </div>
        </div>
    );
});
MediaSkeleton.displayName = 'MediaSkeleton';
export default MediaSkeleton;