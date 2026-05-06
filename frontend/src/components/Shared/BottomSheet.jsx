/**
 * BottomSheet - Enhanced Mobile-First Bottom Sheet with Gestures
 *
 * Features:
 * - Drag-to-dismiss gesture support
 * - Multiple height presets (peek, half, full)
 * - Smooth animations optimized for mobile
 * - Backdrop blur and safe area support
 * - Accessibility improvements
 * - Integration with existing TouchInteractions
 *
 * Can be used as a drop-in replacement for existing bottom sheets
 */
"use client";
import React, {
    useRef,
    useCallback,
    useEffect,
    useState,
    useMemo,
    forwardRef,
    useImperativeHandle,
    memo
} from 'react';
import PropTypes from 'prop-types';
import {X as CloseIcon} from '@/icons';
// Height presets for different use cases
const HEIGHT_PRESETS = {
    peek: '30vh',
    half: '50vh',
    full: '90vh',
    auto: 'auto'
};
// Animation configuration
const ANIMATION_CONFIG = {
    SPRING_DAMPING: 0.8,
    SPRING_STIFFNESS: 100,
    DRAG_THRESHOLD: 50, // Minimum drag distance to dismiss
    VELOCITY_THRESHOLD: 0.5, // Minimum velocity to dismiss
    ANIMATION_DURATION: 300,
};
const BottomSheet = memo(forwardRef(({
                                    isOpen = false,
                                    onClose,
                                    children,
                                    title,
                                    height = HEIGHT_PRESETS.half,
                                    enableDragToDismiss = true,
                                    enableBackdropDismiss = true,
                                    showCloseButton = true,
                                    showDragHandle = true,
                                    backdropBlur = true,
                                    className = '',
                                    contentClassName = '',
                                    snapPoints = [],
                                    initialSnapPoint = 0,
                                    onSnapPointChange,
                                    preventScrollWhenOpen = true,
                                    zIndex = 1000,
                                    ...props
                                }, ref) => {
    // Refs for gesture handling
    const sheetRef = useRef(null);
    const contentRef = useRef(null);
    const dragStateRef = useRef({
        isDragging: false,
        startY: 0,
        currentY: 0,
        velocity: 0,
        startTime: 0
    });
    // State management
    const [isAnimating, setIsAnimating] = useState(false);
    const [currentHeight, setCurrentHeight] = useState(height);
    const [dragOffset, setDragOffset] = useState(0);
    const [currentSnapPoint, setCurrentSnapPoint] = useState(initialSnapPoint);
    // Calculate effective height based on snap points or height prop
    const effectiveHeight = useMemo(() => {
        if (snapPoints.length > 0) {
            return snapPoints[currentSnapPoint] || HEIGHT_PRESETS.half;
        }
        return height;
    }, [snapPoints, currentSnapPoint, height]);
    // Expose imperative methods
    useImperativeHandle(ref, () => ({
        snapTo: (pointIndex) => {
            if (snapPoints.length > 0 && pointIndex < snapPoints.length) {
                setCurrentSnapPoint(pointIndex);
                if (onSnapPointChange) {
                    onSnapPointChange(pointIndex, snapPoints[pointIndex]);
                }
            }
        },
        close: () => {
            handleClose();
        },
        isOpen: () => isOpen
    }), [snapPoints, onSnapPointChange, isOpen]);
    // Handle body scroll prevention
    useEffect(() => {
        if (preventScrollWhenOpen && isOpen) {
            const originalStyle = window.getComputedStyle(document.body).overflow;
            document.body.style.overflow = 'hidden';
            return () => {
                document.body.style.overflow = originalStyle;
            };
        }
    }, [isOpen, preventScrollWhenOpen]);
    // Handle escape key
    useEffect(() => {
        const handleEscKey = (event) => {
            if (event.key === 'Escape' && isOpen) {
                handleClose();
            }
        };
        if (isOpen) {
            document.addEventListener('keydown', handleEscKey);
            return () => document.removeEventListener('keydown', handleEscKey);
        }
    }, [isOpen]);
    // Close handler with animation
    const handleClose = useCallback(() => {
        if (isAnimating) return;
        setIsAnimating(true);
        setDragOffset(0);
        setTimeout(() => {
            setIsAnimating(false);
            onClose();
        }, ANIMATION_CONFIG.ANIMATION_DURATION);
    }, [isAnimating, onClose]);
    // Calculate velocity for better gesture recognition
    const calculateVelocity = useCallback((currentY, startY, currentTime, startTime) => {
        const distance = currentY - startY;
        const time = currentTime - startTime;
        return time > 0 ? distance / time : 0;
    }, []);
    // Handle touch start for drag gesture
    const handleTouchStart = useCallback((event) => {
        if (!enableDragToDismiss) return;
        const touch = event.touches[0];
        dragStateRef.current = {
            isDragging: true,
            startY: touch.clientY,
            currentY: touch.clientY,
            velocity: 0,
            startTime: Date.now()
        };
    }, [enableDragToDismiss]);
    // Handle touch move for drag gesture
    const handleTouchMove = useCallback((event) => {
        if (!enableDragToDismiss || !dragStateRef.current.isDragging) return;
        const touch = event.touches[0];
        const currentY = touch.clientY;
        const deltaY = currentY - dragStateRef.current.startY;
        // Only allow downward dragging to dismiss
        if (deltaY > 0) {
            dragStateRef.current.currentY = currentY;
            dragStateRef.current.velocity = calculateVelocity(
                currentY,
                dragStateRef.current.startY,
                Date.now(),
                dragStateRef.current.startTime
            );
            setDragOffset(deltaY);
            // Prevent default scrolling when dragging
            event.preventDefault();
        }
    }, [enableDragToDismiss, calculateVelocity]);
    // Handle touch end for drag gesture
    const handleTouchEnd = useCallback(() => {
        if (!enableDragToDismiss || !dragStateRef.current.isDragging) return;
        const {currentY, startY, velocity} = dragStateRef.current;
        const deltaY = currentY - startY;
        // Determine if should dismiss based on distance or velocity
        const shouldDismiss =
            deltaY > ANIMATION_CONFIG.DRAG_THRESHOLD ||
            velocity > ANIMATION_CONFIG.VELOCITY_THRESHOLD;
        if (shouldDismiss) {
            handleClose();
        } else {
            // Snap back to original position
            setDragOffset(0);
        }
        // Reset drag state
        dragStateRef.current.isDragging = false;
    }, [enableDragToDismiss, handleClose]);
    // Handle backdrop click
    const handleBackdropClick = useCallback((event) => {
        if (enableBackdropDismiss && event.target === event.currentTarget) {
            handleClose();
        }
    }, [enableBackdropDismiss, handleClose]);
    // Render nothing if not open and not animating
    if (!isOpen && !isAnimating) return null;
    // Calculate transform based on drag offset
    const transform = useMemo(() => {
        if (dragOffset > 0) {
            // Apply some resistance to the drag
            const resistance = Math.min(dragOffset / 3, 100);
            return `translateY(${resistance}px)`;
        }
        return 'translateY(0)';
    }, [dragOffset]);
    // Calculate opacity based on drag progress
    const backdropOpacity = useMemo(() => {
        if (dragOffset > 0) {
            const progress = Math.min(dragOffset / 200, 1);
            return 1 - progress * 0.6;
        }
        return 1;
    }, [dragOffset]);
    return (
        <div
            className={`bottom-sheet-backdrop ${className}`}
            onClick={handleBackdropClick}
            style={{
                position: 'fixed',
                top: 0,
                left: 0,
                right: 0,
                bottom: 0,
                backgroundColor: `rgba(0, 0, 0, ${0.4 * backdropOpacity})`,
                backdropFilter: backdropBlur ? 'blur(4px)' : 'none',
                WebkitBackdropFilter: backdropBlur ? 'blur(4px)' : 'none',
                zIndex,
                display: 'flex',
                alignItems: 'flex-end',
                opacity: isOpen && !isAnimating ? 1 : 0,
                transition: isAnimating ? `opacity ${ANIMATION_CONFIG.ANIMATION_DURATION}ms ease-out` : 'none',
                paddingBottom: 'env(safe-area-inset-bottom)',
            }}
            {...props}
        >
            <div
                ref={sheetRef}
                className={`bottom-sheet-container ${contentClassName}`}
                onTouchStart={handleTouchStart}
                onTouchMove={handleTouchMove}
                onTouchEnd={handleTouchEnd}
                onClick={(e) => e.stopPropagation()}
                style={{
                    width: '100%',
                    maxHeight: effectiveHeight,
                    backgroundColor: 'white',
                    borderTopLeftRadius: '16px',
                    borderTopRightRadius: '16px',
                    transform: isOpen && !isAnimating ? transform : 'translateY(100%)',
                    transition: dragStateRef.current.isDragging
                        ? 'none'
                        : `transform ${ANIMATION_CONFIG.ANIMATION_DURATION}ms ease-out`,
                    touchAction: 'none',
                    overflow: 'hidden',
                    boxShadow: '0 -4px 20px rgba(0, 0, 0, 0.15)',
                }}
                role="dialog"
                aria-modal="true"
                aria-labelledby={title ? 'bottom-sheet-title' : undefined}
            >
                {/* Drag handle */}
                {showDragHandle && (
                    <div
                        style={{
                            display: 'flex',
                            justifyContent: 'center',
                            paddingTop: '8px',
                            paddingBottom: '8px',
                            cursor: 'grab'
                        }}
                    >
                        <div
                            style={{
                                width: '40px',
                                height: '4px',
                                backgroundColor: '#d1d5db',
                                borderRadius: '2px',
                                transition: 'background-color 0.2s ease'
                            }}
                        />
                    </div>
                )}
                {/* Header */}
                {(title || showCloseButton) && (
                    <div
                        style={{
                            display: 'flex',
                            justifyContent: 'space-between',
                            alignItems: 'center',
                            padding: '16px 20px 8px',
                            borderBottom: '1px solid #f3f4f6'
                        }}
                    >
                        {title && (
                            <h3
                                id="bottom-sheet-title"
                                style={{
                                    margin: 0,
                                    fontSize: '18px',
                                    fontWeight: '600',
                                    color: '#1f2937'
                                }}
                            >
                                {title}
                            </h3>
                        )}
                        {showCloseButton && (
                            <button
                                onClick={handleClose}
                                style={{
                                    background: 'none',
                                    border: 'none',
                                    padding: '8px',
                                    cursor: 'pointer',
                                    borderRadius: '8px',
                                    display: 'flex',
                                    alignItems: 'center',
                                    justifyContent: 'center',
                                    color: '#6b7280',
                                    transition: 'background-color 0.2s ease, color 0.2s ease'
                                }}
                                onMouseEnter={(e) => {
                                    e.target.style.backgroundColor = '#f3f4f6';
                                    e.target.style.color = '#374151';
                                }}
                                onMouseLeave={(e) => {
                                    e.target.style.backgroundColor = 'transparent';
                                    e.target.style.color = '#6b7280';
                                }}
                                aria-label="Close bottom sheet"
                            >
                                <CloseIcon size={20}/>
                            </button>
                        )}
                    </div>
                )}
                {/* Content */}
                <div
                    ref={contentRef}
                    style={{
                        flex: 1,
                        overflowY: 'auto',
                        WebkitOverflowScrolling: 'touch',
                        paddingBottom: 'env(safe-area-inset-bottom)'
                    }}
                >
                    {children}
                </div>
            </div>
        </div>
    );
}));
BottomSheet.displayName = 'BottomSheet';
BottomSheet.propTypes = {
    isOpen: PropTypes.bool,
    onClose: PropTypes.func.isRequired,
    children: PropTypes.node.isRequired,
    title: PropTypes.string,
    height: PropTypes.string,
    enableDragToDismiss: PropTypes.bool,
    enableBackdropDismiss: PropTypes.bool,
    showCloseButton: PropTypes.bool,
    showDragHandle: PropTypes.bool,
    backdropBlur: PropTypes.bool,
    className: PropTypes.string,
    contentClassName: PropTypes.string,
    snapPoints: PropTypes.arrayOf(PropTypes.string),
    initialSnapPoint: PropTypes.number,
    onSnapPointChange: PropTypes.func,
    preventScrollWhenOpen: PropTypes.bool,
    zIndex: PropTypes.number,
};
export default BottomSheet;