// File: src/api/productApi.jsx
import axiosInstance from './axiosInstance';
import axios from 'axios';
// Enhanced cache implementation with size limits and versioning
const CACHE_VERSION = '1.0';
const MAX_CACHE_SIZE = 1000; // Maximum number of cached items
const CACHE_TTL = 5 * 60 * 1000; // 5 minutes cache TTL
class CacheManager {
  constructor() {
    this.cache = new Map();
    this.timestamps = new Map();
    this.size = 0;
  }
  set(key, value) {
    // Check if we need to evict old entries
    if (this.size >= MAX_CACHE_SIZE) {
      this.evictOldest();
    }
    const cacheKey = `${CACHE_VERSION}:${key}`;
    this.cache.set(cacheKey, value);
    this.timestamps.set(cacheKey, Date.now());
    this.size++;
  }
  get(key) {
    const cacheKey = `${CACHE_VERSION}:${key}`;
    const timestamp = this.timestamps.get(cacheKey);
    // Check if entry exists and is not expired
    if (timestamp && Date.now() - timestamp < CACHE_TTL) {
      return this.cache.get(cacheKey);
    }
    // Remove expired entry
    if (timestamp) {
      this.delete(cacheKey);
    }
    return null;
  }
  has(key) {
    const cacheKey = `${CACHE_VERSION}:${key}`;
    const timestamp = this.timestamps.get(cacheKey);
    return timestamp && Date.now() - timestamp < CACHE_TTL;
  }
  delete(key) {
    const cacheKey = `${CACHE_VERSION}:${key}`;
    this.cache.delete(cacheKey);
    this.timestamps.delete(cacheKey);
    this.size--;
  }
  evictOldest() {
    let oldestKey = null;
    let oldestTime = Infinity;
    for (const [key, timestamp] of this.timestamps) {
      if (timestamp < oldestTime) {
        oldestTime = timestamp;
        oldestKey = key;
      }
    }
    if (oldestKey) {
      this.delete(oldestKey);
    }
  }
  clear() {
    this.cache.clear();
    this.timestamps.clear();
    this.size = 0;
  }
}
// Initialize cache manager
const cacheManager = new CacheManager();
// Enhanced request cache with versioning and size limits
const requestCache = cacheManager;
// AbortController registry to manage in-flight requests
const abortControllers = new Map();
/**
 * Flexible search/fetch method matching your backend's GetProductsWithFilters fields.
 * Adjusted with aborting capabilities and caching.
 */
