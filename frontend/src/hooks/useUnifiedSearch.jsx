import { useCallback, useMemo, useRef, useEffect } from "react";
import { useQuery, useInfiniteQuery, useQueryClient } from "@tanstack/react-query";
import { useSelector, useDispatch } from "react-redux";
import { shallowEqual } from "react-redux";
import { setFilters } from "../redux/slices/listingFiltersSlice";
import * as searchApi from "../api/searchApi";
import { QUERY_KEYS, CACHE_TIMES } from "../lib/reactQuery";
import { getEntityTypeFromContentType } from "../utils/entityMapper";
import debugFilters from "../utils/debugFilters";

/**
 * Unified search hook that combines useSearch and useFilteredSearch
 * - Supports both single page and infinite scroll queries
 * - Integrates with Redux for filter state management
 * - Handles all entity types through a unified interface
 * - Implements debouncing, caching, and prefetching
 * 
 * @param {Object} options - Hook configuration
 * @param {string} options.entityType - Entity type to search
 * @param {Object} options.filters - Override filters (if not using Redux)
 * @param {boolean} options.infinite - Use infinite query (default: false)
 * @param {boolean} options.useReduxFilters - Use Redux filter state (default: true)
 * @param {number} options.pageSize - Page size for pagination
 * @param {boolean} options.enabled - Enable/disable query
 * @param {Function} options.onSuccess - Success callback
 * @param {Function} options.onError - Error callback
 * @returns {Object} Search state and utilities
 */
