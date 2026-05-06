// File: src/components/MediaGallery/MediaGalleryMobile.jsx
"use client";
import { ChevronLeft, ChevronRight, FaPlay } from '../../../icons';
import React, { useState, useCallback, useMemo, lazy, Suspense, memo } from "react";
import PropTypes from "prop-types";
import {useSwipeable} from "react-swipeable";
import Spinner from "../UI/Spinner";
import MediaDisplay from "./MediaDisplay";
import styles from "./MediaGalleryMobile.module.css";
// Lazy load the lightbox for better performance
const LightboxViewer = lazy(() => import("./LightboxViewer").then(module => ({ default: module.default })));
/**
 * Helper to safely wrap index to stay within bounds
 */
function wrapIndex(index, length) {
    if (length === 0) return 0;
    return ((index % length) + length) % length;
}
/**
 * MediaGalleryMobile - A mobile-optimized media gallery with thumbnails and swipe support
 *
 * Features:
 * - Thumbnail navigation column
 * - Swipe gestures for mobile
 * - Fullscreen lightbox view
 * - Video support with play/pause controls
 * - Accessible controls and keyboard navigation
 * - Responsive design
 *
 * @param {Object} props - Component props
 * @param {Array} props.gallery - Array of media objects with type, src, alt
 * @param {number} props.currentMediaIndex - Current active media index
 * @param {Function} props.setCurrentMediaIndex - Function to update current media index
 * @param {string} props.itemTitle - Title of the parent item (for accessibility)
 */
