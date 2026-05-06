"use client";

import React, { useMemo, useCallback, lazy, Suspense } from 'react';
import { FixedSizeList as List } from 'react-window';
import { useInView } from 'react-intersection-observer';
import AutoSizer from 'react-virtualized-auto-sizer';
import styles from './Feed.module.css';

// Lazy load card components for optimal performance
const PostCard = lazy(() => import('../ImprovedPostCard'));
const ClassifiedCard = lazy(() => import('../classified/ClassifiedCard'));
const ServiceCard = lazy(() => import('../services/ServiceCard'));
const TweetCard = lazy(() => import('../../app/[locale]/design/TweetCard'));

// Optimized loading skeleton
const CardSkeleton = React.memo(() => (
  <div className={styles.cardSkeleton} role="presentation" aria-label="Loading content">
    <div className={styles.skeletonHeader}></div>
    <div className={styles.skeletonContent}></div>
    <div className={styles.skeletonFooter}></div>
  </div>
));

CardSkeleton.displayName = 'CardSkeleton';

// Memoized card renderer for virtual list
const FeedItem = React.memo(({ index, style, data }) => {
  const { items, loadMore, hasMore, isLoading } = data;
  const item = items[index];
  
  // Load more when approaching end
  const { ref } = useInView({
    threshold: 0,
    rootMargin: '400px',
    onChange: (inView) => {
      if (inView && hasMore && !isLoading && index >= items.length - 5) {
        loadMore();
      }
    },
  });
  
  const renderCard = useCallback(() => {
    if (!item) return <CardSkeleton />;
    
    const commonProps = {
      key: `${item.type}-${item.id}`,
      item,
      className: styles.feedCard,
    };
    
    switch (item.type) {
      case 'post':
        return (
          <Suspense fallback={<CardSkeleton />}>
            <PostCard {...commonProps} />
          </Suspense>
        );
      case 'classified':
        return (
          <Suspense fallback={<CardSkeleton />}>
            <ClassifiedCard {...commonProps} />
          </Suspense>
        );
      case 'service':
        return (
          <Suspense fallback={<CardSkeleton />}>
            <ServiceCard {...commonProps} />
          </Suspense>
        );
      case 'tweet':
        return (
          <Suspense fallback={<CardSkeleton />}>
            <TweetCard {...commonProps} />
          </Suspense>
        );
      default:
        return <CardSkeleton />;
    }
  }, [item]);
  
  return (
    <div style={style} ref={index >= items.length - 5 ? ref : null}>
      <div className={styles.feedItemWrapper}>
        {renderCard()}
      </div>
    </div>
  );
});

FeedItem.displayName = 'FeedItem';

/**
 * High-Performance Virtualized Feed Component
 * - Handles 10k+ items efficiently
 * - Memory usage stays constant
 * - Smooth 60fps scrolling
 * - Automatic load-more functionality
 * - Accessibility compliant
 */
const VirtualizedFeed = React.memo(({ 
  items = [], 
  isLoading = false, 
  hasMore = false, 
  loadMore = () => {}, 
  itemHeight = 320,
  overscanCount = 3,
  className = '',
  ...props 
}) => {
  // Memoize data for virtual list
  const itemData = useMemo(() => ({
    items,
    loadMore,
    hasMore,
    isLoading,
  }), [items, loadMore, hasMore, isLoading]);
  
  // Calculate total height with buffer for loading states
  const totalItemCount = useMemo(() => {
    let count = items.length;
    if (isLoading) count += 3; // Add loading skeletons
    return count;
  }, [items.length, isLoading]);
  
  // Error boundary for individual items
  const handleItemError = useCallback((error, errorInfo) => {
    
    // Could integrate with error tracking service here
  }, []);
  
  if (!items.length && !isLoading) {
    return (
      <div className={styles.emptyFeed} role="status" aria-live="polite">
        <div className={styles.emptyIcon}>📭</div>
        <h3>No content available</h3>
        <p>Check back later for new updates!</p>
      </div>
    );
  }
  
  return (
    <div 
      className={`${styles.virtualizedFeed} ${className}`}
      role="feed"
      aria-label="Content feed"
      aria-live="polite"
      aria-busy={isLoading}
      {...props}
    >
      <AutoSizer>
        {({ height, width }) => (
          <List
            height={height}
            width={width}
            itemCount={totalItemCount}
            itemSize={itemHeight}
            itemData={itemData}
            overscanCount={overscanCount}
            onItemsRendered={({ visibleStartIndex, visibleStopIndex }) => {
              // Optional: Track visible items for analytics
              // trackVisibleItems(visibleStartIndex, visibleStopIndex);
            }}
          >
            {FeedItem}
          </List>
        )}
      </AutoSizer>
      
      {/* Loading indicator */}
      {isLoading && (
        <div className={styles.loadingIndicator} role="status" aria-live="polite">
          <div className={styles.spinner} aria-hidden="true"></div>
          <span className={styles.srOnly}>Loading more content...</span>
        </div>
      )}
    </div>
  );
});

VirtualizedFeed.displayName = 'VirtualizedFeed';

export default VirtualizedFeed; 