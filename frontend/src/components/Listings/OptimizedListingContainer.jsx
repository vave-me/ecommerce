"use client";
import React, {
    memo,
    useCallback,
    useMemo,
    useState,
    useRef
} from 'react';
import {useQuery, useInfiniteQuery} from '@tanstack/react-query';
import {useEntityFilters} from '../../hooks/useEntityFilters';
import {useCategoryData} from '../../hooks/useCategoryData';
import {useDispatch} from 'react-redux';
import {setFilters} from '../../redux/slices/listingFiltersSlice';
import {useAuthUser} from '../../context/OptimizedAuthContext';
import {fetchListings} from '../../api/searchApi';
import LoadingSpinner from '../Utils/LoadingSpinner';
import InfiniteScroll from '../Shared/InfiniteScroll';
/**
 * MemoizedListItem - Prevents unnecessary re-renders of list items
 */
const MemoizedListItem = memo(({item, onSelect}) => {
    // Use callback to prevent recreation on each render
    const handleClick = useCallback(() => {
        onSelect(item.id);
    }, [item.id, onSelect]);
    return (
        <div className="listing-item" onClick={handleClick}>
            <h3>{item.title}</h3>
            <p>{item.description.substring(0, 100)}...</p>
            <div className="item-footer">
                <span className="price">${item.price}</span>
                <span className="location">{item.location}</span>
            </div>
        </div>
    );
});
MemoizedListItem.displayName = 'MemoizedListItem';
/**
 * OptimizedListingContainer - Enhanced with mobile-optimized infinite scroll
 * Maintains full desktop compatibility with pagination fallback
 */
