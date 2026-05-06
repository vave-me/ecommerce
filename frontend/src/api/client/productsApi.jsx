// File: src/api/productApi.js (Recommended new name)

// Assuming this is your configured Axios client instance
import axiosInstance from '../axiosInstance';

// --- Helper Functions (Consider moving to a shared utils file) ---

/**
 * Helper to safely encode URI components.
 * @param {string | number | undefined | null} component - The path component to encode.
 * @returns {string} The encoded component.
 */
const safeEncode = (component) => {
    if (component === null || typeof component === 'undefined') {
        
        return ''; // Or throw? Depends on whether it's truly optional upstream.
    }
    return encodeURIComponent(String(component));
};

/**
 * Helper to create a clean query params object, removing null/undefined/empty strings.
 * @param {object} filters - Input filters object.
 * @returns {object} Cleaned filters object for query params.
 */
function cleanFilters(filters = {}) {
    const cleaned = {};
    for (const key in filters) {
        const value = filters[key];
        // Keep 0, false, but remove null, undefined, empty strings
        if (value !== null && typeof value !== 'undefined' && value !== '') {
            cleaned[key] = value;
        }
    }
    return cleaned;
}

// --- Endpoint Constants ---
const PRODUCTS_ENDPOINT = '/products';
// Note: Variants endpoint is nested under /products/ according to your paths
const VARIANTS_ENDPOINT_BASE = `${PRODUCTS_ENDPOINT}/variants`;

// --- Product Functions ---

/**
 * ADD PRODUCT
 * POST /products
 */
export const addProduct = async (productData) => {
    const endpoint = PRODUCTS_ENDPOINT;
    if (!productData || typeof productData !== 'object') {
        throw new Error('Valid productData object is required for addProduct.');
    }
    try {
        const response = await axiosInstance.post(endpoint, productData);
        return response.data;
    } catch (error) {
        // Error details logged for debugging
        throw error;
    }
};

/**
 * GET PRODUCT BY ID
 */
export const getProductById = async (productId) => {
    const endpoint = `/products/${productId}`;
    try {
        const response = await axiosInstance.get(endpoint);
        return response.data;
    } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle
    }
};

/**
 * REMOVE A PRODUCT
 * DELETE /products/{productId}
 */
export const removeProduct = async (productId, userId = '') => {
    if (!productId) throw new Error('productId is required for removeProduct.');

    const endpoint = `${PRODUCTS_ENDPOINT}/${safeEncode(productId)}`;
    const queryParams = cleanFilters({userId}); // Pass optional userId as query param
    try {
        // For DELETE, query params are usually in the config object
        const response = await axiosInstance.delete(endpoint, {params: queryParams});
        return response.data; // Often empty or a success message
    } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle
    }
};

/**
 * UPDATE A PRODUCT (Generic)
 * PATCH /products/{productId}
 */
export const updateProduct = async (productId, updateData) => {
    if (!productId) throw new Error('productId is required for updateProduct.');
    if (!updateData || typeof updateData !== 'object') {
        throw new Error('Valid updateData object is required for updateProduct.');
    }

    const endpoint = `${PRODUCTS_ENDPOINT}/${safeEncode(productId)}`;
    try {
        const response = await axiosInstance.post(endpoint, updateData);
        return response.data;
    } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle
    }
};

/**
 * REBRAND A PRODUCT
 * PATCH /products/{productId}/rebrand
 */
export const rebrandProduct = async (productId, rebrandData) => {
    if (!productId) throw new Error('productId is required for rebrandProduct.');
    if (!rebrandData || typeof rebrandData !== 'object') {
        throw new Error('Valid rebrandData object is required for rebrandProduct.');
    }

    const endpoint = `${PRODUCTS_ENDPOINT}/${safeEncode(productId)}/rebrand`;
    try {
        const response = await axiosInstance.patch(endpoint, rebrandData);
        return response.data;
    } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle
    }
};

/**
 * UPDATE PRODUCT PRICE
 * PATCH /products/{productId}/price
 */
