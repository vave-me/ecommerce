// src/api/categoriesApi.js
// Use the centralized Axios instance configured with the environment variable baseURL
import axiosInstance, { ssrAxiosInstance } from './axiosInstance';
import { categoryCache } from '../utils/apiCache';
/* ---------------------------------------------------------------------------
 * M A P P E R   U T I L S
 * ------------------------------------------------------------------------- */
/**
 * Common mapper to unify API response format for a single category object.
 * Maps description to name and ensures subcategories is an array.
 * @param {Object} cat - Category object from API (expects fields like id, description, name, slug, etc.)
 * @returns {Object | null} Normalized category object or null if input is invalid.
 */
function mapCategoryFromApi(cat) {
    if (!cat || typeof cat !== 'object') {
        return null;
    }
    // Ensure description/name are strings before trimming, provide default
    const description = typeof cat.description === 'string' ? cat.description.trim() : '';
    const nameFallback = typeof cat.name === 'string' ? cat.name.trim() : '';
    return {
        ...cat, // Spread original properties first
        // Prioritize description, then name, finally a default placeholder
        name: description || nameFallback || 'Unnamed Category',
        // Safely map subcategories if they exist and are an array
        subcategories: Array.isArray(cat.subcategories)
            ? cat.subcategories.map(mapCategoryFromApi).filter(Boolean) // Recursively map & filter nulls
            : [], // Default to empty array if subcategories are missing or not an array
    };
}
/**
 * Parse numeric values and map categories from common paginated API responses.
 * @param {Object} data - API response data. Expects { categories: [], totalCount: "string", totalPages: "string", currentPage: "string" }
 * @returns {Object} Formatted object with parsed numbers and mapped categories.
 */
function parsePaginatedCategories(data) {
    if (!data || typeof data !== 'object') {
        // Return a default structure matching expected output
        return {categories: [], totalCount: 0, totalPages: 0, currentPage: 0};
    }
    // Support different possible API response formats
    const categories = Array.isArray(data.categories) ? data.categories : [];
    // Parse string or numeric values, with fallbacks for missing data
    const totalCount = parseInt(data.totalCount  || 0, 10);
    const totalPages = parseInt(data.totalPages  || 0, 10);
    const currentPage = parseInt(data.currentPage || 0, 10);
    // Map categories and filter out nulls
    const mappedCategories = categories.map(mapCategoryFromApi).filter(Boolean);
    return {
        categories: mappedCategories,
        totalCount,
        totalPages,
        currentPage,
    };
}
/* ---------------------------------------------------------------------------
 * C A T E G O R Y   E N D P O I N T S
 * ------------------------------------------------------------------------- */
// Base path relative to the baseURL configured in axiosInstance
const CATEGORIES_ENDPOINT = '/categories';
// Cache for category data to avoid redundant API calls
let localCategoryCache = {};
/**
 * Clear the category cache when needed
 */
export function clearCategoryCache() {
    localCategoryCache = {};
}
/**
 * Prefetch multiple category types in parallel
 * @param {Array} categoryTypes - Array of category type strings to fetch
 * @param {string} locale - Language code for localized categories
 * @returns {Promise<Array>} - Results from Promise.allSettled
 */
export async function prefetchCategoryTypes(categoryTypes = [], locale = 'en') {
    if (!categoryTypes.length) return [];
    const promises = categoryTypes.map(type => {
        return fetchMainCategories({ categoryType: type, lang: locale })
          .then(result => {
            // Cache the result
            localCategoryCache[`${type}-${locale}`] = result;
            return result;
          })
          .catch(error => {
            throw error;
          });
    });
    return Promise.allSettled(promises);
}
/**
 * Fetch all categories (paginated) with optional filters.
 * Corresponds to: GET /api/categories (operationId: getCategories)
 * @param {Object} filters - Query parameters (e.g., { lang, categoryType, page, pageSize, sortBy, sortOrder })
 * @returns {Promise<Object>} Formatted paginated response from parsePaginatedCategories.
 * @throws {Error} Re-throws Axios error if API call fails.
 */
export async function getCategories(filters = {}) {
    const endpoint = CATEGORIES_ENDPOINT; // Endpoint: /api/categories
    try {
        const {data} = await axiosInstance.get(endpoint, {params: filters});
        return parsePaginatedCategories(data);
    } catch (error) {
        throw error;
    }
}
/**
 * Fetch a specific category by ID.
 * Corresponds to: GET /api/categories/{id} (operationId: getCategory)
 * @param {string} id - The ID of the category to fetch.
 * @param {Object} [options] - Optional parameters like { userId }.
 * @returns {Promise<Object | null>} Mapped category object or null.
 * @throws {Error} If id is missing or API call fails.
 */
