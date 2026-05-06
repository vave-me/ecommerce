// File: src/components/MediaGallery/MediaDisplay.jsx
"use client";
import { FaExpand, FaPause, FaPlay } from '../../../utils/iconImports';
import React, { useState, useEffect, useRef, useCallback, memo } from "react";
import PropTypes from "prop-types";
import styles from "./MediaDisplay.module.css";
/**
 * MediaDisplay - Renders either an image or video with appropriate controls
 *
 * Features:
 * - Lazy video loading (only mounts video element when user plays)
 * - Image lightbox support
 * - Play/pause controls for videos
 * - Accessible controls and keyboard navigation
 * - Fallback handling for missing media
 *
 * @param {Object} props - Component props
 * @param {Object} props.media - Media object with type, src, alt
 * @param {Function} props.onImageClick - Handler for image click (lightbox)
 * @param {boolean} props.isPlaying - Whether video is currently playing
 * @param {Function} props.setIsPlaying - Handler to set playing state
 * @param {string} props.itemTitle - Title of the parent item (for accessibility)
 */
function MediaDisplay({
                          media,
                          onImageClick,
                          placeholderImage,
                          isPlaying,
                          setIsPlaying,
                          itemTitle,
                      }) {
    const videoRef = useRef(null);
    const [videoMounted, setVideoMounted] = useState(false);
    const [videoError, setVideoError] = useState(false);
    // Mount video element when play is requested
    useEffect(() => {
        if (isPlaying && !videoMounted && media?.type === "video") {
            setVideoMounted(true);
        }
    }, [isPlaying, videoMounted, media]);
    // Control video playback based on isPlaying prop
    useEffect(() => {
        if (!videoRef.current || media?.type !== "video") return;
        if (isPlaying) {
            videoRef.current
                .play()
                .catch((err) => {
                    setIsPlaying(false);
                    setVideoError(true);
                });
        } else if (videoRef.current) {
            videoRef.current.pause();
        }
    }, [isPlaying, media, setIsPlaying]);
    // Handle image click for lightbox
    const handleImageClick = useCallback(() => {
        if (onImageClick && media?.type === "image") {
            onImageClick();
        }
    }, [onImageClick, media]);
    // Handle keyboard navigation for image lightbox
    const handleImageKeyDown = useCallback(
        (e) => {
            if ((e.key === "Enter" || e.key === " ") && onImageClick && media?.type === "image") {
                e.preventDefault();
                onImageClick();
            }
        },
        [onImageClick, media]
    );
    // Handle play/pause toggle
    const handlePlayToggle = useCallback(
        (e) => {
            e.stopPropagation();
            setIsPlaying(!isPlaying);
        },
        [isPlaying, setIsPlaying]
    );
    // Handle video errors
    const handleVideoError = useCallback(
        (e) => {
            setVideoError(true);
            setIsPlaying(false);
        },
        [setIsPlaying]
    );
    // No media or invalid media
    if (!media || !media.src) {
        return (
            <div
                className={styles.mediaError}
                aria-label="Media not available"
            >
                Media not available
            </div>
        );
    }
    // VIDEO CASE 1: Type is video but we have no thumbnail and it's not mounted yet
    if (media.type === "video" && !media.poster && !videoMounted) {
        return (
            <div
                className={styles.blackThumbWrapper}
                aria-label="Video thumbnail"
            >
                <button
                    className={styles.playOverlay}
                    onClick={() => setIsPlaying(true)}
                    onKeyDown={(e) => {
                        if (e.key === "Enter" || e.key === " ") {
                            e.preventDefault();
                            setIsPlaying(true);
                        }
                    }}
                    aria-label="Play video"
                    type="button"
                >
                    <FaPlay className={styles.playIcon} aria-hidden="true"/>
                    <span className="sr-only">Play video</span>
                </button>
            </div>
        );
    }
    // IMAGE CASE: Render image with lightbox capability
    if (media.type === "image") {
        return (
            <div
                className={styles.mediaContainer}
                aria-label={media.alt || `${itemTitle} image`}
            >
                <img
                    className={styles.styledImage}
                    src={media.src}
                    alt={media.alt || `${itemTitle} image`}
                    onClick={handleImageClick}
                    onKeyDown={handleImageKeyDown}
                    tabIndex={onImageClick ? 0 : -1}
                    loading="lazy"
                    onError={(e) => {
                        e.currentTarget.src = placeholderImage || "/images/placeholder.webp";
                    }}
                />
                {onImageClick && (
                    <button
                        className={styles.expandButton}
                        onClick={handleImageClick}
                        aria-label="View fullscreen"
                        type="button"
                    >
                        <FaExpand aria-hidden="true"/>
                    </button>
                )}
            </div>
        );
    }
    // VIDEO CASE 2: Video with thumbnail (not yet playing)
    if (media.type === "video" && !videoMounted) {
        return (
            <div
                className={styles.mediaContainer}
                aria-label={media.alt || `${itemTitle} video`}
            >
                <div className={styles.thumbWrapper}>
                    <img
                        className={styles.videoThumb}
                        src={media.poster || placeholderImage}
                        alt=""
                        aria-hidden="true"
                        loading="lazy"
                        onError={(e) => {
                            e.currentTarget.src = placeholderImage || "/images/video-icon.webp";
                        }}
                    />
                    <button
                        className={styles.playOverlay}
                        onClick={() => setIsPlaying(true)}
                        onKeyDown={(e) => {
                            if (e.key === "Enter" || e.key === " ") {
                                e.preventDefault();
                                setIsPlaying(true);
                            }
                        }}
                        aria-label="Play video"
                        type="button"
                    >
                        <FaPlay className={styles.playIcon} aria-hidden="true"/>
                        <span className="sr-only">Play video</span>
                    </button>
                </div>
            </div>
        );
    }
    // VIDEO CASE 3: Active video player
    if (media.type === "video" && videoMounted) {
        return (
            <div
                className={styles.mediaContainer}
                aria-label={media.alt || `${itemTitle} video player`}
            >
                <div className={styles.videoWrapper}>
                    {videoError ? (
                        <div className={styles.videoError}>
                            <p>Video playback error</p>
                            <button
                                onClick={() => {
                                    setVideoError(false);
                                    setIsPlaying(true);
                                }}
                                className={styles.retryButton}
                                aria-label="Retry playing video"
                                type="button"
                            >
                                Retry
                            </button>
                        </div>
                    ) : (
                        <>
                            <video
                                className={styles.videoElement}
                                ref={videoRef}
                                src={media.src}
                                poster={media.poster}
                                preload="metadata"
                                muted
                                playsInline
                                onError={handleVideoError}
                                aria-label={media.alt || `${itemTitle} video`}
                            />
                            <button
                                className={isPlaying ? styles.pauseOverlay : styles.playOverlay}
                                onClick={handlePlayToggle}
                                onKeyDown={(e) => {
                                    if (e.key === "Enter" || e.key === " ") {
                                        e.preventDefault();
                                        handlePlayToggle(e);
                                    }
                                }}
                                aria-label={isPlaying ? "Pause video" : "Play video"}
                                type="button"
                            >
                                {isPlaying ? (
                                    <FaPause className={styles.pauseIcon} aria-hidden="true"/>
                                ) : (
                                    <FaPlay className={styles.playIcon} aria-hidden="true"/>
                                )}
                                <span className="sr-only">
                  {isPlaying ? "Pause video" : "Play video"}
                </span>
                            </button>
                        </>
                    )}
                </div>
            </div>
        );
    }
    // Fallback for unknown media type
    return (
        <div className={styles.mediaError} aria-label="Unsupported media">
            Unsupported media type
        </div>
    );
}
MediaDisplay.propTypes = {
    media: PropTypes.shape({
        type: PropTypes.oneOf(["image", "video"]),
        src: PropTypes.string.isRequired,
        alt: PropTypes.string,
        poster: PropTypes.string,
    }),
    onImageClick: PropTypes.func,
    placeholderImage: PropTypes.string,
    isPlaying: PropTypes.bool,
    setIsPlaying: PropTypes.func,
    itemTitle: PropTypes.string,
};
MediaDisplay.defaultProps = {
    media: null,
    onImageClick: null,
    placeholderImage: "/images/placeholder.webp",
    isPlaying: false,
    setIsPlaying: () => {
    },
    itemTitle: "Item",
};
export default memo(MediaDisplay);