export const updateProductPrice = async (productId, priceData) => {
    if (!productId) throw new Error('productId is required for updateProductPrice.');
    if (!priceData || typeof priceData !== 'object') { // Validate priceData structure further if needed
        throw new Error('Valid priceData object is required for updateProductPrice.');
    }

    const endpoint = `${PRODUCTS_ENDPOINT}/${safeEncode(productId)}/price`;
    try {
        const response = await axiosInstance.patch(endpoint, priceData);
        return response.data;
    } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle
    }
};

// --- Product State Changes (Archive, Sold, Leased, Pawned) ---

const createProductStateUpdater = (action) => async (productId, actionData = {}) => {
    if (!productId) throw new Error(`productId is required for ${action} action.`);
    if (typeof actionData !== 'object') {
        throw new Error(`Valid actionData object is required for ${action} action.`);
    }

    const endpoint = `${PRODUCTS_ENDPOINT}/${safeEncode(productId)}/${action}`;
    try {
        const response = await axiosInstance.patch(endpoint, actionData);
        return response.data;
    } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle
    }
};

export const archiveProduct = createProductStateUpdater('archive');
export const markProductSold = createProductStateUpdater('sold');
export const markProductLeased = createProductStateUpdater('lease');
export const markProductPawned = createProductStateUpdater('pawn');

/**
 * ADJUST PRODUCT STOCK
 * PATCH /products/{productId}/stock
 */
export const adjustProductStock = async (productId, stockData) => {
    if (!productId) throw new Error('productId is required for adjustProductStock.');
    if (!stockData || typeof stockData !== 'object') { // e.g., { change: 5 } or { absolute: 100 }
        throw new Error('Valid stockData object is required for adjustProductStock.');
    }

    const endpoint = `${PRODUCTS_ENDPOINT}/${safeEncode(productId)}/stock`;
    try {
        const response = await axiosInstance.patch(endpoint, stockData);
        return response.data;
    } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle
    }
};

// --- Variant Functions ---

/**
 * GET VARIANTS for a Product
 * GET /products/{productId}/variants
 */
export const getVariants = async (productId, filters = {}) => {
    if (!productId) throw new Error('productId is required for getVariants.');

    const endpoint = `${PRODUCTS_ENDPOINT}/${safeEncode(productId)}/variants`;
    const queryParams = cleanFilters(filters);
    try {
        const response = await axiosInstance.get(endpoint, {params: queryParams});
        // Add pagination parsing if needed
        return response.data;
    } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle
    }
};

/**
 * ADD VARIANT
 * POST /products/variants
 */
export const addVariant = async (variantData) => {
    // Note: Endpoint is /products/variants, not /products/{prodId}/variants
    const endpoint = VARIANTS_ENDPOINT_BASE;
    if (!variantData || typeof variantData !== 'object') {
        throw new Error('Valid variantData object is required for addVariant.');
    }
    // variantData should likely include the parent productId
    if (!variantData.productId) { // Adjust field name if necessary
        
    }

    try {
        const response = await axiosInstance.post(endpoint, variantData);
        return response.data;
    } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle
    }
};

/**
 * REMOVE VARIANT
 * DELETE /products/variants/{variantId}
 */
export const removeVariant = async (variantId, userId = '') => {
    if (!variantId) throw new Error('variantId is required for removeVariant.');

    const endpoint = `${VARIANTS_ENDPOINT_BASE}/${safeEncode(variantId)}`;
    const queryParams = cleanFilters({userId});
    try {
        const response = await axiosInstance.delete(endpoint, {params: queryParams});
        return response.data;
    } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle
    }
};

// --- Variant State/Property Changes ---

const createVariantUpdater = (action) => async (variantId, actionData = {}) => {
    if (!variantId) throw new Error(`variantId is required for variant ${action} action.`);
    if (typeof actionData !== 'object') {
        throw new Error(`Valid actionData object is required for variant ${action} action.`);
    }

    const endpoint = `${VARIANTS_ENDPOINT_BASE}/${safeEncode(variantId)}/${action}`;
    try {
        const response = await axiosInstance.patch(endpoint, actionData);
        return response.data;
    } catch (error) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', error);
        }
        throw error; // Re-throw for caller to handle
    }
};

export const archiveVariant = createVariantUpdater('archive');
export const decreaseVariantPrice = createVariantUpdater('decreasePrice');
export const increaseVariantPrice = createVariantUpdater('increasePrice');
export const adjustVariantStock = createVariantUpdater('stock');