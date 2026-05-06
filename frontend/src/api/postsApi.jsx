// File: src/api/postApi.js (Recommended new name)
import axiosInstance from './axiosInstance';
const POSTS_ENDPOINT = '/posts'; // Base path relative to axiosInstance.baseURL
/**
 * Helper to safely encode URI components.
 * @param {string | number | undefined | null} component - The path component to encode.
 * @returns {string} The encoded component.
 */
const safeEncode = (component) => {
    if (component === null || typeof component === 'undefined') {
        // Depending on context, might want to throw or return specific value
        return '';
    }
    return encodeURIComponent(String(component));
};
/**
 * Helper to create a clean payload/params object, removing null/undefined/empty strings/empty arrays.
 * Also handles parsing numbers where specified.
 * @param {object} filters - Input filters object.
 * @param {object} numberFields - Keys to parse as numbers { key: 'int'/'float', ... }.
 * @returns {object} Cleaned filters object.
 */
function cleanFilters(filters = {}, numberFields = {}) {
    const cleaned = {};
    for (const key in filters) {
        const value = filters[key];
        // Skip if value is null, undefined, or an empty string
        if (value === null || typeof value === 'undefined' || value === '') continue;
        // Skip if value is an empty array
        if (Array.isArray(value) && value.length === 0) continue;
        // Parse numbers if specified
        if (numberFields[key]) {
            const numType = numberFields[key];
            const parsedValue = numType === 'float' ? parseFloat(value) : parseInt(value, 10);
            // Keep the parsed number only if it's not NaN
            if (!isNaN(parsedValue)) {
                cleaned[key] = parsedValue;
            } else {
            }
        } else {
            // Keep non-empty arrays and other non-empty/null/undefined values
            cleaned[key] = value;
        }
    }
    return cleaned;
}
/**
 * SEARCH POSTS VIA DEDICATED SEARCH ENDPOINT (POST)
 * POST /search/search-posts
 * Note: Uses a different endpoint base '/search/' compared to others.
 */
export const fetchPostsByFilters = async (filters = {}) => {
    // Endpoint seems specific to search, adjust if needed
    const endpoint = '/search/search-posts';
    // Define which fields should be treated as numbers
    const numberFields = {
        lat: 'float',
        lng: 'float',
        radius: 'int',
        page: 'int',
        pageSize: 'int',
        offset: 'int', // Note: page/pageSize and offset/limit are often mutually exclusive
        limit: 'int'
    };
    // Clean the filters, parse numbers, and remove empty values
    const requestBody = cleanFilters(filters, numberFields);
    try {
        const response = await axiosInstance.post(endpoint, requestBody);
        return response.data; // Assuming API returns data directly
    } catch (error) {
        throw error;
    }
};
/**
 * Parse pagination data from response, ensuring numbers.
 * @param {object} data - The response data object.
 * @returns {{totalCount: number, totalPages: number, currentPage: number}}
 */
function parsePagination(data = {}) {
    return {
        totalCount: parseInt(data.totalCount, 10) || 0,
        totalPages: parseInt(data.totalPages, 10) || 0,
        // Often pages are 1-indexed, adjust if API uses 0-index
        currentPage: parseInt(data.currentPage, 10) || 1,
    };
}
/**
 * RETRIEVE ALL POSTS (PAGINATED)
 * GET /posts
 */
export const getPosts = async (filters = {}) => {
    const endpoint = POSTS_ENDPOINT;
    const queryParams = cleanFilters(filters, {page: 'int', pageSize: 'int'});
    try {
        const response = await axiosInstance.get(endpoint, {params: queryParams});
        const data = response.data || {};
        return {
            posts: data.posts || [],
            ...parsePagination(data)
        };
    } catch (error) {
        throw error;
    }
};
/**
 * RETRIEVE POSTS BY USER ID (PAGINATED)
 * GET /posts/users/{userId}
 */
export const getUserPosts = async (userId, filters = {}) => {
    if (!userId) throw new Error('userId is required for getUserPosts.');
    const endpoint = `${POSTS_ENDPOINT}/users/${safeEncode(userId)}`;
    const queryParams = cleanFilters(filters, {page: 'int', pageSize: 'int'});
    try {
        const response = await axiosInstance.get(endpoint, {params: queryParams});
        const data = response.data || {};
        return {
            posts: data.posts || [],
            ...parsePagination(data)
        };
    } catch (error) {
        throw error;
    }
};
/**
 * RETRIEVE POSTS BY CATEGORY ID
 * GET /posts/categories/{categoryId}
 */
