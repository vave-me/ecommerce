import { useState, useEffect, useCallback, useRef } from 'react';
/**
 * Custom hook for scroll-based navigation hiding
 * Uses modern Intersection Observer API and optimized scroll detection
 * 
 * @param {number} threshold - Scroll threshold to start hiding/showing (default: 100px)
 * @param {number} scrollThreshold - Minimum scroll distance to trigger hide/show (default: 50px)
 * @returns {Object} { isNavVisible, isScrollingDown, scrollY }
 */
export const useScrollNavigation = (threshold = 100, scrollThreshold = 50) => {
    const [isNavVisible, setIsNavVisible] = useState(true);
    const [isScrollingDown, setIsScrollingDown] = useState(false);
    const [scrollY, setScrollY] = useState(0);
    const lastScrollY = useRef(0);
    const ticking = useRef(false);
    const observerRef = useRef(null);
    const sentinelRef = useRef(null);
    // Optimized scroll handler using requestAnimationFrame
    const updateScrollPosition = useCallback(() => {
        const currentScrollY = window.scrollY;
        setScrollY(currentScrollY);
        // Determine scroll direction
        const scrollingDown = currentScrollY > lastScrollY.current;
        setIsScrollingDown(scrollingDown);
        // Only update visibility if scroll distance is significant
        const scrollDelta = Math.abs(currentScrollY - lastScrollY.current);
        if (scrollDelta > scrollThreshold) {
            if (currentScrollY < threshold) {
                // Always show navigation when near top
                setIsNavVisible(true);
            } else {
                // Hide when scrolling down, show when scrolling up
                setIsNavVisible(!scrollingDown);
            }
        }
        lastScrollY.current = currentScrollY;
        ticking.current = false;
    }, [threshold, scrollThreshold]);
    // Throttled scroll handler
    const handleScroll = useCallback(() => {
        if (!ticking.current) {
            requestAnimationFrame(updateScrollPosition);
            ticking.current = true;
        }
    }, [updateScrollPosition]);
    // Intersection Observer setup for performance optimization
    useEffect(() => {
        // Create sentinel element for intersection observer
        const sentinel = document.createElement('div');
        sentinel.style.position = 'absolute';
        sentinel.style.top = `${threshold}px`;
        sentinel.style.height = '1px';
        sentinel.style.width = '1px';
        sentinel.style.opacity = '0';
        sentinel.style.pointerEvents = 'none';
        sentinel.setAttribute('data-scroll-sentinel', 'true');
        document.body.appendChild(sentinel);
        sentinelRef.current = sentinel;
        // Intersection Observer for threshold detection
        const observer = new IntersectionObserver(
            (entries) => {
                entries.forEach((entry) => {
                    if (entry.isIntersecting) {
                        // Near top of page
                        setIsNavVisible(true);
                    }
                });
            },
            {
                root: null,
                rootMargin: '0px',
                threshold: 0
            }
        );
        observer.observe(sentinel);
        observerRef.current = observer;
        return () => {
            observer.unobserve(sentinel);
            observer.disconnect();
            document.body.removeChild(sentinel);
        };
    }, [threshold]);
    // Scroll event listener
    useEffect(() => {
        // Passive listener for better performance
        window.addEventListener('scroll', handleScroll, { passive: true });
        // Set initial scroll position
        setScrollY(window.scrollY);
        lastScrollY.current = window.scrollY;
        return () => {
            window.removeEventListener('scroll', handleScroll);
        };
    }, [handleScroll]);
    // Handle visibility changes (browser tabs, etc.)
    useEffect(() => {
        const handleVisibilityChange = () => {
            if (document.visibilityState === 'visible') {
                // Reset scroll position when returning to tab
                const currentScrollY = window.scrollY;
                setScrollY(currentScrollY);
                lastScrollY.current = currentScrollY;
            }
        };
        document.addEventListener('visibilitychange', handleVisibilityChange);
        return () => document.removeEventListener('visibilitychange', handleVisibilityChange);
    }, []);
    return {
        isNavVisible,
        isScrollingDown,
        scrollY
    };
};
export default useScrollNavigation; 