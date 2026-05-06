// MediaDisplayEditable.jsx
import { FaPause, FaPlay } from '../../utils/iconImports';
import React, {useState, useRef, useEffect, memo} from 'react';
import styles from './MediaDisplayEditable.module.css';
function MediaDisplayEditable({
                          media,
                          onImageClick,
                          placeholderImage,
                          isPlaying,
                          setIsPlaying,
                      }) {
    const videoRef = useRef(null);
    const [videoMounted, setVideoMounted] = useState(false);
    useEffect(() => {
        // Only mount <video> once user decides to play
        if (isPlaying && !videoMounted) {
            setVideoMounted(true);
        }
    }, [isPlaying, videoMounted]);
    useEffect(() => {
        // Control playback
        if (!videoRef.current) return;
        if (isPlaying) {
            videoRef.current.play().catch((err) => {
            });
        } else {
            videoRef.current.pause();
        }
    }, [isPlaying]);
    // === 1) Fallback: If 'video' + no images, display black container with icon ===
    if (media.type === 'video' && media.noImages === true && !videoMounted) {
        return (
            <div className={styles.mediaContainer} aria-label="No Images - Video Fallback">
                <div className={styles.blackThumbWrapper}>
                    <div 
                        className={styles.playOverlay}
                        onClick={(e) => {
                            e.stopPropagation();
                            setIsPlaying(true);
                        }}
                        tabIndex={0}
                        aria-label="Play Video"
                    >
                        <FaPlay className={styles.playIcon} size={50}/>
                    </div>
                </div>
            </div>
        );
    }
    // === 2) For normal images ===
    if (media.type === 'image') {
        return (
            <div className={styles.mediaContainer} aria-label="Image Display">
                <img
                    className={styles.styledImage}
                    src={media.src}
                    alt={media.alt || ''}
                    onClick={onImageClick}
                    tabIndex={0}
                    onKeyUp={(e) => {
                        if (e.key === 'Enter' || e.key === ' ') onImageClick();
                    }}
                    onError={(e) => {
                        e.currentTarget.src =
                            placeholderImage || '/images/video-icon.webp';
                    }}
                />
            </div>
        );
    }
    // === 3) For videos with normal fallback or mounted video ===
    if (media.type === 'video') {
        return (
            <div className={styles.mediaContainer} aria-label="Video Player">
                {!videoMounted ? (
                    <div className={styles.thumbWrapper}>
                        <img
                            className={styles.videoThumb}
                            src={media.poster || placeholderImage}
                            alt={media.alt || 'Video Thumbnail'}
                        />
                        {!isPlaying && (
                            <div
                                className={styles.playOverlay}
                                onClick={(e) => {
                                    e.stopPropagation();
                                    setIsPlaying(true);
                                }}
                                tabIndex={0}
                                aria-label="Play Video"
                            >
                                <FaPlay className={styles.playIcon}/>
                            </div>
                        )}
                    </div>
                ) : (
                    <div className={styles.videoWrapper}>
                        <video
                            className={styles.videoElement}
                            ref={videoRef}
                            src={media.src}
                            preload="metadata"
                            muted
                            loop={false}
                            controls={false}
                            playsInline
                            onError={(e) => }
                        />
                        {isPlaying ? (
                            <div
                                className={styles.pauseOverlay}
                                onClick={(e) => {
                                    e.stopPropagation();
                                    setIsPlaying(false);
                                }}
                                tabIndex={0}
                                aria-label="Pause Video"
                            >
                                <FaPause className={styles.pauseIcon}/>
                            </div>
                        ) : (
                            <div
                                className={styles.playOverlay}
                                onClick={(e) => {
                                    e.stopPropagation();
                                    setIsPlaying(true);
                                }}
                                tabIndex={0}
                                aria-label="Play Video"
                            >
                                <FaPlay className={styles.playIcon}/>
                            </div>
                        )}
                    </div>
                )}
            </div>
        );
    }
    // If media type unrecognized
    return null;
}
export default memo(MediaDisplayEditable);