const OptimizedListingContainer = memo(({entityType = 'products'}) => {
    const dispatch = useDispatch();
    const user = useAuthUser();
    const infiniteScrollRef = useRef(null);
    // Get filter data with memoized selectors
    const {filters, getQueryKey} = useEntityFilters(entityType);
    // Get category data without context re-renders
    const {categoryData} = useCategoryData([entityType]);
    const categories = useMemo(() =>
            categoryData[entityType]?.categories || [],
        [categoryData, entityType]
    );
    // Determine if we should use infinite scroll (mobile-first approach)
    const [useInfiniteScroll, setUseInfiniteScroll] = useState(() => {
        if (typeof window === 'undefined') return false;
        return window.innerWidth <= 768 || 'ontouchstart' in window;
    });
    // Enhanced fetch function for infinite scroll
    const fetchListingsPage = useCallback(async ({pageParam = 1}) => {
        const result = await fetchListings({
            entityType,
            filters: {...filters, page: pageParam},
            userId: user?.userId
        });
        return {
            ...result,
            currentPage: pageParam,
            nextPage: result.pagination?.hasNext ? pageParam + 1 : null,
        };
    }, [entityType, filters, user?.userId]);
    // Infinite query for mobile/touch devices
    const infiniteQueryResult = useInfiniteQuery({
        queryKey: [...getQueryKey(), 'infinite'],
        queryFn: fetchListingsPage,
        getNextPageParam: (lastPage) => lastPage.nextPage,
        keepPreviousData: true,
        enabled: useInfiniteScroll,
    });
    // Standard query for desktop fallback
    const standardQueryResult = useQuery({
        queryKey: getQueryKey(),
        queryFn: () => fetchListings({
            entityType,
            filters,
            userId: user?.userId
        }),
        keepPreviousData: true,
        enabled: !useInfiniteScroll,
    });
    // Use appropriate query result
    const queryResult = useInfiniteScroll ? infiniteQueryResult : standardQueryResult;
    // Flatten infinite scroll data or use standard data
    const listings = useMemo(() => {
        if (useInfiniteScroll && infiniteQueryResult.data) {
            return infiniteQueryResult.data.pages.flatMap(page => page.listings || []);
        }
        return standardQueryResult.data?.listings || [];
    }, [useInfiniteScroll, infiniteQueryResult.data, standardQueryResult.data]);
    // Pagination data for desktop
    const paginationData = useMemo(() => {
        if (useInfiniteScroll) return null;
        return standardQueryResult.data?.pagination;
    }, [useInfiniteScroll, standardQueryResult.data]);
    // Memoize item selection handler
    const handleSelectItem = useCallback((itemId) => {
        // Implementation would dispatch an action or navigate
    }, []);
    // Handle load more for infinite scroll
    const handleLoadMore = useCallback(() => {
        if (useInfiniteScroll && infiniteQueryResult.hasNextPage && !infiniteQueryResult.isFetchingNextPage) {
            infiniteQueryResult.fetchNextPage();
        }
    }, [useInfiniteScroll, infiniteQueryResult]);
    // Handle refresh (pull-to-refresh on mobile)
    const handleRefresh = useCallback(() => {
        if (useInfiniteScroll) {
            infiniteQueryResult.refetch();
        } else {
            standardQueryResult.refetch();
        }
    }, [useInfiniteScroll, infiniteQueryResult, standardQueryResult]);
    // Handle pagination for desktop
    const handlePageChange = useCallback((page) => {
        dispatch(setFilters({
            page,
            offset: (page - 1) * filters.pageSize
        }));
    }, [dispatch, filters.pageSize]);
    // Handle viewport resize to switch between infinite scroll and pagination
    React.useEffect(() => {
        const handleResize = () => {
            const shouldUseInfinite = window.innerWidth <= 768 || 'ontouchstart' in window;
            if (shouldUseInfinite !== useInfiniteScroll) {
                setUseInfiniteScroll(shouldUseInfinite);
            }
        };
        window.addEventListener('resize', handleResize);
        return () => window.removeEventListener('resize', handleResize);
    }, [useInfiniteScroll]);
    // Conditional rendering based on request state
    if (queryResult.isLoading && !queryResult.data) {
        return <LoadingSpinner/>;
    }
    if (queryResult.isError) {
        return <div className="error-message">Error: {queryResult.error?.message}</div>;
    }
    // Empty state
    if (listings.length === 0) {
        return (
            <div className="empty-state">
                <h2>No {entityType} found</h2>
                <p>Try adjusting your filters or search criteria</p>
            </div>
        );
    }
    // Render item component
    const renderItem = useCallback((item, index) => (
        <MemoizedListItem
            key={item.id}
            item={item}
            onSelect={handleSelectItem}
        />
    ), [handleSelectItem]);
    return (
        <div className="optimized-listing-container">
            <div className="listings-header">
                <h1>{entityType.charAt(0).toUpperCase() + entityType.slice(1)}</h1>
                <div className="category-chips">
                    {categories.map(category => (
                        <span
                            key={category.id}
                            className={`category-chip ${filters.category === category.id ? 'active' : ''}`}
                            onClick={() => dispatch(setFilters({category: category.id}))}
                        >
              {category.name}
            </span>
                    ))}
                </div>
            </div>
            {useInfiniteScroll ? (
                // Mobile: Infinite scroll with pull-to-refresh
                <InfiniteScroll
                    ref={infiniteScrollRef}
                    items={listings}
                    renderItem={renderItem}
                    loadMore={handleLoadMore}
                    hasNextPage={infiniteQueryResult.hasNextPage}
                    isLoading={infiniteQueryResult.isFetchingNextPage}
                    error={infiniteQueryResult.error}
                    onRefresh={handleRefresh}
                    refreshing={infiniteQueryResult.isRefetching}
                    enablePullToRefresh={true}
                    enableVirtualization={listings.length > 100}
                    className="listings-infinite-scroll"
                    style={{height: 'calc(100vh - 200px)'}}
                />
            ) : (
                // Desktop: Traditional grid with pagination
                <>
                    <div className="listings-grid">
                        {listings.map(item => (
                            <MemoizedListItem
                                key={item.id}
                                item={item}
                                onSelect={handleSelectItem}
                            />
                        ))}
                    </div>
                    {paginationData && (
                        <div className="pagination">
                            <button
                                disabled={filters.page <= 1}
                                onClick={() => handlePageChange(filters.page - 1)}
                            >
                                Previous
                            </button>
                            <span>Page {filters.page} of {paginationData.totalPages}</span>
                            <button
                                disabled={filters.page >= paginationData.totalPages}
                                onClick={() => handlePageChange(filters.page + 1)}
                            >
                                Next
                            </button>
                        </div>
                    )}
                </>
            )}
        </div>
    );
});
OptimizedListingContainer.displayName = 'OptimizedListingContainer';
export default memo(OptimizedListingContainer); 