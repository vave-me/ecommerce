"use client";

import React, { memo, useEffect, Suspense } from 'react';
import { useInView } from 'react-intersection-observer';
import { useFeed } from '../../hooks/useFeed';
import FeedItem from './FeedItem.client';
import FeedItemServer from './FeedItemServer';
import ClassifiedCard from '../classified/ClassifiedCard';
import PostCard from '../PostCard/PostCard';
import AIResponseCard from './AIResponseCard';
import FeedErrorState from './FeedErrorState';
import FeedEmptyState from './FeedEmptyState';
import FeedLoadingState from './FeedLoadingState';
import { useDispatch, useSelector } from 'react-redux';
import { resetFilters } from '../../redux/slices/listingFiltersSlice';
import styles from '../../app/page.module.css';

/**
 * Unified FeedDisplay Component
 * Combines functionality from FeedDisplay.client.jsx, EnhancedFeedDisplay.jsx, and OptimizedFeedDisplay.jsx
 * 
 * Features:
 * - Server-side rendering support with initial items
 * - Client-side infinite scroll
 * - AI response integration
 * - Filter management
 * - Optimized rendering with memoization
 * 
 * @param {Object} props
 * @param {Array} props.initialItems - Initial items for server-side rendering
 * @param {number} props.skipInitialItems - Number of items to skip (for hybrid rendering)
 * @param {boolean} props.showAIResponses - Whether to show AI responses in feed
 * @param {boolean} props.enableFilters - Whether to enable filter controls
 * @param {Object} props.feedParams - Additional parameters for feed queries
 */
const FeedDisplay = memo(function FeedDisplay({ 
  initialItems = [],
  skipInitialItems = 0,
  showAIResponses = false,
  enableFilters = true,
  feedParams = {}
}) {
  const { items, isLoading, hasMore, loadMore, error, refetch } = useFeed({ 
    entityType: feedParams?.entityType 
  });
  const dispatch = useDispatch();
  const filters = useSelector(state => state.listingFilters);
  
  const { ref: loadMoreRef, inView } = useInView({
    threshold: 0.1,
    rootMargin: '200px',
    triggerOnce: false
  });
  
  // Check if any filters are active
  const hasActiveFilters = React.useMemo(() => {
    if (!filters || !enableFilters) return false;
    return Object.entries(filters).some(([key, value]) => {
      if (key === 'page' || key === 'limit') return false;
      if (Array.isArray(value)) return value.length > 0;
      return value !== null && value !== undefined && value !== '';
    });
  }, [filters, enableFilters]);
  
  // Items to render (no AI response merging - keep it simple)
  const feedItems = items;
  
  // Trigger infinite scroll
  useEffect(() => {
    if (inView && hasMore && !isLoading) {
      loadMore();
    }
  }, [inView, hasMore, isLoading, loadMore]);
  
  // Determine which items to render
  const itemsToRender = React.useMemo(() => {
    return skipInitialItems > 0 ? feedItems.slice(skipInitialItems) : feedItems;
  }, [feedItems, skipInitialItems]);
  
  // Render appropriate card based on item type
  const renderFeedItem = (item, index) => {
    // Remove AI response handling - keep it simple
    
    const entityType = item.entityType || item.entity_type || item.typeOfPost || item.type;
    const entityId = item[entityType]?.id || item.id;
    const adjustedIndex = index + skipInitialItems;
    const uniqueKey = entityId
      ? `client-${entityType}-${entityId}-${adjustedIndex}`
      : `client-${entityType}-${item.createdAt}-${adjustedIndex}`;
    
    // Product/Classified Card
    if (['product', 'vehicle', 'property', 'job', 'listing'].includes(entityType)) {
      const productData = item[entityType] || item;
      return (
        <div key={uniqueKey} className={styles.feedItemContainer}>
          <ClassifiedCard product={productData} />
        </div>
      );
    }
    
    // Post Card
    if (entityType === 'post' || entityType === 'news') {
      const postData = item[entityType] || item;
      return (
        <div key={uniqueKey} className={styles.feedItemContainer}>
          <PostCard post={postData} />
        </div>
      );
    }
    
    // Default Feed Item
    return (
      <div key={uniqueKey} className={styles.feedItemContainer}>
        <FeedItem item={item} />
      </div>
    );
  };
  
  const handleClearFilters = React.useCallback(() => {
    dispatch(resetFilters());
  }, [dispatch]);
  
  const handleRefresh = React.useCallback(() => {
    refetch();
  }, [refetch]);
  
  // Error state
  if (error) {
    return <FeedErrorState error={error} onRetry={handleRefresh} />;
  }
  
  // Empty state
  if (!isLoading && feedItems.length === 0 && !hasMore) {
    return (
      <FeedEmptyState 
        hasFilters={hasActiveFilters}
        onClearFilters={hasActiveFilters && enableFilters ? handleClearFilters : null}
        onRefresh={handleRefresh}
      />
    );
  }
  
  // Initial loading state
  if (isLoading && feedItems.length === 0) {
    return <FeedLoadingState message="Loading feed..." />;
  }
  
  return (
    <section className={styles.feedSection}>
      {/* Server-rendered initial items */}
      {initialItems.length > 0 && (
        <div className="initial-feed-items">
          {initialItems.map((item, index) => {
            const entityId = item[item.entityType]?.id || item.id;
            const uniqueKey = entityId
              ? `server-${item.entityType}-${entityId}-${index}`
              : `server-${item.entityType}-${item.createdAt}-${index}`;
            
            return (
              <Suspense 
                key={uniqueKey} 
                fallback={
                  <div className="feed-item-loading">
                    <div className="loading-skeleton"></div>
                  </div>
                }
              >
                <FeedItemServer item={item} />
              </Suspense>
            );
          })}
        </div>
      )}
      
      {/* Client-rendered items with infinite scroll */}
      <div className={styles.initialFeedItems}>
        {itemsToRender.map(renderFeedItem)}
      </div>
      
      {/* Loading indicator */}
      {isLoading && feedItems.length > 0 && (
        <FeedLoadingState message="Loading more items..." />
      )}
      
      {/* Infinite scroll trigger */}
      <div ref={loadMoreRef} style={{ height: '20px' }} />
      
      {/* End of feed indicator */}
      {!hasMore && feedItems.length > 0 && (
        <div className={styles.endOfFeed}>
          You've reached the end of the feed.
        </div>
      )}
    </section>
  );
}, (prevProps, nextProps) => {
  // Custom comparison for performance optimization
  const itemsEqual = prevProps.initialItems.length === nextProps.initialItems.length &&
    prevProps.initialItems.every((item, index) => {
      const nextItem = nextProps.initialItems[index];
      return item?.id === nextItem?.id && 
             item?.entityType === nextItem?.entityType &&
             item?.createdAt === nextItem?.createdAt;
    });
  
  return (
    itemsEqual &&
    prevProps.skipInitialItems === nextProps.skipInitialItems &&
    prevProps.showAIResponses === nextProps.showAIResponses &&
    prevProps.enableFilters === nextProps.enableFilters &&
    JSON.stringify(prevProps.feedParams) === JSON.stringify(nextProps.feedParams)
  );
});

export default FeedDisplay;