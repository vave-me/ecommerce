"use client";
import { useEffect } from 'react';

/**
 * useIntersectionAutoPlayPause
 *
 * A custom hook that automatically plays a video
 * when it enters the viewport and pauses it when
 * it leaves the viewport, using IntersectionObserver.
 *
 * @param {React.RefObject<HTMLVideoElement>} videoRef - A ref to the <video> element
 * @param {number} threshold - The fraction of the video that needs to be in view to trigger playing
 */
export function useIntersectionAutoPlayPause(videoRef, threshold = 0.7) {
    useEffect(() => {
        if (!videoRef?.current) return;

        const handlePlayPause = (entries) => {
            entries.forEach((entry) => {
                if (entry.isIntersecting) {
                    videoRef.current.play().catch((err) => {
                        // In some browsers, user interaction might be needed to play
                        // Error: 'Auto-play error:', err...
                    });
                } else {
                    videoRef.current.pause();
                }
            });
        };

        const observer = new IntersectionObserver(handlePlayPause, { threshold });
        observer.observe(videoRef.current);

        // Cleanup
        return () => {
            if (videoRef.current) {
                observer.unobserve(videoRef.current);
            }
        };
    }, [videoRef, threshold]);
}
