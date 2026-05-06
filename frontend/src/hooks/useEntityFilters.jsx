import { useCallback, useMemo } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { useSelector, shallowEqual } from 'react-redux';
import { 
  selectDisplayFilters, selectListingType, selectPaginationFilters, 
  selectPriceFilters, selectProductSpecificFilters, selectSearchFilters 
} from '../redux/selectors';
import { QUERY_KEYS } from '../lib/reactQuery';
/**
 * Optimized hook to manage entity filters with efficient caching
 * - Uses Redux for global filter UI state
 * - Integrates with React Query for data fetching
 * - Prevents duplicate data in Redux and React Query
 * - Properly memoizes objects to prevent unnecessary renders
 * 
 * @param {string} entityType Type of entity (products, posts, etc)
 * @returns {Object} Filter management utilities
 */
export function useEntityFilters(entityType) {
  const queryClient = useQueryClient();
  // Get all filter sections from Redux using a single selector call with shallowEqual
  const {
    listingType,
    displayFilters,
    paginationFilters,
    searchFilters,
    priceFilters,
    productFilters
  } = useSelector(state => ({
    listingType: selectListingType(state),
    displayFilters: selectDisplayFilters(state),
    paginationFilters: selectPaginationFilters(state),
    searchFilters: selectSearchFilters(state),
    priceFilters: selectPriceFilters(state),
    productFilters: selectProductSpecificFilters(state)
  }), shallowEqual); // Important: Use shallowEqual to prevent unnecessary renders
  // Compute the effective entity type (from props or Redux state) - memoized
  const effectiveEntityType = useMemo(() => 
    entityType || listingType, 
    [entityType, listingType]
  );
  // Merge filters based on entity type - properly memoized
  const allFilters = useMemo(() => ({
    ...displayFilters,
    ...paginationFilters,
    ...searchFilters,
    ...priceFilters,
    ...(effectiveEntityType === 'products' ? productFilters : {})
  }), [
    displayFilters, 
    paginationFilters, 
    searchFilters, 
    priceFilters, 
    productFilters, 
    effectiveEntityType
  ]);
  // Create stable query key generator
  const getQueryKey = useCallback((customFilters = {}) => {
    // Use standardized query key format
    return QUERY_KEYS.listings({
      type: effectiveEntityType,
      ...allFilters,
      ...customFilters
    });
  }, [effectiveEntityType, allFilters]);
  // Invalidate cache when filters change - stable reference
  const invalidateFilters = useCallback(() => {
    queryClient.invalidateQueries({
      queryKey: QUERY_KEYS.listings({ entityType: effectiveEntityType }),
    });
  }, [queryClient, effectiveEntityType]);
  // Prefetch data with the current filters - stable reference
  const prefetchWithFilters = useCallback(async (customFilters = {}, queryFn) => {
    if (!queryFn) return;
    const queryKey = getQueryKey(customFilters);
    await queryClient.prefetchQuery({
      queryKey,
      queryFn: () => queryFn({
        entityType: effectiveEntityType,
        filters: { ...allFilters, ...customFilters }
      }),
      staleTime: 30000 // 30 seconds
    });
  }, [queryClient, getQueryKey, effectiveEntityType, allFilters]);
  // Return stable object
  return {
    filters: allFilters,
    entityType: effectiveEntityType,
    getQueryKey,
    invalidateFilters,
    prefetchWithFilters,
  };
} 