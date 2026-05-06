// src/api/client/mediaApi.jsx
import axiosInstance from '../axiosInstance';

/**
 * CREATE A NEW MEDIA
 * POST /api/media
 * Body -> mediapbCreateMediaRequest
 */
export const createMedia = async (mediaData) => {
    const js = JSON.stringify(mediaData);
    try {
        const response = await axiosInstance.post('/media', js);
        // Response -> mediapbCreateMediaResponse { id: string }
        return response.data;
    } catch (error) {
        if (process.env.NODE_ENV === 'development') {
            // Error creating media logged for debugging
        }
        throw error;
    }
};

export const updateMedia = async (mediaData) => {
    const js = JSON.stringify(mediaData);
    try {
        const response = await axiosInstance.post('/media/update', js);
        // Response -> mediapbCreateMediaResponse { id: string }
        return response.data;
    } catch (error) {
        if (process.env.NODE_ENV === 'development') {
            // Error updating media logged for debugging
        }
        throw error;
    }
};

/**
 * ADD IMAGE
 * POST /api/media/image
 * Body -> mediapbAddImageRequest
 */
export const addImage = async (imageData) => {
    try {
        const response = await axiosInstance.post('/media/image', imageData);
        // Response -> mediapbAddImageResponse { url: string, viewUrl: string }
        let data = response.data;
        
        // Development-only URL transformation for MinIO
        if (process.env.NODE_ENV === 'development' && data.viewUrl) {
            data.viewUrl = data.viewUrl.replace(
                'https://minio-api.sfx-markt.de',
                'http://localhost:9096'
            );
        }
        
        return data;
    } catch (error) {
        if (process.env.NODE_ENV === 'development') {
            // Error logged for debugging
        }
        throw error;
    }
};

/**
 * ADD VIDEO
 * POST /api/media/video
 * Body -> mediapbAddVideoRequest
 */
export const addVideo = async (videoData) => {
    try {
        const response = await axiosInstance.post('/media/video', videoData);
        // Response -> mediapbAddVideoResponse { url: string, viewUrl: string }
        let data = response.data;
        
        // Development-only URL transformation for MinIO
        if (process.env.NODE_ENV === 'development' && data.viewUrl) {
            data.viewUrl = data.viewUrl.replace(
                'https://minio-api.sfx-markt.de',
                'http://localhost:9096'
            );
        }
        
        return data;
    } catch (error) {
        if (process.env.NODE_ENV === 'development') {
            // Error logged for debugging
        }
        throw error;
    }
};

/**
 * GET MEDIA BY ITEM (CUSTOM ENDPOINT)
 * GET /api/media/item/{itemId}
 * Response -> mediapbGetMediaByItemResponse { media: {...} }
 */
export const getMediaByItem = async (itemId) => {
    try {
        const response = await axiosInstance.get(`/media/item/${itemId}`);
        // Returns { media: { id, itemId, itemType, userId, ... } }
        return response.data;
    } catch (error) {
        // Handle 404 and 500 gracefully - item might not have media yet
        if (error.response?.status === 404 || error.response?.status === 500) {
            // Log informational message only in development for 500 errors
            if (error.response?.status === 500 && process.env.NODE_ENV === 'development') {
                // Server error logged for debugging
            }
            return { media: null };
        }
        
        // Log other errors for debugging
        if (process.env.NODE_ENV === 'development') {
            // Error details logged for debugging
        }
        throw error;
    }
};

/**
 * GET ALL ITEM IMAGES
 * GET /api/media/item/{itemId}/image
 * Response -> mediapbGetAllItemImagesResponse { images: [...] }
 */
export const getAllItemImages = async (itemId) => {
    try {
        const response = await axiosInstance.get(`/media/item/${itemId}/image`);
        // Returns { images: [...] }
        return response.data;
    } catch (error) {
        // Handle 404 and 500 gracefully - item might not have images yet
        if (error.response?.status === 404 || error.response?.status === 500) {
            // Log informational message only in development for 500 errors
            if (error.response?.status === 500 && process.env.NODE_ENV === 'development') {
                // Server error logged for debugging
            }
            return { images: [] };
        }
        
        // Log other errors for debugging
        if (process.env.NODE_ENV === 'development') {
            // Error details logged for debugging
        }
        throw error;
    }
};

