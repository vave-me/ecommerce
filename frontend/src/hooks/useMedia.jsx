// hooks/usePostMedia.js - fixed version
import {useEffect, useState} from "react";
import {getMediaByItem, isMediaResponseSuccess, getMediaErrorMessage} from "../api/mediaApi";
/**
 * Process raw media API response into a standardized format
 * @param {Object} response - Raw API response from getMediaByItem
 * @returns {Array} Processed array of media objects
 */
export function processMediaResponse(response) {
    // Check if the response has the expected structure
    if (!response?.media?.mediaOrder || !Array.isArray(response.media.mediaOrder)) {
        return [];
    }
    return response.media.mediaOrder.map((mItem, index) => {
        // Safety check for valid URL
        const url = mItem.url || '';
        // Determine media type based on URL extension or path pattern
        let type = 'image';
        const lowerUrl = url.toLowerCase();
        if (
            lowerUrl.endsWith('.mp4') ||
            lowerUrl.endsWith('.mov') ||
            lowerUrl.endsWith('.webm') ||
            lowerUrl.includes('/video/') ||
            (mItem.type && mItem.type.toLowerCase().includes('video'))
        ) {
            type = 'video';
        }
        return {
            id: mItem.mediaItemId || `media-${index}`,
            type,
            src: url,
            alt: mItem.altText || `Media ${index + 1}`,
            displayOrder: mItem.displayOrder || index,
            poster: type === 'video' ? mItem.posterUrl : undefined
        };
    });
}
/**
 * Custom hook to fetch and process media for a post
 * @param {string} itemId - ID of the post to fetch media for
 * @returns {Object} Media data with loading and error states
 */
export function useMedia(itemId) {
    const [media, setMedia] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [thumbnail, setThumbnail] = useState(null);
    const [mediaId, setMediaId] = useState(null);
    useEffect(() => {
        // Reset states when itemId changes
        if (!itemId) {
            setMedia([]);
            setThumbnail(null);
            setMediaId(null);
            setError(null);
            return;
        }
        let isMounted = true;
        setLoading(true);
        getMediaByItem(itemId)
            .then(response => {
                if (!isMounted) return;
                // Check if the API response was successful
                if (!isMediaResponseSuccess(response)) {
                    // Handle API error response
                    const errorMessage = getMediaErrorMessage(response);
                    if (response?.severity !== 'warning') {
                        // Only set error for non-warning issues (404s are warnings)
                        setError(new Error(errorMessage));
                    }
                    // For 404s and warnings, just continue with empty media
                    setMedia([]);
                    setThumbnail(null);
                    setMediaId(null);
                    return;
                }
                // Store the media container ID for potential future use
                if (response?.media?.id) {
                    setMediaId(response.media.id);
                }
                const processedMedia = processMediaResponse(response);
                setMedia(processedMedia);
                // Set the first image as thumbnail
                const firstImage = processedMedia.find(item => item.type === 'image');
                if (firstImage) {
                    setThumbnail(firstImage.src);
                } else if (processedMedia.length > 0) {
                    // If no images, use the first media item as thumbnail (could be a video)
                    setThumbnail(processedMedia[0].src);
                } else {
                    setThumbnail(null);
                }
            })
            .catch(err => {
                // This should rarely happen with the new API design,
                // but keep for backward compatibility
                if (isMounted) {
                    setError(err);
                }
            })
            .finally(() => {
                if (isMounted) {
                    setLoading(false);
                }
            });
        return () => {
            isMounted = false;
        };
    }, [itemId]);
    return {
        media,
        loading,
        error,
        thumbnail,
        mediaId, // Added mediaId to the returned values
        isEmpty: media.length === 0
    };
}