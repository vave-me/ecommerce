"use client";
import React, { useRef, useState, useEffect } from 'react';
import { Play } from '@/icons';
import styles from './VideoPlayer.module.css';

const VideoPlayer = ({ 
    media, 
    index, 
    isPlaying, 
    onToggle, 
    resolvedMediaLength, 
    productName,
    className = ''
}) => {
    const videoRef = useRef(null);
    const containerRef = useRef(null);
    const [showControls, setShowControls] = useState(true);
    const [hasLoaded, setHasLoaded] = useState(false);
    const [thumbnailUrl, setThumbnailUrl] = useState(null);
    
    // Handle play/pause from external state
    useEffect(() => {
        if (!videoRef.current) return;
        
        if (isPlaying) {
            videoRef.current.play().catch(err => {
                // Error: 'Video play error:', err...
            });
        } else {
            videoRef.current.pause();
        }
    }, [isPlaying]);
    
    // Show/hide controls based on play state
    useEffect(() => {
        setShowControls(!isPlaying);
    }, [isPlaying]);
    
    const handleVideoClick = (e) => {
        e.preventDefault();
        e.stopPropagation();
        onToggle(index, videoRef);
    };
    
    const handleVideoEnded = () => {
        // Show controls when video ends
        setShowControls(true);
        // Update playing state
        onToggle(index, { current: null }); // Pass null ref to just update state
    };
    
    const handleLoadedData = () => {
        setHasLoaded(true);
        
        // Generate thumbnail from first frame if not provided
        if (!media.thumbnail && videoRef.current) {
            try {
                const canvas = document.createElement('canvas');
                canvas.width = videoRef.current.videoWidth;
                canvas.height = videoRef.current.videoHeight;
                const ctx = canvas.getContext('2d');
                ctx.drawImage(videoRef.current, 0, 0, canvas.width, canvas.height);
                setThumbnailUrl(canvas.toDataURL());
            } catch (err) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', err);
        }
    }
        }
    };
    
    // Intersection Observer to pause video when out of view
    useEffect(() => {
        if (!containerRef.current || !videoRef.current) return;
        
        const observer = new IntersectionObserver(
            (entries) => {
                entries.forEach(entry => {
                    if (!entry.isIntersecting && isPlaying) {
                        // Pause video when it goes out of view
                        videoRef.current?.pause();
                        onToggle(index, { current: null });
                    }
                });
            },
            { threshold: 0.1 }
        );
        
        observer.observe(containerRef.current);
        
        return () => {
            observer.disconnect();
        };
    }, [isPlaying, index, onToggle]);
    
    return (
        <div ref={containerRef} className={`${styles.videoContainer} ${className}`}>
            <video 
                ref={videoRef}
                src={media.url}
                muted
                playsInline
                onEnded={handleVideoEnded}
                onLoadedData={handleLoadedData}
                onClick={handleVideoClick}
                className={styles.video}
                style={{ width: '100%', height: '100%', objectFit: 'cover' }}
                onError={(e) => {}}
                aria-label={`${productName} video ${index + 1} of ${resolvedMediaLength}`}
                data-video-index={index}
                poster={media.thumbnail || thumbnailUrl || ''}
            />
            
            {/* Custom Play Button Overlay */}
            {showControls && hasLoaded && (
                <button 
                    className={styles.playButton}
                    onClick={handleVideoClick}
                    aria-label={`Play ${productName} video`}
                    type="button"
                >
                    <div className={styles.playIconWrapper}>
                        <Play size={24} />
                    </div>
                </button>
            )}
            
            {/* Loading state */}
            {!hasLoaded && (
                <div className={styles.loadingOverlay}>
                    <div className={styles.spinner} />
                </div>
            )}
        </div>
    );
};

export default VideoPlayer;