// File: src/api/productApi.jsx
import axiosInstance from './axiosInstance';
// Base path for all deal-related endpoints
const DEALS_API_BASE_URL = '/deals'; // Assuming '/deals' is relative to axiosInstance.baseURL
/**
 * Fetches all deals, optionally applying filters.
 * GET /deals
 */
export const getAllDeals = async (filters = {}) => {
    const endpoint = DEALS_API_BASE_URL;
    try {
        //  // Added debug log
        const response = await axiosInstance.get(endpoint, {params: filters});
        return response.data;
    } catch (error) {
        // Re-throw the error so UI layer or data fetching hooks (like React Query/SWR) can handle it
        throw error;
    }
};
/**
 * Fetches a specific deal by ID.
 * GET /deals/{id}
 */
export const getDeal = async (id) => {
    if (!id) throw new Error('Deal ID is required for getDeal.');
    const safeDealId = encodeURIComponent(id);
    const endpoint = `${DEALS_API_BASE_URL}/${safeDealId}`;
    try {
        const response = await axiosInstance.get(endpoint);
        return response.data;
    } catch (error) {
        throw error;
    }
};
/**
 * Fetches variants for a specific product ID.
 * GET /deals/{productId}/variants
 */
export const getVariants = async (productId, filters = {}) => {
    // Use encodeURIComponent to safely handle special characters in path parameters
    const safeProductId = encodeURIComponent(productId);
    const endpoint = `${DEALS_API_BASE_URL}/${safeProductId}/variants`; // Corrected template literal syntax
    try {
        // const response = await axiosInstance.get(endpoint, {params: filters});
        return response.data;
    } catch (error) {
        throw error;
    }
};
/**
 * RETRIEVE PRODUCTS BY CATEGORY ID
 * GET /deals/categories/{categoryId}
 */
export const getDealsByCategory = async (categoryId, filters = {}) => {
    const safeCategoryId = encodeURIComponent(categoryId);
    const endpoint = `${DEALS_API_BASE_URL}/categories/${safeCategoryId}`; // Corrected template literal syntax
    try {
        // const response = await axiosInstance.get(endpoint, {params: filters});
        return response.data;
    } catch (error) {
        throw error;
    }
};
/**
 * RETRIEVE PRODUCTS BY CATEGORY SLUG
 * GET /deals/categories/slug/{categorySlug}
 */
export const getDealsBySlug = async (slug, filters = {}) => {
    const safeSlug = encodeURIComponent(slug);
    // Ensure the API endpoint structure matches exactly: /deals/categories/slug/{slug}
    const endpoint = `${DEALS_API_BASE_URL}/categories/slug/${safeSlug}`; // Corrected template literal syntax
    try {
        // const response = await axiosInstance.get(endpoint, {params: filters});
        return response.data;
    } catch (error) {
        throw error;
    }
};
/**
 * RETRIEVE A USER'S CATALOG (Assuming it's related to deals)
 * GET /deals/catalog
 */
export const getDealCatalog = async (filters = {}) => {
    // Assuming the endpoint is literally '/deals/catalog' relative to the base URL
    const endpoint = `${DEALS_API_BASE_URL}/catalog`; // Corrected template literal syntax
    try {
        //const response = await axiosInstance.get(endpoint, {params: filters});
        return response.data;
    } catch (error) {
        throw error;
    }
};
// Consider renaming file to dealsApi.js or deals.api.js if it only contains deal-related functions.
// Also, `.jsx` extension isn't necessary if the file contains no JSX syntax.