export function useUnifiedSearch(options = {}) {
  const {
    entityType: explicitEntityType,
    filters: explicitFilters,
    infinite = false,
    useReduxFilters = true,
    pageSize = 20,
    staleTime = CACHE_TIMES.DYNAMIC || 30000,
    gcTime = CACHE_TIMES.STANDARD || 5 * 60 * 1000,
    enabled = true,
    onSuccess,
    onError,
  } = options;

  const dispatch = useDispatch();
  const queryClient = useQueryClient();
  const abortControllerRef = useRef(null);

  // Get filters from Redux or use explicit filters
  const reduxFilters = useSelector((state) => state.listingFilters, shallowEqual);
  const filters = useReduxFilters ? reduxFilters : (explicitFilters || {});

  // Determine effective entity type
  const effectiveEntityType = useMemo(() => {
    if (explicitEntityType) return explicitEntityType;
    
    if (filters.contentType) {
      const mappedType = getEntityTypeFromContentType(filters.contentType);
      if (mappedType) return mappedType;
    }
    
    if (filters.entityType) return filters.entityType;
    
    if (filters.listingType) {
      return filters.listingType.endsWith("s") 
        ? filters.listingType.slice(0, -1) 
        : filters.listingType;
    }
    
    return "product";
  }, [explicitEntityType, filters.contentType, filters.entityType, filters.listingType]);

  // API method mapping
  const apiMethodMap = useMemo(() => ({
    product: searchApi.searchProductsWithFilters,
    post: searchApi.searchPostsWithFilters,
    service: searchApi.searchServicesWithFilters,
    video: searchApi.searchProductsWithFilters,
    tweet: searchApi.searchPostsWithFilters,
    job: searchApi.searchPostsWithFilters,
    deal: searchApi.searchProductsWithFilters,
    short: searchApi.searchProductsWithFilters,
    vehicle: searchApi.searchProductsWithFilters,
    property: searchApi.searchProductsWithFilters,
    news: searchApi.searchPostsWithFilters,
  }), []);

  // Get API method
  const apiMethod = useMemo(() => {
    return apiMethodMap[effectiveEntityType] || searchApi.unifiedSearch;
  }, [effectiveEntityType, apiMethodMap]);

  // Clean filters
  const cleanFilters = useCallback((filtersToClean) => {
    return Object.entries(filtersToClean).reduce((acc, [key, value]) => {
      if (value !== undefined && value !== null && value !== "") {
        if (Array.isArray(value) && value.length === 0) return acc;
        acc[key] = value;
      }
      return acc;
    }, {});
  }, []);

  // Transform filters for API
  const apiFilters = useMemo(() => {
    const cleaned = cleanFilters(filters);
    return transformFiltersForApi(cleaned, effectiveEntityType);
  }, [filters, effectiveEntityType, cleanFilters]);

  // Create query key
  const queryKey = useMemo(() => {
    const keyBase = ['search'];
    if (infinite) keyBase.push('infinite');
    keyBase.push(effectiveEntityType);
    keyBase.push(cleanFilters(filters));
    if (!infinite && filters.page) keyBase.push({ page: filters.page });
    return keyBase;
  }, [infinite, effectiveEntityType, filters, cleanFilters]);

  // Query function for single page
  const singlePageQueryFn = useCallback(async ({ signal }) => {
    abortControllerRef.current = new AbortController();
    const combinedSignal = signal || abortControllerRef.current.signal;

    try {
      debugFilters.logFilterState('useUnifiedSearch.queryFn', {
        entityType: effectiveEntityType,
        apiFilters
      });

      const result = await apiMethod(apiFilters, { signal: combinedSignal });
      return result;
    } catch (error) {
      if (error.name === "AbortError") return null;
      throw error;
    }
  }, [effectiveEntityType, apiFilters, apiMethod]);

  // Query function for infinite scroll
  const infiniteQueryFn = useCallback(async ({ pageParam = 1, signal }) => {
    abortControllerRef.current = new AbortController();
    const combinedSignal = signal || abortControllerRef.current.signal;

    try {
      const paginatedFilters = { ...apiFilters, page: pageParam, pageSize };
      const response = await apiMethod(paginatedFilters, { signal: combinedSignal });
      
      // Normalize response
      const items = response[effectiveEntityType] || response.items || response.results || [];
      const hasMore = response.hasMore !== false && items.length >= pageSize;
      
      return {
        items,
        hasMore,
        page: pageParam,
        totalCount: response.totalCount || items.length,
      };
    } catch (error) {
      if (error.name === "AbortError") return null;
      // Error: `[useUnifiedSearch] Error:`, error...
      return {
        items: [],
        hasMore: false,
        page: pageParam,
        totalCount: 0,
        error: error.message
      };
    }
  }, [apiFilters, pageSize, apiMethod, effectiveEntityType]);

  // Single page query
  const singleQuery = useQuery({
    queryKey,
    queryFn: singlePageQueryFn,
    enabled: enabled && !infinite,
    staleTime,
    gcTime,
    keepPreviousData: true,
    retry: (failureCount, error) => {
      if (error?.response?.status === 404 || error?.response?.status === 401) {
        return false;
      }
      return failureCount < 2;
    },
    onSuccess,
    onError,
  });

  // Infinite query
  const infiniteQuery = useInfiniteQuery({
    queryKey,
    queryFn: infiniteQueryFn,
    getNextPageParam: (lastPage) => {
      if (!lastPage?.hasMore || lastPage.error) return undefined;
      return lastPage.page + 1;
    },
    enabled: enabled && infinite,
    staleTime,
    gcTime,
    refetchOnWindowFocus: false,
    retry: (failureCount, error) => {
      if (error?.response?.status === 404 || error?.response?.status === 401) {
        return false;
      }
      return failureCount < 2;
    },
    onSuccess,
    onError,
  });

  // Select active query
  const activeQuery = infinite ? infiniteQuery : singleQuery;

  // Flatten items for infinite query
  const items = useMemo(() => {
    if (infinite && infiniteQuery.data?.pages) {
      return infiniteQuery.data.pages.flatMap(page => page?.items || []);
    }
    
    // For single query, return items from various possible response formats
    const data = singleQuery.data;
    if (!data) return [];
    
    return data[effectiveEntityType] || data.items || data.results || [];
  }, [infinite, infiniteQuery.data, singleQuery.data, effectiveEntityType]);

  // Update filters (Redux integration)
  const updateFilters = useCallback((newFilters) => {
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }
    
    if (useReduxFilters) {
      dispatch(setFilters(newFilters));
    }
  }, [dispatch, useReduxFilters]);

  // Prefetch results
  const prefetchResults = useCallback((anticipatedFilters) => {
    const mergedFilters = { ...filters, ...anticipatedFilters };
    const prefetchEntityType = anticipatedFilters.entityType || 
                              (anticipatedFilters.contentType && getEntityTypeFromContentType(anticipatedFilters.contentType)) ||
                              effectiveEntityType;
    
    const prefetchApiMethod = apiMethodMap[prefetchEntityType] || searchApi.unifiedSearch;
    const prefetchApiFilters = transformFiltersForApi(mergedFilters, prefetchEntityType);
    
    const prefetchKey = ['search'];
    if (infinite) prefetchKey.push('infinite');
    prefetchKey.push(prefetchEntityType);
    prefetchKey.push(cleanFilters(mergedFilters));
    
    if (JSON.stringify(prefetchKey) === JSON.stringify(queryKey)) return;
    
    queryClient.prefetchQuery({
      queryKey: prefetchKey,
      queryFn: () => prefetchApiMethod(prefetchApiFilters),
      staleTime,
    });
  }, [queryClient, filters, effectiveEntityType, apiMethodMap, staleTime, queryKey, cleanFilters, infinite]);

  // Cancel search
  const cancelSearch = useCallback(() => {
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
      abortControllerRef.current = null;
    }
  }, []);

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }
    };
  }, []);

  return {
    // Data
    items,
    data: activeQuery.data,
    
    // State
    isLoading: activeQuery.isLoading,
    isError: activeQuery.isError,
    error: activeQuery.error,
    isFetching: activeQuery.isFetching,
    
    // Infinite scroll specific
    hasNextPage: infinite ? infiniteQuery.hasNextPage : false,
    isFetchingNextPage: infinite ? infiniteQuery.isFetchingNextPage : false,
    fetchNextPage: infinite ? infiniteQuery.fetchNextPage : undefined,
    
    // Actions
    refetch: activeQuery.refetch,
    updateFilters: useReduxFilters ? updateFilters : undefined,
    prefetchResults,
    cancelSearch,
    
    // Additional info
    entityType: effectiveEntityType,
    filters,
    totalCount: infinite 
      ? infiniteQuery.data?.pages?.[0]?.totalCount || 0
      : singleQuery.data?.totalCount || items.length,
  };
}