const route = {
    baseURL: process.env.NEXT_PUBLIC_API_BASE_URL || 'http://192.168.178.84:8080/api', // Use env variable
    headers: {
        'Content-Type': 'application/json',
    },
}
// Enhanced cache key generation
const generateCacheKey = (path, params = {}) => {
  const sortedParams = Object.keys(params)
    .sort()
    .reduce((acc, key) => {
      acc[key] = params[key];
      return acc;
    }, {});
  return `${path}:${JSON.stringify(sortedParams)}`;
};
// Function to abort any active request with the same ID
const abortPreviousRequest = (requestId) => {
    if (abortControllers.has(requestId)) {
        const controller = abortControllers.get(requestId);
        controller.abort();
        abortControllers.delete(requestId);
    }
};
export const fetchProductsByFilters = async (filters = {}, requestId = 'products-filter') => {
    // Abort any previous request with the same ID
    abortPreviousRequest(requestId);
    // Create new AbortController for this request
    const abortController = new AbortController();
    abortControllers.set(requestId, abortController);
    try {
        // Destructure the filters with sensible defaults:
        const {
            name = '',
            category = '',
            minPrice = '',
            maxPrice = '',
            condition = '',
            brand = '',
            model = '',
            status = '',
            tags = [],
            manageStock = false,
            minStock = 0,
            maxStock = 9999,
            SKU = '',
            negotiable = false,
            userType = '',
            middlemanService = '',
            hasVariants = false,
            shippingCost = '',
            // Dimensions & weight
            minWeight = '',
            maxWeight = '',
            minHeight = '',
            maxHeight = '',
            minWidth = '',
            maxWidth = '',
            minDepth = '',
            maxDepth = '',
            lat = 0,
            lng = 0,
            radius = 0,
            page = 1,
            pageSize = 20,
            offset = 0,
            limit = 20,
            sortBy = '',
            sortOrder = 'asc',
            skipCache = false, // New parameter to bypass cache
        } = filters;
        // Construct the request body.
        // These field names match your snippet: s.app.GetProductsWithFilters => "Name, Category, MinPrice, ..."
        const requestBody = {
            name: name,
            category: category,
            minPrice: parseFloat(minPrice),
            maxPrice: parseFloat(maxPrice),
            brand: brand,
            condition: condition,
            model: model,
            status: status,
            tags: tags, // array
            manageStock: manageStock,
            minStock: parseInt(minStock, 10),
            maxStock: parseInt(maxStock, 10),
            sku: SKU,
            negotiable: negotiable,
            userType: userType,
            middlemanService: middlemanService,
            hasVariants: hasVariants,
            shippingCost: shippingCost,
            minWeight: parseFloat(minWeight),
            maxWeight: parseFloat(maxWeight),
            minHeight: parseFloat(minHeight),
            maxHeight: parseFloat(maxHeight),
            minWidth: parseFloat(minWidth),
            maxWidth: parseFloat(maxWidth),
            minDepth: parseFloat(minDepth),
            maxDepth: parseFloat(maxDepth), lat: parseFloat(lat),
            lng: parseFloat(lng),
            radius: parseInt(radius, 10),
            page: parseInt(page, 10),
            pageSize: parseInt(pageSize, 10),
            offset: parseInt(offset, 10),
            limit: parseInt(limit, 10),
            sortBy: sortBy,
            sortOrder: sortOrder,
        };
        // Remove empty, null, or undefined fields to keep the request body clean
        Object.keys(requestBody).forEach((key) => {
            const val = requestBody[key];
            // If val is '', null, undefined, or an empty array => remove
            if (
                val === '' ||
                val === null ||
                val === undefined ||
                (Array.isArray(val) && val.length === 0)
            ) {
                delete requestBody[key];
            }
        });
        // POST /search/search-products with the cleaned request body and abort signal
        const response = await axiosInstance.post(
            '/search/search-products', 
            requestBody,
            { signal: abortController.signal }
        );
        return response.data;
    } catch (error) {
        // Don't report aborted requests as errors
        if (axios.isCancel(error)) {
            throw error;
        }
        // Handle 404 errors more gracefully
        if (error.response?.status === 404) {
            return { products: [], totalCount: 0, totalPages: 0, currentPage: 1 };
        }
        throw error;
    } finally {
        // Clean up the abort controller
        abortControllers.delete(requestId);
    }
};
/**
 * RETRIEVE ALL PRODUCTS
 * GET /api/products
 * With caching and abort capabilities
 */
export const getAllProducts = async (filters = {}, requestId = 'all-products') => {
    // Generate cache key for this request
    const cacheKey = generateCacheKey('/products', filters);
    // Check cache first if not explicitly skipped
    if (!filters.skipCache && requestCache.has(cacheKey)) {
        const cachedData = requestCache.get(cacheKey);
        if (cachedData && Date.now() - cachedData.timestamp < CACHE_TTL) {
            return cachedData.data;
        }
        // Cache expired, remove it
        requestCache.delete(cacheKey);
    }
    // Abort any previous request with the same ID
    abortPreviousRequest(requestId);
    // Create new AbortController
    const abortController = new AbortController();
    abortControllers.set(requestId, abortController);
    try {
        const response = await axiosInstance.get(
            '/products', 
            {
                params: filters,
                signal: abortController.signal
            }
        );
        // Cache the result
        requestCache.set(cacheKey, {
            data: response.data,
            timestamp: Date.now()
        });
        return response.data;
    } catch (error) {
        if (axios.isCancel(error)) {
            throw error;
        }
        // Handle 404 errors more gracefully
        if (error.response?.status === 404) {
            return { products: [], totalCount: 0, totalPages: 0, currentPage: 1 };
        }
        throw error;
    } finally {
        abortControllers.delete(requestId);
    }
};
/**
 * RETRIEVE PRODUCTS BY CATEGORY
 * GET /api/products/categories/{categoryId}
 * Accepts query params: page, pageSize, sortBy, sortOrder
 */
