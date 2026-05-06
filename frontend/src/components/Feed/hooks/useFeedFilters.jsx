/**
 * FEED FILTERS HOOK
 * Manages filter updates and Redux integration with optimized performance
 */

import { useCallback, useMemo } from 'react';
import { useDispatch } from 'react-redux';
import { setFilters } from '../../../redux/slices/listingFiltersSlice';

export const useFeedFilters = (feedState, feedActions) => {
  const dispatch = useDispatch();
  const { loadFeed } = feedActions;

  // Optimized filter update with Redux sync
  const updateFilters = useCallback(async (newFilters) => {
    try {
      // Clean and prepare filters
      const cleanedFilters = feedActions.cleanFilters(newFilters);
      
      // Reset to page 1 for new filters
      const filtersWithPage = {
        ...cleanedFilters,
        page: 1
      };

      // Update Redux store
      dispatch(setFilters(filtersWithPage));

      // Load feed with new filters
      await loadFeed(filtersWithPage);
      
    } catch (error) {
      // Error: 'Filter update error:', error...
      feedState.setError('Failed to update filters');
    }
  }, [feedActions, dispatch, loadFeed, feedState]);

  // Optimized filter reset
  const resetFilters = useCallback(async () => {
    const defaultFilters = { page: 1 };
    
    try {
      dispatch(setFilters(defaultFilters));
      await loadFeed(defaultFilters);
    } catch (error) {
      // Error: 'Filter reset error:', error...
      feedState.setError('Failed to reset filters');
    }
  }, [dispatch, loadFeed, feedState]);

  // Get current filter state
  const getCurrentFilters = useCallback(() => {
    return { ...feedState.filterParams };
  }, [feedState.filterParams]);

  // Check if filters are active
  const hasActiveFilters = useMemo(() => {
    const params = feedState.filterParams;
    return Object.keys(params).some(key => 
      key !== 'page' && 
      key !== 'pageSize' && 
      params[key] !== null && 
      params[key] !== undefined && 
      params[key] !== ''
    );
  }, [feedState.filterParams]);

  // Memoized filter functions
  const filterFunctions = useMemo(() => ({
    updateFilters,
    resetFilters,
    getCurrentFilters,
    hasActiveFilters
  }), [updateFilters, resetFilters, getCurrentFilters, hasActiveFilters]);

  return filterFunctions;
};