function MediaGalleryMobile({
                                gallery,
                                currentMediaIndex,
                                setCurrentMediaIndex,
                                itemTitle,
                            }) {
    const [isLightboxOpen, setIsLightboxOpen] = useState(false);
    const [activeVideoIndex, setActiveVideoIndex] = useState(null);
    // Filter out invalid media items
    const validGallery = useMemo(() => {
        if (!Array.isArray(gallery)) return [];
        return gallery.filter((item) => item && item.src);
    }, [gallery]);
    // Get current media with safety checks
    const currentMedia = useMemo(() => {
        if (!validGallery.length) return null;
        const safeIndex = wrapIndex(currentMediaIndex, validGallery.length);
        return validGallery[safeIndex];
    }, [currentMediaIndex, validGallery]);
    // Check if gallery has only one item
    const isSingleMedia = validGallery.length <= 1;
    // Navigation handlers
    const handleNext = useCallback(() => {
        if (isSingleMedia) return;
        setCurrentMediaIndex((prev) => wrapIndex(prev + 1, validGallery.length));
        setActiveVideoIndex(null);
    }, [isSingleMedia, setCurrentMediaIndex, validGallery.length]);
    const handlePrev = useCallback(() => {
        if (isSingleMedia) return;
        setCurrentMediaIndex((prev) => wrapIndex(prev - 1, validGallery.length));
        setActiveVideoIndex(null);
    }, [isSingleMedia, setCurrentMediaIndex, validGallery.length]);
    // Handle thumbnail selection
    const handleThumbnailClick = useCallback(
        (index) => {
            setCurrentMediaIndex(index);
            setActiveVideoIndex(null);
        },
        [setCurrentMediaIndex]
    );
    // Lightbox handlers
    const handleOpenLightbox = useCallback(() => {
        if (currentMedia?.type === "image") {
            setIsLightboxOpen(true);
            setActiveVideoIndex(null);
        }
    }, [currentMedia]);
    const handleCloseLightbox = useCallback(() => {
        setIsLightboxOpen(false);
    }, []);
    // Set up swipe handlers for mobile
    const swipeHandlers = useSwipeable({
        onSwipedLeft: handleNext,
        onSwipedRight: handlePrev,
        preventDefaultTouchmoveEvent: !isSingleMedia,
        trackMouse: !isSingleMedia,
        // Disable if only one media item
        disabled: isSingleMedia,
    });
    // Video play/pause handler
    const handleVideoPlayToggle = useCallback(
        (isPlaying) => {
            setActiveVideoIndex(isPlaying ? currentMediaIndex : null);
        },
        [currentMediaIndex]
    );
    // Handle keyboard navigation
    const handleKeyDown = useCallback(
        (e) => {
            if (e.key === "ArrowRight" || e.key === "ArrowDown") {
                handleNext();
                e.preventDefault();
            } else if (e.key === "ArrowLeft" || e.key === "ArrowUp") {
                handlePrev();
                e.preventDefault();
            }
        },
        [handleNext, handlePrev]
    );
    // If no gallery, display fallback message
    if (!validGallery.length) {
        return (
            <div className={styles.fallbackMedia} aria-label="No media available">
                No media available
            </div>
        );
    }
    return (
        <div
            className={styles.galleryContainer}
            role="region"
            aria-label={`Media gallery for ${itemTitle}`}
            onKeyDown={handleKeyDown}
            tabIndex={0}
        >
            {/* Thumbnail column */}
            <div className={styles.thumbnailColumn} aria-label="Thumbnail navigation">
                {validGallery.map((item, index) => {
                    const isActive = index === currentMediaIndex;
                    const thumbnailSrc =
                        item.type === "video" ? (item.poster || "/images/video-icon.webp") : item.src;
                    const altText = item.alt || `${itemTitle} media ${index + 1}`;
                    return (
                        <button
                            key={`thumbnail-${index}`}
                            type="button"
                            className={`${styles.thumbnailButton} ${
                                isActive ? styles.thumbnailButtonActive : ""
                            }`}
                            onClick={() => handleThumbnailClick(index)}
                            aria-label={`View ${item.type} ${index + 1} of ${validGallery.length}`}
                            aria-current={isActive ? "true" : "false"}
                            title={altText}
                        >
                            <img
                                className={styles.thumbImage}
                                src={thumbnailSrc}
                                alt=""
                                aria-hidden="true"
                                loading="lazy"
                                onError={(e) => {
                                    e.currentTarget.src = "/images/video-icon.webp";
                                }}
                            />
                            {item.type === "video" && (
                                <div className={styles.videoPlayOverlay} aria-hidden="true">
                                    <FaPlay/>
                                </div>
                            )}
                        </button>
                    );
                })}
            </div>
            {/* Main media display area */}
            <div
                className={styles.mainMediaArea}
                {...swipeHandlers}
                aria-live="polite"
            >
                {/* Navigation buttons (only show if more than one media) */}
                {!isSingleMedia && (
                    <>
                        <button
                            className={`${styles.navButton} ${styles.navButtonPrev}`}
                            onClick={handlePrev}
                            aria-label="Previous media"
                            type="button"
                        >
                            <ChevronLeft/>
                        </button>
                        <button
                            className={`${styles.navButton} ${styles.navButtonNext}`}
                            onClick={handleNext}
                            aria-label="Next media"
                            type="button"
                        >
                            <ChevronRight/>
                        </button>
                    </>
                )}
                {/* Media content */}
                <div className={styles.mediaWrapper}>
                    <MediaDisplay
                        media={currentMedia}
                        onImageClick={handleOpenLightbox}
                        isPlaying={activeVideoIndex === currentMediaIndex}
                        setIsPlaying={handleVideoPlayToggle}
                        itemTitle={itemTitle}
                    />
                </div>
                {/* Lightbox (lazy loaded) */}
                {isLightboxOpen && (
                    <Suspense
                        fallback={
                            <div className={styles.lightboxFallback}>
                                <Spinner/>
                            </div>
                        }
                    >
                        <LightboxViewer
                            slides={validGallery.map((item, i) => ({
                                src: item.src,
                                type: item.type,
                                alt: item.alt || `${itemTitle} media ${i + 1}`,
                            }))}
                            currentIndex={currentMediaIndex}
                            onClose={handleCloseLightbox}
                            onNext={handleNext}
                            onPrev={handlePrev}
                            title={itemTitle}
                        />
                    </Suspense>
                )}
            </div>
        </div>
    );
}
MediaGalleryMobile.propTypes = {
    gallery: PropTypes.arrayOf(
        PropTypes.shape({
            type: PropTypes.oneOf(["image", "video"]).isRequired,
            src: PropTypes.string.isRequired,
            alt: PropTypes.string,
            poster: PropTypes.string,
        })
    ),
    currentMediaIndex: PropTypes.number,
    setCurrentMediaIndex: PropTypes.func.isRequired,
    itemTitle: PropTypes.string,
};
MediaGalleryMobile.defaultProps = {
    gallery: [],
    currentMediaIndex: 0,
    itemTitle: "Item",
};
export default memo(MediaGalleryMobile);