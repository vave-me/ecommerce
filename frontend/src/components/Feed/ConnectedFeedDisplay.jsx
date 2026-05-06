"use client";

import React, { useEffect, useRef, memo } from 'react';
import { useInView } from 'react-intersection-observer';
import { useFeed } from '../../hooks/useFeed';
import FeedItem from './FeedItem.client';
import PostCard from '../PostCard/PostCard';
import { ClassifiedCard } from '../classified';
import AIResponseCard from './AIResponseCard'
import ExpandableVaveAds from '../ads/ExpandableVaveAds';
import styles from '../../app/page.module.css';
import { feedLog } from '../../utils/productionConsoleCleanup';
import { useSelector } from 'react-redux';
import { selectShowUnifiedComposer } from '../../redux/slices/uiPreferencesSlice';
import { useIsMobile } from '../../hooks/useMobileDetection';

/**
 * Infinite scroll trigger component
 */
const InfiniteScrollTrigger = memo(function InfiniteScrollTrigger({ onLoadMore, hasMore, isLoading }) {
  const triggerRef = useRef(null);
  
  useEffect(() => {
    const trigger = triggerRef.current;
    if (!trigger || !onLoadMore || !hasMore || isLoading) return;
    
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting) {
          if (process.env.NODE_ENV === 'development') {
            feedLog.scroll('🎯 InfiniteScrollTrigger activated');
          }
          onLoadMore();
        }
      },
      {
        threshold: 0.1,
        rootMargin: '200px'
      }
    );
    
    observer.observe(trigger);
    
    return () => {
      observer.disconnect();
    };
  }, [onLoadMore, hasMore, isLoading]);
  
  return (
    <div
      ref={triggerRef}
      className="load-more-trigger"
      style={{
        height: '20px',
        margin: '20px 0',
        backgroundColor: process.env.NODE_ENV === 'development' ? '#ff000020' : 'transparent',
        border: process.env.NODE_ENV === 'development' ? '2px dashed #ff0000' : 'none',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        fontSize: '12px',
        color: process.env.NODE_ENV === 'development' ? '#ff0000' : 'transparent'
      }}
      aria-label="Loading more content"
    >
      {process.env.NODE_ENV === 'development' && '🎯 INFINITE SCROLL TRIGGER'}
    </div>
  );
});

/**
 * Unified ConnectedFeedDisplay Component
 * Merges functionality from ConnectedFeedDisplay.jsx and ConnectedFeedDisplay.enhanced.jsx
 * 
 * Features:
 * - Feed provider integration
 * - AI response support
 * - Mobile responsive with composer button
 * - Infinite scroll
 * - Debug logging (development only)
 * 
 * @param {Object} props
 * @param {boolean} props.showComposer - Whether to show the composer
 * @param {boolean} props.showAIResponses - Whether to integrate AI responses
 * @param {boolean} props.useEnhancedComposer - Use enhanced composer with AI features
 * @param {string} props.composerPlaceholder - Placeholder text for composer
 * @param {Function} props.onPostCreate - Callback when post is created
 */
