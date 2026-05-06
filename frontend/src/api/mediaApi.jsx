// File: src/api/mediaApi.js (Recommended new name)
import axiosInstance from './axiosInstance';
const MEDIA_ENDPOINT = '/media'; // Base path relative to axiosInstance.baseURL
/**
 * Helper to safely encode URI components.
 * @param {string | number} component - The path component to encode.
 * @returns {string} The encoded component.
 */
const safeEncode = (component) => {
    if (component === null || typeof component === 'undefined') {
        return ''; // Or throw an error if this should be disallowed
    }
    return encodeURIComponent(String(component));
};
/**
 * Enhanced error handler for media API calls
 * @param {Error} error - The error object
 * @param {string} endpoint - The endpoint that failed
 * @param {string} operation - Description of the operation
 * @returns {Object} Standardized error response
 */
const handleMediaError = (error, endpoint, operation) => {
    const errorDetails = {
        success: false,
        operation,
        endpoint,
        timestamp: new Date().toISOString()
    };
    if (error.response) {
        // Server responded with error status
        const { status, data } = error.response;
        errorDetails.status = status;
        errorDetails.message = data?.message || `Server error (${status})`;
        errorDetails.data = data;
        // Handle specific HTTP errors
        switch (status) {
            case 404:
                errorDetails.userMessage = 'Media not found';
                errorDetails.severity = 'warning';
                break;
            case 500:
                errorDetails.userMessage = 'Server error - media temporarily unavailable';
                errorDetails.severity = 'error';
                break;
            case 403:
                errorDetails.userMessage = 'Access denied to media';
                errorDetails.severity = 'error';
                break;
            case 401:
                errorDetails.userMessage = 'Authentication required';
                errorDetails.severity = 'error';
                break;
            default:
                errorDetails.userMessage = 'Unable to load media';
                errorDetails.severity = 'error';
        }
        // Log only in development with safe data serialization
        if (process.env.NODE_ENV === 'development') {
            // Safely serialize response data to avoid empty {} logs
            let safeData = 'No response data';
            if (data) {
                if (typeof data === 'object') {
                    try {
                        const jsonString = JSON.stringify(data);
                        safeData = jsonString !== '{}' ? jsonString : 'Empty response object';
                    } catch (e) {
                        safeData = '[Non-serializable object]';
                    }
                } else {
                    safeData = String(data);
                }
            }
        }
    } else if (error.request) {
        // Network error
        errorDetails.userMessage = 'Network error - please check your connection';
        errorDetails.severity = 'error';
        errorDetails.network = true;
        if (process.env.NODE_ENV === 'development') {
        }
    } else {
        // Other error
        errorDetails.userMessage = 'Unexpected error occurred';
        errorDetails.severity = 'error';
        errorDetails.message = error.message;
        if (process.env.NODE_ENV === 'development') {
        }
    }
    return errorDetails;
};
/**
 * GET MEDIA INFO BY ITEM ID
 * GET /media/item/{itemId}
 * Returns { media: { id, itemId, itemType, userId, ... } } or null if not found
 */
export const getMediaByItem = async (itemId) => {
    if (!itemId) {
        throw new Error('itemId is required for getMediaByItem.');
    }
    const endpoint = `${MEDIA_ENDPOINT}/item/${safeEncode(itemId)}`;
    try {
        const response = await axiosInstance.get(endpoint);
        return { success: true, ...response.data };
    } catch (error) {
        // For 404 and 500, return null instead of throwing (media might not exist yet)
        if (error.response?.status === 404 || error.response?.status === 500) {
            // Log informational message only in development for 500 errors
            if (error.response?.status === 500 && process.env.NODE_ENV === 'development') {
            }
            return { success: true, media: null };
        }
        // For other errors, handle normally
        const errorResponse = handleMediaError(error, endpoint, 'getMediaByItem');
        return errorResponse;
    }
};
/**
 * GET ALL IMAGES BELONGING TO AN ITEM
 * GET /media/item/{itemId}/image
 * Response -> { images: [ { ... }, ... ] }
 */
export const getAllItemImages = async (itemId) => {
    if (!itemId) {
        throw new Error('itemId is required for getAllItemImages.');
    }
    const endpoint = `${MEDIA_ENDPOINT}/item/${safeEncode(itemId)}/image`;
    try {
        const response = await axiosInstance.get(endpoint);
        return { success: true, ...response.data };
    } catch (error) {
        // For 404 and 500, return empty array instead of throwing
        if (error.response?.status === 404 || error.response?.status === 500) {
            // Log informational message only in development for 500 errors
            if (error.response?.status === 500 && process.env.NODE_ENV === 'development') {
            }
            return { success: true, images: [] };
        }
        // For other errors, handle normally
        const errorResponse = handleMediaError(error, endpoint, 'getAllItemImages');
        return errorResponse;
    }
};
/**
 * GET ALL VIDEOS BELONGING TO AN ITEM
 * GET /media/item/{itemId}/video
 * Response -> { videos: [ { ... }, ... ] }
 */
