"use client";
import { useEffect, useRef, useState, useCallback } from 'react';
/**
 * High-performance intersection observer hook
 * Features:
 * - Lazy loading support
 * - Visibility tracking
 * - Performance optimized
 * - Multiple threshold support
 */
export function useIntersectionObserver({
  threshold = 0.1,
  root = null,
  rootMargin = '0px',
  freezeOnceVisible = false,
  initialIsIntersecting = false,
  onIntersect,
} = {}) {
  const [isIntersecting, setIsIntersecting] = useState(initialIsIntersecting);
  const [entry, setEntry] = useState(null);
  const elementRef = useRef(null);
  const observerRef = useRef(null);
  const setElement = useCallback((element) => {
    elementRef.current = element;
  }, []);
  useEffect(() => {
    const element = elementRef.current;
    // Skip if no element or already frozen
    if (!element || (freezeOnceVisible && isIntersecting)) {
      return;
    }
    // Create observer if it doesn't exist
    if (!observerRef.current) {
      observerRef.current = new IntersectionObserver(
        (entries) => {
          const [entry] = entries;
          const isElementIntersecting = entry.isIntersecting;
          setEntry(entry);
          setIsIntersecting(isElementIntersecting);
          // Call callback if provided
          if (onIntersect) {
            onIntersect(entry, isElementIntersecting);
          }
          // Unobserve if frozen once visible
          if (freezeOnceVisible && isElementIntersecting && observerRef.current) {
            observerRef.current.unobserve(element);
          }
        },
        {
          threshold,
          root,
          rootMargin,
        }
      );
    }
    // Start observing
    observerRef.current.observe(element);
    // Cleanup
    return () => {
      if (observerRef.current && element) {
        observerRef.current.unobserve(element);
      }
    };
  }, [threshold, root, rootMargin, freezeOnceVisible, isIntersecting, onIntersect]);
  // Cleanup observer on unmount
  useEffect(() => {
    return () => {
      if (observerRef.current) {
        observerRef.current.disconnect();
      }
    };
  }, []);
  return {
    isIntersecting,
    entry,
    ref: setElement,
  };
}
/**
 * Hook for lazy loading images
 */
export function useLazyImage(src, options = {}) {
  const [imageSrc, setImageSrc] = useState(null);
  const [isLoaded, setIsLoaded] = useState(false);
  const [hasError, setHasError] = useState(false);
  const { isIntersecting, ref } = useIntersectionObserver({
    freezeOnceVisible: true,
    threshold: 0.1,
    ...options,
  });
  useEffect(() => {
    if (isIntersecting && src && !imageSrc) {
      setImageSrc(src);
    }
  }, [isIntersecting, src, imageSrc]);
  const handleLoad = useCallback(() => {
    setIsLoaded(true);
    setHasError(false);
  }, []);
  const handleError = useCallback(() => {
    setHasError(true);
    setIsLoaded(false);
  }, []);
  return {
    ref,
    src: imageSrc,
    isLoaded,
    hasError,
    onLoad: handleLoad,
    onError: handleError,
  };
}
/**
 * Hook for infinite scrolling
 */
export function useInfiniteScroll({
  hasNextPage = false,
  isFetchingNextPage = false,
  fetchNextPage,
  threshold = 0.8,
  rootMargin = '100px',
} = {}) {
  const { isIntersecting, ref } = useIntersectionObserver({
    threshold,
    rootMargin,
  });
  useEffect(() => {
    if (
      isIntersecting && 
      hasNextPage && 
      !isFetchingNextPage && 
      fetchNextPage
    ) {
      fetchNextPage();
    }
  }, [isIntersecting, hasNextPage, isFetchingNextPage, fetchNextPage]);
  return {
    ref,
    isIntersecting,
  };
}
/**
 * Hook for tracking element visibility
 */
export function useVisibilityTracker(onVisible, onHidden, options = {}) {
  const [isVisible, setIsVisible] = useState(false);
  const { isIntersecting, ref } = useIntersectionObserver({
    threshold: 0.5,
    onIntersect: (entry, intersecting) => {
      setIsVisible(intersecting);
      if (intersecting && onVisible) {
        onVisible(entry);
      } else if (!intersecting && onHidden) {
        onHidden(entry);
      }
    },
    ...options,
  });
  return {
    ref,
    isVisible,
    isIntersecting,
  };
} 