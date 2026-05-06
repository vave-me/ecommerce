import { useInfiniteQuery } from '@tanstack/react-query';
import { useSelector } from 'react-redux';
import { useMemo, useCallback } from 'react';
import { unifiedSearch } from '../api/searchApi';

/**
 * Single unified feed hook
 * - Uses Redux for filter state
 * - Uses React Query for data fetching
 * - No context providers, no duplicates
 */
export function useFeed(options = {}) {
  // Get filters from Redux (single source of truth)
  const filters = useSelector(state => state.listingFilters);

  // Query configuration
  const {
    enabled = true,
    entityType = null, // Specific entity type filter
  } = options;
  
  // Build query key
  const queryKey = useMemo(() => {
    const key = ['feed'];
    if (entityType) key.push(entityType);
    key.push(filters);
    return key;
  }, [filters, entityType]);
  
  // Query function
  const queryFn = useCallback(async ({ pageParam = 1 }) => {
    try {
      // ALWAYS use unifiedSearch
      const searchParams = {
        searchTerm: filters.searchText || undefined,
        entityTypes: entityType ? [entityType] : filters.entityTypes || ['product', 'post', 'service'],
        page: pageParam,
        pageSize: filters.pageSize || 20,
        categorySlug: filters.categorySlug || undefined,
        minPrice: filters.minPrice ? parseInt(filters.minPrice) : undefined,
        maxPrice: filters.maxPrice ? parseInt(filters.maxPrice) : undefined,
        userType: filters.userType || undefined,
        negotiable: filters.negotiable,
        sortBy: filters.sortBy || undefined,
        sortOrder: filters.sortOrder || undefined
      };
      
      // Remove undefined values
      Object.keys(searchParams).forEach(key => 
        searchParams[key] === undefined && delete searchParams[key]
      );

      const response = await unifiedSearch(searchParams);
      
      // Transform response to feed format
      return {
        items: response.results || [],
        hasMore: response.hasMore !== false,
        page: pageParam,
        totalCount: response.totalCount || 0,
      };
    } catch (error) {
      // Error: '[useFeed] Error fetching feed:', error...
      // Return empty result on error to prevent infinite retries
      return {
        items: [],
        hasMore: false,
        page: pageParam,
        totalCount: 0,
        error: error.message
      };
    }
  }, [filters, entityType]);
  
  // Use React Query infinite query
  const query = useInfiniteQuery({
    queryKey,
    queryFn,
    getNextPageParam: (lastPage) => {
      if (!lastPage.hasMore || lastPage.error) return undefined;
      return lastPage.page + 1;
    },
    enabled,
    staleTime: 30000,
    gcTime: 5 * 60 * 1000,
    refetchOnWindowFocus: false, // Prevent refetch on window focus
    retry: (failureCount, error) => {
      // Don't retry on 404 or 401 errors
      if (error?.response?.status === 404 || error?.response?.status === 401) {
        return false;
      }
      return failureCount < 2;
    },
    retryDelay: 1000,
  });
  
  // Flatten pages into single array
  const items = useMemo(() => {
    if (!query.data?.pages) return [];
    return query.data.pages.flatMap(page => page.items || []);
  }, [query.data]);
  
  // Handle errors gracefully
  const error = query.error || query.data?.pages?.find(p => p.error)?.error;
  
  return {
    // Data
    items,
    feedItems: items, // Alias for backward compatibility
    
    // State
    isLoading: query.isLoading,
    error,
    hasMore: query.hasNextPage && !error,
    
    // Actions
    loadMore: query.fetchNextPage,
    refetch: query.refetch,
    
    // Additional info
    isFetchingNextPage: query.isFetchingNextPage,
    totalCount: query.data?.pages?.[0]?.totalCount || 0,
    isError: query.isError || !!error,
  };
}

// Export as default
export default useFeed;