/**
 * GET ALL ITEM VIDEOS
 * GET /api/media/item/{itemId}/video
 * Response -> mediapbGetAllItemVideosResponse { videos: [...] }
 */
export const getAllItemVideos = async (itemId) => {
    try {
        const response = await axiosInstance.get(`/media/item/${itemId}/video`);
        // Returns { videos: [...] }
        return response.data;
    } catch (error) {
        // Handle 404 and 500 gracefully - item might not have videos yet
        if (error.response?.status === 404 || error.response?.status === 500) {
            // Log informational message only in development for 500 errors
            if (error.response?.status === 500 && process.env.NODE_ENV === 'development') {
                // Server error logged for debugging
            }
            return { videos: [] };
        }
        
        // Log other errors for debugging
        if (process.env.NODE_ENV === 'development') {
            // Error details logged for debugging
        }
        throw error;
    }
};

/**
 * GET MEDIA BY ID
 * GET /api/media/{mediaId}
 * Response -> mediapbGetMediaResponse { media: {...} }
 */
export const getMedia = async (mediaId) => {
    try {
        const response = await axiosInstance.get(`/media/${mediaId}`);
        // Returns { media: { id, itemId, itemType, userId, ... } }
        return response.data;
    } catch (error) {
        // Handle 404 and 500 gracefully - media might not exist
        if (error.response?.status === 404 || error.response?.status === 500) {

            return { media: null };
        }
        
        // Log other errors for debugging
        if (process.env.NODE_ENV === 'development') {
            // Error logged for debugging
        }
        throw error;
    }
};

/**
 * REMOVE MEDIA
 * DELETE /api/media/{mediaId}?itemId=XYZ (optional query param)
 * Response -> mediapbRemoveMediaResponse { id: string }
 */
export const removeMedia = async (mediaId, itemId) => {
    try {
        const response = await axiosInstance.delete(`/media/${mediaId}`, {
            params: {itemId},
        });
        return response.data;
    } catch (error) {
        if (process.env.NODE_ENV === 'development') {
            // Error logged for debugging
        }
        throw error;
    }
};

/**
 * GET ALL IMAGES BELONGING TO THE MEDIA
 * GET /api/media/{mediaId}/image
 * No query params in the Swagger
 * Response -> mediapbGetAllMediaImagesResponse { images: [ ... ] }
 */
export const getAllMediaImages = async (mediaId) => {
    try {
        const response = await axiosInstance.get(`/media/${mediaId}/image`);
        return response.data; // { images: [ { ... }, ... ] }
    } catch (error) {
        // Handle 404 and 500 gracefully - media might not have images yet
        if (error.response?.status === 404 || error.response?.status === 500) {
            if (process.env.NODE_ENV === 'development') {
                // Server error - returning empty array
            }
            return { images: [] };
        }
        
        // Log other errors for debugging
        if (process.env.NODE_ENV === 'development') {
            // Error: `[MediaAPI] Error fetching media images for ${medi
        }
        throw error;
    }
};

/**
 * REMOVE IMAGE
 * DELETE /api/media/{mediaId}/image/{imageId}
 * Response -> mediapbRemoveImageResponse { id: string }
 */
export const removeImage = async (mediaId, imageId) => {
    try {
        const response = await axiosInstance.delete(`/media/${mediaId}/image/${imageId}`);
        return response.data;
    } catch (error) {
        if (process.env.NODE_ENV === 'development') {
            // Error: `[MediaAPI] Error removing image ${imageId} from m
        }
        throw error;
    }
};

/**
 * GET ALL VIDEOS BELONGING TO THE MEDIA
 * GET /api/media/{mediaId}/video
 * No query params in the Swagger
 * Response -> mediapbGetAllMediaVideosResponse { videos: [ ... ] }
 */