/**
 * Transform filters for API consumption based on entity type
 * (Reused from useFilteredSearch)
 */
function transformFiltersForApi(filters, entityType) {
  const apiFilters = {
    sortOrder: filters.sortOrder || 'asc'
  };
  
  // Map searchText to name
  if (filters.searchText) {
    apiFilters.name = filters.searchText;
  }
  
  // Common filters
  if (filters.tags?.length > 0) apiFilters.tags = filters.tags;
  if (filters.location) apiFilters.location = filters.location;
  if (filters.page) apiFilters.page = Number(filters.page);
  if (filters.pageSize) apiFilters.pageSize = Number(filters.pageSize);
  if (filters.limit) apiFilters.limit = Number(filters.limit);
  if (filters.offset) apiFilters.offset = Number(filters.offset);
  if (filters.sortBy) apiFilters.sortBy = filters.sortBy;
  
  // Price filters
  if (filters.minPrice !== undefined && filters.minPrice !== '') {
    apiFilters.minPrice = Number(filters.minPrice);
  }
  if (filters.maxPrice !== undefined && filters.maxPrice !== '') {
    apiFilters.maxPrice = Number(filters.maxPrice);
  }
  
  // Category handling
  if (filters.categoryID) {
    apiFilters.categoryId = filters.categoryID;
  }
  if (filters.categorySlug) {
    apiFilters.categorySlug = filters.categorySlug;
  }
  if (!filters.categoryID && !filters.categorySlug && filters.category) {
    if (isNaN(filters.category)) {
      apiFilters.categorySlug = filters.category;
    } else {
      apiFilters.categoryId = filters.category;
    }
  }
  
  // Entity-specific filters (simplified for brevity)
  switch (entityType) {
    case 'product':
      if (filters.brand) apiFilters.brand = filters.brand;
      if (filters.condition) apiFilters.condition = filters.condition;
      if (filters.negotiable !== undefined) apiFilters.negotiable = filters.negotiable;
      if (filters.userType) apiFilters.userType = filters.userType;
      break;
      
    case 'post':
      if (filters.content) apiFilters.content = filters.content;
      if (filters.title) apiFilters.title = filters.title;
      if (filters.status) apiFilters.status = filters.status;
      break;
      
    case 'service':
      if (filters.serviceType) apiFilters.serviceType = filters.serviceType;
      if (filters.availability) apiFilters.availability = filters.availability;
      if (filters.isOnline !== undefined) apiFilters.isOnline = filters.isOnline;
      break;
  }
  
  return apiFilters;
}

// Export convenience hooks for backward compatibility
export const useSearch = (entityType, filters, options) => 
  useUnifiedSearch({ entityType, filters, useReduxFilters: false, ...options });

export const useFilteredSearch = (options) => 
  useUnifiedSearch({ useReduxFilters: true, ...options });

export const useInfiniteSearch = (entityType, filters, options) => 
  useUnifiedSearch({ entityType, filters, infinite: true, useReduxFilters: false, ...options });

/**
 * User catalog search hook
 * @param {string} userId - User ID
 * @param {Object} filters - Search filters
 * @param {Object} options - Query options
 */
export function useUserCatalog(userId, filters = {}, options = {}) {
  const {
    enabled = !!userId,
    staleTime = 30000,
    gcTime = 5 * 60 * 1000,
  } = options;

  const queryKey = useMemo(() => 
    ['catalog', userId, filters],
    [userId, filters]
  );

  const queryFn = useCallback(async () => {
    if (!userId) throw new Error('User ID is required');
    return searchApi.getCatalog(userId, { ...filters, ...options });
  }, [userId, filters, options]);

  return useQuery({
    queryKey,
    queryFn,
    enabled,
    staleTime,
    gcTime,
  });
}

/**
 * Quick search hook for search suggestions
 * @param {string} searchTerm - Search term
 * @param {Object} options - Query options
 */
export function useSearchSuggestions(searchTerm, options = {}) {
  const {
    entityType = 'products',
    enabled = searchTerm?.length > 2,
    staleTime = 60000,
    debounceMs = 300,
  } = options;

  const queryKey = useMemo(() => 
    ['search', 'suggestions', entityType, searchTerm],
    [entityType, searchTerm]
  );

  const queryFn = useCallback(async () => {
    // Quick search with limited results
    const response = await searchApi.unifiedSearch({
      searchTerm,
      entityTypes: [entityType],
      page: 1,
      pageSize: 5
    });
    
    return response.results || [];
  }, [entityType, searchTerm]);

  return useQuery({
    queryKey,
    queryFn,
    enabled,
    staleTime,
    gcTime: staleTime * 2,
  });
}

// Default export
export default useUnifiedSearch;