export const getAllItemVideos = async (itemId) => {
    if (!itemId) {
        throw new Error('itemId is required for getAllItemVideos.');
    }
    const endpoint = `${MEDIA_ENDPOINT}/item/${safeEncode(itemId)}/video`;
    try {
        const response = await axiosInstance.get(endpoint);
        return { success: true, ...response.data };
    } catch (error) {
        // For 404 and 500, return empty array instead of throwing
        if (error.response?.status === 404 || error.response?.status === 500) {
            // Log informational message only in development for 500 errors
            if (error.response?.status === 500 && process.env.NODE_ENV === 'development') {
            }
            return { success: true, videos: [] };
        }
        // For other errors, handle normally
        const errorResponse = handleMediaError(error, endpoint, 'getAllItemVideos');
        return errorResponse;
    }
};
/**
 * GET MEDIA INFO BY MEDIA ID
 * GET /media/{mediaId}
 * Response -> { media: {...} }
 */
export const getMedia = async (mediaId) => {
    if (!mediaId) {
        throw new Error('mediaId is required for getMedia.');
    }
    const endpoint = `${MEDIA_ENDPOINT}/${safeEncode(mediaId)}`;
    try {
        const response = await axiosInstance.get(endpoint);
        return { success: true, ...response.data };
    } catch (error) {
        // For 404 and 500, return null instead of throwing (media might not exist yet)
        if (error.response?.status === 404 || error.response?.status === 500) {
            return { success: true, media: null };
        }
        // For other errors, handle normally
        const errorResponse = handleMediaError(error, endpoint, 'getMedia');
        return errorResponse;
    }
};
/**
 * GET ALL IMAGES BELONGING TO A MEDIA RECORD
 * GET /media/{mediaId}/image
 * Response -> { images: [ ... ] }
 */
export const getAllMediaImages = async (mediaId) => {
    if (!mediaId) {
        throw new Error('mediaId is required for getAllMediaImages.');
    }
    const endpoint = `${MEDIA_ENDPOINT}/${safeEncode(mediaId)}/image`;
    try {
        const response = await axiosInstance.get(endpoint);
        return { success: true, ...response.data };
    } catch (error) {
        // For 404 and 500, return empty array instead of throwing
        if (error.response?.status === 404 || error.response?.status === 500) {
            return { success: true, images: [] };
        }
        // For other errors, handle normally
        const errorResponse = handleMediaError(error, endpoint, 'getAllMediaImages');
        return errorResponse;
    }
};
/**
 * GET ALL VIDEOS BELONGING TO A MEDIA RECORD
 * GET /media/{mediaId}/video
 * Response -> { videos: [ ... ] }
 */
export const getAllMediaVideos = async (mediaId) => {
    if (!mediaId) {
        throw new Error('mediaId is required for getAllMediaVideos.');
    }
    const endpoint = `${MEDIA_ENDPOINT}/${safeEncode(mediaId)}/video`;
    try {
        const response = await axiosInstance.get(endpoint);
        return { success: true, ...response.data };
    } catch (error) {
        // For 404 and 500, return empty array instead of throwing
        if (error.response?.status === 404 || error.response?.status === 500) {
            return { success: true, videos: [] };
        }
        // For other errors, handle normally
        const errorResponse = handleMediaError(error, endpoint, 'getAllMediaVideos');
        return errorResponse;
    }
};
/**
 * GET ALL VIDEOS (PAGINATED)
 * GET /media/videos
 * Query params: page, pageSize, sortBy, sortOrder
 * Response -> { videos: [...], totalCount, currentPage, totalPages }
 */
export const getAllVideos = async (filters = {}) => {
    const endpoint = `${MEDIA_ENDPOINT}/videos`;
    // Only pass defined filters to avoid sending `undefined` as query params
    const validFilters = Object.entries(filters).reduce((acc, [key, value]) => {
        if (value !== undefined && value !== null) {
            acc[key] = value;
        }
        return acc;
    }, {});
    try {
        const response = await axiosInstance.get(endpoint, { params: validFilters });
        return { success: true, ...response.data };
    } catch (error) {
        // For 404 and 500, return empty results instead of throwing
        if (error.response?.status === 404 || error.response?.status === 500) {
            return { 
                success: true, 
                videos: [], 
                totalCount: 0, 
                currentPage: 1, 
                totalPages: 0 
            };
        }
        // For other errors, handle normally
        const errorResponse = handleMediaError(error, endpoint, 'getAllVideos');
        return errorResponse;
    }
};
/**
 * Utility function to check if a media API response was successful
 * @param {Object} response - Response from any media API function
 * @returns {boolean} True if successful, false if error
 */
export const isMediaResponseSuccess = (response) => {
    return response && response.success === true;
};
/**
 * Utility function to get user-friendly error message from media API response
 * @param {Object} response - Response from any media API function
 * @returns {string} User-friendly error message
 */
export const getMediaErrorMessage = (response) => {
    if (isMediaResponseSuccess(response)) {
        return null;
    }
    return response?.userMessage || 'Unable to load media';
};