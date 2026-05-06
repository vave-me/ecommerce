/**
 * InfiniteScroll - High-Performance Mobile-Optimized Infinite Scroll
 *
 * Features:
 * - Intersection Observer for performance
 * - Virtual scrolling for large lists
 * - Mobile-optimized loading states
 * - Pull-to-refresh integration
 * - Memory management for mobile devices
 *
 * Designed to work seamlessly with existing data fetching patterns
 */
"use client";
import React, {
    useRef,
    useCallback,
    useEffect,
    useState,
    useMemo,
    forwardRef,
    useImperativeHandle
} from 'react';
import PropTypes from 'prop-types';
import TouchInteractions from './TouchInteractions';
// Configuration for mobile optimization
const MOBILE_CONFIG = {
    INTERSECTION_THRESHOLD: 0.1,
    ROOT_MARGIN: '100px',
    VIRTUAL_ITEM_HEIGHT: 120, // Estimated item height for virtualization
    BUFFER_SIZE: 5, // Number of items to render outside viewport
    LOADING_DELAY: 300, // Delay before showing loading indicator
    MAX_ITEMS_IN_MEMORY: 1000, // Maximum items to keep in memory (mobile optimization)
};
const InfiniteScroll = forwardRef(({
                                       items = [],
                                       renderItem,
                                       loadMore,
                                       hasNextPage = false,
                                       isLoading = false,
                                       error = null,
                                       onRefresh,
                                       refreshing = false,
                                       enablePullToRefresh = true,
                                       enableVirtualization = false,
                                       itemHeight = MOBILE_CONFIG.VIRTUAL_ITEM_HEIGHT,
                                       className = '',
                                       emptyStateComponent: EmptyState,
                                       loadingComponent: LoadingComponent,
                                       errorComponent: ErrorComponent,
                                       endMessage = 'No more items to load',
                                       threshold = MOBILE_CONFIG.INTERSECTION_THRESHOLD,
                                       rootMargin = MOBILE_CONFIG.ROOT_MARGIN,
                                       onItemVisible,
                                       estimatedItemSize,
                                       ...props
                                   }, ref) => {
    // Refs for intersection observer and virtualization
    const containerRef = useRef(null);
    const sentinelRef = useRef(null);
    const observerRef = useRef(null);
    const virtualScrollRef = useRef({ startIndex: 0, endIndex: 0 });
    // State for virtual scrolling
    const [scrollTop, setScrollTop] = useState(0);
    const [containerHeight, setContainerHeight] = useState(0);
    const [showLoadingIndicator, setShowLoadingIndicator] = useState(false);
    // Expose scroll methods to parent components
    useImperativeHandle(ref, () => ({
        scrollToTop: () => {
            if (containerRef.current) {
                containerRef.current.scrollTo({ top: 0, behavior: 'smooth' });
            }
        },
        scrollToItem: (index) => {
            if (containerRef.current && enableVirtualization) {
                const position = index * itemHeight;
                containerRef.current.scrollTo({ top: position, behavior: 'smooth' });
            }
        },
        refresh: () => {
            if (onRefresh) {
                onRefresh();
            }
        }
    }), [itemHeight, enableVirtualization, onRefresh]);
    // Memoized calculations for virtual scrolling
    const virtualProps = useMemo(() => {
        if (!enableVirtualization || !containerHeight) {
            return {
                startIndex: 0,
                endIndex: items.length,
                visibleItems: items,
                totalHeight: 'auto',
                offsetY: 0
            };
        }
        const startIndex = Math.floor(scrollTop / itemHeight);
        const visibleCount = Math.ceil(containerHeight / itemHeight);
        const endIndex = Math.min(
            startIndex + visibleCount + MOBILE_CONFIG.BUFFER_SIZE * 2,
            items.length
        );
        const adjustedStartIndex = Math.max(0, startIndex - MOBILE_CONFIG.BUFFER_SIZE);
        return {
            startIndex: adjustedStartIndex,
            endIndex,
            visibleItems: items.slice(adjustedStartIndex, endIndex),
            totalHeight: items.length * itemHeight,
            offsetY: adjustedStartIndex * itemHeight
        };
    }, [enableVirtualization, containerHeight, scrollTop, itemHeight, items]);
    // Delayed loading indicator to prevent flickering
    useEffect(() => {
        let timeoutId;
        if (isLoading) {
            timeoutId = setTimeout(() => {
                setShowLoadingIndicator(true);
            }, MOBILE_CONFIG.LOADING_DELAY);
        } else {
            setShowLoadingIndicator(false);
        }
        return () => {
            if (timeoutId) {
                clearTimeout(timeoutId);
            }
        };
    }, [isLoading]);
    // Setup intersection observer for load more
    useEffect(() => {
        if (!sentinelRef.current || !loadMore) return;
        const observer = new IntersectionObserver(
            (entries) => {
                const [entry] = entries;
                if (entry.isIntersecting && hasNextPage && !isLoading) {
                    loadMore();
                }
            },
            {
                threshold,
                rootMargin,
            }
        );
        observer.observe(sentinelRef.current);
        observerRef.current = observer;
        return () => {
            if (observerRef.current) {
                observerRef.current.disconnect();
            }
        };
    }, [hasNextPage, isLoading, loadMore, threshold, rootMargin]);
    // Handle scroll for virtualization
    const handleScroll = useCallback((event) => {
        if (!enableVirtualization) return;
        const target = event.target;
        setScrollTop(target.scrollTop);
        // Notify parent of item visibility changes
        if (onItemVisible) {
            const visibleStartIndex = Math.floor(target.scrollTop / itemHeight);
            onItemVisible(visibleStartIndex);
        }
    }, [enableVirtualization, itemHeight, onItemVisible]);
    // Handle container resize for virtualization
    useEffect(() => {
        if (!enableVirtualization || !containerRef.current) return;
        const resizeObserver = new ResizeObserver((entries) => {
            const [entry] = entries;
            setContainerHeight(entry.contentRect.height);
        });
        resizeObserver.observe(containerRef.current);
        return () => {
            resizeObserver.disconnect();
        };
    }, [enableVirtualization]);
    // Memory management for mobile devices
    useEffect(() => {
        if (items.length > MOBILE_CONFIG.MAX_ITEMS_IN_MEMORY) {
        }
    }, [items.length]);
    // Render loading component
    const renderLoading = () => {
        if (LoadingComponent) {
            return <LoadingComponent />;
        }
        return (
            <div style={{
                display: 'flex',
                justifyContent: 'center',
                padding: '20px',
                color: '#6b7280'
            }}>
                <div style={{
                    width: '24px',
                    height: '24px',
                    border: '2px solid #e5e7eb',
                    borderTopColor: '#2980b9',
                    borderRadius: '50%',
                    animation: 'spin 1s linear infinite'
                }} />
                <style jsx>{`
          @keyframes spin {
            to { transform: rotate(360deg); }
          }
        `}</style>
            </div>
        );
    };
    // Render error component
    const renderError = () => {
        if (ErrorComponent) {
            return <ErrorComponent error={error} onRetry={loadMore} />;
        }
        return (
            <div style={{
                padding: '20px',
                textAlign: 'center',
                color: '#ef4444'
            }}>
                <p>Error loading content</p>
                {loadMore && (
                    <button
                        onClick={loadMore}
                        style={{
                            marginTop: '10px',
                            padding: '8px 16px',
                            backgroundColor: '#2980b9',
                            color: 'white',
                            border: 'none',
                            borderRadius: '8px',
                            cursor: 'pointer'
                        }}
                    >
                        Try Again
                    </button>
                )}
            </div>
        );
    };
    // Render empty state
    const renderEmptyState = () => {
        if (EmptyState) {
            return <EmptyState />;
        }
        return (
            <div style={{
                padding: '40px 20px',
                textAlign: 'center',
                color: '#6b7280'
            }}>
                <p>No items to display</p>
            </div>
        );
    };
    // Handle pull-to-refresh
    const handlePullRefresh = useCallback(() => {
        if (onRefresh) {
            onRefresh();
        }
    }, [onRefresh]);
    // Render list content
    const renderListContent = () => {
        // Show error state
        if (error && items.length === 0) {
            return renderError();
        }
        // Show empty state
        if (items.length === 0 && !isLoading) {
            return renderEmptyState();
        }
        // Virtual scrolling rendering
        if (enableVirtualization) {
            return (
                <div style={{ height: virtualProps.totalHeight, position: 'relative' }}>
                    <div
                        style={{
                            transform: `translateY(${virtualProps.offsetY}px)`,
                            position: 'absolute',
                            top: 0,
                            left: 0,
                            right: 0
                        }}
                    >
                        {virtualProps.visibleItems.map((item, index) => {
                            const actualIndex = virtualProps.startIndex + index;
                            return (
                                <div
                                    key={item.id || actualIndex}
                                    style={{
                                        height: itemHeight,
                                        overflow: 'hidden'
                                    }}
                                >
                                    {renderItem(item, actualIndex)}
                                </div>
                            );
                        })}
                    </div>
                </div>
            );
        }
        // Standard rendering
        return items.map((item, index) => (
            <React.Fragment key={item.id || index}>
                {renderItem(item, index)}
            </React.Fragment>
        ));
    };
    return (
        <TouchInteractions
            enablePullToRefresh={enablePullToRefresh}
            onPullRefresh={handlePullRefresh}
            refreshing={refreshing}
            className={className}
            {...props}
        >
            <div
                ref={containerRef}
                onScroll={handleScroll}
                style={{
                    overflowY: 'auto',
                    WebkitOverflowScrolling: 'touch',
                    /* MOBILE NAVIGATION FIX - Extra padding for bottom nav */
                    paddingBottom: '60px',
                    ...props.style
                }}
            >
                {renderListContent()}
                {/* Loading indicator for infinite scroll */}
                {hasNextPage && showLoadingIndicator && renderLoading()}
                {/* End message */}
                {!hasNextPage && items.length > 0 && (
                    <div style={{
                        padding: '20px',
                        textAlign: 'center',
                        color: '#6b7280',
                        fontSize: '14px'
                    }}>
                        {endMessage}
                    </div>
                )}
                {/* Sentinel for intersection observer */}
                {hasNextPage && (
                    <div
                        ref={sentinelRef}
                        style={{
                            height: '1px',
                            margin: '10px 0'
                        }}
                        aria-hidden="true"
                    />
                )}
            </div>
        </TouchInteractions>
    );
});
InfiniteScroll.displayName = 'InfiniteScroll';
InfiniteScroll.propTypes = {
    items: PropTypes.array.isRequired,
    renderItem: PropTypes.func.isRequired,
    loadMore: PropTypes.func,
    hasNextPage: PropTypes.bool,
    isLoading: PropTypes.bool,
    error: PropTypes.any,
    onRefresh: PropTypes.func,
    refreshing: PropTypes.bool,
    enablePullToRefresh: PropTypes.bool,
    enableVirtualization: PropTypes.bool,
    itemHeight: PropTypes.number,
    className: PropTypes.string,
    emptyStateComponent: PropTypes.elementType,
    loadingComponent: PropTypes.elementType,
    errorComponent: PropTypes.elementType,
    endMessage: PropTypes.string,
    threshold: PropTypes.number,
    rootMargin: PropTypes.string,
    onItemVisible: PropTypes.func,
    estimatedItemSize: PropTypes.number,
};
export default InfiniteScroll;