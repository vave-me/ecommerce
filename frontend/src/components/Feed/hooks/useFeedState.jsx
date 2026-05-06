/**
 * FEED STATE HOOK
 * Manages core feed state with optimized performance
 */

import { useState, useEffect, useRef } from 'react';
import { useCleanup } from '../../../hooks/useCleanup';

export const useFeedState = (initialFeedItems = [], initialHasMore = true, initialParams = {}) => {
  const [feedItems, setFeedItems] = useState(initialFeedItems);
  const [isLoading, setIsLoading] = useState(false);
  const [hasMore, setHasMore] = useState(initialHasMore);
  const [page, setPage] = useState(initialParams?.page || 1);
  const [filterParams, setFilterParams] = useState(initialParams);
  const [error, setError] = useState(null);
  const [emptyResponseCount, setEmptyResponseCount] = useState(0);
  
  // Performance optimization refs
  const mountedRef = useRef(true);
  const lastRequestRef = useRef(null);
  const { isMounted } = useCleanup();

  // Initialize feed items from props
  useEffect(() => {
    if (initialFeedItems.length > 0 && feedItems.length === 0) {
      setFeedItems(initialFeedItems);
      setFilterParams(initialParams);
      setPage(initialParams?.page || 1);
      setHasMore(initialHasMore);
    }
  }, [initialFeedItems, initialParams, initialHasMore, feedItems.length]);

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      mountedRef.current = false;
    };
  }, []);

  return {
    // State
    feedItems,
    isLoading,
    hasMore,
    page,
    filterParams,
    error,
    emptyResponseCount,
    
    // Setters
    setFeedItems,
    setIsLoading,
    setHasMore,
    setPage,
    setFilterParams,
    setError,
    setEmptyResponseCount,
    
    // Refs
    mountedRef,
    lastRequestRef,
    isMounted
  };
}; 