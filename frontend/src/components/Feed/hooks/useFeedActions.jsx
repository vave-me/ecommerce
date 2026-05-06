/**
 * FEED ACTIONS HOOK
 * Manages feed loading and data operations with optimized performance
 */

import { useCallback, useMemo } from 'react';

export const useFeedActions = (feedState) => {
  const {
    setFeedItems,
    setIsLoading,
    setHasMore,
    setPage,
    setError,
    setEmptyResponseCount,
    mountedRef,
    lastRequestRef,
    isMounted
  } = feedState;

  // Optimized API handler with dynamic imports
  const getEntityApiHandler = useCallback(async (entityType, params) => {
    // Dynamic import to reduce initial bundle size
    const apiModule = await import('../../../api/searchApi');
    const {
      searchDealsWithFilters,
      searchProductsWithFilters,
      searchServicesWithFilters,
      searchJobsWithFilters,
      searchVehiclesWithFilters,
      searchPropertiesWithFilters,
      searchPostsWithFilters,
      unifiedSearch
    } = apiModule;

    // Optimized API handler mapping
    const apiHandlers = {
      deal: searchDealsWithFilters,
      product: searchProductsWithFilters,
      service: searchServicesWithFilters,
      job: searchJobsWithFilters,
      vehicle: searchVehiclesWithFilters,
      property: searchPropertiesWithFilters,
      post: searchPostsWithFilters,
      news: searchPostsWithFilters,
      tweet: searchPostsWithFilters,
      video: searchProductsWithFilters,
      short: searchProductsWithFilters
    };

    // If entityType is null or not in apiHandlers, use unifiedSearch
    if (!entityType || !apiHandlers[entityType]) {
      // For unified feed, we need to transform the params to include entity types
      const unifiedParams = {
        feedType: params.feedType || 'latest',
        entityTypes: ['product', 'post', 'service'],
        page: params.page || 1,
        pageSize: params.pageSize || 20,
        ...params
      };
      const response = await unifiedSearch(unifiedParams);
      // Transform unified search response to match feed format
      return {
        data: response.results || [],
        hasMore: response.hasMore !== false,
        totalCount: response.totalCount || 0
      };
    }
    
    // Use specific handler for single entity type
    const handler = apiHandlers[entityType];
    return handler(params);
  }, []);

  // Optimized parameter cleaning
  const cleanFilters = useCallback((filters = {}) => {
    const cleaned = {};
    const numericFields = ['minPrice', 'maxPrice', 'minYear', 'maxYear', 'minArea', 'maxArea', 'minStock', 'maxStock'];
    
    for (const [key, value] of Object.entries(filters)) {
      // Skip empty values
      if (value === null || value === undefined || value === '') continue;
      if (Array.isArray(value) && value.length === 0) continue;

      // Handle numeric fields
      if (numericFields.includes(key)) {
        const parsedValue = Number(value);
        if (!isNaN(parsedValue)) {
          cleaned[key] = parsedValue;
        }
      } else {
        cleaned[key] = value;
      }
    }
    return cleaned;
  }, []);

  // Optimized parameter processing
  const getRequiredParams = useCallback((entityType, allParams) => {
    const processedParams = {
      ...allParams,
      sortOrder: allParams.sortOrder || 'asc'
    };

    // Remove empty values efficiently
    return Object.fromEntries(
      Object.entries(processedParams).filter(([, value]) => 
        value !== undefined && value !== null && value !== ''
      )
    );
  }, []);

  // Core feed loading function
  const loadFeed = useCallback(async (paramsToLoad = {}) => {
    // Prevent multiple simultaneous requests
    if (!isMounted() || !mountedRef.current) return;

    const requestId = Date.now();
    lastRequestRef.current = requestId;

    try {
      setIsLoading(true);
      setError(null);

      // Merge and clean parameters
      const effectiveParams = {
        ...feedState.filterParams,
        ...paramsToLoad,
        page: paramsToLoad.page || 1
      };

      // Handle pageSize normalization
      if (effectiveParams.page_size && !effectiveParams.pageSize) {
        effectiveParams.pageSize = Number(effectiveParams.page_size);
        delete effectiveParams.page_size;
      }

      const cleanedParams = cleanFilters(effectiveParams);
      // For unified feed, we should not default to a single entity type
      // If no contentType is specified, use unified feed to show all entity types
      const entityType = cleanedParams.contentType || null;
      const requiredParams = getRequiredParams(entityType, cleanedParams);

      // Check if this request is still current
      if (lastRequestRef.current !== requestId) return;

      // Make API call - getEntityApiHandler already returns the promise/response
      const response = await getEntityApiHandler(entityType, requiredParams);

      // Check again if request is still current
      if (lastRequestRef.current !== requestId || !isMounted()) return;

      if (response?.data) {
        const newItems = Array.isArray(response.data) ? response.data : [];
        
        if (effectiveParams.page === 1) {
          // Replace feed items for new filter
          setFeedItems(newItems);
          setEmptyResponseCount(newItems.length === 0 ? 1 : 0);
        } else {
          // Append for pagination
          setFeedItems(prev => [...prev, ...newItems]);
          setEmptyResponseCount(prev => newItems.length === 0 ? prev + 1 : 0);
        }

        setPage(effectiveParams.page);
        setHasMore(newItems.length > 0 && response.hasMore !== false);
      } else {
        if (effectiveParams.page === 1) {
          setFeedItems([]);
        }
        setHasMore(false);
      }

      // Update filter params
      if (effectiveParams.page === 1) {
        setFilterParams(cleanedParams);
      }

    } catch (error) {
      if (lastRequestRef.current === requestId && isMounted()) {
        // Error: 'Feed loading error:', error...
        setError(error.message || 'Failed to load feed');
        
        if (paramsToLoad.page === 1) {
          setFeedItems([]);
        }
        setHasMore(false);
      }
    } finally {
      if (lastRequestRef.current === requestId && isMounted()) {
        setIsLoading(false);
      }
    }
  }, [feedState.filterParams, cleanFilters, getRequiredParams, getEntityApiHandler, isMounted, mountedRef, setIsLoading, setError, setFeedItems, setEmptyResponseCount, setPage, setHasMore, setFilterParams]);

  // Load more with pagination
  const loadMore = useCallback(async () => {
    if (!feedState.hasMore || feedState.isLoading || feedState.emptyResponseCount >= 3) return;
    
    await loadFeed({
      ...feedState.filterParams,
      page: feedState.page + 1
    });
  }, [feedState.hasMore, feedState.isLoading, feedState.emptyResponseCount, feedState.filterParams, feedState.page, loadFeed]);

  // Memoized actions object
  const actions = useMemo(() => ({
    loadFeed,
    loadMore,
    cleanFilters,
    getRequiredParams,
    getEntityApiHandler
  }), [loadFeed, loadMore, cleanFilters, getRequiredParams, getEntityApiHandler]);

  return actions;
}; 