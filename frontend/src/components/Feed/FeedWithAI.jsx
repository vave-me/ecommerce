"use client";

import React, { useEffect, useMemo, lazy, Suspense, memo } from 'react';
import { useInView } from 'react-intersection-observer';
import { Bot, Sparkles } from '@/icons';
import styles from './Feed.module.css';
import aiStyles from './FeedWithAI.module.css';

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
 * Enhanced Feed Component with AI Results Support
 * Can display regular feed items or AI-sourced results
 * 
 * @param {Object} props
 * @param {string} props.entityType - Optional entity type filter
 * @param {Object} props.aiResults - Optional AI results to display
 * @param {boolean} props.showAIIndicators - Whether to show AI source badges
 * @param {Function} props.onItemClick - Callback when item is clicked
 */
const FeedWithAI = memo(({ 
    entityType = null, 
    aiResults = null,
    showAIIndicators = true,
    onItemClick = null,
    className = ''
}) => {
    // Use the single unified feed hook if no AI results
    const useFeed = !aiResults ? require('../../hooks/useFeed').useFeed : null;
    const feedData = !aiResults ? useFeed?.({ entityType }) : null;
    
    // Determine data source
    const isAIMode = !!aiResults;
    const feedItems = isAIMode 
        ? (aiResults.items || aiResults.Data?.items || [])
        : (feedData?.items || []);
    
    const isLoading = isAIMode ? false : (feedData?.isLoading || false);
    const error = isAIMode ? null : (feedData?.error || null);
    const hasMore = isAIMode ? false : (feedData?.hasMore || false);
    const loadMore = isAIMode ? null : (feedData?.loadMore || (() => {}));
    const isFetchingMore = isAIMode ? false : (feedData?.isFetchingNextPage || false);
    
    // Setup intersection observer for infinite scroll (only for regular feed)
    const { ref, inView } = useInView({
        threshold: 0,
        rootMargin: '200px',
        skip: isAIMode
    });
    
    // Fetch next page when the loading trigger comes into view
    useEffect(() => {
        if (!isAIMode && inView && hasMore && !isFetchingMore && loadMore) {
            loadMore();
        }
    }, [inView, hasMore, isFetchingMore, loadMore, isAIMode]);
    
    // Function to render the appropriate card based on item type
    const renderItem = (item, index) => {
        if (!item) return null;
        
        // Determine item type from various possible fields
        const itemType = item.entityType || item.entity_type || item.contentType || item.type;
        const itemData = item[itemType] || item;
        const isAISourced = item.source === 'assistant' || item.source === 'ai';
        
        // Wrapper for AI-sourced items
        const AIWrapper = ({ children }) => {
            if (!isAISourced || !showAIIndicators) return children;
            
            return (
                <div className={aiStyles.aiItemWrapper}>
                    <div className={aiStyles.aiIndicator}>
                        <Bot size={14} />
                        <span>AI Match</span>
                        {item.relevanceScore && (
                            <span className={aiStyles.relevanceScore}>
                                {Math.round(item.relevanceScore * 100)}%
                            </span>
                        )}
                    </div>
                    {children}
                </div>
            );
        };
        
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
            
            const handleClick = () => {
                if (onItemClick) {
                    onItemClick(item);
                }
            };
            
            return (
                <AIWrapper key={item.id || `item-${index}`}>
                    <Suspense fallback={<CardSkeleton />}>
                        {itemType === 'product' || ['vehicle', 'property', 'job', 'listing'].includes(itemType) ? (
                            <CardComponent product={itemData} onClick={handleClick} />
                        ) : itemType === 'post' || itemType === 'news' ? (
                            <PostCard post={itemData} onClick={handleClick} />
                        ) : itemType === 'service' ? (
                            <ServiceCard service={itemData} onClick={handleClick} />
                        ) : itemType === 'tweet' ? (
                            <TweetCard tweet={itemData} onClick={handleClick} />
                        ) : itemType === 'video' ? (
                            <VideoCard video={itemData} onClick={handleClick} />
                        ) : itemType === 'short' ? (
                            <ShortVideo short={itemData} onClick={handleClick} />
                        ) : null}
                    </Suspense>
                </AIWrapper>
            );
        } catch (error) {
            // Error: 'Error rendering feed item:', error...
            return (
                <div key={item.id || `error-${index}`} className="feed-item-error">
                    Error rendering item
                </div>
            );
        }
    };
    
    // Error state
    if (error) {
        return (
            <div className={`${styles.errorContainer} ${className}`}>
                <h3>Error loading feed</h3>
                <p>{error.message || 'Something went wrong'}</p>
            </div>
        );
    }
    
    // Empty state
    if (!isLoading && feedItems.length === 0) {
        return (
            <div className={`${styles.emptyContainer} ${className}`}>
                {isAIMode ? (
                    <>
                        <Sparkles size={48} className={styles.emptyIcon} />
                        <h3>No results found</h3>
                        <p>Try refining your search or asking a different question</p>
                    </>
                ) : (
                    <>
                        <h3>No items to display</h3>
                        <p>Check back later for new content</p>
                    </>
                )}
            </div>
        );
    }
    
    // AI Results Header
    const renderAIHeader = () => {
        if (!isAIMode || !aiResults.Data?.metadata) return null;
        
        const metadata = aiResults.Data.metadata;
        return (
            <div className={aiStyles.aiResultsHeader}>
                <div className={aiStyles.aiResultsInfo}>
                    <Sparkles size={20} />
                    <h3>AI Search Results</h3>
                    {metadata.totalCount && (
                        <span className={aiStyles.resultCount}>
                            {metadata.totalCount} results
                        </span>
                    )}
                </div>
                {metadata.queryInterpretation && (
                    <p className={aiStyles.queryInterpretation}>
                        Showing results for: "{metadata.queryInterpretation}"
                    </p>
                )}
            </div>
        );
    };
    
    return (
        <div className={`${styles.feedContainer} ${className} ${isAIMode ? aiStyles.aiFeedContainer : ''}`}>
            {renderAIHeader()}
            
            <div className={styles.feedGrid}>
                {feedItems.map((item, index) => renderItem(item, index))}
                
                {/* Loading skeleton for initial load */}
                {isLoading && !feedItems.length && (
                    <>
                        {[...Array(6)].map((_, i) => (
                            <CardSkeleton key={`skeleton-${i}`} />
                        ))}
                    </>
                )}
                
                {/* Infinite scroll trigger (only for regular feed) */}
                {!isAIMode && hasMore && (
                    <div ref={ref} className={styles.loadMoreTrigger}>
                        {isFetchingMore && (
                            <div className={styles.loadingMore}>
                                <div className={styles.spinner} />
                                <span>Loading more...</span>
                            </div>
                        )}
                    </div>
                )}
            </div>
            
            {/* AI Context Footer */}
            {isAIMode && aiResults.Response && (
                <div className={aiStyles.aiContextFooter}>
                    <div className={aiStyles.aiContextHeader}>
                        <Bot size={18} />
                        <span>Assistant's Summary</span>
                    </div>
                    <p>{aiResults.Response}</p>
                </div>
            )}
        </div>
    );
});

FeedWithAI.displayName = 'FeedWithAI';

export default FeedWithAI;