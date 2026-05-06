// File: app/videos/[itemId]/ItemVideosPageClient.jsx
"use client";
import React, { useEffect, useRef, useState, memo, useCallback } from 'react';
import {useIntersectionAutoPlayPause} from "./useIntersectionAutoPlayPause";

// Reuse your existing Intersection Observer hook, or define it here

// (If you have that in a separate file, or keep it inline—either is fine.)

const styles = {
    errorContainer: { /* same style objects as before */ },
    errorText: { /* ... */ },
    container: { /* ... */ },
    heading: { /* ... */ },
    noVideos: { /* ... */ },
    snapContainer: { /* ... */ },
    pageTitle: { /* ... */ },
    videoSnapItem: { /* ... */ },
    cardHeader: { /* ... */ },
    videoTitle: { /* ... */ },
    muteButton: { /* ... */ },
    videoEl: { /* ... */ },
};

/* ------------------------------------------
   CLIENT COMPONENT - OPTIMIZED with React.memo
   ------------------------------------------ */
const ItemVideosPageClient = memo(function ItemVideosPageClient({ itemId, videos, errorMessage }) {
    // If there's an error, display it
    if (errorMessage) {
        return (
            <div style={styles.errorContainer}>
                <h2 style={styles.errorText}>Error</h2>
                <p>{errorMessage}</p>
            </div>
        );
    }

    // If no videos found, show message
    if (!videos || videos.length === 0) {
        return (
            <div style={styles.container}>
                <h1 style={styles.heading}>Short Video Feed (itemId: {itemId})</h1>
                <p style={styles.noVideos}>No videos found for item {itemId}.</p>
            </div>
        );
    }

    // Render the videos in a scroll-snapping container
    return (
        <div style={styles.snapContainer}>
            <h1 style={styles.pageTitle}>Short Video Feed (itemId: {itemId})</h1>
            {videos.map((vid) => (
                <VideoCard key={vid.id} video={vid} />
            ))}
        </div>
    );
}, (prevProps, nextProps) => {
    // Custom comparison for optimal performance
    return (
        prevProps.itemId === nextProps.itemId &&
        prevProps.errorMessage === nextProps.errorMessage &&
        prevProps.videos?.length === nextProps.videos?.length &&
        prevProps.videos?.every((video, index) => 
            video.id === nextProps.videos?.[index]?.id
        )
    );
});

/* ------------------------------------------
   VIDEO CARD (client subcomponent) - OPTIMIZED with React.memo
   ------------------------------------------ */
const VideoCard = memo(function VideoCard({ video }) {
    const videoRef = useRef(null);
    const [isMuted, setIsMuted] = useState(true);

    // Handle auto-play/pause
    useIntersectionAutoPlayPause(videoRef);

    const toggleMute = useCallback(() => setIsMuted((prev) => !prev), []);

    return (
        <div style={styles.videoSnapItem}>
            <div style={styles.cardHeader}>
                <h3 style={styles.videoTitle}>
                    {video.metadata || `Untitled Video (ID: ${video.id})`}
                </h3>
                <button onClick={toggleMute} style={styles.muteButton}>
                    {isMuted ? 'Unmute' : 'Mute'}
                </button>
            </div>

            <video
                ref={videoRef}
                src={video.url}
                poster={video.thumbnail || ''}
                loop
                muted={isMuted}
                controls
                style={styles.videoEl}
            />
        </div>
    );
}, (prevProps, nextProps) => {
    // Custom comparison for video props
    return (
        prevProps.video.id === nextProps.video.id &&
        prevProps.video.url === nextProps.video.url &&
        prevProps.video.metadata === nextProps.video.metadata &&
        prevProps.video.thumbnail === nextProps.video.thumbnail
    );
});

export default ItemVideosPageClient;
