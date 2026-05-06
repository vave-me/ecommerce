/**
 * TouchInteractions - Mobile-First Touch Gesture Component
 *
 * Provides reusable touch interactions including:
 * - Swipe gestures (left, right, up, down)
 * - Pull-to-refresh functionality
 * - Touch feedback and haptics
 * - Long press detection
 *
 * Designed for high performance and mobile optimization
 */
"use client";
import React, { useRef, useCallback, useState, useEffect, memo } from 'react';
import PropTypes from 'prop-types';
// Touch interaction thresholds (optimized for mobile)
const TOUCH_THRESHOLDS = {
    SWIPE_MIN_DISTANCE: 50,      // Minimum distance for swipe recognition
    SWIPE_MAX_TIME: 500,         // Maximum time for swipe (ms)
    PULL_REFRESH_DISTANCE: 80,   // Distance to trigger refresh
    LONG_PRESS_DURATION: 500,    // Long press duration (ms)
    VELOCITY_THRESHOLD: 0.3,     // Minimum velocity for swipe
};
const TouchInteractions = memo(function TouchInteractions({
                               children,
                               onSwipeLeft,
                               onSwipeRight,
                               onSwipeUp,
                               onSwipeDown,
                               onPullRefresh,
                               onLongPress,
                               onTouchStart,
                               onTouchEnd,
                               enablePullToRefresh = false,
                               enableSwipeGestures = true,
                               enableLongPress = false,
                               pullRefreshThreshold = TOUCH_THRESHOLDS.PULL_REFRESH_DISTANCE,
                               className = '',
                               refreshing = false,
                               disabled = false,
                               hapticFeedback = true,
                               ...props
                           }) {
    // Touch state management
    const touchStartRef = useRef(null);
    const touchTimeRef = useRef(null);
    const longPressTimerRef = useRef(null);
    const containerRef = useRef(null);
    // Component state
    const [isDragging, setIsDragging] = useState(false);
    const [pullDistance, setPullDistance] = useState(0);
    const [isPulling, setIsPulling] = useState(false);
    const [touchFeedback, setTouchFeedback] = useState(false);
    // Cleanup long press timer
    useEffect(() => {
        return () => {
            if (longPressTimerRef.current) {
                clearTimeout(longPressTimerRef.current);
            }
        };
    }, []);
    // Haptic feedback helper
    const triggerHaptic = useCallback((type = 'light') => {
        if (!hapticFeedback || !navigator.vibrate) return;
        try {
            switch (type) {
                case 'light':
                    navigator.vibrate(10);
                    break;
                case 'medium':
                    navigator.vibrate(20);
                    break;
                case 'heavy':
                    navigator.vibrate([10, 10, 10]);
                    break;
                default:
                    navigator.vibrate(10);
            }
        } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
    }, [hapticFeedback]);
    // Calculate touch velocity
    const calculateVelocity = useCallback((startTouch, endTouch, duration) => {
        const distance = Math.sqrt(
            Math.pow(endTouch.clientX - startTouch.clientX, 2) +
            Math.pow(endTouch.clientY - startTouch.clientY, 2)
        );
        return distance / duration;
    }, []);
    // Handle touch start
    const handleTouchStart = useCallback((event) => {
        if (disabled) return;
        const touch = event.touches[0];
        touchStartRef.current = {
            x: touch.clientX,
            y: touch.clientY,
            timestamp: Date.now(),
        };
        touchTimeRef.current = Date.now();
        setTouchFeedback(true);
        setIsDragging(false);
        // Start long press timer if enabled
        if (enableLongPress && onLongPress) {
            longPressTimerRef.current = setTimeout(() => {
                triggerHaptic('medium');
                onLongPress(event);
            }, TOUCH_THRESHOLDS.LONG_PRESS_DURATION);
        }
        // Call external touch start handler
        if (onTouchStart) {
            onTouchStart(event);
        }
    }, [disabled, enableLongPress, onLongPress, onTouchStart, triggerHaptic]);
    // Handle touch move
    const handleTouchMove = useCallback((event) => {
        if (disabled || !touchStartRef.current) return;
        const touch = event.touches[0];
        const deltaX = touch.clientX - touchStartRef.current.x;
        const deltaY = touch.clientY - touchStartRef.current.y;
        const distance = Math.sqrt(deltaX * deltaX + deltaY * deltaY);
        // Cancel long press if touch moves significantly
        if (distance > 10 && longPressTimerRef.current) {
            clearTimeout(longPressTimerRef.current);
            longPressTimerRef.current = null;
        }
        // Handle pull-to-refresh
        if (enablePullToRefresh && deltaY > 0 && !isDragging) {
            event.preventDefault();
            setIsPulling(true);
            setPullDistance(Math.min(deltaY, pullRefreshThreshold * 1.5));
            // Provide haptic feedback when reaching threshold
            if (deltaY >= pullRefreshThreshold && !refreshing) {
                triggerHaptic('light');
            }
        }
        // Set dragging state for other gestures
        if (distance > 10) {
            setIsDragging(true);
        }
    }, [disabled, enablePullToRefresh, isDragging, pullRefreshThreshold, refreshing, triggerHaptic]);
    // Handle touch end
    const handleTouchEnd = useCallback((event) => {
        if (disabled || !touchStartRef.current) return;
        const touch = event.changedTouches[0];
        const deltaX = touch.clientX - touchStartRef.current.x;
        const deltaY = touch.clientY - touchStartRef.current.y;
        const duration = Date.now() - touchTimeRef.current;
        const velocity = calculateVelocity(touchStartRef.current, touch, duration);
        // Clear long press timer
        if (longPressTimerRef.current) {
            clearTimeout(longPressTimerRef.current);
            longPressTimerRef.current = null;
        }
        // Reset touch feedback
        setTouchFeedback(false);
        // Handle pull-to-refresh
        if (isPulling) {
            if (pullDistance >= pullRefreshThreshold && onPullRefresh && !refreshing) {
                triggerHaptic('medium');
                onPullRefresh();
            }
            setIsPulling(false);
            setPullDistance(0);
        }
        // Handle swipe gestures
        if (enableSwipeGestures && isDragging && velocity > TOUCH_THRESHOLDS.VELOCITY_THRESHOLD) {
            const absX = Math.abs(deltaX);
            const absY = Math.abs(deltaY);
            const minDistance = TOUCH_THRESHOLDS.SWIPE_MIN_DISTANCE;
            if (duration < TOUCH_THRESHOLDS.SWIPE_MAX_TIME) {
                if (absX > absY && absX > minDistance) {
                    // Horizontal swipe
                    triggerHaptic('light');
                    if (deltaX > 0 && onSwipeRight) {
                        onSwipeRight(event, { distance: absX, velocity });
                    } else if (deltaX < 0 && onSwipeLeft) {
                        onSwipeLeft(event, { distance: absX, velocity });
                    }
                } else if (absY > absX && absY > minDistance) {
                    // Vertical swipe
                    triggerHaptic('light');
                    if (deltaY > 0 && onSwipeDown) {
                        onSwipeDown(event, { distance: absY, velocity });
                    } else if (deltaY < 0 && onSwipeUp) {
                        onSwipeUp(event, { distance: absY, velocity });
                    }
                }
            }
        }
        // Reset state
        setIsDragging(false);
        touchStartRef.current = null;
        // Call external touch end handler
        if (onTouchEnd) {
            onTouchEnd(event);
        }
    }, [
        disabled, calculateVelocity, isPulling, pullDistance, pullRefreshThreshold,
        onPullRefresh, refreshing, triggerHaptic, enableSwipeGestures, isDragging,
        onSwipeRight, onSwipeLeft, onSwipeDown, onSwipeUp, onTouchEnd
    ]);
    // Render pull-to-refresh indicator
    const renderPullRefreshIndicator = () => {
        if (!enablePullToRefresh || (!isPulling && !refreshing)) return null;
        const progress = Math.min(pullDistance / pullRefreshThreshold, 1);
        const rotation = refreshing ? 360 : progress * 180;
        return (
            <div
                style={{
                    position: 'absolute',
                    top: -40,
                    left: '50%',
                    transform: 'translateX(-50%)',
                    opacity: isPulling || refreshing ? 1 : 0,
                    transition: 'opacity 0.2s ease',
                    zIndex: 1000,
                }}
            >
                <div
                    style={{
                        width: 24,
                        height: 24,
                        borderRadius: '50%',
                        border: '2px solid #e5e7eb',
                        borderTopColor: '#2980b9',
                        transform: `rotate(${rotation}deg)`,
                        transition: refreshing ? 'none' : 'transform 0.1s ease',
                        animation: refreshing ? 'spin 1s linear infinite' : 'none',
                    }}
                />
                <style jsx>{`
          @keyframes spin {
            to {
              transform: rotate(360deg);
            }
          }
        `}</style>
            </div>
        );
    };
    return (
        <div
            ref={containerRef}
            className={`touch-interactions ${className}`}
            onTouchStart={handleTouchStart}
            onTouchMove={handleTouchMove}
            onTouchEnd={handleTouchEnd}
            style={{
                position: 'relative',
                transform: isPulling ? `translateY(${Math.min(pullDistance / 2, 20)}px)` : 'none',
                transition: isPulling ? 'none' : 'transform 0.2s ease',
                touchAction: enablePullToRefresh ? 'pan-x pan-down' : 'manipulation',
                WebkitTouchCallout: 'none',
                WebkitUserSelect: 'none',
                userSelect: 'none',
                opacity: touchFeedback && !isDragging ? 0.9 : 1,
                backgroundColor: touchFeedback && !isDragging ? 'rgba(79, 70, 229, 0.05)' : 'transparent',
                borderRadius: '8px',
                transition: 'opacity 0.1s ease, background-color 0.1s ease',
                ...props.style,
            }}
            {...props}
        >
            {renderPullRefreshIndicator()}
            {children}
        </div>
    );
});
TouchInteractions.propTypes = {
    children: PropTypes.node.isRequired,
    onSwipeLeft: PropTypes.func,
    onSwipeRight: PropTypes.func,
    onSwipeUp: PropTypes.func,
    onSwipeDown: PropTypes.func,
    onPullRefresh: PropTypes.func,
    onLongPress: PropTypes.func,
    onTouchStart: PropTypes.func,
    onTouchEnd: PropTypes.func,
    enablePullToRefresh: PropTypes.bool,
    enableSwipeGestures: PropTypes.bool,
    enableLongPress: PropTypes.bool,
    pullRefreshThreshold: PropTypes.number,
    className: PropTypes.string,
    refreshing: PropTypes.bool,
    disabled: PropTypes.bool,
    hapticFeedback: PropTypes.bool,
};
export default TouchInteractions;