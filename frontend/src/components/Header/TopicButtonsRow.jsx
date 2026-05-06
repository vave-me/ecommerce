"use client";
import React, { forwardRef, useEffect, useRef, useState, useCallback } from 'react';
import PropTypes from 'prop-types';
import { useTranslations } from 'next-intl';
import TopicButton from './TopicButton';
import styles from './TopicButtonsRow.module.css';
/**
 * TopicButtonsRow Component
 * Scrollable row of topic buttons with navigation behavior and auto-scroll
 * Extracted from SelectTopic to improve component modularity
 */
const TopicButtonsRow = forwardRef(({ 
    translatedTopics, 
    activeTopicValue, 
    openTopicValue, 
    topicButtonRefs, 
    onTopicToggle,
    enableAutoScroll = true,
    autoScrollSpeed = 0.5, // pixels per frame
    autoScrollPauseDuration = 2000, // ms to pause at each end
    autoScrollDisableOnInteraction = true,
    enablePerformanceMonitoring = false, // Enable for development/debugging
    performanceThreshold = 50 // Minimum FPS threshold
}, ref) => {
    const t = useTranslations('Topics');
    const animationRef = useRef(null);
    const [isAutoScrollActive, setIsAutoScrollActive] = useState(false);
    const [autoScrollDirection, setAutoScrollDirection] = useState(1); // 1 for right, -1 for left
    const [autoScrollPaused, setAutoScrollPaused] = useState(false);
    const [userInteracted, setUserInteracted] = useState(false);
    const [hasOverflow, setHasOverflow] = useState(false);
    const pauseTimeoutRef = useRef(null);
    const lastTimestamp = useRef(0);
    // Performance monitoring
    const performanceDataRef = useRef({
        frameCount: 0,
        lastFpsCheck: 0,
        lowFpsWarningShown: false
    });
    // Check if scrolling is needed (content overflows)
    const checkScrollNeeded = useCallback(() => {
        if (!ref?.current) return false;
        const container = ref.current;
        const overflow = container.scrollWidth > container.clientWidth;
        setHasOverflow(overflow);
        return overflow;
    }, [ref]);
    // Update overflow state on resize - using optimized resize handling
    useEffect(() => {
        // Check on mount and when dependencies change
        checkScrollNeeded();
        // Use a single optimized resize listener that debounces calls
        let resizeTimeoutId = null;
        const handleResize = () => {
            if (resizeTimeoutId) clearTimeout(resizeTimeoutId);
            resizeTimeoutId = setTimeout(checkScrollNeeded, 100); // Debounce resize
        };
        window.addEventListener('resize', handleResize, { passive: true });
        return () => {
            if (resizeTimeoutId) clearTimeout(resizeTimeoutId);
            window.removeEventListener('resize', handleResize);
        };
    }, [checkScrollNeeded]);
    // Generate CSS classes for scroll row
    const getScrollRowClasses = useCallback(() => {
        let classes = styles.scrollRow;
        if (hasOverflow) {
            classes += ` ${styles.hasOverflow}`;
        }
        if (isAutoScrollActive) {
            classes += ` ${styles.autoScrolling}`;
        } else if (hasOverflow && userInteracted) {
            // Show enhanced scrollbar when user has interacted and auto-scroll is off
            classes += ` ${styles.showScrollbar}`;
        }
        return classes;
    }, [hasOverflow, isAutoScrollActive, userInteracted]);
    // Performance monitoring function
    const monitorPerformance = useCallback((timestamp) => {
        if (!enablePerformanceMonitoring) return;
        const perfData = performanceDataRef.current;
        perfData.frameCount++;
        // Check FPS every second
        if (timestamp - perfData.lastFpsCheck >= 1000) {
            const fps = perfData.frameCount;
            perfData.frameCount = 0;
            perfData.lastFpsCheck = timestamp;
            if (fps < performanceThreshold && !perfData.lowFpsWarningShown) {
                perfData.lowFpsWarningShown = true;
            }
        }
    }, [enablePerformanceMonitoring, performanceThreshold]);
    // Auto-scroll animation function
    const performAutoScroll = useCallback((timestamp) => {
        if (!ref?.current || !isAutoScrollActive || autoScrollPaused || userInteracted) {
            return;
        }
        const container = ref.current;
        const scrollNeeded = checkScrollNeeded();
        if (!scrollNeeded) {
            setIsAutoScrollActive(false);
            return;
        }
        // Throttle animation to ~60fps for performance
        if (timestamp - lastTimestamp.current < 16) {
            animationRef.current = requestAnimationFrame(performAutoScroll);
            return;
        }
        lastTimestamp.current = timestamp;
        // Monitor performance
        monitorPerformance(timestamp);
        const maxScrollLeft = container.scrollWidth - container.clientWidth;
        const currentScroll = container.scrollLeft;
        // Calculate next scroll position
        const nextScroll = currentScroll + (autoScrollSpeed * autoScrollDirection);
        // Check boundaries and reverse direction with pause
        if (nextScroll >= maxScrollLeft) {
            container.scrollLeft = maxScrollLeft;
            setAutoScrollDirection(-1);
            setAutoScrollPaused(true);
            pauseTimeoutRef.current = setTimeout(() => {
                setAutoScrollPaused(false);
            }, autoScrollPauseDuration);
        } else if (nextScroll <= 0) {
            container.scrollLeft = 0;
            setAutoScrollDirection(1);
            setAutoScrollPaused(true);
            pauseTimeoutRef.current = setTimeout(() => {
                setAutoScrollPaused(false);
            }, autoScrollPauseDuration);
        } else {
            container.scrollLeft = nextScroll;
        }
        // Continue animation
        animationRef.current = requestAnimationFrame(performAutoScroll);
    }, [ref, isAutoScrollActive, autoScrollPaused, userInteracted, autoScrollSpeed, autoScrollDirection, autoScrollPauseDuration, checkScrollNeeded, monitorPerformance]);
    // Start auto-scroll
    const startAutoScroll = useCallback(() => {
        if (!enableAutoScroll || !checkScrollNeeded()) return;
        setIsAutoScrollActive(true);
        setUserInteracted(false);
        animationRef.current = requestAnimationFrame(performAutoScroll);
    }, [enableAutoScroll, checkScrollNeeded, performAutoScroll]);
    // Stop auto-scroll
    const stopAutoScroll = useCallback(() => {
        setIsAutoScrollActive(false);
        if (animationRef.current) {
            cancelAnimationFrame(animationRef.current);
            animationRef.current = null;
        }
        if (pauseTimeoutRef.current) {
            clearTimeout(pauseTimeoutRef.current);
            pauseTimeoutRef.current = null;
        }
    }, []);
    // Handle user interaction
    const handleUserInteraction = useCallback(() => {
        if (autoScrollDisableOnInteraction) {
            setUserInteracted(true);
            stopAutoScroll();
        }
    }, [autoScrollDisableOnInteraction, stopAutoScroll]);
    // Initialize auto-scroll when component mounts and content changes
    useEffect(() => {
        if (!enableAutoScroll) return;
        // Delay start to allow layout to settle
        const initTimer = setTimeout(() => {
            if (checkScrollNeeded() && !userInteracted) {
                startAutoScroll();
            }
        }, 1000); // 1 second delay before starting
        return () => {
            clearTimeout(initTimer);
            stopAutoScroll();
        };
    }, [enableAutoScroll, translatedTopics.length, checkScrollNeeded, userInteracted, startAutoScroll, stopAutoScroll]);
    // Cleanup on unmount
    useEffect(() => {
        return () => {
            stopAutoScroll();
        };
    }, [stopAutoScroll]);
    // Handle visibility change (pause when tab is not visible)
    useEffect(() => {
        const handleVisibilityChange = () => {
            if (document.hidden) {
                stopAutoScroll();
            } else if (enableAutoScroll && !userInteracted && checkScrollNeeded()) {
                setTimeout(startAutoScroll, 500);
            }
        };
        document.addEventListener('visibilitychange', handleVisibilityChange);
        return () => document.removeEventListener('visibilitychange', handleVisibilityChange);
    }, [enableAutoScroll, userInteracted, checkScrollNeeded, startAutoScroll, stopAutoScroll]);
    return (
        <div
            ref={ref}
            className={getScrollRowClasses()}
            role="navigation"
            aria-label={t('topicsAriaLabel')}
            onMouseEnter={handleUserInteraction}
            onTouchStart={handleUserInteraction}
            onScroll={handleUserInteraction}
            onWheel={handleUserInteraction}
        >
            {translatedTopics.map((topic) => {
                const isActive = topic.value === activeTopicValue;
                const isOpen = topic.value === openTopicValue;
                const dropdownId = `topic-dropdown-${topic.value}`;
                const handleClick = (e) => {
                    handleUserInteraction(); // Stop auto-scroll on button click
                    onTopicToggle(topic.value, e.currentTarget);
                };
                const handleKeyDown = (e) => {
                    if (e.key === 'Enter' || e.key === ' ') {
                        handleUserInteraction(); // Stop auto-scroll on keyboard interaction
                        onTopicToggle(topic.value, e.currentTarget);
                    }
                };
                return (
                    <TopicButton
                        key={topic.value}
                        ref={(el) => (topicButtonRefs.current[topic.value] = el)}
                        topic={topic}
                        isActive={isActive}
                        isOpen={isOpen}
                        onClick={handleClick}
                        onKeyDown={handleKeyDown}
                        dropdownId={dropdownId}
                    />
                );
            })}
        </div>
    );
});
TopicButtonsRow.displayName = 'TopicButtonsRow';
TopicButtonsRow.propTypes = {
    translatedTopics: PropTypes.arrayOf(PropTypes.shape({
        value: PropTypes.string.isRequired,
        label: PropTypes.string.isRequired,
        badge: PropTypes.string,
    })).isRequired,
    activeTopicValue: PropTypes.string.isRequired,
    openTopicValue: PropTypes.string,
    topicButtonRefs: PropTypes.object.isRequired,
    onTopicToggle: PropTypes.func.isRequired,
    enableAutoScroll: PropTypes.bool,
    autoScrollSpeed: PropTypes.number,
    autoScrollPauseDuration: PropTypes.number,
    autoScrollDisableOnInteraction: PropTypes.bool,
    enablePerformanceMonitoring: PropTypes.bool,
    performanceThreshold: PropTypes.number,
};
export default TopicButtonsRow; 