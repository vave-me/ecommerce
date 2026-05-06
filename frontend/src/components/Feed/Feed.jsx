"use client";

import React, { useEffect, useMemo, lazy, Suspense, memo } from 'react';
import { useInView } from 'react-intersection-observer';
import styles from './Feed.module.css';

// Lazy load card components for better performance
const PostCard = lazy(() => import('../ImprovedPostCard'));
const ClassifiedCard = lazy(() => import('../classified/ClassifiedCard'));
const ServiceCard = lazy(() => import('../services/ServiceCard'));
const TweetCard = lazy(() => import('../../app/[locale]/design/TweetCard'));
const VideoCard = lazy(() => import('../../app/[locale]/design/videos/page'));
const ShortVideo = lazy(() => import('../../app/[locale]/design/shortsSingle/page'));

// Card loading fallback component
const CardSkeleton = memo(() => (
  <div className={styles.cardSkeleton}>
    <div className={styles.skeletonHeader}></div>
    <div className={styles.skeletonContent}></div>
    <div className={styles.skeletonFooter}></div>
  </div>
));

CardSkeleton.displayName = 'CardSkeleton';

/**
 * Feed Component - displays feed items with infinite scroll
 * Uses the unified useFeed hook for data
 * 
 * @param {Object} props
 * @param {string} props.entityType - Optional entity type filter
 */
const Feed = memo(({ entityType = null }) => {
  // Use the single unified feed hook
  const useFeed = require('../../hooks/useFeed').useFeed;
  const data = useFeed({ entityType });
  
  // Extract data from hook
  const feedItems = data?.items || [];
  const isLoading = data?.isLoading || false;
  const error = data?.error || null;
  const hasMore = data?.hasMore || false;
  const loadMore = data?.loadMore || (() => {});
  const isFetchingMore = data?.isFetchingNextPage || false;
  
  // Setup intersection observer for infinite scroll
  const { ref, inView } = useInView({
    threshold: 0,
    rootMargin: '200px'
  });
  
  // Fetch next page when the loading trigger comes into view
  useEffect(() => {
    if (inView && hasMore && !isFetchingMore && loadMore) {
      loadMore();
    }
  }, [inView, hasMore, isFetchingMore, loadMore]);
  
  // Function to render the appropriate card based on item type
  const renderItem = (item) => {
    if (!item) return null;
    
    // Determine item type from various possible fields
    const itemType = item.entityType || item.entity_type || item.contentType || item.type;
    const itemData = item[itemType] || item;
    
    try {
      const CardComponent = (() => {
        switch (itemType) {
          case 'product':
          case 'vehicle':
          case 'property':
          case 'job':
          case 'listing':
            return ClassifiedCard;
          case 'post':
          case 'news':
            return PostCard;
          case 'service':
            return ServiceCard;
          case 'tweet':
            return TweetCard;
          case 'video':
            return VideoCard;
          case 'short':
            return ShortVideo;
          default:
            return null;
        }
      })();
      
      if (!CardComponent) {
        return <div key={item.id} className="feed-item-error">Unknown item type: {itemType}</div>;
      }
      
      return (
        <Suspense key={item.id} fallback={<CardSkeleton />}>
          {itemType === 'product' || ['vehicle', 'property', 'job', 'listing'].includes(itemType) ? (
            <CardComponent product={itemData} />
          ) : itemType === 'post' || itemType === 'news' ? (
            <CardComponent post={itemData} />
          ) : (
            <CardComponent {...{[itemType]: itemData}} />
          )}
        </Suspense>
      );
    } catch (err) {
      // Error: 'Error rendering feed item:', err, item...
      return <div key={item.id} className="feed-item-error">Error rendering item</div>;
    }
  };
  
  // Handle error state
  if (error) {
    return (
      <div className={styles.feedContainer}>
        <div className={styles.errorState}>
          <h3>Error loading feed</h3>
          <p>{error.message || 'Something went wrong'}</p>
          <button onClick={() => window.location.reload()}>
            Retry
          </button>
        </div>
      </div>
    );
  }
  
  // Handle initial loading state
  if (isLoading && !feedItems.length) {
    return (
      <div className={styles.feedContainer}>
        <div className={styles.loadingGrid}>
          {[...Array(6)].map((_, i) => (
            <CardSkeleton key={`skeleton-${i}`} />
          ))}
        </div>
      </div>
    );
  }
  
  // Handle empty state
  if (!isLoading && !feedItems.length) {
    return (
      <div className={styles.feedContainer}>
        <div className={styles.emptyState}>
          <h3>No items found</h3>
          <p>Try adjusting your filters or check back later</p>
        </div>
      </div>
    );
  }
  
  return (
    <div className={styles.feedContainer}>
      <div className={styles.feedGrid}>
        {feedItems.map(renderItem)}
      </div>
      
      {/* Infinite scroll trigger */}
      {hasMore && (
        <div ref={ref} className={styles.loadMoreTrigger}>
          {isFetchingMore && (
            <div className={styles.loadingMore}>
              <CardSkeleton />
              <CardSkeleton />
              <CardSkeleton />
            </div>
          )}
        </div>
      )}
      
      {/* End of feed indicator */}
      {!hasMore && feedItems.length > 0 && (
        <div className={styles.endOfFeed}>
          You've reached the end of the feed
        </div>
      )}
    </div>
  );
});

Feed.displayName = 'Feed';

export default Feed;