export const getProductsByCategory = async (categoryId, filters = {}) => {
    try {
        const response = await axiosInstance.get(
            `/products/categories/${encodeURIComponent(categoryId)}`,
            {params: filters}
        );
        return response.data;
    } catch (error) {
        throw error;
    }
};
/**
 * RETRIEVE PRODUCTS BY SLUG
 * GET /api/products/categories/{categoryId}
 * Accepts query params: page, pageSize, sortBy, sortOrder
 */
export const getProductsBySlug = async (slug, filters = {}) => {
    try {
        const response = await axiosInstance.get(
            `/products/categories/slug/${encodeURIComponent(slug)}`,
            {params: filters}
        );
        return response.data;
    } catch (error) {
        throw error;
    }
};
/**
 * RETRIEVE A USER'S CATALOG
 * GET /api/products/catalog
 * Accepts query params: userId, page, pageSize, sortBy, sortOrder
 */
export const getCatalog = async (filters = {}) => {
    try {
        // If userId is passed as a parameter or in filters, include it
        const response = await axiosInstance.get('/products/catalog', {
            params: filters,
        });
        return response.data;
    } catch (error) {
        throw error;
    }
};
/**
 * GET A SPECIFIC PRODUCT
 * GET /api/products/{id}
 * FALLBACK METHOD
 */
export const getProduct = async (productId) => {
    try {
        const response = await axiosInstance.get(`/products/${encodeURIComponent(productId)}`);
        return response.data;
    } catch (error) {
        throw error;
    }
};
/**
 * GET VARIANTS for a Product
 * GET /api/products/{productId}/variants
 */
export const getVariants = async (productId, filters = {}) => {
    try {
        const response = await axiosInstance.get(
            `/products/${encodeURIComponent(productId)}/variants`,
            {params: filters}
        );
        return response.data;
    } catch (error) {
        throw error;
    }
};
/**
 * ADD VARIANT
 * POST /api/products/variants
 */
export const addVariant = async (variantData) => {
    try {
        const response = await axiosInstance.post('/products/variants', variantData);
        return response.data;
    } catch (error) {
        throw error;
    }
};
/**
 * GET SPECIFIC VARIANT
 * GET /api/products/variants/{variantId}
 */
export const getVariant = async (variantId, userId = '') => {
    try {
        const params = userId ? {userId} : {};
        const response = await axiosInstance.get(
            `/products/variants/${encodeURIComponent(variantId)}`,
            {params}
        );
        return response.data;
    } catch (error) {
        throw error;
    }
};
/**
 * RETRIEVE A USER'S PUBLIC CATALOG
 * GET /api/products/{userId}/catalog
 * Accepts query params: page, pageSize, sortBy, sortOrder
 * This is a public API endpoint that doesn't require authentication
 */
export const getPublicCatalog = async (userId, filters = {}) => {
    try {
        // Use direct axios instead of axiosInstance since this is a public API
        const response = await axiosInstance.get(`/products/${encodeURIComponent(userId)}/catalog`, {
            params: filters,
        });
        return response.data;
    } catch (error) {
        throw error;
    }
};
// Enhanced request function with improved caching
export const getProducts = async (params = {}) => {
  const cacheKey = generateCacheKey('/products', params);
  const cachedData = requestCache.get(cacheKey);
  if (cachedData) {
    return cachedData;
  }
  try {
    const response = await axiosInstance.get('/products', { params });
    requestCache.set(cacheKey, response.data);
    return response.data;
  } catch (error) {
    throw error;
  }
};
// Add cache warming function
export const warmProductCache = async (commonParams = []) => {
  try {
    await Promise.all(
      commonParams.map(params => getProducts(params))
    );
  } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
};
// Add cache invalidation function
export const invalidateProductCache = (params = {}) => {
  const cacheKey = generateCacheKey('/products', params);
  requestCache.delete(cacheKey);
};
