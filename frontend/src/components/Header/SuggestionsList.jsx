"use client";
import React, { useEffect, useState, useRef, useCallback, useMemo, memo } from "react";
import { createPortal } from "react-dom";
import { Search, MapPin, Loader, X } from "@/icons";
import { useId } from "react";
import styles from "./SearchBar.module.css";
/**
 * SuggestionsList Component - Mobile-First Optimized
 * Production-ready autocomplete with superior mobile UX
 * Features: Touch targets, swipe gestures, progressive loading, haptic feedback
 */
const SuggestionsList = memo(function SuggestionsList({
    suggestions,
    query,
    loading,
    error,
    activeIndex,
    onItemClick,
    onDismiss,
    suggestionsRef,
    highlightMatch,
    type = "product",
    isMobile = false,
    anchorElement = null
}) {
    const listboxId = useId();
    const [buttonRect, setButtonRect] = useState(null);
    const [touchStart, setTouchStart] = useState(null);
    const [isDragging, setIsDragging] = useState(false);
    const [translateY, setTranslateY] = useState(0);
    const [isVisible, setIsVisible] = useState(true);
    // Progressive loading state
    const [displayCount, setDisplayCount] = useState(isMobile ? 8 : 10);
    const [isExpanded, setIsExpanded] = useState(false);
    const containerRef = useRef(null);
    const initialTouchY = useRef(0);
    const lastTouchTime = useRef(0);
    // Progressive loading for mobile performance - with safety check
    const visibleSuggestions = useMemo(() => {
        // Ensure suggestions is always an array
        const safeSuggestions = Array.isArray(suggestions) ? suggestions : [];
        if (!isMobile || isExpanded) return safeSuggestions;
        return safeSuggestions.slice(0, displayCount);
    }, [suggestions, displayCount, isExpanded, isMobile]);
    const hasMoreSuggestions = Array.isArray(suggestions) && suggestions.length > displayCount && !isExpanded;
    // Optimized position calculation for desktop
    useEffect(() => {
        if (anchorElement && !isMobile) {
            const updateButtonRect = () => {
                const rect = anchorElement.getBoundingClientRect();
                setButtonRect(rect);
            };
            updateButtonRect();
            const handleScroll = () => requestAnimationFrame(updateButtonRect);
            window.addEventListener('scroll', handleScroll, { passive: true });
            window.addEventListener('resize', handleScroll, { passive: true });
            return () => {
                window.removeEventListener('scroll', handleScroll);
                window.removeEventListener('resize', handleScroll);
            };
        } else {
            setButtonRect(null);
        }
    }, [anchorElement, isMobile]);
    // Mobile-specific touch handlers with haptic feedback
    const handleTouchStart = useCallback((e) => {
        if (!isMobile) return;
        const touch = e.touches[0];
        setTouchStart({ x: touch.clientX, y: touch.clientY });
        initialTouchY.current = touch.clientY;
        setIsDragging(false);
        setTranslateY(0);
        // Haptic feedback on touch start
        if ('vibrate' in navigator) {
            navigator.vibrate(1);
        }
    }, [isMobile]);
    const handleTouchMove = useCallback((e) => {
        if (!isMobile || !touchStart) return;
        const touch = e.touches[0];
        const deltaY = touch.clientY - touchStart.y;
        // Only handle downward swipes for dismissal
        if (deltaY > 10) {
            setIsDragging(true);
            setTranslateY(Math.max(0, deltaY * 0.6)); // Resistance effect
            // Prevent default scrolling when dragging
            e.preventDefault();
        }
    }, [isMobile, touchStart]);
    const handleTouchEnd = useCallback((e) => {
        if (!isMobile || !touchStart) return;
        const touch = e.changedTouches[0];
        const deltaY = touch.clientY - touchStart.y;
        const deltaTime = Date.now() - lastTouchTime.current;
        const velocity = Math.abs(deltaY) / deltaTime;
        // Dismiss on significant swipe down or fast swipe
        if ((deltaY > 100 && velocity > 0.5) || deltaY > 200) {
            setIsVisible(false);
            // Haptic feedback on dismiss
            if ('vibrate' in navigator) {
                navigator.vibrate([2, 1, 2]);
            }
            setTimeout(() => {
                onDismiss?.();
            }, 200);
        } else {
            // Bounce back animation
            setTranslateY(0);
        }
        setTouchStart(null);
        setIsDragging(false);
        lastTouchTime.current = Date.now();
    }, [isMobile, touchStart, onDismiss]);
    // Keyboard handling for mobile
    useEffect(() => {
        if (!isMobile) return;
        const handleKeyDown = (e) => {
            if (e.key === 'Escape') {
                onDismiss?.();
            }
        };
        document.addEventListener('keydown', handleKeyDown);
        return () => document.removeEventListener('keydown', handleKeyDown);
    }, [isMobile, onDismiss]);
    // Load more suggestions progressively
    const handleLoadMore = useCallback(() => {
        if (hasMoreSuggestions && Array.isArray(suggestions)) {
            setDisplayCount(prev => Math.min(prev + 8, suggestions.length));
            // Haptic feedback
            if ('vibrate' in navigator) {
                navigator.vibrate(3);
            }
        }
    }, [hasMoreSuggestions, suggestions]);
    // Expand all suggestions
    const handleShowAll = useCallback(() => {
        setIsExpanded(true);
        setDisplayCount(Array.isArray(suggestions) ? suggestions.length : 0);
    }, [suggestions]);
    // Don't render if no content
    if (!visibleSuggestions.length && !loading && !error && query.trim() === "") {
        return null;
    }
    // Positioning for both mobile and desktop
    const getDropdownStyle = () => {
        if (!buttonRect) {
            return {
                position: 'fixed',
                top: '70px',
                left: '50%',
                transform: 'translateX(-50%)',
                width: 'min(90vw, 400px)',
                zIndex: 9999999,
                maxHeight: '70vh',
            };
        }
        const baseStyle = {
            position: 'fixed',
            top: buttonRect.bottom + 8,
            left: buttonRect.left,
            width: Math.max(buttonRect.width, isMobile ? '95vw' : 280),
            zIndex: 9999999,
            maxHeight: '70vh',
        };
        if (isMobile) {
            return {
                ...baseStyle,
                left: '2.5vw', // Center on mobile with padding
                width: '95vw',
                transform: `translateY(${translateY}px)`,
                transition: isDragging ? 'none' : 'transform 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
                opacity: isVisible ? 1 : 0,
            };
        }
        return baseStyle;
    };
    const suggestionsList = (
        <div
            ref={containerRef}
            className={`${styles.suggestionsList} ${isMobile ? styles.suggestionsMobile : ''}`}
            style={getDropdownStyle()}
            onTouchStart={handleTouchStart}
            onTouchMove={handleTouchMove}
            onTouchEnd={handleTouchEnd}
        >
            {/* Mobile header with dismiss */}
            {isMobile && (
                <div className={styles.mobileHeader}>
                    <button 
                        className={styles.dismissButton}
                        onClick={onDismiss}
                        aria-label="Close suggestions"
                    >
                        <X size={20} />
                    </button>
                </div>
            )}
            <ul
                id={`${type}-suggestions-${listboxId}`}
                ref={suggestionsRef}
                role="listbox"
                className={styles.suggestionListContent}
                aria-label={`${type === 'product' ? 'Search' : 'Location'} suggestions`}
                tabIndex={-1}
            >
                {/* Loading state */}
                {loading && (
                    <li role="option" className={styles.loadingItem}>
                        <Loader size={20} className={styles.suggestionIcon} />
                        <span>Finding suggestions...</span>
                    </li>
                )}
                {/* Error state */}
                {error && (
                    <li role="option" className={styles.errorItem}>
                        <span>Unable to load suggestions</span>
                    </li>
                )}
                {/* Empty results */}
                {!loading && !error && Array.isArray(suggestions) && suggestions.length === 0 && query.trim() !== "" && (
                    <li className={styles.suggestionEmpty} role="status">
                        <div className={styles.emptyContent}>
                            <Search size={24} className={styles.emptyIcon} />
                            <span>No {type === 'product' ? 'suggestions' : 'locations'} found</span>
                            <small>Try a different search term</small>
                        </div>
                    </li>
                )}
                {/* Suggestion items */}
                {!loading && !error && visibleSuggestions.map((suggestion, index) => {
                    const isActive = index === activeIndex;
                    let displayText;
                    if (type === 'product') {
                        displayText = suggestion.name || suggestion.text || suggestion.label || '';
                    } else if (type === 'location') {
                        displayText = typeof suggestion === 'string' ? suggestion : 
                                    suggestion.suggestedCity || suggestion.name || '';
                    } else {
                        displayText = suggestion.text || suggestion.name || suggestion.label || '';
                    }
                    displayText = String(displayText || '');
                    return (
                        <li
                            key={suggestion.id || suggestion.name || displayText || index}
                            role="option"
                            aria-selected={isActive}
                            className={`${styles.suggestionItem} ${isActive ? styles.active : ''}`}
                            onClick={() => {
                                // Haptic feedback on selection
                                if (isMobile && 'vibrate' in navigator) {
                                    navigator.vibrate(5);
                                }
                                onItemClick(suggestion);
                            }}
                        >
                            <div className={styles.suggestionIcon}>
                                {type === 'location' ? (
                                    <MapPin size={isMobile ? 20 : 16} />
                                ) : (
                                    <Search size={isMobile ? 20 : 16} />
                                )}
                            </div>
                            <div className={styles.suggestionContent}>
                                <span className={styles.suggestionText}>
                                    {highlightMatch && displayText ? 
                                        highlightMatch(displayText, query) : displayText}
                                </span>
                                {type === 'product' && suggestion.category && (
                                    <span className={styles.suggestionCategory}>
                                        in {suggestion.category}
                                    </span>
                                )}
                            </div>
                        </li>
                    );
                })}
            </ul>
            {/* Progressive loading controls for mobile */}
            {isMobile && hasMoreSuggestions && (
                <div className={styles.loadMoreSection}>
                    <button 
                        className={styles.loadMoreButton}
                        onClick={handleLoadMore}
                    >
                        Show {Math.min(8, (Array.isArray(suggestions) ? suggestions.length : 0) - displayCount)} more
                    </button>
                    <button 
                        className={styles.showAllButton}
                        onClick={handleShowAll}
                    >
                        Show all {Array.isArray(suggestions) ? suggestions.length : 0}
                    </button>
                </div>
            )}
        </div>
    );
    // Always use portal for consistent positioning
    if (typeof window !== 'undefined') {
        try {
            return createPortal(suggestionsList, document.body);
        } catch (error) {
            return suggestionsList;
        }
    }
    return suggestionsList;
});
export default SuggestionsList; 