const ConnectedFeedDisplay = memo(function ConnectedFeedDisplay({ 
  showComposer = true,
  showAIResponses = false,
  useEnhancedComposer = false,
  composerPlaceholder = "What are you looking for today?",
  onPostCreate
}) {
  const { items: feedItems, isLoading, hasMore, loadMore, error } = useFeed();
  const showUnifiedComposer = useSelector(selectShowUnifiedComposer);
  const isMobile = useIsMobile();
  
  const { ref: loadMoreRef, inView } = useInView({
    threshold: 0.1,
    rootMargin: '100px',
  });
  
  // Use feed items directly (no AI merging)
  const displayItems = feedItems;
  
  // Trigger infinite scroll with ref
  useEffect(() => {
    if (inView && !isLoading && hasMore) {
      loadMore();
    }
  }, [inView, isLoading, hasMore, loadMore]);
  
  // Debug logging (development only)
  useEffect(() => {
    if (process.env.NODE_ENV === 'development') {
      feedLog.update('Feed state update:', {
        feedItemsCount: feedItems?.length || 0,
        isLoading,
        hasMore,
        loadMore: typeof loadMore,
        showAIResponses,
        displayItemsCount: displayItems?.length || 0
      });
    }
  }, [feedItems, isLoading, hasMore, loadMore, showAIResponses, displayItems]);
  
  const renderItem = (item, index) => {
    // AI Response
    if (showAIResponses && item.type === 'ai-response') {
      return (
        <AIResponseCard
          key={item.id}
          response={item.response}
          query={item.query}
          timestamp={item.timestamp}
          metadata={item.metadata}
        />
      );
    }
    
    const entityType = item.entityType || item.entity_type || item.typeOfPost || 'unknown';
    const entityId = item[entityType]?.id || item.id;
    const uniqueKey = entityId
      ? `feed-${entityType}-${entityId}-${index}`
      : `feed-${entityType}-${item.createdAt}-${index}`;
    
    // Route to appropriate component based on entity type
    switch (entityType) {
      case 'product':
      case 'vehicle':
      case 'property':
      case 'job':
      case 'listing':
        const productData = item[entityType] || item;
        return (
          <div key={uniqueKey} className={styles.feedItemContainer}>
            <ClassifiedCard product={productData} />
          </div>
        );
      
      case 'post':
      case 'news':
        const postData = item[entityType] || item;
        return (
          <div key={uniqueKey} className={styles.feedItemContainer}>
            <PostCard post={postData} />
          </div>
        );
      
      default:
        return (
          <div key={uniqueKey} className={styles.feedItemContainer}>
            <FeedItem item={item} />
          </div>
        );
    }
  };
  
  // Configure composer props based on enhanced mode
  const composerProps = {
    placeholder: composerPlaceholder,
    onPostCreate,
    showTemplates: useEnhancedComposer,
    enableChat: useEnhancedComposer,
    autoFocusOnAI: false
  };
  
  if (error) {
    return (
      <div className={styles.errorFallback}>
        <p>Error loading feed: {error.message}</p>
        <button onClick={() => window.location.reload()}>
          Retry
        </button>
      </div>
    );
  }
  
  if (!isLoading && displayItems.length === 0) {
    return (
      <>
        {/* Show only ads when feed is empty */}
        {showComposer && !isMobile && showUnifiedComposer && (
          <ExpandableVaveAds />
        )}

        <div className={styles.emptyFeed}>
          <h3>No items found</h3>
          <p>Try adjusting your filters or search criteria</p>
        </div>
      </>
    );
  }
  
  return (
    <div className={styles.feedSection}>
      {/* Desktop: Show only ads */}
      {showComposer && !isMobile && showUnifiedComposer && (
        <ExpandableVaveAds />
      )}

      {/* Feed items container */}
      <div className={styles.initialFeedItems}>
        {displayItems.map((item, index) => renderItem(item, index))}
      </div>
      
      {/* Loading state */}
      {isLoading && (
        <div className={styles.feedLoading}>
          <div className={styles.spinner}></div>
          <p>Loading{displayItems.length > 0 ? ' more items' : ' feed'}...</p>
        </div>
      )}
      
      {/* End of feed indicator */}
      {!hasMore && displayItems.length > 0 && (
        <div className={styles.endOfFeed}>
          You've reached the end of the feed.
        </div>
      )}
      
      {/* Development mode: Manual load more button */}
      {process.env.NODE_ENV === 'development' && hasMore && !isLoading && (
        <div style={{ textAlign: 'center', margin: '20px 0' }}>
          <button
            onClick={loadMore}
            style={{
              padding: '10px 20px',
              backgroundColor: '#2980b9',
              color: 'white',
              border: 'none',
              borderRadius: '8px',
              cursor: 'pointer',
              marginBottom: '10px'
            }}
          >
            🔄 Manual Load More
          </button>
          <div style={{ fontSize: '12px', color: '#666', marginTop: '5px' }}>
            hasMore: {hasMore ? 'true' : 'false'} | Items: {displayItems.length}
          </div>
        </div>
      )}
      
      {/* Infinite scroll trigger */}
      {hasMore && !isLoading && <InfiniteScrollTrigger onLoadMore={loadMore} hasMore={hasMore} isLoading={isLoading} />}
      
      {/* Alternative infinite scroll trigger with ref */}
      <div ref={loadMoreRef} style={{ height: '20px' }} />
    </div>
  );
});

export default ConnectedFeedDisplay;