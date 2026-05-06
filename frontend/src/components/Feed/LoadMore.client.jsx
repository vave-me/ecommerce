"use client";

import React, { useCallback, useRef, memo } from 'react';
import styles from './LoadMore.module.css';

const LoadMore = memo(({ onLoadMore, isLoading, hasMore }) => {
    const observer = useRef();
    const lastItemRef = useCallback(node => {
        if (isLoading) return;
        if (observer.current) observer.current.disconnect();
        observer.current = new IntersectionObserver(entries => {
            if (entries[0].isIntersecting && hasMore) {
                onLoadMore();
            }
        });
        if (node) observer.current.observe(node);
    }, [isLoading, hasMore, onLoadMore]);

    return (
        <div ref={lastItemRef} className={styles.loadMoreContainer}>
            {isLoading ? (
                <div className={styles.skeletonLoading}>
                    {[...Array(3)].map((_, index) => (
                        <div key={index} className={styles.skeletonItem}>
                            <div className={styles.skeletonHeader}>
                                <div className={styles.skeletonAvatar} />
                                <div className={styles.skeletonInfo}>
                                    <div className={styles.skeletonName} />
                                    <div className={styles.skeletonTime} />
                                </div>
                            </div>
                            <div className={styles.skeletonContent} />
                            <div className={styles.skeletonMedia} />
                            <div className={styles.skeletonActions} />
                        </div>
                    ))}
                </div>
            ) : hasMore ? (
                <button 
                    onClick={onLoadMore}
                    className={styles.loadMoreButton}
                >
                    Load More
                </button>
            ) : null}
        </div>
    );
});

LoadMore.displayName = 'LoadMore';

export default LoadMore; 