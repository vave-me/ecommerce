import { useCallback, useState } from 'react';
import { useQueryClient } from '@tanstack/react-query';
/**
 * Custom hook that triggers data prefetching when hovering over an element
 * 
 * @param {Function} prefetchFn - Function to execute on hover
 * @param {number} hoverDelay - Delay in ms before triggering prefetch (default: 300ms)
 * @returns {Object} Object containing event handlers for hover
 */
export function usePrefetchOnHover(prefetchFn, hoverDelay = 300) {
  const [timeoutId, setTimeoutId] = useState(null);
  const [hasPrefetched, setHasPrefetched] = useState(false);
  const handleMouseEnter = useCallback(() => {
    if (hasPrefetched) return;
    const id = setTimeout(() => {
      prefetchFn();
      setHasPrefetched(true);
    }, hoverDelay);
    setTimeoutId(id);
  }, [prefetchFn, hoverDelay, hasPrefetched]);
  const handleMouseLeave = useCallback(() => {
    if (timeoutId) {
      clearTimeout(timeoutId);
      setTimeoutId(null);
    }
  }, [timeoutId]);
  return {
    onMouseEnter: handleMouseEnter,
    onMouseLeave: handleMouseLeave,
    hasPrefetched
  };
}
/**
 * Hook to prefetch images when hovering
 * 
 * @param {Array} imageUrls - Array of image URLs to prefetch
 * @param {number} hoverDelay - Delay in ms before starting prefetch
 * @returns {Object} Object containing event handlers for hover
 */
export function useImagePrefetch(imageUrls = [], hoverDelay = 300) {
  const prefetchImages = useCallback(() => {
    if (!imageUrls.length) return;
    // Prefetch first 3 images immediately
    imageUrls.slice(0, 3).forEach(url => {
      if (!url) return;
      const img = new Image();
      img.src = url;
    });
    // Prefetch remaining images with delay
    if (imageUrls.length > 3) {
      setTimeout(() => {
        imageUrls.slice(3).forEach(url => {
          if (!url) return;
          const img = new Image();
          img.src = url;
        });
      }, 500);
    }
  }, [imageUrls]);
  return usePrefetchOnHover(prefetchImages, hoverDelay);
}
/**
 * Hook to prefetch data with React Query when hovering
 * 
 * @param {Array} queryKey - React Query key for the data to prefetch
 * @param {Function} queryFn - Function to fetch the data
 * @param {number} hoverDelay - Delay in ms before starting prefetch
 * @returns {Object} Object containing event handlers for hover
 */
export function useQueryPrefetch(queryKey, queryFn, hoverDelay = 300) {
  const queryClient = useQueryClient();
  const prefetchQuery = useCallback(() => {
    queryClient.prefetchQuery({
      queryKey,
      queryFn,
      staleTime: 5 * 60 * 1000, // 5 minutes
    });
  }, [queryClient, queryKey, queryFn]);
  return usePrefetchOnHover(prefetchQuery, hoverDelay);
} 