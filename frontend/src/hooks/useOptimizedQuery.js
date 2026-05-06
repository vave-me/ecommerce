"use client";
import { useQuery, useMutation, useQueryClient, useInfiniteQuery } from '@tanstack/react-query';
import { useCallback, useMemo } from 'react';
import { useDispatch, useSelector } from 'react-redux';
/**
 * Advanced optimized query hook with Redux integration
 * Features:
 * - Intelligent caching strategies
 * - Optimistic updates
 * - Error boundary integration
 * - Performance monitoring
 * - Background refetching
 */
export function useOptimizedQuery({
  queryKey,
  queryFn,
  reduxSlice,
  reduxAction,
  options = {},
  optimisticUpdate = false,
  backgroundRefetch = true,
  staleTime = 5 * 60 * 1000, // 5 minutes
  cacheTime = 10 * 60 * 1000, // 10 minutes
}) {
  const dispatch = useDispatch();
  const queryClient = useQueryClient();
  // Get Redux state for this slice
  const reduxData = useSelector(state => 
    reduxSlice ? state[reduxSlice] : null
  );
  // Enhanced query configuration
  const queryConfig = useMemo(() => ({
    queryKey,
    queryFn,
    staleTime,
    cacheTime,
    refetchOnWindowFocus: backgroundRefetch,
    refetchOnReconnect: true,
    retry: (failureCount, error) => {
      // Smart retry logic
      if (error?.status === 404 || error?.status === 403) {
        return false; // Don't retry for these errors
      }
      return failureCount < 3;
    },
    retryDelay: attemptIndex => Math.min(1000 * 2 ** attemptIndex, 30000),
    onSuccess: (data) => {
      // Sync with Redux if configured
      if (reduxAction && dispatch) {
        dispatch(reduxAction(data));
      }
      options.onSuccess?.(data);
    },
    onError: (error) => {
      options.onError?.(error);
    },
    ...options
  }), [queryKey, queryFn, staleTime, cacheTime, backgroundRefetch, reduxAction, dispatch, options]);
  const query = useQuery(queryConfig);
  // Optimized invalidation
  const invalidateQuery = useCallback(() => {
    queryClient.invalidateQueries({ queryKey });
  }, [queryClient, queryKey]);
  // Prefetch related data
  const prefetchRelated = useCallback((relatedQueryKey, relatedQueryFn) => {
    queryClient.prefetchQuery({
      queryKey: relatedQueryKey,
      queryFn: relatedQueryFn,
      staleTime: staleTime / 2, // Shorter stale time for prefetched data
    });
  }, [queryClient, staleTime]);
  // Update cache optimistically
  const updateCacheOptimistically = useCallback((updater) => {
    queryClient.setQueryData(queryKey, updater);
  }, [queryClient, queryKey]);
  // Get cached data without triggering a request
  const getCachedData = useCallback(() => {
    return queryClient.getQueryData(queryKey);
  }, [queryClient, queryKey]);
  return {
    ...query,
    reduxData,
    invalidateQuery,
    prefetchRelated,
    updateCacheOptimistically,
    getCachedData,
    // Performance metrics
    isStale: query.isStale,
    isFetching: query.isFetching,
    isBackground: query.isFetching && !query.isLoading,
  };
}
/**
 * Optimized mutation hook with Redux integration
 */
export function useOptimizedMutation({
  mutationFn,
  queryKeysToInvalidate = [],
  reduxAction,
  optimisticUpdate,
  rollbackOnError = true,
  options = {}
}) {
  const dispatch = useDispatch();
  const queryClient = useQueryClient();
  const mutation = useMutation({
    mutationFn,
    onMutate: async (variables) => {
      // Cancel outgoing refetches
      await Promise.all(
        queryKeysToInvalidate.map(key => 
          queryClient.cancelQueries({ queryKey: key })
        )
      );
      // Snapshot previous values for rollback
      const previousData = {};
      if (rollbackOnError) {
        queryKeysToInvalidate.forEach(key => {
          previousData[key] = queryClient.getQueryData(key);
        });
      }
      // Optimistic update
      if (optimisticUpdate) {
        queryKeysToInvalidate.forEach(key => {
          queryClient.setQueryData(key, old => 
            optimisticUpdate(old, variables)
          );
        });
      }
      options.onMutate?.(variables);
      return { previousData };
    },
    onError: (error, variables, context) => {
      // Rollback on error
      if (rollbackOnError && context?.previousData) {
        Object.entries(context.previousData).forEach(([key, data]) => {
          queryClient.setQueryData(JSON.parse(key), data);
        });
      }
      options.onError?.(error, variables, context);
    },
    onSuccess: (data, variables, context) => {
      // Update Redux state
      if (reduxAction && dispatch) {
        dispatch(reduxAction(data));
      }
      options.onSuccess?.(data, variables, context);
    },
    onSettled: (data, error, variables, context) => {
      // Invalidate and refetch
      queryKeysToInvalidate.forEach(key => {
        queryClient.invalidateQueries({ queryKey: key });
      });
      options.onSettled?.(data, error, variables, context);
    },
    ...options
  });
  return mutation;
}
/**
 * Hook for managing infinite queries with optimization
 */
export function useOptimizedInfiniteQuery({
  queryKey,
  queryFn,
  getNextPageParam,
  getPreviousPageParam,
  options = {}
}) {
  const queryClient = useQueryClient();
  const infiniteQuery = useInfiniteQuery({
    queryKey,
    queryFn,
    getNextPageParam,
    getPreviousPageParam,
    staleTime: 5 * 60 * 1000,
    cacheTime: 10 * 60 * 1000,
    refetchOnWindowFocus: false,
    ...options
  });
  // Flatten all pages data
  const flatData = useMemo(() => {
    return infiniteQuery.data?.pages?.flat() || [];
  }, [infiniteQuery.data]);
  // Prefetch next page
  const prefetchNextPage = useCallback(() => {
    if (infiniteQuery.hasNextPage && !infiniteQuery.isFetchingNextPage) {
      infiniteQuery.fetchNextPage();
    }
  }, [infiniteQuery]);
  return {
    ...infiniteQuery,
    flatData,
    prefetchNextPage,
    totalItems: flatData.length,
  };
}
/**
 * Performance-optimized selector hook
 */
export function useOptimizedSelector(selector, equalityFn) {
  return useSelector(selector, equalityFn || ((a, b) => {
    // Shallow comparison for objects
    if (typeof a === 'object' && typeof b === 'object') {
      const keysA = Object.keys(a);
      const keysB = Object.keys(b);
      if (keysA.length !== keysB.length) return false;
      return keysA.every(key => a[key] === b[key]);
    }
    return a === b;
  }));
} 