export const getAllMediaVideos = async (mediaId) => {
    try {
        const response = await axiosInstance.get(`/media/${mediaId}/video`);
        return response.data; // { videos: [ ... ] }
    } catch (error) {
        // Handle 404 and 500 gracefully - media might not have videos yet
        if (error.response?.status === 404 || error.response?.status === 500) {
            if (process.env.NODE_ENV === 'development') {
                // Server error - returning empty array
            }
            return { videos: [] };
        }
        
        // Log other errors for debugging
        if (process.env.NODE_ENV === 'development') {
            // Error: `[MediaAPI] Error fetching media videos for ${medi
        }
        throw error;
    }
};

/**
 * GET ALL VIDEOS
 * GET /api/media/videos
 * Query params: page, pageSize, sortBy, sortOrder
 * Response -> mediapbGetAllVideosResponse { videos: [...], totalCount, currentPage, totalPages }
 */
export const getAllVideos = async (filters = {}) => {
    const validFilters = Object.entries(filters).reduce((acc, [key, value]) => {
        if (value !== undefined && value !== null) {
            acc[key] = value;
        }
        return acc;
    }, {});

    try {
        const response = await axiosInstance.get('/media/videos', { params: validFilters });
        return response.data;
    } catch (error) {
        // Handle 404 and 500 gracefully - no videos found
        if (error.response?.status === 404 || error.response?.status === 500) {
            if (process.env.NODE_ENV === 'development') {
                // Server error - returning empty results
            }
            return { 
                videos: [], 
                totalCount: 0, 
                currentPage: 1, 
                totalPages: 0 
            };
        }
        
        // Log other errors for debugging
        if (process.env.NODE_ENV === 'development') {
            // Error logged for debugging
        }
        throw error;
    }
};

/**
 * REMOVE VIDEO
 * DELETE /api/media/{mediaId}/video/{videoId}
 * Response -> mediapbRemoveVideoResponse { id: string }
 */
export const removeVideo = async (mediaId, videoId) => {
    try {
        const response = await axiosInstance.delete(`/media/${mediaId}/video/${videoId}`);
        return response.data;
    } catch (error) {
        if (process.env.NODE_ENV === 'development') {
            // Error: `[MediaAPI] Error removing video ${videoId} from m
        }
        throw error;
    }
};

/**
 * Utility function to check if a media API response was successful and has data
 * @param {Object} response - Response from any media API function
 * @returns {boolean} True if response contains valid data
 */
export const isMediaResponseSuccess = (response) => {
    // Check if response exists and contains data (not just null)
    return response && (
        (response.media !== null && response.media !== undefined) || 
        (response.images && response.images.length > 0) || 
        (response.videos && response.videos.length > 0) ||
        response.id || // for create/update responses
        response.url // for add image/video responses
    );
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

/**
 * Utility function to safely get media from a response
 * @param {Object} response - Response from getMediaByItem or getMedia
 * @returns {Object|null} Media object or null if not found
 */
export const getMediaFromResponse = (response) => {
    return response?.media || null;
};

/**
 * Utility function to safely get images from a response
 * @param {Object} response - Response from getAllItemImages or getAllMediaImages
 * @returns {Array} Array of images or empty array if none found
 */
export const getImagesFromResponse = (response) => {
    return response?.images || [];
};

/**
 * Utility function to safely get videos from a response
 * @param {Object} response - Response from getAllItemVideos or getAllMediaVideos
 * @returns {Array} Array of videos or empty array if none found
 */
export const getVideosFromResponse = (response) => {
    return response?.videos || [];
};

/**
 * Utility function to check if an item has any media
 * @param {Object} response - Response from getMediaByItem
 * @returns {boolean} True if item has media
 */
export const hasMedia = (response) => {
    const media = getMediaFromResponse(response);
    return media && (
        (media.mediaOrder && media.mediaOrder.length > 0) ||
        media.images ||
        media.videos
    );
};