export const getPostsByCategory = async (categoryId, filters = {}) => {
    if (!categoryId) throw new Error('categoryId is required for getPostsByCategory.');
    const endpoint = `${POSTS_ENDPOINT}/categories/${safeEncode(categoryId)}`;
    const queryParams = cleanFilters(filters, {page: 'int', pageSize: 'int'});
    try {
        const response = await axiosInstance.get(endpoint, {params: queryParams});
        // Assuming response structure includes pagination and posts array
        const data = response.data || {};
        return {
            posts: data.posts || [],
            ...parsePagination(data)
        };
    } catch (error) {
        throw error;
    }
};
/**
 * RETRIEVE POSTS BY CATEGORY SLUG
 * GET /posts/categories/slug/{slug}
 */
export const getPostsBySlug = async (slug, filters = {}) => {
    if (!slug) throw new Error('slug is required for getPostsBySlug.');
    const endpoint = `${POSTS_ENDPOINT}/categories/slug/${safeEncode(slug)}`;
    const queryParams = cleanFilters(filters, {page: 'int', pageSize: 'int'});
    try {
        const response = await axiosInstance.get(endpoint, {params: queryParams});
        // Assuming response structure includes pagination and posts array
        const data = response.data || {};
        return {
            posts: data.posts || [],
            ...parsePagination(data)
        };
    } catch (error) {
        throw error;
    }
};
/**
 * RETRIEVE A USER'S CATALOG (Specific Post Collection?)
 * GET /posts/catalog
 * Expects userId as a query parameter.
 */
export const getCatalog = async (userId, filters = {}) => {
    // If userId is mandatory for this endpoint, enforce it
    if (!userId) throw new Error('userId is required for getCatalog.');
    const endpoint = `${POSTS_ENDPOINT}/catalog`;
    // Clean other filters and ensure userId is included
    const queryParams = cleanFilters({...filters, userId}, {page: 'int', pageSize: 'int'});
    try {
        const response = await axiosInstance.get(endpoint, {params: queryParams});
        // Assuming response matches postspbGetCatalogResponse => { posts: [...], totalCount, ... }
        const data = response.data || {};
        return {
            posts: data.posts || [],
            ...parsePagination(data)
        };
    } catch (error) {
        throw error;
    }
};
/**
 * GET A SPECIFIC POST BY ID
 * GET /posts/{postId}
 */
export const getPost = async (postId) => {
    if (!postId) throw new Error('postId is required for getPost.');
    const endpoint = `${POSTS_ENDPOINT}/${safeEncode(postId)}`;
    try {
        const response = await axiosInstance.get(endpoint);
        // Assuming response format { post: {...} } or similar
        return response.data; // Or response.data.post if nested
    } catch (error) {
        throw error;
    }
};
/**
 * GET VARIANTS for a Post
 * GET /posts/{postId}/variants
 */
export const getVariants = async (postId, filters = {}) => {
    if (!postId) throw new Error('postId is required for getVariants.');
    const endpoint = `${POSTS_ENDPOINT}/${safeEncode(postId)}/variants`;
    const queryParams = cleanFilters(filters, {page: 'int', pageSize: 'int'}); // Add userId if needed & passed in filters
    try {
        const response = await axiosInstance.get(endpoint, {params: queryParams});
        // response => postspbGetVariantsResponse => { variants: [...], totalCount, ... }
        const data = response.data || {};
        return {
            variants: data.variants || [], // Adjust key based on actual API response
            ...parsePagination(data)
        };
    } catch (error) {
        throw error;
    }
};
/**
 * GET SPECIFIC VARIANT BY ID
 * GET /posts/variants/{variantId}
 */
export const getVariant = async (variantId, userId = '') => {
    if (!variantId) throw new Error('variantId is required for getVariant.');
    // Note: Endpoint seems nested under /posts/, not top-level /variants/
    const endpoint = `${POSTS_ENDPOINT}/variants/${safeEncode(variantId)}`;
    // Clean potential userId filter
    const queryParams = cleanFilters({userId});
    try {
        const response = await axiosInstance.get(endpoint, {params: queryParams});
        // response => postspbGetVariantResponse => { variant: {...} }
        return response.data; // Or response.data.variant if nested
    } catch (error) {
        throw error;
    }
};