export async function fetchCategory(id, {userId} = {}) {
    if (!id) {
        const missingIdError = new Error('Category ID is required for fetchCategory.');
        throw missingIdError;
    }
    const safeId = encodeURIComponent(id);
    const endpoint = `${CATEGORIES_ENDPOINT}/${safeId}`; // Endpoint: /api/categories/{id}
    const params = userId ? {userId} : {};
    try {
        const {data} = await axiosInstance.get(endpoint, {params});
        // Adapt based on single item response structure (e.g., { category: {...} } or just {...})
        const cat = data?.category || data;
        return cat ? mapCategoryFromApi(cat) : null;
    } catch (error) {
        throw error;
    }
}
/**
 * Fetch main (top-level) categories, supports filtering by categoryType.
 * Corresponds to: GET /api/categories/main (operationId: getMainCategories)
 * @param {Object} filters - Query parameters (e.g., { categoryType, lang, page, pageSize, sortBy, sortOrder })
 * @returns {Promise<Object>} Formatted paginated response from parsePaginatedCategories.
 * @throws {Error} Re-throws Axios error if API call fails.
 */
export async function fetchMainCategories(filters = {}) {
    // *** FIXED Endpoint: Use /main as it accepts categoryType filter according to Swagger ***
    const endpoint = `${CATEGORIES_ENDPOINT}/main`; // Endpoint: /api/categories/main
    // Use SSR-optimized instance on server side
    const axiosInstanceToUse = typeof window === 'undefined' ? ssrAxiosInstance : axiosInstance;
    
    // Use cache for client-side requests
    if (typeof window !== 'undefined') {
        try {
            const cached = await categoryCache.getOrFetch(
                endpoint,
                filters,
                async () => {
                    const {data} = await axiosInstanceToUse.get(endpoint, {params: filters});
                    return data;
                }
            );
            return parsePaginatedCategories(cached);
        } catch (error) {
            throw error;
        }
    }
    
    // Direct request for SSR
    try {
        const {data} = await axiosInstanceToUse.get(endpoint, {params: filters});
        return parsePaginatedCategories(data);
    } catch (error) {
        throw error;
    }
}
/**
 * SSR-safe version of fetchMainCategories with aggressive timeout and fallback
 * @param {Object} filters - Query parameters  
 * @param {number} timeoutMs - Custom timeout in milliseconds (default: 3000)
 * @returns {Promise<Object>} Categories data or fallback structure
 */
export async function fetchMainCategoriesSSR(filters = {}, timeoutMs = 3000) {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), timeoutMs);
    try {
        const endpoint = `${CATEGORIES_ENDPOINT}/main`;
        const {data} = await ssrAxiosInstance.get(endpoint, {
            params: filters,
            signal: controller.signal
        });
        clearTimeout(timeoutId);
        return parsePaginatedCategories(data);
    } catch (error) {
        clearTimeout(timeoutId);
        // Return fallback data structure instead of throwing
        if (process.env.NODE_ENV === 'production') {
        }
        return {
            categories: [],
            totalCount: 0,
            totalPages: 0,
            currentPage: 0
        };
    }
}
/**
 * Alias for fetchMainCategories to maintain backward compatibility
 * @param {Object} filters - Query parameters
 * @returns {Promise<Object>} Paginated category data
 */
export async function getMainCategories(filters = {}) {
    return fetchMainCategories(filters);
}
/**
 * Fetch subcategories for a specific parent category ID.
 * Corresponds to: GET /api/categories/subcategory (operationId: getSubCategories)
 * @param {string} parentId - The ID of the parent category.
 * @param {Object} [filters] - Optional query parameters (e.g., { lang, page, pageSize, sortBy, sortOrder })
 * @returns {Promise<Array<Object>>} Array of mapped subcategory objects.
 * @throws {Error} If parentId is missing or API call fails.
 */
export async function fetchSubCategories(parentId, filters = {}) {
    if (!parentId) {
        const missingParentIdError = new Error('parentId is required to fetch subcategories.');
        throw missingParentIdError;
    }
    const endpoint = `${CATEGORIES_ENDPOINT}/subcategory`; // Endpoint: /api/categories/subcategory
    // Combine parentId with other filters for the query parameters
    const params = {parentCategoryId: parentId, ...filters};
    try {
        const {data} = await axiosInstance.get(endpoint, {params});
        // Expecting { categories: [...] } structure in response data
        const categories = Array.isArray(data?.categories) ? data.categories : [];
        return categories.map(mapCategoryFromApi).filter(Boolean); // Map and filter nulls
    } catch (error) {
        throw error;
    }
}