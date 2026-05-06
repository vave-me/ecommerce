// hooks/usePostInteractions.js
import {useState, useCallback} from 'react';
/**
 * Hook to manage post interactions state
 * @returns {Object} State and handlers for post interactions
 */
export function useInteractions() {
    const [currentMediaIndex, setCurrentMediaIndex] = useState(0);
    const [activeVideoIndex, setActiveVideoIndex] = useState(null);
    const [isProcessing, setIsProcessing] = useState(false);
    /**
     * Handle media navigation (next/prev)
     * @param {string} direction - Direction to navigate ('next' or 'prev')
     * @param {number} galleryLength - Length of the gallery array
     */
    const handleMediaNavigation = useCallback((direction, galleryLength) => {
        if (!galleryLength) return;
        setCurrentMediaIndex(prev => {
            if (direction === 'next') {
                return (prev + 1) % galleryLength;
            } else {
                return (prev - 1 + galleryLength) % galleryLength;
            }
        });
        // Reset active video when changing media
        setActiveVideoIndex(null);
    }, []);
    return {
        currentMediaIndex,
        setCurrentMediaIndex,
        activeVideoIndex,
        setActiveVideoIndex,
        isProcessing,
        setIsProcessing,
        handleMediaNavigation,
    };
}